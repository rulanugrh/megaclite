package config

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
)

type App struct {
	Server struct {
		ApiPort  string
		Host     string
		Secret   string
		ViewPort string
	}

	Database struct {
		Host string
		Name string
		Port string
		User string
		Pass string
	}

	Observability struct {
		Host string
		Port string
	}

	Store *session.Store
}

var app *App

func GetConfig() *App {
	if app == nil {
		app = initConfig()
	}

	return app
}

func initConfig() *App {
	conf := App{}
	if err := godotenv.Load(); err != nil {
		// Adding database config default
		conf.Database.Host = "localhost"
		conf.Database.Name = "db_store"
		conf.Database.Port = "3306"
		conf.Database.Pass = ""
		conf.Database.User = "root"

		// Adding observability host and port default
		conf.Observability.Host = "localhost"
		conf.Observability.Port = "4137"

		// Adding server host and port config default
		conf.Server.Host = "localhost"
		conf.Server.ApiPort = "3000"
		conf.Server.ViewPort = "8080"
		conf.Server.Secret = "HelloWorldThisIsAdministrator"

		return &conf
	}

	conf.Database.Host = os.Getenv("DATABASE_HOST")
	conf.Database.Port = os.Getenv("DATABASE_PORT")
	conf.Database.Name = os.Getenv("DATABASE_NAME")
	conf.Database.Pass = os.Getenv("DATABASE_PASS")
	conf.Database.User = os.Getenv("DATABASE_USER")

	conf.Observability.Host = os.Getenv("OTEL_HOST")
	conf.Observability.Port = os.Getenv("OTEL_PORT")

	conf.Server.Host = os.Getenv("SERVER_HOST")
	conf.Server.ApiPort = os.Getenv("SERVER_API_PORT")
	conf.Server.ViewPort = os.Getenv("SERVER_VIEW_PORT")

	conf.Server.Secret = os.Getenv("SERVER_SECRET")
	conf.Store = session.New(
		session.Config{
			CookieHTTPOnly: true,
			Expiration:     24 * time.Hour,
		},
	)
	return &conf

}
