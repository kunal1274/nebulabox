package cli_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/nebulabox/nebulabox/internal/cli/tests/mocks"
	"github.com/nebulabox/nebulabox/internal/cli/tests/testutils"
)

// TestBuildCommand tests the `nebulabox build` command
func TestBuildCommand(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	mockAPI.RegisterHandler("POST", "/api/images/build", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Dockerfile string `json:"dockerfile"`
			Tag        string `json:"tag"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Build initiated",
			"tag":     req.Tag,
		})
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
	stdout, stderr, err := runner.Run("build", ".", "--tag", "test-image:latest")

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
	}
	testutils.AssertOutputContains(t, stdout, "Build")
	t.Log("Build command test complete")
}

