package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/services"
)

type MetricHandler struct {
	metricService services.MetricService
}

func NewMetricHandler(metricService *services.MetricService) *MetricHandler {
	return &MetricHandler{metricService: *metricService}
}

func isJSONRequest(c *gin.Context) bool {
	return c.GetHeader("Content-type") == "application/json"
}

func writeError(c *gin.Context, err error) {
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

// Обновление метрики
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
		m, err := h.metricService.GetMetric(mName)
		if err != nil {
			writeError(c, err)
			return
		}
		dto, err := entities.NewMetricDTO(m.GetName(), m.GetType().String(), m.GetStringValue())
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
	} else {
		c.String(http.StatusOK, "Metric %s updated successfully\n", mName)
	}
}

// Получение метрики
func (h *MetricHandler) GetMetric(c *gin.Context) {
	var mName string

	if isJSONRequest(c) {
		var dto *entities.MetricDTO

		err := json.NewDecoder(c.Request.Body).Decode(&dto)
		if err != nil {
			writeError(c, err)
			return
		}

		mName = dto.ID

	} else {
		mName = c.Param("name")
	}

	m, err := h.metricService.GetMetric(mName)
	if err != nil {
		writeError(c, err)
		return
	}

	if isJSONRequest(c) {
		dto, err := entities.NewMetricDTO(m.GetName(), m.GetType().String(), m.GetStringValue())
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
	metrics := h.metricService.GetAllMetrics()
	res := gin.H{
		"metrics": metrics,
	}

	c.HTML(http.StatusOK, "metrics.html", res)
}
