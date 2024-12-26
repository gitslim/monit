package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/handlers"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/services"
)

func Start(ctx context.Context, cfg *conf.Config, log *logging.Logger, metricService *services.MetricService) {
	// Gin engine handler
	r, err := handlers.CreateGinEngine(cfg, log, gin.ReleaseMode, "templates/*", metricService)
	if err != nil {
		panic("Creating gin engine failed")
	}

	// сервер
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	// запуск в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	log.Infof("Server is running on %v\n", cfg.Addr)

	// Gracefull shutdown
	// таймаут ожидания завершения
	gracefulTimeout := 5 * time.Second

	// Создаем контекст с тайм-аутом для завершения работы сервера
	ctx, cancel := context.WithTimeout(ctx, gracefulTimeout)
	defer cancel()

	// канал для получения сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Инициируем graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:", err)
	}

	log.Info("Server exited")
}
