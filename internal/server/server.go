package server

import (
	"log"
	"net/http"

	"github.com/gitslim/monit/internal/handlers"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("Server is starting at %s...\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}

// Инициализация сервера с нужными обработчиками
func New(addr string, metricsHandler *handlers.MetricsHandler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}", metricsHandler.UpdateMetrics)

	return NewServer(addr, mux)
}
