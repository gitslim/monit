package sender

import (
	"context"
	"time"

	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/logging"
)

// RunSendMetricsWorker запуск воркера отправки метрик
func RunSendMetricsWorker(ctx context.Context, log *logging.Logger, wp *worker.WorkerPool) {
	// Таймер для периодической отправки метрик
	reportTicker := time.NewTicker(time.Duration(wp.Cfg.ReportInterval * uint64(time.Second)))
	defer reportTicker.Stop()

	// Создаем пустой батч метрик
	batch := []*entities.MetricDTO{}

	for {
		select {
		case metric := <-wp.Metrics:
			// Добавляем метрику в батч
			batch = append(batch, &metric)
		case <-ctx.Done():
			return
		case <-reportTicker.C:
			// Отправляем батч метрик
			err := SendMetrics(ctx, wp.Cfg, wp.Client, batch, false)
			if err != nil {
				log.Errorf("Send metrics failed: %v\n", err)
			}
		}
	}
}
