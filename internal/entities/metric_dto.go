package entities

import (
	"strconv"

	"github.com/gitslim/monit/internal/errs"
)

// MetricDTO содержит данные о метрике
type MetricDTO struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewCounterMetricDTO создаёт новую метрику типа counter
func NewCounterMetricDTO(mName, mValue string) (*MetricDTO, error) {
	value, err := strconv.ParseInt(mValue, 10, 64)
	if err != nil {
		return nil, errs.ErrInvalidMetricValue
	}

	return &MetricDTO{
		ID:    mName,
		MType: "counter",
		Delta: &value,
	}, nil
}

// NewGaugeMetricDTO создаёт новую метрику типа gauge
func NewGaugeMetricDTO(mName, mValue string) (*MetricDTO, error) {
	value, err := strconv.ParseFloat(mValue, 64)
	if err != nil {
		return nil, errs.ErrInvalidMetricValue
	}

	return &MetricDTO{
		ID:    mName,
		MType: "gauge",
		Value: &value,
	}, nil
}

// NewMetricDTO создаёт новую метрику заданного типа
func NewMetricDTO(mName, mType string, mValue any) (*MetricDTO, error) {
	t, err := GetMetricType(mType)
	if err != nil {
		return nil, errs.ErrInvalidMetricType
	}

	var mDto *MetricDTO

	switch t {
	case Counter:
		v, ok := mValue.(int64)
		if !ok {
			return nil, errs.ErrInvalidMetricValue
		}
		dto, err := NewCounterMetricDTO(mName, strconv.FormatInt(v, 10))
		if err != nil {
			return nil, err
		}
		mDto = dto

	case Gauge:
		v, ok := mValue.(float64)
		if !ok {
			return nil, errs.ErrInvalidMetricValue
		}
		dto, err := NewGaugeMetricDTO(mName, strconv.FormatFloat(v, 'f', -1, 64))
		if err != nil {
			return nil, err
		}
		mDto = dto

	default:
		return nil, errs.ErrInvalidMetricType
	}
	return mDto, nil
}
