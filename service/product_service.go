package service

import (
	"context"
	"math"
	"product-crud/dto/app"
	"product-crud/dto/request"
	"product-crud/dto/response"
	"product-crud/models"
	"product-crud/repository"
	errorUtil "product-crud/util/error"
	"product-crud/util/logger"
	"time"
)

type IProductService interface {
	GetAll(pagination app.Pagination) app.PaginatedResult[response.GetProductResponse]
	GetById(productId uint) response.GetProductResponse
	AddProduct(productInput request.ProductAddRequest, userId uint) response.GetProductResponse
	UpdateProduct(productId uint, productInput request.ProductUpdateRequest, userId uint) response.GetProductResponse
	DeleteProduct(productId uint, userId uint)
}

type ProductService struct {
	productRepository repository.IProductRepository
	userRepository    repository.IUserRepository
}

func NewProductService(productRepository repository.IProductRepository, userRepository repository.IUserRepository) ProductService {
	logger.Info("Initializing product service..")
	return ProductService{
		productRepository: productRepository,
		userRepository:    userRepository,
	}
}

func (ps ProductService) GetAll(pagination app.Pagination) app.PaginatedResult[response.GetProductResponse] {
	logger.Info("Getting all product from repository")
	var count int64

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	products, err := ps.productRepository.GetAllProduct(ctx, &pagination, &count)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	var productDatas []response.GetProductResponse
	for _, x := range products {
		productDatas = append(productDatas, *response.NewGetProductResponse(x))
	}

	return app.PaginatedResult[response.GetProductResponse]{
		Items:      productDatas,
		Page:       pagination.Page,
		Size:       len(productDatas),
		TotalItems: int(count),
		TotalPage:  int(math.Ceil(float64(count) / float64(pagination.Limit))),
	}
}

func (ps ProductService) GetById(productId uint) response.GetProductResponse {
	logger.Info("Getting product from repository")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product, err := ps.productRepository.GetByProductId(ctx, productId)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	return *response.NewGetProductResponse(product)
}

func (ps ProductService) AddProduct(productInput request.ProductAddRequest, userId uint) response.GetProductResponse {
	logger.Info(`Adding new product, product = %+v, user_id = %+v`, productInput, userId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := ps.userRepository.GetByUserId(ctx, userId)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	product := models.Product{
		ProductName:        productInput.ProductName,
		ProductDescription: productInput.ProductDescription,
		Photo:              productInput.Photo,
		UploaderId:         user.ID,
	}

	createdProduct, err := ps.productRepository.AddProduct(ctx, product)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	return *response.NewGetProductResponse(createdProduct)
}

func (ps ProductService) UpdateProduct(productId uint, productInput request.ProductUpdateRequest, userId uint) response.GetProductResponse {
	logger.Info(`Updating product, product = %+v, user_id = %d`, productInput, userId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product, err := ps.productRepository.GetByProductId(ctx, productId)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	if product.UploaderId != userId {
		err := errorUtil.Unauthorized("user is not allowed to modify this product")
		logger.Error("Error : %v", err)
		panic(err)
	}
	product.ProductName = productInput.ProductName
	product.ProductDescription = productInput.ProductDescription
	product.Photo = productInput.Photo

	updatedProduct, err := ps.productRepository.UpdateProduct(ctx, *product)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	return *response.NewGetProductResponse(updatedProduct)
}

func (ps ProductService) DeleteProduct(productId uint, userId uint) {
	logger.Info(`Deleting product, product_id = %d, user_id = %d`, productId, userId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product, err := ps.productRepository.GetByProductId(ctx, productId)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}

	if product.UploaderId != userId {
		err := errorUtil.Unauthorized("user is not allowed to modify this product")
		logger.Error("Error : %v", err)
		panic(err)
	}

	if err := ps.productRepository.DeleteProduct(ctx, productId); err != nil {
		logger.Error("Error : %v", err)
		panic(err)
	}
}

var _ IProductService = (*ProductService)(nil)
