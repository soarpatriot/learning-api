package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Topic{}, &Question{}, &Answer{})
	r := gin.Default()
	RegisterRoutes(r, db)
	return r, db
}

func TestListTopics(t *testing.T) {
	r, _ := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateTopic(t *testing.T) {
	r, _ := setupTestRouter()
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
	r, db := setupTestRouter()
	topic := Topic{Name: "Test", Description: "D", Explaination: "E"}
	db.Create(&topic)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateTopic(t *testing.T) {
	r, db := setupTestRouter()
	topic := Topic{Name: "Test", Description: "D", Explaination: "E"}
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
	r, db := setupTestRouter()
	topic := Topic{Name: "Test", Description: "D", Explaination: "E"}
	db.Create(&topic)
	req, _ := http.NewRequest("DELETE", "/topics/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestCreateQuestion(t *testing.T) {
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
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
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	req, _ := http.NewRequest("GET", "/questions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateQuestion(t *testing.T) {
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := Question{Content: "Q1", Weight: 1, TopicID: 1}
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
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	req, _ := http.NewRequest("DELETE", "/questions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestCreateAnswer(t *testing.T) {
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := Question{Content: "Q1", Weight: 1, TopicID: 1}
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
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	answer := Answer{Content: "A1", Correct: true, QuestionID: 1}
	db.Create(&answer)
	req, _ := http.NewRequest("GET", "/answers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateAnswer(t *testing.T) {
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	answer := Answer{Content: "A1", Correct: true, QuestionID: 1}
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
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	answer := Answer{Content: "A1", Correct: true, QuestionID: 1}
	db.Create(&answer)
	req, _ := http.NewRequest("DELETE", "/answers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestGetQuestionsWithAnswers(t *testing.T) {
	r, db := setupTestRouter()
	topic := Topic{Name: "T", Description: "D", Explaination: "E"}
	db.Create(&topic)
	question := Question{Content: "Q1", Weight: 1, TopicID: 1}
	db.Create(&question)
	answer := Answer{Content: "A1", Correct: true, QuestionID: 1}
	db.Create(&answer)
	req, _ := http.NewRequest("GET", "/topics/1/questions-answers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}