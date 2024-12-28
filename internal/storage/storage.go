package storage

import (
	"github.com/gitslim/monit/internal/entities"
)

// Storager определяет интерфейс для работы с хранилищем метрик.
type Storager interface {
	// UpdadateOrCreateMetric обновляет или создает метрику.
	UpdateOrCreateMetric(mName string, mType entities.MetricType, mValue interface{}) error
	// BatchUpdateOrCreateMetrics обновляет или создает метрики.
	BatchUpdateOrCreateMetrics([]*entities.MetricDTO) error
	// GetMetric получает метрику.
	GetMetric(mName string, mType string) (entities.Metric, error)
	// GetAllMetrics получает все метрики.
	GetAllMetrics() (map[string]entities.Metric, error)
	// Ping проверяет соединение с хранилищем.
	Ping() error
}
