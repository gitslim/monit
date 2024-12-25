package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
	r.POST("/update/", metricHandler.UpdateMetric)

	type metric struct {
		typ   string
		name  string
		param string
		val   string
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name: "valid counter",
			metric: metric{
				typ:   "counter",
				name:  "c1",
				param: "delta",
				val:   "100",
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
				param: "value",
				val:   "100",
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
				param: "value",
				val:   "100",
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
				param: "value",
				val:   "100",
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
				param: "value",
				val:   "abc",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+":plain", func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/%s/%s", tt.metric.typ, tt.metric.name, tt.metric.val)
			req, err := http.NewRequest(http.MethodPost, url, nil)
			assert.NoError(t, err)

			req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypePlain)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()
			res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
		t.Run(tt.name+":json", func(t *testing.T) {
			url := "/update/"
			jsonData := fmt.Sprintf("{\"id\":\"%s\",\"type\":\"%s\",\"%s\":%s}", tt.metric.name, tt.metric.typ, tt.metric.param, tt.metric.val)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonData))
			assert.NoError(t, err)

			req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()
			res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}

func TestBatchUpdateMetrics(t *testing.T) {
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
	r.POST("/updates/", metricHandler.BatchUpdateMetrics)

	type metric struct {
		typ   string
		name  string
		param string
		val   string
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name: "valid counter",
			metric: metric{
				typ:   "counter",
				name:  "c1",
				param: "delta",
				val:   "100",
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
				param: "value",
				val:   "100",
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
				param: "value",
				val:   "100",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	t.Run("BatchUpdate", func(t *testing.T) {
		url := "/updates/"
		jsonData := "["
		for i, tt := range tests {
			jsonData = jsonData + fmt.Sprintf("{\"id\":\"%s\",\"type\":\"%s\",\"%s\":%s}", tt.metric.name, tt.metric.typ, tt.metric.param, tt.metric.val)
			if i != len(tests)-1 {
				jsonData += ","
			}
		}
		jsonData += "]"
		fmt.Println(jsonData)

		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonData))
		assert.NoError(t, err)

		req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
