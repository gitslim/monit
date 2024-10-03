package entities

import (
	"encoding/json"
	"fmt"

	"github.com/gitslim/monit/internal/errs"
)

// MetricType - типы метрик
type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

var metricTypeMap = map[string]MetricType{
	"gauge":   Gauge,
	"counter": Counter,
}

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

// GetMetricType получает тип метрики из строки
func GetMetricType(metricTypeStr string) (MetricType, error) {
	if metricType, ok := metricTypeMap[metricTypeStr]; ok {
		return metricType, nil
	}
	return -1, errs.ErrInvalidMetricType
}

// Metric - интерфейс для метрик
type Metric interface {
	GetName() string
	GetType() MetricType
	GetValue() interface{}
	GetStringValue() string
	SetValue(interface{}) error
	String() string
}

// GaugeMetric - реализация Gauge метрики
type GaugeMetric struct {
	Name  string
	Value float64
}

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

func NewGaugeMetric(name string) *GaugeMetric {
	return &GaugeMetric{Name: name}
}

func (g *GaugeMetric) GetName() string {
	return g.Name
}

func (g *GaugeMetric) GetType() MetricType {
	return Gauge
}

func (g *GaugeMetric) GetValue() interface{} {
	return g.Value
}

func (g *GaugeMetric) GetStringValue() string {
	return fmt.Sprintf("%v", g.GetValue())
}

func (g *GaugeMetric) SetValue(value interface{}) error {
	if v, ok := value.(float64); ok {
		g.Value = v
		return nil
	}

	return errs.ErrInvalidMetricValue
}

func (g *GaugeMetric) String() string {
	return g.GetStringValue()
}

// CounterMetric - реализация Counter метрики
type CounterMetric struct {
	Name  string
	Value int64
}

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

func NewCounterMetric(name string) *CounterMetric {
	return &CounterMetric{Name: name}
}

func (c *CounterMetric) GetName() string {
	return c.Name
}

func (c *CounterMetric) GetType() MetricType {
	return Counter
}

func (c *CounterMetric) GetValue() interface{} {
	return c.Value
}

func (c *CounterMetric) GetStringValue() string {
	return fmt.Sprintf("%v", c.GetValue())
}

func (c *CounterMetric) SetValue(value interface{}) error {
	if v, ok := value.(int64); ok {
		c.Value += v
		return nil
	}
	return errs.ErrInvalidMetricValue
}

func (c *CounterMetric) String() string {
	return fmt.Sprintf("%v", c.GetValue())
}
