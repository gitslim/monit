package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	log, err := logging.NewLogger()
	if err != nil {
		log.Fatal("Failed init logger")
	}
	cfg := &conf.Config{
		StoreInterval:   0,
		FileStoragePath: "/tmp/.monit/memstorage.json",
		Restore:         false,
	}

	conf, err := services.WithMemStorage(context.Background(), log, cfg, make(chan<- error))
	assert.NoError(t, err)

	metricService, err := services.NewMetricService(conf)
	assert.NoError(t, err)

	metricHandler := NewMetricHandler(metricService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/update/:type/:name/:value", metricHandler.UpdateMetric)

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
				name:  "c1",
				value: "100",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "valid gauge",
			metric: metric{
				typ:   "gauge",
				name:  "g1",
				value: "100",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "invalid type",
			metric: metric{
				typ:   "foo",
				name:  "some",
				value: "100",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "empty name",
			metric: metric{
				typ:   "gauge",
				name:  "",
				value: "100",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "non digital value",
			metric: metric{
				typ:   "gauge",
				name:  "some",
				value: "abc",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/%s/%s", tt.metric.typ, tt.metric.name, tt.metric.value)
			req, err := http.NewRequest(http.MethodPost, url, nil)
			assert.NoError(t, err)

			req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypePlain)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()
			res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}
