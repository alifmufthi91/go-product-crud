package controller

import (
	"context"
	"product-crud/app"
	"product-crud/cache"
	ERROR_CONSTANT "product-crud/constant"
	"product-crud/dto/request"
	resp "product-crud/dto/response"
	"product-crud/service"
	"product-crud/util"
	"product-crud/util/logger"
	responseUtil "product-crud/util/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	defer responseUtil.ErrorHandling(c)

	logger.Info("Get all product request")
	pagination := util.GeneratePaginationFromRequest(c)

	hash := util.HashFromStruct(pagination)
	key := "GetAllProduct:all:" + hash

	var products app.PaginatedResult[resp.GetProductResponse]
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	if c.DefaultQuery("no_cache", "0") == "0" {
		err := cache.Get(ctx, key, &products)
		if err != nil {
			logger.Error("Error : %v", err)
			panic(ERROR_CONSTANT.INTERNAL_ERROR)
		}
	}

	isFromCache := false
	if !products.IsEmpty() {
		isFromCache = true
	} else {
		products = pc.productService.GetAll(pagination)
		cache.Set(ctx, key, products)
	}

	logger.Info("Get all product success")
	responseUtil.Success(c, products, isFromCache)
}

func (pc ProductController) GetProductById(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)

	logger.Info(`Get product by id request, id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("Error : %v", err)
		panic(ERROR_CONSTANT.INTERNAL_ERROR)
	}

	key := "GetProductById:" + c.Param("id")

	var product resp.GetProductResponse
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	if c.DefaultQuery("no_cache", "0") == "0" {
		err := cache.Get(ctx, key, &product)
		if err != nil {
			logger.Error("Error : %v", err)
			panic(ERROR_CONSTANT.INTERNAL_ERROR)
		}
	}

	isFromCache := false
	if !product.IsEmpty() {
		isFromCache = true
	} else {
		product = pc.productService.GetById(uint(id))
		cache.Set(ctx, key, product)
	}

	logger.Info(`Get product by id, id = %s success`, c.Param("id"))
	responseUtil.Success(c, product, isFromCache)
}

func (pc ProductController) AddProduct(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)

	logger.Info(`Add new product request`)
	var request request.ProductAddRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		panic(err)
	}
	user, err := util.GetUserClaims(c)
	if err != nil {
		panic(err)
	}
	product := pc.productService.AddProduct(request, user.UserId)

	logger.Info(`Add new product success`)
	responseUtil.Success(c, product, false)
}

func (pc ProductController) UpdateProduct(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)

	logger.Info(`Update product of id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("Error : %v", err)
		panic(ERROR_CONSTANT.INTERNAL_ERROR)
	}
	var request request.ProductUpdateRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		panic(err)
	}
	user, err := util.GetUserClaims(c)
	if err != nil {
		panic(err)
	}
	product := pc.productService.UpdateProduct(uint(id), request, user.UserId)

	logger.Info(`Update product of id = %s success`, c.Param("id"))
	responseUtil.Success(c, product, false)
}

func (pc ProductController) DeleteProduct(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)

	logger.Info(`Delete product of id = %s`, c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("Error : %v", err)
		panic(ERROR_CONSTANT.INTERNAL_ERROR)
	}
	user, err := util.GetUserClaims(c)
	if err != nil {
		panic(err)
	}
	pc.productService.DeleteProduct(uint(id), user.UserId)

	logger.Info(`Delete product of id = %s success`, c.Param("id"))
	responseUtil.Success(c, nil, false)
}

var _ IProductController = (*ProductController)(nil)
