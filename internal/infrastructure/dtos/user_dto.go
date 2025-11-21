package dtos

import "github.com/google/uuid"

type UserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IsActive bool      `json:"is_active"`
}

type UserTeamDTO struct {
	User     *UserDTO `json:"user"`
	TeamName string   `json:"team_name"`
}
