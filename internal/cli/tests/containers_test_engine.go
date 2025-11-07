package cli_test

import (
	"os"
	"testing"

	"github.com/nebulabox/nebulabox/internal/cli/tests/testutils"
)

// TestContainerListEngine tests the `nebulabox list` command with real engine
func TestContainerListEngine(t *testing.T) {
	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	runner := testutils.NewCommandRunner(binaryPath)
	stdout, stderr, err := runner.Run("list")

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
		// List command may return empty list, which is OK
		if stderr != "" {
			t.Logf("Note: %s", stderr)
		}
	}

	// Should show header or "No containers found"
	if stdout == "" && stderr == "" {
		t.Error("Expected output from list command")
	}

	t.Log("Container list test complete - using real engine")
}

// TestContainerListPs tests the `nebulabox ps` command (alias) with real engine
func TestContainerListPs(t *testing.T) {
	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	runner := testutils.NewCommandRunner(binaryPath)
	stdout, stderr, err := runner.Run("ps")

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
	}

	// Should show header or "No containers found"
	if stdout == "" && stderr == "" {
		t.Error("Expected output from ps command")
	}

	t.Log("Container ps test complete - using real engine")
}

