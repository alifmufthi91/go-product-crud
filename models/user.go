package models

import (
	app "product-crud/app"

	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(50)"`
	LastName  string `gorm:"type:varchar(50)"`
	Email     string `gorm:"type:varchar(100);unique_index"`
	Password  []byte
	Products  []Product `gorm:"foreignKey:UploaderId"`
}

func (u *User) UserToUser() app.User {
	var productDatas []app.Product
	for _, x := range u.Products {
		productDatas = append(productDatas, x.ProductToProduct())
	}
	return app.User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Products:  productDatas,
	}
}

func (res User) IsEmpty() bool {
	return cmp.Equal(res, User{})
}
