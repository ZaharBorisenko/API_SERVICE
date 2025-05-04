package main

import (
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/handlers"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/models"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "gorm.io/driver/sqlite"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	if err := models.MigrateAll(db); err != nil {
		log.Fatal("Could not migrate models")
	}

	app := fiber.New()

	bookHandler := handlers.NewBookHandler(db)
	authorHandler := handlers.NewAuthorHandler(db)

	bookHandler.SetupRoutes(app)
	authorHandler.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
