package models

// Unit tests for models/token.go

import (
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
	db.AutoMigrate(&User{}, &Token{})
	SetDB(db) // Set the global db variable for model methods
	return db
}

// Add this helper to set the global db variable
func SetDB(testDB *gorm.DB) {
	db = testDB
}

func TestGenerateToken_JWTClaims(t *testing.T) {
	os.Setenv("API_SECRET", "testsecret")
	os.Setenv("API_KEY", "testkey")
	_ = setupTestDB()
	openID := "jwt_claims_openid"
	token, err := generateToken(openID)
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
	if claims["open_id"] != openID {
		t.Errorf("open_id claim mismatch: got %v", claims["open_id"])
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
	tok, err := (&Token{}).FindOrCreateUserToken(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok == nil || tok.AccessToken == "" {
		t.Fatal("Token should be created and not empty")
	}
	var user User
	if err := db.Where("open_id = ?", openid).First(&user).Error; err != nil {
		t.Fatal("User should be created in DB")
	}
	if user.SessionKey != sessionKey {
		t.Errorf("SessionKey mismatch: got %v", user.SessionKey)
	}
}

func TestGenJWTToken_Expiration(t *testing.T) {
	secret := "expire_secret"
	openID := "expire_openid"
	expires := time.Second * 1
	tkn, err := genJWTToken(secret, openID, expires)
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
	if !ok {
		t.Fatal("claims not mapclaims")
	}
	if claims["open_id"] != openID {
		t.Errorf("open_id claim mismatch: got %v", claims["open_id"])
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
