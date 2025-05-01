package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return db, err
	}

	return db, nil

}
