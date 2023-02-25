package service

import (
	"bytes"
	"context"
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

type IUserService interface {
	GetAll(pagination *app.Pagination) *app.PaginatedResult[app.User]
	GetById(userId uint) *app.User
	Register(userInput validation.RegisterUser) *app.User
	Login(userInput validation.LoginUser) *string
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

func (us UserService) GetAll(pagination *app.Pagination) *app.PaginatedResult[app.User] {
	logger.Info("Getting all user from repository")
	var count int64

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := us.userRepository.GetAllUser(ctx, pagination, &count)
	if err != nil {
		panic(err)
	}

	var userDatas []app.User
	for _, x := range users {
		userDatas = append(userDatas, x.UserToUser())
	}
	paginatedResult := app.PaginatedResult[app.User]{
		Items:      userDatas,
		Page:       pagination.Page,
		Size:       len(userDatas),
		TotalItems: int(count),
		TotalPage:  int(math.Ceil(float64(count) / float64(pagination.Limit))),
	}

	return &paginatedResult
}

func (us UserService) GetById(userId uint) *app.User {
	logger.Info("Getting user from repository")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := us.userRepository.GetByUserId(ctx, userId)
	if err != nil {
		panic(err)
	}

	userData := user.UserToUser()
	return &userData
}

func (us UserService) Register(userInput validation.RegisterUser) *app.User {
	logger.Info(`Registering new user, user = %+v`, userInput)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existing, err := us.userRepository.IsExistingEmail(ctx, userInput.Email)
	if err != nil {
		panic(err)
	}
	if *existing {
		panic(errors.New("email is already exists"))
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
	userData := createdUser.UserToUser()
	return &userData
}

func (us UserService) Login(userInput validation.LoginUser) *string {
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
		panic(errors.New(`user Password is wrong`))
	}

	sign := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), app.UserClaims{
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
		panic(err)
	}
	return &token
}

var _ IUserService = (*UserService)(nil)
