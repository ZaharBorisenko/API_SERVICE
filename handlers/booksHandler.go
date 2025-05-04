package handlers

import (
	"errors"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/models"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/pagination"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
)

type BookHandler struct {
	DB *gorm.DB
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	return &BookHandler{DB: db}
}

func (r *BookHandler) GetBooks(context *fiber.Ctx) error {
	var books []models.Book

	pagination := pagination.Pagination{
		Limit: context.QueryInt("limit", 10),
		Page:  context.QueryInt("page", 1),
		Sort:  context.Query("sort", "id asc"),
	}

	r.DB.Scopes(pagination.Paginate(books, r.DB)).Preload("Author").Find(&books)
	pagination.Rows = books

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book pagination successfully",
		"data":    pagination,
	})
}
func (r *BookHandler) GetBookById(context *fiber.Ctx) error {
	book := &models.Book{}
	id := context.Params("id")
	if id == "" {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}

	result := r.DB.Preload("Author").First(&book, id)

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
func (r *BookHandler) CreateBook(context *fiber.Ctx) error {
	singleBook := models.Book{}
	if err := context.BodyParser(&singleBook); err == nil {
		if err := r.DB.Create(&singleBook).Error; err != nil {
			return context.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "failed to create book", "error": err.Error()})
		}

		return context.Status(http.StatusCreated).JSON(
			&fiber.Map{
				"message": "book created successfully",
				"data":    singleBook,
			})
	}

	multipleBooks := []models.Book{}
	if err := context.BodyParser(&multipleBooks); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "invalid request format", "error": err.Error()})
	}

	if len(multipleBooks) == 0 {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "empty books list provided"})
	}

	if err := r.DB.Create(&multipleBooks).Error; err != nil {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "failed to create books", "error": err.Error()})
	}

	return context.Status(http.StatusCreated).JSON(
		&fiber.Map{
			"message": "books created successfully",
			"count":   len(multipleBooks),
			"data":    multipleBooks,
		})
}
func (r *BookHandler) DeleteBook(context *fiber.Ctx) error {
	book := &models.Book{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
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
		"data":    book,
	})

	return nil
}

func (r *BookHandler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/books")
	api.Get("/", r.GetBooks)
	api.Get("/:id", r.GetBookById)
	api.Post("/create", r.CreateBook)
	api.Delete("/delete/:id", r.DeleteBook)
}
