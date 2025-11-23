package e2e

import (
	"context"
	"encoding/json"
	"testing"
)

func TestTeamAndUserLifecycle(t *testing.T) {
	ctx := context.Background()
	t.Parallel()

	teamName := "DevTeam" + generateRandomString(4)
	userA := TeamMember{UserID: "uA" + generateRandomString(4), Username: "Alice", IsActive: true}
	userB := TeamMember{UserID: "uB" + generateRandomString(4), Username: "Bob", IsActive: true}
	userC := TeamMember{UserID: "uC" + generateRandomString(4), Username: "Charlie", IsActive: false}

	t.Logf("1. Creating team %s with users A, B, C...", teamName)
	createdTeam := createTeamHelper(t, ctx, teamName, []TeamMember{userA, userB, userC})

	if createdTeam.TeamName != teamName {
		t.Fatalf("Expected team name %s, got %s", teamName, createdTeam.TeamName)
	}
	if len(createdTeam.Members) != 3 {
		t.Fatalf("Expected 3 members in team, got %d", len(createdTeam.Members))
	}

	t.Log("2. Retrieving team details...")
	bodyTeam := mustGetJSON(t, ctx, "/team/get?team_name="+teamName)
	var retrievedTeam Team
	if err := json.Unmarshal(bodyTeam, &retrievedTeam); err != nil {
		t.Fatalf("Failed to unmarshal retrieved team: %v", err)
	}

	var activeCount int
	for _, member := range retrievedTeam.Members {
		if member.UserID == userA.UserID && !member.IsActive {
			t.Errorf("User A (%s) should be active", userA.UserID)
		}
		if member.UserID == userC.UserID && member.IsActive {
			t.Errorf("User C (%s) should be inactive", userC.UserID)
		}
		if member.IsActive {
			activeCount++
		}
	}
	if activeCount != 2 {
		t.Fatalf("Expected 2 active members (A, B), got %d", activeCount)
	}

	t.Logf("3. Deactivating User B (%s)...", userB.UserID)
	reqActive := SetActiveRequest{UserId: userB.UserID, IsActive: false}
	bodyUpdate := mustPostJSON(t, ctx, "/users/setIsActive", reqActive)

	var activeResp SetActiveResponse
	if err := json.Unmarshal(bodyUpdate, &activeResp); err != nil {
		t.Fatalf("Failed to unmarshal updated user: %v", err)
	}
	if activeResp.User.IsActive {
		t.Fatalf("Expected User B to be inactive, got active")
	}

	t.Log("4. Checking team details again...")
	bodyTeam2 := mustGetJSON(t, ctx, "/team/get?team_name="+teamName)
	var retrievedTeam2 Team
	err := json.Unmarshal(bodyTeam2, &retrievedTeam2)
	if err != nil {
		t.Fatalf("Failed to unmarshal retrieved team: %v", err)
		return
	}

	activeCount = 0
	for _, member := range retrievedTeam2.Members {
		if member.IsActive {
			activeCount++
		}
	}
	if activeCount != 1 {
		t.Fatalf("Expected 1 active member (A) after deactivating B, got %d", activeCount)
	}
}
