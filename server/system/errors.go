package system

import "errors"

// Sentinel errors handlers translate into HTTP status codes.
var (
	ErrNotFound           = errors.New("not found")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrConflict           = errors.New("conflict")
)
