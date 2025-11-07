package runtime

import (
	"fmt"
	"path/filepath"
)

// OverlayManager manages OverlayFS mounts for container filesystems
type OverlayManager struct {
	basePath string
}

// NewOverlayManager creates a new OverlayFS manager
func NewOverlayManager(basePath string) *OverlayManager {
	if basePath == "" {
		basePath = "/var/lib/nebula-runtime"
	}

	return &OverlayManager{
		basePath: basePath,
	}
}

// CreateOverlay creates an OverlayFS mount for a container
func (om *OverlayManager) CreateOverlay(containerID string, layers []string) (string, error) {
	// Mock implementation - in production this would:
	// 1. Create lowerdir (read-only layers)
	// 2. Create upperdir (writable layer)
	// 3. Create workdir (work directory)
	// 4. Mount OverlayFS using mount(2) with MS_MOVE
	// 5. Return mount point path

	// Structure:
	// - basePath/containers/{id}/lower (symlink to image layers)
	// - basePath/containers/{id}/upper (writable layer)
	// - basePath/containers/{id}/work (work directory)
	// - basePath/containers/{id}/merged (merged mount point)

	mountPoint := filepath.Join(om.basePath, "containers", containerID, "merged")

	// In production, would execute:
	// mount -t overlay overlay \
	//   -o lowerdir=lower1:lower2,upperdir=upper,workdir=work \
	//   merged

	_ = layers // Suppress unused warning
	_ = mountPoint

	return mountPoint, nil
}

// RemoveOverlay removes an OverlayFS mount
func (om *OverlayManager) RemoveOverlay(containerID string) error {
	// Mock implementation - in production would:
	// 1. Unmount the overlay
	// 2. Remove upper, work, and merged directories

	mountPoint := filepath.Join(om.basePath, "containers", containerID, "merged")
	_ = mountPoint

	return nil
}

// CreateLayer creates a new layer in the image
func (om *OverlayManager) CreateLayer(imageID string) (string, error) {
	layerPath := filepath.Join(om.basePath, "images", imageID, "layers", fmt.Sprintf("layer-%d", len([]string{})))
	return layerPath, nil
}

// GetLayerPath returns the path to an image layer
func (om *OverlayManager) GetLayerPath(imageID, layerDigest string) string {
	return filepath.Join(om.basePath, "images", imageID, "layers", layerDigest)
}

// MountImageLayers mounts image layers in preparation for container creation
func (om *OverlayManager) MountImageLayers(containerID string, imageLayers []ImageLayer) ([]string, error) {
	// Returns paths to lower directories
	lowerDirs := make([]string, len(imageLayers))
	for i, layer := range imageLayers {
		lowerDirs[i] = om.GetLayerPath(imageLayers[0].Digest, layer.Digest)
	}
	return lowerDirs, nil
}

