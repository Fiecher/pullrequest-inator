package dtos

import (
	"pullrequest-inator/internal/infrastructure/encoding"
	"pullrequest-inator/internal/infrastructure/models"
)

func ModelToPullRequestDTO(pr *models.PullRequest, statusName string) *PullRequest {
	pullRequestID := encoding.EncodeID(pr.ID)
	authorID := encoding.EncodeID(pr.AuthorID)

	return &PullRequest{
		PullRequestId:     pullRequestID,
		PullRequestName:   pr.Title,
		AuthorId:          authorID,
		Status:            PullRequestStatus(statusName),
		AssignedReviewers: idsToStrings(pr.ReviewersIDs),
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}
func idsToStrings(ids []int64) []string {
	strings := make([]string, len(ids))
	for i, id := range ids {
		s := encoding.EncodeID(id)
		strings[i] = s
	}
	return strings
}
