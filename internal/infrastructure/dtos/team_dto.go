package dtos

import "github.com/google/uuid"

type TeamDTO struct {
	ID      uuid.UUID  `json:"id"`
	Name    string     `json:"name"`
	UserIDs uuid.UUIDs `json:"user_ids"`
}

type TeamUsersDTO struct {
	ID        uuid.UUID  `json:"id"`
	TeamName  string     `json:"team_name"`
	TeamUsers []*UserDTO `json:"team_users"`
}
