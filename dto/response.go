package dto

import (
	"github.com/ZaharBorisenko/GOLAND_API_BOOKS/models"
	"github.com/google/uuid"
	"time"
)

type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BookResponse struct {
	Author      AuthorPartial      `json:"author"`
	ID          uuid.UUID          `json:"id"`
	AuthorID    uuid.UUID          `json:"author_id"`
	Title       string             `json:"title"`
	Categories  []CategoryResponse `json:"categories"`
	CoverURL    string             `json:"cover_url"`
	Publisher   string             `json:"publisher"`
	Description string             `json:"description"`
	NumberPages uint               `json:"number_pages"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type AuthorPartial struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CoverURL  string `json:"cover_url"`
}

type BooksListResponse struct {
	Books      []BookResponse `json:"books"`
	Pagination *Pagination    `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalRows  int `json:"total_rows"`
	TotalPages int `json:"total_pages"`
}

func ToBookResponse(book models.Book) BookResponse {

	categories := make([]CategoryResponse, len(book.Categories))
	for i, category := range book.Categories {
		categories[i] = CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		}
	}

	return BookResponse{
		Author: AuthorPartial{
			FirstName: book.Author.FirstName,
			LastName:  book.Author.LastName,
			CoverURL:  book.Author.CoverURL,
		},
		ID:          book.ID,
		AuthorID:    book.AuthorID,
		Title:       book.Title,
		Categories:  categories,
		CoverURL:    book.CoverURL,
		Publisher:   book.Publisher,
		Description: book.Description,
		NumberPages: book.NumberPages,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}
