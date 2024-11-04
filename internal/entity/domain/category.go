package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	MailID   uint   `json:"mail_id"`
	Category string `json:"category"`
	UserID   uint   `json:"user_id"`
	Mails    Mail   `json:"mail" gorm:"foreignKey:MailID;references:ID"`
	Users    User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
