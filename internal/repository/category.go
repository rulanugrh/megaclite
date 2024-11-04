package repository

import (
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type CategoryInterface interface {
	Create(req domain.Category) (*domain.Category, error)
	Update(id uint, req domain.Category) (*domain.Category, error)
	Delete(id uint) error
}

type category struct {
	connection config.Database
}

func NewCategoryRepository(config config.Database) CategoryInterface {
	return &category{
		connection: config,
	}
}

func (c *category) Create(req domain.Category) (*domain.Category, error) {
	var response domain.Category

	err := c.connection.DB.Exec("INSERT INTO categories(name, description) VALUES(?,?)",
		req.Name,
		req.Description,
	).Find(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot create new category")
	}

	return &response, nil
}

func (c *category) Update(id uint, req domain.Category) (*domain.Category, error) {
	var response domain.Category

	err := c.connection.DB.Exec("UPDATE categories SET name = ?, description = ? WHERE id = ?",
		req.Name, req.Description, id,
	).Find(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot update category with this id")
	}

	return &response, nil
}

func (c *category) Delete(id uint) error {
	err := c.connection.DB.Exec("DELETE FROM categories WHERE id = ?", id).Error
	if err != nil {
		return web.InternalServerError("Cannot delete category with this ID")
	}

	return nil
}
