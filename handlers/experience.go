package handlers

import (
	"learning-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ExperienceRequest struct {
	TopicID   uint   `json:"topic_id"`
	AnswerIDs []uint `json:"answer_ids"`
}

func CreateExperience(c *gin.Context) {
	var req ExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate required fields
	if req.TopicID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic_id is required"})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	experience := models.Experience{}
	err := experience.CreateWithReplies(req.TopicID, currentUser.(models.User).ID, req.AnswerIDs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Experience and replies created successfully"})
}

func GetExperience(c *gin.Context) {
	db := models.GetDB()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid experience id"})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(models.User)

	var experience models.Experience
	err = db.Preload("Replies").Preload("User").Preload("Topic.Questions.Answers").First(&experience, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "experience not found"})
		return
	}

	if experience.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not your experience"})
		return
	}

	type experienceResponse struct {
		ID        uint           `json:"id"`
		TopicID   uint           `json:"topic_id"`
		UserID    uint           `json:"user_id"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		Replies   []models.Reply `json:"replies"`
		Topic     models.Topic   `json:"topic"`
	}

	experience.MarkCheckedAnswers()

	resp := experienceResponse{
		ID:        experience.ID,
		TopicID:   experience.TopicID,
		UserID:    experience.UserID,
		CreatedAt: experience.CreatedAt,
		UpdatedAt: experience.UpdatedAt,
		Replies:   experience.Replies,
		Topic:     experience.Topic,
	}

	c.JSON(http.StatusOK, resp)
}
