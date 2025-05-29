package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"learning-api/handlers"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/topics", func(c *gin.Context) { handlers.ListTopics(c, db) })
	r.POST("/topics", func(c *gin.Context) { handlers.CreateTopic(c, db) })
	r.GET("/topics/:id", func(c *gin.Context) { handlers.GetTopic(c, db) })
	r.PUT("/topics/:id", func(c *gin.Context) { handlers.UpdateTopic(c, db) })
	r.DELETE("/topics/:id", func(c *gin.Context) { handlers.DeleteTopic(c, db) })

	r.GET("/questions", func(c *gin.Context) { handlers.ListQuestions(c, db) })
	r.POST("/questions", func(c *gin.Context) { handlers.CreateQuestion(c, db) })
	r.GET("/questions/:id", func(c *gin.Context) { handlers.GetQuestion(c, db) })
	r.PUT("/questions/:id", func(c *gin.Context) { handlers.UpdateQuestion(c, db) })
	r.DELETE("/questions/:id", func(c *gin.Context) { handlers.DeleteQuestion(c, db) })

	r.GET("/answers", func(c *gin.Context) { handlers.ListAnswers(c, db) })
	r.POST("/answers", func(c *gin.Context) { handlers.CreateAnswer(c, db) })
	r.GET("/answers/:id", func(c *gin.Context) { handlers.GetAnswer(c, db) })
	r.PUT("/answers/:id", func(c *gin.Context) { handlers.UpdateAnswer(c, db) })
	r.DELETE("/answers/:id", func(c *gin.Context) { handlers.DeleteAnswer(c, db) })

	r.GET("/topics/:id/questions-answers", func(c *gin.Context) { handlers.GetQuestionsWithAnswers(c, db) })
}
