package testhelpers

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/server/engine"
	"github.com/gitslim/monit/internal/services"
)

// CreateServerMock создает мок-сервер для тестирования.
func CreateServerMock() (*gin.Engine, error) {
	log, err := logging.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v\n", err))
	}
	cfg := &conf.Config{
		StoreInterval:   0,
		FileStoragePath: "/tmp/.monit/memstorage.json",
		Restore:         false,
	}

	conf, err := services.WithMemStorage(context.Background(), log, cfg, make(chan<- error))
	if err != nil {
		return nil, err
	}

	metricService, err := services.NewMetricService(conf)
	if err != nil {
		return nil, err
	}

	return engine.CreateGinEngine(cfg, log, gin.ReleaseMode, metricService)
}

func createListener() (l net.Listener, close func()) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	return l, func() {
		_ = l.Close()
	}
}

func StartServerMock(r *gin.Engine) *http.Server {
	// Создаем сервер на свободном порту.
	l, _ := createListener()

	srv := &http.Server{
		Addr:    l.Addr().String(),
		Handler: r,
	}

	go func() {
		err := srv.Serve(l)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %w", err))
		}
	}()

	return srv
}

func StopServerMock(srv *http.Server) {
	// Завершаем работу сервера.
	err := srv.Shutdown(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to shutdown server: %w", err))
	}
}
