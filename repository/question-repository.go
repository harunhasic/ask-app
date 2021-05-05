package repository

import (
	"github.com/mop/entity"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	InsertQuestion(q entity.Question) entity.Question
	UpdateQuestion(q entity.Question) entity.Question
	DeleteQuestion(q entity.Question)
	AllQuestions() []entity.Question
	FindQuestionByID(questionID uint64) entity.Question
	GetNumberOfLikesForQuestion() []entity.Question
	Like(questionID uint64, userID uint64)
	DeleteLike(questionID uint64, userID uint64)
	QuestionPage(questionID uint64, userID uint64) entity.QuestionResult
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
	db.connection.Preload("User").Preload("Answers").Find(&question, questionID)
	return question
}

func (db *questionConnection) AllQuestions() []entity.Question {
	var questions []entity.Question
	db.connection.Preload("User").Limit(20).Order("created_at desc").Find(&questions)
	return questions
}

func (db *questionConnection) GetNumberOfLikesForQuestion() []entity.Question {
	var questions []entity.Question
	db.connection.Limit(5).Order("num_of_likes desc").Find(&questions)
	return questions
}

func (db *questionConnection) QuestionPage(questionID uint64, userID uint64) entity.QuestionResult {
	var question entity.QuestionResult
	db.connection.Raw(" select q.body as Body, q.id as QuestionID, q.num_of_likes as NumOfLikes, COUNT(ql.question_id) as NumOfLikes,(select COUNT(*) > 0 FROM question_likes where user_id =  ? and question_id = q.id) AS IsLiked,(q.user_id = ? ) as IsEditable from questions q left join question_likes ql on q.id = ql.question_id where q.id = ?", userID, userID, questionID).Scan(&question)
	return question
}

func (db *questionConnection) Like(questiondId uint64, userID uint64) {
	var questionLikes entity.QuestionLikes
	questionLikes.UserID = userID
	questionLikes.QuestionID = questiondId
	db.connection.Save(&questionLikes)
}

func (db *questionConnection) DeleteLike(questiondId uint64, userID uint64) {
	var questionLikes entity.QuestionLikes
	questionLikes.UserID = userID
	questionLikes.QuestionID = questiondId

	db.connection.Exec("DELETE FROM QuestionLikes WHERE question_id = ? AND user_id = ?", questiondId, userID)
	db.connection.Delete(&questionLikes)
}

// func (db *questionConnection) GetTop5LikedQuestions(){
// 	var questions []entity.Question
// 	db.connection.Preload("QuestionLikes").Find(&questions)
// 	foreach question in questions
// 		question.numberoflikes = question.QuestionLikes.count

// 	return question.Limit(5).Order("num_of_likes desc")
// }

// &questions.QuestionLikes.Where()
// var questionsvslikes | QuestionLikes

// foreach querstion in questins
// 	question.numberoflikes =  questionsvslikes.where(questionid = question.id).countdb //GetNumberOfLikesForQuestion
// 	question.isliked = questionsvslikes.any(userid == logovaniuserid && questionid == question.id)

// func Dislike(questiondId, userid){
// 	db.connection/QuestionLikes .remove(questiondi, userid )
// db.connection.question.numoflikes++

// }
