package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Image struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FileName   string    `gorm:"not null" json:"file_name"`
	URL        string    `gorm:"not null" json:"url"`
	UploadDate time.Time `json:"upload_date"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Setting up the foreign key relationship
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
