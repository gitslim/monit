package main

import (
	"log"

	"github.com/gitslim/monit/internal/handlers"
	"github.com/gitslim/monit/internal/server"
	"github.com/gitslim/monit/internal/storage"
)

func main() {
	// Парсинг конфига
	cfg, err := server.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	// Инициализация хранилища
	memStorage := storage.NewMemStorage()

	// Инициализация обработчика
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	// Инициализация сервера
	srv := server.New(cfg.Addr, metricsHandler)

	// Запуск сервера
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
