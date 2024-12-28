package entities

import (
	"encoding/json"
	"fmt"

	"github.com/gitslim/monit/internal/errs"
)

const (
	Gauge MetricType = iota
	Counter
)

// MetricType представляет тип метрики.
type MetricType int

func (m MetricType) String() string {
	switch m {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	default:
		return ""
	}
}

var metricTypeMap = map[string]MetricType{
	"gauge":   Gauge,
	"counter": Counter,
}

// GetMetricType получает тип метрики из строки.
func GetMetricType(metricTypeStr string) (MetricType, error) {
	if metricType, ok := metricTypeMap[metricTypeStr]; ok {
		return metricType, nil
	}
	return -1, errs.ErrInvalidMetricType
}

// Metric представляет интерфейс для работы с метриками.
type Metric interface {
	// GetName Возвращает имя метрики.
	GetName() string
	// GetType возвращает тип метрики.
	GetType() MetricType
	// GetValue возвращает значение метрики.
	GetValue() interface{}
	// GetStringValue возвращает строковое представление значения метрики.
	GetStringValue() string
	// SetValue устанавливает новое значение метрики.
	SetValue(interface{}) error
	// GetStringValue возвращает строковое представление значения метрики.
	String() string
}

// GaugeMetric реализация метрики Gauge.
type GaugeMetric struct {
	Name  string
	Value float64
}

// MarshalJSON возвращает JSON-сериализованную метрику.
func (g GaugeMetric) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  string     `json:"name"`
		Value float64    `json:"value"`
		Type  MetricType `json:"type"`
	}{
		Name:  g.Name,
		Value: g.Value,
		Type:  g.GetType(),
	})
}

// UnmarshalJSON возвращает JSON-десериализованную метрику.
func (g *GaugeMetric) UnmarshalJSON(data []byte) error {
	var temp struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	g.Name = temp.Name
	g.Value = temp.Value
	return nil
}

// NewGaugeMetric создает новую метрику GaugeMetric.
func NewGaugeMetric(name string) *GaugeMetric {
	return &GaugeMetric{Name: name}
}

// NewGaugeMetricFromDTO создает новую метрику GaugeMetric из DTO.
func NewGaugeMetricFromDTO(dto *MetricDTO) (*GaugeMetric, error) {
	if dto.ID == "" || dto.Value == nil {
		return nil, fmt.Errorf("invalid gauge metric dto: %v", dto)
	}
	return &GaugeMetric{
		Name:  dto.ID,
		Value: *dto.Value,
	}, nil
}

// GetName возвращает имя метрики.
func (g *GaugeMetric) GetName() string {
	return g.Name
}

// GetType возвращает тип метрики.
func (g *GaugeMetric) GetType() MetricType {
	return Gauge
}

// GetValue возвращает значение метрики.
func (g *GaugeMetric) GetValue() interface{} {
	return g.Value
}

// GetStringValue возвращает строковое представление значения метрики.
func (g *GaugeMetric) GetStringValue() string {
	return fmt.Sprintf("%v", g.GetValue())
}

// SetValue устанавливает новое значение метрики.
func (g *GaugeMetric) SetValue(value interface{}) error {
	if v, ok := value.(float64); ok {
		g.Value = v
		return nil
	}

	return errs.ErrInvalidMetricValue
}

// String возвращает строковое представление значения метрики.
func (g *GaugeMetric) String() string {
	return g.GetStringValue()
}

// CounterMetric реализация метрики Counter.
type CounterMetric struct {
	Name  string
	Value int64
}

// MarshalJSON возвращает JSON-сериализованную метрику.
func (c CounterMetric) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  string     `json:"name"`
		Value int64      `json:"value"`
		Type  MetricType `json:"type"`
	}{
		Name:  c.Name,
		Value: c.Value,
		Type:  c.GetType(),
	})
}

// UnmarshalJSON возвращает JSON-десериализованную метрику.
func (c *CounterMetric) UnmarshalJSON(data []byte) error {
	var temp struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	c.Name = temp.Name
	c.Value = temp.Value
	return nil
}

// NewCounterMetric создает новую метрику CounterMetric.
func NewCounterMetric(name string) *CounterMetric {
	return &CounterMetric{Name: name}
}

// NewCounterMetricFromDTO создает новую метрику CounterMetric из DTO.
func NewCounterMetricFromDTO(dto *MetricDTO) (*CounterMetric, error) {
	if dto.ID == "" || dto.Delta == nil {
		return nil, fmt.Errorf("invalid counter metric dto: %v", dto)
	}
	return &CounterMetric{
		Name:  dto.ID,
		Value: *dto.Delta,
	}, nil
}

// GetName возвращает имя метрики.
func (c *CounterMetric) GetName() string {
	return c.Name
}

// GetType возвращает тип метрики.
func (c *CounterMetric) GetType() MetricType {
	return Counter
}

// GetValue возвращает значение метрики.
func (c *CounterMetric) GetValue() interface{} {
	return c.Value
}

// GetStringValue возвращает строковое представление значения метрики.
func (c *CounterMetric) GetStringValue() string {
	return fmt.Sprintf("%v", c.GetValue())
}

// SetValue устанавливает значение метрики.
func (c *CounterMetric) SetValue(value interface{}) error {
	if v, ok := value.(int64); ok {
		c.Value += v
		return nil
	}
	return errs.ErrInvalidMetricValue
}

// String возвращает строковое представление значения метрики.
func (c *CounterMetric) String() string {
	return fmt.Sprintf("%v", c.GetValue())
}
