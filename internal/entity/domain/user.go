package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50)"`
	Email    string `gorm:"type:varchar(50)"`
	Password string `gorm:"type:varchar(256)"`
	Avatar   string `gorm:"type:varchar(200)"`
	Address  string `gorm:"type:text"`
}

type Register struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
