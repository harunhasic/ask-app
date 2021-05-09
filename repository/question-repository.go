package repository

import (
	"github.com/mop/entity"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	InsertQuestion(q entity.Question) entity.Question
	UpdateQuestion(q entity.Question) entity.Question
	DeleteQuestion(q entity.Question)
	AllQuestions(question *entity.Question, pagination *entity.Pagination) []entity.Question
	FindQuestionByID(questionID uint64) entity.Question
	GetNumberOfLikesForQuestion() []entity.Question
	Like(questionID uint64, userID uint64)
	DeleteLike(questionID uint64, userID uint64)
	QuestionPage(questionID uint64, userID uint64) entity.QuestionResult
	GetUserQuestions(userId uint64, question *entity.Question, pagination *entity.Pagination) []entity.Question
}

type questionConnection struct {
	connection *gorm.DB
}

//NewQuestionRepository creates an instance of the QuestionRepository
func NewQuestionRepository(dbConn *gorm.DB) QuestionRepository {
	return &questionConnection{
		connection: dbConn,
	}
}

func (db *questionConnection) InsertQuestion(q entity.Question) entity.Question {
	db.connection.Save(&q)
	db.connection.Preload("User").Find(&q)
	return q
}

func (db *questionConnection) UpdateQuestion(q entity.Question) entity.Question {
	db.connection.Save(&q)
	db.connection.Preload("User").Find(&q)
	return q
}

func (db *questionConnection) DeleteQuestion(q entity.Question) {
	db.connection.Delete(&q)
}

func (db *questionConnection) FindQuestionByID(questionID uint64) entity.Question {
	var question entity.Question
	db.connection.Preload("User").Preload("Answers").Preload("Answers.User").Find(&question, questionID)
	return question
}

func (db *questionConnection) AllQuestions(question *entity.Question, pagination *entity.Pagination) []entity.Question {
	var questions []entity.Question
	offset := (pagination.Page - 1) * pagination.Limit
	db.connection.Preload("User").Limit(pagination.Limit).Order(pagination.Sort).Offset(offset).Find(&questions)
	return questions
}

func (db *questionConnection) GetNumberOfLikesForQuestion() []entity.Question {
	var questions []entity.Question
	db.connection.Raw(" select distinct q.body as Body, q.id as id,(select COUNT(*) FROM question_likes where question_id = q.id) AS NumOfLikes from questions q left join question_likes ql on q.id = ql.question_id order by NumOfLikes desc limit 5").Scan(&questions)
	return questions
}

func (db *questionConnection) QuestionPage(questionID uint64, userID uint64) entity.QuestionResult {
	var question entity.QuestionResult
	db.connection.Raw(" select q.body as Body, q.id as QuestionID,(select COUNT(*) FROM question_likes where question_id = q.id) AS NumOfLikes,(select COUNT(*) > 0 FROM question_likes where user_id =  ? and question_id = q.id) AS IsLiked,(q.user_id = ? ) as IsEditable from questions q left join question_likes ql on q.id = ql.question_id where q.id = ?", userID, userID, questionID).Scan(&question)
	return question
}

func (db *questionConnection) GetUserQuestions(userId uint64, question *entity.Question, pagination *entity.Pagination) []entity.Question {
	var questions []entity.Question
	offset := (pagination.Page - 1) * pagination.Limit
	db.connection.Where("user_id = ?", userId).Limit(pagination.Limit).Order(pagination.Sort).Offset(offset).Find(&questions)
	return questions
}

func (db *questionConnection) Like(questiondId uint64, userID uint64) {
	var questionLikes entity.QuestionLikes
	questionLikes.UserID = userID
	questionLikes.QuestionID = questiondId
	db.connection.Save(&questionLikes)
}

func (db *questionConnection) DeleteLike(questionId uint64, userID uint64) {
	var questionLikes entity.QuestionLikes
	questionLikes.UserID = userID
	questionLikes.QuestionID = questionId
	db.connection.Where("question_id = ? and user_id = ?", questionId, userID).Delete(&questionLikes)
}
