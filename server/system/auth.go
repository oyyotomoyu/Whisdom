package system

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password for storage.
func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword reports whether plain matches a hash produced by
// HashPassword.
func VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// AuthenticateUser validates credentials and returns the matching user.
// Nonexistent accounts and wrong passwords both return ErrInvalidCredentials
// so a caller can't use login to enumerate registered emails; a disabled
// account is only revealed once the password has already been proven.
func (s *Store) AuthenticateUser(email, password string) (*User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		// Run the hash comparison anyway against a fixed dummy hash so the
		// response time doesn't leak whether the email exists.
		VerifyPassword(dummyHash, password)
		return nil, ErrInvalidCredentials
	}
	if !VerifyPassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	if !user.Active {
		return nil, ErrUserDisabled
	}
	return user, nil
}

// dummyHash is a valid bcrypt hash of an unguessed constant, used only to
// keep AuthenticateUser's timing consistent for unknown emails.
const dummyHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5UlZ6DjHKz5G8k5m3v6iVEqQwzS0S"

// --- JWT (HS256) -----------------------------------------------------------
//
// A small hand-rolled implementation is used instead of a dependency: the
// platform only needs to issue and verify its own short-lived tokens, so the
// full JOSE surface (multiple algorithms, JWKS, etc.) would be unused
// complexity.

// TokenType distinguishes access tokens from refresh tokens so one can't be
// used in place of the other.
type TokenType string

const (
	TokenAccess  TokenType = "access"
	TokenRefresh TokenType = "refresh"
)

// Claims is the JWT payload Whisdom issues.
type Claims struct {
	Subject  string    `json:"sub"`
	Type     TokenType `json:"type"`
	IssuedAt int64     `json:"iat"`
	Expires  int64     `json:"exp"`
}

var jwtHeader = base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

// IssueToken signs a new JWT for userID, valid for ttl.
func IssueToken(secret, userID string, tokenType TokenType, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Subject:  userID,
		Type:     tokenType,
		IssuedAt: now.Unix(),
		Expires:  now.Add(ttl).Unix(),
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("encode claims: %w", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := jwtHeader + "." + payload
	sig := sign(secret, signingInput)
	return signingInput + "." + sig, nil
}

// ParseToken verifies a JWT's signature and expiry and returns its claims.
func ParseToken(secret, token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSig := sign(secret, signingInput)
	if !hmac.Equal([]byte(expectedSig), []byte(parts[2])) {
		return Claims{}, ErrInvalidToken
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}

	if time.Now().Unix() > claims.Expires {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}

func sign(secret, signingInput string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
