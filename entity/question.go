package entity

import "time"

//Question represents the question table in database
type Question struct {
	ID            uint64           `gorm:"primary_key:auto_increment" json:"id"`
	Body          string           `gorm:"type:text" json:"body"`
	UserID        uint64           `gorm:"not null" json:"user_id"`
	User          User             `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
	NumOfLikes    int              `json:"num_of_likes"`
	Answers       *[]Answer        `json:"answers"`
	QuestionLikes *[]QuestionLikes `json:"question_likes,omitempty"`
	IsLiked       bool             `json:"is_liked"`
	CreatedAt     time.Time        `json:"created"`
	UpdatedAt     time.Time        `json:"updated"`
}
