package cli_test

import (
	"os"
	"testing"

	"github.com/nebulabox/nebulabox/internal/cli/tests/testutils"
)

// TestImageList tests the `nebulabox images` command with real engine
func TestImageList(t *testing.T) {
	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	runner := testutils.NewCommandRunner(binaryPath)
	stdout, stderr, err := runner.Run("images")

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
		// Images command may return empty list, which is OK
		if stderr != "" && !contains(stderr, "No images found") {
			t.Errorf("Unexpected error: %v", err)
		}
	}

	// Should show header or "No images found"
	if stdout == "" && stderr == "" {
		t.Error("Expected output from images command")
	}

	t.Log("Image list test complete - using real engine")
}

// TestImageDelete tests the `nebulabox rmi` command with real engine
func TestImageDelete(t *testing.T) {
	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	// Try to delete a non-existent image (should show error, not crash)
	runner := testutils.NewCommandRunner(binaryPath)
	stdout, stderr, err := runner.Run("rmi", "nonexistent-image:latest")

	// Should handle error gracefully
	if err != nil {
		// Error is expected for non-existent image
		t.Logf("Expected error for non-existent image: %v", err)
	}

	// Should show some output (error message or success)
	if stdout == "" && stderr == "" {
		t.Error("Expected output from rmi command")
	}

	t.Log("Image delete test complete - using real engine")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

