package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB     *gorm.DB
	config *App
}

func InitDatabase(app *App) *Database {
	return &Database{config: app}
}

func (conn *Database) Connection() *gorm.DB {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		conn.config.Database.User,
		conn.config.Database.Pass,
		conn.config.Database.Host,
		conn.config.Database.Port,
		conn.config.Database.Name,
	)

	// Setting for logger query to database
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		log.Printf("Error while connect to DB: %s", err.Error())
	}

	sql, err := db.DB()
	if err != nil {
		log.Printf("Error while set for configuration DB: %s", err.Error())
	}

	// Set Max Open Connection to DB
	sql.SetMaxOpenConns(100)
	// Set Idle Connection DB
	sql.SetMaxIdleConns(10)
	// Set Max Lifetime for Connection
	sql.SetConnMaxLifetime(time.Since(time.Now().Add(30 * time.Minute)))
	// Max Lifetime for idle connection
	sql.SetConnMaxIdleTime(time.Since(time.Now().Add(1 * time.Minute)))

	conn.DB = db
	return db
}
