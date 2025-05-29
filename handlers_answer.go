package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListAnswers(c *gin.Context, db *gorm.DB) {
	var answers []Answer
	if err := db.Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, answers)
}

func CreateAnswer(c *gin.Context, db *gorm.DB) {
	var answer Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, answer)
}

func GetAnswer(c *gin.Context, db *gorm.DB) {
	var answer Answer
	id := c.Param("id")
	if err := db.First(&answer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}
	c.JSON(http.StatusOK, answer)
}

func UpdateAnswer(c *gin.Context, db *gorm.DB) {
	var answer Answer
	id := c.Param("id")
	if err := db.First(&answer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer.ID = 0 // Prevent ID overwrite
	if err := db.Model(&Answer{}).Where("id = ?", id).Updates(answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, answer)
}

func DeleteAnswer(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Delete(&Answer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
