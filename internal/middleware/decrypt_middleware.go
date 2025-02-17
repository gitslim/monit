package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/security"
)

// DecryptMiddleware - мидлварь для расшифровки тела запроса.
func DecryptMiddleware(privateKeyPath string) (gin.HandlerFunc, error) {
	privateKey, err := security.ReadRSAPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	// Возвращаем функцию-мидлварь.
	return func(c *gin.Context) {
		// Читаем зашифрованное тело запроса
		encryptedBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer func() { _ = c.Request.Body.Close() }()

		// Расшифровываем тело запроса
		decryptedData, err := security.DecryptRSA(privateKey, encryptedBody)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(decryptedData))

		c.Next()
	}, nil
}
