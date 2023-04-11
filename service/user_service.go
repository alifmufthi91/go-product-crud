package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"product-crud/config"
	"product-crud/dto/app"
	"product-crud/dto/request"
	"product-crud/dto/response"
	"product-crud/models"
	"product-crud/repository"
	errorUtil "product-crud/util/error"
	"product-crud/util/logger"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

type IUserService interface {
	GetAll(pagination app.Pagination) (app.PaginatedResult[response.GetUserResponse], error)
	GetById(userId uint) (response.GetUserResponse, error)
	Register(userInput request.UserRegisterRequest) (response.GetUserResponse, error)
	Login(userInput request.UserLoginRequest) (string, error)
}

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) UserService {
	logger.Info("Initializing user service..")
	return UserService{
		userRepository: userRepository,
	}
}

func (us UserService) GetAll(pagination app.Pagination) (app.PaginatedResult[response.GetUserResponse], error) {
	logger.Info("Getting all user from repository")
	var count int64
	var result app.PaginatedResult[response.GetUserResponse]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := us.userRepository.GetAllUser(ctx, pagination, &count)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}

	var userDatas []response.GetUserResponse
	for _, x := range users {
		userDatas = append(userDatas, response.NewGetUserResponse(x))
	}
	result = app.PaginatedResult[response.GetUserResponse]{
		Items:      userDatas,
		Page:       pagination.Page,
		Size:       len(userDatas),
		TotalItems: int(count),
		TotalPage:  int(math.Ceil(float64(count) / float64(pagination.Limit))),
	}
	return result, nil
}

func (us UserService) GetById(userId uint) (response.GetUserResponse, error) {
	logger.Info("Getting user from repository")

	var result response.GetUserResponse
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nullableUser, err := us.userRepository.GetByUserId(ctx, userId)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}
	if !nullableUser.Valid {
		logger.Error("Invalid user")
		return result, errors.New("Invalid user")
	}

	result = response.NewGetUserResponse(nullableUser.Stuff)
	return result, nil
}

func (us UserService) Register(userInput request.UserRegisterRequest) (response.GetUserResponse, error) {
	logger.Info(`Registering new user, user = %+v`, userInput)

	var result response.GetUserResponse
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existing, err := us.userRepository.IsExistingEmail(ctx, userInput.Email)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}
	if existing {
		err := errorUtil.ParamIllegal("email is already exists")
		logger.Error("Error : %v", err)
		return result, err
	}

	bv := []byte(userInput.Password)
	hasher := sha256.New()
	_, err = hasher.Write(bv)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}

	user := models.User{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Email:     userInput.Email,
		Password:  hasher.Sum(nil),
	}

	nullableUser, err := us.userRepository.AddUser(ctx, user)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}
	result = response.NewGetUserResponse(nullableUser.Stuff)
	return result, nil
}

func (us UserService) Login(userInput request.UserLoginRequest) (string, error) {
	logger.Info(`Login user by email, email = %s`, userInput.Email)

	var token string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nullableUser, err := us.userRepository.GetByEmail(ctx, userInput.Email)
	if err != nil {
		logger.Error("Error : %v", err)
		return token, err
	}
	if !nullableUser.Valid {
		logger.Error("Invalid user")
		return token, errors.New("Invalid user")
	}

	bv := []byte(userInput.Password)
	hasher := sha256.New()
	_, err = hasher.Write(bv)
	if err != nil {
		logger.Error("Error : %v", err)
		return token, err
	}
	user := nullableUser.Stuff
	if !bytes.Equal(user.Password, hasher.Sum(nil)) {
		err := errorUtil.ParamIllegal("user password is incorrect")
		logger.Error("Error : %v", err)
		return token, err
	}

	sign := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), app.UserClaims{
		UserId:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  fmt.Sprintf(`%s %s`, user.FirstName, user.LastName),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 999999).Unix(),
			Issuer:    "test",
		},
	})
	token, err = sign.SignedString([]byte(config.GetEnv().JWTSECRET))
	if err != nil {
		logger.Error("Error : %v", err)
		return token, err
	}
	return token, nil
}

var _ IUserService = (*UserService)(nil)
