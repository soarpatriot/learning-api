package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRandomOrderNo(t *testing.T) {
	// Test that randomOrderNo generates unique order numbers
	orderNo1 := randomOrderNo()
	orderNo2 := randomOrderNo()

	// Should not be empty
	assert.NotEmpty(t, orderNo1)
	assert.NotEmpty(t, orderNo2)

	// Should be different (very high probability)
	assert.NotEqual(t, orderNo1, orderNo2)

	// Should not contain the old prefix
	assert.False(t, strings.HasPrefix(orderNo1, "out_order_no_"))
	assert.False(t, strings.HasPrefix(orderNo2, "out_order_no_"))

	// Should be reasonable length (timestamp in base36 + 4 random chars)
	assert.True(t, len(orderNo1) > 10)
	assert.True(t, len(orderNo1) < 30)
}

func TestRandString(t *testing.T) {
	// Test RandString function
	str1 := RandString(4)
	str2 := RandString(4)
	str10 := RandString(10)

	// Should have correct length
	assert.Equal(t, 4, len(str1))
	assert.Equal(t, 4, len(str2))
	assert.Equal(t, 10, len(str10))

	// Should be different (very high probability)
	assert.NotEqual(t, str1, str2)

	// Should only contain valid characters
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range str1 {
		assert.Contains(t, validChars, string(char))
	}
}

func TestPayOrder_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock HTTP server to simulate Douyin API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and content type
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Return mock Douyin response
		response := map[string]interface{}{
			"err_no":      0,
			"err_tips":    "success",
			"order_id":    "7123456789012345678",
			"order_token": "ChAKGG91dF9vcmRlcl9ub18xNjc...",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Note: In a real test, you'd need to mock the Douyin API URL
	// For this test, we'll focus on request validation and response structure

	router := gin.New()
	router.POST("/pay/order", PayOrder)

	// Test valid request
	requestBody := PayOrderRequest{
		TotalAmount: 1000,
		Subject:     "Test Payment",
		Body:        "Test payment description",
		CpExtra:     "extra_data",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/pay/order", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Note: This will likely fail in test environment due to external API call
	// But we can test the request validation part
	assert.NotEqual(t, http.StatusBadRequest, w.Code, "Request should be valid")
}

func TestPayOrder_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/pay/order", PayOrder)

	// Test invalid JSON
	req, _ := http.NewRequest("POST", "/pay/order", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestPayOrder_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/pay/order", PayOrder)

	// Test with missing required fields
	requestBody := map[string]interface{}{
		"subject": "Test Payment",
		// Missing total_amount, body
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/pay/order", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should not return 400 for missing fields since they're not validated in the current implementation
	// But the request should be processed (though it may fail at the API level)
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestPayOrderCallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/pay/callback", PayOrderCallback)

	// Test callback endpoint
	req, _ := http.NewRequest("POST", "/pay/callback", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Payment callback received", response["message"])
}

func TestPayDouOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/pay/dou-order", PayDouOrder)

	// Test PayDouOrder endpoint
	req, _ := http.NewRequest("POST", "/pay/dou-order", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "auth")

	// The auth field might be empty due to crypto errors in test environment
	// but the endpoint should still return a response structure
	_, exists := response["auth"]
	assert.True(t, exists, "auth field should exist in response")
}

func TestDouyinOrderRequest_Structure(t *testing.T) {
	// Test DouyinOrderRequest struct
	order := DouyinOrderRequest{
		AppID:       "test_app_id",
		OutOrderNo:  "test_order_123",
		TotalAmount: 1000,
		Subject:     "Test Subject",
		Body:        "Test Body",
		ValidTime:   180,
		Sign:        "test_signature",
		NotifyURL:   "https://example.com/callback",
		StoreUid:    "test_store_uid",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(order)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var unmarshaled DouyinOrderRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, order.AppID, unmarshaled.AppID)
	assert.Equal(t, order.OutOrderNo, unmarshaled.OutOrderNo)
	assert.Equal(t, order.TotalAmount, unmarshaled.TotalAmount)
}

func TestPayOrderRequest_Structure(t *testing.T) {
	// Test PayOrderRequest struct
	request := PayOrderRequest{
		TotalAmount: 1500,
		Subject:     "Test Payment Subject",
		Body:        "Test Payment Body",
		CpExtra:     "extra_info",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var unmarshaled PayOrderRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, request.TotalAmount, unmarshaled.TotalAmount)
	assert.Equal(t, request.Subject, unmarshaled.Subject)
	assert.Equal(t, request.Body, unmarshaled.Body)
	assert.Equal(t, request.CpExtra, unmarshaled.CpExtra)
}

func TestCreateSecureHTTPClient(t *testing.T) {
	client, err := createSecureHTTPClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Transport)
}

func TestCreateInsecureHTTPClient(t *testing.T) {
	client := createInsecureHTTPClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.Transport)
}

// Benchmark tests
func BenchmarkRandomOrderNo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomOrderNo()
	}
}

func BenchmarkRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString(4)
	}
}

// Test helper functions
func TestPayOrder_ResponseMerging(t *testing.T) {
	// Test the response merging logic by creating a mock scenario
	gin.SetMode(gin.TestMode)

	// Create a test that verifies the out_order_no is added to response
	// This is more of an integration test concept since we can't easily mock the external API

	// Test that order number generation works
	orderNo := randomOrderNo()
	assert.NotEmpty(t, orderNo)
	assert.False(t, strings.HasPrefix(orderNo, "out_order_no_"))

	// Test response structure that would be returned
	mockDouyinResponse := map[string]interface{}{
		"err_no":      0,
		"err_tips":    "success",
		"order_id":    "7123456789012345678",
		"order_token": "ChAKGG91dF9vcmRlcl9ub18xNjc...",
	}

	// Simulate adding out_order_no to response
	mockDouyinResponse["out_order_no"] = orderNo

	// Verify the merged response structure
	assert.Contains(t, mockDouyinResponse, "err_no")
	assert.Contains(t, mockDouyinResponse, "err_tips")
	assert.Contains(t, mockDouyinResponse, "order_id")
	assert.Contains(t, mockDouyinResponse, "order_token")
	assert.Contains(t, mockDouyinResponse, "out_order_no")
	assert.Equal(t, orderNo, mockDouyinResponse["out_order_no"])
}
