// Модуль для запуска сервера метрик.
package main

import (
	"context"
	"fmt"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация логгера.
	log, err := logging.NewLogger()
	if err != nil {
		// Логгер еще недоступен поэтому fmt...
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Парсинг конфига.
	cfg, err := conf.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	log.Debugf("Server config: %+v", cfg)

	// Инициализация хранилища.
	var metricConf services.MetricServiceConf
	if cfg.DatabaseDSN != "" {
		log.Debug("Using postgres storage")
		metricConf, err = services.WithPGStorage(ctx, log, cfg)
		if err != nil {
			log.Fatalf("Postgres storage configuration failed: %v", err)
		}
	} else {
		log.Debug("Using memory storage")
		errCh := make(chan error)
		metricConf, err = services.WithMemStorage(ctx, log, cfg, errCh)
		if err != nil {
			log.Fatalf("Memory storage configuration failed: %v", err)
		}

		// Обработка ошибки бэкапа.
		go func() {
			<-errCh
			cancel()
		}()
	}

	// Инициализация сервиса метрик.
	svc, err := services.NewMetricService(metricConf)
	if err != nil {
		log.Fatalf("Metric service initialization failed: %v", err)
	}

	// Запуск pprof сервера.
	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	// Запуск сервера.
	server.Start(ctx, cfg, log, svc)
}
