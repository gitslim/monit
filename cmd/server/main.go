// Команда server запускает сервер метрик.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "net/http/pprof"

	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/server/grpcengine"
	"github.com/gitslim/monit/internal/server/httpengine"
	"github.com/gitslim/monit/internal/services"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", getOrDefault(buildVersion))
	fmt.Printf("Build date: %s\n", getOrDefault(buildDate))
	fmt.Printf("Build commit: %s\n", getOrDefault(buildCommit))
}

func getOrDefault(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func main() {
	// Вывод информации о билде.
	printBuildInfo()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// Инициализация логгера.
	log, err := logging.NewLogger()
	if err != nil {
		// Логгер еще недоступен поэтому panic...
		panic(fmt.Sprintf("Failed to initialize logger: %v\n", err))
	}

	// Парсинг конфига.
	cfg, err := conf.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	log.Infof("Monit server config: %+v", cfg)

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
			stop()
		}()
	}

	// Инициализация сервиса метрик.
	svc, err := services.NewMetricService(metricConf)
	if err != nil {
		log.Fatalf("Metric service initialization failed: %v", err)
	}

	var wg sync.WaitGroup

	// Запуск http сервера.
	wg.Add(1)
	if cfg.GRPC {
		go grpcengine.StartGRPCServer(ctx, cfg, log, &wg, svc)
	} else {
		go httpengine.StartHTTPServer(ctx, cfg, log, &wg, svc)
	}

	// Ожидание завершения.
	wg.Wait()

	log.Info("Monit stopped")
}
