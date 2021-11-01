package service

import (
	"ibf-benevolence/entity"
	"ibf-benevolence/model"
	"ibf-benevolence/repository"
	"ibf-benevolence/util/logger"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"
)

type UserService interface {
	GetAll() ([]entity.User, error)
	GetById(userId string) (*entity.User, error)
	Register(userInput model.UserRegisterInput) (*entity.User, error)
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
		return nil, err
	}
	return users, nil
}

func (us userService) GetById(userId string) (*entity.User, error) {
	logger.Info("Getting user from repository")
	user, err := us.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us userService) Register(userInput model.UserRegisterInput) (*entity.User, error) {
	logger.Info("register new user to repository")
	id, err := uuid.NewV1()
	if err != nil {
		return nil, err
	}
	user := entity.User{UserId: id.String(), Status: "INACTIVE", CreatedAt: time.Now().UnixMilli()}
	copier.Copy(&user, &userInput)
	err = us.userRepository.AddUser(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
