package main

import (
	"log"

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
	storage := storage.NewMemStorage()

	// Запуск сервера
	if err := server.Start(cfg.Addr, storage); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
