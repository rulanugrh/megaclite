package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/service"
)

type UserInterface interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
}

type user struct {
	service service.UserInterface
}

func NewUserHandler(service service.UserInterface) UserInterface {
	return &user{
		service: service,
	}
}

// @Summary register new account
// @ID register
// @Tags users
// @Accept json
// @Produce json
// @Param request body domain.Register true "request body for new user"
// @Router /api/user/register [post]
// @Success 201 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (u *user) Register(c *fiber.Ctx) error {
	var request domain.Register
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot parsing body request"))
	}

	data, err := u.service.Register(request)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(201).JSON(web.Created("Success create new account", data))
}

// @Summary login user
// @ID login
// @Tags users
// @Accept json
// @Produce json
// @Param request body domain.Login true "request body for login existing account"
// @Route /api/user/login [post]
// @Succes 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (u *user) Login(c *fiber.Ctx) error {
	var request domain.Login
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot parsing body request"))
	}

	data, err := u.service.Login(request)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(201).JSON(web.Success("Success login into account", data))
}

// @Summary get user by emails
// @ID get_by_emails
// @Tags users
// @Accept json
// @Produce json
// @Param emails path string true "Emails User"
// @Route /api/user/{emails} [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (u *user) Get(c *fiber.Ctx) error {
	emails := c.Params("emails")

	data, err := u.service.GetEmail(emails)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(201).JSON(web.Success("Success get account", data))
}
