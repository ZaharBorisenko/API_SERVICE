package dto

import "github.com/google/uuid"

type CreateBookRequest struct {
	AuthorID    uuid.UUID   `json:"author_id"`
	Title       string      `json:"title"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
	CoverURL    string      `json:"cover_url"`
	Publisher   string      `json:"publisher"`
	Description string      `json:"description"`
	NumberPages uint        `json:"number_pages"`
}
