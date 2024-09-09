package main

import (
	"fmt"
	"log"

	"github.com/gitslim/monit/internal/config"
	"github.com/gitslim/monit/internal/handlers"
	"github.com/gitslim/monit/internal/server"
	"github.com/gitslim/monit/internal/storage"
)

func main() {
	// Парсинг конфига
	conf, err := config.Parse()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	// Инициализация хранилища
	memStorage := storage.NewMemStorage()

	// Инициализация обработчика
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	// Инициализация сервера
	srv := server.New(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port), metricsHandler)

	// Запуск сервера
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
