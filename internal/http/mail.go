package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/service"
)

type MailInterface interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type mail struct {
	service service.MailInterface
}

func NewMailHandler(service service.MailInterface) MailInterface {
	return &mail{
		service: service,
	}
}

func (m *mail) Create(c *fiber.Ctx) error {
	var request domain.Mail
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Error while parsing request"))
	}

	data, err := m.service.Create(request)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(201).JSON(web.Created("Success create new Mail", data))
}

func (m *mail) GetAll(c *fiber.Ctx) error {
	email := c.Params("email")

	data, err := m.service.Get(email)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(200).JSON(web.Success("Success get all mails by this email", data))
}

func (m *mail) GetByID(c *fiber.Ctx) error {
	getID := c.Params("id")

	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON("Error while parsing ID")
	}

	data, err := m.service.FindByID(uint(id))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(web.Success("Success get mails by this id", data))
}

func (m *mail) Delete(c *fiber.Ctx) error {
	getID := c.Params("id")

	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON("Error while parsing ID")
	}

	err = m.service.Delete(uint(id))
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(200).JSON(web.Success("Success delete this email", nil))
}
