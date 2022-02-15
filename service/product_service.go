package service

import (
	"errors"
	"product-crud/app"
	"product-crud/models"
	"product-crud/repository"
	"product-crud/util/logger"
	"product-crud/validation"
)

type ProductService interface {
	GetAll() ([]app.Product, error)
	GetById(productId string) (*app.Product, error)
	AddProduct(productInput validation.AddProduct, userId uint) (*app.Product, error)
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

func (ps productService) GetAll() ([]app.Product, error) {
	logger.Info("Getting all product from repository")
	products, err := ps.productRepository.GetAllProduct()
	if err != nil {
		return nil, err
	}
	var productDatas []app.Product
	for _, x := range products {
		productDatas = append(productDatas, x.ProductToProduct())
	}

	return productDatas, nil
}

func (ps productService) GetById(productId string) (*app.Product, error) {
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
