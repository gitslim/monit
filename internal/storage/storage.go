package storage

import (
	"fmt"
	"sync"

	"github.com/gitslim/monit/internal/repositories"
)

// MemStorage реализует интерфейс MetricsRepository и хранит метрики
type MemStorage struct {
	mu      sync.Mutex
	metrics map[string]repositories.Metric
}

// Создаем новое хранилище
func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]repositories.Metric),
	}
}

// Обновление метрики
func (s *MemStorage) UpdateMetric(metricType repositories.MetricType, name string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	metric, exists := s.metrics[name]
	if !exists {
		// Если метрика не существует, создаем новую в зависимости от типа
		switch metricType {
		case repositories.GaugeType:
			metric = repositories.NewGaugeMetric()
		case repositories.CounterType:
			metric = repositories.NewCounterMetric()
		default:
			return fmt.Errorf("unknown metric type: %s", metricType)
		}
		s.metrics[name] = metric
	}

	return metric.Update(value)
}

// Получение метрики
func (s *MemStorage) GetMetric(name string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metric, ok := s.metrics[name]
	if !ok {
		return "", false
	}
	return metric.Get(), true
}
