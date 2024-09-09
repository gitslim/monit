package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Табличный тест для функции sendMetric
func TestSendMetricTableDriven(t *testing.T) {
	// Создаем mock HTTP-сервер для обработки запросов
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Путь запроса должен совпадать с тем, что ожидается в тесте
		if r.URL.Path == "/update/gauge/test_metric/123.45" {
			w.WriteHeader(http.StatusOK) // Возвращаем успех для конкретного пути
		} else if r.URL.Path == "/update/counter/test_counter/10" {
			w.WriteHeader(http.StatusOK) // Возвращаем успех для другого пути
		} else {
			w.WriteHeader(http.StatusBadRequest) // Все остальные пути считаем ошибочными
		}
	}))
	defer server.Close()

	// Заменим глобальный serverURL на URL mock-сервера
	//serverURL = server.URL

	// Таблица тестов
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		value       string
		expectedErr bool
	}{
		{
			name:        "Valid gauge metric",
			metricType:  "gauge",
			metricName:  "test_metric",
			value:       "123.45",
			expectedErr: false,
		},
		{
			name:        "Valid counter metric",
			metricType:  "counter",
			metricName:  "test_counter",
			value:       "10",
			expectedErr: false,
		},
		{
			name:        "Invalid path metric",
			metricType:  "gauge",
			metricName:  "invalid_metric",
			value:       "invalid_value",
			expectedErr: true,
		},
		{
			name:        "Invalid metric type",
			metricType:  "invalid_type",
			metricName:  "test_metric",
			value:       "123.45",
			expectedErr: true,
		},
	}

	// Проход по каждому тестовому случаю
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sendMetric(server.URL, tt.metricType, tt.metricName, tt.value)

			if (err != nil) != tt.expectedErr {
				t.Errorf("sendMetric() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
