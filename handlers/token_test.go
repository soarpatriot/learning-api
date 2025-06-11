package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"learning-api/helpers"
	"learning-api/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (m *mockHandlerDouyinClient) V2Jscode2session(req *openApiSdkClient.V2Jscode2sessionRequest) (*openApiSdkClient.V2Jscode2sessionResponse, error) {
	return &openApiSdkClient.V2Jscode2sessionResponse{
		Data: &openApiSdkClient.V2Jscode2sessionResponseData{
			Openid:     ptrString("mock_openid"),
			Unionid:    ptrString("mock_unionid"),
			SessionKey: ptrString("mock_sessionkey"),
		},
		ErrNo: ptrInt64(0),
	}, nil
}

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Token{})
	return db
}

func setModelsDB(db *gorm.DB) {
	models.SetDB(db)
}

func TestPostToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setModelsDB(setupTestDB())
	r := gin.Default()
	r.POST("/token", PostToken)

	// Patch newDouyinClientFunc to return a mock client
	newDouyinClientFunc = func() helpers.ThirdPartyClient {
		return &mockHandlerDouyinClient{}
	}

	// Valid request
	body, _ := json.Marshal(map[string]string{"code": "mock_code"})
	req, _ := http.NewRequest("POST", "/token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Invalid request (missing code)
	req, _ = http.NewRequest("POST", "/token", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestPostRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	setModelsDB(db)
	r := gin.Default()
	r.POST("/refresh-token", PostRefreshToken)

	// Generate a real access token and refresh token
	token := models.NewToken()
	if err := token.GenTokenWithDate(); err != nil {
		t.Fatal("Failed to generate token:", err)
	}

	if err := db.Create(token).Error; err != nil {
		fmt.Println("Error creating token:", err)
		t.Fatal("Failed to create token in test database")
	}

	// Valid request
	body, _ := json.Marshal(map[string]string{"access_token": token.AccessToken, "refresh_token": token.RefreshToken})
	req, _ := http.NewRequest("POST", "/refresh-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//conver the w.Body to a token struct
	var responseToken models.Token
	fmt.Println("Response Body:", w.Body.String())
	if err := json.Unmarshal(w.Body.Bytes(), &responseToken); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if responseToken.AccessToken == "" || responseToken.RefreshToken == "" {
		t.Errorf("Expected non-empty access and refresh tokens in response")
	}
	if responseToken.RefreshTokenExpiresIn != token.RefreshTokenExpiresIn-token.AccessTokenExpiresIn {
		t.Errorf("Expected RefreshTokenExpiresIn to be %d, got %d", token.RefreshTokenExpiresIn-token.AccessTokenExpiresIn, responseToken.RefreshTokenExpiresIn)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Body)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Invalid request (missing tokens)
	req, _ = http.NewRequest("POST", "/refresh-token", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	// Expired refresh token
	expiredToken := &models.Token{
		AccessToken:           "expired_access_token",
		RefreshToken:          "expired_refresh_token",
		RefreshTokenExpiresIn: 3600, // 1 hour
		CreatedAt:             time.Now().Add(-2 * time.Hour),
	}

	db.Create(expiredToken)
	req, _ = http.NewRequest("POST", "/refresh-token", bytes.NewBuffer([]byte(`{"access_token":"expired_access_token","refresh_token":"expired_refresh_token"}`)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// mockHandlerDouyinClient implements helpers.ThirdPartyClient
// and returns a mock token for testing
type mockHandlerDouyinClient struct{}

func (m *mockHandlerDouyinClient) Jscode2session(code string, anonymousCode string) (*models.Token, error) {
	return &models.Token{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
	}, nil
}

func ptrString(s string) *string { return &s }
func ptrInt64(i int64) *int64    { return &i }
