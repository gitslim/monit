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
	// Создаем конфигурацию
	cfg := &conf.Config{
		PollInterval: 1,
		RateLimit:    5,
	}

	// Канал сбора метрик
	metricsCh := make(chan entities.MetricDTO, 100)

	// Инициализация пула worker'ов
	wp := worker.NewWorkerPool(cfg)
	wp.Metrics = metricsCh

	// Логгер
	log, err := logging.NewLogger()
	assert.NoError(t, err)

	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Добавление worker'ов сбора метрик
	wp.AddWorker(func() {
		CollectRuntimeMetrics(ctx, log, wp)
	})
	wp.AddWorker(func() {
		CollectSystemMetrics(ctx, log, wp)
	})

	// Даем время на сбор метрик
	time.Sleep(1 * time.Second)

	// Завершаем пул worker'ов
	cancel()
	wp.WaitClose()

	// Забираем метрики из канала
	collected := make(map[string]bool)
	for metric := range metricsCh {
		collected[metric.ID] = true
	}

	// Список ожидаемых метрик
	expected := []string{
		"Alloc", "BuckHashSys", "HeapAlloc", "RandomValue", "PollCount",
		"TotalMemory", "FreeMemory", "CPUutilization1",
	}

	// Проверяем что все метрики собрались
	for _, metricName := range expected {
		assert.Contains(t, collected, metricName, "Metric %s should be collected", metricName)
	}
}
