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
	Replies   []Reply   `gorm:"foreignKey:ExperienceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"replies"`
	User      User
	Topic     Topic
}

func (e *Experience) CreateWithReplies(topicID uint, userID uint, answerIds []uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		e.TopicID = topicID
		e.UserID = userID

		var replies []Reply
		for _, answerID := range answerIds {
			reply := &Reply{
				AnswerID: answerID,
			}
			replies = append(replies, *reply)
		}
		e.Replies = replies
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		return nil
	})
}

func (e *Experience) MarkCheckedAnswers() {
	checkedMap := map[uint]struct{}{}
	for _, reply := range e.Replies {
		checkedMap[reply.AnswerID] = struct{}{}
	}
	for qi := range e.Topic.Questions {
		for ai := range e.Topic.Questions[qi].Answers {
			ans := &e.Topic.Questions[qi].Answers[ai]
			if _, ok := checkedMap[ans.ID]; ok {
				ans.Checked = true
			} else {
				ans.Checked = false
			}
		}
	}
}
