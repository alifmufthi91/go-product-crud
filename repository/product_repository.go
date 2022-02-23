package repository

import (
	"product-crud/app"
	"product-crud/database"
	"product-crud/models"
	"product-crud/util/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	*gorm.DB
}

type ProductRepository interface {
	GetAllProduct(pagination *app.Pagination, count *int64) ([]*models.Product, error)
	GetByProductId(productId uint) (*models.Product, error)
	AddProduct(product models.Product) (*models.Product, error)
	UpdateProduct(product models.Product) (*models.Product, error)
	DeleteProduct(productId uint) error
}

func NewProductRepository() ProductRepository {
	logger.Info("Initializing product repository..")
	dbconn := database.DBConnection()
	return productRepository{
		DB: dbconn,
	}
}

func (repo productRepository) GetAllProduct(pagination *app.Pagination, count *int64) ([]*models.Product, error) {
	products := []*models.Product{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Preload(clause.Associations).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Find(&products).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (repo productRepository) GetByProductId(id uint) (*models.Product, error) {
	product := models.Product{}
	result := repo.Preload(clause.Associations).First(&product, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (repo productRepository) AddProduct(product models.Product) (*models.Product, error) {
	result := repo.Create(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	result = repo.Preload(clause.Associations).First(&product, "id = ?", product.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (repo productRepository) UpdateProduct(product models.Product) (*models.Product, error) {
	result := repo.Updates(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	result = repo.Preload(clause.Associations).First(&product, "id = ?", product.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (repo productRepository) DeleteProduct(productId uint) error {
	result := repo.Delete(&models.Product{}, productId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
