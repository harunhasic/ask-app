package dto

//UserUpdateDTO is used by client for the update/PUT request
//We make sure that the mail must be an email and that the password is minimum 5 char long
type UserUpdateDTO struct {
	ID        uint64 `json:"id" form:"id"`
	Firstname string `json:"firstname" form:"firstname"`
	Lastname  string `json:"lastname" form:"lastname"`
	Email     string `json:"email" form:"email" binding:"required,email"`
	Password  string `json:"password,omitempty" form:"password,omitempty" binding:"min=5"`
}
