package handlers

import (
	"errors"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/dto"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/models"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
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
	var totalRows int64

	pagination := pagination.Pagination{
		Limit: context.QueryInt("limit", 10),
		Page:  context.QueryInt("page", 1),
		Sort:  context.Query("sort", "id asc"),
	}

	if err := r.DB.Model(&models.Book{}).Count(&totalRows).Error; err != nil {
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "failed to count books",
			"error":   err.Error(),
		})
	}

	if err := r.DB.
		Scopes(pagination.Paginate(books, r.DB)).
		Preload("Author").Preload("Categories").
		Find(&books).Error; err != nil {
		return context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "failed to fetch books",
			"error":   err.Error(),
		})
	}

	bookResponses := make([]dto.BookResponse, len(books))
	for i, book := range books {
		bookResponses[i] = dto.ToBookResponse(book)
	}

	response := fiber.Map{
		"message": "books fetched successfully",
		"data": fiber.Map{
			"books": bookResponses,
			"pagination": fiber.Map{
				"page":        pagination.Page,
				"limit":       pagination.Limit,
				"total_rows":  totalRows,
				"total_pages": int(math.Ceil(float64(totalRows) / float64(pagination.Limit))),
			},
		},
	}

	return context.Status(http.StatusOK).JSON(response)
}

func (r *BookHandler) GetBookById(context *fiber.Ctx) error {
	book := models.Book{}
	id, _ := uuid.Parse(context.Params("id"))

	result := r.DB.
		Preload("Author").
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		First(&book, "id = ?", id)

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
		"data":    dto.ToBookResponse(book),
	})
}

func (r *BookHandler) CreateBook(context *fiber.Ctx) error {

	var request dto.CreateBookRequest
	if err := context.BodyParser(&request); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request format",
			"error":   err.Error(),
		})
	}

	// Создаем книгу
	book := models.Book{
		AuthorID:    request.AuthorID,
		Title:       request.Title,
		CoverURL:    request.CoverURL,
		Publisher:   request.Publisher,
		Description: request.Description,
		NumberPages: request.NumberPages,
	}

	// Находим категории по их ID
	var categories []models.Categories
	if len(request.CategoryIDs) > 0 {
		if err := r.DB.Where("id IN ?", request.CategoryIDs).Find(&categories).Error; err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to find categories",
				"error":   err.Error(),
			})
		}
	}

	// Создаем книгу в базе данных
	if err := r.DB.Create(&book).Error; err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to create book",
			"error":   err.Error(),
		})
	}

	// Добавляем связи с категориями
	if len(categories) > 0 {
		if err := r.DB.Model(&book).Association("Categories").Append(categories); err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to associate categories with book",
				"error":   err.Error(),
			})
		}
	}

	// Загружаем книгу с категориями для ответа
	var createdBook models.Book
	if err := r.DB.Preload("Categories").First(&createdBook, book.ID).Error; err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch created book",
			"error":   err.Error(),
		})
	}

	return context.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "book created successfully",
		"data":    createdBook,
	})
}

func (r *BookHandler) DeleteBook(context *fiber.Ctx) error {
	book := models.Book{}
	id, _ := uuid.Parse(context.Params("id"))

	err := r.DB.Where("id = ?", id).Delete(&book)

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

func (r *BookHandler) UpdateBook(context *fiber.Ctx) error {
	updateBook := models.Book{}
	id, _ := uuid.Parse(context.Params("id"))

	if err := context.BodyParser(&updateBook); err != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid request format",
		})
	}

	result := r.DB.Where("id = ?", id).Updates(&updateBook)

	if result.Error != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"messages": "Invalid request",
		})
	}
	return context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book updated successfully", "data": updateBook})
}

func (r *BookHandler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/books")
	api.Get("/", r.GetBooks)
	api.Get("/:id", r.GetBookById)
	api.Post("/create", r.CreateBook)
	api.Put("/update/:id", r.UpdateBook)
	api.Delete("/delete/:id", r.DeleteBook)
}
