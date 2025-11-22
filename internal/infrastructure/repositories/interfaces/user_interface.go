package repositories

import (
	"pullrequest-inator/internal/infrastructure/models"
)

type User interface {
	Repository[models.User, int64]
}
