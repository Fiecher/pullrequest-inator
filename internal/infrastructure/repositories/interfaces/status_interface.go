package repositories

import (
	"context"
	"pullrequest-inator/internal/infrastructure/models"
)

type Status interface {
	FindByID(ctx context.Context, id int64) (*models.Status, error)
	FindAll(ctx context.Context) ([]*models.Status, error)
}
