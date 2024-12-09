package domain

import (
	"gorm.io/gorm"
)

type MultiAttachment []string

type Mail struct {
	gorm.Model
	From       string
	To         string
	Message    string `gorm:"type:text"`
	Title      string `gorm:"type:varchar(100)"`
	Subtitle   string `gorm:"type:varchar(100)"`
	Attachment string `gorm:"type:text"`
	Star       bool   `gorm:"default:false"`
}

type MailLabel struct {
	gorm.Model
	CategoryID uint     `json:"category_id" form:"category_id"`
	MailID     uint     `json:"mail_id" form:"mail_id"`
	UserID     uint     `json:"user_id" form:"user_id"`
	Mail       Mail     `json:"mail" gorm:"foreignKey:MailID;references:ID"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
}

type MailRegister struct {
	From       string `json:"from" form:"from"`
	To         string `json:"to" form:"to" validate:"required"`
	Message    string `json:"message" form:"message" validate:"required"`
	Title      string `json:"title" form:"title" validate:"required"`
	Subtitle   string `json:"subtitle" form:"subtitle" validate:"required"`
	Attachment string `json:"attachment" form:"attachment"`
	Star       bool   `json:"star" form:"star"`
}

type MailLabelRegister struct {
	CategoryID uint `json:"category_id" form:"category_id" validate:"required"`
	MailID     uint `json:"mail_id" form:"mail_id" validate:"required"`
	UserID     uint `json:"user_id" form:"user_id" validate:"required"`
}
