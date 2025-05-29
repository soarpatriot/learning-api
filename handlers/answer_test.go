package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"learning-api/models"
)

func setupTestRouterAnswer() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Topic{}, &models.Question{}, &models.Answer{})
	r := gin.Default()
	RegisterAnswerRoutes(r, db)
	return r, db
}

func RegisterAnswerRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/answers", func(c *gin.Context) { ListAnswers(c, db) })
	r.POST("/answers", func(c *gin.Context) { CreateAnswer(c, db) })
	r.GET("/answers/:id", func(c *gin.Context) { GetAnswer(c, db) })
	r.PUT("/answers/:id", func(c *gin.Context) { UpdateAnswer(c, db) })
	r.DELETE("/answers/:id", func(c *gin.Context) { DeleteAnswer(c, db) })
}

func TestCreateAnswer(t *testing.T) {
	r, db := setupTestRouterAnswer()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := models.Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	body := `{"content":"A1","correct":true,"question_id":1}`
	req, _ := http.NewRequest("POST", "/answers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestGetAnswer(t *testing.T) {
	r, db := setupTestRouterAnswer()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := models.Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	answer := models.Answer{Content: "A1", Correct: true, QuestionID: 1}
	db.Create(&answer)
	req, _ := http.NewRequest("GET", "/answers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateAnswer(t *testing.T) {
	r, db := setupTestRouterAnswer()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := models.Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	answer := models.Answer{Content: "A1", Correct: true, QuestionID: 1}
	db.Create(&answer)
	body := `{"content":"A2","correct":false,"question_id":1}`
	req, _ := http.NewRequest("PUT", "/answers/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteAnswer(t *testing.T) {
	r, db := setupTestRouterAnswer()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := models.Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	answer := models.Answer{Content: "A1", Correct: true, QuestionID: 1}
	db.Create(&answer)
	req, _ := http.NewRequest("DELETE", "/answers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}
