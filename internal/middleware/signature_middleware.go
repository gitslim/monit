package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/gitslim/monit/internal/logging"
)

// signatureResponseWriter представляет собой обертку над gin.ResponseWriter, которая добавляет хэш SHA-256 в заголовок ответа.
type signatureResponseWriter struct {
	gin.ResponseWriter
	key string
}

// Write реализует интерфейс gin.ResponseWriter.
func (w *signatureResponseWriter) Write(data []byte) (int, error) {
	respHash, err := makeHash(&data, w.key)
	if err != nil {
		return 0, err
	}
	w.ResponseWriter.Header().Set(httpconst.HeaderHashSHA256, respHash)
	return w.ResponseWriter.Write(data)
}

// SignatureMiddleware добавляет хэш SHA-256 в заголовок ответа, если он указан в запросе.
func SignatureMiddleware(log *logging.Logger, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем наличие хэша SHA-256 в заголовке запроса.
		headerHash := c.GetHeader(httpconst.HeaderHashSHA256)
		if headerHash != "" {
			var buf []byte
			_, err := c.Request.Body.Read(buf)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			reqHash, err := makeHash(&buf, key)
			if err != nil || reqHash != headerHash {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		rrw := &signatureResponseWriter{ResponseWriter: c.Writer, key: key}
		c.Writer = rrw

		c.Next()
	}
}

// makeHash вычисляет хэш для данных по алгоритму SHA256.
func makeHash(buf *[]byte, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(*buf)
	return hex.EncodeToString(h.Sum(nil)), nil
}
