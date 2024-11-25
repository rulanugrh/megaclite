package repository

import (
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type LabelInterface interface {
	Create(req domain.MailLabelRegister) (*domain.MailLabel, error)
	Get(id uint) (*domain.MailLabel, error)
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
	err := l.connection.DB.Exec("INSERT INTO mail_labels(category_id, user_id, mail_id) VALUES (?,?,?)",
		req.CategoryID,
		req.UserID,
		req.MailID,
	).Find(&response).Error

	if err != nil {
		return nil, web.InternalServerError("cannot add new labels for mails")
	}

	return &response, nil
}

func (l *label) Get(id uint) (*domain.MailLabel, error) {
	var response domain.MailLabel
	err := l.connection.DB.Raw("SELECT * FROM mail_labels WHERE id = ?", id).Scan(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot get mail with this id")
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
	err := l.connection.DB.Raw("SELECT * FROM mail_labels WHERE category_id = ? AND user_id = ?", categoryID, userID).Scan(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot get mail with this id")
	}

	return &response, nil
}
