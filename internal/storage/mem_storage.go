package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/logging"
)

type MemStorage struct {
	mu         sync.RWMutex
	metrics    map[string]entities.Metric
	syncBackup bool
	backupFd   *os.File
}

func (s *MemStorage) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for name, metric := range s.metrics {
		tmp[name] = map[string]interface{}{
			"name":  metric.GetName(),
			"value": metric.GetValue(),
			"type":  metric.GetType(),
		}
	}
	return json.Marshal(tmp)
}

func (s *MemStorage) UnmarshalJSON(data []byte) error {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	s.metrics = make(map[string]entities.Metric)
	for name, raw := range temp {
		var metricType struct {
			Type entities.MetricType `json:"type"`
		}
		if err := json.Unmarshal(raw, &metricType); err != nil {
			return err
		}

		switch metricType.Type {
		case entities.Gauge:
			var gauge entities.GaugeMetric
			if err := json.Unmarshal(raw, &gauge); err == nil {
				s.metrics[name] = &gauge
			}
		case entities.Counter:
			var counter entities.CounterMetric
			if err := json.Unmarshal(raw, &counter); err == nil {
				s.metrics[name] = &counter
			}
		default:
			return fmt.Errorf("unknown metric type: %s", metricType.Type)
		}
	}
	return nil
}

func NewMemStorage(syncBackup bool, fd *os.File) *MemStorage {
	return &MemStorage{
		metrics:    make(map[string]entities.Metric),
		syncBackup: syncBackup,
		backupFd:   fd,
	}
}

// UpdateOrCreateMetric обновляет значение метрики, если метрика отстутствует то создает ее
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
	if s.syncBackup {
		if err := s.SaveToFile(s.backupFd); err != nil {
			return err
		}
	}
	return nil
}

// GetMetric получает метрику по имени
func (s *MemStorage) GetMetric(mName string) (entities.Metric, error) {
	if metric, exists := s.metrics[mName]; exists {
		return metric, nil
	}
	return nil, errs.ErrMetricNotFound
}

// GetAllMetrics получает все метрики
func (s *MemStorage) GetAllMetrics() map[string]entities.Metric {
	return s.metrics
}

// LoadFromFile - загружает данные из файла в хранилище
func (s *MemStorage) LoadFromFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// var metrics map[string]entities.Metric
	var storage MemStorage
	err = json.Unmarshal(data, &storage)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.metrics = storage.metrics
	s.mu.Unlock()

	return nil
}

// SaveToFile - сохраняет данные хранилища в файл
func (s *MemStorage) SaveToFile(fd *os.File) error {
	// s.mu.RLock()
	// defer s.mu.RUnlock()

	data, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	// очищаем файл
	if err := fd.Truncate(0); err != nil {
		return err
	}

	// переходим в начало
	if _, err := fd.Seek(0, 0); err != nil {
		return err
	}
	_, err = fd.Write(data)
	return err
}

// StartPeriodicBackup - запускает периодическое сохранение данных на диск
func (s *MemStorage) StartPeriodicBackup(ctx context.Context, log *logging.Logger, fd *os.File, interval time.Duration, errChan chan<- error) {
	defer fd.Close()

	for {
		select {
		case <-ctx.Done():
			log.Debug("MemStorage backup stopped")
			return
		case <-time.After(interval):
			if err := s.SaveToFile(fd); err != nil {
				log.Errorf("MemStorage backup error: %v", err)
				errChan <- err
				return
			}
			log.Debug("MemStorage backup success")
		}
	}
}

// CreateBackupFile создает файл для записи бэкапа
func CreateBackupFile(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)
	if dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}

	return fd, nil
}
