package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var BaseURL string

func init() {
	BaseURL = "http://localhost:18080"
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func mustPostJSON(t *testing.T, ctx context.Context, path string, reqBody interface{}) []byte {
	t.Helper()

	var bodyBytes []byte
	var err error
	if reqBody != nil {
		bodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, BaseURL+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute POST request to %s: %v", path, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Logf("Response Body: %s", string(respBody))
		var errResp ErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error.Code != "" {
			t.Fatalf("POST %s failed with status %d. Code: %s, Msg: %s", path, resp.StatusCode, errResp.Error.Code, errResp.Error.Message)
		}
		t.Fatalf("POST %s failed with status %d. Body: %s", path, resp.StatusCode, string(respBody))
	}

	return respBody
}

func mustGetJSON(t *testing.T, ctx context.Context, path string) []byte {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, BaseURL+path, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute GET request to %s: %v", path, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Logf("Response Body: %s", string(respBody))
		var errResp ErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error.Code != "" {
			t.Fatalf("GET %s failed with status %d. Code: %s, Msg: %s", path, resp.StatusCode, errResp.Error.Code, errResp.Error.Message)
		}
		t.Fatalf("GET %s failed with status %d. Body: %s", path, resp.StatusCode, string(respBody))
	}

	return respBody
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createTeamHelper(t *testing.T, ctx context.Context, teamName string, members []TeamMember) *Team {
	t.Helper()

	req := Team{TeamName: teamName, Members: members}
	body := mustPostJSON(t, ctx, "/team/add", req)

	var resp struct {
		Team Team `json:"team"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("Failed to unmarshal created team response: %v", err)
	}
	return &resp.Team
}
