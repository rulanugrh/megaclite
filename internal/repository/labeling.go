package repository

import (
	"time"

	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type LabelInterface interface {
	Create(req domain.MailLabelRegister) (*domain.MailLabel, error)
	UpdateLabel(id uint, categoryID uint) (*domain.MailLabel, error)
	GetByCategory(categoryID uint, userID uint) (*[]domain.MailLabel, error)
}

type label struct {
	connection config.Database
}

func NewLabelMailRepository(config config.Database) LabelInterface {
	return &label{
		connection: config,
	}
}

func (l *label) Create(req domain.MailLabelRegister) (*domain.MailLabel, error) {
	var response domain.MailLabel
	err := l.connection.DB.Exec("INSERT INTO mail_labels(created_at, updated_at, category_id, user_id, mail_id) VALUES (?,?,?,?,?)",
		time.Now(),
		time.Now(),
		req.CategoryID,
		req.UserID,
		req.MailID,
	).First(&response).Error

	if err != nil {
		return nil, web.InternalServerError("cannot add new labels for mails")
	}

	return &response, nil
}

func (l *label) UpdateLabel(id uint, categoryID uint) (*domain.MailLabel, error) {
	var response domain.MailLabel
	err := l.connection.DB.Exec("UPDATE mail_labels SET category_id = ? WHERE id = ?", categoryID, id).Find(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot update labels for this mail ID")
	}

	return &response, nil
}

func (l *label) GetByCategory(categoryID uint, userID uint) (*[]domain.MailLabel, error) {
	var response []domain.MailLabel
	err := l.connection.DB.Where("category_id = ?", categoryID).Where("user_id = ?", userID).Preload("Mail").Find(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot get mail with this id")
	}

	return &response, nil
}
