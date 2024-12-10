package handler

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/service"
	"github.com/rulanugrh/megaclite/view"
	"github.com/sujit-baniya/flash"
)

type LabelingView interface {
	Add(c *fiber.Ctx) error
	SpamView(c *fiber.Ctx) error
	TrashView(c *fiber.Ctx) error
	FavoriteView(c *fiber.Ctx) error
}

type labeling struct {
	service    service.LabelingInterface
	middleware middleware.JWTInterface
}

func NewLabelView(service service.LabelingInterface) LabelingView {
	return &labeling{
		service:    service,
		middleware: middleware.NewJWTMiddleware(),
	}
}

func (l *labeling) Add(c *fiber.Ctx) error {
	id_category := c.Params("categoryID")
	if id_category == "" {
		return flash.WithError(c, fiber.Map{
			"type":    "error",
			"message": "Sorry ID category nil",
		}).Redirect("/home")
	}

	id_email := c.Params("id")
	if id_category == "" {
		return flash.WithError(c, fiber.Map{
			"type":    "error",
			"message": "Sorry ID mail nil",
		}).Redirect("/home")
	}

	id, _ := strconv.Atoi(id_category)
	idmail, _ := strconv.Atoi(id_email)

	token := c.Locals("Authorization").(string)
	get_id, err := l.middleware.GetUserID(token)
	if get_id == 0 && err != nil {
		return flash.WithError(c, fiber.Map{
			"type":    "error",
			"message": "Sorry User ID invalid",
		}).Redirect("/home")
	}

	request := domain.MailLabelRegister{
		CategoryID: uint(id),
		MailID:     uint(idmail),
		UserID:     get_id,
	}

	_, err = l.service.Create(request)
	if err != nil {
		return flash.WithError(c, fiber.Map{
			"type":    "error",
			"message": "Sorry cannot add mail label",
		}).Redirect("/home")
	}

	return flash.WithSuccess(c, fiber.Map{
		"type":    "success",
		"message": "Success add mail to label",
	}).Redirect("/home")
}

func (l *labeling) SpamView(c *fiber.Ctx) error {
	msgError := fiber.Map{
		"type": "error",
	}

	token := c.Locals("Authorization").(string)
	getUID, err := l.middleware.GetUserID(token)
	if err != nil {
		msgError["message"] = err.Error()
		return flash.WithError(c, msgError).Redirect("/home")
	}

	var check bool = getUID != 0
	if !check {
		msgError["message"] = "Sorry you token is invalid"
		return flash.WithError(c, msgError).Redirect("/home")
	}

	data, err := l.service.FindByCategory(4, getUID)
	if err != nil {
		msgError["message"] = "Cannot get Sent Mail"
		return flash.WithError(c, msgError).Redirect("/home")
	}

	index := view.MailViewIndex(*data)
	views := view.MailView("| Spam Mail", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	return handler(c)
}

func (l *labeling) TrashView(c *fiber.Ctx) error {
	msgError := fiber.Map{
		"type": "error",
	}

	token := c.Locals("Authorization").(string)
	getUID, err := l.middleware.GetUserID(token)
	if err != nil {
		msgError["message"] = err.Error()
		return flash.WithError(c, msgError).Redirect("/home")
	}

	var check bool = getUID != 0
	if !check {
		msgError["message"] = "Sorry you token is invalid"
		return flash.WithError(c, msgError).Redirect("/home")
	}

	data, err := l.service.FindByCategory(3, getUID)
	if err != nil {
		msgError["message"] = "Cannot get Sent Mail"
		return flash.WithError(c, msgError).Redirect("/home")
	}

	index := view.MailViewIndex(*data)
	views := view.MailView("| Spam Mail", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	return handler(c)
}

func (l *labeling) FavoriteView(c *fiber.Ctx) error {
	msgError := fiber.Map{
		"type": "error",
	}

	token := c.Locals("Authorization").(string)
	getUID, err := l.middleware.GetUserID(token)
	if err != nil {
		msgError["message"] = err.Error()
		return flash.WithError(c, msgError).Redirect("/home")
	}

	var check bool = getUID != 0
	if !check {
		msgError["message"] = "Sorry you token is invalid"
		return flash.WithError(c, msgError).Redirect("/home")
	}

	data, err := l.service.FindByCategory(1, getUID)
	if err != nil {
		msgError["message"] = "Cannot get Sent Mail"
		return flash.WithError(c, msgError).Redirect("/home")
	}

	index := view.MailViewIndex(*data)
	views := view.MailView("| Spam Mail", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	return handler(c)
}
