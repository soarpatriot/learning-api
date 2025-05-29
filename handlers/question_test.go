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

func setupTestRouterQuestion() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Topic{}, &models.Question{}, &models.Answer{})
	r := gin.Default()
	RegisterQuestionRoutes(r, db)
	return r, db
}

func RegisterQuestionRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/questions", func(c *gin.Context) { ListQuestions(c, db) })
	r.POST("/questions", func(c *gin.Context) { CreateQuestion(c, db) })
	r.GET("/questions/:id", func(c *gin.Context) { GetQuestion(c, db) })
	r.PUT("/questions/:id", func(c *gin.Context) { UpdateQuestion(c, db) })
	r.DELETE("/questions/:id", func(c *gin.Context) { DeleteQuestion(c, db) })
	r.GET("/topics/:id/questions-answers", func(c *gin.Context) { GetQuestionsWithAnswers(c, db) })
}

func TestCreateQuestion(t *testing.T) {
	r, db := setupTestRouterQuestion()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	body := `{"content":"Q1","weight":1,"topic_id":1}`
	req, _ := http.NewRequest("POST", "/questions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestGetQuestion(t *testing.T) {
	r, db := setupTestRouterQuestion()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := models.Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	req, _ := http.NewRequest("GET", "/questions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateQuestion(t *testing.T) {
	r, db := setupTestRouterQuestion()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := models.Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	body := `{"content":"Q2","weight":2,"topic_id":1}`
	req, _ := http.NewRequest("PUT", "/questions/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteQuestion(t *testing.T) {
	r, db := setupTestRouterQuestion()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := models.Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	req, _ := http.NewRequest("DELETE", "/questions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestGetQuestionsWithAnswers(t *testing.T) {
	r, db := setupTestRouterQuestion()
	topic := models.Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := models.Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	answer := models.Answer{Content: "A1", Correct: true, QuestionID: 1}
	db.Create(&answer)
	req, _ := http.NewRequest("GET", "/topics/1/questions-answers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
