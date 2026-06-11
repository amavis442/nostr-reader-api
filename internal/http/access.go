// ABOUTME: IP allowlist middleware for the API. Loopback is always permitted;
// ABOUTME: any other client IP must appear in the configured allowlist.
package http

import (
	"net"
	"net/http"
	"strings"
)

// clientIP extracts the host portion from a RemoteAddr (host:port). A bare host
// without a port is returned unchanged.
func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

// toIPSet turns a slice of configured IPs into a lookup set, ignoring blanks.
func toIPSet(ips []string) map[string]struct{} {
	set := make(map[string]struct{}, len(ips))
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			set[ip] = struct{}{}
		}
	}
	return set
}

// ipAllowed reports whether a client IP may access the API. Loopback is always
// allowed (the app is local-first); any other IP must be in the allowlist. An
// empty allowlist therefore means loopback-only.
func ipAllowed(ip string, allowed map[string]struct{}) bool {
	if parsed := net.ParseIP(ip); parsed != nil && parsed.IsLoopback() {
		return true
	}
	_, ok := allowed[ip]
	return ok
}

// ipAllowlist is middleware that rejects clients whose IP is neither loopback
// nor present in the configured allowlist with a 403.
func ipAllowlist(allowed []string) func(http.Handler) http.Handler {
	set := toIPSet(allowed)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !ipAllowed(clientIP(r.RemoteAddr), set) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
