package storage

import (
	"github.com/gitslim/monit/internal/entities"
)

// Storage - интерфейс для работы с хранилищем метрик
type Storage interface {
	UpdateOrCreateMetric(string, entities.MetricType, interface{}) error
	GetMetric(name string) (entities.Metric, error)
	GetAllMetrics() map[string]entities.Metric
}
