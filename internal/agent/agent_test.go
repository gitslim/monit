package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMetricTableDriven(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/update/gauge/test_metric/123.45" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/update/counter/test_counter/10" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sendMetric(server.URL, tt.metricType, tt.metricName, tt.value)

			if (err != nil) != tt.expectedErr {
				t.Errorf("sendMetric() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
