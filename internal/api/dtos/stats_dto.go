package dtos

type StatsResponse struct {
	TotalPullRequests  int             `json:"total_pull_requests"`
	OpenPullRequests   int             `json:"open_pull_requests"`
	MergedPullRequests int             `json:"merged_pull_requests"`
	ReviewerStats      []ReviewerStats `json:"reviewer_stats"`
}

type ReviewerStats struct {
	ReviewerID    string `json:"reviewer_id"`
	Username      string `json:"username"`
	AssignedCount int    `json:"assigned_count"`
}
