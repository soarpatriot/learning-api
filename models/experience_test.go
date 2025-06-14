package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestExperience_CreateWithReplies(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Experience{}, &Reply{})
	SetDB(db)

	e := &Experience{}
	replies := []uint{1, 2, 3}
	err := e.CreateWithReplies(10, 20, replies)
	assert.NoError(t, err)
	assert.Equal(t, uint(10), e.TopicID)
	assert.Equal(t, uint(20), e.UserID)

	var count int64
	db.Model(&Experience{}).Count(&count)
	assert.Equal(t, int64(1), count)
	db.Model(&Reply{}).Count(&count)
	assert.Equal(t, int64(3), count)
}

func TestExperience_CreateWithReplies_CreateError(t *testing.T) {
	// To simulate error, use a closed DB or invalid table
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	SetDB(db)
	// Don't migrate tables to force error
	e := &Experience{}
	replies := []uint{1, 2}
	err := e.CreateWithReplies(10, 20, replies)
	assert.Error(t, err)
}
