// Package engine предоставляет функции для создания и настройки Gin engine.
package engine

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

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
func CreateGinEngine(cfg *conf.Config, log *logging.Logger, ginMode string, metricService *services.MetricService) (g *gin.Engine, err error) {
	// Создаем gin engine.
	gin.SetMode(ginMode)
	r := gin.New()

	// Обработка паники gin.
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("gin panic: %v", rec) // nolint:errcheck
		}
	}()

	// Middlewares.
	r.Use(middleware.GzipMiddleware())
	r.Use(middleware.LoggerMiddleware(log))
	if cfg.Key != "" {
		log.Debug("Using signature middleware")
		r.Use(middleware.SignatureMiddleware(log, cfg.Key))
	}

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
