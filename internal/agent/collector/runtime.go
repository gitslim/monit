package collector

import (
	"context"
	"math/rand/v2"
	"runtime"
	"time"

	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/logging"
)

// CollectRuntimeMetrics собирает метрики информации о системе и отправляет их в канал wp.Metrics.
func CollectRuntimeMetrics(ctx context.Context, log *logging.Logger, wp *worker.WorkerPool) {
	// Таймер для периодического сбора метрик.
	pollTicker := time.NewTicker(time.Duration(wp.Cfg.PollInterval * uint64(time.Second)))
	defer pollTicker.Stop()

	var memStats runtime.MemStats

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			// Сбор статистики памяти.
			runtime.ReadMemStats(&memStats)

			// Подготовка gauges.
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

			// Добавляем случайное значение RandomValue.
			gauges["RandomValue"] = rand.Float64()

			// Отправка gauges в канал.
			for k, v := range gauges {
				gauge, err := entities.NewMetricDTO(k, "gauge", v)
				if err != nil {
					log.Errorf("Failed to create gauge DTO: %k", k)
				} else {
					wp.Metrics <- *gauge
				}
			}
			
			// Отправка counter в канал.
			counter, err := entities.NewMetricDTO("PollCount", "counter", int64(1))
			if err != nil {
				log.Error("Failed to create counter DTO: PollCount")
			}
			wp.Metrics <- *counter
		}
	}
}
