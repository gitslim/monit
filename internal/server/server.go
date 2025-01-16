// Package server содержит логику создания и работы сервера метрик.
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/server/engine"
	"github.com/gitslim/monit/internal/services"
)

// Start запускает сервер.
func Start(ctx context.Context, cfg *conf.Config, log *logging.Logger, metricService *services.MetricService) {
	// Создание gin engine.
	r, err := engine.CreateGinEngine(cfg, log, gin.ReleaseMode, metricService)
	if err != nil {
		panic("Creating gin engine failed")
	}

	// Создаем сервер.
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	// Запуск сервера в горутине.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	log.Infof("Server is running on %v\n", cfg.Addr)

	// Gracefull shutdown.
	// Таймаут ожидания завершения работы сервера.
	gracefulTimeout := 5 * time.Second

	// Создаем контекст с тайм-аутом для завершения работы сервера.
	ctx, cancel := context.WithTimeout(ctx, gracefulTimeout)
	defer cancel()

	// Канал для получения сигналов.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Инициируем graceful shutdown.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:", err)
	}

	log.Info("Server exited")
}
