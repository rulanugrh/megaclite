package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	handler "github.com/rulanugrh/megaclite/internal/http"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
	"github.com/rulanugrh/megaclite/internal/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
	userHandler := handler.NewUserHandler(userService)

	// Initialize Mail Komponen
	mailRepository := repository.NewMailRepository(*connectionDB)
	mailService := service.NewMailService(mailRepository, middleware)
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

	// parsing argument command
	args := os.Args[1]
	if args == "migration" {
		err := connectionDB.DB.AutoMigrate(&domain.Category{}, &domain.User{}, &domain.Mail{}, &domain.MailLabel{})
		if err != nil {
			log.Fatal("Error while migration data: " + err.Error())
		}
	} else if args == "seed" {
		seeder(db, *conf)
	} else if args == "serve" {
		application(mailHandler, userHandler, categoryHandler, labelingHandler, app, conf)
	} else {
		help_command()
	}

}

func application(mail handler.MailInterface, user handler.UserInterface, category handler.CategoryInterface, labeling handler.LabelingInterface, app *fiber.App, config *config.App) {
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
		log.Fatal("Error while running server: " + err.Error())
	}

	log.Println("App running at: " + config.Server.Host + ":" + config.Server.Port)
}

func seeder(db *gorm.DB, conf config.App) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(conf.Administrator.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashed password: " + err.Error())
	}

	hashed_user, err := bcrypt.GenerateFromPassword([]byte(conf.Dummy.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashed password: " + err.Error())
	}

	admin := domain.User{
		Username: "Administrator Megaclite",
		Password: string(hashed),
		Email:    conf.Administrator.Email,
		Avatar:   "https://i.pinimg.com/736x/f9/7d/cc/f97dcc1f4c7d2e4ceb47b57dc13060c1.jpg",
		Address:  "JL. Kemuning 200 Apalotega",
	}

	dummy_user := domain.User{
		Username: "Kyora",
		Password: string(hashed_user),
		Email:    conf.Dummy.Email,
		Avatar:   "https://i.pinimg.com/736x/99/73/e7/9973e72bde7835e070e7a8c795522ffb.jpg",
		Address:  "JL. Penuh Hambatan No. 201",
	}

	err = db.Create(&admin).Error
	if err != nil {
		log.Fatal("Error while seeding administrator: " + err.Error())
	}

	err = db.Create(&dummy_user).Error
	if err != nil {
		log.Fatal("Error while seeding dummy user: " + err.Error())
	}

	log.Println("Success seeding data into database")
}

func help_command() {
	content := [][]string{
		{"help", "to show help message"},
		{"migration", "to migration table into database"},
		{"seed", "to seed dummy data into database"},
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
