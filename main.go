package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"fmt"
)

func getDSN() string {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	db := os.Getenv("MYSQL_DB")
	if user == "" { user = "root" }
	if pass == "" { pass = "22143521" }
	if host == "" { host = "127.0.0.1:3306" }
	if db == "" { db = "learning" }
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, db)
}

func main() {
	profile := os.Getenv("PROFILE")
	if profile == "" { profile = "dev" }
	if profile == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	dsn := getDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Topic{}, &Question{}, &Answer{})

	// Insert example data if tables are empty
	var count int64
	db.Model(&Topic{}).Count(&count)
	if count == 0 {
		topic := Topic{Name: "Go Basics", Description: "Learn Go basics", Explaination: "Covers variables, loops, etc."}
		db.Create(&topic)
		question := Question{Content: "What is a goroutine?", Weight: 1, TopicID: topic.ID}
		db.Create(&question)
		answer1 := Answer{Content: "A lightweight thread managed by Go runtime", Correct: true, QuestionID: question.ID}
		answer2 := Answer{Content: "A type of variable", Correct: false, QuestionID: question.ID}
		db.Create(&answer1)
		db.Create(&answer2)
	}

	r := gin.Default()
	RegisterRoutes(r, db)
	r.Run()
}
