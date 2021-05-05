package dto

//QuestionUpdateDTO is the model that the client uses when we want to update a question
type AnswerUpdateDTO struct {
	ID         uint64 `json:"id" form:"id" binding:"required"`
	Body       string `json:"body" form:"body" binding:"required"`
	QuestionID uint64 `json:"question_id,omitempty" form:"question_id,omitempty"`
	UserID     uint64 `json:"user_id,omitempty" form:"user_id,omitempty"`
}

//QuestionCreateDTO is the model that the client uses when we create a new question
type AnswerCreateDTO struct {
	Body       string `json:"body" form:"body" binding:"required"`
	QuestionID uint64 `json:"question_id,omitempty" form:"question_id,omitempty" binding:"required"`
	UserID     uint64 `json:"user_id,omitempty" form:"user_id,omitempty"`
}
