package service

import (
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/repository"
)

type CategoryInterface interface {
	Create(req domain.Category) (*web.Category, error)
	Update(id uint, req domain.Category) (*web.Category, error)
	Delete(id uint) error
}

type category struct {
	repository repository.CategoryInterface
}

func NewCategoryService(repository repository.CategoryInterface) CategoryInterface {
	return &category{
		repository: repository,
	}
}

func (c *category) Create(req domain.Category) (*web.Category, error) {
	data, err := c.repository.Create(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.Category{
		Name:        data.Name,
		Description: data.Description,
	}

	return &response, nil
}

func (c *category) Update(id uint, req domain.Category) (*web.Category, error) {
	data, err := c.repository.Update(id, req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.Category{
		Name:        data.Name,
		Description: data.Description,
	}

	return &response, nil
}

func (c *category) Delete(id uint) error {
	err := c.repository.Delete(id)
	if err != nil {
		return web.InternalServerError(err.Error())
	}

	return nil
}
