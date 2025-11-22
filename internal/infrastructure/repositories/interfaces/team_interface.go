package repositories

import (
	"context"
	"pullrequest-inator/internal/api/dtos"
	"pullrequest-inator/internal/infrastructure/models"
)

type Team interface {
	Repository[models.Team, int64]
	FindByName(ctx context.Context, name string) (*models.Team, error)
	FindByUserID(ctx context.Context, userID int64) (*models.Team, error)
	CreateWithUsers(ctx context.Context, teamReq *dtos.Team) error
}
