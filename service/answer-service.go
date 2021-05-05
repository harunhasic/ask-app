package service

import (
	"fmt"
	"log"

	"github.com/mashingan/smapping"
	"github.com/mop/dto"
	"github.com/mop/entity"
	"github.com/mop/repository"
)

//Interface that represents the methods a answer service needs to have
type AnswerService interface {
	Insert(a dto.AnswerCreateDTO) entity.Answer
	Update(a dto.AnswerUpdateDTO) entity.Answer
	Delete(a entity.Answer)
	All() []entity.Answer
	FindById(bookID uint64) entity.Answer
	IsAllowedToEdit(userID string, answerID uint64) bool
	UsersWithMostAnswers() []entity.Answer
}

type answerService struct {
	answerRepository repository.AnswerRepository
}

//Creates a new instance of the answer service
func NewAnswerService(answerRepo repository.AnswerRepository) AnswerService {
	return &answerService{
		answerRepository: answerRepo,
	}
}

//Implementation of the service Insert method which saves the answer
func (service *answerService) Insert(a dto.AnswerCreateDTO) entity.Answer {
	answer := entity.Answer{}
	err := smapping.FillStruct(&answer, smapping.MapFields(&a))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.answerRepository.InsertAnswer(answer)
	return res
}

func (service *answerService) Update(a dto.AnswerUpdateDTO) entity.Answer {
	answer := entity.Answer{}
	err := smapping.FillStruct(&answer, smapping.MapFields(&a))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.answerRepository.UpdateAnswer(answer)
	return res
}

func (service *answerService) Delete(a entity.Answer) {
	service.answerRepository.DeleteAnswer(a)
}

func (service *answerService) All() []entity.Answer {
	return service.answerRepository.AllAnswers()
}

func (service *answerService) UsersWithMostAnswers() []entity.Answer {
	return service.answerRepository.UsersWithMostAnswers()
}

func (service *answerService) FindById(answerID uint64) entity.Answer {
	return service.answerRepository.FindAnswerByID(answerID)
}

func (service *answerService) IsAllowedToEdit(userID string, answerID uint64) bool {
	a := service.answerRepository.FindAnswerByID(answerID)
	id := fmt.Sprintf("%v", a.UserID)
	return userID == id
}
