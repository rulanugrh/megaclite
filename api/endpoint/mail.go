package endpoint

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/service"
)

type MailInterface interface {
	Create(c *fiber.Ctx) error
	Sent(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Inbox(c *fiber.Ctx) error
}

type mail struct {
	service    service.MailInterface
	middleware middleware.JWTInterface
}

func NewMailAPI(service service.MailInterface) MailInterface {
	return &mail{
		service:    service,
		middleware: middleware.NewJWTMiddleware(),
	}
}

// @Summary create new mail
// @ID createMail
// @Tags mails
// @Accept json
// @Produce json
// @Param request body domain.MailRegister true "request body for create mail"
// @Router /api/mail/create [post]
// @Success 200 {object} web.Response
// @Failure 500 {object} web.Response
// @Failure 400 {object} web.Response
func (m *mail) Create(c *fiber.Ctx) error {
	getToken := c.Get("Authorization")
	email, err := m.middleware.GetEmail(getToken)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// parser body request to entity
	var request domain.MailRegister
	request.From = email
	err = c.BodyParser(&request)
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

// @Summary get all sent email
// @ID getSendMail
// @Tags mails
// @Accept json
// @Produce json
// @Router /api/mail/send [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (m *mail) Sent(c *fiber.Ctx) error {
	// get email from parameter
	getToken := c.Get("Authorization")
	email, err := m.middleware.GetEmail(getToken)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// process get data from email parameter
	data, err := m.service.Sent(email)
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
// @Router /api/mail/find/{id} [get]
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
// @Router /api/mail/delete/{id} [delete]
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

// @Summary get all inbox email
// @ID getAll
// @Tags mails
// @Accept json
// @Produce json
// @Router /api/mail/inbox [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (m *mail) Inbox(c *fiber.Ctx) error {
	// get email from parameter
	getToken := c.Get("Authorization")
	email, err := m.middleware.GetEmail(getToken)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// process get data from email parameter
	data, err := m.service.Inbox(email)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success
	return c.Status(200).JSON(web.Success("Success get all mails by this email", data))
}
