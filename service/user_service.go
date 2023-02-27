package service

import (
	"bytes"
	"context"
	"crypto/sha256"
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
	GetAll(pagination app.Pagination) app.PaginatedResult[response.GetUserResponse]
	GetById(userId uint) response.GetUserResponse
	Register(userInput request.UserRegisterRequest) response.GetUserResponse
	Login(userInput request.UserLoginRequest) string
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

func (us UserService) GetAll(pagination app.Pagination) app.PaginatedResult[response.GetUserResponse] {
	logger.Info("Getting all user from repository")
	var count int64

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := us.userRepository.GetAllUser(ctx, pagination, &count)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	var userDatas []response.GetUserResponse
	for _, x := range users {
		userDatas = append(userDatas, *response.NewGetUserResponse(x))
	}
	return app.PaginatedResult[response.GetUserResponse]{
		Items:      userDatas,
		Page:       pagination.Page,
		Size:       len(userDatas),
		TotalItems: int(count),
		TotalPage:  int(math.Ceil(float64(count) / float64(pagination.Limit))),
	}
}

func (us UserService) GetById(userId uint) response.GetUserResponse {
	logger.Info("Getting user from repository")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := us.userRepository.GetByUserId(ctx, userId)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	return *response.NewGetUserResponse(user)
}

func (us UserService) Register(userInput request.UserRegisterRequest) response.GetUserResponse {
	logger.Info(`Registering new user, user = %+v`, userInput)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existing, err := us.userRepository.IsExistingEmail(ctx, userInput.Email)
	if err != nil {
		panic(err)
	}
	if *existing {
		panic(errorUtil.ParamIllegal("email is already exists"))
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

	createdUser, err := us.userRepository.AddUser(ctx, user)
	if err != nil {
		panic(err)
	}
	return *response.NewGetUserResponse(createdUser)
}

func (us UserService) Login(userInput request.UserLoginRequest) string {
	logger.Info(`Login user by email, email = %s`, userInput.Email)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := us.userRepository.GetByEmail(ctx, userInput.Email)
	if err != nil {
		panic(err)
	}
	bv := []byte(userInput.Password)
	hasher := sha256.New()
	hasher.Write(bv)

	if !bytes.Equal(user.Password, bv) {
		panic(errorUtil.ParamIllegal("user password is incorrect"))
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
	token, err := sign.SignedString([]byte(config.GetEnv().JWTSECRET))
	if err != nil {
		panic(err)
	}
	return token
}

var _ IUserService = (*UserService)(nil)
