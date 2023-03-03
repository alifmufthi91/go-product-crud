package repository

import (
	"context"
	"product-crud/dto/app"
	"product-crud/models"
	errorUtil "product-crud/util/error"
	"product-crud/util/logger"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAllUser(ctx context.Context, pagination app.Pagination, count *int64) ([]*models.User, error)
	GetByUserId(ctx context.Context, userId uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	AddUser(ctx context.Context, user models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)
	IsExistingEmail(ctx context.Context, email string) (*bool, error)
}

type UserRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	logger.Info("New user repository..")
	return UserRepository{
		DB: db,
	}
}

func (repo UserRepository) GetAllUser(ctx context.Context, pagination app.Pagination, count *int64) ([]*models.User, error) {
	var users []*models.User
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Preload("Products").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.WithContext(ctx).Find(&users).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (repo UserRepository) GetByUserId(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := repo.WithContext(ctx).Preload("Products").First(&user, "users.id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			panic(errorUtil.DataNotFound("user is not found"))
		}
		return nil, result.Error
	}

	return &user, nil
}

func (repo UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := repo.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			panic(errorUtil.DataNotFound("user is not found"))
		}
		return nil, result.Error
	}

	return &user, nil
}

func (repo UserRepository) IsExistingEmail(ctx context.Context, email string) (*bool, error) {
	var exists bool
	err := repo.WithContext(ctx).Model(models.User{}).Select("count(*) > 0").Where("email = ?", email).Find(&exists).Error
	if err != nil {
		return nil, err
	}
	return &exists, nil
}

func (repo UserRepository) AddUser(ctx context.Context, user models.User) (*models.User, error) {
	result := repo.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo UserRepository) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	result := repo.WithContext(ctx).Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

var _ IUserRepository = (*UserRepository)(nil)
