package containerd

import (
	"context"
	"testing"
)

// BenchmarkListContainers benchmarks container listing
func BenchmarkListContainers(b *testing.B) {
	ctx := context.Background()
	client, err := NewClient()
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ListContainers(ctx)
	}
}

// BenchmarkPullImage benchmarks image pulling
func BenchmarkPullImage(b *testing.B) {
	ctx := context.Background()
	client, err := NewClient()
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.PullImage(ctx, "alpine:latest")
	}
}

// BenchmarkCreateContainer benchmarks container creation
func BenchmarkCreateContainer(b *testing.B) {
	ctx := context.Background()
	client, err := NewClient()
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}
	
	opts := &ContainerOptions{
		Name:        "bench-container",
		Ports:       map[string]string{"80": "8080"},
		Environment: map[string]string{"ENV": "test"},
		Detach:      true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opts.Name = "bench-container-" + string(rune(i))
		_, _ = client.CreateContainer(ctx, "alpine:latest", opts.Name, opts)
	}
}

