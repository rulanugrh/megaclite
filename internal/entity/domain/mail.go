package domain

import "github.com/lib/pq"

type Mail struct {
	From       string
	To         string
	Message    string         `gorm:"type:text"`
	Title      string         `gorm:"type:varchar(100)"`
	Subtitle   string         `gorm:"type:varchar(100)"`
	Attachment pq.StringArray `gorm:"type:text[]"`
	Star       bool           `gorm:"default:false"`
}
