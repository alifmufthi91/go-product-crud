package service

import (
	"context"
	"errors"
	"math"
	"product-crud/dto/app"
	"product-crud/dto/request"
	"product-crud/dto/response"
	"product-crud/models"
	"product-crud/repository"
	"product-crud/util/errorhandler"
	"product-crud/util/logger"
	"time"
)

type IProductService interface {
	GetAll(pagination app.Pagination) (app.PaginatedResult[response.GetProductResponse], error)
	GetById(productId uint) (response.GetProductResponse, error)
	AddProduct(productInput request.ProductAddRequest, userId uint) (response.GetProductResponse, error)
	UpdateProduct(productId uint, productInput request.ProductUpdateRequest, userId uint) (response.GetProductResponse, error)
	DeleteProduct(productId uint, userId uint) error
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

func (ps ProductService) GetAll(pagination app.Pagination) (app.PaginatedResult[response.GetProductResponse], error) {
	logger.Info("Getting all product from repository")
	var count int64

	var result app.PaginatedResult[response.GetProductResponse]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	products, err := ps.productRepository.GetAllProduct(ctx, &pagination, &count)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}

	var productDatas []response.GetProductResponse
	for _, x := range products {
		productDatas = append(productDatas, response.NewGetProductResponse(x))
	}

	result = app.PaginatedResult[response.GetProductResponse]{
		Items:      productDatas,
		Page:       pagination.Page,
		Size:       len(productDatas),
		TotalItems: int(count),
		TotalPage:  int(math.Ceil(float64(count) / float64(pagination.Limit))),
	}
	return result, nil
}

func (ps ProductService) GetById(productId uint) (response.GetProductResponse, error) {
	logger.Info("Getting product from repository")

	var result response.GetProductResponse
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nullableProduct, err := ps.productRepository.GetByProductId(ctx, productId)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}
	if !nullableProduct.Valid {
		logger.Error("Invalid product")
		return result, errors.New("Invalid product")
	}
	result = response.NewGetProductResponse(nullableProduct.Stuff)
	return result, nil
}

func (ps ProductService) AddProduct(productInput request.ProductAddRequest, userId uint) (response.GetProductResponse, error) {
	logger.Info(`Adding new product, product = %+v, user_id = %+v`, productInput, userId)

	var result response.GetProductResponse
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product := models.Product{
		ProductName:        productInput.ProductName,
		ProductDescription: productInput.ProductDescription,
		Photo:              productInput.Photo,
		UploaderId:         userId,
	}

	nullableProduct, err := ps.productRepository.AddProduct(ctx, product)
	if !nullableProduct.Valid {
		logger.Error("Invalid product")
		return result, errors.New("Invalid product")
	}
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}
	result = response.NewGetProductResponse(nullableProduct.Stuff)
	return result, nil
}

func (ps ProductService) UpdateProduct(productId uint, productInput request.ProductUpdateRequest, userId uint) (response.GetProductResponse, error) {
	logger.Info(`Updating product, product = %+v, user_id = %d`, productInput, userId)

	var result response.GetProductResponse
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nullableProduct, err := ps.productRepository.GetByProductId(ctx, productId)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}
	if !nullableProduct.Valid {
		logger.Error("Invalid product")
		return result, errors.New("Invalid product")
	}

	product := nullableProduct.Stuff
	if product.UploaderId != userId {
		err := errorhandler.Unauthorized("user is not allowed to modify this product")
		logger.Error("Error : %v", err)
		return result, err
	}
	product.ProductName = productInput.ProductName
	product.ProductDescription = productInput.ProductDescription
	product.Photo = productInput.Photo

	updatedProduct, err := ps.productRepository.UpdateProduct(ctx, product)
	if err != nil {
		logger.Error("Error : %v", err)
		return result, err
	}
	if !updatedProduct.Valid {
		logger.Error("Invalid product")
		return result, errors.New("Invalid product")
	}
	result = response.NewGetProductResponse(updatedProduct.Stuff)
	return result, nil
}

func (ps ProductService) DeleteProduct(productId uint, userId uint) error {
	logger.Info(`Deleting product, product_id = %d, user_id = %d`, productId, userId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nullableProduct, err := ps.productRepository.GetByProductId(ctx, productId)
	if err != nil {
		logger.Error("Error : %v", err)
		return err
	}
	if !nullableProduct.Valid {
		logger.Error("Invalid product")
		return errors.New("Invalid product")
	}

	product := nullableProduct.Stuff
	if product.UploaderId != userId {
		err := errorhandler.Unauthorized("user is not allowed to modify this product")
		logger.Error("Error : %v", err)
		return err
	}

	if err := ps.productRepository.DeleteProduct(ctx, productId); err != nil {
		logger.Error("Error : %v", err)
		return err
	}
	return nil
}

var _ IProductService = (*ProductService)(nil)
