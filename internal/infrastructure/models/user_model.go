package models

import (
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	TeamIDs []int64
}
