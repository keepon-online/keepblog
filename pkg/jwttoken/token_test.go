package jwttoken

import (
	"testing"
	"time"
)

// withTestSecret 保证测试不依赖外部配置即可签发 token
func withTestSecret(t *testing.T) {
	t.Helper()
	if len(MySecret) == 0 {
		MySecret = []byte("test-secret-0123456789abcdef0123456")
	}
}

func TestCreateAndParseToken(t *testing.T) {
	withTestSecret(t)

	token, err := CreateToken("admin")
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("CreateToken returned empty token")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.Username != "admin" {
		t.Fatalf("expected username admin, got %q", claims.Username)
	}
	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
		t.Fatal("token should carry a future expiry")
	}
}

func TestParseInvalidToken(t *testing.T) {
	withTestSecret(t)

	if _, err := ParseToken("not-a-token"); err == nil {
		t.Fatal("expected error for malformed token")
	}
}

func TestAccessTokenAndRefreshToken(t *testing.T) {
	withTestSecret(t)

	access, err := AccessToken("admin")
	if err != nil {
		t.Fatalf("AccessToken failed: %v", err)
	}
	refresh, err := RefreshToken("admin")
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	if access == "" || refresh == "" || access == refresh {
		t.Fatal("access and refresh tokens should be non-empty and different")
	}
}
