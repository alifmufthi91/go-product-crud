package repository

import (
	"context"
	"product-crud/dto/app"
	"product-crud/models"
	"product-crud/util/errorhandler"
	"product-crud/util/logger"

	"gorm.io/gorm"
)

type IProductRepository interface {
	GetAllProduct(ctx context.Context, pagination *app.Pagination, count *int64) ([]models.Product, error)
	GetByProductId(ctx context.Context, productId uint) (app.NullableStuff[models.Product], error)
	AddProduct(ctx context.Context, product models.Product) (app.NullableStuff[models.Product], error)
	UpdateProduct(ctx context.Context, product models.Product) (app.NullableStuff[models.Product], error)
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

func (repo ProductRepository) GetAllProduct(ctx context.Context, pagination *app.Pagination, count *int64) ([]models.Product, error) {
	var products []models.Product
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Joins("Uploader").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.WithContext(ctx).Find(&products).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (repo ProductRepository) GetByProductId(ctx context.Context, id uint) (app.NullableStuff[models.Product], error) {
	var product models.Product
	var nullableProduct app.NullableStuff[models.Product]
	result := repo.WithContext(ctx).Joins("Uploader").First(&product, "products.id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nullableProduct, errorhandler.DataNotFound("product is not found")
		}
		return nullableProduct, result.Error
	}
	nullableProduct = app.NewNullableStuff(product)
	return nullableProduct, nil
}

func (repo ProductRepository) AddProduct(ctx context.Context, product models.Product) (app.NullableStuff[models.Product], error) {
	result := repo.WithContext(ctx).Create(&product)
	var nullableProduct app.NullableStuff[models.Product]
	if result.Error != nil {
		return nullableProduct, result.Error
	}
	nullableProduct = app.NewNullableStuff(product)
	return nullableProduct, nil
}

func (repo ProductRepository) UpdateProduct(ctx context.Context, product models.Product) (app.NullableStuff[models.Product], error) {
	result := repo.WithContext(ctx).Updates(&product)
	var nullableProduct app.NullableStuff[models.Product]
	if result.Error != nil {
		return nullableProduct, result.Error
	}

	nullableProduct = app.NewNullableStuff(product)
	return nullableProduct, nil
}

func (repo ProductRepository) DeleteProduct(ctx context.Context, productId uint) error {
	result := repo.WithContext(ctx).Delete(&models.Product{}, productId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var _ IProductRepository = (*ProductRepository)(nil)
