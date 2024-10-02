package entities

import (
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
