package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/mop/dto"
	"github.com/mop/entity"
	"github.com/mop/repository"
)

//UserService is the contract for the UserService and lists all the methods
type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userID string) entity.User
	GetAnswers() []entity.MostAnswersDTO
}

type userService struct {
	userRepository repository.UserRepository
}

//NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed to map %v:", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) entity.User {
	return service.userRepository.ProfileUser(userID)
}

func (service *userService) GetAnswers() []entity.MostAnswersDTO {
	return service.userRepository.GetAnswers()
}
