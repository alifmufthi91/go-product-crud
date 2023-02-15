package repository

import (
	"product-crud/app"
	"product-crud/models"
	"product-crud/util/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserRepository interface {
	GetAllUser(pagination *app.Pagination, count *int64) []*models.User
	GetByUserId(userId uint) *models.User
	GetByEmail(email string) *models.User
	AddUser(user models.User) *models.User
	IsExistingEmail(email string) *bool
}

type UserRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	logger.Info("Initializing user repository..")
	return UserRepository{
		DB: db,
	}
}

func (repo UserRepository) GetAllUser(pagination *app.Pagination, count *int64) []*models.User {
	users := []*models.User{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Preload("Products.Uploader").Preload(clause.Associations).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Find(&users).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		panic(result.Error)
	}

	return users
}

func (repo UserRepository) GetByUserId(id uint) *models.User {
	user := models.User{}
	result := repo.Preload("Products.Uploader").Preload(clause.Associations).First(&user, "users.id = ?", id)
	if result.Error != nil {
		panic(result.Error)
	}

	return &user
}

func (repo UserRepository) GetByEmail(email string) *models.User {
	user := models.User{}
	result := repo.Preload(clause.Associations).First(&user, "email = ?", email)
	if result.Error != nil {
		panic(result.Error)
	}

	return &user
}

func (repo UserRepository) IsExistingEmail(email string) *bool {
	var exists bool
	err := repo.Model(models.User{}).Select("count(*) > 0").Where("email = ?", email).Find(&exists).Error
	if err != nil {
		panic(err)
	}
	return &exists
}

func (repo UserRepository) AddUser(user models.User) *models.User {
	result := repo.Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}

	return &user
}

var _ IUserRepository = (*UserRepository)(nil)
