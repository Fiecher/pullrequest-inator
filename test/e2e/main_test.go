package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var BaseURL1 = "http://localhost:18080"

func TestMain(m *testing.M) {
	ctx := context.Background()
	e2eDir := filepath.Join("..", "e2e")

	cmdDown := exec.CommandContext(ctx, "docker", "compose", "-f", "docker-compose.test.yml", "down", "-v")
	cmdDown.Dir = e2eDir
	cmdDown.Run()

	cmdUp := exec.CommandContext(ctx, "docker", "compose", "-f", "docker-compose.test.yml", "up", "--build", "-d")
	cmdUp.Dir = e2eDir
	if out, err := cmdUp.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("failed to start containers: %v\n%s", err, out))
	}

	if err := waitForHealth(BaseURL1+"/health", 30*time.Second); err != nil {
		panic(fmt.Sprintf("service is not healthy: %v", err))
	}

	code := m.Run()

	cmdDown2 := exec.CommandContext(ctx, "docker", "compose", "-f", "docker-compose.test.yml", "down", "-v")
	cmdDown2.Dir = e2eDir
	cmdDown2.Run()

	os.Exit(code)
}

func waitForHealth(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			if err := resp.Body.Close(); err != nil {
				return err
			}
			return nil
		}
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				return err
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", url)
}
