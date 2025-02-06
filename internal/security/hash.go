package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// MakeSHA256Hash расчитывает хэш по алгоритму SHA256.
func MakeSHA256Hash(rc io.Reader, key string) (string, error) {
	var buf []byte
	_, err := rc.Read(buf)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil)), nil
}
