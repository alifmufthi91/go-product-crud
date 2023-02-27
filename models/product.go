package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductName        string `gorm:"type:varchar(100)"`
	ProductDescription string
	Photo              string `gorm:"type:varchar(100)"`
	UploaderId         uint
	Uploader           *User `gorm:"foreignKey:UploaderId"`
}
