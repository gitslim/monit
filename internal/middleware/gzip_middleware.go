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

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func isCompressionAcceptable(c *gin.Context) bool {
	return strings.Contains(c.GetHeader(httpconst.HeaderAcceptEncoding), httpconst.AcceptEncodingGzip) // TODO: better check
}

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

func isRequestCompressed(c *gin.Context) bool {
	return c.GetHeader(httpconst.HeaderContentEncoding) == httpconst.ContentEncodingGzip
}

func GzipMiddleware() gin.HandlerFunc {
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
	return func(c *gin.Context) {
		if isRequestCompressed(c) {
			gzReader := gzipReaderPool.Get().(*gzip.Reader)
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
