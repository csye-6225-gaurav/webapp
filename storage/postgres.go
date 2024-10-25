package storage

import (
	"fmt"

	"github.com/csye-6225-gaurav/webapp/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	Host    string
	Port    string
	Pass    string
	User    string
	DBname  string
	SSLMode string
}

func NewConnection(config *Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Pass, config.DBname, config.Port, config.SSLMode,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.AutoMigrate(&models.User{}, &models.Image{})
	if err != nil {
		return err
	}
	return nil
}
