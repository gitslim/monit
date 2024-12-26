package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/middleware"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/services"
)

func CreateGinEngine(cfg *conf.Config, log *logging.Logger, ginMode string, templatesGlob string, metricService *services.MetricService) (*gin.Engine, error) {
	// Gin engine
	r := gin.New()

	gin.SetMode(ginMode)

	// middlewares
	r.Use(middleware.GzipMiddleware())
	r.Use(middleware.LoggerMiddleware(log))
	if cfg.Key != "" {
		log.Debug("Using signature middleware")
		r.Use(middleware.SignatureMiddleware(log, cfg.Key))
	}

	r.LoadHTMLGlob(templatesGlob)

	// создание хендлера
	metricHandler := NewMetricHandler(metricService)

	// роуты
	r.GET("/", metricHandler.ListMetrics)
	r.POST("/update/", metricHandler.UpdateMetric)
	r.POST("/updates/", metricHandler.BatchUpdateMetrics)
	r.POST("/value/", metricHandler.GetMetric)
	r.GET("/value/:type/:name", metricHandler.GetMetric)
	r.POST("/update/:type/:name/:value", metricHandler.UpdateMetric)
	r.GET("/ping", metricHandler.PingStorage)

	return r, nil
}
