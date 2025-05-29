package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListQuestions(c *gin.Context, db *gorm.DB) {
	var questions []Question
	if err := db.Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

func CreateQuestion(c *gin.Context, db *gorm.DB) {
	var question Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, question)
}

func GetQuestion(c *gin.Context, db *gorm.DB) {
	var question Question
	id := c.Param("id")
	if err := db.First(&question, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, question)
}

func UpdateQuestion(c *gin.Context, db *gorm.DB) {
	var question Question
	id := c.Param("id")
	if err := db.First(&question, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question.ID = 0 // Prevent ID overwrite
	if err := db.Model(&Question{}).Where("id = ?", id).Updates(question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, question)
}

func DeleteQuestion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Delete(&Question{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func GetQuestionsWithAnswers(c *gin.Context, db *gorm.DB) {
	var questions []Question
	topicID := c.Param("id")
	if err := db.Preload("Answers").Where("topic_id = ?", topicID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}
