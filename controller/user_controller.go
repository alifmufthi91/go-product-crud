package controller

import (
	"encoding/json"
	"errors"
	"ibf-benevolence/controller/response"
	"ibf-benevolence/model"
	"ibf-benevolence/service"
	"ibf-benevolence/util/logger"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetAllUser(*gin.Context)
	GetUserById(*gin.Context)
	AddUser(c *gin.Context)
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
	logger.Info("Get all user requested")
	users, err := uc.userService.GetAll()
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	response.Success(c, users)
}

func (uc userController) GetUserById(c *gin.Context) {
	logger.Info("Get user by id requested")
	id := c.Param("id")
	user, err := uc.userService.GetById(id)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	response.Success(c, user)
}

func (uc userController) AddUser(c *gin.Context) {
	logger.Info("Add user requested")
	var input model.UserRegisterInput
	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	user, err := uc.userService.Register(input)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	response.Success(c, user)
}
