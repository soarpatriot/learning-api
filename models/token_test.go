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
	return db
}

func TestGenerateToken_JWTClaims(t *testing.T) {
	os.Setenv("API_SECRET", "testsecret")
	os.Setenv("API_KEY", "testkey")
	db = setupTestDB()
	openID := "jwt_claims_openid"
	token := generateToken(openID)
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
	db = setupTestDB()
	openid := "openid_test"
	unionid := "unionid_test"
	sessionKey := "sessionkey_test"
	data := &openApiSdkClient.V2Jscode2sessionResponseData{}
	data.SetOpenid(openid)
	data.SetUnionid(unionid)
	data.SetSessionKey(sessionKey)
	tok, err := findOrCreateUserToken(data)
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
	tok := &Token{}
	secret := "expire_secret"
	openID := "expire_openid"
	expires := time.Second * 1
	tkn, err := tok.genJWTToken(secret, openID, expires)
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

// Mock client for Jscode2session test

type mockDouyinClient struct{}

func (m *mockDouyinClient) V2Jscode2session(req *openApiSdkClient.V2Jscode2sessionRequest) (*openApiSdkClient.V2Jscode2sessionResponse, error) {
	return &openApiSdkClient.V2Jscode2sessionResponse{
		Data: &openApiSdkClient.V2Jscode2sessionResponseData{
			Openid:     ptrString("mock_openid"),
			Unionid:    ptrString("mock_unionid"),
			SessionKey: ptrString("mock_sessionkey"),
		},
		ErrNo: ptrInt64(0),
	}, nil
}

func TestJscode2session(t *testing.T) {
	os.Setenv("API_SECRET", "testsecret")
	os.Setenv("API_KEY", "testkey")
	os.Setenv("APP_ID", "testappid")

	// Patch generateSdkClientFunc to return our mock client
	origGen := generateSdkClientFunc
	generateSdkClientFunc = func() (DouyinClient, error) {
		return &mockDouyinClient{}, nil
	}
	defer func() { generateSdkClientFunc = origGen }()

	db = setupTestDB()
	code := "mock_code"
	token, err := Jscode2session(code)
	if err != nil {
		t.Fatalf("Jscode2session failed: %v", err)
	}
	if token == "" {
		t.Error("Expected non-empty token from Jscode2session")
	}
}

func ptrString(s string) *string { return &s }
func ptrInt64(i int64) *int64    { return &i }
