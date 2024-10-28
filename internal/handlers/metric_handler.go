package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/gitslim/monit/internal/services"
)

type MetricHandler struct {
	metricService services.MetricService
}

func NewMetricHandler(metricService *services.MetricService) *MetricHandler {
	return &MetricHandler{metricService: *metricService}
}

func isJSONRequest(c *gin.Context) bool {
	return c.GetHeader(httpconst.HeaderContentType) == httpconst.ContentTypeJSON
}

func writeError(c *gin.Context, err error) {
	fmt.Printf("GetMetric error: %v\n", err)
	var e *errs.Error
	if errors.As(err, &e) {
		if isJSONRequest(c) {
			c.JSON(e.Code, e.Error())
		} else {
			c.String(e.Code, e.Error())
		}
		return
	}

	if isJSONRequest(c) {
		c.JSON(http.StatusInternalServerError, e.Error())
	} else {
		c.String(http.StatusInternalServerError, "Internal server error")
	}
}

// UpdateMetric обновляет метрику
func (h *MetricHandler) UpdateMetric(c *gin.Context) {
	var mType, mName, mValue string

	if isJSONRequest(c) {
		var dto *entities.MetricDTO

		err := json.NewDecoder(c.Request.Body).Decode(&dto)
		if err != nil {
			writeError(c, err)
			return
		}

		mType = dto.MType
		mName = dto.ID

		if dto.Delta != nil {
			mValue = strconv.FormatInt(*dto.Delta, 10)
		} else if dto.Value != nil {
			mValue = strconv.FormatFloat(*dto.Value, 'f', -1, 64)
		}

	} else {
		mType = c.Param("type")
		mName = c.Param("name")
		mValue = c.Param("value")
	}

	if err := h.metricService.UpdateMetric(mName, mType, mValue); err != nil {
		writeError(c, err)
		return
	}
	if isJSONRequest(c) {
		m, err := h.metricService.GetMetric(mName, mType)
		if err != nil {
			writeError(c, err)
			return
		}
		dto, err := entities.NewMetricDTO(m.GetName(), m.GetType().String(), m.GetValue())
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	} else {
		c.String(http.StatusOK, "Metric %s updated successfully\n", mName)
	}
}

// BatchUpdateMetrics обновляет метрики батчами
func (h *MetricHandler) BatchUpdateMetrics(c *gin.Context) {
	var metrics []*entities.MetricDTO

	err := json.NewDecoder(c.Request.Body).Decode(&metrics)
	if err != nil {
		writeError(c, err)
		return
	}

	if err := h.metricService.BatchUpdateMetrics(metrics); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// Получение метрики
func (h *MetricHandler) GetMetric(c *gin.Context) {
	var mName, mType string

	if isJSONRequest(c) {
		var dto *entities.MetricDTO

		err := json.NewDecoder(c.Request.Body).Decode(&dto)
		if err != nil {
			writeError(c, err)
			return
		}

		mName = dto.ID
		mType = dto.MType

	} else {
		mName = c.Param("name")
		mType = c.Param("type")
	}

	m, err := h.metricService.GetMetric(mName, mType)
	if err != nil {
		writeError(c, err)
		return
	}

	if isJSONRequest(c) {
		dto, err := entities.NewMetricDTO(m.GetName(), m.GetType().String(), m.GetValue())
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	} else {
		c.String(http.StatusOK, "%v", m.GetValue())
	}
}

// HTML со списком метрик
func (h *MetricHandler) ListMetrics(c *gin.Context) {
	metrics, err := h.metricService.GetAllMetrics()
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	res := gin.H{
		"metrics": metrics,
	}

	c.HTML(http.StatusOK, "metrics.html", res)
}

func (h *MetricHandler) PingStorage(c *gin.Context) {
	if err := h.metricService.PingStorage(); err != nil {
		c.String(http.StatusInternalServerError, "error")
	} else {
		c.String(http.StatusOK, "ok")
	}
}
