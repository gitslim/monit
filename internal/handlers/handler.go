// /internal/handlers/handler.go
package handlers

import (
	"fmt"
	"net/http"
	"strings"

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

	path := strings.TrimPrefix(r.URL.Path, "/update/")
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	metricType := repositories.MetricType(parts[0])
	metricName := parts[1]
	metricValue := parts[2]

	fmt.Printf("update metric request: type: %s name: %s value: %s\n", metricType, metricName, metricValue)

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
