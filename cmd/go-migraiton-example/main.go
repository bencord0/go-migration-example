package main

import (
	"net/url"
	"os"

	"github.com/bencord0/go-migration-example/internal/models"
	"github.com/bencord0/go-migration-example/internal/views"
	"github.com/bencord0/go-migration-example/internal/views/admin"

	. "github.com/bencord0/webframework"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
}

func setupDB() *models.Database {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	parsedUrl, err := url.Parse(databaseUrl)
	if err != nil {
		log.Fatal("Failed to parse database url: ", err)
	}

	db, err := models.NewDatabase(parsedUrl)
	if err != nil {
		log.Fatal("Couldn't create database: ", err)
	}

	return db
}

func main() {
	db := setupDB()

	settings := Settings{
		Addr: "0:8000",
		Urls: []Url{
			{"/users", views.GetUsers(db)},
			{"/users/add", views.AddUser(db)},
			{"/admin/migrations", admin.ListMigrations(db)},
			{"/admin/migrations/apply", admin.ApplyMigrations(db)},
		},
	}

	application := NewApplication(settings)
	log.Printf("Running application on port 8000")
	application.Run()
}
