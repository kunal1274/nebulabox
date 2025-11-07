package cli_test

import (
	"os"
	"testing"

	"github.com/nebulabox/nebulabox/internal/cli/tests/testutils"
)

// TestVersionCommand tests the `nebulabox version` command
func TestVersionCommand(t *testing.T) {
	binaryPath := "../../../bin/nebulabox"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("CLI binary not found at %s. Run 'make build-cli-test' first.", binaryPath)
		return
	}

	runner := testutils.NewCommandRunner(binaryPath)
	stdout, stderr, err := runner.Run("version")

	if err != nil {
		t.Logf("Command stderr: %s", stderr)
	}

	testutils.AssertOutputContains(t, stdout, "NebulaBox")
	testutils.AssertOutputContains(t, stdout, "Version")
	t.Log("Version command test complete")
}

