package middleware

import (
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/httpconst"
)

// isCompressionAcceptable возвращает true, если клиент может принимать сжатый ответ.
func isCompressionAcceptable(c *gin.Context) bool {
	return strings.Contains(c.GetHeader(httpconst.HeaderAcceptEncoding), httpconst.AcceptEncodingGzip) // TODO: better check
}

// isContentTypeCompressable возвращает true, если тип содержимого может быть сжат.
func isContentTypeCompressable(c *gin.Context) bool {
	supportedContentTypes := []string{httpconst.ContentTypeJSON, httpconst.ContentTypeHTML}
	ct := c.GetHeader(httpconst.HeaderContentType)
	for _, v := range supportedContentTypes {
		if strings.Contains(v, ct) {
			return true
		}
	}
	return false
}

// isRequestCompressed возвращает true, если запрос сжат.
func isRequestCompressed(c *gin.Context) bool {
	return c.GetHeader(httpconst.HeaderContentEncoding) == httpconst.ContentEncodingGzip
}

// GzipMiddleware прозрачно управляет gzip-сжатием в зависимости от поддержки клиентом.
func GzipMiddleware() gin.HandlerFunc {
	// Создаем пул для ускорения работы gzip и уменьшения расхода памяти.
	pool := NewGzipPool()

	// Возвращаем функцию-мидлварь.
	return func(c *gin.Context) {
		if isRequestCompressed(c) {
			gzReader := pool.GetReader(c.Request.Body)
			defer func() {
				pool.PutReader(gzReader)
			}()
			c.Request.Body = io.NopCloser(gzReader)
		}

		if isCompressionAcceptable(c) && isContentTypeCompressable(c) {
			gzWriter := pool.GetWriter(c.Writer)
			defer func() {
				pool.PutWriter(gzWriter)
			}()

			c.Header(httpconst.HeaderContentEncoding, httpconst.ContentEncodingGzip)

			gzResponseWriter := pool.GetResponseWriter(gzWriter, c.Writer)
			c.Writer = gzResponseWriter
		}
		c.Next()
	}
}
