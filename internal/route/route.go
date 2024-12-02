package route

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/config"
	handler "github.com/rulanugrh/megaclite/internal/http"
	"github.com/rulanugrh/megaclite/internal/middleware"
)

type RouteViewInterface interface {
	Run(user handler.UserView)
}

type viewRoute struct {
	conf             *config.App
	app              *fiber.App
	commonMiddleware middleware.CommonMiddlewareInterface
}

func NewViewRoute(conf *config.App, middleware middleware.CommonMiddlewareInterface) RouteViewInterface {
	return &viewRoute{
		app: fiber.New(fiber.Config{
			AppName: "Webmail Megaclite",
		}),
		conf:             conf,
		commonMiddleware: middleware,
	}
}

func (v *viewRoute) Run(user handler.UserView) {
	v.app.Static("/assets", "./view/assets")
	// Views User
	v.app.Get("/", user.LoginView)
	v.app.Post("/", user.LoginView)
	v.app.Get("/register", user.RegisterView)
	v.app.Post("/register", user.RegisterView)

	// Home Index
	v.app.Get("/home", v.commonMiddleware.ViewMiddleware, user.HomeView)
	// Running Application
	err := v.app.Listen(fmt.Sprintf("%s:%s", v.conf.Server.Host, v.conf.Server.ViewPort))
	if err != nil {
		log.Fatal("Error while running server: " + err.Error())
	}

	log.Println("App running at: " + v.conf.Server.Host + ":" + v.conf.Server.ViewPort)
}