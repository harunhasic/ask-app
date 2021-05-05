package entity

type QuestionResult struct {
	QuestionID int    `json:"id, omitempty"`
	Body       string `json:"body"`
	NumOfLikes int    `json:"num_of_likes"`
	IsLiked    bool   `json:"is_liked"`
	IsEditable bool   `json:"is_editable"`
}
