package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/logging"
)

// MemStorage хранилище метрик в памяти.
type MemStorage struct {
	metrics          sync.Map
	shouldBackupSync bool
	backupWriter     io.Writer
}

// MarshalJSON сериализует данные в json.
func (s *MemStorage) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	s.metrics.Range(func(key, value interface{}) bool {
		metric := value.(entities.Metric)
		tmp[key.(string)] = map[string]interface{}{
			"name":  metric.GetName(),
			"value": metric.GetValue(),
			"type":  metric.GetType(),
		}
		return true
	})
	return json.Marshal(tmp)
}

// UnmarshalJSON десериализует данные из json.
func (s *MemStorage) UnmarshalJSON(data []byte) error {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

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
				s.metrics.Store(name, &gauge)
			} else {
				return err
			}
		case entities.Counter:
			var counter entities.CounterMetric
			if err := json.Unmarshal(raw, &counter); err == nil {
				s.metrics.Store(name, &counter)
			} else {
				return err
			}
		default:
			return fmt.Errorf("unknown metric type: %s", metricType.Type)
		}
	}

	return nil
}

// NewMemStorage - создает новое хранилище метрик в памяти.
func NewMemStorage(shouldBackupSync bool, backupWriter io.Writer) *MemStorage {
	return &MemStorage{
		metrics:          sync.Map{},
		shouldBackupSync: shouldBackupSync,
		backupWriter:     backupWriter,
	}
}

// UpdateOrCreateMetric обновляет значение метрики, если метрика отстутствует то создает ее.
func (s *MemStorage) UpdateOrCreateMetric(mName string, mType entities.MetricType, value interface{}) error {
	var m entities.Metric

	metric, err := s.GetMetric(mName, mType.String())
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

	s.metrics.Store(mName, m)

	if s.shouldBackupSync {
		if err := s.WriteBackup(s.backupWriter); err != nil {
			return err
		}
	}
	return nil
}

// GetMetric получает метрику по имени.
func (s *MemStorage) GetMetric(mName string, mType string) (entities.Metric, error) {
	if metric, exists := s.metrics.Load(mName); exists {
		return metric.(entities.Metric), nil
	}
	return nil, errs.ErrMetricNotFound
}

// GetAllMetrics получает все метрики.
func (s *MemStorage) GetAllMetrics() (map[string]entities.Metric, error) {
	metrics := make(map[string]entities.Metric)
	s.metrics.Range(func(key, value interface{}) bool {
		metrics[key.(string)] = value.(entities.Metric)
		return true
	})
	return metrics, nil
}

// LoadFromFile загружает данные из файла в хранилище.
func (s *MemStorage) LoadFromFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	return nil
}

// WriteBackup сохраняет данные хранилища в файл.
func (s *MemStorage) WriteBackup(w io.Writer) error {
	data, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	if file, ok := w.(*os.File); ok {
		// Очищаем файл.
		if err = file.Truncate(0); err != nil {
			return err
		}

		// Переходим в начало.
		if _, err = file.Seek(0, 0); err != nil {
			return err
		}
	}
	_, err = w.Write(data)
	return err
}

// StartPeriodicBackup запускает периодическое сохранение данных в файл на диске.
func (s *MemStorage) StartPeriodicBackup(ctx context.Context, log *logging.Logger, fd *os.File, interval time.Duration, errChan chan<- error) {
	defer fd.Close()

	for {
		select {
		case <-ctx.Done():
			log.Debug("MemStorage backup stopped")
			return
		case <-time.After(interval):
			if err := s.WriteBackup(fd); err != nil {
				log.Errorf("MemStorage backup error: %v", err)
				errChan <- err
				return
			}
			log.Debug("MemStorage backup success")
		}
	}
}

// CreateBackupFile создает файл для записи файла бэкапа.
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

// Ping проверяет соединение с хранилищем.
func (s *MemStorage) Ping() error {
	return nil
}

// BatchUpdateOrCreateMetrics обновляет данные в хранилище батчами.
func (s *MemStorage) BatchUpdateOrCreateMetrics(metrics []*entities.MetricDTO) error {
	for _, dto := range metrics {

		var m entities.Metric

		mType, err := entities.GetMetricType(dto.MType)
		if err != nil {
			fmt.Printf("Unknown metric type: %v\n", err)
		}
		switch mType {
		case entities.Gauge:
			m, err = entities.NewGaugeMetricFromDTO(dto)
			if err != nil {
				fmt.Printf("Error creating gauge metric: %v\n", err)
				continue
			}
		case entities.Counter:
			m, err = entities.NewCounterMetricFromDTO(dto)
			if err != nil {
				fmt.Printf("Error creating counter metric: %v\n", err)
				continue
			}
		default:
			fmt.Printf("Unknown metric type: %v\n", mType)
			continue
		}

		s.metrics.Store(m.GetName(), m)

		if s.shouldBackupSync {
			if err := s.WriteBackup(s.backupWriter); err != nil {
				return err
			}
		}
	}
	return nil
}
