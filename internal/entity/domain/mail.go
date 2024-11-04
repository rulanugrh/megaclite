package domain

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Mail struct {
	gorm.Model
	From       string
	To         string
	Message    string         `gorm:"type:text"`
	Title      string         `gorm:"type:varchar(100)"`
	Subtitle   string         `gorm:"type:varchar(100)"`
	Attachment pq.StringArray `gorm:"type:text[];default:null"`
	Star       bool           `gorm:"default:false"`
}

type MailLabel struct {
	gorm.Model
	CategoryID uint     `json:"category_id" form:"category_id"`
	MailID     uint     `json:"mail_id" form:"mail_id"`
	UserID     uint     `json:"user_id" form:"user_id"`
	Mail       Mail     `json:"mail" gorm:"foreignKey:MailID;references:ID"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
}
