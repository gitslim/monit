package worker

import (
	"context"
	"time"

	"github.com/gitslim/monit/internal/agent/sender"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/logging"
)

// RunSendMetricsWorker запуск воркера отправки метрик.
func RunSendMetricsWorker(ctx context.Context, log *logging.Logger, wp *WorkerPool) {
	// Таймер для периодической отправки метрик.
	reportTicker := time.NewTicker(time.Duration(wp.Cfg.ReportInterval * uint64(time.Second)))
	defer reportTicker.Stop()

	// Создаем пустой батч метрик.
	batch := []*entities.MetricDTO{}

	for {
		select {
		case metric := <-wp.Metrics:
			// Добавляем метрику в батч.
			batch = append(batch, &metric)
		case <-ctx.Done():
			return
		case <-reportTicker.C:
			// Отправляем батч метрик.
			if wp.Cfg.GRPC {
				err := sender.GRPCSendMetrics(ctx, wp.Cfg, *wp.GRPCClient, batch, true)
				if err != nil {
					log.Errorf("GRPC send metrics failed: %v\n", err)
				}
			} else {
				err := sender.HTTPSendMetrics(ctx, wp.Cfg, wp.HTTPClient, batch, true)
				if err != nil {
					log.Errorf("HTTP send metrics failed: %v\n", err)
				}
			}
		}
	}
}
