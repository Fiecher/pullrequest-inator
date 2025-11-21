package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"pullrequest-manager/internal/domain/models"
	"pullrequest-manager/internal/infrastructure/dtos"
	"pullrequest-manager/internal/infrastructure/repositories"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPRAlreadyExists    = errors.New("pull request already exists")
	ErrUserNotReviewer    = errors.New("user is not a reviewer")
	ErrNoReviewCandidates = errors.New("no users available to review")
	ErrPRAlreadyMerged    = errors.New("cannot change PR state because already merged")
)

type DefaultPullRequestService struct {
	userRepo   repositories.User
	prRepo     repositories.PullRequest
	teamRepo   repositories.Team
	statusRepo repositories.Status
}

func NewDefaultPullRequestService(
	userRepo repositories.User,
	prRepo repositories.PullRequest,
	teamRepo repositories.Team,
	statusRepo repositories.Status,
) (*DefaultPullRequestService, error) {
	if userRepo == nil || prRepo == nil || teamRepo == nil || statusRepo == nil {
		return nil, errors.New("repositories cannot be nil")
	}

	return &DefaultPullRequestService{
		userRepo:   userRepo,
		prRepo:     prRepo,
		teamRepo:   teamRepo,
		statusRepo: statusRepo,
	}, nil
}

func (s *DefaultPullRequestService) CreateWithReviewers(ctx context.Context, pr *models.PullRequest) (*dtos.PullRequestDTO, error) {
	existing, err := s.prRepo.FindByID(ctx, pr.ID)
	if err != nil {
		return nil, fmt.Errorf("check for existing PR: %w", err)
	}
	if existing != nil {
		return nil, ErrPRAlreadyExists
	}

	if _, err = s.userRepo.FindByID(ctx, pr.AuthorID); err != nil {
		return nil, fmt.Errorf("find author: %w", err)
	}

	team, err := s.teamRepo.FindByUserID(ctx, pr.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("find team for author: %w", err)
	}

	activeUsers := []uuid.UUID{}
	for _, uid := range team.UserIDs {
		if uid == pr.AuthorID {
			continue
		}
		u, err := s.userRepo.FindByID(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("find user %s for team %s: %w", uid, team.ID, err)
		}
		if u.IsActive {
			activeUsers = append(activeUsers, uid)
		}
	}

	if len(activeUsers) == 0 {
		return nil, ErrNoReviewCandidates
	}

	pr.ReviewersIDs = chooseRandomUsers(activeUsers, 2)

	statuses, err := s.statusRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all statuses: %w", err)
	}
	var openStatusID uuid.UUID
	for _, st := range statuses {
		if st.Name == "OPEN" {
			openStatusID = st.ID
			break
		}
	}
	if openStatusID == uuid.Nil {
		return nil, fmt.Errorf("status 'OPEN' not found in database")
	}
	pr.StatusID = openStatusID

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, fmt.Errorf("create pull request: %w", err)
	}

	return convertPullRequestToDTO(pr, statuses), nil
}

func (s *DefaultPullRequestService) ReassignReviewer(ctx context.Context, userID uuid.UUID, prID uuid.UUID) (*dtos.ReassignReviewerResponse, error) {
	pr, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("find PR for reassignment: %w", err)
	}

	status, err := s.statusRepo.FindByID(ctx, pr.StatusID)
	if err != nil {
		return nil, fmt.Errorf("find status for PR: %w", err)
	}
	if status != nil && status.Name == "MERGED" {
		return nil, ErrPRAlreadyMerged
	}

	isReviewer := false
	reviewerIndex := -1
	for i, rid := range pr.ReviewersIDs {
		if rid == userID {
			isReviewer = true
			reviewerIndex = i
			break
		}
	}
	if !isReviewer {
		return nil, ErrUserNotReviewer
	}

	team, err := s.teamRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find team for reviewer %s: %w", userID, err)
	}

	candidates := []uuid.UUID{}
	for _, uid := range team.UserIDs {
		if uid == pr.AuthorID || contains(pr.ReviewersIDs, uid) {
			continue
		}
		u, err := s.userRepo.FindByID(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("find user %s for team %s: %w", uid, team.ID, err)
		}
		if u.IsActive {
			candidates = append(candidates, uid)
		}
	}

	if len(candidates) == 0 {
		return nil, ErrNoReviewCandidates
	}

	newReviewer := candidates[rand.Intn(len(candidates))]
	pr.ReviewersIDs[reviewerIndex] = newReviewer

	if err := s.prRepo.Update(ctx, pr); err != nil {
		return nil, fmt.Errorf("update pull request after reassignment: %w", err)
	}

	statuses, err := s.statusRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get statuses for response DTO: %w", err)
	}
	prDTO := convertPullRequestToDTO(pr, statuses)

	return &dtos.ReassignReviewerResponse{
		NewReviewerID: newReviewer,
		PullRequest:   *prDTO,
	}, nil
}

func (s *DefaultPullRequestService) MarkAsMerged(ctx context.Context, prID uuid.UUID) (*dtos.PullRequestDTO, error) {
	pr, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("find PR to mark as merged: %w", err)
	}

	statuses, err := s.statusRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all statuses: %w", err)
	}
	var mergedStatusID uuid.UUID
	for _, st := range statuses {
		if st.Name == "MERGED" {
			mergedStatusID = st.ID
			break
		}
	}
	if mergedStatusID == uuid.Nil {
		return nil, fmt.Errorf("status 'MERGED' not found in database")
	}

	if pr.StatusID == mergedStatusID {
		return convertPullRequestToDTO(pr, statuses), nil
	}

	now := time.Now()
	pr.StatusID = mergedStatusID
	pr.MergedAt = &now

	if err := s.prRepo.Update(ctx, pr); err != nil {
		return nil, fmt.Errorf("update pull request to merged: %w", err)
	}

	return convertPullRequestToDTO(pr, statuses), nil
}

func chooseRandomUsers(userIDs []uuid.UUID, max int) []uuid.UUID {
	n := len(userIDs)
	if n <= max {
		return append([]uuid.UUID{}, userIDs...)
	}

	selected := make([]uuid.UUID, 0, max)
	perm := rand.Perm(n)
	for i := 0; i < max; i++ {
		selected = append(selected, userIDs[perm[i]])
	}
	return selected
}

func contains(slice []uuid.UUID, item uuid.UUID) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func convertPullRequestToDTO(pr *models.PullRequest, statuses []*models.Status) *dtos.PullRequestDTO {
	var statusDTO dtos.StatusDTO
	statusFound := false

	if statuses != nil {
		for _, st := range statuses {
			if st.ID == pr.StatusID {
				statusDTO = dtos.StatusDTO{
					ID:   st.ID,
					Name: st.Name,
				}
				statusFound = true
				break
			}
		}
	}

	if !statusFound {
		statusDTO = dtos.StatusDTO{
			ID:   pr.StatusID,
			Name: "",
		}
	}

	return &dtos.PullRequestDTO{
		ID:          pr.ID,
		Title:       pr.Title,
		AuthorID:    pr.AuthorID,
		Status:      statusDTO,
		MergedAt:    pr.MergedAt,
		ReviewerIDs: pr.ReviewersIDs,
	}
}
