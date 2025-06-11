package models_test

// Unit tests for models/token.go

import (
	"learning-api/models"
	"os"
	"testing"
	"time"

	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Token{})
	models.SetDB(db) // Set the global db variable for model methods
	return db
}

func TestGenerateToken_JWTClaims(t *testing.T) {
	os.Setenv("CLIENT_SECRET", "testsecret")
	os.Setenv("CLIENT_KEY", "testkey")
	_ = setupTestDB()
	token := models.NewToken()
	err := token.GenTokenWithDate()
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken == "" || token.RefreshToken == "" {
		t.Fatal("Tokens should not be empty")
	}
	parsed, err := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("testsecret"), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("AccessToken is not valid: %v", err)
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims not mapclaims")
	}
	if claims["iss"] != "learning" {
		t.Errorf("issuer claim mismatch: got %v", claims["iss"])
	}
	if claims["aud"] != "douyin" {
		t.Errorf("audience claim mismatch: got %v", claims["aud"])
	}
}

func TestFindOrCreateUserToken_CreatesUserAndToken(t *testing.T) {
	db := setupTestDB()
	openid := "openid_test"
	unionid := "unionid_test"
	sessionKey := "sessionkey_test"
	data := &openApiSdkClient.V2Jscode2sessionResponseData{}
	data.SetOpenid(openid)
	data.SetUnionid(unionid)
	data.SetSessionKey(sessionKey)
	tok, err := models.FindOrCreateUserToken(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok == nil || tok.AccessToken == "" || tok.RefreshToken == "" {
		t.Fatal("Token should be created and not empty")
	}
	var user models.User
	if err := db.Where("open_id = ?", openid).First(&user).Error; err != nil {
		t.Fatal("User should be created in DB")
	}
	if user.SessionKey != sessionKey {
		t.Errorf("SessionKey mismatch: got %v", user.SessionKey)
	}
}

func TestFindOrCreateUserToken_NewAndExistingUser(t *testing.T) {
	db := setupTestDB()
	openid := "new_openid_test"
	unionid := "new_unionid_test"
	sessionKey := "new_sessionkey_test"
	data := &openApiSdkClient.V2Jscode2sessionResponseData{}
	data.SetOpenid(openid)
	data.SetUnionid(unionid)
	data.SetSessionKey(sessionKey)

	// First time: create new user and token
	tok, err := models.FindOrCreateUserToken(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok == nil || tok.AccessToken == "" || tok.RefreshToken == "" {
		t.Fatal("Token should be created and not empty")
	}
	var user models.User
	if err := db.Where("open_id = ?", openid).First(&user).Error; err != nil {
		t.Fatal("User should be created in DB")
	}
	if user.SessionKey != sessionKey {
		t.Errorf("SessionKey mismatch: got %v", user.SessionKey)
	}
	db.Preload("Tokens").First(&user, "open_id = ?", openid)
	if len(user.Tokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(user.Tokens))
	}

	// Second time: save existing user with new token
	newSessionKey := "updated_sessionkey_test"
	data.SetSessionKey(newSessionKey)
	tok, err = models.FindOrCreateUserToken(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok == nil || tok.AccessToken == "" || tok.RefreshToken == "" {
		t.Fatal("Token should be created and not empty")
	}
	if err := db.Where("open_id = ?", openid).First(&user).Error; err != nil {
		t.Fatal("User should exist in DB")
	}
	if user.SessionKey != newSessionKey {
		t.Errorf("SessionKey mismatch: got %v", user.SessionKey)
	}
	db.Preload("Tokens").First(&user, "open_id = ?", openid)
	if len(user.Tokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(user.Tokens))
	}
}

func TestGenJWTToken_Expiration(t *testing.T) {
	secret := "expire_secret"
	expires := time.Second * 1
	tkn, err := models.GenJWTToken(secret, expires)
	if err != nil {
		t.Fatalf("failed to generate jwt: %v", err)
	}
	parsed, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("JWT not valid: %v", err)
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if claims["exp"] == nil {
		t.Error("Expiration claim should be present")
	}
	if !ok {
		t.Fatal("claims not mapclaims")
	}
	// Wait for expiration
	time.Sleep(2 * time.Second)
	parsed, err = jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err == nil && parsed.Valid {
		t.Error("JWT should be expired but is still valid")
	}
}
