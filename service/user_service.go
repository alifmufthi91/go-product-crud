package service

import (
	"fmt"
	"ibf-benevolence/entity"
	"ibf-benevolence/repository"
)

type UserService interface {
	GetAll() []entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService() UserService {
	fmt.Println("Initializing user service")
	ur := repository.NewUserRepository()
	return userService{
		userRepository: ur,
	}
}

func (us userService) GetAll() []entity.User {
	fmt.Println("Getting all user from repository")
	return us.userRepository.FindAllUser()
}
