package worker

import (
	"net/http"
	"sync"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/entities"
)

// WorkerPool - Пул worker'ов
type WorkerPool struct {
	Metrics chan entities.MetricDTO
	WG      *sync.WaitGroup
	Cfg     *conf.Config
	Client  *http.Client
}

// Start - Запуск пула worker'ов
func (w *WorkerPool) Start(f func()) {
	for i := 0; i < int(w.Cfg.RateLimit); i++ {
		w.WG.Add(1)
		go f()
	}
}

// Wait - Ожидание завершения всех worker'ов
func (w *WorkerPool) Wait() {
	w.WG.Wait()
}

// NewWorkerPool - Создание пула worker'ов
func NewWorkerPool(cfg *conf.Config) *WorkerPool {
	return &WorkerPool{
		Metrics: make(chan entities.MetricDTO, cfg.RateLimit),
		WG:      &sync.WaitGroup{},
		Cfg:     cfg,
		Client:  &http.Client{},
	}
}
