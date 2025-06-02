package main

import (
	"fmt"
	"learning-api/config"
	"learning-api/models"
	"learning-api/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getDSN(cfg config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLUser, cfg.MySQLPassword, cfg.MySQLHost, cfg.MySQLDB)
}

func main() {
	cfg := config.LoadConfig()
	if cfg.Profile == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	dsn := getDSN(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Topic{}, &models.Question{}, &models.Answer{})

	// Insert example data if tables are empty
	var count int64
	db.Model(&models.Topic{}).Count(&count)
	if count == 0 {
		topic := models.Topic{Name: "Go Basics", Description: "Learn Go basics", Explaination: "Covers variables, loops, etc."}
		db.Create(&topic)
		question := models.Question{Content: "What is a goroutine?", Weight: 1, TopicID: topic.ID}
		db.Create(&question)
		answer1 := models.Answer{Content: "A lightweight thread managed by Go runtime", Correct: true, QuestionID: question.ID}
		answer2 := models.Answer{Content: "A type of variable", Correct: false, QuestionID: question.ID}
		db.Create(&answer1)
		db.Create(&answer2)
	}

	r := gin.Default()
	routes.RegisterRoutes(r, db)
	r.Run()
}
