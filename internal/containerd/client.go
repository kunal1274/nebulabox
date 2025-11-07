package containerd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Client wraps containerd client with NebulaBox-specific functionality
type Client struct {
	ctx        context.Context
	realClient *RealClient
	mockMode   bool
}

// Container represents a NebulaBox container
type Container struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Created time.Time
}

// NewClient creates a new containerd client
func NewClient() (*Client, error) {
	// Check if we should use real containerd or mock mode
	useReal := os.Getenv("NEBULABOX_REAL_CONTAINERD") == "true"
	
	if useReal {
		logrus.Info("🔧 Initializing NebulaBox containerd client (real mode)")
		realClient, err := NewRealClient("nebulabox")
		if err != nil {
			logrus.Warnf("Failed to connect to real containerd, falling back to mock mode: %v", err)
			return &Client{
				ctx:      context.Background(),
				mockMode: true,
			}, nil
		}
		
		return &Client{
			ctx:        context.Background(),
			realClient: realClient,
			mockMode:   false,
		}, nil
	}
	
	logrus.Info("🔧 Initializing NebulaBox containerd client (mock mode)")
	return &Client{
		ctx:      context.Background(),
		mockMode: true,
	}, nil
}

// IsMock returns true if running in mock mode
func (c *Client) IsMock() bool {
	return c.mockMode
}

// PullImage pulls an image from registry
func (c *Client) PullImage(ctx context.Context, image string) error {
	if c.realClient != nil {
		return c.realClient.PullImage(ctx, image)
	}
	
	// Mock implementation
	logrus.Infof("🔄 Pulling image: %s", image)
	time.Sleep(2 * time.Second)
	logrus.Infof("✅ Successfully pulled image: %s", image)
	return nil
}

// CreateContainer creates a new container
func (c *Client) CreateContainer(ctx context.Context, image, name string, opts *ContainerOptions) (*Container, error) {
	if c.realClient != nil {
		return c.realClient.CreateContainer(ctx, image, name, opts)
	}
	
	// Mock implementation
	container := &Container{
		ID:      generateMockID(),
		Name:    name,
		Image:   image,
		Status:  "created",
		Created: time.Now(),
	}
	
	logrus.Infof("🔄 Creating container: %s (%s)", name, container.ID)
	time.Sleep(500 * time.Millisecond)
	logrus.Infof("✅ Container created: %s (%s)", name, container.ID)
	return container, nil
}

// StartContainer starts a container
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	if c.realClient != nil {
		return c.realClient.StartContainer(ctx, containerID)
	}
	
	// Mock implementation
	logrus.Infof("🔄 Starting container: %s", containerID)
	time.Sleep(1 * time.Second)
	logrus.Infof("✅ Container started: %s", containerID)
	return nil
}

// StopContainer stops a container
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	if c.realClient != nil {
		return c.realClient.StopContainer(ctx, containerID)
	}
	
	// Mock implementation
	logrus.Infof("🔄 Stopping container: %s", containerID)
	time.Sleep(300 * time.Millisecond)
	logrus.Infof("✅ Container stopped: %s", containerID)
	return nil
}

// ListContainers lists all containers
func (c *Client) ListContainers(ctx context.Context) ([]*Container, error) {
	if c.realClient != nil {
		return c.realClient.ListContainers(ctx)
	}
	
	// Mock implementation
	logrus.Info("🔄 Listing containers")
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
	if c.realClient != nil {
		return c.realClient.GetContainerLogs(ctx, containerID)
	}
	
	// Mock implementation
	logrus.Infof("🔄 Getting logs for container: %s", containerID)
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
	if c.realClient != nil {
		return c.realClient.Close()
	}
	return nil
}

// ContainerOptions represents options for container creation
type ContainerOptions struct {
	Name        string
	Ports       map[string]string
	Environment map[string]string
	Volumes     map[string]string
	Detach      bool
    HealthCheck *HealthCheckOptions
    Network     string
}

// HealthCheckOptions defines container health probe configuration
type HealthCheckOptions struct {
    Type            string   // "http", "tcp", "cmd"
    HTTPPath        string
    HTTPPort        string
    TCPPort         string
    Command         []string
    IntervalSeconds int
    TimeoutSeconds  int
    Retries         int
    StartPeriodSec  int
}

// ContainerHealth represents current health state
type ContainerHealth struct {
    Status      string    // healthy, unhealthy, starting, unknown
    LastChecked time.Time
    FailingStreak int
    Message     string
}

// GetContainerHealth returns container health
func (c *Client) GetContainerHealth(ctx context.Context, containerID string) (*ContainerHealth, error) {
    if c.realClient != nil {
        return c.realClient.GetContainerHealth(ctx, containerID)
    }
    // Mock: containers created <30s ago are "starting", else "healthy"
    now := time.Now()
    status := "healthy"
    msg := "probe passed"
    // Without persistence, just alternate by id hash
    if len(containerID) > 0 && (containerID[len(containerID)-1]%2) == 1 {
        status = "starting"
        msg = "initializing"
    }
    return &ContainerHealth{
        Status:       status,
        LastChecked:  now,
        FailingStreak: 0,
        Message:      msg,
    }, nil
}

// generateMockID generates a mock container ID
func generateMockID() string {
	return fmt.Sprintf("mock-%d", time.Now().UnixNano()%10000)
}