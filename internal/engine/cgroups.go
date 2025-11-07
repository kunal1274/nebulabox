package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CGroupManager manages cgroups for containers
type CGroupManager struct {
	cgroupRoot string
	cgroups    map[string]string // containerID -> cgroup path
}

// NewCGroupManager creates a new cgroup manager
func NewCGroupManager() (*CGroupManager, error) {
	// Detect cgroup version and root path
	cgroupRoot := "/sys/fs/cgroup"
	if _, err := os.Stat(cgroupRoot); os.IsNotExist(err) {
		// Try cgroup v2 path
		cgroupRoot = "/sys/fs/cgroup/unified"
	}

	return &CGroupManager{
		cgroupRoot: cgroupRoot,
		cgroups:    make(map[string]string),
	}, nil
}

// CreateCGroup creates a new cgroup for a container
func (cg *CGroupManager) CreateCGroup(containerID string, limits *ResourceLimits) (string, error) {
	// Create cgroup path
	cgroupPath := filepath.Join(cg.cgroupRoot, "nebulabox", containerID)

	// Create the directory
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create cgroup directory: %w", err)
	}

	// Apply resource limits if provided
	if limits != nil {
		if err := cg.applyLimits(cgroupPath, limits); err != nil {
			return "", fmt.Errorf("failed to apply limits: %w", err)
		}
	}

	cg.cgroups[containerID] = cgroupPath
	return cgroupPath, nil
}

// applyLimits applies resource limits to a cgroup
func (cg *CGroupManager) applyLimits(cgroupPath string, limits *ResourceLimits) error {
	// CPU limits
	if limits.CPUShares > 0 {
		if err := cg.writeCGroupFile(cgroupPath, "cpu.shares", strconv.FormatInt(limits.CPUShares, 10)); err != nil {
			return err
		}
	}

	if limits.CPUQuota > 0 {
		// CPU quota is typically in microseconds
		if err := cg.writeCGroupFile(cgroupPath, "cpu.cfs_quota_us", strconv.FormatInt(limits.CPUQuota, 10)); err != nil {
			return err
		}
	}

	// Memory limits
	if limits.MemoryLimit > 0 {
		if err := cg.writeCGroupFile(cgroupPath, "memory.limit_in_bytes", strconv.FormatInt(limits.MemoryLimit, 10)); err != nil {
			return err
		}
	}

	if limits.MemorySwap > 0 {
		if err := cg.writeCGroupFile(cgroupPath, "memory.memsw.limit_in_bytes", strconv.FormatInt(limits.MemorySwap, 10)); err != nil {
			return err
		}
	}

	// PIDs limit
	if limits.PidsLimit > 0 {
		if err := cg.writeCGroupFile(cgroupPath, "pids.max", strconv.FormatInt(limits.PidsLimit, 10)); err != nil {
			return err
		}
	}

	// Blkio limits
	if limits.BlkioWeight > 0 {
		if err := cg.writeCGroupFile(cgroupPath, "blkio.weight", strconv.FormatUint(uint64(limits.BlkioWeight), 10)); err != nil {
			return err
		}
	}

	return nil
}

// writeCGroupFile writes a value to a cgroup control file
func (cg *CGroupManager) writeCGroupFile(cgroupPath, filename, value string) error {
	filePath := filepath.Join(cgroupPath, filename)
	
	// Check if file exists (cgroup v1 vs v2 have different structures)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Try cgroup v2 path
		filePath = filepath.Join(cgroupPath, strings.Replace(filename, ".", ".", -1))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// File doesn't exist, skip (might not be supported)
			return nil
		}
	}

	return os.WriteFile(filePath, []byte(value), 0644)
}

// AddProcess adds a process to a cgroup
func (cg *CGroupManager) AddProcess(cgroupPath string, pid int) error {
	procsFile := filepath.Join(cgroupPath, "cgroup.procs")
	return os.WriteFile(procsFile, []byte(strconv.Itoa(pid)), 0644)
}

// RemoveCGroup removes a cgroup
func (cg *CGroupManager) RemoveCGroup(cgroupPath string) error {
	// Remove all processes from cgroup first
	procsFile := filepath.Join(cgroupPath, "cgroup.procs")
	if data, err := os.ReadFile(procsFile); err == nil {
		// Move processes to parent cgroup
		parentPath := filepath.Dir(cgroupPath)
		parentProcsFile := filepath.Join(parentPath, "cgroup.procs")
		os.WriteFile(parentProcsFile, data, 0644)
	}

	// Remove the cgroup directory
	if err := os.RemoveAll(cgroupPath); err != nil {
		return fmt.Errorf("failed to remove cgroup: %w", err)
	}

	// Remove from tracking
	for id, path := range cg.cgroups {
		if path == cgroupPath {
			delete(cg.cgroups, id)
			break
		}
	}

	return nil
}

// GetCGroupPath returns the cgroup path for a container
func (cg *CGroupManager) GetCGroupPath(containerID string) (string, bool) {
	path, exists := cg.cgroups[containerID]
	return path, exists
}

// SetCPUQuota sets CPU quota for a cgroup
func (cg *CGroupManager) SetCPUQuota(cgroupPath string, quota int64) error {
	return cg.writeCGroupFile(cgroupPath, "cpu.cfs_quota_us", strconv.FormatInt(quota, 10))
}

// SetMemoryLimit sets memory limit for a cgroup
func (cg *CGroupManager) SetMemoryLimit(cgroupPath string, limit int64) error {
	return cg.writeCGroupFile(cgroupPath, "memory.limit_in_bytes", strconv.FormatInt(limit, 10))
}

// GetMemoryUsage returns current memory usage
func (cg *CGroupManager) GetMemoryUsage(cgroupPath string) (int64, error) {
	usageFile := filepath.Join(cgroupPath, "memory.usage_in_bytes")
	data, err := os.ReadFile(usageFile)
	if err != nil {
		return 0, err
	}

	usage, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0, err
	}

	return usage, nil
}

// GetCPUUsage returns current CPU usage
func (cg *CGroupManager) GetCPUUsage(cgroupPath string) (float64, error) {
	// CPU usage is more complex, requires reading cpuacct stats
	// This is a simplified version
	return 0.0, nil
}

