package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                 uint   `gorm:"primarykey"`
	ProductName        string `gorm:"type:varchar(100)"`
	ProductDescription string
	Photo              string `gorm:"type:varchar(100)"`
	UploaderId         uint
	Uploader           User `gorm:"foreignKey:UploaderId"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}
