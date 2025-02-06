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
	"github.com/gitslim/monit/internal/security"
)

// encryptData шифрует данные перед отправкой.
func encryptData(cfg *conf.Config, data []byte) ([]byte, error) {
	// Если задан ключ шифрования, шифруем данные.
	if cfg.CryptoKey != "" {
		pubKey, err := security.ReadRSAPublicKeyFromFile(cfg.CryptoKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read public key: %v", err)
		}
		data, err = security.EncryptRSA(pubKey, data)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt body: %v", err)
		}
	}
	return data, nil
}

// signRequest подписывает запрос перед отправкой.
func signRequest(req *http.Request, cfg *conf.Config) error {
	// Если задан ключ, подписываем данные.
	if cfg.Key != "" {
		body, err := req.GetBody()
		if err != nil {
			return err
		}
		defer func() { _ = body.Close() }()

		// Подписываем данные.
		hash, err := security.MakeSHA256Hash(body, cfg.Key)
		if err != nil {
			return err
		}

		// Записываем хэш в заголовок.
		req.Header.Set(httpconst.HeaderHashSHA256, hash)
	}
	return nil
}

// sendJSON отправляет метрики в формате JSON батчем или по одной.
func sendJSON(ctx context.Context, cfg *conf.Config, client *http.Client, url string, jsonData []byte) error {
	// Шифруем данные если необходимо.
	jsonData, err := encryptData(cfg, jsonData)
	if err != nil {
		return fmt.Errorf("failed to encrypt body: %v", err)
	}

	// Сжимаем данные в gzip
	buf, err := compressGzip(jsonData, gzip.BestSpeed)
	if err != nil {
		return fmt.Errorf("failed to compress with gzip: %v", err)
	}

	// Таймаут запроса.
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
	req.Header.Set(httpconst.HeaderContentEncoding, httpconst.ContentEncodingGzip)

	// Подписываем запрос если необходимо.
	err = signRequest(req, cfg)
	if err != nil {
		return fmt.Errorf("failed to sign request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// SendMetrics отправляет метрики на сервер в формате JSON батчем или по одной.
func SendMetrics(ctx context.Context, cfg *conf.Config, client *http.Client, metrics []*entities.MetricDTO, batch bool) error {
	serverURL := fmt.Sprintf("http://%s", cfg.Addr)

	// Ретраи при сбое.
	return retry.Retry(func() error {
		var url string
		var jsonData []byte
		var err error

		if batch {
			// Отправляем батч метрик.
			url = fmt.Sprintf("%s/updates/", serverURL)
			jsonData, err = json.Marshal(metrics)
			if err != nil {
				return err
			}
			return sendJSON(ctx, cfg, client, url, jsonData)
		} else {
			// Отправляем метрики по одной.
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
