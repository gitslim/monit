package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
)

// Функция для сбора метрик из пакета runtime
func gatherRuntimeMetrics() map[string]float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := map[string]float64{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": float64(memStats.GCCPUFraction),
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
	}

	// Добавляем случайное значение RandomValue
	metrics["RandomValue"] = rand.Float64()

	return metrics
}

func compressGzip(data []byte, level int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	gzWriter, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %v", err)
	}

	_, err = gzWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %v", err)
	}

	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}
	return &buf, nil
}

// Функция для отправки метрики на сервер
func sendMetric(client *http.Client, serverURL string, mType, mName, mValue string) error {
	url := fmt.Sprintf("%s/update/", serverURL)

	dto, err := entities.NewMetricDTO(mName, mType, mValue)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(&dto)
	if err != nil {
		return errs.ErrInternal
	}

	buf, err := compressGzip(jsonData, gzip.BestSpeed)
	if err != nil {
		return fmt.Errorf("failed to compress with gzip: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(dto)
	if err != nil {
		return fmt.Errorf("json decode error")
	}

	// var v interface{}
	// if dto.Delta != nil {
	// 	v = *dto.Delta
	// }
	// if dto.Value != nil {
	// 	v = *dto.Value
	// }
	// fmt.Printf("name: %v, type: %v, value: %v\n", dto.ID, dto.MType, v)

	return nil
}

// Функция для отправки всех метрик на сервер
func sendMetrics(client *http.Client, serverURL string, metrics map[string]float64, counter int64) {
	// Отправляем метрики типа gauge
	for name, value := range metrics {
		err := sendMetric(client, serverURL, "gauge", name, strconv.FormatFloat(value, 'f', -1, 64))
		if err != nil {
			log.Printf("Error sending gauge metric %s: %v\n", name, err)
		}
	}

	// Отправляем метрику PollCount типа counter
	err := sendMetric(client, serverURL, "counter", "PollCount", strconv.FormatInt(counter, 10))
	if err != nil {
		log.Printf("Error sending counter metric PollCount: %v\n", err)
	}
}

func Start(cfg *Config) {
	serverURL := fmt.Sprintf("http://%s", cfg.Addr)
	pollInterval := time.Duration(cfg.PollInterval * float64(time.Second))
	reportInterval := time.Duration(cfg.ReportInterval * float64(time.Second))

	metrics := make(map[string]float64)
	lastReportTime := time.Now() // Время последней отправки метрик
	var pollCount int64          // Счётчик обновлений PollCount

	client := &http.Client{}

	for {
		newMetrics := gatherRuntimeMetrics()
		for key, value := range newMetrics {
			metrics[key] = value
		}
		pollCount++ // Увеличиваем счетчик PollCount

		// Если прошло 10 секунд с момента последней отправки, отправляем метрики
		if time.Since(lastReportTime) >= reportInterval {
			sendMetrics(client, serverURL, metrics, pollCount)
			lastReportTime = time.Now() // Обновляем время последней отправки
		}

		// Собираем метрики из runtime каждые 2 секунды (pollInterval)
		time.Sleep(pollInterval)
	}
}
