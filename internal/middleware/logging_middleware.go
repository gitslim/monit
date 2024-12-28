package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/logging"
)

// LoggerMiddleware логгирует запросы и ответы, вычисляет время выполнения запроса и размер ответа.
func LoggerMiddleware(log *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		// Логгируем запрос.
		log.Info("req",
			"URI", c.Request.RequestURI,
			"method", c.Request.Method,
			"latency", latency)

		status := c.Writer.Status()
		size := c.Writer.Size()
		// Логгируем ответ.
		log.Info("res",
			"status", status,
			"size", size)
	}
}
