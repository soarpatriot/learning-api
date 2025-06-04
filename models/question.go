package models

import "time"

// Question represents a question entity
type Question struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"size:1000" json:"content"`
	Weight    int       `json:"weight"`
	TopicID   uint      `json:"topic_id"`
	Answers   []Answer  `json:"answers"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
