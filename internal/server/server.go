package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/handlers"
)

type Server struct {
	Addr   string
	Engine *gin.Engine
}

func NewServer(addr string, engine *gin.Engine) *Server {
	return &Server{
		Addr:   addr,
		Engine: engine,
	}
}

func (s *Server) Start() error {
	log.Printf("Server is starting at %s...\n", s.Addr)
	return s.Engine.Run(s.Addr)
}

// Инициализация сервера с нужными обработчиками
func New(addr string, metricsHandler *handlers.MetricsHandler) *Server {
	r := gin.Default()
	r.POST("/update/:type/:name/:value", metricsHandler.UpdateMetrics)

	return NewServer(addr, r)
}
