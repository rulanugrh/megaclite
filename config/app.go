package config

import (
	"os"

	"github.com/joho/godotenv"
)

type App struct {
	Server struct {
		Port string
		Host string
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
		conf.Server.Port = "3000"
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
	conf.Server.Port = os.Getenv("SERVER_PORT")

	return &conf

}
