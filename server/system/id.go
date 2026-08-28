package system

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID returns a short, random, prefixed identifier such as "usr_1a2b3c4d".
// IDs are opaque to clients and never reused across restarts.
func NewID(prefix string) string {
	buf := make([]byte, 6)
	_, _ = rand.Read(buf)
	return prefix + "_" + hex.EncodeToString(buf)
}
