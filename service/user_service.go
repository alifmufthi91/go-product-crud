package service

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"product-crud/app"
	"product-crud/config"
	"product-crud/models"
	"product-crud/repository"
	"product-crud/util/logger"
	"product-crud/validation"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

type UserService interface {
	GetAll(pagination *app.Pagination) (*app.PaginatedResult, error)
	GetById(userId uint) (*app.User, error)
	Register(userInput validation.RegisterUser) (*app.User, error)
	Login(userInput validation.LoginUser) (*string, error)
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

func (us userService) GetAll(pagination *app.Pagination) (*app.PaginatedResult, error) {
	logger.Info("Getting all user from repository")
	var count int64
	users, err := us.userRepository.GetAllUser(pagination, &count)
	if err != nil {
		return nil, err
	}
	var userDatas []app.User
	for _, x := range users {
		userDatas = append(userDatas, x.UserToUser())
	}
	paginatedResult := app.PaginatedResult{
		Items:      userDatas,
		Page:       pagination.Page,
		Size:       len(userDatas),
		TotalItems: int(count),
		TotalPage:  int(math.Ceil(float64(count) / float64(pagination.Limit))),
	}

	return &paginatedResult, nil
}

func (us userService) GetById(userId uint) (*app.User, error) {
	logger.Info("Getting user from repository")
	user, err := us.userRepository.GetByUserId(userId)
	if err != nil {
		return nil, err
	}
	userData := user.UserToUser()
	return &userData, nil
}

func (us userService) Register(userInput validation.RegisterUser) (*app.User, error) {
	logger.Info(`Registering new user, user = %+v`, userInput)
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

func (us userService) Login(userInput validation.LoginUser) (*string, error) {
	logger.Info(`Login user by email, email = %s`, userInput.Email)
	user, err := us.userRepository.GetByEmail(userInput.Email)
	if err != nil {
		return nil, err
	}

	bv := []byte(userInput.Password)
	hasher := sha256.New()
	hasher.Write(bv)

	if !bytes.Equal(user.Password, bv) {
		return nil, errors.New(`user Password is wrong`)
	}

	sign := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), app.MyCustomClaims{
		UserId:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  fmt.Sprintf(`%s %s`, user.FirstName, user.LastName),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(),
			Issuer:    "test",
		},
	})
	token, err := sign.SignedString([]byte(config.Env.JWTSECRET))
	if err != nil {
		return nil, err
	}
	return &token, nil
}
