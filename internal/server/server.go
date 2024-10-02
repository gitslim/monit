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
	"github.com/gitslim/monit/internal/middleware"
	"github.com/gitslim/monit/internal/services"
	"go.uber.org/zap"
)

func Start(addr string, sugar *zap.SugaredLogger, metricService *services.MetricService) {
	// таймаут ожидания завершения
	gracefulTimeout := 5 * time.Second

	// роутер
	r := gin.New()

	// логгирование через middleware
	r.Use(middleware.LoggerMiddleware(sugar))
	gin.SetMode(gin.ReleaseMode)

	r.LoadHTMLGlob("templates/*")

	// создание хендлера
	metricHandler := handlers.NewMetricHandler(metricService)

	// роуты
	r.GET("/", metricHandler.ListMetrics)
	r.POST("/update/", metricHandler.UpdateMetric)
	r.POST("/value/", metricHandler.GetMetric)
	r.GET("/value/:type/:name", metricHandler.GetMetric)
	r.POST("/update/:type/:name/:value", metricHandler.UpdateMetric)

	// сервер
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// запуск в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	sugar.Infof("Server is running on %v\n", addr)

	// канал для получения сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Shutting down server...")

	// Создаем контекст с тайм-аутом для завершения работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	// Инициируем graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server forced to shutdown:", err)
	}

	sugar.Info("Server exited")
}
