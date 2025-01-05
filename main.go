package main

import (
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
	"github.com/rulanugrh/megaclite/api"
	"github.com/rulanugrh/megaclite/api/endpoint"
	"github.com/rulanugrh/megaclite/config"
	handler "github.com/rulanugrh/megaclite/internal/http"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
	"github.com/rulanugrh/megaclite/internal/route"
	"github.com/rulanugrh/megaclite/internal/service"
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
	connectionDB.Connection()

	// Initialize Middleware
	pgpMiddleware := middleware.NewPGPMiddleware(pgp, *conf)
	common := middleware.NewCommonMiddleware(conf)

	// Initialize API and Webview
	apiRoute := api.NewAPIRoutes(conf, common)
	viewRoute := route.NewViewRoute(conf, common)

	// Initialize User Komponen
	userRepository := repository.NewUserRepository(*connectionDB)
	userService := service.NewUserService(userRepository, pgpMiddleware)
	userAPI := endpoint.NewUserAPI(userService)
	userView := handler.NewUserView(userService, conf)

	// Initialize Mail Komponen
	mailRepository := repository.NewMailRepository(*connectionDB)
	mailService := service.NewMailService(mailRepository, pgpMiddleware)
	mailAPI := endpoint.NewMailAPI(mailService)
	mailView := handler.NewMailView(mailService)

	// Initialize Category Komponen
	categoryRepository := repository.NewCategoryRepository(*connectionDB)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryAPI := endpoint.NewCategoryAPI(categoryService)

	// Initialize Labeling Komponen
	labelingRepository := repository.NewLabelMailRepository(*connectionDB)
	labelingService := service.NewLabelMailService(labelingRepository)
	labelingAPI := endpoint.NewLabelMailAPI(labelingService)
	labelingView := handler.NewLabelView(labelingService)

	// parsing argument command
	args := os.Args[1]
	if args == "migration" {
		connectionDB.Migration()
	} else if args == "seed" {
		connectionDB.Seeder()
	} else if args == "api" {
		apiRoute.Run(userAPI, labelingAPI, mailAPI, categoryAPI)
	} else if args == "serve" {
		viewRoute.Run(userView, mailView, labelingView)
	} else {
		help_command()
	}

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
