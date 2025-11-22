package repositories

import (
	"context"
	"pullrequest-inator/internal/infrastructure/models"
)

type PullRequest interface {
	Repository[models.PullRequest, int64]
	FindByReviewer(ctx context.Context, userID int64) ([]*models.PullRequest, error)
}
