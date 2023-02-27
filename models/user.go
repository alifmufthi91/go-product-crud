package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primarykey"`
	FirstName string    `gorm:"type:varchar(50)"`
	LastName  string    `gorm:"type:varchar(50)"`
	Email     string    `gorm:"type:varchar(100);unique_index"`
	Password  []byte    `gorm:"->:false;<-:create"`
	Products  []Product `gorm:"foreignKey:UploaderId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
