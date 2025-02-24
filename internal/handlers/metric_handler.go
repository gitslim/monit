package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

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

// req2MetricDTO возвращает MetricDTO из запроса.
func req2MetricDTO(c *gin.Context) (*entities.MetricDTO, error) {

	if isJSONRequest(c) {
		jsonBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}
		var dto entities.MetricDTO
		if err := json.Unmarshal(jsonBytes, &dto); err != nil {
			return nil, err
		}
		return &dto, nil
	} else {
		mName := c.Param("name")
		mType := c.Param("type")
		mValue := c.Param("value")
		dto, err := entities.NewMetricDTO(mName, mType, mValue)
		if err != nil {
			return nil, err
		}
		return dto, nil
	}
}

// req2MetricRequestDTO возвращает MetricRequestDTO из запроса
func req2MetricRequestDTO(c *gin.Context) (*entities.MetricRequestDTO, error) {

	if isJSONRequest(c) {
		jsonBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}
		var dto entities.MetricRequestDTO
		if err := json.Unmarshal(jsonBytes, &dto); err != nil {
			return nil, err
		}
		return &dto, nil
	} else {
		mName := c.Param("name")
		mType := c.Param("type")
		dto, err := entities.NewMetricRequestDTO(mName, mType)
		if err != nil {
			return nil, err
		}
		return dto, nil
	}
}

// UpdateMetric обновляет метрику.
func (h *MetricHandler) UpdateMetric(c *gin.Context) {

	dto, err := req2MetricDTO(c)
	if err != nil {
		writeError(c, errs.ErrBadRequest)
		return
	}

	if dto.ID == "" {
		writeError(c, errs.ErrMetricNotFound)
		return
	}

	if err = h.metricService.UpdateMetric(*dto); err != nil {
		writeError(c, err)
		return
	}

	if isJSONRequest(c) {
		metric, err := h.metricService.GetMetric(dto.ID, dto.MType)
		if err != nil {
			writeError(c, err)
			return
		}
		dto, err = entities.NewMetricDTO(metric.GetName(), metric.GetType().String(), metric.GetValue())
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto)
		return
	}

	c.String(http.StatusOK, "Metric %s updated successfully\n", dto.ID)
}

// BatchUpdateMetrics обновляет метрики батчами.
func (h *MetricHandler) BatchUpdateMetrics(c *gin.Context) {
	var dtos []*entities.MetricDTO

	jsonBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		writeError(c, err)
		return
	}

	if err := json.Unmarshal(jsonBytes, &dtos); err != nil {
		writeError(c, err)
		return
	}

	if err := h.metricService.BatchUpdateMetrics(dtos); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, "{}")
}

// GetMetric возвращает метрику по имени и типу.
func (h *MetricHandler) GetMetric(c *gin.Context) {
	dto, err := req2MetricRequestDTO(c)
	if err != nil {
		writeError(c, errs.ErrBadRequest)
		return
	}

	m, err := h.metricService.GetMetric(dto.ID, dto.MType.String())
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
