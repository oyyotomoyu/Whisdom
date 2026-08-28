package system

import (
	"testing"
	"time"
)

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !VerifyPassword(hash, "correct horse battery staple") {
		t.Error("VerifyPassword should accept the correct password")
	}
	if VerifyPassword(hash, "wrong password") {
		t.Error("VerifyPassword should reject an incorrect password")
	}
}

func TestIssueAndParseToken(t *testing.T) {
	secret := "test-secret"

	token, err := IssueToken(secret, "usr_123", TokenAccess, time.Minute)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.Subject != "usr_123" {
		t.Errorf("Subject = %q, want usr_123", claims.Subject)
	}
	if claims.Type != TokenAccess {
		t.Errorf("Type = %q, want %q", claims.Type, TokenAccess)
	}
}

func TestParseTokenRejectsExpired(t *testing.T) {
	token, err := IssueToken("secret", "usr_123", TokenAccess, -time.Minute)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	if _, err := ParseToken("secret", token); err == nil {
		t.Error("ParseToken should reject an expired token")
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	token, err := IssueToken("secret-a", "usr_123", TokenAccess, time.Minute)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	if _, err := ParseToken("secret-b", token); err == nil {
		t.Error("ParseToken should reject a token signed with a different secret")
	}
}

func TestParseTokenRejectsMalformed(t *testing.T) {
	for _, bad := range []string{"", "not-a-jwt", "a.b", "a.b.c.d"} {
		if _, err := ParseToken("secret", bad); err == nil {
			t.Errorf("ParseToken(%q) should fail", bad)
		}
	}
}

func TestAuthenticateUser(t *testing.T) {
	store := NewStore("/training/default/")
	roles := store.SeedDefaultRoles()
	var userRole *Role
	for _, r := range roles {
		if r.Name == RoleUser {
			userRole = r
		}
	}
	if userRole == nil {
		t.Fatal("seeded roles missing user role")
	}

	hash, _ := HashPassword("s3cret-password")
	if _, err := store.CreateUser("Alex", "alex@example.com", hash, userRole.ID); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	if _, err := store.AuthenticateUser("alex@example.com", "s3cret-password"); err != nil {
		t.Errorf("AuthenticateUser with correct credentials: %v", err)
	}
	if _, err := store.AuthenticateUser("alex@example.com", "wrong"); err != ErrInvalidCredentials {
		t.Errorf("AuthenticateUser wrong password error = %v, want ErrInvalidCredentials", err)
	}
	if _, err := store.AuthenticateUser("nobody@example.com", "whatever"); err != ErrInvalidCredentials {
		t.Errorf("AuthenticateUser unknown email error = %v, want ErrInvalidCredentials", err)
	}
}
