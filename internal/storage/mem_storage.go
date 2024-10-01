package storage

import (
	"fmt"
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

func (s *MemStorage) SetMetric(metric entities.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics[metric.GetName()] = metric
	return nil
}

func (s *MemStorage) GetMetric(name string) (entities.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if metric, exists := s.metrics[name]; exists {
		return metric, nil
	}
	return nil, fmt.Errorf("memstorage GetMetric failed %w", errs.ErrMetricNotFound)
}

func (s *MemStorage) GetAllMetrics() map[string]entities.Metric {
	// 	var result []entities.Metric
	// 	for _, metric := range s.metrics {
	// 		result = append(result, metric)
	// 	}
	// 	return result
	return s.metrics
}
