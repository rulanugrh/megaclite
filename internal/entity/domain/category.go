package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	ListMailLabel []MailLabel `json:"mail_label" form:"mail_label" gorm:"many2many:mail_label_megaclite"`
}

type CategoryRegister struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}

type CategoryUpdate struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}
