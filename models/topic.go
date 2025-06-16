package models

import "time"

// Topic represents a topic entity
type Topic struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Explaination string     `json:"explaination"`
	Questions    []Question `json:"questions"`
	CoverURL     string     `gorm:"type:varchar(1000)" json:"cover_url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
