package controller

import (
	"encoding/json"
	"log"
	"product-crud/app"
	"product-crud/controller/response"
	"product-crud/service"
	"product-crud/util"
	"product-crud/util/logger"
	"product-crud/validation"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductController interface {
	GetAllProduct(c *gin.Context)
	GetProductById(c *gin.Context)
	AddProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type productController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) *productController {
	logger.Info("Initializing product controller..")
	ps := productService
	return &productController{
		productService: ps,
	}
}

func (pc productController) GetAllProduct(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	logger.Info("Get all product request")
	pagination := util.GeneratePaginationFromRequest(c)
	products := pc.productService.GetAll(&pagination)

	logger.Info("Get all product success")
	response.Success(c, products)
}

func (pc productController) GetProductById(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	logger.Info(`Get product by id request, id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	product := pc.productService.GetById(uint(id))

	logger.Info(`Get product by id, id = %s success`, c.Param("id"))
	response.Success(c, product)
}

func (pc productController) AddProduct(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	logger.Info(`Add new product request`)
	var input validation.AddProduct
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
	userClaims, _ := c.Get("user")
	user, ok := userClaims.(*app.MyCustomClaims)
	if !ok {
		panic("Error: userClaims type is not correct")
	}
	product := pc.productService.AddProduct(input, user.UserId)

	logger.Info(`Add new product success`)
	response.Success(c, product)
}

func (pc productController) UpdateProduct(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	logger.Info(`Update product of id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	var input validation.UpdateProduct
	err = json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		panic(err)
	}
	logger.Info(`Validating request, request = %+v`, input)
	v := validator.New()
	err = v.Struct(input)
	if err != nil {
		panic(err)
	}
	userClaims, _ := c.Get("user")
	user, ok := userClaims.(*app.MyCustomClaims)
	if !ok {
		panic("Error: userClaims type is not correct")
	}
	product := pc.productService.UpdateProduct(uint(id), input, user.UserId)

	logger.Info(`Update product of id = %s success`, c.Param("id"))
	response.Success(c, product)
}

func (pc productController) DeleteProduct(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	logger.Info(`Delete product of id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	userClaims, _ := c.Get("user")
	user, ok := userClaims.(*app.MyCustomClaims)
	if !ok {
		panic("Error: userClaims type is not correct")
	}
	pc.productService.DeleteProduct(uint(id), user.UserId)

	logger.Info(`Delete product of id = %s success`, c.Param("id"))
	response.Success(c, nil)
}
