package collector

import (
	"context"
	"testing"
	"time"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestCollectMetrics(t *testing.T) {
	cfg := &conf.Config{
		PollInterval: 1,
		RateLimit:    5,
	}

	metricsCh := make(chan entities.MetricDTO, 100)
	wp := worker.NewWorkerPool(cfg)
	wp.Metrics = metricsCh

	log, err := logging.NewLogger()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Добавление worker'ов сбора метрик
	wp.AddWorker(func() {
		CollectRuntimeMetrics(ctx, log, wp)
	})
	wp.AddWorker(func() {
		CollectSystemMetrics(ctx, log, wp)
	})

	// Завершаем сбор метрик
	time.Sleep(1 * time.Second)
	cancel()

	// Ждем завершения пула worker'ов
	wp.WaitClose()

	collectedMetrics := make(map[string]bool)
	for metric := range metricsCh {
		collectedMetrics[metric.ID] = true
	}

	expectedMetrics := []string{
		"Alloc", "BuckHashSys", "HeapAlloc", "RandomValue", "PollCount",
		"TotalMemory", "FreeMemory", "CPUutilization1",
	}

	for _, metricName := range expectedMetrics {
		assert.Contains(t, collectedMetrics, metricName, "Metric %s should be collected", metricName)
	}
}
