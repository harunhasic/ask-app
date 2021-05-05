package entity

//User represents the 'user' table in database
type User struct {
	ID        uint64      `gorm:"primary_key:auto_increment" json:"id"`
	Firstname string      `gorm:"type:varchar(255)" json:"firstname"`
	Lastname  string      `gorm:"type:varchar(255)" json:"lastname"`
	Email     string      `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password  string      `gorm:"->;<-;not null" json:"-"`
	Token     string      `gorm:"-" json:"token,omitempty"`
	Questions *[]Question `json:"questions,omitempty"`
	Answers   *[]Answer   `gorm:"foreignKey:UserID" json:"answers,omitempty"`
}
