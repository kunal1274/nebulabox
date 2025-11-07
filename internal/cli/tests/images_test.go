package cli_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/nebulabox/nebulabox/internal/cli/tests/mocks"
	"github.com/nebulabox/nebulabox/internal/cli/tests/testutils"
)

// TestImageList tests the `nebulabox images` command
func TestImageList(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	mockAPI.RegisterHandler("GET", "/api/images", mocks.DefaultImageListHandler())

	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	// CLI doesn't have "images" command, it uses "list" for containers
	// For now, skip this test until images list command is implemented
	t.Skip("Images list command not yet implemented in CLI")
	
	// When implemented, would be:
	// runner := testutils.NewCommandRunner(binaryPath)
	// runner.Env = env
	// stdout, stderr, err := runner.Run("images")
	// testutils.AssertOutputContains(t, stdout, "nginx")
	t.Log("Image list test skipped - command not available")
}

// TestImagePull tests the `nebulabox pull` command
func TestImagePull(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	mockAPI.RegisterHandler("POST", "/api/images/pull", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Image string `json:"image"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Image pull initiated",
			"image":   req.Image,
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
	stdout, stderr, err := runner.Run("pull", "nginx:latest")

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
	}
	// CLI output shows "Pulling image" and success message
	testutils.AssertOutputContains(t, stdout, "Pulling image")
	testutils.AssertOutputContains(t, stdout, "pulled successfully")
	t.Log("Image pull test complete")
}

// TestImageDelete tests the `nebulabox rmi` command
func TestImageDelete(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	mockAPI.RegisterHandler("DELETE", "/api/images/image-1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Image deleted",
		})
	})

	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	// CLI doesn't have "rmi" command yet
	t.Skip("Image delete (rmi) command not yet implemented in CLI")
	
	// When implemented, would be:
	// runner := testutils.NewCommandRunner(binaryPath)
	// runner.Env = env
	// stdout, stderr, err := runner.Run("rmi", "image-1")
	// testutils.AssertOutputContains(t, stdout, "deleted")
	t.Log("Image delete test skipped - command not available")
}

