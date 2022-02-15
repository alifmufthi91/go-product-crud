package repository

import (
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
	GetAllUser() ([]models.User, error)
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

func (repo userRepository) GetAllUser() ([]models.User, error) {
	users := []models.User{}
	result := repo.Preload(clause.Associations).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (repo userRepository) GetByUserId(id uint) (*models.User, error) {
	user := models.User{}
	result := repo.Preload(clause.Associations).First(&user, "id = ?", id)
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
