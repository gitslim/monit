// /internal/handlers/handler.go
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/repositories"
)

type MetricsHandler struct {
	storage repositories.MetricsRepository
}

// Создаем новый обработчик
func NewMetricsHandler(storage repositories.MetricsRepository) *MetricsHandler {
	return &MetricsHandler{storage: storage}
}

// Обработчик обновления метрики
func (h *MetricsHandler) UpdateMetrics(c *gin.Context) {
	metricType := repositories.MetricType(c.Param("type"))
	metricName := c.Param("name")
	metricValue := c.Param("value")

	fmt.Printf("update metric request: type: %s name: %s value: %s\n", metricType, metricName, metricValue)

	if metricType != repositories.GaugeType && metricType != repositories.CounterType {
		c.String(http.StatusBadRequest, "Invalid metric type: %s", metricType)
		return
	}

	if metricName == "" {
		c.String(http.StatusNotFound, "Metric name missing")
		return
	}

	if metricValue == "" {
		c.String(http.StatusBadRequest, "Metric value missing")
		return
	}

	if err := h.storage.UpdateMetric(metricType, metricName, metricValue); err != nil {
		c.String(http.StatusBadRequest, "Failed to update metric: %v", err)
		return
	}

	c.String(http.StatusOK, "Metric %s updated successfully\n", metricName)
}

// Обработчик получения метрики
func (h *MetricsHandler) GetMetric(c *gin.Context) {
	metricType := repositories.MetricType(c.Param("type"))
	metricName := c.Param("name")

	if metricType != repositories.GaugeType && metricType != repositories.CounterType {
		c.String(http.StatusBadRequest, "Invalid metric type: %s", metricType)
		return
	}

	if metricName == "" {
		c.String(http.StatusNotFound, "Metric name missing")
		return
	}

	val, exists := h.storage.GetMetric(metricName)
	if !exists {
		c.String(http.StatusNotFound, "No such metric: %v", metricName)
		return
	}

	c.String(http.StatusOK, val)
}

// Обработчик вывода списка метрик
func (h *MetricsHandler) ListMetrics(c *gin.Context) {
	res := gin.H{
		"metrics": h.storage.ListMetrics(),
	}

	c.HTML(http.StatusOK, "metrics.html", res)
}
