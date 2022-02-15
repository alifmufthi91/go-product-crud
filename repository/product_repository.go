package repository

import (
	"product-crud/database"
	"product-crud/models"
	"product-crud/util/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	db *gorm.DB
}

type ProductRepository interface {
	GetAllProduct() ([]models.Product, error)
	GetByProductId(productId uint) (*models.Product, error)
	AddProduct(product models.Product) (*models.Product, error)
}

func NewProductRepository() ProductRepository {
	logger.Info("Initializing product repository..")
	dbconn := database.DBConnection()
	return productRepository{
		db: dbconn,
	}
}

func (repo productRepository) GetAllProduct() ([]models.Product, error) {
	products := []models.Product{}
	result := repo.db.Preload(clause.Associations).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (repo productRepository) GetByProductId(id uint) (*models.Product, error) {
	product := models.Product{}
	result := repo.db.Preload(clause.Associations).First(&product, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (repo productRepository) AddProduct(product models.Product) (*models.Product, error) {
	result := repo.db.Create(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	result = repo.db.Preload(clause.Associations).First(&product, "id = ?", product.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}
