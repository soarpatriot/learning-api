package models

import (
	"time"
)

type Reply struct {
	ID           uint      `gorm:"primaryKey"`
	ExperienceID uint      `gorm:"not null"`
	AnswerID     uint      `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
