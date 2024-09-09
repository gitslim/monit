package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := NewMetricHandler(memStorage)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/update/:type/:name/:value", metricsHandler.UpdateMetric)

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
			name: "non digital value",
			metric: metric{
				typ:   "gauge",
				name:  "some",
				value: "abc",
			},
			want: want{
				statusCode: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/%s/%s", tt.metric.typ, tt.metric.name, tt.metric.value)
			req, err := http.NewRequest(http.MethodPost, url, nil)
			assert.NoError(t, err)

			req.Header.Add("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()
			res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}
