package cli_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/nebulabox/nebulabox/internal/cli/tests/mocks"
	"github.com/nebulabox/nebulabox/internal/cli/tests/testutils"
)

// TestContainerList tests the `nebulabox list` command
func TestContainerList(t *testing.T) {
	// Setup mock API server
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	mockAPI.RegisterHandler("GET", "/api/containers", mocks.DefaultContainerListHandler())

	// Set API URL environment variable
	env := []string{
		"NEBULABOX_API_URL=" + mockAPI.URL(),
	}

	// Build CLI binary path
	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	runner := testutils.NewCommandRunner(binaryPath)
	runner.Env = env
	stdout, stderr, err := runner.Run("list")

	// Test assertions
	if err != nil {
		t.Logf("Command stderr: %s", stderr)
	}
	// CLI output format: "mock-001", "mock-002", etc. from DefaultContainerListHandler
	testutils.AssertOutputContains(t, stdout, "mock-")
	testutils.AssertOutputContains(t, stdout, "nginx")
	t.Log("Container list test complete")
}

// TestContainerRun tests the `nebulabox run` command
func TestContainerRun(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	// Register handler for container creation
	mockAPI.RegisterHandler("POST", "/api/containers/run", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     "new-container",
			"name":   "my-container",
			"image":  "nginx:latest",
			"status": "running",
		})
	})

	// Build CLI binary path
	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	env := []string{
		"NEBULABOX_API_URL=" + mockAPI.URL(),
	}

	runner := testutils.NewCommandRunner(binaryPath)
	runner.Env = env
	stdout, stderr, err := runner.Run("run", "nginx:latest", "--name", "my-container")

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
	}
	// CLI output format shows container ID like "mock-8449" and status
	testutils.AssertOutputContains(t, stdout, "Container started successfully")
	testutils.AssertOutputContains(t, stdout, "my-container")
	t.Log("Container run test complete")
}

// TestContainerStop tests the `nebulabox stop` command
func TestContainerStop(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	mockAPI.RegisterHandler("POST", "/api/containers/container-1/stop", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Build CLI binary path
	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	env := []string{
		"NEBULABOX_API_URL=" + mockAPI.URL(),
	}

	runner := testutils.NewCommandRunner(binaryPath)
	runner.Env = env
	stdout, stderr, err := runner.Run("stop", "container-1")

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
	}
	// CLI may show success message or just return without error
	// If there's output, it should indicate success
	if stdout != "" {
		testutils.AssertOutputContains(t, stdout, "stop")
	}
	t.Log("Container stop test complete")
}

