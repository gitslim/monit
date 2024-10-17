package entities

import (
	"strconv"

	"github.com/gitslim/monit/internal/errs"
)

type MetricDTO struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewCounterMetricDTO(mName, mType, mValue string) (*MetricDTO, error) {
	value, err := strconv.ParseInt(mValue, 10, 64)
	if err != nil {
		return nil, errs.ErrInvalidMetricValue
	}

	return &MetricDTO{
		ID:    mName,
		MType: mType,
		Delta: &value,
	}, nil
}

func NewGaugeMetricDTO(mName, mType, mValue string) (*MetricDTO, error) {
	value, err := strconv.ParseFloat(mValue, 64)
	if err != nil {
		return nil, errs.ErrInvalidMetricValue
	}

	return &MetricDTO{
		ID:    mName,
		MType: mType,
		Value: &value,
	}, nil
}

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
		dto, err := NewCounterMetricDTO(mName, mType, strconv.FormatInt(v, 10))
		if err != nil {
			return nil, err
		}
		mDto = dto

	case Gauge:
		v, ok := mValue.(float64)
		if !ok {
			return nil, errs.ErrInvalidMetricValue
		}
		dto, err := NewGaugeMetricDTO(mName, mType, strconv.FormatFloat(v, 'f', -1, 64))
		if err != nil {
			return nil, err
		}
		mDto = dto

	default:
		return nil, errs.ErrInvalidMetricType
	}
	return mDto, nil
}
