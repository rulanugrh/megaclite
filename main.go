package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rulanugrh/megaclite/api"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	handler "github.com/rulanugrh/megaclite/internal/http"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
	"github.com/rulanugrh/megaclite/internal/service"
	"gorm.io/gorm"
)

var Store *session.Store

// @title Megaclite API Documentation
// @version 1.0
// @description Documentation for API OpenPGP with HTMX
// @termsOfService https://swagger.io/terms

// @contact.name Kyora
// @contact.url https://github.com/rulanugrh
// @contact.email rulanugrh@proton.me

// @license.name MIT
// @host localhost:4000
// @BasePath /api/
// @securityDefinition.basic BasicAuth
func main() {
	// Initialize PGP
	pgp := crypto.PGPWithProfile(profile.RFC4880())

	//  Initialize Config and Connection
	conf := config.GetConfig()
	connectionDB := config.InitDatabase(conf)
	db := connectionDB.Connection()

	// Initialize Middleware
	middleware := middleware.NewPGPMiddleware(pgp, *conf)

	// Initialize User Komponen
	userRepository := repository.NewUserRepository(*connectionDB)
	userService := service.NewUserService(userRepository, middleware)
	userAPI := api.NewUserAPI(userService)
	userView := handler.NewUserView(userService)

	// Initialize Mail Komponen
	mailRepository := repository.NewMailRepository(*connectionDB)
	mailService := service.NewMailService(mailRepository, middleware)
	mailAPI := api.NewMailAPI(mailService)

	// Initialize Category Komponen
	categoryRepository := repository.NewCategoryRepository(*connectionDB)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryAPI := api.NewCategoryAPI(categoryService)

	// Initialize Labeling Komponen
	labelingRepository := repository.NewLabelMailRepository(*connectionDB)
	labelingService := service.NewLabelMailService(labelingRepository)
	labelingAPI := api.NewLabelMailAPI(labelingService)

	app := fiber.New(fiber.Config{
		AppName: "PGP with Golang",
	})

	// parsing argument command
	args := os.Args[1]
	if args == "migration" {
		err := connectionDB.DB.AutoMigrate(&domain.Category{}, &domain.User{}, &domain.Mail{}, &domain.MailLabel{})
		if err != nil {
			log.Fatal("Error while migration data: " + err.Error())
		}
	} else if args == "seed" {
		seeder(db)
	} else if args == "api" {
		webAPI(mailAPI, userAPI, categoryAPI, labelingAPI, app, conf)
	} else if args == "serve" {
		webView(userView, app, conf)
	} else {
		help_command()
	}

}

func webAPI(mail api.MailInterface, user api.UserInterface, category api.CategoryInterface, labeling api.LabelingInterface, app *fiber.App, config *config.App) {
	// Route Group for Mail API
	mailRoutes := app.Group("/api/mail")
	mailRoutes.Post("/", mail.Create)
	mailRoutes.Get("/find/:id", mail.GetByID)
	mailRoutes.Get("/getall", mail.GetAll)
	mailRoutes.Delete("/delete/:id", mail.GetByID)

	// Routing API for User
	userRoutes := app.Group("/api/user")
	userRoutes.Post("/register", user.Register)
	userRoutes.Post("/login", user.Login)
	userRoutes.Get("/:email", user.Get)

	// Routing API For Category
	categoryRoutes := app.Group("/api/category")
	categoryRoutes.Post("/", category.Create)
	categoryRoutes.Delete("/:id", category.Delete)
	categoryRoutes.Put("/:id", category.Update)

	// Routing API For Labeling
	labelingRoutes := app.Group("/api/labeling")
	labelingRoutes.Post("/", labeling.Create)
	labelingRoutes.Get("/:categoryID/:user_id", labeling.FindByCategory)
	labelingRoutes.Get("/get/:id", labeling.FindByID)
	labelingRoutes.Put("/update/:id/:categoryID", labeling.UpdateLabel)

	// Running Application
	err := app.Listen(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.ApiPort))
	if err != nil {
		log.Fatal("Error while running server: " + err.Error())
	}

	log.Println("App running at: " + config.Server.Host + ":" + config.Server.ApiPort)
}

func webView(user handler.UserView, app *fiber.App, conf *config.App) {
	config.Store = session.New(
		session.Config{
			CookieHTTPOnly: true,
			Expiration:     24 * time.Hour,
		},
	)
	app.Static("/assets", "./view/assets")
	// Views User
	app.Get("/", user.LoginView)
	app.Post("/", user.LoginView)
	app.Get("/register", user.RegisterView)
	app.Post("/register", user.RegisterView)

	// Home Index
	app.Get("/home", middleware.ViewMiddleware, user.HomeView)
	// Running Application
	err := app.Listen(fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.ViewPort))
	if err != nil {
		log.Fatal("Error while running server: " + err.Error())
	}

	log.Println("App running at: " + conf.Server.Host + ":" + conf.Server.ViewPort)

}

func seeder(db *gorm.DB) {
	favoriteCategory := domain.Category{
		Name:        "Favorite",
		Description: "Favorite Mail",
	}

	archiveCategory := domain.Category{
		Name:        "Archive",
		Description: "Archive Mail",
	}

	trashCategory := domain.Category{
		Name:        "Trash",
		Description: "Trash Mail",
	}

	spamCategory := domain.Category{
		Name:        "Spam",
		Description: "Spam Mail",
	}

	err := db.Create(&favoriteCategory).Error
	if err != nil {
		log.Fatal("Something error while migrate: " + err.Error())
	}

	err = db.Create(&archiveCategory).Error
	if err != nil {
		log.Fatal("Something error while migrate: " + err.Error())
	}

	err = db.Create(&trashCategory).Error
	if err != nil {
		log.Fatal("Something error while migrate: " + err.Error())
	}

	err = db.Create(&spamCategory).Error
	if err != nil {
		log.Fatal("Something error while migrate: " + err.Error())
	}

	log.Println("Success seeding data into database")
}

func help_command() {
	content := [][]string{
		{"help", "to show help message"},
		{"migration", "to migration table into database"},
		{"seed", "to seed dummy data into database"},
		{"api", "to serve API"},
		{"serve", "to serve application"},
	}

	example := "\nexample: go run main.go help\n"
	max := len(content[0][0])
	for _, part := range content {
		length := len(part[0])
		if length > max {
			max = length
		}
	}

	var builder strings.Builder
	const space = 4
	for _, part := range content {
		builder.WriteString(part[0])
		spacer := (max - len(part[0])) + space
		for spacer > 0 {
			builder.WriteByte(' ')
			spacer--
		}

		builder.WriteString(part[1])
		builder.WriteByte('\n')
	}

	println(builder.String()[:builder.Len()-1])
	println(example)
}
