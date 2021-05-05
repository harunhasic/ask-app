package repository

import (
	"github.com/mop/entity"
	"gorm.io/gorm"
)

type AnswerRepository interface {
	InsertAnswer(a entity.Answer) entity.Answer
	UpdateAnswer(a entity.Answer) entity.Answer
	DeleteAnswer(a entity.Answer)
	AllAnswers() []entity.Answer
	FindAnswerByID(questionID uint64) entity.Answer
	UsersWithMostAnswers() []entity.Answer
}

type answerConnection struct {
	connection *gorm.DB
}

//NewQuestionRepository creates an instance of the QuestionRepository
func NewAnswerRepository(dbConn *gorm.DB) AnswerRepository {
	return &answerConnection{
		connection: dbConn,
	}
}

func (db *answerConnection) InsertAnswer(a entity.Answer) entity.Answer {
	db.connection.Save(&a)
	db.connection.Preload("User").Find(&a)
	return a
}

func (db *answerConnection) UpdateAnswer(a entity.Answer) entity.Answer {
	db.connection.Save(&a)
	db.connection.Preload("User").Find(&a)
	return a
}

func (db *answerConnection) DeleteAnswer(a entity.Answer) {
	db.connection.Delete(&a)
}

func (db *answerConnection) FindAnswerByID(answerID uint64) entity.Answer {
	var answer entity.Answer
	db.connection.Preload("User").Find(&answer, answerID)
	return answer
}

func (db *answerConnection) AllAnswers() []entity.Answer {
	var answers []entity.Answer
	db.connection.Preload("User").Preload("Question").Find(&answers)
	return answers
}

func (db *answerConnection) UsersWithMostAnswers() []entity.Answer {
	var answers []entity.Answer
	db.connection.Preload("User").Order("user_id desc").Find(&answers)
	return answers
}
