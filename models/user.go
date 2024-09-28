package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `form:"first_name" json:"first_name,omitempty"`
	LastName  string    `form:"last_name" json:"last_name,omitempty"`
	Password  string    `gorm:"notNull" form:"password" json:"-" binding:"required"`
	Email     string    `gorm:"type:varchar(254);unique; notNull" form:"email" json:"email,omitempty" binding:"required"`
}

func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
