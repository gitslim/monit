package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/httpconst"
)

// gzipResponseWriter представляет ответ, который будет сжат gzip-компрессией
type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

// Write реализует интерфейс io.Writer
func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

// isCompressionAcceptable возвращает true, если клиент может принимать сжатый ответ
func isCompressionAcceptable(c *gin.Context) bool {
	return strings.Contains(c.GetHeader(httpconst.HeaderAcceptEncoding), httpconst.AcceptEncodingGzip) // TODO: better check
}

// isContentTypeCompressable возвращает true, если тип содержимого может быть сжат
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

// isRequestCompressed возвращает true, если запрос сжат
func isRequestCompressed(c *gin.Context) bool {
	return c.GetHeader(httpconst.HeaderContentEncoding) == httpconst.ContentEncodingGzip
}

// GzipMiddleware возвращает функцию-мидлварь для сжатия ответов gzip
func GzipMiddleware() gin.HandlerFunc {
	// Создаем пулы для ускорения работы
	gzipWriterPool := sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}

	gzipReaderPool := sync.Pool{
		New: func() interface{} {
			r, _ := gzip.NewReader(nil)
			return r
		},
	}

	var gzipResponseWriterPool = sync.Pool{
		New: func() interface{} {
			return &gzipResponseWriter{}
		},
	}
	// возвращаем функцию-мидлварь
	return func(c *gin.Context) {
		if isRequestCompressed(c) {
			fmt.Printf("c: %+v\n", c)
			gzReader := gzipReaderPool.Get().(*gzip.Reader)
			fmt.Printf("GZREADER: %v\n", gzReader)
			err := gzReader.Reset(c.Request.Body)
			if err != nil {
				fmt.Printf("Gzreader error: %v\n", err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			defer func() {
				gzReader.Close()
				gzipReaderPool.Put(gzReader)
			}()

			c.Request.Body = io.NopCloser(gzReader)
		}

		if isCompressionAcceptable(c) && isContentTypeCompressable(c) {
			gzWriter := gzipWriterPool.Get().(*gzip.Writer)
			gzWriter.Reset(c.Writer)
			defer func() {
				gzWriter.Close()
				gzipWriterPool.Put(gzWriter)
			}()

			c.Header(httpconst.HeaderContentEncoding, httpconst.ContentEncodingGzip)

			gzResponseWriter := gzipResponseWriterPool.Get().(*gzipResponseWriter)
			gzResponseWriter.Writer = gzWriter
			gzResponseWriter.ResponseWriter = c.Writer
			c.Writer = gzResponseWriter
			defer func() {
				gzipResponseWriterPool.Put(gzResponseWriter)
			}()
		}
		c.Next()
	}
}
