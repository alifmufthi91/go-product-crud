package repository

import (
	"context"
	"product-crud/dto/app"
	"product-crud/models"
	errorUtil "product-crud/util/error"
	"product-crud/util/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IProductRepository interface {
	GetAllProduct(ctx context.Context, pagination *app.Pagination, count *int64) ([]*models.Product, error)
	GetByProductId(ctx context.Context, productId uint) (*models.Product, error)
	AddProduct(ctx context.Context, product models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, productId uint) error
}

type ProductRepository struct {
	*gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	logger.Info("New product repository..")
	return ProductRepository{
		DB: db,
	}
}

func (repo ProductRepository) GetAllProduct(ctx context.Context, pagination *app.Pagination, count *int64) ([]*models.Product, error) {
	products := []*models.Product{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Preload(clause.Associations).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.WithContext(ctx).Find(&products).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (repo ProductRepository) GetByProductId(ctx context.Context, id uint) (*models.Product, error) {
	product := models.Product{}
	result := repo.WithContext(ctx).Preload(clause.Associations).First(&product, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			panic(errorUtil.DataNotFound("product is not found"))
		}
		return nil, result.Error
	}

	return &product, nil
}

func (repo ProductRepository) AddProduct(ctx context.Context, product models.Product) (*models.Product, error) {
	result := repo.WithContext(ctx).Create(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	// result = repo.WithContext(ctx).Preload(clause.Associations).First(&product, "id = ?", product.ID)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }
	return &product, nil
}

func (repo ProductRepository) UpdateProduct(ctx context.Context, product models.Product) (*models.Product, error) {
	result := repo.WithContext(ctx).Updates(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	// result = repo.WithContext(ctx).Preload(clause.Associations).First(&product, "id = ?", product.ID)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }
	return &product, nil
}

func (repo ProductRepository) DeleteProduct(ctx context.Context, productId uint) error {
	result := repo.WithContext(ctx).Delete(&models.Product{}, productId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var _ IProductRepository = (*ProductRepository)(nil)
