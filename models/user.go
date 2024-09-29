package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `form:"first_name" json:"first_name,omitempty"`
	LastName  string    `form:"last_name" json:"last_name,omitempty"`
	Password  string    `gorm:"notNull" form:"password" json:"password,omitempty"`
	Email     string    `gorm:"type:varchar(254);unique; notNull" form:"email" json:"email,omitempty"`
}

type UpdateUser struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password"`
}
