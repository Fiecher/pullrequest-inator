package dtos

import "github.com/google/uuid"

type StatusDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
