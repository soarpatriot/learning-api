package models

import "time"

// User represents a user entity
type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OpenID     string    `json:"open_id"`
	UnionID    string    `json:"union_id"`
	SessionKey string    `json:"session_key"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Avatar     string    `json:"avatar"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Tokens     []Token   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokens"`
}
