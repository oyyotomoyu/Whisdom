package logs

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

type contextKey struct{}

// RequestLogger is a Service bound to one request's client IP and
// authenticated user ID, so handlers only ever supply status and content.
type RequestLogger struct {
	svc    *Service
	ip     string
	userID string
}

// Log writes one record for the bound request. It is nil-safe: calling Log
// on a RequestLogger obtained outside a request (missing middleware) is a
// no-op rather than a panic, since audit logging must never take down a
// request. Write failures are reported to stderr.
func (l *RequestLogger) Log(status, content string) {
	if l == nil || l.svc == nil {
		return
	}
	if err := l.svc.Log(status, l.ip, l.userID, content); err != nil {
		fmt.Fprintf(os.Stderr, "logs: failed to write record: %v\n", err)
	}
}

// UserIDFunc extracts the authenticated user ID (if any) from a request.
// Supplied by the apis package so this package stays decoupled from the
// auth context key.
type UserIDFunc func(*http.Request) string

// Middleware binds a RequestLogger carrying the client IP and, when
// available, the authenticated user ID into the request context. It must run
// after authentication middleware so userIDFn can see the resolved user.
func Middleware(svc *Service, userIDFn UserIDFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := ""
			if userIDFn != nil {
				userID = userIDFn(r)
			}
			rl := &RequestLogger{svc: svc, ip: clientIP(r), userID: userID}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, rl)))
		})
	}
}

// FromContext returns the RequestLogger bound by Middleware, or nil if none
// is present.
func FromContext(ctx context.Context) *RequestLogger {
	rl, _ := ctx.Value(contextKey{}).(*RequestLogger)
	return rl
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		for _, part := range strings.Split(fwd, ",") {
			if part = strings.TrimSpace(part); part != "" {
				return part
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
