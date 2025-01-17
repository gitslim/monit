package sender_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/sender"
	"github.com/gitslim/monit/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMetrics(t *testing.T) {
	// запуск мок-сервера
	engine, err := testhelpers.CreateServerMock()
	require.NoError(t, err)
	srv := testhelpers.StartServerMock(engine)
	// defer testhelpers.StopServerMock(srv)

	// создание конфига
	cfg := &conf.Config{
		Addr: srv.Addr,
	}

	// создание клиента
	client := &http.Client{}

	sendFunc := func(data string, batch bool) error {
		metrics, err := testhelpers.JsonToMetricDTO(data)
		assert.NoError(t, err)
		return sender.SendMetrics(context.Background(), cfg, client, metrics, batch)
	}

	tests := []struct {
		name     string
		jsonData string
		batch    bool
		key      string
		wantErr  bool
	}{
		{
			jsonData: `[{"id": "test_gauge", "type": "gauge", "value": 100.0}]`,
			batch:    true,
			key:      "",
			wantErr:  false,
		},
		{
			jsonData: `[{"id": "test_gauge", "type": "gauge", "value": 100.0}]`,
			batch:    false,
			key:      "",
			wantErr:  false,
		},
		{
			jsonData: `[{"id": "test_counter", "type": "counter", "value": 10}]`,
			batch:    true,
			key:      "some-key",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := sendFunc(tt.jsonData, tt.batch)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SendMetrics() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SendMetrics() succeeded unexpectedly")
			}
		})
	}
}
