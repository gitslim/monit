package sender

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

// compressGzip - Сжатие данных методом gzip
func compressGzip(data []byte, level int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	gzWriter, err := gzip.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %v", err)
	}
	defer gzWriter.Close()

	_, err = gzWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %v", err)
	}

	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}
	return &buf, nil
}
