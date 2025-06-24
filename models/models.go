package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitTestDB() *gorm.DB {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	database.AutoMigrate(&User{}, &Token{}, &Order{})
	return database
}

func SetDB(database *gorm.DB) {
	db = database
}

func GetDB() *gorm.DB {
	return db
}
