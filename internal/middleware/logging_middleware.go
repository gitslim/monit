package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		// логируем запрос
		sugar.Infow("req",
			"URI", c.Request.RequestURI,
			"method", c.Request.Method,
			"latency", latency)

		status := c.Writer.Status()
		size := c.Writer.Size()
		// логируем ответ
		sugar.Infow("resp",
			"status", status,
			"size", size)
	}
}
