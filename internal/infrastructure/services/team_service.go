package services

import (
	"context"
	"errors"
	"fmt"
	"pullrequest-inator/internal/api/dtos"
	"pullrequest-inator/internal/infrastructure/repositories/interfaces"
	"pullrequest-inator/internal/infrastructure/repositories/pg"

	"github.com/google/uuid"
)

var (
	ErrTeamExists      = errors.New("team already exists")
	ErrNotFound        = errors.New("resource not found")
	ErrFalseUserInTeam = errors.New("detected deleted user in team")
)

type TeamService struct {
	teamRepo repositories.Team
	userRepo repositories.User
}

func NewTeamService(teamRepo repositories.Team, userRepo repositories.User) (*TeamService, error) {
	if teamRepo == nil {
		return nil, errors.New("teamRepository cannot be nil")
	}
	if userRepo == nil {
		return nil, errors.New("userRepository cannot be nil")
	}

	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}, nil
}

func (s *TeamService) CreateTeamWithUsers(ctx context.Context, teamReq *dtos.Team) error {
	if _, err := s.teamRepo.FindByName(ctx, teamReq.TeamName); err == nil {
		return ErrTeamExists
	}

	return s.teamRepo.CreateWithUsers(ctx, teamReq)
}

func (s *TeamService) GetTeamByName(ctx context.Context, teamName string) (*dtos.Team, error) {
	teamModel, err := s.teamRepo.FindByName(ctx, teamName)
	if errors.Is(err, pg.ErrTeamNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find team: %w", err)
	}

	var members []dtos.TeamMember
	for _, userID := range teamModel.UserIDs {
		user, err := s.userRepo.FindByID(ctx, userID)
		if errors.Is(err, pg.ErrUserNotFound) {
			return nil, ErrFalseUserInTeam
		}
		if err != nil {
			return nil, fmt.Errorf("find user %s: %w", userID, err)
		}

		members = append(members, dtos.TeamMember{
			UserId:   user.ID.String(),
			Username: user.Username,
			IsActive: user.IsActive,
		})
	}

	return &dtos.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}

func (s *TeamService) SetUserActiveByID(ctx context.Context, userID uuid.UUID, active bool) (*dtos.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if errors.Is(err, pg.ErrUserNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	user.IsActive = active
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	team, err := s.teamRepo.FindByUserID(ctx, userID)
	var teamName string
	if err != nil && !errors.Is(err, pg.ErrTeamNotFound) {
		return nil, fmt.Errorf("find team for user: %w", err)
	}
	if team != nil {
		teamName = team.Name
	}

	return &dtos.User{
		UserId:   user.ID.String(),
		Username: user.Username,
		IsActive: user.IsActive,
		TeamName: teamName,
	}, nil
}
