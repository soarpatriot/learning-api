package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Experience struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TopicID   uint      `gorm:"not null" json:"topic_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Replies   []Reply   `gorm:"foreignKey:ExperienceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"replies"`
	Order     *Order    `gorm:"foreignKey:ExperienceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order"`
	User      User      `json:"user"`
	Topic     Topic     `json:"topic"`
}

// IsPaid returns true if there's an order with status paid or confirmed
func (e *Experience) IsPaid() bool {
	if e.Order != nil {
		return e.Order.Status == OrderStatusPaid || e.Order.Status == OrderStatusConfirmed
	}
	return false
}

// Paid returns the paid status as a virtual property for JSON serialization
func (e *Experience) Paid() bool {
	return e.IsPaid()
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

func (e *Experience) TimeAgoZh() string {
	delta := time.Since(e.CreatedAt)
	if delta < time.Minute {
		return "刚刚"
	} else if delta < time.Hour {
		return fmt.Sprintf("%d 分钟前", int(delta.Minutes()))
	} else if delta < 24*time.Hour {
		return fmt.Sprintf("%d 小时前", int(delta.Hours()))
	} else if delta < 30*24*time.Hour {
		return fmt.Sprintf("%d 天前", int(delta.Hours()/24))
	} else if delta < 12*30*24*time.Hour {
		return fmt.Sprintf("%d 个月前", int(delta.Hours()/(24*30)))
	}
	return fmt.Sprintf("%d 年前", int(delta.Hours()/(24*365)))
}

type MyExperienceResponse struct {
	ID        uint      `json:"id"`
	TopicID   uint      `json:"topic_id"`
	UserID    uint      `json:"user_id"`
	Paid      bool      `json:"paid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Replies   []Reply   `json:"replies"`
	Topic     Topic     `json:"topic"`
	TimeAgoZh string    `json:"time_ago_zh"`
}

func ToMyExperienceResponses(experiences []Experience) []MyExperienceResponse {
	resp := make([]MyExperienceResponse, 0, len(experiences))
	for _, exp := range experiences {
		resp = append(resp, MyExperienceResponse{
			ID:        exp.ID,
			TopicID:   exp.TopicID,
			UserID:    exp.UserID,
			Paid:      exp.Paid(),
			CreatedAt: exp.CreatedAt,
			UpdatedAt: exp.UpdatedAt,
			Replies:   exp.Replies,
			Topic:     exp.Topic,
			TimeAgoZh: exp.TimeAgoZh(),
		})
	}
	return resp
}
