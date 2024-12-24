package agent

import (
	"context"

	"github.com/gitslim/monit/internal/agent/collector"
	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/sender"
	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/logging"
)

// Start - Запуск агента сбора метрик
func Start(ctx context.Context, cfg *conf.Config, log *logging.Logger) {
	log.Info("Monit agent started")

	// Создание пула worker'ов
	wp := worker.NewWorkerPool(cfg)

	// Запуск worker'ов отсылки метрик
	wp.Start(func() {
		sender.SendMetricsWorker(ctx, log, wp)
	})

	// Сбор рантайм метрик
	go collector.CollectRuntimeMetrics(ctx, log, wp)

	// Сбор системных метрик
	go collector.CollectSystemMetrics(ctx, log, wp)

	// Ожидание завершения пула worker'ов
	wp.Wait()
}
