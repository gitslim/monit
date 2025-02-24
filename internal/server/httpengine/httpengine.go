// Package httpengine предоставляет функции для создания и настройки Gin engine.
package httpengine

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/handlers"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/middleware"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/services"
)

func getTemplateGlob() (string, error) {
	// Получаем путь к файлу, из которого вызывается функция.
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get runtime caller file path")
	}

	dir := filepath.Dir(file)
	// Получаем путь к папке с шаблонами.
	parts := strings.Split(dir, "/")
	dir = strings.Join(parts[:len(parts)-3], "/")
	return fmt.Sprintf("%s/templates/*", dir), nil
}

// CreateGinEngine создает и настраивает Gin engine с использованием конфигурации, логгера, режима Gin и шаблонов HTML.
func CreateGinEngine(cfg *conf.Config, log *logging.Logger, ginMode string, metricService *services.MetricService) (g *gin.Engine, e error) {
	// Создаем gin engine.
	gin.SetMode(ginMode)
	r := gin.New()

	// Обработка паники gin.
	defer func() {
		if rec := recover(); rec != nil {
			e = fmt.Errorf("gin panic: %v", rec) // nolint:errcheck
		}
	}()

	// Миддлвари
	if cfg.TrustedSubnet != "" {
		log.Info("Using trusted subnet middleware")
		r.Use(middleware.TrustedSubnetMiddleware(log, cfg.TrustedSubnet))
	}

	log.Info("Using gzip middleware")
	r.Use(middleware.GzipMiddleware())

	if cfg.CryptoKey != "" {
		log.Info("Using decrypt middleware")
		dmw, err := middleware.DecryptMiddleware(cfg.CryptoKey)
		if err != nil {
			return nil, err
		}
		r.Use(dmw)
	}

	if cfg.Key != "" {
		log.Info("Using signature middleware")
		r.Use(middleware.SignatureMiddleware(log, cfg.Key))
	}

	log.Info("Using logger middleware")
	r.Use(middleware.LoggerMiddleware(log))

	// Загрузка шаблонов HTML.
	t, err := getTemplateGlob()
	if err != nil {
		return nil, err
	}
	r.LoadHTMLGlob(t)

	// Создание хендлера.
	metricHandler := handlers.NewMetricHandler(metricService)

	// Роуты.
	r.GET("/", metricHandler.ListMetrics)
	r.POST("/update/", metricHandler.UpdateMetric)
	r.POST("/updates/", metricHandler.BatchUpdateMetrics)
	r.POST("/value/", metricHandler.GetMetric)
	r.GET("/value/:type/:name", metricHandler.GetMetric)
	r.POST("/update/:type/:name/:value", metricHandler.UpdateMetric)
	r.GET("/ping", metricHandler.PingStorage)

	return r, err
}

// StartHTTPServer запускает http сервер.
func StartHTTPServer(ctx context.Context, cfg *conf.Config, log *logging.Logger, wg *sync.WaitGroup, metricService *services.MetricService) {
	defer wg.Done()

	// Создание gin engine.
	r, err := CreateGinEngine(cfg, log, gin.ReleaseMode, metricService)
	if err != nil {
		log.Fatalf("HTTP server failed to create gin engine: %v", err)
	}

	// Регистрация pprof.
	pprof.Register(r)

	// Создаем сервер.
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	// Запуск сервера в горутине.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}()

	log.Infof("HTTP server is running on %v", cfg.Addr)

	// Ожидаем сигнал завершения.
	<-ctx.Done()

	// Инициируем graceful shutdown.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server forced to shutdown:", err)
		return
	}

	log.Info("HTTP server stopped gracefully.")

}
