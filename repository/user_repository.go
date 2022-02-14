package repository

import (
	"log"
	"product-crud/database"
	"product-crud/models"
	"product-crud/util/logger"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	GetAllUser() ([]models.User, error)
	GetByUserId(userId string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	AddUser(user models.User) (*models.User, error)
	IsExistingEmail(email string) (*bool, error)
}

func NewUserRepository() UserRepository {
	logger.Info("Initializing user repository..")
	dbconn := database.DBConnection()
	return userRepository{
		db: dbconn,
	}
}

func (repo userRepository) GetAllUser() ([]models.User, error) {
	logger.Info("Get all user in database")
	users := []models.User{}
	result := repo.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (repo userRepository) GetByUserId(id string) (*models.User, error) {
	logger.Info("Get user in database by id")
	user := models.User{}
	result := repo.db.First(&user, "user_id = ?", id)
	log.Printf("%v", user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo userRepository) GetByEmail(email string) (*models.User, error) {
	logger.Info("Get user in database by email")
	user := models.User{}
	result := repo.db.First(&user, "email = ?", email)
	log.Printf("%v", user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo userRepository) IsExistingEmail(email string) (*bool, error) {
	logger.Info("Find existing user in database by email")
	var exists bool
	err := repo.db.Model(models.User{}).Select("count(*) > 0").Where("email = ?", email).Find(&exists).Error
	if err != nil {
		return nil, err
	}
	return &exists, nil
}

func (repo userRepository) AddUser(user models.User) (*models.User, error) {
	logger.Info("Add new user to database")
	result := repo.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
