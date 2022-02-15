package models

import (
	"product-crud/app"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductName        string `gorm:"type:varchar(100)"`
	ProductDescription string
	Photo              string `gorm:"type:varchar(100)"`
	UploaderId         uint
	Uploader           User `gorm:"foreignKey:UploaderId"`
}

func (p *Product) ProductToProduct() app.Product {
	return app.Product{
		ID:                 p.ID,
		ProductName:        p.ProductName,
		ProductDescription: p.ProductDescription,
		Photo:              p.Photo,
		UploaderId:         p.UploaderId,
		Uploader:           p.Uploader.UserToUser(),
	}
}
