// Модуль worker управляет пулом worker'ов.
package worker

import (
	"net/http"
	"sync"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/entities"
)

// WorkerPool определяет пул worker'ов.
type WorkerPool struct {
	Metrics chan entities.MetricDTO
	WG      *sync.WaitGroup
	Cfg     *conf.Config
	Client  *http.Client
}

// Start запускает пул worker'ов.
func (w *WorkerPool) Start(f func()) {
	for i := 0; i < int(w.Cfg.RateLimit); i++ {
		w.WG.Add(1)
		go func() {
			defer w.WG.Done()
			f()
		}()
	}
}

// AddWorker добавляет worker'а в пул.
func (w *WorkerPool) AddWorker(f func()) {
	w.WG.Add(1)
	go func() {
		defer w.WG.Done()
		f()
	}()
}

// WaitClose ожидает завершения всех worker'ов и закрывает канал.
func (w *WorkerPool) WaitClose() {
	w.WG.Wait()
	close(w.Metrics)
}

// NewWorkerPool создает пул worker'ов.
func NewWorkerPool(cfg *conf.Config) *WorkerPool {
	return &WorkerPool{
		Metrics: make(chan entities.MetricDTO, cfg.RateLimit),
		WG:      &sync.WaitGroup{},
		Cfg:     cfg,
		Client:  &http.Client{},
	}
}
