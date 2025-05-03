package main

import (
	"errors"
	"fmt"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/models"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/pagination"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"os"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := &models.Book{}

	if err := context.BodyParser(&book); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "CreateBook : request failed"})
	}

	if err := r.DB.Create(&book).Error; err != nil {
		// Если ошибка создания - возвращаем статус 400
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "not create book!"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book added"})

	return nil
}

func (r *Repository) CreateBooks(context *fiber.Ctx) error {
	books := []models.Book{}

	if err := context.BodyParser(&books); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "CreateBooks: invalid request format"})
	}

	if err := r.DB.Create(&books).Error; err != nil {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "failed to create books"})
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "books added successfully",
			"count":   len(books),
		})

	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	book := &models.Book{}
	id := context.Params("id")
	if id == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}

	fmt.Println("ID BOOK:", id)

	result := r.DB.First(&book, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return context.Status(http.StatusNotFound).JSON(&fiber.Map{
				"message": "book not found",
			})
		}
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "database error",
		})
	}

	return context.Status(http.StatusOK).JSON(fiber.Map{
		"message": "book fetched successfully",
		"data":    book,
	})
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	book := &models.Book{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "iid cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(book, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Error! not delete book",
		})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books successfully delete",
	})

	return nil
}

func paginate(value interface{}, pagination *pagination.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	var books []models.Book

	pagination := pagination.Pagination{
		Limit: context.QueryInt("limit", 5),
		Page:  context.QueryInt("page", 1),
		Sort:  context.Query("sort", "id asc"),
	}

	r.DB.Scopes(paginate(books, &pagination, r.DB)).Find(&books)

	pagination.Rows = books

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book pagination successfully",
		"data":    pagination,
	})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/books", r.GetBooks)
	api.Get("/get_books/:id", r.GetBookById)
	api.Post("/create_book", r.CreateBook)
	api.Post("/create_books", r.CreateBooks)
	api.Delete("/delete_book/:id", r.DeleteBook)
}

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
		log.Fatal("No connect database")
	}
	err = models.MigrateBook(db)

	if err != nil {
		log.Fatal(err)
	}

	r := Repository{DB: db}
	app := fiber.New()
	r.SetupRoutes(app)

	app.Listen(":8080")
}
