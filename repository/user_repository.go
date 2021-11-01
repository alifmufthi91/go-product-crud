package repository

import (
	"ibf-benevolence/database"
	"ibf-benevolence/entity"
	"ibf-benevolence/util/logger"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	FindAllUser() ([]entity.User, error)
	FindByUserId(userId string) (*entity.User, error)
	AddUser(user entity.User) error
}

func NewUserRepository() UserRepository {
	logger.Info("Initializing user repository..")
	dbconn := database.DBConnection()
	return userRepository{
		db: dbconn,
	}
}

func (repo userRepository) FindAllUser() ([]entity.User, error) {
	logger.Info("Find all user in database")
	users := []entity.User{}
	result := repo.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (repo userRepository) FindByUserId(id string) (*entity.User, error) {
	logger.Info("Find user in database")
	user := entity.User{}
	result := repo.db.First(&user, "user_id = ?", id)
	log.Printf("%v", user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo userRepository) AddUser(user entity.User) error {
	logger.Info("Add new user to database")
	result := repo.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
