package services

import (
	"context"
	"pullrequest-inator/internal/api/dtos"
)

type User interface {
	RegisterUser(ctx context.Context, user *dtos.User) error
	UnregisterUserByID(ctx context.Context, userID int64) error
	ListUsers(ctx context.Context) ([]*dtos.User, error)
	SetUserActive(ctx context.Context, userID int64, active bool) (*dtos.User, error)
}
