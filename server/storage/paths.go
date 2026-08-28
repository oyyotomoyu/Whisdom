// Package storage owns file path validation and on-disk persistence for
// uploaded training materials.
package storage

import (
	"errors"
	"path"
	"strings"
)

// ValidateDestinationPath enforces the rules shared by the training-material
// config path and per-material destination paths: it must be an absolute,
// traversal-free virtual path.
func ValidateDestinationPath(p string) error {
	if strings.TrimSpace(p) == "" {
		return errors.New("path is required")
	}
	if strings.ContainsRune(p, 0) {
		return errors.New("path cannot contain null bytes")
	}
	if !strings.HasPrefix(p, "/") {
		return errors.New("path must start with /")
	}
	if strings.Contains(p, "..") {
		return errors.New("path cannot contain \"..\"")
	}
	return nil
}

// NormalizeDestinationPath cleans a validated virtual path, preserving a
// trailing slash when the caller supplied one (destination paths are
// conventionally directory-like).
func NormalizeDestinationPath(p string) string {
	cleaned := path.Clean(p)
	if cleaned != "/" && strings.HasSuffix(p, "/") && !strings.HasSuffix(cleaned, "/") {
		cleaned += "/"
	}
	return cleaned
}
