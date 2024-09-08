package main

import (
	"log"

	"github.com/gitslim/monit/internal/handlers"
	"github.com/gitslim/monit/internal/server"
	"github.com/gitslim/monit/internal/storage"
)

func main() {
	// Инициализация хранилища
	memStorage := storage.NewMemStorage()

	// Инициализация обработчика
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	// Инициализация и запуск сервера
	srv := server.InitServer(":8080", metricsHandler)

	// Запуск сервера
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
