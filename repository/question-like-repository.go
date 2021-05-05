package repository

// type QuestionLikeRepository interface {
// 	FindUserAndQuestion(questionID int64, userID int64) bool
// 	AlreadyLiked(questionID int64, userID int64) bool
// }

// type likeConnection struct {
// 	connection *gorm.DB
// }

// func NewQuestionLikeRepository(dbConn *gorm.DB) QuestionLikeRepository {
// 	return &likeConnection{
// 		connection: dbConn,
// 	}
// }

// func (db *questionConnection) AddLike(questionID int64, userID int64) {
// 	var question entity.Question
// 	db.connection.Find(&question, "id ? = AND user_id ? = ", questionID, userID)
// 	HandleLike(question)
// 	return question
// }

// func HandleLike(question entity.Question) {
// 	if question.IsLiked == false {
// 		question.IsLiked = true
// 		question.NumOfLikes++
// 	} else {
// 		question.IsLiked = false
// 		question.NumOfLikes--
// 	}
// }
