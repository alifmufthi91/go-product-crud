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
	GetAll(pagination *app.Pagination) (*app.PaginatedResult, error)
	GetById(productId uint) (*app.Product, error)
	AddProduct(productInput validation.AddProduct, userId uint) (*app.Product, error)
	UpdateProduct(productId uint, productInput validation.UpdateProduct, userId uint) (*app.Product, error)
	DeleteProduct(productId uint, userId uint) error
}

type productService struct {
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
}

func NewProductService() ProductService {
	logger.Info("Initializing product service..")
	ur := repository.NewProductRepository()
	us := repository.NewUserRepository()
	return productService{
		productRepository: ur,
		userRepository:    us,
	}
}

func (ps productService) GetAll(pagination *app.Pagination) (*app.PaginatedResult, error) {
	logger.Info("Getting all product from repository")
	var count int64
	products, err := ps.productRepository.GetAllProduct(pagination, &count)
	if err != nil {
		return nil, err
	}
	logger.Info(`count: %+d`, count)
	var productDatas []app.Product
	for _, x := range *products {
		productDatas = append(productDatas, x.ProductToProduct())
	}

	paginatedResult := app.PaginatedResult{
		Items:      productDatas,
		Page:       pagination.Page,
		Size:       len(productDatas),
		TotalItems: int(count),
		TotalPage:  int(math.Ceil(float64(count) / float64(pagination.Limit))),
	}

	return &paginatedResult, nil
}

func (ps productService) GetById(productId uint) (*app.Product, error) {
	logger.Info("Getting product from repository")
	product, err := ps.productRepository.GetByProductId(productId)
	if err != nil {
		return nil, err
	}
	productData := product.ProductToProduct()
	return &productData, nil
}

func (ps productService) AddProduct(productInput validation.AddProduct, userId uint) (*app.Product, error) {
	logger.Info(`Adding new product, product = %+v, user_id = %+v`, productInput, userId)
	user, _ := ps.userRepository.GetByUserId(userId)
	if user == nil {
		return nil, errors.New("user is not exists")
	}

	product := models.Product{
		ProductName:        productInput.ProductName,
		ProductDescription: productInput.ProductDescription,
		Photo:              productInput.Photo,
		UploaderId:         user.ID,
	}
	createdProduct, err := ps.productRepository.AddProduct(product)
	if err != nil {
		return nil, err
	}

	logger.Info(`product data = %+v`, createdProduct)
	productData := createdProduct.ProductToProduct()
	return &productData, nil
}

func (ps productService) UpdateProduct(productId uint, productInput validation.UpdateProduct, userId uint) (*app.Product, error) {
	logger.Info(`Updating product, product = %+v, user_id = %d`, productInput, userId)
	product, _ := ps.productRepository.GetByProductId(productId)
	if product == nil {
		return nil, errors.New("product is not exists")
	}
	if product.UploaderId != userId {
		return nil, errors.New("user is not allowed to modify this product")
	}
	product.ProductName = productInput.ProductName
	product.ProductDescription = productInput.ProductDescription
	product.Photo = productInput.Photo
	updatedProduct, err := ps.productRepository.UpdateProduct(*product)
	if err != nil {
		return nil, err
	}

	logger.Info(`product data = %+v`, updatedProduct)
	productData := updatedProduct.ProductToProduct()
	return &productData, nil
}

func (ps productService) DeleteProduct(productId uint, userId uint) error {
	logger.Info(`Deleting product, product_id = %d, user_id = %d`, productId, userId)
	product, _ := ps.productRepository.GetByProductId(productId)
	if product == nil {
		return errors.New("product is not exists")
	}
	if product.UploaderId != userId {
		return errors.New("user is not allowed to modify this product")
	}
	err := ps.productRepository.DeleteProduct(productId)
	if err != nil {
		return err
	}
	return nil
}
