package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/handlers"
	"github.com/gitslim/monit/internal/repositories"
)

func Start(addr string, storage repositories.MetricRepository) error {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Создание хендлера
	metricHandler := handlers.NewMetricHandler(storage)

	// роуты
	r.GET("/", metricHandler.ListMetrics)
	r.GET("/value/:type/:name", metricHandler.GetMetric)
	r.POST("/update/:type/:name/:value", metricHandler.UpdateMetric)

	return r.Run(addr)
}
