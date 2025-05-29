package main

// Topic represents a topic entity
type Topic struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Explaination string `json:"explaination"`
	Questions   []Question `json:"questions"`
}

// Question represents a question entity
type Question struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Content   string `gorm:"size:1000" json:"content"`
	Weight    int    `json:"weight"`
	TopicID   uint   `json:"topic_id"`
	Answers   []Answer `json:"answers"`
}

// Answer represents an answer entity
type Answer struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Content    string `gorm:"size:1000" json:"content"`
	Correct    bool   `json:"correct"`
	QuestionID uint   `json:"question_id"`
}
