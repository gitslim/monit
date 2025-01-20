package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/handlers"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/gitslim/monit/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

// TestMemStorageNonBatchUpdate тестирует обновление метрик в хранилище в памяти с одиночной отправкой
func TestMemStorageNonBatchUpdate(t *testing.T) {
	r, teardown, err := testhelpers.CreateServerMock(false)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()
	updateMetricNonBatch(t, r)
}

// TestMemStorageBatchUpdate тестирует обновление метрик в хранилище в памяти с пакетной отправкой
func TestMemStorageBatchUpdate(t *testing.T) {
	r, teardown, err := testhelpers.CreateServerMock(false)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()
	updateMetricsBatch(t, r)
}

// TestPGStorageNonBatchUpdate тестирует обновление метрик в хранилище в базе данных с одиночной отправкой
func TestPGStorageNonBatchUpdate(t *testing.T) {
	r, teardown, err := testhelpers.CreateServerMock(false)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()
	updateMetricNonBatch(t, r)
}

// TestPGStorageBatchUpdate тестирует обновление метрик в хранилище в базе данных с пакетной отправкой
func TestPGStorageBatchUpdate(t *testing.T) {
	r, teardown, err := testhelpers.CreateServerMock(false)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()
	updateMetricsBatch(t, r)
}

// updateMetricNonBatch тестирует обновление метрик по одной.
func updateMetricNonBatch(t *testing.T, r *gin.Engine) {
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
			req, err2 := http.NewRequest(http.MethodPost, url, nil)
			assert.NoError(t, err2)

			req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypePlain)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()
			err := res.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
		t.Run(tt.name+":json", func(t *testing.T) {
			url := "/update/"
			jsonData := fmt.Sprintf("{\"id\":\"%s\",\"type\":\"%s\",\"%s\":%s}", tt.metric.name, tt.metric.typ, tt.metric.param, tt.metric.val)

			req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonData))
			assert.NoError(t, err)

			req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()
			err = res.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}

// updateMetricsBatch тестирует пакетное обновление, получение значений метрик и их список.
func updateMetricsBatch(t *testing.T, r *gin.Engine) {
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
	t.Run("Update", func(t *testing.T) {
		url := "/updates/"
		jsonData := "["
		for i, tt := range tests {
			jsonData = jsonData + fmt.Sprintf("{\"id\":\"%s\",\"type\":\"%s\",\"%s\":%s}", tt.metric.name, tt.metric.typ, tt.metric.param, tt.metric.val)
			if i != len(tests)-1 {
				jsonData += ","
			}
		}
		jsonData += "]"

		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonData))
		assert.NoError(t, err)

		req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		err = res.Body.Close()
		assert.NoError(t, err)
	})
	t.Run("Get", func(t *testing.T) {
		url := "/value/"
		tt := tests[0]
		jsonData := fmt.Sprintf("{\"id\":\"%s\",\"type\":\"%s\"}", tt.metric.name, tt.metric.typ)

		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonData))
		assert.NoError(t, err)

		req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var dto *entities.MetricDTO
		err = json.NewDecoder(res.Body).Decode(&dto)
		assert.NoError(t, err)
		err = res.Body.Close()
		assert.NoError(t, err)

		switch tt.metric.typ {
		case "counter":
			assert.Equal(t, tt.metric.val, strconv.FormatInt(*dto.Delta, 10))
		case "gauge":
			assert.Equal(t, tt.metric.val, strconv.FormatFloat(*dto.Value, 'f', -1, 64))
		}
	})
	t.Run("List", func(t *testing.T) {
		url := "/"
		tt := tests[0]
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypePlain)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		err = res.Body.Close()
		assert.NoError(t, err)

		assert.Contains(t, string(body), tt.metric.name)
	})
}

// ExampleMetricHandler_ListMetrics пример получения html страницы со списком метрик.
func ExampleMetricHandler_ListMetrics() {
	// создаем сервер
	r, teardown, err := testhelpers.CreateServerMock(false)
	if err != nil {
		panic(err)
	}
	defer teardown()

	// посылаем запрос на обновление метрики
	w := httptest.NewRecorder()
	payload := `{"id":"Alloc","type":"gauge","value":12345.67}`
	req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(payload))
	req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
	r.ServeHTTP(w, req)

	// получаем страницу со списком метрик
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	// проверяем, что метрика есть на странице
	fmt.Println(strings.Contains(w.Body.String(), "12345.67"))

	// Output: true
}

// ExampleMetricHandler_UpdateMetric - пример обновления метрики
func ExampleMetricHandler_UpdateMetric() {
	var _ handlers.MetricHandler // for import

	// создаем сервер
	r, teardown, err := testhelpers.CreateServerMock(false)
	if err != nil {
		panic(err)
	}
	defer teardown()

	// отправляем запрос на обновление метрики
	w := httptest.NewRecorder()
	payload := `{"id":"Alloc","type":"gauge","value":12345.67}`
	req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(payload))
	req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
	r.ServeHTTP(w, req)

	// проверяем, что код ответа равен 200
	fmt.Println(w.Code == http.StatusOK)

	// Output: true
}

// ExampleMetricHandler_GetMetric пример получения метртики
func ExampleMetricHandler_GetMetric() {
	// создаем сервер
	r, teardown, err := testhelpers.CreateServerMock(false)
	if err != nil {
		panic(err)
	}
	defer teardown()

	// отправляем запрос на обновление метрики
	w := httptest.NewRecorder()
	payload := `{"id":"Alloc","type":"gauge","value":12345.67}`
	req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(payload))
	req.Header.Add(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
	r.ServeHTTP(w, req)

	// отправляем запрос на получение метрики
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/value/gauge/Alloc", nil)
	r.ServeHTTP(w, req)

	// выводим результат
	fmt.Println(w.Body.String())

	// Output: 12345.67
}
