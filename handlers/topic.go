package handlers

import (
	"learning-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListTopics(c *gin.Context, db *gorm.DB) {
	var topics []models.Topic
	if err := db.Find(&topics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, topics)
}

func CreateTopic(c *gin.Context, db *gorm.DB) {
	var topic models.Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&topic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, topic)
}

func GetTopic(c *gin.Context, db *gorm.DB) {
	var topic models.Topic
	id := c.Param("id")
	if err := db.Preload("Questions.Answers").First(&topic, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	c.JSON(http.StatusOK, topic)
}

func UpdateTopic(c *gin.Context, db *gorm.DB) {
	var topic models.Topic
	id := c.Param("id")
	if err := db.First(&topic, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topic.ID = 0 // Prevent ID overwrite
	if err := db.Model(&models.Topic{}).Where("id = ?", id).Updates(topic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, topic)
}

func DeleteTopic(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Delete(&models.Topic{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
