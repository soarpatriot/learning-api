package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPostToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/token", PostToken)

	// Test missing code param
	w := httptest.NewRecorder()
	reqBody := bytes.NewBufferString(`{}`)
	req, _ := http.NewRequest("POST", "/token", reqBody)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	// Test valid code param
	w = httptest.NewRecorder()
	body, _ := json.Marshal(map[string]string{"code": "abc123"})
	req, _ = http.NewRequest("POST", "/token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
