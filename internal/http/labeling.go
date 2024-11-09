package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/service"
)

type LabelingInterface interface {
	Create(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
	FindByCategory(c *fiber.Ctx) error
	UpdateLabel(c *fiber.Ctx) error
}

type labeling struct {
	service service.LabelingInterface
}

func NewLabelMailHandler(service service.LabelingInterface) LabelingInterface {
	return &labeling{
		service: service,
	}
}

func (l *labeling) Create(c *fiber.Ctx) error {
	var request domain.MailLabel
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parser Body Request"))
	}

	data, err := l.service.Create(request)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(201).JSON(web.Created("Success Create New Label Mail", data))
}

func (l *labeling) FindByID(c *fiber.Ctx) error {
	getID := c.Params("id")
	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing ID"))
	}

	data, err := l.service.FindByID(uint(id))
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(200).JSON(web.Success("Success Get Mail with This ID", data))
}

func (l *labeling) FindByCategory(c *fiber.Ctx) error {
	getCategory := c.Params("category")
	getUID := c.Params("user_id")

	userID, err := strconv.Atoi(getUID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing User ID"))
	}

	categoryID, err := strconv.Atoi(getCategory)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Category ID"))
	}

	data, err := l.service.FindByCategory(uint(categoryID), uint(userID))
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(200).JSON(web.Success("Success Get Mail with This Category", data))
}

func (l *labeling) UpdateLabel(c *fiber.Ctx) error {
	var request domain.MailLabel
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parser Body Request"))
	}

	getID := c.Params("id")
	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Request ID"))
	}

	getCategoryID := c.Params("categoryID")
	categoryID, err := strconv.Atoi(getCategoryID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Request ID"))
	}

	data, err := l.service.UpdateLabel(uint(id), uint(categoryID))
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(200).JSON(web.Success("Success Update Label Mail", data))
}
