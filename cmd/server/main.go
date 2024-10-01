package main

import (
	"log"

	"github.com/gitslim/monit/internal/server"
	"github.com/gitslim/monit/internal/services"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логгера
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// Очищаем буферы логгера
	defer logger.Sync()

	sugar := logger.Sugar()

	// Парсинг конфига
	cfg, err := server.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	// Инициализация сервиса метрик
	metricService, err := services.NewMetricService(services.WithMemStorage())
	if err != nil {
		log.Fatalf("Metric service initialization failed: %v", err)
	}

	// Запуск сервера
	server.Start(cfg.Addr, sugar, metricService)
}
