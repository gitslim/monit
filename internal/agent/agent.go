package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/retry"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// gatherRuntimeMetrics - Функция для сбора метрик из пакета runtime
func gatherRuntimeMetrics(ctx context.Context, cfg *conf.Config, log *logging.Logger, metrics chan<- entities.MetricDTO) {
	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval * uint64(time.Second)))
	var memStats runtime.MemStats

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			runtime.ReadMemStats(&memStats)

			gauges := map[string]float64{
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
			gauges["RandomValue"] = rand.Float64()

			for k, v := range gauges {
				gauge, err := entities.NewMetricDTO(k, "gauge", v)
				if err != nil {
					log.Errorf("Failed to create gauge DTO: %k", k)
				} else {
					// log.Debugf("Gauge: %+v", *gauge)
					metrics <- *gauge
				}
			}
			counter, err := entities.NewMetricDTO("PollCount", "counter", int64(1))
			if err != nil {
				log.Error("Failed to create counter DTO: PollCount")
			}
			// log.Debugf("PollCount: %+v", *counter)
			metrics <- *counter
		}
	}
}

func gatherSystemMetrics(ctx context.Context, cfg *conf.Config, log *logging.Logger, metrics chan<- entities.MetricDTO) {
	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval * uint64(time.Second)))
	var metric *entities.MetricDTO

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			vMem, err := mem.VirtualMemory()
			if err != nil {
				log.Errorf("failed to get memory info: %v", err)
			} else {
				metric, err = entities.NewMetricDTO("TotalMemory", "gauge", float64(vMem.Total))
				if err != nil {
					log.Error("Failed to create gauge DTO: TotalMemory")
				} else {
					metrics <- *metric
				}

				metric, err = entities.NewMetricDTO("FreeMemory", "gauge", float64(vMem.Free))
				if err != nil {
					log.Error("Failed to create gauge DTO: FreeMemory")
				} else {
					metrics <- *metric
				}
			}

			cpuPercents, err := cpu.Percent(0, true)
			if err != nil {
				log.Errorf("failed to get CPU info: %v", err)
			} else {
				for i, cpuPercent := range cpuPercents {
					metricName := fmt.Sprintf("CPUutilization%d", i+1)
					metric, err = entities.NewMetricDTO(metricName, "gauge", cpuPercent)
					if err != nil {
						log.Errorf("Failed to create gauge DTO: %v", metricName)
					} else {
						metrics <- *metric
					}

				}
			}
		}
	}
}

// compressGzip - Сжатие данных методом gzip
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

// makeHash - расчет хеша sha256
func makeHash(rc io.ReadCloser, key string) (string, error) {
	var buf []byte
	_, err := rc.Read(buf)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// sendJSON - Отправка метрик в формате JSON батчем или по одной
func sendJSON(cfg *conf.Config, client *http.Client, url string, jsonData []byte) error {
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

	if cfg.Key != "" {
		body, err := req.GetBody()
		if err != nil {
			return err
		}
		hash, err := makeHash(body, cfg.Key)
		if err != nil {
			return err
		}
		req.Header.Set(httpconst.HeaderHashSHA256, hash)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// if batch {
	// 	var v []entities.MetricDTO
	// 	err = json.NewDecoder(resp.Body).Decode(&v)
	// } else {
	// 	var v entities.MetricDTO
	// 	err = json.NewDecoder(resp.Body).Decode(&v)

	// }

	// if err != nil {
	// 	return fmt.Errorf("json decode error")
	// }

	return nil
}

// sendMetrics - Функция для отправки всех метрик на сервер
func sendMetrics(cfg *conf.Config, client *http.Client, metrics []*entities.MetricDTO, batch bool) error {
	serverURL := fmt.Sprintf("http://%s", cfg.Addr)

	return retry.Retry(func() error {
		var url string
		var jsonData []byte
		var err error

		if batch {
			// Отправляем батч метрик
			url = fmt.Sprintf("%s/updates/", serverURL)
			jsonData, err = json.Marshal(metrics)
			if err != nil {
				return err
			}
			return sendJSON(cfg, client, url, jsonData)
		} else {
			// Отправляем метрики по одной
			url = fmt.Sprintf("%s/update/", serverURL)

			for _, metric := range metrics {
				jsonData, err = json.Marshal(&metric)
				if err != nil {
					return err
				}
				err := sendJSON(cfg, client, url, jsonData)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}, 3)
}

// sendMetricsWorker - Воркер для отправки метрик
func sendMetricsWorker(ctx context.Context, cfg *conf.Config, log *logging.Logger, client *http.Client, metrics <-chan entities.MetricDTO, wg *sync.WaitGroup) {
	defer wg.Done()

	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval * uint64(time.Second)))
	batch := []*entities.MetricDTO{}

	for {
		select {
		case metric := <-metrics:
			batch = append(batch, &metric)
		case <-ctx.Done():
			return
		case <-reportTicker.C:
			err := sendMetrics(cfg, client, batch, false)
			if err != nil {
				log.Errorf("Send metrics failed: %v\n", err)
			}
		}
	}
}

// Start - Запуск агента
func Start(ctx context.Context, cfg *conf.Config, log *logging.Logger) {
	// pollInterval := time.Duration(cfg.PollInterval * uint64(time.Second))
	// reportInterval := time.Duration(cfg.ReportInterval * uint64(time.Second))

	client := &http.Client{}
	log.Info("Agent started")

	metrics := make(chan entities.MetricDTO, cfg.RateLimit)
	var wg sync.WaitGroup

	// Запуск worker pool
	for i := 0; i < int(cfg.RateLimit); i++ {
		wg.Add(1)
		go sendMetricsWorker(ctx, cfg, log, client, metrics, &wg)
	}

	// Горутина для сбора рантайм метрик
	go gatherRuntimeMetrics(ctx, cfg, log, metrics)

	// Горутина для сбора системных метрик
	go gatherSystemMetrics(ctx, cfg, log, metrics)

	// Ожидание завершения
	wg.Wait()

	// lastReportTime := time.Now() // Время последней отправки метрик
	// var pollCount int64          // Счётчик обновлений PollCount

	// for {
	// 	var metrics []*entities.MetricDTO
	// 	newMetrics := gatherRuntimeMetrics()
	// 	for key, value := range newMetrics {
	// 		dto, err := entities.NewMetricDTO(key, "gauge", value)
	// 		if err != nil {
	// 			log.Errorf("Failed to create gauge DTO: %s\n", key)
	// 		}
	// 		metrics = append(metrics, dto)
	// 	}
	// 	pollCount++ // Увеличиваем счетчик PollCount

	// 	// Если прошло 10 секунд с момента последней отправки, отправляем метрики
	// 	if time.Since(lastReportTime) >= reportInterval {
	// 		dto, err := entities.NewMetricDTO("PollCount", "counter", pollCount)
	// 		if err != nil {
	// 			log.Errorf("Failed to create counter DTO: PollCount")
	// 		}
	// 		metrics = append(metrics, dto)
	// 		err = sendMetrics(cfg, client, metrics, true)
	// 		if err != nil {
	// 			fmt.Printf("Send metrics failed: %v\n", err)
	// 		}
	// 		lastReportTime = time.Now() // Обновляем время последней отправки
	// 	}

	// 	// Собираем метрики из runtime каждые 2 секунды (pollInterval)
	// 	time.Sleep(pollInterval)
	//	}
}
