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
	pagination, err := util.GeneratePaginationFromRequest(c)
	if err != nil {
		panic(err)
	}

	hash, err := util.HashFromStruct(pagination)
	if err != nil {
		panic(err)
	}
	key := "GetAllProduct:all:" + *hash

	var products app.PaginatedResult[resp.GetProductResponse]
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	if c.DefaultQuery("no_cache", "0") == "0" {
		err := cache.Get(ctx, key, &products)
		if err != nil {
			logger.Error("Error : %v", err)
		}
	}

	isFromCache := false
	if !products.IsEmpty() {
		isFromCache = true
	} else {
		val, err := pc.productService.GetAll(*pagination)
		if err != nil {
			panic(err)
		}
		products = *val
		go func() {
			ctx, cancel := context.WithTimeout(c, 3*time.Second)
			defer cancel()
			err := cache.Set(ctx, key, products)
			if err != nil {
				logger.Error("Error : %v", err)
			}
		}()
	}

	logger.Info("Get all product success")
	responseUtil.Ok(c, products, isFromCache)
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
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	if c.DefaultQuery("no_cache", "0") == "0" {
		err := cache.Get(ctx, key, &product)
		if err != nil {
			logger.Error("Error : %v", err)
		}
	}

	isFromCache := false
	if !product.IsEmpty() {
		isFromCache = true
	} else {
		val, err := pc.productService.GetById(uint(id))
		if err != nil {
			panic(err)
		}
		product = *val
		go func() {
			ctx, cancel := context.WithTimeout(c, 3*time.Second)
			defer cancel()
			err := cache.Set(ctx, key, product)
			if err != nil {
				logger.Error("Error : %v", err)
			}
		}()
	}

	logger.Info(`Get product by id, id = %s success`, c.Param("id"))
	responseUtil.Ok(c, product, isFromCache)
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
	product, err := pc.productService.AddProduct(request, user.UserId)
	if err != nil {
		panic(err)
	}

	logger.Info(`Add new product success`)
	responseUtil.Ok(c, *product, false)
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
	product, err := pc.productService.UpdateProduct(uint(id), request, user.UserId)
	if err != nil {
		panic(err)
	}
	logger.Info(`Update product of id = %s success`, c.Param("id"))
	responseUtil.Ok(c, *product, false)
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
	responseUtil.Ok(c, nil, false)
}

var _ IProductController = (*ProductController)(nil)
