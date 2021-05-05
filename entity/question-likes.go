package entity

type QuestionLikes struct {
	ID         uint64   `gorm:"primary_key:auto_increment" json:"id"`
	QuestionID uint64   `gorm:"not null" json:"question_id"`
	Question   Question `gorm:"foreignkey:QuestionID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
	UserID     uint64   `gorm:"not null" json:"user_id"`
	User       User     `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
}
