package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PullRequestDTO struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	AuthorID    uuid.UUID  `json:"author_id"`
	Status      StatusDTO  `json:"status"`
	MergedAt    *time.Time `json:"merged_at,omitempty"`
	ReviewerIDs uuid.UUIDs `json:"reviewer_ids"`
}

type PullRequestNewReviewerDTO struct {
	PullRequest   *PullRequestDTO `json:"pull_request"`
	NewReviewerID uuid.UUID       `json:"new_reviewer_id"`
}

type ReassignReviewerResponse struct {
	NewReviewerID uuid.UUID      `json:"new_reviewer_id"`
	PullRequest   PullRequestDTO `json:"pull_request"`
}
