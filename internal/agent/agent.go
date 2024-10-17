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
	"time"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/httpconst"
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

	gzWriter, err := gzip.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %v", err)
	}
	defer gzWriter.Close()

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
func sendMetric(client *http.Client, serverURL string, dto *entities.MetricDTO) error {
	url := fmt.Sprintf("%s/update/", serverURL)

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

	req.Header.Set(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
	req.Header.Set(httpconst.HeaderContentEncoding, httpconst.ContentEncodingGzip)

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

	return nil
}

// Функция для отправки всех метрик на сервер
func sendMetrics(client *http.Client, serverURL string, metrics []*entities.MetricDTO) {
	// Отправляем метрики по одной
	for _, metric := range metrics {
		err := sendMetric(client, serverURL, metric)
		if err != nil {
			log.Printf("Error sending metric: %+v\n", metric)
		}
	}
}

func Start(cfg *conf.Config) {
	serverURL := fmt.Sprintf("http://%s", cfg.Addr)
	pollInterval := time.Duration(cfg.PollInterval * uint64(time.Second))
	reportInterval := time.Duration(cfg.ReportInterval * uint64(time.Second))

	lastReportTime := time.Now() // Время последней отправки метрик
	var pollCount int64          // Счётчик обновлений PollCount

	client := &http.Client{}

	for {
		var metrics []*entities.MetricDTO
		newMetrics := gatherRuntimeMetrics()
		for key, value := range newMetrics {
			dto, err := entities.NewMetricDTO(key, "gauge", value)
			if err != nil {
				log.Printf("Failed to create gauge DTO: %s\n", key)
			}
			metrics = append(metrics, dto)
		}
		pollCount++ // Увеличиваем счетчик PollCount

		// Если прошло 10 секунд с момента последней отправки, отправляем метрики
		if time.Since(lastReportTime) >= reportInterval {
			dto, err := entities.NewMetricDTO("PollCount", "counter", pollCount)
			if err != nil {
				log.Println("Failed to create counter DTO: PollCount")
			}
			metrics = append(metrics, dto)
			sendMetrics(client, serverURL, metrics)
			lastReportTime = time.Now() // Обновляем время последней отправки
		}

		// Собираем метрики из runtime каждые 2 секунды (pollInterval)
		time.Sleep(pollInterval)
	}
}
