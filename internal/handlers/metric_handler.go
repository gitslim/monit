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

// MetricHandler представляет обработчик метрик.
type MetricHandler struct {
	metricService services.MetricService
}

// NewMetricHandler создает обработчик метрик.
func NewMetricHandler(metricService *services.MetricService) *MetricHandler {
	return &MetricHandler{metricService: *metricService}
}

// isJSONRequest проверяет является ли запрос JSON.
func isJSONRequest(c *gin.Context) bool {
	return c.GetHeader(httpconst.HeaderContentType) == httpconst.ContentTypeJSON
}

// writeError записывает ошибку в ответ сервера.
func writeError(c *gin.Context, err error) {
	fmt.Printf("Error: %v\n", err)
	var e *errs.Error
	if errors.As(err, &e) {
		if isJSONRequest(c) {
			c.JSON(e.Code, e.Error())
		} else {
			c.String(e.Code, e.Error())
		}
		return
	}

	msg := "Internal server error"
	if isJSONRequest(c) {
		c.JSON(http.StatusInternalServerError, msg)
	} else {
		c.String(http.StatusInternalServerError, msg)
	}
}

// UpdateMetric обновляет метрику.
func (h *MetricHandler) UpdateMetric(c *gin.Context) {
	var mType, mName, mValue string

	if isJSONRequest(c) {
		var dto *entities.MetricDTO

		err := json.NewDecoder(c.Request.Body).Decode(&dto)
		if err != nil {
			writeError(c, errs.ErrBadRequest)
			return
		}

		mType = dto.MType
		mName = dto.ID

		if mType == "counter" && dto.Delta != nil {
			mValue = strconv.FormatInt(*dto.Delta, 10)
		} else if mType == "gauge" && dto.Value != nil {
			mValue = strconv.FormatFloat(*dto.Value, 'f', -1, 64)
		} else {
			writeError(c, errs.ErrBadRequest)
			return
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

// BatchUpdateMetrics обновляет метрики батчами.
func (h *MetricHandler) BatchUpdateMetrics(c *gin.Context) {
	var metrics []*entities.MetricDTO

	err := json.NewDecoder(c.Request.Body).Decode(&metrics)
	if err != nil {
		writeError(c, fmt.Errorf("error decoding JSON: %w", err))
		return
	}
	fmt.Printf("Metrics: %v\n", metrics)

	if err := h.metricService.BatchUpdateMetrics(metrics); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetMetric возвращает метрику по имени и типу.
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

// ListMetrics возвращает список метрик в виде HTML-страницы.
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

// PingStorage проверяет соединение с хранилищем.
func (h *MetricHandler) PingStorage(c *gin.Context) {
	if err := h.metricService.PingStorage(); err != nil {
		c.String(http.StatusInternalServerError, "error")
	} else {
		c.String(http.StatusOK, "ok")
	}
}
