package controller

import (
	"encoding/json"
	"errors"
	"product-crud/controller/response"
	"product-crud/service"
	"product-crud/util/logger"
	"product-crud/validation"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController interface {
	GetAllUser(c *gin.Context)
	GetUserById(c *gin.Context)
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

func (uc userController) GetAllUser(c *gin.Context) {
	logger.Info("Get all user request")
	users, err := uc.userService.GetAll()
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	logger.Info("Get all user success")
	response.Success(c, users)
}

func (uc userController) GetUserById(c *gin.Context) {
	logger.Info(`Get user by id, id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	user, err := uc.userService.GetById(uint(id))
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	logger.Info(`Get user by id, id = %s success`, c.Param("id"))
	response.Success(c, user)
}

func (uc userController) RegisterUser(c *gin.Context) {
	logger.Info(`Register new user request`)
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
	logger.Info(`Register new user success`)
	response.Success(c, user)
}

func (uc userController) LoginUser(c *gin.Context) {
	logger.Info(`Login User request`)
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
	logger.Info(`Login User success`)
	response.Success(c, user)
}
