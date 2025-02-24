package worker

import (
	"context"
	"net/http"
	"sync"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/sender"
	"github.com/gitslim/monit/internal/entities"
	pb "github.com/gitslim/monit/pkg/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// WorkerPool определяет пул worker'ов.
type WorkerPool struct {
	Metrics    chan entities.MetricDTO
	WG         *sync.WaitGroup
	Cfg        *conf.Config
	HTTPClient *http.Client
	GRPCClient *pb.MonitClient
	once       sync.Once // Для безопасного закрытия канала Metrics
}

// Start запускает пул worker'ов с поддержкой контекста.
func (w *WorkerPool) Start(ctx context.Context, f func(ctx context.Context)) {
	for i := 0; i < int(w.Cfg.RateLimit); i++ {
		w.WG.Add(1)
		go func() {
			defer w.WG.Done()
			f(ctx) // Передаём контекст в worker
		}()
	}
}

// AddWorker добавляет worker'а в пул с поддержкой контекста.
func (w *WorkerPool) AddWorker(ctx context.Context, f func(ctx context.Context)) {
	w.WG.Add(1)
	go func() {
		defer w.WG.Done()
		f(ctx)
	}()
}

// WaitClose ожидает завершения всех worker'ов и закрывает канал метрик.
func (w *WorkerPool) WaitClose() {
	// Ожидаем завершения всех worker'ов
	w.WG.Wait()

	// Закрываем канал метрик
	w.once.Do(func() {
		close(w.Metrics)
	})
}

// NewWorkerPool создает пул worker'ов.
func NewWorkerPool(cfg *conf.Config) (*WorkerPool, error) {
	var httpClient http.Client
	var grpcClient pb.MonitClient

	if cfg.GRPC {
		conn, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		grpcClient = pb.NewMonitClient(conn)
	} else {
		httpClient = http.Client{
			Transport: sender.NewCustomTransport(), // Используем CustomTransport
		}
	}
	return &WorkerPool{
		Metrics:    make(chan entities.MetricDTO, cfg.RateLimit),
		WG:         &sync.WaitGroup{},
		Cfg:        cfg,
		HTTPClient: &httpClient,
		GRPCClient: &grpcClient,
	}, nil
}
