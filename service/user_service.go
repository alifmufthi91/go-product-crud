package service

import (
	"crypto/sha256"
	"errors"
	"product-crud/app"
	"product-crud/models"
	"product-crud/repository"
	"product-crud/util/logger"
	"product-crud/validation"
)

type UserService interface {
	// GetAll() ([]app.User, error)
	// GetById(userId string) (*app.User, error)
	Register(userInput validation.RegisterUser) (*app.User, error)
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

// func (us userService) GetAll() ([]app.User, error) {
// 	logger.Info("Getting all user from repository")
// 	users, err := us.userRepository.FindAllUser()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }

// func (us userService) GetById(userId string) (*app.User, error) {
// 	logger.Info("Getting user from repository")
// 	user, err := us.userRepository.FindByUserId(userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

func (us userService) Register(userInput validation.RegisterUser) (*app.User, error) {
	logger.Info("registering new user")
	existing, _ := us.userRepository.IsExistingEmail(userInput.Email)
	if *existing {
		return nil, errors.New("email is already exists")
	}
	bv := []byte(userInput.Password)
	hasher := sha256.New()
	hasher.Write(bv)

	user := models.User{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Email:     userInput.Email,
		Password:  bv,
	}
	createdUser, err := us.userRepository.AddUser(user)
	if err != nil {
		return nil, err
	}
	userData := createdUser.UserToUser()
	return &userData, nil
}
