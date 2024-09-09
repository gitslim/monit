package storage

import (
	"sync"
)

// MemStorage реализует интерфейс MetricsRepository
type MemStorage struct {
	mu       sync.Mutex
	gauges   map[string]float64
	counters map[string]int64
}

// Конструктор
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauges[name] = value
}

func (s *MemStorage) UpdateCounter(name string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[name] += value
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {
	val, ok := s.gauges[name]
	return val, ok
}

func (s *MemStorage) GetCounter(name string) (int64, bool) {
	val, ok := s.counters[name]
	return val, ok
}

func (s *MemStorage) ListGauges() map[string]float64 {
	return s.gauges
}

func (s *MemStorage) ListCounters() map[string]int64 {
	return s.counters
}
