package testutils

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// CommandRunner helps execute CLI commands in tests
type CommandRunner struct {
	BinaryPath string
	Env        []string
}

// NewCommandRunner creates a new command runner
func NewCommandRunner(binaryPath string) *CommandRunner {
	return &CommandRunner{
		BinaryPath: binaryPath,
		Env:        os.Environ(),
	}
}

// Run executes a command and returns stdout, stderr, and error
func (cr *CommandRunner) Run(args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(cr.BinaryPath, args...)
	cmd.Env = cr.Env

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

// RunWithInput executes a command with input
func (cr *CommandRunner) RunWithInput(input string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(cr.BinaryPath, args...)
	cmd.Env = cr.Env

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", "", err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()

	err = cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

// AssertOutputContains checks if output contains expected text
func AssertOutputContains(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', but got: %s", expected, output)
	}
}

// AssertOutputNotContains checks if output does not contain text
func AssertOutputNotContains(t *testing.T, output, unexpected string) {
	t.Helper()
	if strings.Contains(output, unexpected) {
		t.Errorf("Expected output to not contain '%s', but got: %s", unexpected, output)
	}
}

// AssertExitCode checks if command exited with expected code
func AssertExitCode(t *testing.T, err error, expectedCode int) {
	t.Helper()
	if err == nil && expectedCode != 0 {
		t.Errorf("Expected exit code %d, but command succeeded", expectedCode)
		return
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != expectedCode {
				t.Errorf("Expected exit code %d, but got %d", expectedCode, exitErr.ExitCode())
			}
		} else {
			t.Errorf("Command failed with unexpected error: %v", err)
		}
	}
}

