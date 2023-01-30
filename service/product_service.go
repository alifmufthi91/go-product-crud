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

type ProductService interface {
	GetAll(pagination *app.Pagination) *app.PaginatedResult[app.Product]
	GetById(productId uint) *app.Product
	AddProduct(productInput validation.AddProduct, userId uint) *app.Product
	UpdateProduct(productId uint, productInput validation.UpdateProduct, userId uint) *app.Product
	DeleteProduct(productId uint, userId uint)
}

type productService struct {
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
}

func NewProductService(productRepository repository.ProductRepository, userRepository repository.UserRepository) *productService {
	logger.Info("Initializing product service..")
	return &productService{
		productRepository: productRepository,
		userRepository:    userRepository,
	}
}

func (ps productService) GetAll(pagination *app.Pagination) *app.PaginatedResult[app.Product] {
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

func (ps productService) GetById(productId uint) *app.Product {
	logger.Info("Getting product from repository")
	product := ps.productRepository.GetByProductId(productId)
	productData := product.ProductToProduct()
	return &productData
}

func (ps productService) AddProduct(productInput validation.AddProduct, userId uint) *app.Product {
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

func (ps productService) UpdateProduct(productId uint, productInput validation.UpdateProduct, userId uint) *app.Product {
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

func (ps productService) DeleteProduct(productId uint, userId uint) {
	logger.Info(`Deleting product, product_id = %d, user_id = %d`, productId, userId)
	product := ps.productRepository.GetByProductId(productId)
	if product.UploaderId != userId {
		panic(errors.New("user is not allowed to modify this product"))
	}
	ps.productRepository.DeleteProduct(productId)
}

var _ ProductService = (*productService)(nil)
