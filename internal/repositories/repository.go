package repositories

type MetricRepository interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	ListGauges() map[string]float64
	ListCounters() map[string]int64
}
