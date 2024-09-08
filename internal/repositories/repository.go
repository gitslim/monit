package repositories

type MetricsRepository interface {
	UpdateMetric(metricType MetricType, name string, value string) error
	GetMetric(name string) (string, bool)
}
