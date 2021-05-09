package dto

import "time"

//QuestionUpdateDTO is the model that the client uses when we want to update a question
type QuestionUpdateDTO struct {
	ID         uint64    `json:"id" form:"id" binding:"required"`
	Body       string    `json:"body" form:"body" binding:"required"`
	NumOfLikes int       `json:"num_of_likes" form:"num_of_likes"`
	UserID     uint64    `json:"user_id,omitempty" form:"user_id,omitempty"`
	CreatedAt  time.Time `json:"created"`
	UpdatedAt  time.Time `json:"updated"`
}

//QuestionCreateDTO is the model that the client uses when we create a new question
type QuestionCreateDTO struct {
	Body       string    `json:"body" form:"body" binding:"required"`
	NumOfLikes int       `json:"num_of_likes" form:"num_of_likes"`
	UserID     uint64    `json:"user_id,omitempty" form:"user_id,omitempty"`
	CreatedAt  time.Time `json:"created"`
	UpdatedAt  time.Time `json:"updated"`
}
