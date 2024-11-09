package main

import (
	"fmt"
	"log"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/config"
	handler "github.com/rulanugrh/megaclite/internal/http"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
	"github.com/rulanugrh/megaclite/internal/service"
)

func main() {
	// Initialize PGP
	pgp := crypto.PGP()

	//  Initialize Config and Connection
	conf := config.GetConfig()
	connectionDB := config.InitDatabase(conf)
	connectionDB.Connection()

	// Initialize Middleware
	middleware := middleware.NewPGPMiddleware(pgp)

	// Initialize User Komponen
	userRepository := repository.NewUserRepository(*connectionDB)
	userService := service.NewUserService(userRepository, middleware)
	userHandler := handler.NewUserHandler(userService)

	// Initialize Mail Komponen
	mailRepository := repository.NewMailRepository(*connectionDB)
	mailService := service.NewMailService(mailRepository)
	mailHandler := handler.NewMailHandler(mailService)

	// Initialize Category Komponen
	categoryRepository := repository.NewCategoryRepository(*connectionDB)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Initialize Labeling Komponen
	labelingRepository := repository.NewLabelMailRepository(*connectionDB)
	labelingService := service.NewLabelMailService(labelingRepository)
	labelingHandler := handler.NewLabelMailHandler(labelingService)

	app := fiber.New(fiber.Config{
		AppName: "PGP with Golang",
	})

	err := application(mailHandler, userHandler, categoryHandler, labelingHandler, app, conf)
	if err != nil {
		log.Fatal("Error While Connection to App: " + err.Error())
	}
}

func application(mail handler.MailInterface, user handler.UserInterface, category handler.CategoryInterface, labeling handler.LabelingInterface, app *fiber.App, config *config.App) error {
	// Route Group for Mail Handler
	mailRoutes := app.Group("/api/mail")
	mailRoutes.Post("/", mail.Create)
	mailRoutes.Get("/:id", mail.GetByID)
	mailRoutes.Get("/getall", mail.GetAll)
	mailRoutes.Delete("/delete/:id", mail.GetByID)

	// Routing Handler for User
	userRoutes := app.Group("/api/user")
	userRoutes.Post("/register", user.Register)
	userRoutes.Post("/login", user.Login)
	userRoutes.Get("/:email", user.Get)

	// Routing Handler For Category
	categoryRoutes := app.Group("/api/category")
	categoryRoutes.Post("/", category.Create)
	categoryRoutes.Delete("/:id", category.Delete)
	categoryRoutes.Put("/:id", category.Update)

	// Routing Handler For Labeling
	labelingRoutes := app.Group("/api/labeling")
	labelingRoutes.Post("/", labeling.Create)
	labelingRoutes.Get("/:categoryID/:user_id", labeling.FindByCategory)
	labelingRoutes.Get("/get/:id", labeling.FindByID)
	labelingRoutes.Put("/update/:id/:categoryID", labeling.UpdateLabel)

	// Running Application
	err := app.Listen(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port))
	if err != nil {
		return err
	}

	return nil
}
