package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/monit/internal/httpconst"
	"github.com/gitslim/monit/internal/logging"
)

func TrustedSubnetMiddleware(log *logging.Logger, allowedSubnet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.GetHeader(httpconst.HeaderXRealIP)

		if isAllowedSubnet(ip, allowedSubnet) {
			c.Next()
			return
		}

		log.Debugf("ip not allowed: %s", ip)
		c.AbortWithStatus(http.StatusForbidden)
	}
}

// isAllowedSubnet проверяет, что ip принадлежит подсети
func isAllowedSubnet(ip string, allowedSubnet string) bool {
	_, subnet, err := net.ParseCIDR(allowedSubnet)
	if err != nil {
		return false
	}
	return subnet.Contains(net.ParseIP(ip))
}
