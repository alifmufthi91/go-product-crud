package models

import (
	app "product-crud/app"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(50)"`
	LastName  string `gorm:"type:varchar(50)"`
	Email     string `gorm:"type:varchar(100);unique_index"`
	Password  []byte
}

func (u *User) UserToUser() app.User {
	return app.User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}
