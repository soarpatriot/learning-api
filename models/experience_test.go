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
	answerIDs := []uint{1, 2, 3}
	err := e.CreateWithReplies(10, 20, answerIDs)
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

func TestMarkCheckedAnswers(t *testing.T) {
	e := &Experience{
		Replies: []Reply{
			{AnswerID: 2},
			{AnswerID: 4},
		},
		Topic: Topic{
			Questions: []Question{
				{
					Answers: []Answer{
						{ID: 1},
						{ID: 2},
						{ID: 3},
					},
				},
				{
					Answers: []Answer{
						{ID: 4},
						{ID: 5},
					},
				},
			},
		},
	}
	e.MarkCheckedAnswers()

	// Check first question answers
	assert.False(t, e.Topic.Questions[0].Answers[0].Checked) // ID 1
	assert.True(t, e.Topic.Questions[0].Answers[1].Checked)  // ID 2
	assert.False(t, e.Topic.Questions[0].Answers[2].Checked) // ID 3
	// Check second question answers
	assert.True(t, e.Topic.Questions[1].Answers[0].Checked)  // ID 4
	assert.False(t, e.Topic.Questions[1].Answers[1].Checked) // ID 5
}
