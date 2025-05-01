package models

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Publisher   string    `json:"publisher"`
	Description string    `json:"description"`
	NumberPages uint      `json:"number_pages"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func MigrateBook(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}
