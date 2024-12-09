package repository

import (
	"log"

	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type MailInterface interface {
	Create(req domain.MailRegister) (*domain.Mail, error)
	Get(id uint) (*domain.Mail, error)
	Delete(id uint) error
	Sent(from string) (*[]domain.Mail, error)
	Inbox(to string) (*[]domain.Mail, error)
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
	request := domain.Mail{
		From:       req.From,
		To:         req.To,
		Title:      req.Title,
		Subtitle:   req.Subtitle,
		Attachment: req.Attachment,
		Message:    req.Message,
	}
	err := m.connection.DB.Create(&request).Error

	if err != nil {
		log.Println(err)
		return nil, web.InternalServerError("Cannot create mail")
	}

	return &request, nil
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

func (m *mail) Sent(from string) (*[]domain.Mail, error) {
	var response []domain.Mail
	err := m.connection.DB.Raw("SELECT * FROM mails WHERE `from` = ?", from).Scan(&response).Error
	if err != nil {
		return nil, web.InternalServerError("Cannot get all sent mail")
	}

	return &response, nil
}

func (m *mail) Inbox(to string) (*[]domain.Mail, error) {
	var response []domain.Mail
	err := m.connection.DB.Raw("SELECT * FROM mails WHERE `to` = ?", to).Scan(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Canot get all inbox mail")
	}

	return &response, nil
}
