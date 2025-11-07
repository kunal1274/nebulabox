package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRuntime(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "storage")
	statePath := filepath.Join(tmpDir, "state")

	// Create runtime
	runtime, err := NewRuntime(storagePath, statePath)
	if err != nil {
		t.Fatalf("Failed to create runtime: %v", err)
	}

	if runtime == nil {
		t.Fatal("Runtime is nil")
	}

	// Test Info
	info, err := runtime.Info(context.Background())
	if err != nil {
		t.Fatalf("Failed to get runtime info: %v", err)
	}

	if info.Name != "NebulaBox Runtime" {
		t.Errorf("Expected runtime name 'NebulaBox Runtime', got '%s'", info.Name)
	}

	// Test Version
	version, err := runtime.Version(context.Background())
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if version == "" {
		t.Error("Version should not be empty")
	}
}

func TestRuntimeStoragePaths(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "storage")
	statePath := filepath.Join(tmpDir, "state")

	_, err := NewRuntime(storagePath, statePath)
	if err != nil {
		t.Fatalf("Failed to create runtime: %v", err)
	}

	// Verify storage directories were created
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		t.Errorf("Storage path was not created: %s", storagePath)
	}
}

func TestListContainersEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "storage")
	statePath := filepath.Join(tmpDir, "state")

	runtime, err := NewRuntime(storagePath, statePath)
	if err != nil {
		t.Fatalf("Failed to create runtime: %v", err)
	}

	containers, err := runtime.ListContainers(context.Background())
	if err != nil {
		t.Fatalf("Failed to list containers: %v", err)
	}

	if len(containers) != 0 {
		t.Errorf("Expected 0 containers, got %d", len(containers))
	}
}

func TestListImagesEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "storage")
	statePath := filepath.Join(tmpDir, "state")

	runtime, err := NewRuntime(storagePath, statePath)
	if err != nil {
		t.Fatalf("Failed to create runtime: %v", err)
	}

	images, err := runtime.ListImages(context.Background())
	if err != nil {
		t.Fatalf("Failed to list images: %v", err)
	}

	if len(images) != 0 {
		t.Errorf("Expected 0 images, got %d", len(images))
	}
}

func TestListGroupsEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "storage")
	statePath := filepath.Join(tmpDir, "state")

	runtime, err := NewRuntime(storagePath, statePath)
	if err != nil {
		t.Fatalf("Failed to create runtime: %v", err)
	}

	groups, err := runtime.ListGroups(context.Background())
	if err != nil {
		t.Fatalf("Failed to list groups: %v", err)
	}

	if len(groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(groups))
	}
}

