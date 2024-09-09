// /internal/handlers/handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/repositories"
)

type MetricHandler struct {
	storage repositories.MetricRepository
}

// Конструктор
func NewMetricHandler(storage repositories.MetricRepository) *MetricHandler {
	return &MetricHandler{storage: storage}
}

// Обновление метрики
func (h *MetricHandler) UpdateMetric(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")

	if metricName == "" {
		c.String(http.StatusNotFound, "Metric name missing")
		return
	}

	// if metricValue == "" {
	// 	c.String(http.StatusBadRequest, "Metric value missing")
	// 	return
	// }

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid value for gauge")
			return
		}
		h.storage.UpdateGauge(metricName, value)

	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid value for counter")
			return
		}
		h.storage.UpdateCounter(metricName, value)

	default:
		c.String(http.StatusBadRequest, "Invalid metric type")
		return
	}

	c.String(http.StatusOK, "Metric %s updated successfully\n", metricName)
}

// Получение метрики
func (h *MetricHandler) GetMetric(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	if metricName == "" {
		c.String(http.StatusNotFound, "Metric name missing")
		return
	}

	switch metricType {
	case "gauge":
		val, exists := h.storage.GetGauge(metricName)
		if !exists {
			c.String(http.StatusNotFound, "No such metric: %v", metricName)
			return
		}
		c.String(http.StatusOK, "%v", val)

	case "counter":
		val, exists := h.storage.GetCounter(metricName)
		if !exists {
			c.String(http.StatusNotFound, "No such metric: %v", metricName)
			return
		}
		c.String(http.StatusOK, "%v", val)

	default:
		c.String(http.StatusBadRequest, "Invalid metric type")
		return
	}

}

// HTML со списком метрик
func (h *MetricHandler) ListMetrics(c *gin.Context) {
	res := gin.H{
		"gauges":   h.storage.ListGauges(),
		"counters": h.storage.ListCounters(),
	}

	c.HTML(http.StatusOK, "metrics.html", res)
}
