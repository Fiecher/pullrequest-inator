package services

import (
	"context"
	"pullrequest-inator/internal/api/dtos"
)

type PullRequest interface {
	CreatePullRequest(ctx context.Context, pr *dtos.PullRequest) (*dtos.PullRequest, error)
	ReassignReviewer(ctx context.Context, userID int64, prID int64) (*dtos.ReassignReviewerResponse, error)
	FindPullRequestsByReviewer(ctx context.Context, userID int64) ([]*dtos.PullRequest, error)
	MarkAsMerged(ctx context.Context, prID int64) (*dtos.PullRequest, error)
	GetUserReviews(ctx context.Context, userID int64) (*dtos.UserGetReviewResponse, error)
	CreateWithReviewers(ctx context.Context, prID int64, prName string, authorID int64) (*dtos.PullRequest, error)
}
