// Package transport содержит реализацию кастомного Transport для http
package transport

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/gitslim/monit/internal/httpconst"
)

// CustomTransport расширяет http.Transport
type CustomTransport struct {
	Transport http.RoundTripper
	ip        net.IP
	mu        sync.RWMutex
}

// NewCustomTransport создает кастомный Transport
func NewCustomTransport() *CustomTransport {
	// Создаем кастомный Transport
	transport := &CustomTransport{}
	transport.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Устанавливаем соединение
			conn, err := net.Dial(network, addr)
			if err != nil {
				return nil, err
			}

			// Получаем локальный IP
			localAddr := conn.LocalAddr().(*net.TCPAddr)
			ip := localAddr.IP

			// Сохраняем IP в CustomTransport
			transport.mu.Lock()
			transport.ip = ip
			transport.mu.Unlock()

			// Возвращаем соединение
			return conn, nil
		},
	}

	return transport
}

// RoundTrip перехватывает запрос и устанавливает заголовок X-Real-IP
func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.RLock()
	ip := t.ip
	t.mu.RUnlock()

	// Если IP установлен, добавляем заголовок X-Real-IP
	if ip != nil {
		req.Header.Set(httpconst.HeaderXRealIP, ip.String())
	}

	// Выполняем запрос через дефолтный Transport
	return t.Transport.RoundTrip(req)
}
