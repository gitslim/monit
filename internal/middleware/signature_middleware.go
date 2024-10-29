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

type signatureResponseWriter struct {
	gin.ResponseWriter
	key string
}

func (w *signatureResponseWriter) Write(data []byte) (int, error) {
	respHash, err := makeHash(&data, w.key)
	if err != nil {
		return 0, err
	}
	w.ResponseWriter.Header().Set(httpconst.HeaderHashSHA256, respHash)
	return w.ResponseWriter.Write(data)
}

func SignatureMiddleware(log *logging.Logger, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

func makeHash(buf *[]byte, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(*buf)
	return hex.EncodeToString(h.Sum(nil)), nil
}
