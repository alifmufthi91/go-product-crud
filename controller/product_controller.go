package controller

import (
	"encoding/json"
	"errors"
	"product-crud/app"
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

type IProductController interface {
	GetAllProduct(c *gin.Context)
	GetProductById(c *gin.Context)
	AddProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type ProductController struct {
	productService service.IProductService
}

func NewProductController(productService service.IProductService) ProductController {
	logger.Info("Initializing product controller..")
	return ProductController{
		productService: productService,
	}
}

func (pc ProductController) GetAllProduct(c *gin.Context) {
	defer response.ErrorHandling(c)

	logger.Info("Get all product request")
	pagination := util.GeneratePaginationFromRequest(c)

	hash := util.HashFromStruct(pagination)
	key := "GetAllProduct:all:" + hash

	var products = app.PaginatedResult[app.Product]{}
	if c.DefaultQuery("no_cache", "0") == "0" {
		err := cache.Get(key, &products)
		if err != nil {
			panic(err)
		}
	}

	isFromCache := false
	if !products.IsEmpty() {
		isFromCache = true
	} else {
		products = *pc.productService.GetAll(&pagination)
		cache.Set(key, products)
	}

	logger.Info("Get all product success")
	response.Success(c, products, isFromCache)
}

func (pc ProductController) GetProductById(c *gin.Context) {
	defer response.ErrorHandling(c)

	logger.Info(`Get product by id request, id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	key := "GetProductById:" + c.Param("id")

	var product = app.Product{}
	if c.DefaultQuery("no_cache", "0") == "0" {
		err := cache.Get(key, &product)
		if err != nil {
			panic(err)
		}
	}

	isFromCache := false
	if !product.IsEmpty() {
		isFromCache = true
	} else {
		product = *pc.productService.GetById(uint(id))
		cache.Set(key, product)
	}

	logger.Info(`Get product by id, id = %s success`, c.Param("id"))
	response.Success(c, product, isFromCache)
}

func (pc ProductController) AddProduct(c *gin.Context) {
	defer response.ErrorHandling(c)

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
		panic(errors.New("userClaims type is not correct"))
	}
	product := pc.productService.AddProduct(input, user.UserId)

	logger.Info(`Add new product success`)
	response.Success(c, product, false)
}

func (pc ProductController) UpdateProduct(c *gin.Context) {
	defer response.ErrorHandling(c)

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
		panic(errors.New("userClaims type is not correct"))
	}
	product := pc.productService.UpdateProduct(uint(id), input, user.UserId)

	logger.Info(`Update product of id = %s success`, c.Param("id"))
	response.Success(c, product, false)
}

func (pc ProductController) DeleteProduct(c *gin.Context) {
	defer response.ErrorHandling(c)

	logger.Info(`Delete product of id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	userClaims, _ := c.Get("user")
	user, ok := userClaims.(*app.MyCustomClaims)
	if !ok {
		panic(errors.New("userClaims type is not correct"))
	}
	pc.productService.DeleteProduct(uint(id), user.UserId)

	logger.Info(`Delete product of id = %s success`, c.Param("id"))
	response.Success(c, nil, false)
}

var _ IProductController = (*ProductController)(nil)
