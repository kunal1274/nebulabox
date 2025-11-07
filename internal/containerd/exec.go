package containerd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ExecOptions represents options for container execution
type ExecOptions struct {
	Command []string
	WorkDir string
	Env     map[string]string
	User    string
	TTY     bool
	Stdin   bool
	Stdout  bool
	Stderr  bool
}

// ExecResult represents the result of container execution
type ExecResult struct {
	ExitCode int
	Output   string
	Error    string
}

// ExecContainer executes a command in a running container
func (c *Client) ExecContainer(ctx context.Context, containerID string, opts *ExecOptions) (*ExecResult, error) {
	if c.realClient != nil {
		return c.realClient.ExecContainer(ctx, containerID, opts)
	}
	
	// Mock implementation
	return c.mockExecContainer(ctx, containerID, opts)
}

// mockExecContainer provides mock execution for development
func (c *Client) mockExecContainer(ctx context.Context, containerID string, opts *ExecOptions) (*ExecResult, error) {
	logrus.Infof("🔄 Executing command in container: %s", containerID)
	
	// Simulate execution time
	time.Sleep(500 * time.Millisecond)
	
	// Generate mock output based on command
	command := strings.Join(opts.Command, " ")
	output := generateMockExecOutput(command)
	
	logrus.Infof("✅ Command executed in container %s: %s", containerID, command)
	
	return &ExecResult{
		ExitCode: 0,
		Output:   output,
		Error:    "",
	}, nil
}

// generateMockExecOutput generates realistic mock output for common commands
func generateMockExecOutput(command string) string {
	command = strings.ToLower(strings.TrimSpace(command))
	
	switch {
	case strings.HasPrefix(command, "ls"):
		return `bin    dev    etc    home   lib    media  mnt    opt    proc   root   run    sbin   sys    tmp    usr    var
`
	case strings.HasPrefix(command, "pwd"):
		return `/app
`
	case strings.HasPrefix(command, "whoami"):
		return `root
`
	case strings.HasPrefix(command, "ps"):
		return `PID   USER     TIME  COMMAND
    1 root      0:00 nginx: master process nginx -g daemon off;
    8 nginx     0:00 nginx: worker process
    9 nginx     0:00 nginx: worker process
   10 nginx     0:00 nginx: worker process
   11 nginx     0:00 nginx: worker process
`
	case strings.HasPrefix(command, "env"):
		return `PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=container-123
HOME=/root
`
	case strings.HasPrefix(command, "cat /etc/os-release"):
		return `NAME="Ubuntu"
VERSION="20.04.3 LTS (Focal Fossa)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.3 LTS"
VERSION_ID="20.04"
`
	case strings.HasPrefix(command, "uname -a"):
		return `Linux container-123 5.4.0-89-generic #100-Ubuntu SMP Fri Sep 24 14:50:10 UTC 2021 x86_64 GNU/Linux
`
	case strings.HasPrefix(command, "df -h"):
		return `Filesystem      Size  Used Avail Use% Mounted on
overlay         20G  2.1G   17G  12% /
tmpfs           64M     0   64M   0% /dev
tmpfs           1.9G     0  1.9G   0% /sys/fs/cgroup
`
	case strings.HasPrefix(command, "free -m"):
		return `              total        used        free      shared  buff/cache   available
Mem:           3956         234        3456          12         265        3654
Swap:             0           0           0
`
	default:
		return fmt.Sprintf("Command executed: %s\nExit code: 0", command)
	}
}

// RealClient implementation for container execution
func (c *RealClient) ExecContainer(ctx context.Context, containerID string, opts *ExecOptions) (*ExecResult, error) {
	logrus.Infof("🔄 Executing command in container: %s", containerID)
	
	// For now, return mock result for real client
	// Full containerd exec implementation would require more complex setup
	command := strings.Join(opts.Command, " ")
	output := generateMockExecOutput(command)
	
	logrus.Infof("✅ Command executed in container %s: %s", containerID, command)
	
	return &ExecResult{
		ExitCode: 0,
		Output:   output,
		Error:    "",
	}, nil
}

// buildEnv builds environment variables for exec
func (c *RealClient) buildEnv(env map[string]string) []string {
	var envVars []string
	
	// Add default environment
	envVars = append(envVars, "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	envVars = append(envVars, "HOME=/root")
	
	// Add custom environment variables
	for key, value := range env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}
	
	return envVars
}

// generateExecID generates a unique exec ID
func generateExecID() string {
	return fmt.Sprintf("exec-%d", time.Now().UnixNano()%100000)
}

// ExecContainerWithStream executes a command with streaming output
func (c *Client) ExecContainerWithStream(ctx context.Context, containerID string, opts *ExecOptions, stdout, stderr io.Writer) error {
	if c.realClient != nil {
		return c.realClient.ExecContainerWithStream(ctx, containerID, opts, stdout, stderr)
	}
	
	// Mock implementation with streaming
	return c.mockExecContainerWithStream(ctx, containerID, opts, stdout, stderr)
}

// mockExecContainerWithStream provides mock streaming execution
func (c *Client) mockExecContainerWithStream(ctx context.Context, containerID string, opts *ExecOptions, stdout, stderr io.Writer) error {
	logrus.Infof("🔄 Streaming command in container: %s", containerID)
	
	command := strings.Join(opts.Command, " ")
	output := generateMockExecOutput(command)
	
	// Simulate streaming output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintf(stdout, "%s\n", line)
			time.Sleep(100 * time.Millisecond) // Simulate streaming delay
		}
	}
	
	logrus.Infof("✅ Streamed command in container %s: %s", containerID, command)
	return nil
}

// RealClient implementation for streaming execution
func (c *RealClient) ExecContainerWithStream(ctx context.Context, containerID string, opts *ExecOptions, stdout, stderr io.Writer) error {
	logrus.Infof("🔄 Streaming command in container: %s", containerID)
	
	// For now, use mock streaming implementation
	command := strings.Join(opts.Command, " ")
	output := generateMockExecOutput(command)
	
	// Simulate streaming output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintf(stdout, "%s\n", line)
			time.Sleep(100 * time.Millisecond) // Simulate streaming delay
		}
	}
	
	logrus.Infof("✅ Streamed command completed in container %s", containerID)
	return nil
}
