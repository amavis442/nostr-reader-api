// ABOUTME: Tests for the API access-control helpers (client IP parsing and the
// ABOUTME: loopback-plus-allowlist decision used by the IP allowlist middleware).
package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientIP(t *testing.T) {
	cases := map[string]string{
		"1.2.3.4:5678": "1.2.3.4",
		"[::1]:8080":   "::1",
		"127.0.0.1":    "127.0.0.1", // bare host, no port
	}
	for in, want := range cases {
		if got := clientIP(in); got != want {
			t.Errorf("clientIP(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIPAllowed(t *testing.T) {
	allowed := toIPSet([]string{"192.168.1.50"})

	if !ipAllowed("127.0.0.1", allowed) {
		t.Error("loopback IPv4 should always be allowed")
	}
	if !ipAllowed("::1", allowed) {
		t.Error("loopback IPv6 should always be allowed")
	}
	if !ipAllowed("192.168.1.50", allowed) {
		t.Error("configured IP should be allowed")
	}
	if ipAllowed("192.168.1.99", allowed) {
		t.Error("unlisted IP should be denied")
	}
}

func TestIPAllowedEmptyListIsLoopbackOnly(t *testing.T) {
	empty := toIPSet(nil)

	if !ipAllowed("127.0.0.1", empty) {
		t.Error("loopback should be allowed with an empty allowlist")
	}
	if ipAllowed("192.168.1.50", empty) {
		t.Error("empty allowlist must reject any non-loopback IP")
	}
}

func TestIPAllowlistMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := ipAllowlist([]string{"203.0.113.7"})(next)

	cases := []struct {
		remoteAddr string
		wantStatus int
	}{
		{"127.0.0.1:1234", http.StatusOK},
		{"203.0.113.7:1234", http.StatusOK},
		{"203.0.113.8:1234", http.StatusForbidden},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "/api/getrelays", nil)
		req.RemoteAddr = tc.remoteAddr
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != tc.wantStatus {
			t.Errorf("RemoteAddr %s: got status %d, want %d", tc.remoteAddr, rec.Code, tc.wantStatus)
		}
	}
}
