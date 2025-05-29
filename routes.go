package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/topics", func(c *gin.Context) { ListTopics(c, db) })
	r.POST("/topics", func(c *gin.Context) { CreateTopic(c, db) })
	r.GET("/topics/:id", func(c *gin.Context) { GetTopic(c, db) })
	r.PUT("/topics/:id", func(c *gin.Context) { UpdateTopic(c, db) })
	r.DELETE("/topics/:id", func(c *gin.Context) { DeleteTopic(c, db) })

	r.GET("/questions", func(c *gin.Context) { ListQuestions(c, db) })
	r.POST("/questions", func(c *gin.Context) { CreateQuestion(c, db) })
	r.GET("/questions/:id", func(c *gin.Context) { GetQuestion(c, db) })
	r.PUT("/questions/:id", func(c *gin.Context) { UpdateQuestion(c, db) })
	r.DELETE("/questions/:id", func(c *gin.Context) { DeleteQuestion(c, db) })

	r.GET("/answers", func(c *gin.Context) { ListAnswers(c, db) })
	r.POST("/answers", func(c *gin.Context) { CreateAnswer(c, db) })
	r.GET("/answers/:id", func(c *gin.Context) { GetAnswer(c, db) })
	r.PUT("/answers/:id", func(c *gin.Context) { UpdateAnswer(c, db) })
	r.DELETE("/answers/:id", func(c *gin.Context) { DeleteAnswer(c, db) })

	r.GET("/topics/:id/questions-answers", func(c *gin.Context) { GetQuestionsWithAnswers(c, db) })
}
