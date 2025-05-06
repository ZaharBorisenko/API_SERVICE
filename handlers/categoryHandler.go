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

type CategoryHandler struct {
	DB *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

func (c *CategoryHandler) GetCategories(context *fiber.Ctx) error {
	categories := []models.Categories{}

	pagination := pagination.Pagination{
		Limit: context.QueryInt("limit", 10),
		Page:  context.QueryInt("page", 1),
		Sort:  context.Query("sort", "id asc"),
	}

	c.DB.Scopes(pagination.Paginate(categories, c.DB)).Find(&categories)
	pagination.Rows = categories

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "categories pagination successfully",
		"data":    pagination,
	})
}

func (c *CategoryHandler) GetCategoryById(context *fiber.Ctx) error {
	category := models.Categories{}
	id, _ := uuid.Parse(context.Params("id"))

	result := c.DB.Where("id = ?", id).Find(&category)

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

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "category fetched successfully",
		"data":    category,
	})
}

func (c *CategoryHandler) CreateCategory(context *fiber.Ctx) error {
	category := models.Categories{}

	if err := context.BodyParser(&category); err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Category request failed"})
	}

	if err := c.DB.Create(&category).Error; err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "category was not created",
		})
		return err
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "author added",
		"data":    category,
	})

}

func (c *CategoryHandler) UpdateCategory(context *fiber.Ctx) error {

	categoryUpdate := models.Categories{}
	id, _ := uuid.Parse(context.Params("id"))

	if err := context.BodyParser(&categoryUpdate); err != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid request format",
		})
	}

	err := c.DB.Where("id = ?", id).Updates(&categoryUpdate)

	if err.Error != nil {
		return context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"messages": "Invalid request",
		})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "category updates",
		"data":    categoryUpdate,
	})
}

func (c *CategoryHandler) DeleteCategory(context *fiber.Ctx) error {
	category := models.Categories{}
	id, _ := uuid.Parse(context.Params("id"))

	err := c.DB.Where("id = ?", id).Delete(&category)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Error! not delete category",
		})
		return err.Error
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "category delete",
		"data":    category,
	})
}

func (c *CategoryHandler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/categories")
	api.Get("/", c.GetCategories)
	api.Post("/create", c.CreateCategory)
	api.Get("/:id", c.GetCategoryById)
	api.Put("/update/:id", c.UpdateCategory)
	api.Delete("/delete/:id", c.DeleteCategory)
}
