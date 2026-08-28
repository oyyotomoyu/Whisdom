// Package apis handles HTTP request/response concerns: routing,
// authentication, permission checks, and translating system package errors
// into JSON responses. Business logic stays in the system package; handlers
// here should stay thin.
package apis

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/whisdom/server/system"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

type errorBody struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorBody{Error: message})
}

// writeSystemError maps a system package sentinel error to the right HTTP
// status. Unrecognized errors become 500s with a generic message so internal
// details never reach the client.
func writeSystemError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, system.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, system.ErrPermissionDenied):
		writeError(w, http.StatusForbidden, "permission denied")
	case errors.Is(err, system.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, system.ErrInvalidToken):
		writeError(w, http.StatusUnauthorized, "invalid or expired token")
	case errors.Is(err, system.ErrUserDisabled):
		writeError(w, http.StatusForbidden, "account is disabled")
	case errors.Is(err, system.ErrConflict):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusBadRequest, err.Error())
	}
}

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
