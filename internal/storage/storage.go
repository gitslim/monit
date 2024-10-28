package storage

import (
	"github.com/gitslim/monit/internal/entities"
)

// Storage - интерфейс для работы с хранилищем метрик
type Storage interface {
	UpdateOrCreateMetric(mName string, mType entities.MetricType, mValue interface{}) error
	BatchUpdateOrCreateMetrics([]*entities.MetricDTO) error
	GetMetric(mName string, mType string) (entities.Metric, error)
	GetAllMetrics() (map[string]entities.Metric, error)
	Ping() error
}
