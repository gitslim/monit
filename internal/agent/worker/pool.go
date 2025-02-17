// Package worker управляет пулом worker'ов.
package worker

import (
	"context"
	"net/http"
	"sync"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/agent/transport"
	"github.com/gitslim/monit/internal/entities"
)

// WorkerPool определяет пул worker'ов.
type WorkerPool struct {
	Metrics chan entities.MetricDTO
	WG      *sync.WaitGroup
	Cfg     *conf.Config
	Client  *http.Client
	once    sync.Once // Для безопасного закрытия канала Metrics
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

// Stop останавливает пул worker'ов.
func (w *WorkerPool) Stop() {
	w.once.Do(func() {
		close(w.Metrics) // Закрываем канал метрик
	})
}

// Wait ожидает завершения всех worker'ов.
func (w *WorkerPool) Wait() {
	w.WG.Wait()
}

// NewWorkerPool создает пул worker'ов.
func NewWorkerPool(cfg *conf.Config) *WorkerPool {
	return &WorkerPool{
		Metrics: make(chan entities.MetricDTO, cfg.RateLimit),
		WG:      &sync.WaitGroup{},
		Cfg:     cfg,
		Client: &http.Client{
			Transport: transport.NewCustomTransport(), // Используем CustomTransport
		},
	}
}
