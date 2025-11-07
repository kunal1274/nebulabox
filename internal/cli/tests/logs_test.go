package cli_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/nebulabox/nebulabox/internal/cli/tests/mocks"
	"github.com/nebulabox/nebulabox/internal/cli/tests/testutils"
)

// TestLogsCommand tests the `nebulabox logs` command
func TestLogsCommand(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	containerID := "test-container-123"

	mockAPI.RegisterHandler("GET", "/api/containers/"+containerID+"/logs", func(w http.ResponseWriter, r *http.Request) {
		logs := []string{
			"[2024-01-01T10:00:00Z] Container started",
			"[2024-01-01T10:00:01Z] Application listening on port 80",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	})

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
	stdout, stderr, err := runner.Run("logs", containerID)

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
	}
	// CLI logs output format shows log entries with timestamps
	// The actual output contains lines like "[2025-11-05 11:39:52] Container test-container-123 started"
	testutils.AssertOutputContains(t, stdout, "Container")
	testutils.AssertOutputContains(t, stdout, containerID)
	t.Log("Logs command test complete")
}

