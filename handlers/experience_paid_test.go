package handlers

import (
	"bytes"
	"encoding/json"
	"learning-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMarkExperiencePaid(t *testing.T) {
	// Setup test database
	db := models.InitTestDB()
	db.AutoMigrate(&models.User{}, &models.Topic{}, &models.Experience{}, &models.Order{})
	models.SetDB(db)

	// Create test user
	user := models.User{
		OpenID: "test_user",
		Name:   "Test User",
	}
	db.Create(&user)

	// Create test topic
	topic := models.Topic{
		Name: "Test Topic",
	}
	db.Create(&topic)

	// Create test experience
	experience := models.Experience{
		TopicID: topic.ID,
		UserID:  user.ID,
	}
	db.Create(&experience)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add middleware to set current user
	router.Use(func(c *gin.Context) {
		c.Set("currentUser", user)
		c.Next()
	})

	router.POST("/experiences/:id/paid", MarkExperiencePaid)

	t.Run("successful payment", func(t *testing.T) {
		// Prepare request body
		requestBody := MarkPaidRequest{
			OrderNo:    "ORD20250624001",
			OutOrderNo: "PAY20250624001",
			Price:      1000,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// Create request
		req, _ := http.NewRequest("POST", "/experiences/1/paid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check response structure
		assert.Equal(t, float64(1), response["id"])
		assert.Equal(t, float64(topic.ID), response["topic_id"])
		assert.Equal(t, float64(user.ID), response["user_id"])
		assert.Equal(t, true, response["paid"])

		// Check order details
		order, ok := response["order"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "ORD20250624001", order["order_no"])
		assert.Equal(t, "PAY20250624001", order["out_order_no"])
		assert.Equal(t, float64(1000), order["price"])
		assert.Equal(t, "paid", order["status"])

		// Verify in database
		var dbExperience models.Experience
		db.Preload("Order").First(&dbExperience, experience.ID)
		assert.True(t, dbExperience.Paid())
		assert.NotNil(t, dbExperience.Order)
		assert.Equal(t, models.OrderStatusPaid, dbExperience.Order.Status)
	})

	t.Run("invalid experience ID", func(t *testing.T) {
		requestBody := MarkPaidRequest{
			OrderNo:    "ORD20250624002",
			OutOrderNo: "PAY20250624002",
			Price:      500,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/experiences/invalid/paid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("experience not found", func(t *testing.T) {
		requestBody := MarkPaidRequest{
			OrderNo:    "ORD20250624003",
			OutOrderNo: "PAY20250624003",
			Price:      750,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/experiences/999/paid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"order_no": "ORD20250624004",
			// Missing out_order_no and price
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/experiences/1/paid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("experience already has order", func(t *testing.T) {
		// Create another experience with an existing order
		experience2 := models.Experience{
			TopicID: topic.ID,
			UserID:  user.ID,
		}
		db.Create(&experience2)

		existingOrder := models.Order{
			UserID:       user.ID,
			ExperienceID: experience2.ID,
			Price:        500,
			Status:       models.OrderStatusCreated,
			OrderNo:      "ORD20250624005",
			OutOrderNo:   "PAY20250624005",
		}
		db.Create(&existingOrder)

		requestBody := MarkPaidRequest{
			OrderNo:    "ORD20250624006",
			OutOrderNo: "PAY20250624006",
			Price:      1200,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/experiences/2/paid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
