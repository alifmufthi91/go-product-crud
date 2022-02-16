package controller

import (
	"encoding/json"
	"errors"
	"product-crud/app"
	"product-crud/controller/response"
	"product-crud/service"
	"product-crud/util/logger"
	"product-crud/validation"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductController interface {
	GetAllProduct(*gin.Context)
	GetProductById(*gin.Context)
	AddProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
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

func (pc productController) GetAllProduct(c *gin.Context) {
	logger.Info("Get all product request")
	products, err := pc.productService.GetAll()
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	logger.Info("Get all product success")
	response.Success(c, products)
}

func (pc productController) GetProductById(c *gin.Context) {
	logger.Info(`Get product by id request, id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	product, err := pc.productService.GetById(uint(id))
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	logger.Info(`Get product by id, id = %s success`, c.Param("id"))
	response.Success(c, product)
}

func (pc productController) AddProduct(c *gin.Context) {
	logger.Info(`Add new product request`)
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
	user, ok := userClaims.(*app.MyCustomClaims)
	if !ok {
		logger.Error("Error: userClaims type is not correct")
		response.Fail(c, "Error: userClaims type is not correct")
		return
	}
	product, err := pc.productService.AddProduct(input, user.UserId)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	logger.Info(`Add new product success`)
	response.Success(c, product)
}

func (pc productController) UpdateProduct(c *gin.Context) {
	logger.Info(`Update product of id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	var input validation.UpdateProduct
	err = json.NewDecoder(c.Request.Body).Decode(&input)
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
	user, ok := userClaims.(*app.MyCustomClaims)
	if !ok {
		logger.Error("Error: userClaims type is not correct")
		response.Fail(c, "Error: userClaims type is not correct")
		return
	}
	product, err := pc.productService.UpdateProduct(uint(id), input, user.UserId)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	logger.Info(`Update product of id = %s success`, c.Param("id"))
	response.Success(c, product)
}

func (pc productController) DeleteProduct(c *gin.Context) {
	logger.Info(`Delete product of id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	userClaims, _ := c.Get("user")
	user, ok := userClaims.(*app.MyCustomClaims)
	if !ok {
		logger.Error("Error: userClaims type is not correct")
		response.Fail(c, "Error: userClaims type is not correct")
		return
	}
	err = pc.productService.DeleteProduct(uint(id), user.UserId)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, err.Error())
		return
	}
	logger.Info(`Delete product of id = %s success`, c.Param("id"))
	response.Success(c, nil)
}
