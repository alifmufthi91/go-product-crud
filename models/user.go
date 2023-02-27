package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string    `gorm:"type:varchar(50)"`
	LastName  string    `gorm:"type:varchar(50)"`
	Email     string    `gorm:"type:varchar(100);unique_index"`
	Password  []byte    `gorm:"<-:create"`
	Products  []Product `gorm:"foreignKey:UploaderId"`
}
