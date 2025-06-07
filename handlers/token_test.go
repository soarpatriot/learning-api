package handlers

import (
	"bytes"
	"encoding/json"
	"learning-api/helpers"
	"learning-api/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"unsafe"

	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type handlerMockDouyinClient struct{}

func (m *handlerMockDouyinClient) V2Jscode2session(req *openApiSdkClient.V2Jscode2sessionRequest) (*openApiSdkClient.V2Jscode2sessionResponse, error) {
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
	dbPtr := reflect.ValueOf(&models.Token{}).Elem().FieldByName("db")
	if dbPtr.CanSet() {
		dbPtr.Set(reflect.ValueOf(db))
	} else {
		// fallback: set unexported global via unsafe
		p := reflect.ValueOf(&models.Token{}).Elem().Addr().Pointer()
		ptr := (*gorm.DB)(unsafe.Pointer(p))
		*ptr = *db
	}
}

func TestPostToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setModelsDB(setupTestDB())
	r := gin.Default()
	r.POST("/token", PostToken)

	// Patch newDouyinClientFunc to return a mock client
	orig := newDouyinClientFunc
	newDouyinClientFunc = func() helpers.ThirdPartyClient {
		return &mockHandlerDouyinClient{}
	}
	defer func() { newDouyinClientFunc = orig }()

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

// mockHandlerDouyinClient implements helpers.ThirdPartyClient
// and returns a mock token for testing
type mockHandlerDouyinClient struct{}

func (m *mockHandlerDouyinClient) Jscode2session(code string) (*models.Token, error) {
	return &models.Token{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
	}, nil
}

func ptrString(s string) *string { return &s }
func ptrInt64(i int64) *int64    { return &i }
