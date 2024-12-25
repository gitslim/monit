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
		go func() {
			defer w.WG.Done()
			f()
		}()
	}
}

// AddWorker - Добавление worker'а в пул
func (w *WorkerPool) AddWorker(f func()) {
	w.WG.Add(1)
	go func() {
		defer w.WG.Done()
		f()
	}()
}

// WaitClose - Ожидание завершения всех worker'ов и закрытие канала
func (w *WorkerPool) WaitClose() {
	w.WG.Wait()
	close(w.Metrics)
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
