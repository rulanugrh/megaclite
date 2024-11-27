package repository

import (
	"time"

	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type MailInterface interface {
	Create(req domain.MailRegister) (*domain.Mail, error)
	Get(id uint) (*domain.Mail, error)
	Delete(id uint) error
	GetAll(from string) (*[]domain.Mail, error)
}

type mail struct {
	connection config.Database
}

func NewMailRepository(conn config.Database) MailInterface {
	return &mail{
		connection: conn,
	}
}

func (m *mail) Create(req domain.MailRegister) (*domain.Mail, error) {
	var response domain.Mail
	err := m.connection.DB.Exec("INSERT INTO mails(created_at, updated_at, `from`, `to`, message, title, subtitle) VALUES(?, ?, ?, ?, ?, ?, ?)",
		time.Now(),
		time.Now(),
		req.From,
		req.To,
		req.Message,
		req.Title,
		req.Subtitle,
	).Find(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot create mail")
	}

	return &response, nil
}

func (m *mail) Get(id uint) (*domain.Mail, error) {
	var response domain.Mail
	err := m.connection.DB.Raw("SELECT * FROM mails WHERE id = ?", id).Scan(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot find mail with this id")
	}

	return &response, nil
}

func (m *mail) Delete(id uint) error {
	err := m.connection.DB.Exec("DELETE FROM mails WHERE id = ?", id).Error
	if err != nil {
		return web.InternalServerError("Cannot delete this mail")
	}

	return nil
}

func (m *mail) GetAll(from string) (*[]domain.Mail, error) {
	var response []domain.Mail
	err := m.connection.DB.Raw("SELECT * FROM mails WHERE from = ?", from).Scan(&response).Error
	if err != nil {
		return nil, web.InternalServerError("Cannot get all mail with this email")
	}

	return &response, nil
}
