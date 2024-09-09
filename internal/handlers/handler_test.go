package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitslim/monit/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	memStorage := storage.NewMemStorage()
	handler := NewMetricsHandler(memStorage)
	type metric struct {
		typ   string
		name  string
		value string
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name        string
		contentType string
		metric      metric
		want        want
	}{
		{
			name: "valid counter",
			metric: metric{
				typ:   "counter",
				name:  "some",
				value: "100",
			},
			want: want{
				statusCode: http.StatusOK},
		},
		{
			name: "valid gauge",
			metric: metric{
				typ:   "gauge",
				name:  "some",
				value: "100",
			},
			want: want{
				statusCode: http.StatusOK},
		},
		{
			name: "invalid type",
			metric: metric{
				typ:   "foo",
				name:  "some",
				value: "100",
			},
			want: want{
				statusCode: http.StatusBadRequest},
		},
		{
			name: "empty name",
			metric: metric{
				typ:   "gauge",
				name:  "",
				value: "100",
			},
			want: want{
				statusCode: http.StatusNotFound},
		},
		{
			name: "empty value",
			metric: metric{
				typ:   "gauge",
				name:  "some",
				value: "",
			},
			want: want{
				statusCode: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/%s/%s", tt.metric.typ, tt.metric.name, tt.metric.value)
			request := httptest.NewRequest(http.MethodPost, url, nil)
			request.Header.Add("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			handler.UpdateMetrics(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
