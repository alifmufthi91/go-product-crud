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
	GetAllUser(ctx context.Context, pagination app.Pagination, count *int64) ([]models.User, error)
	GetByUserId(ctx context.Context, userId uint) (app.NullableStuff[models.User], error)
	GetByEmail(ctx context.Context, email string) (app.NullableStuff[models.User], error)
	AddUser(ctx context.Context, user models.User) (app.NullableStuff[models.User], error)
	UpdateUser(ctx context.Context, user models.User) (app.NullableStuff[models.User], error)
	IsExistingEmail(ctx context.Context, email string) (bool, error)
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

func (repo UserRepository) GetAllUser(ctx context.Context, pagination app.Pagination, count *int64) ([]models.User, error) {
	var users []models.User
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Preload("Products").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.WithContext(ctx).Find(&users).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (repo UserRepository) GetByUserId(ctx context.Context, id uint) (app.NullableStuff[models.User], error) {
	var user models.User
	var nullableUser app.NullableStuff[models.User]
	result := repo.WithContext(ctx).Preload("Products").First(&user, "users.id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nullableUser, errorUtil.DataNotFound("user is not found")
		}
		return nullableUser, result.Error
	}
	nullableUser = app.NewNullableStuff(user)
	return nullableUser, nil
}

func (repo UserRepository) GetByEmail(ctx context.Context, email string) (app.NullableStuff[models.User], error) {
	var user models.User
	var nullableUser app.NullableStuff[models.User]
	result := repo.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nullableUser, errorUtil.DataNotFound("user is not found")
		}
		return nullableUser, result.Error
	}
	nullableUser = app.NewNullableStuff(user)
	return nullableUser, nil
}

func (repo UserRepository) IsExistingEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := repo.WithContext(ctx).Model(models.User{}).Select("count(*) > 0").Where("email = ?", email).Find(&exists).Error
	if err != nil {
		return exists, err
	}
	return exists, nil
}

func (repo UserRepository) AddUser(ctx context.Context, user models.User) (app.NullableStuff[models.User], error) {
	result := repo.WithContext(ctx).Create(&user)
	var nullableUser app.NullableStuff[models.User]
	if result.Error != nil {
		return nullableUser, result.Error
	}
	nullableUser = app.NewNullableStuff(user)
	return nullableUser, nil
}

func (repo UserRepository) UpdateUser(ctx context.Context, user models.User) (app.NullableStuff[models.User], error) {
	result := repo.WithContext(ctx).Save(&user)
	var nullableUser app.NullableStuff[models.User]
	if result.Error != nil {
		return nullableUser, result.Error
	}
	nullableUser = app.NewNullableStuff(user)
	return nullableUser, nil
}

var _ IUserRepository = (*UserRepository)(nil)
