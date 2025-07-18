package handlers

import (
	"net"
	"net/http"
)

// trustedSubnetMiddleware проверяет, что X-Real-IP входит в доверенную подсеть
func trustedSubnetMiddleware(trustedSubnet string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				next.ServeHTTP(w, r)
				return
			}
			ipStr := r.Header.Get("X-Real-IP")
			if ipStr == "" {
				http.Error(w, "X-Real-IP required", http.StatusForbidden)
				return
			}
			ip := net.ParseIP(ipStr)
			_, subnet, err := net.ParseCIDR(trustedSubnet)
			if err != nil || ip == nil || !subnet.Contains(ip) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
