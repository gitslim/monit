package sender

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/gitslim/monit/internal/retry"
)

// sendJSON - Отправка метрик в формате JSON батчем или по одной
func sendJSON(ctx context.Context, cfg *conf.Config, client *http.Client, url string, jsonData []byte) error {
	buf, err := compressGzip(jsonData, gzip.BestSpeed)
	if err != nil {
		return fmt.Errorf("failed to compress with gzip: %v", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
	req.Header.Set(httpconst.HeaderContentEncoding, httpconst.ContentEncodingGzip)

	if cfg.Key != "" {
		body, err := req.GetBody()
		if err != nil {
			return err
		}
		hash, err := makeHash(body, cfg.Key)
		if err != nil {
			return err
		}
		req.Header.Set(httpconst.HeaderHashSHA256, hash)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// SendMetrics - Функция для отправки всех метрик на сервер
func SendMetrics(ctx context.Context, cfg *conf.Config, client *http.Client, metrics []*entities.MetricDTO, batch bool) error {
	serverURL := fmt.Sprintf("http://%s", cfg.Addr)

	return retry.Retry(func() error {
		var url string
		var jsonData []byte
		var err error

		if batch {
			// Отправляем батч метрик
			url = fmt.Sprintf("%s/updates/", serverURL)
			jsonData, err = json.Marshal(metrics)
			if err != nil {
				return err
			}
			return sendJSON(ctx, cfg, client, url, jsonData)
		} else {
			// Отправляем метрики по одной
			url = fmt.Sprintf("%s/update/", serverURL)

			for _, metric := range metrics {
				jsonData, err = json.Marshal(&metric)
				if err != nil {
					return err
				}
				err := sendJSON(ctx, cfg, client, url, jsonData)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}, 3)
}
