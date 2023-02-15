package repository

import (
	"product-crud/app"
	"product-crud/models"
	"product-crud/util/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IProductRepository interface {
	GetAllProduct(pagination *app.Pagination, count *int64) []*models.Product
	GetByProductId(productId uint) *models.Product
	AddProduct(product models.Product) *models.Product
	UpdateProduct(product models.Product) *models.Product
	DeleteProduct(productId uint)
}

type ProductRepository struct {
	*gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	logger.Info("Initializing product repository..")
	return ProductRepository{
		DB: db,
	}
}

func (repo ProductRepository) GetAllProduct(pagination *app.Pagination, count *int64) []*models.Product {
	products := []*models.Product{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Preload(clause.Associations).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Find(&products).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		panic(result.Error)
	}

	return products
}

func (repo ProductRepository) GetByProductId(id uint) *models.Product {
	product := models.Product{}
	result := repo.Preload(clause.Associations).First(&product, "id = ?", id)
	if result.Error != nil {
		panic(result.Error)
	}

	return &product
}

func (repo ProductRepository) AddProduct(product models.Product) *models.Product {
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

func (repo ProductRepository) UpdateProduct(product models.Product) *models.Product {
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

func (repo ProductRepository) DeleteProduct(productId uint) {
	result := repo.Delete(&models.Product{}, productId)
	if result.Error != nil {
		panic(result.Error)
	}
}

var _ IProductRepository = (*ProductRepository)(nil)
