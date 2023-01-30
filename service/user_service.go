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
	GetAll(pagination *app.Pagination) *app.PaginatedResult[app.User]
	GetById(userId uint) *app.User
	Register(userInput validation.RegisterUser) *app.User
	Login(userInput validation.LoginUser) *string
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	logger.Info("Initializing user service..")
	return &userService{
		userRepository: userRepository,
	}
}

func (us userService) GetAll(pagination *app.Pagination) *app.PaginatedResult[app.User] {
	logger.Info("Getting all user from repository")
	var count int64
	users := us.userRepository.GetAllUser(pagination, &count)
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

func (us userService) GetById(userId uint) *app.User {
	logger.Info("Getting user from repository")
	user := us.userRepository.GetByUserId(userId)
	userData := user.UserToUser()
	return &userData
}

func (us userService) Register(userInput validation.RegisterUser) *app.User {
	logger.Info(`Registering new user, user = %+v`, userInput)
	existing := us.userRepository.IsExistingEmail(userInput.Email)
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
	createdUser := us.userRepository.AddUser(user)
	userData := createdUser.UserToUser()
	return &userData
}

func (us userService) Login(userInput validation.LoginUser) *string {
	logger.Info(`Login user by email, email = %s`, userInput.Email)
	user := us.userRepository.GetByEmail(userInput.Email)

	bv := []byte(userInput.Password)
	hasher := sha256.New()
	hasher.Write(bv)

	if !bytes.Equal(user.Password, bv) {
		panic(errors.New(`user Password is wrong`))
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
		panic(err)
	}
	return &token
}

var _ UserService = (*userService)(nil)
