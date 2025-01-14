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
	Run(user handler.UserView, mail handler.MailView, label handler.LabelingView)
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

func (v *viewRoute) Run(user handler.UserView, mail handler.MailView, label handler.LabelingView) {
	v.app.Static("/assets", "./view/assets")
	v.app.Get("/", v.commonMiddleware.ViewMiddleware, user.HomeView)
	// Views User
	v.app.Get("/login", user.LoginView)
	v.app.Post("/login", user.LoginView)
	v.app.Get("/register", user.RegisterView)
	v.app.Post("/register", user.RegisterView)
	v.app.Get("/home/profile", v.commonMiddleware.ViewMiddleware, user.ProfileView)
	v.app.Post("/logout", v.commonMiddleware.ViewMiddleware, user.Logout)
	v.app.Post("/update/password", v.commonMiddleware.ViewMiddleware, user.UpdatePassword)
	v.app.Post("/update/profile", v.commonMiddleware.ViewMiddleware, user.UpdateProfile)

	// Home Index
	v.app.Get("/home", v.commonMiddleware.ViewMiddleware, mail.InboxView)
	v.app.Post("/mail/add", v.commonMiddleware.ViewMiddleware, mail.AddMail)
	v.app.Get("/home/sent", v.commonMiddleware.ViewMiddleware, mail.SentView)
	v.app.Get("/home/trash", v.commonMiddleware.ViewMiddleware, label.TrashView)
	v.app.Get("/home/favorite", v.commonMiddleware.ViewMiddleware, label.FavoriteView)
	v.app.Get("/home/spam", v.commonMiddleware.ViewMiddleware, label.SpamView)
	v.app.Get("/home/detail/:id", v.commonMiddleware.ViewMiddleware, mail.GetMail)

	// Label Index
	v.app.Post("/label/add/:categoryID/:id", v.commonMiddleware.ViewMiddleware, label.Add)
	// Running Application
	err := v.app.Listen(fmt.Sprintf("%s:%s", v.conf.Server.Host, v.conf.Server.ViewPort))
	if err != nil {
		log.Fatal("Error while running server: " + err.Error())
	}

	log.Println("App running at: " + v.conf.Server.Host + ":" + v.conf.Server.ViewPort)
}
