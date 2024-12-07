package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/api/endpoint"
	"github.com/rulanugrh/megaclite/config"
)

type APIInterface interface {
	Run(user endpoint.UserInterface, labeling endpoint.LabelingInterface, mail endpoint.MailInterface, category endpoint.CategoryInterface)
}

type Api struct {
	app  *fiber.App
	conf *config.App
}

func NewAPIRoutes(conf *config.App) APIInterface {
	return &Api{
		app: fiber.New(fiber.Config{
			AppName: "API Webmail Megaclite with Golang and Mysql",
		}),
		conf: conf,
	}
}

func (a *Api) Run(user endpoint.UserInterface, labeling endpoint.LabelingInterface, mail endpoint.MailInterface, category endpoint.CategoryInterface) {
	// Routing API for Mail
	mailRoutes := a.app.Group("/api/mail")
	mailRoutes.Post("/", mail.Create)
	mailRoutes.Get("/find/:id", mail.GetByID)
	mailRoutes.Get("/sent", mail.Sent)
	mailRoutes.Get("/inbox", mail.Inbox)
	mailRoutes.Delete("/delete/:id", mail.GetByID)

	// Routing API for User
	userRoutes := a.app.Group("/api/user")
	userRoutes.Post("/register", user.Register)
	userRoutes.Post("/login", user.Login)
	userRoutes.Get("/:email", user.Get)

	// Routing API For Category
	categoryRoutes := a.app.Group("/api/category")
	categoryRoutes.Post("/", category.Create)
	categoryRoutes.Delete("/:id", category.Delete)
	categoryRoutes.Put("/:id", category.Update)

	// Routing API For Labeling
	labelingRoutes := a.app.Group("/api/labeling")
	labelingRoutes.Post("/", labeling.Create)
	labelingRoutes.Get("/:categoryID/:user_id", labeling.FindByCategory)
	labelingRoutes.Put("/update/:id/:categoryID", labeling.UpdateLabel)

	// Running Application
	err := a.app.Listen(fmt.Sprintf("%s:%s", a.conf.Server.Host, a.conf.Server.ApiPort))
	if err != nil {
		log.Fatal("Error while running server: " + err.Error())
	}

	log.Println("App running at: " + a.conf.Server.Host + ":" + a.conf.Server.ApiPort)
}
