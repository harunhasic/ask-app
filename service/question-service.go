package service

import (
	"fmt"
	"log"

	"github.com/mashingan/smapping"
	"github.com/mop/dto"
	"github.com/mop/entity"
	"github.com/mop/repository"
)

//Interface that represents the methods a question service needs to have
type QuestionService interface {
	Insert(q dto.QuestionCreateDTO) entity.Question
	Update(q dto.QuestionUpdateDTO) entity.Question
	Delete(q entity.Question)
	All(question *entity.Question, pagination *entity.Pagination) []entity.Question
	FindById(questionID uint64) entity.Question
	IsAllowedToEdit(userID string, questionID uint64) bool
	GetNumberOfLikesForQuestion() []entity.Question
	QuestionPage(questionID uint64, userID uint64) entity.QuestionResult
	Like(questionID uint64, userID uint64)
	DeleteLike(questionID uint64, userID uint64)
	GetUserQuestions(userId uint64, question *entity.Question, pagination *entity.Pagination) []entity.Question
}

type questionService struct {
	questionRepository repository.QuestionRepository
}

//NewQuestionService creates a new instance of the QuestionService
func NewQuestionService(questionRepo repository.QuestionRepository) QuestionService {
	return &questionService{
		questionRepository: questionRepo,
	}
}

//Implementation of the service Insert method which saves the question
func (service *questionService) Insert(q dto.QuestionCreateDTO) entity.Question {
	question := entity.Question{}
	err := smapping.FillStruct(&question, smapping.MapFields(&q))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.questionRepository.InsertQuestion(question)
	return res
}

func (service *questionService) Update(q dto.QuestionUpdateDTO) entity.Question {
	question := entity.Question{}
	err := smapping.FillStruct(&question, smapping.MapFields(&q))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.questionRepository.UpdateQuestion(question)
	return res
}

func (service *questionService) Delete(q entity.Question) {
	service.questionRepository.DeleteQuestion(q)
}

func (service *questionService) All(question *entity.Question, pagination *entity.Pagination) []entity.Question {
	return service.questionRepository.AllQuestions(question, pagination)
}

func (service *questionService) FindById(questionID uint64) entity.Question {
	return service.questionRepository.FindQuestionByID(questionID)
}

func (service *questionService) IsAllowedToEdit(userID string, questionID uint64) bool {
	q := service.questionRepository.FindQuestionByID(questionID)
	id := fmt.Sprintf("%v", q.UserID)
	return userID == id
}

func (service *questionService) QuestionPage(questionID uint64, userID uint64) entity.QuestionResult {
	return service.questionRepository.QuestionPage(questionID, userID)
}

func (service *questionService) Like(questionID uint64, userID uint64) {
	service.questionRepository.Like(questionID, userID)
}

func (service *questionService) DeleteLike(questionID uint64, userID uint64) {
	service.questionRepository.DeleteLike(questionID, userID)
}

func (service *questionService) GetNumberOfLikesForQuestion() []entity.Question {
	return service.questionRepository.GetNumberOfLikesForQuestion()
}

func (service *questionService) GetUserQuestions(userId uint64, question *entity.Question, pagination *entity.Pagination) []entity.Question {
	return service.questionRepository.GetUserQuestions(userId, question, pagination)
}
