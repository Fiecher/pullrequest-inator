package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestPRAssignmentAndMerge(t *testing.T) {
	ctx := context.Background()
	t.Parallel()
	t.Log("1. Setup: Creating team...")
	author := TeamMember{UserID: "auth" + generateRandomString(4), Username: "Author", IsActive: true}
	rev1 := TeamMember{UserID: "rev1" + generateRandomString(4), Username: "R1", IsActive: true}
	rev2 := TeamMember{UserID: "rev2" + generateRandomString(4), Username: "R2", IsActive: true}
	rev3 := TeamMember{UserID: "rev3" + generateRandomString(4), Username: "R3", IsActive: true}
	inactive := TeamMember{UserID: "inact" + generateRandomString(4), Username: "Ina", IsActive: false}

	teamName := "LargeTeam" + generateRandomString(4)
	createTeamHelper(t, ctx, teamName, []TeamMember{author, rev1, rev2, rev3, inactive})

	t.Log("2. Creating PR and checking auto-assignment...")
	prID := "pr" + generateRandomString(5)
	prReq := CreatePRRequest{
		PullRequestId:   prID,
		PullRequestName: "Feature X",
		AuthorId:        author.UserID,
	}

	bodyPR := mustPostJSON(t, ctx, "/pullRequest/create", prReq)

	var createResp CreatePRResponseWrapper
	if err := json.Unmarshal(bodyPR, &createResp); err != nil {
		t.Fatalf("Failed to unmarshal created PR: %v", err)
	}
	pr := createResp.Pr

	if pr.PullRequestId != prID || pr.Status != "OPEN" {
		t.Fatalf("PR creation failed: %+v", pr)
	}

	if len(pr.AssignedReviewers) != 2 {
		t.Fatalf("Expected 2 reviewers, got %d. List: %v", len(pr.AssignedReviewers), pr.AssignedReviewers)
	}

	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == author.UserID {
			t.Fatalf("Author assigned as reviewer!")
		}
		if reviewerID == inactive.UserID {
			t.Fatalf("Inactive user assigned as reviewer!")
		}
	}

	originalReviewer1 := pr.AssignedReviewers[0]

	t.Logf("3. Reassigning reviewer %s...", originalReviewer1)
	reassignReq := ReassignRequest{
		PullRequestId: prID,
		OldUserId:     originalReviewer1,
	}

	bodyReassign := mustPostJSON(t, ctx, "/pullRequest/reassign", reassignReq)

	var reassignResp ReassignResponse
	if err := json.Unmarshal(bodyReassign, &reassignResp); err != nil {
		t.Fatalf("Failed to unmarshal reassign response: %v", err)
	}

	prAfter := reassignResp.Pr
	if len(prAfter.AssignedReviewers) != 2 {
		t.Fatalf("Expected 2 reviewers, got %d", len(prAfter.AssignedReviewers))
	}

	isOriginalPresent := false
	for _, id := range prAfter.AssignedReviewers {
		if id == originalReviewer1 {
			isOriginalPresent = true
		}
	}
	if isOriginalPresent {
		t.Fatalf("Original reviewer %s was not removed!", originalReviewer1)
	}

	t.Log("4. Marking PR as MERGED...")
	mergeReq := MergePRRequest{PullRequestId: prID}
	bodyMerge := mustPostJSON(t, ctx, "/pullRequest/merge", mergeReq)

	var mergeResp CreatePRResponseWrapper
	if err := json.Unmarshal(bodyMerge, &mergeResp); err != nil {
		t.Fatalf("Failed to unmarshal merge response: %v", err)
	}
	if mergeResp.Pr.Status != "MERGED" {
		t.Fatalf("Expected status MERGED, got %s", mergeResp.Pr.Status)
	}

	t.Log("5. Testing MERGE idempotency...")
	mustPostJSON(t, ctx, "/pullRequest/merge", mergeReq)

	t.Log("6. Testing Reassignment restriction after MERGED...")
	bodyBytes, _ := json.Marshal(reassignReq)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, BaseURL+"/pullRequest/reassign", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("Expected 409 Conflict, got %d", resp.StatusCode)
	}
}

func TestPRAssignmentLessThanTwo(t *testing.T) {
	ctx := context.Background()
	t.Parallel()

	author := TeamMember{UserID: "authS" + generateRandomString(4), Username: "AuthS", IsActive: true}
	rev1 := TeamMember{UserID: "revS" + generateRandomString(4), Username: "RevS", IsActive: true}

	teamName := "SmallTeam" + generateRandomString(4)
	createTeamHelper(t, ctx, teamName, []TeamMember{author, rev1})

	prID := "prSmall" + generateRandomString(4)
	prReq := CreatePRRequest{
		PullRequestId:   prID,
		PullRequestName: "Small PR",
		AuthorId:        author.UserID,
	}

	bodyPR := mustPostJSON(t, ctx, "/pullRequest/create", prReq)
	var createResp CreatePRResponseWrapper
	err := json.Unmarshal(bodyPR, &createResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal created PR: %v", err)
		return
	}

	if len(createResp.Pr.AssignedReviewers) != 1 {
		t.Fatalf("Expected 1 reviewer, got %d", len(createResp.Pr.AssignedReviewers))
	}
	if createResp.Pr.AssignedReviewers[0] != rev1.UserID {
		t.Fatalf("Wrong reviewer assigned. Expected %s, got %s", rev1.UserID, createResp.Pr.AssignedReviewers[0])
	}
}
