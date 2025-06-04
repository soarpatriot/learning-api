package models

import "time"

// Answer represents an answer entity
type Answer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Content    string    `gorm:"size:1000" json:"content"`
	Correct    bool      `json:"correct"`
	QuestionID uint      `json:"question_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
