package services

import (
	"strconv"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/storage"
)

// MetricService сервис для работы с метриками
type MetricService struct {
	storage storage.Storage
}

// MetricServiceConf конфиг для MetricService
type MetricServiceConf func(svc *MetricService) error

// NewMetricService создает новый сервис MetricService, применяя к нему все конфиги
func NewMetricService(cfgs ...MetricServiceConf) (*MetricService, error) {
	svc := &MetricService{}

	for _, cfg := range cfgs {
		err := cfg(svc)
		if err != nil {
			return nil, err
		}
	}
	return svc, nil
}

// WithStorage конфигурирует MetricService с заданным заданный Storage
func WithStorage(stor storage.Storage) MetricServiceConf {
	return func(svc *MetricService) error {
		svc.storage = stor
		return nil
	}
}

// WithMemStorage конфигурирует MetricService c MemStorage
func WithMemStorage() MetricServiceConf {
	storage := storage.NewMemStorage()
	return WithStorage(storage)
}

func (s *MetricService) GetMetric(name string) (entities.Metric, error) {
	val, err := s.storage.GetMetric(name)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("GetMetric: %v\n", val)
	return val, nil
}

func (s *MetricService) SetMetric(mName, mType, mValue string) error {
	var m entities.Metric
	var v interface{}

	if mName == "" || mType == "" || mValue == "" {
		return errs.ErrMetricNotFound
	}

	t, err := entities.GetMetricType(mType)
	if err != nil {
		return err
	}
	// fmt.Printf("SetMetric: %v\n", mValue)

	switch t {
	case entities.Gauge:
		m = entities.NewGaugeMetric(mName)
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return errs.ErrInvalidMetricValue
		}
		v = val
	case entities.Counter:
		m = entities.NewCounterMetric(mName)
		val, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return errs.ErrInvalidMetricValue
		}
		v = val
	default:
		return errs.ErrInvalidMetricType
	}

	if err := m.SetValue(v); err != nil {
		return err
	}
	return s.storage.SetMetric(m)
}

func (s *MetricService) GetAllMetrics() map[string]entities.Metric {
	return s.storage.GetAllMetrics()
}
