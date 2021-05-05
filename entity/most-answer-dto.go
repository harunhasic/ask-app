package entity

type MostAnswersDTO struct {
	UserID          int    `json:"user_id, omitempty"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	NumberOfAnswers int    `json:"number_of_answers"`
}
