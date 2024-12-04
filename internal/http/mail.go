package handler

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/jinzhu/copier"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/service"
	"github.com/rulanugrh/megaclite/view"
	"github.com/sujit-baniya/flash"
)

type MailView interface {
	InboxView(c *fiber.Ctx) error
	ArchiveView(c *fiber.Ctx) error
	SpamView(c *fiber.Ctx) error
	SentView(c *fiber.Ctx) error
	TrashView(c *fiber.Ctx) error
}

type mail struct {
	service    service.MailInterface
	middleware middleware.JWTInterface
}

func NewMailView(service service.MailInterface) MailView {
	return &mail{
		service:    service,
		middleware: middleware.NewJWTMiddleware(),
	}
}

func (m *mail) InboxView(c *fiber.Ctx) error {
	var response []web.GetMail
	msgError := fiber.Map{
		"type": "error",
	}

	token := c.Locals("Authorization").(string)
	getMail, err := m.middleware.GetEmail(token)
	if err != nil {
		msgError["message"] = err.Error()
		return flash.WithError(c, msgError).Redirect("/home")
	}

	var check bool = getMail != ""
	if !check {
		msgError["message"] = "Sorry you token is invalid"
		return flash.WithError(c, msgError).Redirect("/")
	}

	index := view.HomeIndex(response)
	views := view.Home("Inbox", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "GET" {
		data, err := m.service.Inbox(getMail)
		if err != nil {
			msgError["message"] = "Cannot get Inbox Mail"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		err = copier.Copy(&response, data)
		if err != nil {
			msgError["message"] = "Cannot parsing value"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		msgSuccess := fiber.Map{
			"type":    "success",
			"message": "Success Get Inbox Mail",
		}

		return flash.WithSuccess(c, msgSuccess).Next()
	}

	return handler(c)

}

func (m *mail) ArchiveView(c *fiber.Ctx) error {
	var response []web.GetMail
	msgError := fiber.Map{
		"type": "error",
	}

	token := c.Locals("Authorization").(string)
	getMail, err := m.middleware.GetEmail(token)
	if err != nil {
		msgError["message"] = err.Error()
		return flash.WithError(c, msgError).Redirect("/home")
	}

	var check bool = getMail != ""
	if !check {
		msgError["message"] = "Sorry you token is invalid"
		return flash.WithError(c, msgError).Redirect("/")
	}

	index := view.HomeIndex(response)
	views := view.Home("Inbox", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "GET" {
		data, err := m.service.Inbox(getMail)
		if err != nil {
			msgError["message"] = "Cannot get Inbox Mail"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		err = copier.Copy(&response, data)
		if err != nil {
			msgError["message"] = "Cannot parsing value"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		msgSuccess := fiber.Map{
			"type":    "success",
			"message": "Success Get Inbox Mail",
		}

		return flash.WithSuccess(c, msgSuccess).Next()
	}

	return handler(c)
}

func (m *mail) SpamView(c *fiber.Ctx) error {
	var response []web.GetMail
	msgError := fiber.Map{
		"type": "error",
	}

	token := c.Locals("Authorization").(string)
	getMail, err := m.middleware.GetEmail(token)
	if err != nil {
		msgError["message"] = err.Error()
		return flash.WithError(c, msgError).Redirect("/home")
	}

	var check bool = getMail != ""
	if !check {
		msgError["message"] = "Sorry you token is invalid"
		return flash.WithError(c, msgError).Redirect("/")
	}

	index := view.HomeIndex(response)
	views := view.Home("Inbox", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "GET" {
		data, err := m.service.Inbox(getMail)
		if err != nil {
			msgError["message"] = "Cannot get Inbox Mail"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		err = copier.Copy(&response, data)
		if err != nil {
			msgError["message"] = "Cannot parsing value"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		msgSuccess := fiber.Map{
			"type":    "success",
			"message": "Success Get Inbox Mail",
		}

		return flash.WithSuccess(c, msgSuccess).Next()
	}

	return handler(c)
}

func (m *mail) SentView(c *fiber.Ctx) error {
	var response []web.GetMail
	msgError := fiber.Map{
		"type": "error",
	}

	token := c.Locals("Authorization").(string)
	getMail, err := m.middleware.GetEmail(token)
	if err != nil {
		msgError["message"] = err.Error()
		return flash.WithError(c, msgError).Redirect("/home")
	}

	var check bool = getMail != ""
	if !check {
		msgError["message"] = "Sorry you token is invalid"
		return flash.WithError(c, msgError).Redirect("/")
	}

	index := view.HomeIndex(response)
	views := view.Home("Inbox", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "GET" {
		data, err := m.service.Inbox(getMail)
		if err != nil {
			msgError["message"] = "Cannot get Inbox Mail"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		err = copier.Copy(&response, data)
		if err != nil {
			msgError["message"] = "Cannot parsing value"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		msgSuccess := fiber.Map{
			"type":    "success",
			"message": "Success Get Inbox Mail",
		}

		return flash.WithSuccess(c, msgSuccess).Next()
	}

	return handler(c)
}

func (m *mail) TrashView(c *fiber.Ctx) error {
	var response []web.GetMail
	msgError := fiber.Map{
		"type": "error",
	}

	token := c.Locals("Authorization").(string)
	getMail, err := m.middleware.GetEmail(token)
	if err != nil {
		msgError["message"] = err.Error()
		return flash.WithError(c, msgError).Redirect("/home")
	}

	var check bool = getMail != ""
	if !check {
		msgError["message"] = "Sorry you token is invalid"
		return flash.WithError(c, msgError).Redirect("/")
	}

	index := view.HomeIndex(response)
	views := view.Home("Inbox", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "GET" {
		data, err := m.service.Inbox(getMail)
		if err != nil {
			msgError["message"] = "Cannot get Inbox Mail"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		err = copier.Copy(&response, data)
		if err != nil {
			msgError["message"] = "Cannot parsing value"
			return flash.WithError(c, msgError).Redirect("/home")
		}

		msgSuccess := fiber.Map{
			"type":    "success",
			"message": "Success Get Inbox Mail",
		}

		return flash.WithSuccess(c, msgSuccess).Next()
	}

	return handler(c)
}
