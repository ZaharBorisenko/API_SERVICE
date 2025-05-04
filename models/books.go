package models

import (
	"time"
)

type Book struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AuthorID    uint      `json:"author_id"`
	Author      Author    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	CoverURL    string    `json:"cover_url" gorm:"type:text;default:'https://avatars.mds.yandex.net/i?id=e5be5c8d1fe86f031ac75d8cf920d292_l-5324012-images-thumbs&n=13.jpg'"`
	Publisher   string    `json:"publisher"`
	Description string    `json:"description"`
	NumberPages uint      `json:"number_pages"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
