package sender

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

// compressGzip сжимает данные в gzip.
func compressGzip(data []byte, level int) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	w, err := gzip.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %v", err)
	}
	defer func() {
		err2 := w.Close()
		if err2 != nil {
			fmt.Printf("failed to close gzip writer: %v\n", err2)
		}
	}()

	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %v", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}
	return &buf, nil
}
