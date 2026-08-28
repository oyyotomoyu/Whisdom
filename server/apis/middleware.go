package apis

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

type contextKey string

const (
	ctxUserKey   contextKey = "auth_user"
	ctxUserIDKey contextKey = "auth_user_id"
)

// currentUser returns the authenticated user bound by requireAuth, or false
// if the request reached this point without authentication.
func currentUser(ctx context.Context) (*system.AuthUser, bool) {
	u, ok := ctx.Value(ctxUserKey).(*system.AuthUser)
	return u, ok
}

// currentUserID is used by the log middleware to bind the authenticated
// user ID to every log line written during the request.
func currentUserID(r *http.Request) string {
	id, _ := r.Context().Value(ctxUserIDKey).(string)
	return id
}

// requireAuth validates the Authorization bearer token, loads the user, and
// stores it in the request context. Missing/invalid tokens and disabled
// accounts are rejected before the handler runs.
func (a *App) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}

		claims, err := system.ParseToken(a.Config.JWTSecret, token)
		if err != nil || claims.Type != system.TokenAccess {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		user, err := a.Store.GetUser(claims.Subject)
		if err != nil || !user.Active {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		authUser, err := a.Store.ResolveAuthUser(user)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to resolve user")
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserKey, authUser)
		ctx = context.WithValue(ctx, ctxUserIDKey, authUser.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

// requirePermission rejects the request unless the authenticated user (set
// by requireAuth, which must run first) holds perm.
func (a *App) requirePermission(perm system.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := currentUser(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}
			for _, p := range user.Permissions {
				if p == perm {
					next.ServeHTTP(w, r)
					return
				}
			}
			writeError(w, http.StatusForbidden, "permission denied")
		})
	}
}

// hasAnyPermission reports whether user holds at least one of perms. Used
// for the handful of routes docs/server.md grants to more than one
// permission (e.g. the training material path accepts config.manage OR
// materials.path.edit), which the AND-only protected() chain can't express.
func hasAnyPermission(user *system.AuthUser, perms ...system.Permission) bool {
	for _, want := range perms {
		for _, have := range user.Permissions {
			if have == want {
				return true
			}
		}
	}
	return false
}

// recoverMiddleware converts a panicking handler into a 500 response instead
// of taking down the whole server.
func (a *App) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logs.FromContext(r.Context()).Log("error", fmt.Sprintf("panic: %v", rec))
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// public wraps a handler with request logging only — for routes that don't
// require authentication (login, health).
func (a *App) public(h http.HandlerFunc) http.Handler {
	return a.recoverMiddleware(logs.Middleware(a.Logs, currentUserID)(h))
}

// protected wraps a handler with authentication, request logging, and
// (optionally) one or more required permissions, applied in the order
// described by docs/server.md: auth -> log -> permission -> handler.
func (a *App) protected(h http.HandlerFunc, perms ...system.Permission) http.Handler {
	var handler http.Handler = h
	for i := len(perms) - 1; i >= 0; i-- {
		handler = a.requirePermission(perms[i])(handler)
	}
	handler = logs.Middleware(a.Logs, currentUserID)(handler)
	handler = a.requireAuth(handler)
	return a.recoverMiddleware(handler)
}
