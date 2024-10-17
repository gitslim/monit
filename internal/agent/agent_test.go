package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/stretchr/testify/assert"
)

func TestSendMetricTableDriven(t *testing.T) {
	client := &http.Client{}
	dummyJSON, _ := json.Marshal(nil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/update/":
			w.Header().Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
			_, err := w.Write(dummyJSON)
			assert.NoError(t, err)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	tests := []struct {
		name        string
		metricType  string
		metricName  string
		value       any
		expectedErr bool
	}{
		{
			name:        "Valid gauge metric",
			metricType:  "gauge",
			metricName:  "test_metric",
			value:       float64(123.45),
			expectedErr: false,
		},
		{
			name:        "Valid counter metric",
			metricType:  "counter",
			metricName:  "test_counter",
			value:       int64(10),
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
			value:       float64(123.45),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto, err := entities.NewMetricDTO(tt.metricName, tt.metricType, tt.value)
			if (err != nil) != tt.expectedErr {
				t.Errorf("sendMetric() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			err = sendMetric(client, server.URL, dto)
			if (err != nil) != tt.expectedErr {
				t.Errorf("sendMetric() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
