package repository

import (
	"ibf-benevolence/database"
	"ibf-benevolence/entity"
	"ibf-benevolence/util/logger"

	"github.com/jmoiron/sqlx"
)

const (
	querySelectUser = `
		SELECT
			user.user_id,
			user.name,
			user.email,
			user.phone_number_code,
			user.phone_number,
			user.photo_url,
			user.gender,
			user.algo_address,
			user.status,
			user.created_at,
			user.updated_at
		FROM user`

	queryInsertUser = `
		INSERT INTO user
			(user_id,
			name,
			email,
			phone_number_code,
			phone_number,
			photo_url,
			gender,
			algo_address,
			status,
			created_at)
		VALUES (:user_id,
			:name,
			:email,
			:phone_number_code,
			:phone_number,
			:photo_url,
			:gender,
			:algo_address,
			:status,
			:created_at)`

	queryUpdateInventory = `
		UPDATE inventory
		SET
			name = :product_entity_id,
			photo_url = :photo_url,
			gender = :gender
		WHERE user_id = :user_id`
)

type userRepository struct {
	db *sqlx.DB
}

type UserRepository interface {
	FindAllUser() ([]entity.User, error)
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
	err := repo.selectAll(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo userRepository) AddUser(user entity.User) error {
	logger.Info("Add new user to database")
	err := repo.insertUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (repo userRepository) selectAll(dest interface{}) error {
	logger.Info("Querying: " + querySelectUser)
	return repo.db.Select(dest, querySelectUser)
}

func (repo userRepository) insertUser(user entity.User) error {
	logger.Info("Querying: " + queryInsertUser)
	_, err := repo.db.NamedExec(queryInsertUser, user)
	return err
}
