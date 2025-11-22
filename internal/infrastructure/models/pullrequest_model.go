package models

import (
	"time"
)

type PullRequest struct {
	ID        int64      `db:"id"`
	Title     string     `db:"title"`
	AuthorID  int64      `db:"author_id"`
	StatusID  int64      `db:"status_id"`
	MergedAt  *time.Time `db:"merged_at"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`

	ReviewersIDs []int64
}
