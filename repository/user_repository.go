package repository

import (
	"errors"
	"ibf-benevolence/database"
	"ibf-benevolence/entity"
	"ibf-benevolence/util/logger"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	querySelectUser = `
		SELECT
			user.user_id,
			user.algo_address,
			user.status,
			user.created_at,
			user.updated_at
		FROM user`

	queryFindUser = `
		SELECT
			user.user_id,
			user.algo_address,
			user.status,
			user.created_at,
			user.updated_at
		FROM user
		WHERE user_id = ?`

	queryInsertUser = `
		INSERT INTO user
			(user_id,
			algo_address,
			status,
			created_at)
		VALUES (:user_id,
			:name,
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
	Find(userId string) (*entity.User, error)
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

func (repo userRepository) Find(id string) (*entity.User, error) {
	logger.Info("Find user in database")
	users := []entity.User{}
	err := repo.db.Select(&users, queryFindUser, id)
	log.Printf("%v", users)
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, errors.New("user not found")
	}

	return &users[0], nil
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
