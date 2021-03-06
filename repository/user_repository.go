package repository

import (
	"product-crud/app"
	"product-crud/database"
	"product-crud/models"
	"product-crud/util/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	*gorm.DB
}

type UserRepository interface {
	GetAllUser(pagination *app.Pagination, count *int64) ([]*models.User, error)
	GetByUserId(userId uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	AddUser(user models.User) (*models.User, error)
	IsExistingEmail(email string) (*bool, error)
}

func NewUserRepository() UserRepository {
	logger.Info("Initializing user repository..")
	dbconn := database.DBConnection()
	return userRepository{
		DB: dbconn,
	}
}

func (repo userRepository) GetAllUser(pagination *app.Pagination, count *int64) ([]*models.User, error) {
	users := []*models.User{}
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := repo.Preload("Products.Uploader").Preload(clause.Associations).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Find(&users).Limit(-1).Offset(-1).Count(count)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (repo userRepository) GetByUserId(id uint) (*models.User, error) {
	user := models.User{}
	result := repo.Preload("Products.Uploader").Preload(clause.Associations).First(&user, "users.id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo userRepository) GetByEmail(email string) (*models.User, error) {
	user := models.User{}
	result := repo.Preload(clause.Associations).First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo userRepository) IsExistingEmail(email string) (*bool, error) {
	var exists bool
	err := repo.Model(models.User{}).Select("count(*) > 0").Where("email = ?", email).Find(&exists).Error
	if err != nil {
		return nil, err
	}
	return &exists, nil
}

func (repo userRepository) AddUser(user models.User) (*models.User, error) {
	result := repo.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
