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

func TestGetExperience_Success(t *testing.T) {
	r, _, _, _, answers := setupGetExperienceTestDB()
	req, _ := http.NewRequest("GET", "/experience/11", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var resp struct {
		ID    uint `json:"id"`
		Topic struct {
			Questions []struct {
				Answers []struct {
					ID      uint `json:"id"`
					Checked bool `json:"checked"`
				} `json:"answers"`
			} `json:"questions"`
		} `json:"topic"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if resp.ID != 11 {
		t.Errorf("Expected experience ID 11, got %d", resp.ID)
	}
	// Check checked property
	found := false
	for _, q := range resp.Topic.Questions {
		for _, a := range q.Answers {
			if a.ID == answers[1].ID && a.Checked != true {
				t.Errorf("Expected answer %d checked true", a.ID)
			}
			if a.ID == answers[0].ID && a.Checked != false {
				t.Errorf("Expected answer %d checked false", a.ID)
			}
			if a.ID == answers[1].ID {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("Expected answer with ID %d in response", answers[1].ID)
	}
}

func TestGetExperience_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Topic{}, &models.Question{}, &models.Answer{}, &models.Experience{}, &models.Reply{})
	models.SetDB(db)
	user := models.User{ID: 1, Name: "testuser"}
	db.Create(&user)
	otherUser := models.User{ID: 2, Name: "other"}
	db.Create(&otherUser)
	topic := models.Topic{ID: 1, Name: "topic1"}
	db.Create(&topic)
	exp := models.Experience{ID: 12, TopicID: topic.ID, UserID: user.ID}
	db.Create(&exp)
	r := gin.Default()
	r.GET("/experience/:id", func(c *gin.Context) {
		c.Set("currentUser", otherUser)
		GetExperience(c)
	})
	req, _ := http.NewRequest("GET", "/experience/12", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403, got %d", w.Code)
	}
}

func setupGetExperienceTestDB() (*gin.Engine, *gorm.DB, models.User, models.Experience, []models.Answer) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Topic{}, &models.Question{}, &models.Answer{}, &models.Experience{}, &models.Reply{})
	models.SetDB(db)

	user := models.User{ID: 1, Name: "testuser"}
	db.Create(&user)
	topic := models.Topic{ID: 1, Name: "topic1"}
	db.Create(&topic)
	question := models.Question{ID: 1, TopicID: topic.ID, Content: "q1"}
	db.Create(&question)
	answers := []models.Answer{
		{ID: 1, QuestionID: question.ID, Content: "a1"},
		{ID: 2, QuestionID: question.ID, Content: "a2"},
	}
	for _, a := range answers {
		db.Create(&a)
	}
	exp := models.Experience{ID: 11, TopicID: topic.ID, UserID: user.ID}
	db.Create(&exp)
	reply := models.Reply{ExperienceID: exp.ID, AnswerID: 2}
	db.Create(&reply)

	r := gin.Default()
	r.GET("/experience/:id", func(c *gin.Context) {
		c.Set("currentUser", user)
		GetExperience(c)
	})
	return r, db, user, exp, answers
}

func TestGetMyExperiences(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Topic{}, &models.Experience{})
	models.SetDB(db)

	user := models.User{ID: 1, Name: "testuser"}
	db.Create(&user)
	topic := models.Topic{ID: 1, Name: "topic1"}
	db.Create(&topic)
	exp1 := models.Experience{ID: 1, TopicID: topic.ID, UserID: user.ID}
	exp2 := models.Experience{ID: 2, TopicID: topic.ID, UserID: user.ID}
	db.Create(&exp1)
	db.Create(&exp2)

	r := gin.Default()
	r.GET("/experiences/my", func(c *gin.Context) {
		c.Set("currentUser", user)
		GetMyExperiences(c)
	})

	req, _ := http.NewRequest("GET", "/experiences/my", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}

	var resp []struct {
		ID    uint `json:"id"`
		Topic struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"topic"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("Expected 2 experiences, got %d", len(resp))
	}
	for _, e := range resp {
		if e.Topic.ID != topic.ID || e.Topic.Name != topic.Name {
			t.Errorf("Expected topic to be preloaded with correct data")
		}
	}
}
