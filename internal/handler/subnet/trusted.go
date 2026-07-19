package subnet

import (
	"net"
	"net/http"
	"strings"

	"github.com/liebeSonne/shortlink/internal/logger"
)

func NewTrustedSubnetMiddleware(
	next http.HandlerFunc,
	trustedSubnet string,
	logger logger.Logger,
) http.HandlerFunc {
	trustedSubnet = strings.TrimSpace(trustedSubnet)
	var trustedIPNet *net.IPNet
	if trustedSubnet != "" {
		_, ipNet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			logger.Errorf("parse trusted subnet CIDR error: %v", err)
		}
		if err == nil {
			trustedIPNet = ipNet
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if trustedSubnet == "" || trustedIPNet == nil {
			logger.Warnf("Forbidden on trusted subnet is not set")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		xRealIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
		if xRealIP == "" {
			logger.Warnf("Forbidden on X-Real-IP is not set")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		clientIP := net.ParseIP(xRealIP)
		if clientIP == nil {
			logger.Warnf("Forbidden on parse X-Real-IP is nil")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		if !trustedIPNet.Contains(clientIP) {
			logger.Warnf("Forbidden on no contains client ip in trusted subnet")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
