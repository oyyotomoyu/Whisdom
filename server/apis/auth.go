package apis

import (
	"errors"
	"net/http"
	"time"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

const refreshCookieName = "whisdom_refresh_token"

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string           `json:"access_token"`
	User        *system.AuthUser `json:"user"`
}

// handleLogin validates credentials, issues an access token in the response
// body, and sets a long-lived refresh token as an HttpOnly cookie scoped to
// the auth routes.
func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger := logs.FromContext(r.Context())

	user, err := a.Store.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, system.ErrUserDisabled) {
			logger.Log("warning", "login blocked: account disabled ("+req.Username+")")
			writeError(w, http.StatusForbidden, "account is disabled")
			return
		}
		logger.Log("warning", "login failed for "+req.Username)
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	access, refresh, err := a.issueTokenPair(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue session")
		return
	}
	setRefreshCookie(w, refresh, a.Config.RefreshTokenTTL)

	authUser, err := a.Store.ResolveAuthUser(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to resolve user")
		return
	}

	logger.Log("info", "login succeeded: "+user.ID)
	writeJSON(w, http.StatusOK, loginResponse{AccessToken: access, User: authUser})
}

// handleLogout clears the refresh cookie. Access tokens are short-lived and
// not individually revocable in this in-memory implementation, so callers
// should also discard the access token client-side.
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	clearRefreshCookie(w)
	if user, ok := currentUser(r.Context()); ok {
		logs.FromContext(r.Context()).Log("info", "logout: "+user.ID)
	}
	writeJSON(w, http.StatusOK, nil)
}

type refreshResponse struct {
	AccessToken string `json:"access_token"`
}

// handleRefresh exchanges a valid refresh cookie for a new access token and
// rotates the refresh cookie.
func (a *App) handleRefresh(w http.ResponseWriter, r *http.Request) {
	logger := logs.FromContext(r.Context())

	cookie, err := r.Cookie(refreshCookieName)
	if err != nil || cookie.Value == "" {
		logger.Log("warning", "token refresh failed: missing refresh cookie")
		writeError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	claims, err := system.ParseToken(a.Config.JWTSecret, cookie.Value)
	if err != nil || claims.Type != system.TokenRefresh {
		logger.Log("warning", "token refresh failed: invalid refresh token")
		writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	user, err := a.Store.GetUser(claims.Subject)
	if err != nil || !user.Active {
		logger.Log("warning", "token refresh failed: user unavailable")
		writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	access, refresh, err := a.issueTokenPair(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue session")
		return
	}
	setRefreshCookie(w, refresh, a.Config.RefreshTokenTTL)

	logger.Log("info", "token refresh succeeded: "+user.ID)
	writeJSON(w, http.StatusOK, refreshResponse{AccessToken: access})
}

// handleMe returns the authenticated user's profile and resolved
// permissions.
func (a *App) handleMe(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing bearer token")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (a *App) issueTokenPair(userID string) (access, refresh string, err error) {
	access, err = system.IssueToken(a.Config.JWTSecret, userID, system.TokenAccess, a.Config.AccessTokenTTL)
	if err != nil {
		return "", "", err
	}
	refresh, err = system.IssueToken(a.Config.JWTSecret, userID, system.TokenRefresh, a.Config.RefreshTokenTTL)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func setRefreshCookie(w http.ResponseWriter, value string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    value,
		Path:     "/api/v1/auth",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	})
}

func clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/api/v1/auth",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
