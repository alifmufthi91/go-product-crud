package service

import (
	"errors"
	"ibf-benevolence/entity"
	"ibf-benevolence/repository"
	"ibf-benevolence/util/logger"
)

type UserService interface {
	GetAll() ([]entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService() UserService {
	logger.Info("Initializing user service..")
	ur := repository.NewUserRepository()
	return userService{
		userRepository: ur,
	}
}

func (us userService) GetAll() ([]entity.User, error) {
	logger.Info("Getting all user from repository")
	users, err := us.userRepository.FindAllUser()
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed getting all user")
	}
	return users, nil
}
