package handler

import (
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/service"
	"github.com/rulanugrh/megaclite/view"
	"github.com/rulanugrh/megaclite/view/auth"
	"github.com/sujit-baniya/flash"
)

type UserView interface {
	RegisterView(c *fiber.Ctx) error
	LoginView(c *fiber.Ctx) error
	HomeView(c *fiber.Ctx) error
}

type user struct {
	service    service.UserInterface
	middleware middleware.JWTInterface
	conf       *config.App
}

func NewUserView(service service.UserInterface, conf *config.App) UserView {
	return &user{
		service:    service,
		middleware: middleware.NewJWTToken(),
		conf:       conf,
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
			fm := fiber.Map{
				"type":    "error",
				"message": err.Error(),
			}

			return flash.WithError(c, fm).Redirect("/")
		}

		// Mengirim file JSON

		c.Set("Content-Disposition", "attachment; filename=keygen.pgp")
		c.Set("Content-Type", "text/plain")
		return c.Send([]byte(data.Private))
	}

	return handler(c)
}

func (u *user) LoginView(c *fiber.Ctx) error {

	index := auth.LoginIndex()
	views := auth.Login("Login View", false, flash.Get(c), index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	if c.Method() == "POST" {
		request := domain.Login{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}

		keygen, err := c.FormFile("file")
		if err != nil {
			log.Println(err.Error())
		}

		read, err := keygen.Open()
		if err != nil {
			log.Println("Cannot open file request")
		}

		content, err := io.ReadAll(read)
		if err != nil {
			log.Println("Cannot read content file")
		}

		data, err := u.service.Login(request, string(content))
		if err != nil {
			fm := fiber.Map{
				"type":    "error",
				"message": err.Error(),
			}

			return flash.WithError(c, fm).Redirect("/")
		}

		token, err := u.middleware.GenerateToken(*data)
		if err != nil {
			fm := fiber.Map{
				"type":    "error",
				"message": err.Error(),
			}

			return flash.WithError(c, fm).Redirect("/")
		}

		session, err := u.conf.Store.Get(c)
		if err != nil {
			fm := fiber.Map{
				"type":    "error",
				"message": "Cannot create new Session",
			}

			return flash.WithError(c, fm).Redirect("/")
		}

		session.Set("Authorization", token)
		err = session.Save()

		if err != nil {
			return flash.WithError(c, fiber.Map{
				"type":    "error",
				"message": "Cannot save session",
			}).Redirect("/")
		}

		fm := fiber.Map{
			"type":    "success",
			"message": "Success Login Account",
		}

		return flash.WithSuccess(c, fm).Redirect("/home")
	}

	return handler(c)
}
func (u *user) HomeView(c *fiber.Ctx) error {
	token := c.Locals("Authorization").(string)
	getMail, err := u.middleware.GetEmail(token)

	if err != nil {
		fm := fiber.Map{
			"type":    "error",
			"message": "Cannot get token jwt",
		}

		return flash.WithError(c, fm).Redirect("/")
	}

	index := view.HomeIndex(getMail)
	views := view.Home("Dashboard", false, flash.Get(c), index)

	handler := adaptor.HTTPHandler(templ.Handler(views))

	return handler(c)
}
