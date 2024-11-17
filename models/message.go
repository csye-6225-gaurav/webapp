package models

import "github.com/gofrs/uuid"

type Message struct {
	Email string    `json:"email"`
	Token uuid.UUID `json:"token"`
}
