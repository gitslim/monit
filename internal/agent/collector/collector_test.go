package collector

import (
	"context"
	"sync"
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
	wp := &worker.WorkerPool{
		Metrics: metricsCh,
		WG:      &sync.WaitGroup{},
		Cfg:     cfg,
	}

	log, err := logging.NewLogger()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go CollectRuntimeMetrics(ctx, log, wp)
	go CollectSystemMetrics(ctx, log, wp)

	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
	close(metricsCh)

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
