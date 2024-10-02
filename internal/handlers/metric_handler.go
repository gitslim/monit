package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/services"
)

type MetricHandler struct {
	metricService services.MetricService
}

func NewMetricHandler(metricService *services.MetricService) *MetricHandler {
	return &MetricHandler{metricService: *metricService}
}

func writeError(c *gin.Context, err error) {
	var e *errs.Error
	if errors.As(err, &e) {
		// fmt.Printf("errs: %v\n", e)
		c.String(e.Code, e.Error())
		return
	}
	c.String(http.StatusInternalServerError, "Internal server error")
}

// Обновление метрики
func (h *MetricHandler) UpdateMetric(c *gin.Context) {
	mType := c.Param("type")
	mName := c.Param("name")
	mValue := c.Param("value")

	if err := h.metricService.UpdateMetric(mName, mType, mValue); err != nil {
		writeError(c, err)
		return
	}
	c.String(http.StatusOK, "Metric %s updated successfully\n", mName)
}

// Получение метрики
func (h *MetricHandler) GetMetric(c *gin.Context) {
	// mType := c.Param("type")
	mName := c.Param("name")

	m, err := h.metricService.GetMetric(mName)
	if err != nil {
		writeError(c, err)
		return
	}

	c.String(http.StatusOK, "%v", m.GetValue())
}

// HTML со списком метрик
func (h *MetricHandler) ListMetrics(c *gin.Context) {
	metrics := h.metricService.GetAllMetrics()
	// for name, metric := range metrics {
	// 	fmt.Printf("%s: %s\n", name, metric)
	// }
	res := gin.H{
		"metrics": metrics,
	}

	c.HTML(http.StatusOK, "metrics.html", res)
}
