package handlers

import (
	"learning-api/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouterTopic() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Topic{}, &models.Question{}, &models.Answer{})
	r := gin.Default()
	RegisterTopicRoutes(r, db)
	return r, db
}

func RegisterTopicRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/topics", func(c *gin.Context) { ListTopics(c, db) })
	r.POST("/topics", func(c *gin.Context) { CreateTopic(c, db) })
	r.GET("/topics/:id", func(c *gin.Context) { GetTopic(c, db) })
	r.PUT("/topics/:id", func(c *gin.Context) { UpdateTopic(c, db) })
	r.DELETE("/topics/:id", func(c *gin.Context) { DeleteTopic(c, db) })
}

func TestListTopics(t *testing.T) {
	r, _ := setupTestRouterTopic()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateTopic(t *testing.T) {
	r, _ := setupTestRouterTopic()
	body := `{"name":"Test Topic","description":"Desc","explaination":"Expl"}`
	req, _ := http.NewRequest("POST", "/topics", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestGetTopic(t *testing.T) {
	r, db := setupTestRouterTopic()
	// convert 2025-06-13T22:09:20.750225+08:00 to time.Time
	layout := "2006-01-02T15:04:05.999999999-07:00"
	createdAt, _ := time.Parse(layout, "2025-06-13T22:09:20.750225+08:00")
	updatedAt, _ := time.Parse(layout, "2025-06-13T22:09:20.750225+08:00")

	topic := models.Topic{Name: "Test", Description: "D", Explaination: "E", CreatedAt: createdAt, UpdatedAt: updatedAt}
	db.Create(&topic)

	question := models.Question{Content: "Test Question", TopicID: topic.ID, CreatedAt: createdAt, UpdatedAt: updatedAt}
	db.Create(&question)

	answer := models.Answer{Content: "Test Answer", QuestionID: question.ID, CreatedAt: createdAt, UpdatedAt: updatedAt}
	db.Create(&answer)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	// only check the name, description, explaination, and questions, and answers content
	// as the created_at and updated_at are dynamic
	if !strings.Contains(w.Body.String(), `"name":"Test"`) ||
		!strings.Contains(w.Body.String(), `"description":"D"`) ||
		!strings.Contains(w.Body.String(), `"explaination":"E"`) ||
		!strings.Contains(w.Body.String(), `"content":"Test Question"`) ||
		!strings.Contains(w.Body.String(), `"content":"Test Answer"`) {
		t.Errorf("Response body does not contain expected topic details")
	}
}

func TestUpdateTopic(t *testing.T) {
	r, db := setupTestRouterTopic()
	topic := models.Topic{Name: "Test", Description: "D", Explaination: "E"}
	db.Create(&topic)
	body := `{"name":"Updated","description":"D2","explaination":"E2"}`
	req, _ := http.NewRequest("PUT", "/topics/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteTopic(t *testing.T) {
	r, db := setupTestRouterTopic()
	topic := models.Topic{Name: "Test", Description: "D", Explaination: "E"}
	db.Create(&topic)
	req, _ := http.NewRequest("DELETE", "/topics/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}
