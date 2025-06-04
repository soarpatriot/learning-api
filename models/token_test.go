package models

import (
	"testing"
	"time"
)

func TestGenJWTAccessToken(t *testing.T) {
	token := &Token{UserID: 1}
	secret := "testsecret"
	expiresIn := time.Hour
	jwtStr, err := token.GenJWTAccessToken(secret, expiresIn)
	if err != nil {
		t.Fatalf("GenJWTAccessToken failed: %v", err)
	}
	if jwtStr == "" {
		t.Error("Expected non-empty JWT access token string")
	}
}

func TestGenRefreshToken(t *testing.T) {
	token := &Token{}
	length := 32
	rt, err := token.GenRefreshToken(length)
	if err != nil {
		t.Fatalf("GenRefreshToken failed: %v", err)
	}
	if len(rt) == 0 {
		t.Error("Expected non-empty refresh token string")
	}
}

func TestGenerateAToken(t *testing.T) {
	secret := "testsecret"
	userID := uint(42)
	token := generate_a_token(secret, userID)
	if token == nil {
		t.Fatal("generate_a_token returned nil")
	}
	if token.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, token.UserID)
	}
	if token.AccessToken == "" {
		t.Error("Expected non-empty AccessToken")
	}
	if token.RefreshToken == "" {
		t.Error("Expected non-empty RefreshToken")
	}
	if token.AccessTokenExpiresIn != 3600 {
		t.Errorf("Expected AccessTokenExpiresIn 3600, got %d", token.AccessTokenExpiresIn)
	}
	if token.RefreshTokenExpiresIn != 7200 {
		t.Errorf("Expected RefreshTokenExpiresIn 7200, got %d", token.RefreshTokenExpiresIn)
	}
}

func TestTempUser(t *testing.T) {
	user := temp_user()
	if user == nil {
		t.Fatal("temp_user returned nil")
	}
	if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
		t.Error("Expected CreatedAt and UpdatedAt to be set")
	}
}

func TestJscode2session(t *testing.T) {
	if db == nil {
		t.Skip("db is nil, skipping integration test for Jscode2session")
	}
	// Example: test with an invalid code (should fail gracefully)
	_, err := Jscode2session("invalid_code_for_test")
	if err == nil {
		t.Error("Expected error for invalid code, got nil")
	}
}
