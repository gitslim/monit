package services

import (
	"context"
	"strconv"
	"time"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
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
func WithMemStorage(ctx context.Context, log *logging.Logger, cfg *conf.Config, backupErrChan chan<- error) (MetricServiceConf, error) {
	shouldBackupSync := cfg.StoreInterval == 0

	file, err := storage.CreateBackupFile(cfg.FileStoragePath)
	if err != nil {
		return nil, err
	}

	stor := storage.NewMemStorage(shouldBackupSync, file)
	if cfg.Restore {
		// Загружаем данные при запуске
		err := stor.LoadFromFile(cfg.FileStoragePath)
		if err != nil {
			log.Debugf("Failed to load metrics from file: %v", err)
		}
	}

	if cfg.StoreInterval > 0 {
		go stor.StartPeriodicBackup(ctx, log, file, time.Duration(cfg.StoreInterval)*time.Second, backupErrChan)
	}
	return WithStorage(stor), nil
}

func (s *MetricService) GetMetric(name string) (entities.Metric, error) {
	val, err := s.storage.GetMetric(name)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (s *MetricService) UpdateMetric(mName, mType, mValue string) error {
	var v interface{}

	if mName == "" || mType == "" || mValue == "" {
		return errs.ErrMetricNotFound
	}

	t, err := entities.GetMetricType(mType)
	if err != nil {
		return err
	}

	switch t {
	case entities.Gauge:
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return errs.ErrInvalidMetricValue
		}
		v = val
	case entities.Counter:
		val, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return errs.ErrInvalidMetricValue
		}
		v = val
	default:
		return errs.ErrInvalidMetricType
	}

	return s.storage.UpdateOrCreateMetric(mName, t, v)
}

func (s *MetricService) GetAllMetrics() map[string]entities.Metric {
	return s.storage.GetAllMetrics()
}
