package sender

import (
	"context"
	"time"

	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/logging"
)

// SendMetricsWorker - Воркер для отправки метрик
func SendMetricsWorker(ctx context.Context, log *logging.Logger, wp *worker.WorkerPool) {
	defer wp.WG.Done()

	reportTicker := time.NewTicker(time.Duration(wp.Cfg.ReportInterval * uint64(time.Second)))
	batch := []*entities.MetricDTO{}

	for {
		select {
		case metric := <-wp.Metrics:
			batch = append(batch, &metric)
		case <-ctx.Done():
			return
		case <-reportTicker.C:
			err := SendMetrics(ctx, wp.Cfg, wp.Client, batch, false)
			if err != nil {
				log.Errorf("Send metrics failed: %v\n", err)
			}
		}
	}
}
