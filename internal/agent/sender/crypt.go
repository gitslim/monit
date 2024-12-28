package sender

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// makeHash расчитывает хэш по алгоритму SHA256.
func makeHash(rc io.ReadCloser, key string) (string, error) {
	var buf []byte
	_, err := rc.Read(buf)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil)), nil
}
