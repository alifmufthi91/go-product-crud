package controller

import (
	"context"
	"product-crud/cache"
	ERROR_CONSTANT "product-crud/constant/error_constant"
	"product-crud/dto/app"
	"product-crud/dto/request"
	resp "product-crud/dto/response"
	"product-crud/service"
	"product-crud/util"
	"product-crud/util/logger"
	responseUtil "product-crud/util/response"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type IUserController interface {
	GetAllUser(c *gin.Context)
	GetUserById(c *gin.Context)
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	GetAllUserRequestCounter(c *gin.Context)
}

type UserController struct {
	userService service.IUserService
}

var getAllUserRequestCalled uint64

func NewUserController(userService service.IUserService) UserController {
	logger.Info("Initializing user controller..")
	return UserController{
		userService: userService,
	}
}

func (uc UserController) GetAllUser(c *gin.Context) {
	atomic.AddUint64(&getAllUserRequestCalled, 1)
	defer responseUtil.ErrorHandling(c)

	logger.Info("Get all user request")

	pagination, err := util.GeneratePaginationFromRequest(c)
	if err != nil {
		panic(err)
	}
	hash, err := util.HashFromStruct(pagination)
	if err != nil {
		panic(err)
	}
	key := "GetAllUser:all:" + *hash

	var users app.PaginatedResult[resp.GetUserResponse]
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	if c.DefaultQuery("no_cache", "0") == "0" {
		err := cache.Get(ctx, key, &users)
		if err != nil {
			logger.Error("Error : %v", err)
			// panic(ERROR_CONSTANT.INTERNAL_ERROR)
		}
	}
	isFromCache := false
	if !users.IsEmpty() {
		isFromCache = true
	} else {
		val, err := uc.userService.GetAll(*pagination)
		if err != nil {
			panic(err)
		}
		users = *val
		go func() {
			ctx, cancel := context.WithTimeout(c, 3*time.Second)
			defer cancel()
			err := cache.Set(ctx, key, users)
			if err != nil {
				logger.Error("Error : %v", err)
			}
		}()
	}
	logger.Info("Get all user success")
	responseUtil.Ok(c, users, isFromCache)
}

func (uc UserController) GetUserById(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)

	logger.Info(`Get user by id, id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logger.Error("Error : %v", err)
		panic(ERROR_CONSTANT.INTERNAL_ERROR)
	}

	key := "GetUserById:" + c.Param("id")

	var user resp.GetUserResponse
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	if c.DefaultQuery("no_cache", "0") == "0" {
		err := cache.Get(ctx, key, &user)
		if err != nil {
			logger.Error("Error : %v", err)
			// panic(ERROR_CONSTANT.INTERNAL_ERROR)
		}
	}

	isFromCache := false
	if !user.IsEmpty() {
		isFromCache = true
	} else {
		val, err := uc.userService.GetById(uint(id))
		if err != nil {
			panic(err)
		}
		user = *val
		go func() {
			ctx, cancel := context.WithTimeout(c, 3*time.Second)
			defer cancel()
			err := cache.Set(ctx, key, user)
			if err != nil {
				logger.Error("Error : %v", err)
			}
		}()
	}

	logger.Info(`Get user by id, id = %s success`, c.Param("id"))
	responseUtil.Ok(c, user, isFromCache)
}

func (uc UserController) RegisterUser(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)

	logger.Info(`Register new user request`)
	var request request.UserRegisterRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		panic(err)
	}
	user, err := uc.userService.Register(request)
	if err != nil {
		panic(err)
	}
	logger.Info(`Register new user success`)
	responseUtil.Ok(c, *user, false)
}

func (uc UserController) LoginUser(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)

	logger.Info(`Login User request`)
	var request request.UserLoginRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		panic(err)
	}
	token, err := uc.userService.Login(request)
	if err != nil {
		panic(err)
	}
	logger.Info(`Login User success`)
	responseUtil.Ok(c, *token, false)
}

func (uc UserController) GetAllUserRequestCounter(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)

	responseUtil.Ok(c, atomic.LoadUint64(&getAllUserRequestCalled), false)
}
