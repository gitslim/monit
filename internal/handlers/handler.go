// /internal/handlers/handler.go
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gitslim/monit/internal/repositories"
)

type MetricsHandler struct {
	storage repositories.MetricsRepository
}

// Создаем новый обработчик
func NewMetricsHandler(storage repositories.MetricsRepository) *MetricsHandler {
	return &MetricsHandler{storage: storage}
}

// Обработчик обновления метрик
func (h *MetricsHandler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	metricType := repositories.MetricType(r.PathValue("type"))
	metricName := r.PathValue("name")
	metricValue := r.PathValue("value")

	if metricType != repositories.GaugeType && metricType != repositories.CounterType {
		http.Error(w, fmt.Sprintf("Invalid metric type: %s", metricType), http.StatusBadRequest)
		return
	}

	if metricName == "" {
		http.Error(w, "Metric name missing", http.StatusNotFound)
		return
	}

	if err := h.storage.UpdateMetric(metricType, metricName, metricValue); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update metric: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Metric %s updated successfully\n", metricName)
}
