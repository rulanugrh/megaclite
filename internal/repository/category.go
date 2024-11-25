package repository

import (
	"time"

	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type CategoryInterface interface {
	Create(req domain.CategoryRegister) (*domain.Category, error)
	Update(id uint, req domain.CategoryUpdate) (*domain.Category, error)
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

func (c *category) Create(req domain.CategoryRegister) (*domain.Category, error) {
	var response domain.Category

	err := c.connection.DB.Exec("INSERT INTO categories(created_at, updated_at, name, description) VALUES(?,?,?,?)",
		time.Now(),
		time.Now(),
		req.Name,
		req.Description,
	).Find(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot create new category")
	}

	return &response, nil
}

func (c *category) Update(id uint, req domain.CategoryUpdate) (*domain.Category, error) {
	var response domain.Category

	err := c.connection.DB.Exec("UPDATE categories SET name = ?, description = ?, updated_at = ? WHERE id = ?",
		req.Name, req.Description, time.Now(), id,
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
