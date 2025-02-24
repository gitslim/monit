package entities

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gitslim/monit/internal/errs"
)

// MetricDTO содержит данные о метрике и может быть сериализована в JSON.
type MetricDTO struct {
	ID          string  `json:"id"`                     // имя метрики
	MType       string  `json:"type"`                   // параметр, принимающий значение gauge или counter
	Delta       int64   `json:"delta,omitempty"`        // значение метрики в случае передачи counter
	Value       float64 `json:"value,omitempty"`        // значение метрики в случае передачи gauge
	StringValue string  `json:"string_value,omitempty"` // значение Delta или Value в виде строки
}

// String возвращает строковое представление MetricDTO.
func (s MetricDTO) String() string {
	return fmt.Sprintf("MetricDTO{ID: %s, MType: %s, Delta: %d, Value: %f, StringValue: %s}", s.ID, s.MType, s.Delta, s.Value, s.StringValue)
}

// GetValue возвращает значение метрики в случае передачи gauge.
func (s MetricDTO) GetValue() float64 {
	return s.Value
}

// GetDelta возвращает значение метрики в случае передачи counter.
func (s MetricDTO) GetDelta() int64 {
	return s.Delta
}

// GetStringValue возвращает строковое представление метрики.
func (s MetricDTO) GetStringValue() string {
	return s.StringValue
}

// MarshalJSON возвращает JSON-представление MetricDTO.
func (s MetricDTO) MarshalJSON() ([]byte, error) {
	v := ""

	switch s.MType {
	case "gauge":
		v = strconv.FormatFloat(s.Value, 'f', -1, 64)
	case "counter":
		v = strconv.FormatInt(s.Delta, 10)
	default:
		return nil, errs.ErrInvalidMetricType
	}

	temp := struct {
		ID          string  `json:"id"`
		MType       string  `json:"type"`
		Delta       int64   `json:"delta,omitempty"`
		Value       float64 `json:"value,omitempty"`
		StringValue string  `json:"string_value"`
	}{
		ID:          s.ID,
		MType:       s.MType,
		Delta:       s.Delta,
		Value:       s.Value,
		StringValue: v,
	}

	return json.Marshal(temp)
}

// UnmarshalJSON - функция для декодирования JSON-данных в структуру MetricDTO.
func (s *MetricDTO) UnmarshalJSON(data []byte) error {
	var err error
	var v string

	var temp struct {
		ID          string  `json:"id"`
		Type        string  `json:"type"`
		Value       float64 `json:"value"`
		Delta       int64   `json:"delta"`
		StringValue string  `json:"string_value"`
	}

	if err = json.Unmarshal(data, &temp); err != nil {
		return err
	}

	switch temp.Type {
	case "counter":
		v = fmt.Sprintf("%d", temp.Delta)
	case "gauge":
		v = fmt.Sprintf("%f", temp.Value)
	default:
		return errs.ErrInvalidMetricType
	}

	// Присваиваем значения из temp в структуру s
	s.ID = temp.ID
	s.MType = temp.Type
	s.Delta = temp.Delta
	s.Value = temp.Value
	s.StringValue = v

	return nil
}

// UnmarshalJSON возвращает MetricDTO из JSON-представления.

// NewCounterMetricDTO создаёт новую MetricDTO с типом counter.
func NewCounterMetricDTO(name string, delta int64) (*MetricDTO, error) {
	var value float64 = 0

	dto := &MetricDTO{
		ID:    name,
		MType: "counter",
		Delta: delta,
		Value: value,
	}
	return dto, nil
}

// NewGaugeMetricDTO создаёт новую MetricDTO с типом gauge.
func NewGaugeMetricDTO(name string, value float64) (*MetricDTO, error) {
	var delta int64 = 0

	dto := &MetricDTO{
		ID:    name,
		MType: "gauge",
		Delta: delta,
		Value: value,
	}
	return dto, nil
}

// NewCounterMetricFromDTO создает новую метрику CounterMetric из MetricDTO.
func NewCounterMetricFromDTO(dto *MetricDTO) (*CounterMetric, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("invalid counter metric dto: %v", dto)
	}
	return &CounterMetric{
		Name:  dto.ID,
		Value: dto.Delta,
	}, nil
}

// NewGaugeMetricFromDTO создает новую метрику GaugeMetric из MetricDTO.
func NewGaugeMetricFromDTO(dto *MetricDTO) (*GaugeMetric, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("invalid gauge metric dto: %v", dto)
	}
	return &GaugeMetric{
		Name:  dto.ID,
		Value: dto.Value,
	}, nil
}

// NewMetricDTO создаёт новую MetricDTO заданного типа.
func NewMetricDTO(mName, mType string, mValue any) (*MetricDTO, error) {
	var dto *MetricDTO
	var value float64
	var delta int64

	t, err := GetMetricType(mType)
	if err != nil {
		return nil, errs.ErrInvalidMetricType
	}

	switch t {
	case Counter:
		switch v := mValue.(type) {
		case int64:
			delta = v
		case string:
			delta, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, errs.ErrInvalidMetricValue
			}
		default:
			return nil, errs.ErrInvalidMetricValue
		}

		dto, err = NewCounterMetricDTO(mName, delta)
		if err != nil {
			return nil, err
		}

	case Gauge:
		switch v := mValue.(type) {
		case float64:
			value = v
		case string:
			value, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, errs.ErrInvalidMetricValue
			}
		default:
			return nil, errs.ErrInvalidMetricValue
		}

		dto, err = NewGaugeMetricDTO(mName, value)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errs.ErrInvalidMetricType
	}

	return dto, nil
}
