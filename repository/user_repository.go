package repository

import (
	"fmt"
	"ibf-benevolence/entity"
)

type userRepository struct {
	base BaseRepository
}

type UserRepository interface {
	FindAllUser() ([]entity.User, error)
}

func NewUserRepository() UserRepository {
	fmt.Println("Initializing user repository")
	br := NewRepository("user", "user_id", "usr")
	ur := userRepository{base: br}
	return ur
}

func (repo userRepository) FindAllUser() ([]entity.User, error) {
	fmt.Println("Find all user in database")
	rows, err := repo.base.findAll()
	if err != nil {
		return nil, err
	}
	users := []entity.User{}
	for rows.Next() {
		var r entity.User
		err = rows.Scan(&r.UserId, &r.Name, &r.Email, &r.PhoneNumberCode, &r.PhoneNumber,
			&r.PhotoUrl, &r.Gender, &r.AlgoAddress, &r.Status, &r.CreatedAt, &r.UpdatedAt)
		// fmt.Printf("%+v\n", r)
		if err != nil {
			return nil, err
		}
		users = append(users, r)
	}
	return users, nil
}
