package models

import (
	"time"

	"gorm.io/gorm"
)

type Experience struct {
	ID        uint      `gorm:"primaryKey"`
	TopicID   uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (e *Experience) CreateWithReplies(topicID uint, userID uint, replies []uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		e.TopicID = topicID
		e.UserID = userID

		if err := tx.Create(e).Error; err != nil {
			return err
		}

		for _, answerID := range replies {
			reply := Reply{
				ExperienceID: e.ID,
				AnswerID:     answerID,
			}
			if err := tx.Create(&reply).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
