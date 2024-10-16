package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server"
	"github.com/gitslim/monit/internal/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация логгера
	log, err := logging.NewLogger()
	if err != nil {
		// Логгер еще недоступен поэтому fmt...
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Парсинг конфига
	cfg, err := server.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	log.Debugf("Server config: %+v", cfg)

	// Инициализация хранилища
	backupErrChan := make(chan error)
	metricConf, err := services.WithMemStorage(ctx, log, cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore, backupErrChan)
	if err != nil {
		log.Fatalf("Metric service configuration failed: %v", err)
	}

	// обработка ошибки бэкапа
	go func() {
		<-backupErrChan
		cancel()
	}()

	// Инициализация сервиса метрик
	metricService, err := services.NewMetricService(metricConf)
	if err != nil {
		log.Fatalf("Metric service initialization failed: %v", err)
	}

	// Запуск сервера
	server.Start(ctx, cfg.Addr, log, metricService)
}
