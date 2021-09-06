package service

import (
	"errors"
	"fmt"
	"ibf-benevolence/entity"
	"ibf-benevolence/repository"
)

type UserService interface {
	GetAll() ([]entity.User, error)
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

func (us userService) GetAll() ([]entity.User, error) {
	fmt.Println("Getting all user from repository")
	users, err := us.userRepository.FindAllUser()
	if err != nil {
		return nil, errors.New("failed getting all user")
	}
	return users, nil
}
