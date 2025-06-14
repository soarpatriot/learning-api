package handlers

import (
	"learning-api/models"
	"net/http"

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
