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

	c.JSON(http.StatusOK, experience)
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
	err = db.Preload("Replies").Preload("User").Preload("Order").Preload("Topic.Questions.Answers").First(&experience, id).Error
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
		Paid      bool           `json:"paid"`
		Order     *models.Order  `json:"order"`
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
		Paid:      experience.Paid(),
		Order:     experience.Order,
	}

	c.JSON(http.StatusOK, resp)
}

func GetMyExperiences(c *gin.Context) {
	db := models.GetDB()
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(models.User)

	var experiences []models.Experience
	err := db.Preload("Topic").Where("user_id = ?", user.ID).Find(&experiences).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := models.ToMyExperienceResponses(experiences)
	c.JSON(http.StatusOK, resp)
}

type MarkPaidRequest struct {
	OrderNo    string `json:"order_no" binding:"required"`
	OutOrderNo string `json:"out_order_no" binding:"required"`
	Price      int    `json:"price" binding:"required"`
}

func MarkExperiencePaid(c *gin.Context) {
	db := models.GetDB()

	// Get experience ID from URL parameter
	idStr := c.Param("id")
	experienceID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid experience id"})
		return
	}

	// Load current user
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(models.User)

	// Parse request body
	var req MarkPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Load the experience
	var experience models.Experience
	err = db.Preload("Order").First(&experience, experienceID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "experience not found"})
		return
	}

	// Check if user owns this experience
	if experience.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not your experience"})
		return
	}

	// Check if experience already has an order
	if experience.Order != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "experience already has an order"})
		return
	}

	// Create new order with paid status
	order := models.Order{
		UserID:       user.ID,
		ExperienceID: uint(experienceID),
		Price:        req.Price,
		Status:       models.OrderStatusPaid,
		OrderNo:      req.OrderNo,
		OutOrderNo:   req.OutOrderNo,
	}

	// Save the order
	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order: " + err.Error()})
		return
	}

	// Reload experience with the new order
	err = db.Preload("Order").First(&experience, experienceID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload experience"})
		return
	}

	// Return the updated experience with paid status
	response := gin.H{
		"id":         experience.ID,
		"topic_id":   experience.TopicID,
		"user_id":    experience.UserID,
		"paid":       experience.Paid(),
		"created_at": experience.CreatedAt,
		"updated_at": experience.UpdatedAt,
		"order": gin.H{
			"id":           order.ID,
			"order_no":     order.OrderNo,
			"out_order_no": order.OutOrderNo,
			"price":        order.Price,
			"status":       order.Status.String(),
			"created_at":   order.CreatedAt,
			"updated_at":   order.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}
