package controller

import (
	"encoding/json"
	"product-crud/cache"
	"product-crud/controller/response"
	"product-crud/service"
	"product-crud/util"
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

func NewUserController(userService service.UserService) *userController {
	logger.Info("Initializing user controller..")
	us := userService
	return &userController{
		userService: us,
	}
}

func (uc userController) GetAllUser(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic: %+v", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()

	logger.Info("Get all user request")

	pagination := util.GeneratePaginationFromRequest(c)
	hash := util.HashFromStruct(pagination)
	key := "GetAllUser:all:" + hash

	var cached interface{}
	if c.DefaultQuery("no_cache", "0") == "0" {
		cached = cache.Get(key)
	}
	var users interface{}
	isFromCache := false
	if cached != nil {
		logger.Info("Getting from cache")
		users = cached
		isFromCache = true
	} else {
		users = uc.userService.GetAll(&pagination)
	}
	if cached == nil {
		cache.Set(key, users)
	}
	logger.Info("Get all user success")
	response.Success(c, users, isFromCache)
}

func (uc userController) GetUserById(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic: %+v", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()

	logger.Info(`Get user by id, id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		panic(err)
	}

	key := "GetUserById:" + c.Param("id")

	var cached interface{}
	if c.DefaultQuery("no_cache", "0") == "0" {
		cached = cache.Get(key)
	}

	var user interface{}
	isFromCache := false
	if cached != nil {
		logger.Info("Getting from cache")
		user = cached
		isFromCache = true
	} else {
		user = uc.userService.GetById(uint(id))
	}
	if cached == nil {
		cache.Set(key, user)
	}

	logger.Info(`Get user by id, id = %s success`, c.Param("id"))
	response.Success(c, user, isFromCache)
}

func (uc userController) RegisterUser(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic: %+v", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	logger.Info(`Register new user request`)
	var input validation.RegisterUser
	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		panic(err)
	}
	logger.Info(`Validating request, request = %+v`, input)
	v := validator.New()
	err = v.Struct(input)
	if err != nil {
		panic(err)
	}
	user := uc.userService.Register(input)

	logger.Info(`Register new user success`)
	response.Success(c, user, false)
}

func (uc userController) LoginUser(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic: %+v", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	logger.Info(`Login User request`)
	var input validation.LoginUser
	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		panic(err)
	}
	v := validator.New()
	err = v.Struct(input)
	if err != nil {
		panic(err)
	}
	user := uc.userService.Login(input)

	logger.Info(`Login User success`)
	response.Success(c, user, false)
}

var _ UserController = (*userController)(nil)
