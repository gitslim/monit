package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/sender"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSendMetrics тестирует отправку метрик.
func TestSendMetrics(t *testing.T) {
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
		name       string
		metricType string
		metricName string
		value      any
		marshalErr bool
		sendErr    bool
	}{
		{
			name:       "Valid gauge metric",
			metricType: "gauge",
			metricName: "test_metric",
			value:      float64(123.45),
			marshalErr: false,
			sendErr:    false,
		},
		{
			name:       "Valid counter metric",
			metricType: "counter",
			metricName: "test_counter",
			value:      int64(10),
			marshalErr: false,
			sendErr:    false,
		},
		{
			name:       "Invalid path metric",
			metricType: "gauge",
			metricName: "invalid_metric",
			value:      "invalid_value",
			marshalErr: false,
			sendErr:    false,
		},
		{
			name:       "Invalid metric type",
			metricType: "invalid_type",
			metricName: "test_metric",
			value:      float64(123.45),
			marshalErr: false,
			sendErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var err error

			dto, err := entities.NewMetricDTO(tt.metricName, tt.metricType, tt.value)
			if tt.marshalErr {
				require.Error(t, err)
			}

			metrics := []*entities.MetricDTO{dto}
			cfg := &conf.Config{Addr: server.URL}

			err = sender.SendMetrics(ctx, cfg, client, metrics, false)
			if tt.sendErr {
				require.Error(t, err)
			}

			err = sender.SendMetrics(ctx, cfg, client, metrics, true)
			if tt.sendErr {
				require.Error(t, err)
			}
		})
	}
}
