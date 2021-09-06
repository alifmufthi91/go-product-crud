package controller

import (
	"ibf-benevolence/controller/response"
	"ibf-benevolence/service"
	"ibf-benevolence/util/logger"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetAllUser(*gin.Context)
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
		response.Fail(c, err.Error())
	} else {
		response.Success(c, users)
	}
}
