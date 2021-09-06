package controller

import (
	"fmt"
	"ibf-benevolence/controller/response"
	"ibf-benevolence/service"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetAllUser(*gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController() UserController {
	fmt.Println("Initializing user controller..")
	us := service.NewUserService()
	return userController{
		userService: us,
	}
}

func (uc userController) GetAllUser(c *gin.Context) {
	fmt.Println("Get all user requested")
	users, err := uc.userService.GetAll()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, users)
	}
}
