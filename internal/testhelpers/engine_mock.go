package testhelpers

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/server/engine"
	"github.com/gitslim/monit/internal/services"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// CreateServerMock создает мок-сервер для тестирования.
func CreateServerMock(withPostgres bool) (*gin.Engine, func(), error) {
	ctx := context.Background()
	teardown := func() {}

	log, err := logging.NewLogger()
	if err != nil {
		return nil, nil, err
	}

	var cfg *conf.Config
	var svcConf services.MetricServiceConf

	if withPostgres {
		pgContainer, err2 := postgres.Run(ctx,
			"postgres:17-alpine",
			postgres.WithDatabase("test-db"),
			postgres.WithUsername("postgres"),
			postgres.WithPassword("postgres"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		)
		if err2 != nil {
			return nil, nil, err2
		}

		teardown = func() {
			if err2 = pgContainer.Terminate(ctx); err2 != nil {
				panic("failed to terminate pgContainer")
			}
		}

		dsn, err2 := pgContainer.ConnectionString(ctx, "sslmode=disable")
		if err2 != nil {
			return nil, nil, err2
		}

		cfg = &conf.Config{
			DatabaseDSN: dsn,
		}
		svcConf, err2 = services.WithPGStorage(ctx, log, cfg)
		if err2 != nil {
			return nil, nil, err2
		}

	} else {
		cfg = &conf.Config{
			StoreInterval:   0,
			FileStoragePath: "/tmp/.monit/memstorage.json",
			Restore:         false,
		}
		svcConf, err = services.WithMemStorage(ctx, log, cfg, make(chan<- error))
		if err != nil {
			return nil, nil, err
		}
	}

	metricService, err := services.NewMetricService(svcConf)
	if err != nil {
		return nil, nil, err
	}

	engine, err := engine.CreateGinEngine(cfg, log, gin.ReleaseMode, metricService)
	if err != nil {
		return nil, nil, err
	}
	return engine, teardown, nil
}

func StartServerMock(r *gin.Engine) (*http.Server, func(), error) {
	// Создаем листенер на свободном порту.
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil,nil, err
	}

	// создаем сервер
	srv := &http.Server{
		Addr:    l.Addr().String(),
		Handler: r,
	}

	// запускаем сервер
	go func() {
		_ = srv.Serve(l)
	}()

	teardown := func() {
		_ = srv.Close()
		_ = l.Close()
	}

	return srv, teardown, nil
}
