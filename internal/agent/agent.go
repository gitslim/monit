// Package agent содержит логику создания и работы агента сбора метрик.
package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gitslim/monit/internal/agent/collector"
	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/sender"
	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/logging"
)

// Start запускает агент сбора метрик.
func Start(cfg *conf.Config, log *logging.Logger) {
	log.Info("Monit agent started")

	// Контекст для graceful shutdown с таймаутом.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для сигналов ОС.
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Создание пула worker'ов.
	wp := worker.NewWorkerPool(cfg)

	// Запуск worker'ов отсылки метрик.
	wp.Start(ctx, func(ctx context.Context) {
		sender.RunSendMetricsWorker(ctx, log, wp)
	})

	// Добавление worker'ов сбора метрик.
	wp.AddWorker(ctx, func(ctx context.Context) {
		collector.CollectRuntimeMetrics(ctx, log, wp)
	})
	wp.AddWorker(ctx, func(ctx context.Context) {
		collector.CollectSystemMetrics(ctx, log, wp)
	})

	// Ожидание сигнала завершения.
	go func() {
		quit := <-quitChan
		log.Infof("Received signal: %v, shutting down...", quit)

		cancel() // Отмена контекста
	}()

	wp.WaitClose() // Ожидаем завершения worker'ов
	log.Info("Monit agent stopped.")
}
