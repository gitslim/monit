package middleware

import (
	"compress/gzip"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

// GzipResponseWriter представляет ответ, который будет сжат gzip-компрессией.
type GzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

// Write реализует интерфейс io.Writer.
func (w *GzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

// GzipPool представляет пул для gzip-компрессии.
type GzipPool struct {
	readers         sync.Pool
	writers         sync.Pool
	responseWriters sync.Pool
}

// NewGzipPool создает пул для gzip-компрессии.
func NewGzipPool() *GzipPool {
	return &GzipPool{
		readers:         sync.Pool{},
		writers:         sync.Pool{},
		responseWriters: sync.Pool{},
	}
}

// GetReader возвращает reader для gzip-компрессии из пула.
func (pool *GzipPool) GetReader(src io.Reader) (reader *gzip.Reader) {
	if r := pool.readers.Get(); r != nil {
		reader = r.(*gzip.Reader)
		err := reader.Reset(src)
		if err != nil {
			reader, _ = gzip.NewReader(src)
		}
	} else {
		reader, _ = gzip.NewReader(src)
	}
	return reader
}

// PutReader возвращает reader для gzip-компрессии в пул.
func (pool *GzipPool) PutReader(reader *gzip.Reader) {
	err := reader.Close()
	if err != nil {
		return
	}
	pool.readers.Put(reader)
}

// GetWriter возвращает writer для gzip-компрессии из пула.
func (pool *GzipPool) GetWriter(dst io.Writer) (writer *gzip.Writer) {
	if w := pool.writers.Get(); w != nil {
		writer = w.(*gzip.Writer)
		writer.Reset(dst)
	} else {
		writer, _ = gzip.NewWriterLevel(dst, gzip.BestCompression)
	}
	return writer
}

// PutWriter возвращает writer для gzip-компрессии в пул.
func (pool *GzipPool) PutWriter(writer *gzip.Writer) {
	err := writer.Close()
	if err != nil {
		return
	}
	pool.writers.Put(writer)
}

// GetResponseWriter возвращает response writer для gzip-компрессии из пула.
func (pool *GzipPool) GetResponseWriter(gzipWriter *gzip.Writer, responseWriter gin.ResponseWriter) (writer *GzipResponseWriter) {
	if grw := pool.responseWriters.Get(); grw != nil {
		writer = grw.(*GzipResponseWriter)
		writer.ResponseWriter = responseWriter
	} else {
		writer = &GzipResponseWriter{Writer: gzipWriter, ResponseWriter: responseWriter}
	}
	return writer
}

// PutResponseWriter возвращает response writer для gzip-компрессии в пул.
func (pool *GzipPool) PutResponseWriter(gzipResponseWriter *GzipResponseWriter) {
	pool.responseWriters.Put(gzipResponseWriter)
}
