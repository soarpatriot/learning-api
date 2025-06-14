package handlers

import (
	"bytes"
	"encoding/json"
	"learning-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockUser struct {
	ID uint
}

func setupTestRouterExperience() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Experience{}, &models.Reply{}, &models.User{})
	models.SetDB(db)
	r := gin.Default()
	r.POST("/experience", func(c *gin.Context) {
		// Simulate authentication middleware
		c.Set("currentUser", models.User{ID: 42})
		CreateExperience(c)
	})
	return r, db
}

func TestCreateExperience_Success(t *testing.T) {
	r, db := setupTestRouterExperience()
	// Insert a user for foreign key
	db.Create(&models.User{ID: 42})
	body, _ := json.Marshal(map[string]interface{}{
		"topic_id":   1,
		"answer_ids": []uint{2, 3},
	})
	req, _ := http.NewRequest("POST", "/experience", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Experience and replies created successfully")) {
		t.Errorf("Expected success message, got %s", w.Body.String())
	}
}

func TestCreateExperience_BadRequest(t *testing.T) {
	r, _ := setupTestRouterExperience()
	// Missing topic_id
	body, _ := json.Marshal(map[string]interface{}{})
	req, _ := http.NewRequest("POST", "/experience", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateExperience_Unauthorized(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Experience{}, &models.Reply{}, &models.User{})
	models.SetDB(db)
	r := gin.Default()
	r.POST("/experience", func(c *gin.Context) {
		// Do not set currentUser
		CreateExperience(c)
	})
	body, _ := json.Marshal(map[string]interface{}{
		"topic_id":   1,
		"answer_ids": []uint{2, 3},
	})
	req, _ := http.NewRequest("POST", "/experience", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
