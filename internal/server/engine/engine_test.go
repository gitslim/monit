// Package engine_test содержит тесты для пакета engine.
package engine_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/server/engine"
	"github.com/gitslim/monit/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateGinEngine(t *testing.T) {
	log, err := logging.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v\n", err))
	}
	cfg := &conf.Config{
		StoreInterval:   0,
		FileStoragePath: "/tmp/.monit/memstorage.json",
		Restore:         false,
	}

	svcCfg, err := services.WithMemStorage(context.Background(), log, cfg, make(chan<- error))
	assert.NoError(t, err)

	metricService, err := services.NewMetricService(svcCfg)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		cfg           *conf.Config
		log           *logging.Logger
		ginMode       string
		metricService *services.MetricService
		wantErr       bool
	}{
		{
			name:          "success",
			cfg:           cfg,
			log:           log,
			ginMode:       gin.DebugMode,
			metricService: metricService,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := engine.CreateGinEngine(tt.cfg, tt.log, tt.ginMode, tt.metricService)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateGinEngine() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateGinEngine() succeeded unexpectedly")
			}

			if got == nil {
				t.Error("CreateGinEngine() = nil")
			}
		})
	}
}
