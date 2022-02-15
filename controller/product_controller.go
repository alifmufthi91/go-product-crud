package controller

import (
	"encoding/json"
	"product-crud/app"
	"product-crud/controller/response"
	"product-crud/service"
	"product-crud/util/logger"
	"product-crud/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductController interface {
	// GetAllUser(*gin.Context)
	// GetUserById(*gin.Context)
	AddProduct(c *gin.Context)
	// LoginUser(c *gin.Context)
}

type productController struct {
	productService service.ProductService
}

func NewProductController() ProductController {
	logger.Info("Initializing product controller..")
	ps := service.NewProductService()
	return productController{
		productService: ps,
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

// func (uc userController) GetUserById(c *gin.Context) {
// 	logger.Info(`Get user by id, id = %s`, c.Param("id"))
// 	id := c.Param("id")
// 	user, err := uc.userService.GetById(id)
// 	if err != nil {
// 		logger.Error(err.Error())
// 		response.Fail(c, errors.New("something went wrong").Error())
// 		return
// 	}
// 	response.Success(c, user)
// }

func (pc productController) AddProduct(c *gin.Context) {
	logger.Info(`Add new user`)
	var input validation.AddProduct
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
	userClaims, _ := c.Get("user")
	user := userClaims.(*app.MyCustomClaims)
	product, err := pc.productService.AddProduct(input, user.UserId)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, product)
}
