package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func isCompressionAcceptable(c *gin.Context) bool {
	return strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") // TODO: better check
}

func isContentTypeCompressable(c *gin.Context) bool {
	supportedContentTypes := []string{"application/json", "text/html"}
	ct := c.GetHeader("Content-Type")
	for _, v := range supportedContentTypes {
		if strings.Contains(v, ct) {
			return true
		}
	}
	return false
}

func isRequestCompressed(c *gin.Context) bool {
	return c.GetHeader("Content-Encoding") == "gzip"
}

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isRequestCompressed(c) {
			gzReader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			defer gzReader.Close()

			c.Request.Body = io.NopCloser(gzReader)
		}

		if isCompressionAcceptable(c) && isContentTypeCompressable(c) {
			// gzWriter := gzip.NewWriter(c.Writer)
			gzWriter, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed) // TODO: make shared writer with Reset()
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer gzWriter.Close()

			c.Header("Content-Encoding", "gzip")

			c.Writer = &gzipResponseWriter{Writer: gzWriter, ResponseWriter: c.Writer}

		}
		c.Next()
	}
}
