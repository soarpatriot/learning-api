package main

import (
	"fmt"
	"learning-api/config"
	"learning-api/middlewares"
	"learning-api/models"
	"learning-api/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getDSN(cfg config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLUserName, cfg.MySQLPassword, cfg.MySQLAddress, cfg.MySQLDB)
}

func init() {
	_ = godotenv.Load() // Loads .env from project root if present
}

func main() {
	cfg := config.LoadConfig()
	if cfg.Profile == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	dsn := getDSN(cfg)
	// println("Connecting to database with DSN:", dsn)
	// mask the password in the DSN for security
	linkString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLUserName, "******", cfg.MySQLAddress, cfg.MySQLDB)
	fmt.Println("Connecting to database with DSN:", linkString)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	models.SetDB(db)
	db.AutoMigrate(&models.Topic{}, &models.Question{}, &models.Answer{}, &models.User{}, &models.Token{}, &models.Experience{}, &models.Reply{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r := gin.Default()
	r.Use(middlewares.AuthMiddleware())
	routes.RegisterRoutes(r, db)
	r.Run(":" + port)
}
