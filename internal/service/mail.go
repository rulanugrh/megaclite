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
	Sent(from string) (*[]web.GetMail, error)
	Inbox(to string) (*[]web.GetMail, error)
}

type mail struct {
	repository repository.MailInterface
	validation middleware.IValidation
	pgp        middleware.PGPInterface
}

func NewMailService(repository repository.MailInterface, pgp middleware.PGPInterface) MailInterface {
	return &mail{
		repository: repository,
		validation: middleware.NewValidation(),
		pgp:        pgp,
	}
}

func (m *mail) Create(req domain.MailRegister) (*web.GetDetailMail, error) {
	err := m.validation.Validate(req)
	if err != nil {
		return nil, m.validation.ValidationMessage(err)
	}

	encrypt, err := m.pgp.Encryption(req)
	if err != nil {
		return nil, web.InternalServerError("Error While Encryption Message")
	}

	request := domain.MailRegister{
		From:       req.From,
		To:         req.To,
		Title:      req.Title,
		Message:    string(encrypt),
		Subtitle:   req.Subtitle,
		Star:       req.Star,
		Attachment: req.Attachment,
	}

	data, err := m.repository.Create(request)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetDetailMail{
		From:     data.From,
		To:       data.To,
		Title:    data.Title,
		Message:  req.Message,
		Subtitle: data.Subtitle,
	}

	return &response, nil
}

func (m *mail) FindByID(id uint) (*web.GetDetailMail, error) {
	data, err := m.repository.Get(id)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	decryption, err := m.pgp.Decryption(*data)
	if err != nil {
		return nil, web.InternalServerError("Error while decryption message")
	}

	response := web.GetDetailMail{
		From:      data.From,
		To:        data.To,
		Title:     data.Title,
		Message:   string(decryption),
		Subtitle:  data.Subtitle,
		CreatedAt: data.CreatedAt,
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

func (m *mail) Sent(from string) (*[]web.GetMail, error) {
	var response []web.GetMail

	data, err := m.repository.Sent(from)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	for _, mails := range *data {
		mail := web.GetMail{
			ID:        mails.ID,
			From:      mails.From,
			To:        mails.To,
			Title:     mails.Title,
			Subtitle:  mails.Subtitle,
			CreatedAt: mails.CreatedAt,
		}

		response = append(response, mail)
	}

	return &response, nil
}

func (m *mail) Inbox(from string) (*[]web.GetMail, error) {
	var response []web.GetMail

	data, err := m.repository.Inbox(from)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	for _, mails := range *data {
		mail := web.GetMail{
			ID:        mails.ID,
			From:      mails.From,
			To:        mails.To,
			Title:     mails.Title,
			Subtitle:  mails.Subtitle,
			CreatedAt: mails.CreatedAt,
		}

		response = append(response, mail)
	}

	return &response, nil
}
