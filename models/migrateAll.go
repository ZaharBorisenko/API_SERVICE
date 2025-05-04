package models

import "gorm.io/gorm"

func MigrateAll(db *gorm.DB) error {
	if err := db.AutoMigrate(&Author{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Book{}); err != nil {
		return err
	}
	return nil
}
