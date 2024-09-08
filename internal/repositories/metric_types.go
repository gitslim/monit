package repositories

import "strconv"

// Тип метрики
type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

// Интерфейс метрики
type Metric interface {
	Update(value string) error
	Get() string
}

// Gauge метрика
type GaugeMetric struct {
	value float64
}

func NewGaugeMetric() *GaugeMetric {
	return &GaugeMetric{}
}

func (m *GaugeMetric) Update(value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.value = v
	return nil
}

func (m *GaugeMetric) Get() string {
	return strconv.FormatFloat(m.value, 'f', -1, 64)
}

// Counter метрика
type CounterMetric struct {
	value int64
}

func NewCounterMetric() *CounterMetric {
	return &CounterMetric{}
}

func (m *CounterMetric) Update(value string) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	m.value += v
	return nil
}

func (m *CounterMetric) Get() string {
	return strconv.FormatInt(m.value, 10)
}
