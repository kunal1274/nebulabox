package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// FilesystemManager manages container filesystems
type FilesystemManager struct {
	storagePath string
	rootfs      map[string]string // containerID -> rootfs path
}

// NewFilesystemManager creates a new filesystem manager
func NewFilesystemManager(storagePath string) (*FilesystemManager, error) {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &FilesystemManager{
		storagePath: storagePath,
		rootfs:      make(map[string]string),
	}, nil
}

// CreateRootfs creates a root filesystem for a container
func (fm *FilesystemManager) CreateRootfs(containerID string, image *Image) (string, error) {
	rootfsPath := filepath.Join(fm.storagePath, "containers", containerID, "rootfs")

	// Create rootfs directory
	if err := os.MkdirAll(rootfsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create rootfs directory: %w", err)
	}

	// For now, we'll create a basic structure
	// In a full implementation, we'd:
	// 1. Extract image layers
	// 2. Apply layers using overlay filesystem
	// 3. Setup bind mounts for volumes

	// Create essential directories
	dirs := []string{"bin", "etc", "usr", "var", "tmp", "proc", "sys", "dev"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(rootfsPath, dir), 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	fm.rootfs[containerID] = rootfsPath
	return rootfsPath, nil
}

// RemoveRootfs removes a container's root filesystem
func (fm *FilesystemManager) RemoveRootfs(containerID string) error {
	rootfsPath, exists := fm.rootfs[containerID]
	if !exists {
		return fmt.Errorf("rootfs for container %s not found", containerID)
	}

	// Remove the rootfs directory
	if err := os.RemoveAll(rootfsPath); err != nil {
		return fmt.Errorf("failed to remove rootfs: %w", err)
	}

	delete(fm.rootfs, containerID)
	return nil
}

// SetupOverlay creates an overlay filesystem
func (fm *FilesystemManager) SetupOverlay(containerID string, lowerDirs []string, upperDir, workDir string) error {
	_ = filepath.Join(fm.storagePath, "containers", containerID, "rootfs")

	// Create overlay directories
	if err := os.MkdirAll(upperDir, 0755); err != nil {
		return fmt.Errorf("failed to create upper dir: %w", err)
	}
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return fmt.Errorf("failed to create work dir: %w", err)
	}

	// Mount overlay filesystem
	// This requires root privileges and uses syscall.Mount
	// For POC, we'll use a simplified approach
	_ = fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", 
		joinPaths(lowerDirs), upperDir, workDir)

	// Note: Actual mount requires root
	// syscall.Mount("overlay", rootfsPath, "overlay", 0, options)
	
	// For now, just create the structure
	return nil
}

// SetupBindMount creates a bind mount
func (fm *FilesystemManager) SetupBindMount(source, target string, readonly bool) error {
	// Create target directory if it doesn't exist
	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("failed to create mount target: %w", err)
	}

	// Setup bind mount
	// syscall.Mount(source, target, "", syscall.MS_BIND, "")
	// if readonly {
	//     syscall.Mount("", target, "", syscall.MS_REMOUNT|syscall.MS_BIND|syscall.MS_RDONLY, "")
	// }

	// For POC, we'll track bind mounts but actual mounting requires root
	return nil
}

// SetupVolume creates a volume
func (fm *FilesystemManager) SetupVolume(volumeID string) (string, error) {
	volumePath := filepath.Join(fm.storagePath, "volumes", volumeID)
	if err := os.MkdirAll(volumePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create volume: %w", err)
	}
	return volumePath, nil
}

// RemoveVolume removes a volume
func (fm *FilesystemManager) RemoveVolume(volumeID string) error {
	volumePath := filepath.Join(fm.storagePath, "volumes", volumeID)
	return os.RemoveAll(volumePath)
}

// GetRootfsPath returns the rootfs path for a container
func (fm *FilesystemManager) GetRootfsPath(containerID string) (string, bool) {
	path, exists := fm.rootfs[containerID]
	return path, exists
}

// ExtractImageLayers extracts layers from an image
func (fm *FilesystemManager) ExtractImageLayers(image *Image, targetDir string) error {
	// Extract each layer
	for _, layer := range image.Layers {
		layerPath := filepath.Join(targetDir, layer.ID)
		if err := os.MkdirAll(layerPath, 0755); err != nil {
			return fmt.Errorf("failed to create layer directory: %w", err)
		}

		// In a full implementation, we'd:
		// 1. Download layer if needed
		// 2. Extract tar archive
		// 3. Apply layer to filesystem
	}

	return nil
}

// Chroot changes root to a directory
func Chroot(path string) error {
	return syscall.Chroot(path)
}

// PivotRoot performs a pivot root operation
func PivotRoot(newroot, putold string) error {
	// syscall.PivotRoot(newroot, putold)
	// This requires root privileges
	return nil
}

// joinPaths joins multiple paths with colon separator (for overlay)
func joinPaths(paths []string) string {
	result := ""
	for i, path := range paths {
		if i > 0 {
			result += ":"
		}
		result += path
	}
	return result
}

