package containerd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListContainers(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient()
	assert.NoError(t, err)
	
	containers, err := client.ListContainers(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, containers)
}

func TestPullImage(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient()
	assert.NoError(t, err)
	
	// This will use mock mode, so it should succeed
	err = client.PullImage(ctx, "alpine:latest")
	assert.NoError(t, err)
}

func TestCreateContainer(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient()
	assert.NoError(t, err)
	
	opts := &ContainerOptions{
		Name: "test-container",
		Ports: map[string]string{"80": "8080"},
		Environment: map[string]string{"ENV": "test"},
		Detach: true,
	}
	
	container, err := client.CreateContainer(ctx, "alpine:latest", "test-container", opts)
	assert.NoError(t, err)
	assert.NotNil(t, container)
	assert.Equal(t, "test-container", container.Name)
}

