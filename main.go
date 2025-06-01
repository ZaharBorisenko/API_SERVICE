package main

import (
	"context"
	"fmt"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/handlers"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/models"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	_ "gorm.io/driver/sqlite"
	"log"
	"os"
)

var ctx = context.Background()

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379" // fallback
	}

	cache := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	// Проверка подключения к Redis
	if _, err := cache.Ping(ctx).Result(); err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	fmt.Println("Подключение к Redis успешно")

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
	categoriesHandler := handlers.NewCategoryHandler(db)

	bookHandler.SetupRoutes(app)
	authorHandler.SetupRoutes(app)
	categoriesHandler.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
