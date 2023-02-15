package service

import (
	"errors"
	"math"
	"product-crud/app"
	"product-crud/models"
	"product-crud/repository"
	"product-crud/util/logger"
	"product-crud/validation"
)

type IProductService interface {
	GetAll(pagination *app.Pagination) *app.PaginatedResult[app.Product]
	GetById(productId uint) *app.Product
	AddProduct(productInput validation.AddProduct, userId uint) *app.Product
	UpdateProduct(productId uint, productInput validation.UpdateProduct, userId uint) *app.Product
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

func (ps ProductService) GetAll(pagination *app.Pagination) *app.PaginatedResult[app.Product] {
	logger.Info("Getting all product from repository")
	var count int64
	products := ps.productRepository.GetAllProduct(pagination, &count)

	logger.Info(`count: %+d`, count)
	var productDatas []app.Product
	for _, x := range products {
		productDatas = append(productDatas, x.ProductToProduct())
	}

	paginatedResult := app.PaginatedResult[app.Product]{
		Items:      productDatas,
		Page:       pagination.Page,
		Size:       len(productDatas),
		TotalItems: int(count),
		TotalPage:  int(math.Ceil(float64(count) / float64(pagination.Limit))),
	}

	return &paginatedResult
}

func (ps ProductService) GetById(productId uint) *app.Product {
	logger.Info("Getting product from repository")
	product := ps.productRepository.GetByProductId(productId)
	productData := product.ProductToProduct()
	return &productData
}

func (ps ProductService) AddProduct(productInput validation.AddProduct, userId uint) *app.Product {
	logger.Info(`Adding new product, product = %+v, user_id = %+v`, productInput, userId)
	user := ps.userRepository.GetByUserId(userId)

	product := models.Product{
		ProductName:        productInput.ProductName,
		ProductDescription: productInput.ProductDescription,
		Photo:              productInput.Photo,
		UploaderId:         user.ID,
	}
	createdProduct := ps.productRepository.AddProduct(product)

	productData := createdProduct.ProductToProduct()
	return &productData
}

func (ps ProductService) UpdateProduct(productId uint, productInput validation.UpdateProduct, userId uint) *app.Product {
	logger.Info(`Updating product, product = %+v, user_id = %d`, productInput, userId)
	product := ps.productRepository.GetByProductId(productId)
	if product.UploaderId != userId {
		panic(errors.New("user is not allowed to modify this product"))
	}
	product.ProductName = productInput.ProductName
	product.ProductDescription = productInput.ProductDescription
	product.Photo = productInput.Photo
	updatedProduct := ps.productRepository.UpdateProduct(*product)

	productData := updatedProduct.ProductToProduct()
	return &productData
}

func (ps ProductService) DeleteProduct(productId uint, userId uint) {
	logger.Info(`Deleting product, product_id = %d, user_id = %d`, productId, userId)
	product := ps.productRepository.GetByProductId(productId)
	if product.UploaderId != userId {
		panic(errors.New("user is not allowed to modify this product"))
	}
	ps.productRepository.DeleteProduct(productId)
}

var _ IProductService = (*ProductService)(nil)
