package service

import (
	"ibf-benevolence/entity"
	"ibf-benevolence/model"
	"ibf-benevolence/repository"
	"ibf-benevolence/util/logger"
	"time"

	"github.com/gofrs/uuid"
)

type UserService interface {
	GetAll() ([]entity.User, error)
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

func (us userService) Register(userInput model.UserRegisterInput) (*entity.User, error) {
	logger.Info("register new user to repository")
	id, err := uuid.NewV1()
	if err != nil {
		return nil, err
	}
	user := entity.User{}
	user.UserId = id.String()
	user.Name = userInput.Name
	user.AlgoAddress = userInput.AlgoAddress
	user.Email = userInput.Email
	user.PhoneNumberCode = userInput.PhoneNumberCode
	user.PhoneNumber = userInput.PhoneNumber
	user.PhotoUrl = &userInput.PhotoUrl
	user.Gender = userInput.Gender
	user.Status = "INACTIVE"
	user.CreatedAt = time.Now().UnixMilli()
	err = us.userRepository.AddUser(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
