package handler

import (
	"fmt"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/service"
	"github.com/rulanugrh/megaclite/view"
	"github.com/rulanugrh/megaclite/view/auth"
	"github.com/sujit-baniya/flash"
)

type UserView interface {
	RegisterView(c *fiber.Ctx) error
	LoginView(c *fiber.Ctx) error
	ProfileView(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type user struct {
	service    service.UserInterface
	middleware middleware.JWTInterface
	conf       *config.App
}

func NewUserView(service service.UserInterface, conf *config.App) UserView {
	return &user{
		service:    service,
		middleware: middleware.NewJWTMiddleware(),
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
			return web.RedirectView(c, fmt.Sprintf("something wrong: %s", err.Error()), "/register")
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
			return web.RedirectView(c, err.Error(), "/")
		}

		token, err := u.middleware.GenerateToken(*data)
		if err != nil {
			return web.RedirectView(c, err.Error(), "/")
		}

		session, err := u.conf.Store.Get(c)
		if err != nil {
			return web.RedirectView(c, "cannot create new session", "/")
		}

		session.Set("Authorization", token)
		err = session.Save()

		if err != nil {
			return web.RedirectView(c, "cannot save new session", "/")
		}

		return web.SuccessView(c, "Success Login Account", data, "/home")
	}

	return handler(c)
}

func (u *user) ProfileView(c *fiber.Ctx) error {
	token, checks := c.Locals("Authorization").(string)
	if !checks {
		log.Println(token)
	}
	log.Println(token)
	getEmail, err := u.middleware.GetEmail(token)
	if err != nil {
		log.Println(err)
		return web.RedirectView(c, err.Error(), "/home")
	}

	var check bool = getEmail != ""
	if !check {
		return web.RedirectView(c, "Sorry your token is invalid", "/home")
	}

	data, err := u.service.GetEmail(getEmail)
	if err != nil {
		log.Println(err)
		return web.RedirectView(c, err.Error(), "/home")
	}

	index := view.ProfileIndex(*data)
	views := view.ProfileView("Login View", false, flash.Get(c), check, index)

	handler := adaptor.HTTPHandler(templ.Handler(views))
	return handler(c)
}

func (u *user) Logout(c *fiber.Ctx) error {
	session, err := u.conf.Store.Get(c)
	if err != nil {
		return web.RedirectView(c, "Session Not found", "/login")
	}

	err = session.Destroy()
	if err != nil {
		return web.RedirectView(c, fmt.Sprintf("something wrong: %s", err.Error()), "/home")
	}

	return web.SuccessView(c, "Success Logout", nil, "/login")
}
