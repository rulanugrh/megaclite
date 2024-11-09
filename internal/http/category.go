package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/service"
)

type CategoryInterface interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type category struct {
	service service.CategoryInterface
}

func NewCategoryHandler(service service.CategoryInterface) CategoryInterface {
	return &category{
		service: service,
	}
}

func (ct *category) Create(c *fiber.Ctx) error {
	var request domain.Category
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Request"))
	}

	data, err := ct.service.Create(request)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(201).JSON(web.Created("Success Create New Category", data))
}

func (ct *category) Delete(c *fiber.Ctx) error {
	getID := c.Params("id")

	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing ID"))
	}

	err = ct.service.Delete(uint(id))
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(200).JSON(web.Success("Success Delete Category", nil))
}

func (ct *category) Update(c *fiber.Ctx) error {
	var request domain.Category
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Body Request"))
	}

	getID := c.Params("id")
	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing ID Request"))
	}

	data, err := ct.service.Update(uint(id), request)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(200).JSON(web.Success("Success Update Category", data))
}
