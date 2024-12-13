package handler

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/service"
	"github.com/rulanugrh/megaclite/view"
	"github.com/sujit-baniya/flash"
)

type MailView interface {
	InboxView(c *fiber.Ctx) error
	SentView(c *fiber.Ctx) error
	AddMail(c *fiber.Ctx) error
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

	data, err := m.service.Inbox(getMail)
	if err != nil {
		msgError["message"] = "Cannot get Inbox Mail"
		return flash.WithError(c, msgError).Redirect("/home")
	}

	index := view.HomeIndex(*data, getMail)
	views := view.Home("Inbox", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "POST" {

	}
	return handler(c)

}

func (m *mail) SentView(c *fiber.Ctx) error {
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

	data, err := m.service.Sent(getMail)
	if err != nil {
		msgError["message"] = "Cannot get Sent Mail"
		return flash.WithError(c, msgError).Redirect("/home/sent")
	}

	index := view.MailViewIndex(*data)
	views := view.MailView("| Sent Mail", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	return handler(c)
}

func (m *mail) AddMail(c *fiber.Ctx) error {
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

	if c.Method() == "POST" {

		form, err := c.MultipartForm()
		if err != nil {
			msgError["message"] = "Invalid handler multi part form"
			return flash.WithError(c, msgError).Redirect("/")
		}

		files := form.File["attachments"]

		var filenames []string
		for _, file := range files {
			pwd, _ := os.Getwd()
			err := c.SaveFile(file, fmt.Sprintf("%s\\view\\public\\%d-%s", pwd, time.Now().Unix(), file.Filename))
			if err != nil {
				msgError["message"] = "Error while save attachment"
				log.Println(err.Error())
				return flash.WithError(c, msgError).Next()
			}

			filenames = append(filenames, fmt.Sprintf("%s\\view\\public\\%d-%s", pwd, time.Now().Unix(), file.Filename))
		}

		request := domain.MailRegister{
			From:       getMail,
			To:         c.FormValue("to-people"),
			Message:    c.FormValue("message"),
			Title:      c.FormValue("title"),
			Subtitle:   c.FormValue("subtitle"),
			Attachment: strings.Join(filenames, ","),
		}

		data, err := m.service.Create(request)
		if err != nil {
			log.Println(err.Error())
			return flash.WithError(c, fiber.Map{
				"type":    "error",
				"message": err.Error(),
			}).Next()
		}

		success := fiber.Map{
			"type":    "success",
			"message": "success create mail " + data.From,
		}
		return flash.WithSuccess(c, success).Redirect("/home")
	}

	return flash.WithSuccess(c, fiber.Map{"type": "success"}).Next()

}
