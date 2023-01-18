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
	GetAllProduct(pagination *app.Pagination, count *int64) []*models.Product
	GetByProductId(productId uint) *models.Product
	AddProduct(product models.Product) *models.Product
	UpdateProduct(product models.Product) *models.Product
	DeleteProduct(productId uint)
}

func NewProductRepository() *productRepository {
	logger.Info("Initializing product repository..")
	dbconn := database.DBConnection()
	return &productRepository{
		DB: dbconn,
	}
}

func (repo productRepository) GetAllProduct(pagination *app.Pagination, count *int64) []*models.Product {
	products := []*models.Product{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Preload(clause.Associations).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Find(&products).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		panic(result.Error)
	}

	return products
}

func (repo productRepository) GetByProductId(id uint) *models.Product {
	product := models.Product{}
	result := repo.Preload(clause.Associations).First(&product, "id = ?", id)
	if result.Error != nil {
		panic(result.Error)
	}

	return &product
}

func (repo productRepository) AddProduct(product models.Product) *models.Product {
	result := repo.Create(&product)
	if result.Error != nil {
		panic(result.Error)
	}
	result = repo.Preload(clause.Associations).First(&product, "id = ?", product.ID)
	if result.Error != nil {
		panic(result.Error)
	}
	return &product
}

func (repo productRepository) UpdateProduct(product models.Product) *models.Product {
	result := repo.Updates(&product)
	if result.Error != nil {
		panic(result.Error)
	}
	result = repo.Preload(clause.Associations).First(&product, "id = ?", product.ID)
	if result.Error != nil {
		panic(result.Error)
	}
	return &product
}

func (repo productRepository) DeleteProduct(productId uint) {
	result := repo.Delete(&models.Product{}, productId)
	if result.Error != nil {
		panic(result.Error)
	}
}
