package main

import (
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server"
	"github.com/gitslim/monit/internal/services"
)

func main() {
	// Инициализация логгера
	log, err := logging.NewLogger()
	if err != nil {
		panic("Failed init logger")
	}

	// Парсинг конфига
	cfg, err := server.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	log.Debugf("Server config: %+v", cfg)

	// Инициализация хранилища
	metricConf, err := services.WithMemStorage(log, cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore)
	if err != nil {
		log.Fatalf("Metric service configuration failed: %v", err)
	}

	// Инициализация сервиса метрик
	metricService, err := services.NewMetricService(metricConf)
	if err != nil {
		log.Fatalf("Metric service initialization failed: %v", err)
	}

	// Запуск сервера
	server.Start(cfg.Addr, log, metricService)
}
