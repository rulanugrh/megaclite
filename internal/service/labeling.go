package service

import (
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
)

type LabelingInterface interface {
	Create(req domain.MailLabelRegister) (*web.GetMail, error)
	FindByCategory(id uint, userID uint) (*[]web.GetMail, error)
	UpdateLabel(id uint, categoryID uint) (*web.GetMail, error)
}

type labeling struct {
	repository repository.LabelInterface
	validation middleware.IValidation
}

func NewLabelMailService(repository repository.LabelInterface) LabelingInterface {
	return &labeling{
		repository: repository,
		validation: middleware.NewValidation(),
	}
}

func (l *labeling) Create(req domain.MailLabelRegister) (*web.GetMail, error) {
	err := l.validation.Validate(req)
	if err != nil {
		return nil, l.validation.ValidationMessage(err)
	}

	data, err := l.repository.Create(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetMail{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		Title:     data.Mail.Title,
		Subtitle:  data.Mail.Subtitle,
		From:      data.Mail.From,
		To:        data.Mail.To,
	}

	return &response, nil
}

func (l *labeling) FindByCategory(id uint, userID uint) (*[]web.GetMail, error) {
	var response []web.GetMail
	data, err := l.repository.GetByCategory(id, userID)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	for _, mails := range *data {
		mail := web.GetMail{
			ID:        mails.Mail.ID,
			CreatedAt: mails.Mail.CreatedAt,
			Title:     mails.Mail.Title,
			Subtitle:  mails.Mail.Subtitle,
			From:      mails.Mail.From,
			To:        mails.Mail.To,
		}

		response = append(response, mail)
	}

	return &response, nil
}

func (l *labeling) UpdateLabel(id uint, categoryID uint) (*web.GetMail, error) {
	data, err := l.repository.UpdateLabel(id, categoryID)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetMail{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		Title:     data.Mail.Title,
		Subtitle:  data.Mail.Subtitle,
		From:      data.Mail.From,
		To:        data.Mail.To,
	}

	return &response, nil
}
