package dto

//RegisterDTO is a model thats used when client posts from the /register url
//We make sure that the mail must be an email and that the password is minimum 5 char long
type RegisterDTO struct {
	Firstname string `json:"firstname" form:"name"`
	Lastname  string `json:"lastname" form:"name"`
	Email     string `json:"email" form:"email" binding:"required,email" validate:"email"`
	Password  string `json:"password" form:"password" binding:"required,min=5"`
}
