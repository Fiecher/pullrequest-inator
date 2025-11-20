package repositories

import (
	"pullrequest-manager/internal/domain/models"

	"github.com/google/uuid"
)

type Status interface {
	FindByID(id uuid.UUID) (*models.Status, error)
	FindAll() ([]*models.Status, error)
}
