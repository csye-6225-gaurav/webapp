package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt  time.Time `json:"account_created"`
	UpdatedAt  time.Time `json:"account_updated"`
	FirstName  string    `form:"first_name" json:"first_name,omitempty"`
	LastName   string    `form:"last_name" json:"last_name,omitempty"`
	Password   string    `gorm:"notNull" form:"password" json:"password,omitempty"`
	Email      string    `gorm:"type:varchar(254);unique; notNull" form:"email" json:"email,omitempty"`
	IsVerified bool      `gorm:"type:bool; default:false" json:"is_verified,omitempty"`
	Token      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"token,omitempty"`
}

type ResponseUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"account_created"`
	UpdatedAt time.Time `json:"account_updated"`
	FirstName string    `form:"first_name" json:"first_name,omitempty"`
	LastName  string    `form:"last_name" json:"last_name,omitempty"`
	Email     string    `form:"email" json:"email,omitempty"`
}

type UpdateUser struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password"`
}

type RequestUser struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password"`
}
