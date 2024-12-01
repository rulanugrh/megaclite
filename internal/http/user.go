package handler

import (
	"log"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/service"
	"github.com/rulanugrh/megaclite/view/auth"
	"github.com/sujit-baniya/flash"
)

type UserView interface {
	RegisterView(c *fiber.Ctx) error
	LoginView(c *fiber.Ctx) error
	HomeView(c *fiber.Ctx) error
}

type user struct {
	service service.UserInterface
}

func NewUserView(service service.UserInterface) UserView {
	return &user{
		service: service,
	}
}

func (u *user) RegisterView(c *fiber.Ctx) error {
	index := auth.RegisterIndex()
	views := auth.Register("Register Account", false, flash.Get(c), index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "POST" {
		request := domain.Register{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
			Email:    c.FormValue("email"),
		}

		data, err := u.service.Register(request)
		if err != nil {
			log.Println(err)
			fm := fiber.Map{
				"type":    "error",
				"message": err.Error(),
			}

			return flash.WithError(c, fm).Redirect("/")
		}

		// Mengirim file JSON
		c.Set("Content-Disposition", "attachment; filename=keygen.pgp")
		c.Set("Content-Type", "text/plain")
		c.Send([]byte(data.Private))

		fm := fiber.Map{
			"type":    "success",
			"message": "Success Create Account",
		}

		return flash.WithSuccess(c, fm).Redirect("/login")
	}

	return handler(c)
}

func (u *user) LoginView(c *fiber.Ctx) error {
	index := auth.RegisterIndex()
	views := auth.Register("Register Account", false, flash.Get(c), index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "POST" {
		request := domain.Register{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
			Email:    c.FormValue("email"),
		}

		data, err := u.service.Register(request)
		if err != nil {
			fm := fiber.Map{
				"type":    "error",
				"message": err.Error(),
			}

			return flash.WithError(c, fm).Redirect("/register")
		}

		fm := fiber.Map{
			"type":    "success",
			"message": "Success Create Account",
			"data":    data.Private,
		}

		return flash.WithSuccess(c, fm).Redirect("/login")
	}

	return handler(c)
}
func (u *user) HomeView(c *fiber.Ctx) error {
	index := auth.RegisterIndex()
	views := auth.Register("Register Account", false, flash.Get(c), index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "POST" {
		request := domain.Register{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
			Email:    c.FormValue("email"),
		}

		data, err := u.service.Register(request)
		if err != nil {
			log.Println(err.Error())
			fm := fiber.Map{
				"type":    "error",
				"message": err.Error(),
			}

			return flash.WithError(c, fm).Redirect("/")
		}

		// Mengirim file JSON
		go func() {
			c.Set("Content-Disposition", "attachment; filename=keygen.pgp")
			c.Set("Content-Type", "text/plain")
			c.Send([]byte(data.Private))
		}()

		// Return redirect if succces
		fm := fiber.Map{
			"type":    "success",
			"message": "Success Creat Account",
		}

		return flash.WithSuccess(c, fm).Redirect("/login")

	}

	return handler(c)
}
