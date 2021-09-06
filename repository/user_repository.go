package repository

import (
	"errors"
	"ibf-benevolence/entity"
	"ibf-benevolence/util/logger"
)

type userRepository struct {
	BaseRepository
}

type UserRepository interface {
	FindAllUser() ([]entity.User, error)
}

func NewUserRepository() UserRepository {
	logger.Info("Initializing user repository..")
	br := NewRepository("user", "user_id", "usr")
	ur := userRepository{br}
	return ur
}

func (repo userRepository) FindAllUser() ([]entity.User, error) {
	logger.Info("Find all user in database")
	rows, err := repo.findAll()
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find from database")
	}
	users := []entity.User{}
	for rows.Next() {
		var r entity.User
		err = rows.Scan(&r.UserId, &r.Name, &r.Email, &r.PhoneNumberCode, &r.PhoneNumber,
			&r.PhotoUrl, &r.Gender, &r.AlgoAddress, &r.Status, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to scan data from rows")
		}
		users = append(users, r)
	}
	return users, nil
}
