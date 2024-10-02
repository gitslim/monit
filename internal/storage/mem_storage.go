package storage

import (
	"sync"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
)

type MemStorage struct {
	mu      sync.RWMutex
	metrics map[string]entities.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]entities.Metric),
	}
}

func (s *MemStorage) UpdateOrCreateMetric(mName string, mType entities.MetricType, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var m entities.Metric

	metric, err := s.GetMetric(mName)
	if err != nil {
		switch mType {
		case entities.Gauge:
			m = entities.NewGaugeMetric(mName)
		case entities.Counter:
			m = entities.NewCounterMetric(mName)
		}
	} else {
		m = metric
	}

	if err := m.SetValue(value); err != nil {
		return err
	}

	s.metrics[mName] = m
	return nil
}

func (s *MemStorage) GetMetric(mName string) (entities.Metric, error) {
	if metric, exists := s.metrics[mName]; exists {
		return metric, nil
	}
	return nil, errs.ErrMetricNotFound
}

func (s *MemStorage) GetAllMetrics() map[string]entities.Metric {
	// 	var result []entities.Metric
	// 	for _, metric := range s.metrics {
	// 		result = append(result, metric)
	// 	}
	// 	return result
	return s.metrics
}
