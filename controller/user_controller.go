package controller

import (
	"encoding/json"
	"errors"
	"product-crud/controller/response"
	"product-crud/service"
	"product-crud/util/logger"
	"product-crud/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController interface {
	// GetAllUser(*gin.Context)
	GetUserById(*gin.Context)
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController() UserController {
	logger.Info("Initializing user controller..")
	us := service.NewUserService()
	return userController{
		userService: us,
	}
}

// func (uc userController) GetAllUser(c *gin.Context) {
// 	logger.Info("Get all user requested")
// 	users, err := uc.userService.GetAll()
// 	if err != nil {
// 		logger.Error(err.Error())
// 		response.Fail(c, errors.New("something went wrong").Error())
// 		return
// 	}
// 	response.Success(c, users)
// }

func (uc userController) GetUserById(c *gin.Context) {
	logger.Info(`Get user by id, id = %s`, c.Param("id"))
	id := c.Param("id")
	user, err := uc.userService.GetById(id)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	response.Success(c, user)
}

func (uc userController) RegisterUser(c *gin.Context) {
	logger.Info(`Register new user`)
	var input validation.RegisterUser
	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	logger.Info(`Validating request, request = %+v`, input)
	v := validator.New()
	err = v.Struct(input)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	user, err := uc.userService.Register(input)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (uc userController) LoginUser(c *gin.Context) {
	logger.Info(`Login User`)
	var input validation.LoginUser
	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	v := validator.New()
	err = v.Struct(input)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	user, err := uc.userService.Login(input)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}
