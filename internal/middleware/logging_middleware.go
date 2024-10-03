package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/logging"
)

func LoggerMiddleware(log *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		// логируем запрос
		log.Info("req",
			"URI", c.Request.RequestURI,
			"method", c.Request.Method,
			"latency", latency)

		status := c.Writer.Status()
		size := c.Writer.Size()
		// логируем ответ
		log.Info("resp",
			"status", status,
			"size", size)
	}
}
