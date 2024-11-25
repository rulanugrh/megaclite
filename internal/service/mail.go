package service

import (
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
)

type MailInterface interface {
	Create(req domain.MailRegister) (*web.GetDetailMail, error)
	FindByID(id uint) (*web.GetDetailMail, error)
	Delete(id uint) error
	Get(from string) (*[]web.GetDetailMail, error)
}

type mail struct {
	repository repository.MailInterface
	validation middleware.IValidation
}

func NewMailService(repository repository.MailInterface) MailInterface {
	return &mail{
		repository: repository,
		validation: middleware.NewValidation(),
	}
}

func (m *mail) Create(req domain.MailRegister) (*web.GetDetailMail, error) {
	err := m.validation.Validate(req)
	if err != nil {
		return nil, m.validation.ValidationMessage(err)
	}

	data, err := m.repository.Create(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetDetailMail{
		From:     data.From,
		To:       data.To,
		Title:    data.Title,
		Message:  data.Message,
		Subtitle: data.Subtitle,
	}

	return &response, nil
}

func (m *mail) FindByID(id uint) (*web.GetDetailMail, error) {
	data, err := m.repository.Get(id)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetDetailMail{
		From:     data.From,
		To:       data.To,
		Title:    data.Title,
		Message:  data.Message,
		Subtitle: data.Subtitle,
	}

	return &response, nil
}

func (m *mail) Delete(id uint) error {
	err := m.repository.Delete(id)
	if err != nil {
		return web.InternalServerError(err.Error())
	}

	return nil
}

func (m *mail) Get(from string) (*[]web.GetDetailMail, error) {
	var response []web.GetDetailMail

	data, err := m.repository.GetAll(from)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	for _, mails := range *data {
		mail := web.GetDetailMail{
			From:     mails.From,
			To:       mails.To,
			Title:    mails.Title,
			Subtitle: mails.Subtitle,
			Message:  mails.Subtitle,
		}

		response = append(response, mail)
	}

	return &response, nil
}
