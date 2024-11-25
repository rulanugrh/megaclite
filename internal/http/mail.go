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

// @Summary create new mail
// @ID createMail
// @Tags mails
// @Accept json
// @Produce json
// @Param request body domain.MailRegister true "request body for create mail"
// @Route /api/mail/create [post]
// @Success 200 {object} web.Response
// @Failure 500 {object} web.Response
// @Failure 400 {object} web.Response
func (m *mail) Create(c *fiber.Ctx) error {
	// parser body request to entity
	var request domain.MailRegister
	err := c.BodyParser(&request)
	// check error if have error while create data
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Error while parsing request"))
	}

	// process create to service layer
	data, err := m.service.Create(request)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success
	return c.Status(201).JSON(web.Created("Success create new Mail", data))
}

// @Summary get all emails
// @ID getAll
// @Tags mails
// @Accept json
// @Produce json
// @Param user path string true "user parameter"
// @Route /api/mail/get/{user} [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (m *mail) GetAll(c *fiber.Ctx) error {
	// get email from parameter
	user := c.Params("user")

	// process get data from email parameter
	data, err := m.service.Get(user)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success
	return c.Status(200).JSON(web.Success("Success get all mails by this email", data))
}

// @Summary "get mails by id"
// @ID getByID
// @Tags mails
// @Accept json
// @Produce json
// @Param id path int true "id mail"
// @Route /api/mail/find/{id} [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (m *mail) GetByID(c *fiber.Ctx) error {
	// get id from parameter
	getID := c.Params("id")

	// convert string into integer
	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON("Error while parsing ID")
	}

	// process get data from id parameter
	data, err := m.service.FindByID(uint(id))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// return success
	return c.Status(200).JSON(web.Success("Success get mails by this id", data))
}

// @Summary "delete mail by id"
// @ID deleteMailByID
// @Tags mails
// @Accept json
// @Produce json
// @Param id path int true "id email"
// @Route /api/mail/delete/{id} [delete]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
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
