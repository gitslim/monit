package collector

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/logging"
	"github.com/stretchr/testify/assert"
)

// TesCollectMetrics тестирует сбор метрик.
func TestCollectMetrics(t *testing.T) {
	// Создаем конфигурацию.
	cfg := &conf.Config{
		PollInterval: 1,
		RateLimit:    5,
	}

	// Инициализация пула worker'ов.
	wp := worker.NewWorkerPool(cfg)

	// Создаем логгер.
	log, err := logging.NewLogger()
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Добавление worker'ов сбора метрик.
	wp.AddWorker(ctx, func(ctx context.Context) {
		CollectRuntimeMetrics(ctx, log, wp)
	})
	wp.AddWorker(ctx, func(ctx context.Context) {
		CollectSystemMetrics(ctx, log, wp)
	})

	// Даем время на сбор метрик.
	time.Sleep(3 * time.Second)

	// Забираем метрики из канала.
	collected := &sync.Map{}
	go func() {
		for metric := range wp.Metrics {
			collected.Store(metric.ID, true)
		}
	}()

	// Завершаем работу worker'ов
	cancel()

	// Ждём завершения всех worker'ов.
	wp.WaitClose()

	// Список ожидаемых метрик.
	expected := []string{
		"Alloc", "BuckHashSys", "HeapAlloc", "RandomValue", "PollCount",
		"TotalMemory", "FreeMemory", "CPUutilization1",
	}

	// Проверяем что все метрики собрались.
	for _, metricName := range expected {
		_, ok := collected.Load(metricName)
		assert.True(t, ok, "Metric %s should be collected", metricName)
	}
}
