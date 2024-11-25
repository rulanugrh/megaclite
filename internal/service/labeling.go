package service

import (
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
)

type LabelingInterface interface {
	Create(req domain.MailLabelRegister) (*web.GetMailLabel, error)
	FindByID(id uint) (*web.GetMailLabel, error)
	FindByCategory(id uint, userID uint) (*[]web.GetMailLabel, error)
	UpdateLabel(id uint, categoryID uint) (*web.GetMailLabel, error)
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

func (l *labeling) Create(req domain.MailLabelRegister) (*web.GetMailLabel, error) {
	err := l.validation.Validate(req)
	if err != nil {
		return nil, l.validation.ValidationMessage(err)
	}

	data, err := l.repository.Create(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetMailLabel{
		Category: data.Category.Name,
		Message:  data.Mail.Message,
		Title:    data.Mail.Title,
		Subtitle: data.Mail.Subtitle,
		From:     data.Mail.From,
		To:       data.Mail.To,
	}

	return &response, nil
}

func (l *labeling) FindByID(id uint) (*web.GetMailLabel, error) {
	data, err := l.repository.Get(id)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetMailLabel{
		Category: data.Category.Name,
		Message:  data.Mail.Message,
		Title:    data.Mail.Title,
		Subtitle: data.Mail.Subtitle,
		From:     data.Mail.From,
		To:       data.Mail.To,
	}

	return &response, nil
}

func (l *labeling) FindByCategory(id uint, userID uint) (*[]web.GetMailLabel, error) {
	var response []web.GetMailLabel
	data, err := l.repository.GetByCategory(id, userID)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	for _, mails := range *data {
		mail := web.GetMailLabel{
			Category: mails.Category.Name,
			Message:  mails.Mail.Message,
			Title:    mails.Mail.Title,
			Subtitle: mails.Mail.Subtitle,
			From:     mails.Mail.From,
			To:       mails.Mail.To,
		}

		response = append(response, mail)
	}

	return &response, nil
}

func (l *labeling) UpdateLabel(id uint, categoryID uint) (*web.GetMailLabel, error) {
	data, err := l.repository.UpdateLabel(id, categoryID)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetMailLabel{
		Category: data.Category.Name,
		Message:  data.Mail.Message,
		Title:    data.Mail.Title,
		Subtitle: data.Mail.Subtitle,
		From:     data.Mail.From,
		To:       data.Mail.To,
	}

	return &response, nil
}
