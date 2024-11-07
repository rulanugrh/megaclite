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

func (u *user) Get(c *fiber.Ctx) error {
	emails := c.Params("emails")

	data, err := u.service.GetEmail(emails)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	return c.Status(201).JSON(web.Success("Success get account", data))
}
