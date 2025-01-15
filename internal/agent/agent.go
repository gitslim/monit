// Package agent содержит логику создания и работы агента сбора метрик.
package agent

import (
	"context"

	"github.com/gitslim/monit/internal/agent/collector"
	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/sender"
	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/logging"
)

// Start запускает агент сбора метрик.
func Start(ctx context.Context, cfg *conf.Config, log *logging.Logger) {
	log.Info("Monit agent started")

	// Создание пула worker'ов.
	wp := worker.NewWorkerPool(cfg)

	// Запуск worker'ов отсылки метрик.
	wp.Start(func() {
		sender.RunSendMetricsWorker(ctx, log, wp)
	})

	// Добавление worker'ов сбора метрик.
	wp.AddWorker(func() {
		collector.CollectRuntimeMetrics(ctx, log, wp)
	})
	wp.AddWorker(func() {
		collector.CollectSystemMetrics(ctx, log, wp)
	})

	// Ожидание завершения пула worker'ов.
	wp.WaitClose()
}
