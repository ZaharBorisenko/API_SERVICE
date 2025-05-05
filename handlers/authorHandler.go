package handlers

import (
	"errors"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/models"
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

type AuthorHandler struct {
	DB *gorm.DB
}

func NewAuthorHandler(db *gorm.DB) *AuthorHandler {
	return &AuthorHandler{DB: db}
}

func (a *AuthorHandler) GetAuthor(context *fiber.Ctx) error {
	author := []models.Author{}

	pagination := pagination.Pagination{
		Limit: context.QueryInt("limit", 10),
		Page:  context.QueryInt("page", 1),
		Sort:  context.Query("sort", "id asc"),
	}

	a.DB.Scopes(pagination.Paginate(author, a.DB)).Preload("Books").Find(&author)
	pagination.Rows = author

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "author pagination successfully",
		"data":    pagination,
	})

}
func (a *AuthorHandler) GetAuthorById(context *fiber.Ctx) error {
	author := models.Author{}
	id, _ := uuid.Parse(context.Params("id"))

	result := a.DB.Preload("Books").Where("id = ?", id).First(&author)

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

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "author fetched successfully", "data": author})

}
func (a *AuthorHandler) CreateAuthor(context *fiber.Ctx) error {
	author := models.Author{}

	if err := context.BodyParser(&author); err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Author request failed"})
	}

	if err := a.DB.Create(&author).Error; err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "author was not created"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "author added", "data": author})

	return nil

}
func (a *AuthorHandler) DeleteAuthor(context *fiber.Ctx) error {
	author := models.Author{}
	id, _ := uuid.Parse(context.Params("id"))

	err := a.DB.Where("id = ?", id).Delete(&author)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Error! not delete book",
		})
		return err.Error
	}

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "author successfully delete", "data": author})

}

func (a *AuthorHandler) UpdateAuthor(context *fiber.Ctx) error {
	authorUpdate := models.Author{}
	id, _ := uuid.Parse(context.Params("id"))
	if err := context.BodyParser(&authorUpdate); err != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid request format",
		})
	}

	result := a.DB.Where("id = ?", id).Updates(&authorUpdate)

	if result.Error != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"messages": "Invalid request",
		})
	}
	return context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "author updated successfully", "data": authorUpdate})
}

func (a *AuthorHandler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/authors")
	api.Get("/", a.GetAuthor)
	api.Get("/:id", a.GetAuthorById)
	api.Put("/update/:id", a.UpdateAuthor)
	api.Delete("/delete/:id", a.DeleteAuthor)
}
