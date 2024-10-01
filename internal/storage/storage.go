package storage

import (
	"github.com/gitslim/monit/internal/entities"
)

// Storage - интерфейс для работы с хранилищем метрик
type Storage interface {
	SetMetric(entities.Metric) error
	GetMetric(name string) (entities.Metric, error)
	GetAllMetrics() map[string]entities.Metric
}
