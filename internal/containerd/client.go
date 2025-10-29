package containerd

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Client wraps containerd client with NebulaBox-specific functionality
type Client struct {
	ctx context.Context
}

// Container represents a NebulaBox container
type Container struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Created time.Time
}

// NewClient creates a new containerd client (mock implementation)
func NewClient() (*Client, error) {
	logrus.Info("🔧 Initializing NebulaBox containerd client (mock mode)")
	return &Client{
		ctx: context.Background(),
	}, nil
}

// IsMock returns true if running in mock mode
func (c *Client) IsMock() bool {
	return true // Always mock for now
}

// PullImage pulls an image from registry
func (c *Client) PullImage(ctx context.Context, image string) error {
	logrus.Infof("🔄 Pulling image: %s", image)
	
	// Simulate download time based on image size
	time.Sleep(2 * time.Second)
	
	logrus.Infof("✅ Successfully pulled image: %s", image)
	return nil
}

// CreateContainer creates a new container
func (c *Client) CreateContainer(ctx context.Context, image, name string, opts *ContainerOptions) (*Container, error) {
	container := &Container{
		ID:      generateMockID(),
		Name:    name,
		Image:   image,
		Status:  "created",
		Created: time.Now(),
	}
	
	logrus.Infof("🔄 Creating container: %s (%s)", name, container.ID)
	time.Sleep(500 * time.Millisecond) // Simulate creation time
	
	logrus.Infof("✅ Container created: %s (%s)", name, container.ID)
	return container, nil
}

// StartContainer starts a container
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	logrus.Infof("🔄 Starting container: %s", containerID)
	time.Sleep(1 * time.Second) // Simulate start time
	
	logrus.Infof("✅ Container started: %s", containerID)
	return nil
}

// StopContainer stops a container
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	logrus.Infof("🔄 Stopping container: %s", containerID)
	time.Sleep(300 * time.Millisecond) // Simulate stop time
	
	logrus.Infof("✅ Container stopped: %s", containerID)
	return nil
}

// ListContainers lists all containers
func (c *Client) ListContainers(ctx context.Context) ([]*Container, error) {
	logrus.Info("🔄 Listing containers")
	
	// Return mock containers for demonstration
	containers := []*Container{
		{
			ID:      "mock-001",
			Name:    "web-server",
			Image:   "nginx:latest",
			Status:  "running",
			Created: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:      "mock-002",
			Name:    "db-server",
			Image:   "postgres:13",
			Status:  "running",
			Created: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:      "mock-003",
			Name:    "api-server",
			Image:   "node:18",
			Status:  "stopped",
			Created: time.Now().Add(-30 * time.Minute),
		},
	}
	
	logrus.Infof("✅ Found %d containers", len(containers))
	return containers, nil
}

// GetContainerLogs gets container logs
func (c *Client) GetContainerLogs(ctx context.Context, containerID string) ([]string, error) {
	logrus.Infof("🔄 Getting logs for container: %s", containerID)
	
	// Generate realistic mock logs
	logs := []string{
		fmt.Sprintf("[%s] Container %s started", time.Now().Add(-2*time.Hour).Format("2006-01-02 15:04:05"), containerID),
		fmt.Sprintf("[%s] Application listening on port 80", time.Now().Add(-2*time.Hour+1*time.Minute).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("[%s] Health check passed", time.Now().Add(-2*time.Hour+2*time.Minute).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("[%s] Request processed: GET /", time.Now().Add(-1*time.Hour).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("[%s] Request processed: GET /api/status", time.Now().Add(-30*time.Minute).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("[%s] Container %s running normally", time.Now().Format("2006-01-02 15:04:05"), containerID),
	}
	
	logrus.Infof("✅ Retrieved %d log lines for container %s", len(logs), containerID)
	return logs, nil
}

// Close closes the containerd client
func (c *Client) Close() error {
	logrus.Info("🔄 Closing containerd client")
	return nil
}

// ContainerOptions represents options for container creation
type ContainerOptions struct {
	Name        string
	Ports       map[string]string
	Environment map[string]string
	Volumes     map[string]string
	Detach      bool
}

// generateMockID generates a mock container ID
func generateMockID() string {
	return fmt.Sprintf("mock-%d", time.Now().UnixNano()%10000)
}