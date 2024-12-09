package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/service"
	"github.com/sujit-baniya/flash"
)

type LabelingView interface {
	Add(c *fiber.Ctx) error
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
		}).Next()
	}

	id_email := c.Params("id")
	if id_category == "" {
		return flash.WithError(c, fiber.Map{
			"type":    "error",
			"message": "Sorry ID mail nil",
		}).Next()
	}

	id, _ := strconv.Atoi(id_category)
	idmail, _ := strconv.Atoi(id_email)

	token := c.Locals("Authorization").(string)
	get_id, err := l.middleware.GetUserID(token)
	if get_id == 0 && err != nil {
		return flash.WithError(c, fiber.Map{
			"type":    "error",
			"message": "Sorry User ID invalid",
		}).Next()
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
		}).Next()
	}

	return flash.WithSuccess(c, fiber.Map{
		"type":    "success",
		"message": "Success add mail to label",
	}).Next()
}

// func (l *labeling) SpamView(c *fiber.Ctx) error {

// }
