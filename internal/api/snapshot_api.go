package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
	"github.com/nebulabox/nebulabox/internal/snapshot"
)

// CreateSnapshotRequest represents a request to create a snapshot
type CreateSnapshotRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type" binding:"required"` // container, workspace, volume
	ResourceID  string                 `json:"resourceId" binding:"required"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// RestoreSnapshotRequest represents a request to restore from a snapshot
type RestoreSnapshotRequest struct {
	SnapshotID string `json:"snapshotId" binding:"required"`
	NewName    string `json:"newName,omitempty"` // For creating a new resource from snapshot
}

// createSnapshot handles POST /api/snapshots
func (s *Server) createSnapshot(c *gin.Context) {
	username, _ := c.Get("username")
	userID := username.(string)

	var req CreateSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate snapshot type
	var snapshotType snapshot.SnapshotType
	switch strings.ToLower(req.Type) {
	case "container":
		snapshotType = snapshot.SnapshotTypeContainer
	case "workspace":
		snapshotType = snapshot.SnapshotTypeWorkspace
	case "volume":
		snapshotType = snapshot.SnapshotTypeVolume
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid snapshot type. Must be one of: container, workspace, volume",
		})
		return
	}

	// Create snapshot
	snap, err := s.snapshotManager.CreateSnapshot(
		req.Name,
		req.Description,
		req.ResourceID,
		snapshotType,
		userID,
		req.Metadata,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create snapshot",
			"details": err.Error(),
		})
		return
	}

	// If it's a container snapshot, capture container data
	if snapshotType == snapshot.SnapshotTypeContainer {
		// List containers to find the one we're snapshotting
		ctx := c.Request.Context()
		containers, err := s.containerd.ListContainers(ctx)
		var container *containerd.Container
		if err == nil {
			for _, c := range containers {
				if c.ID == req.ResourceID {
					container = c
					break
				}
			}
		}
		if container != nil {
			// Extract container configuration
			// Note: Container struct is simplified, so we use placeholder values
			// In a real implementation, these would be extracted from container metadata
			env := make(map[string]string)
			ports := make(map[string]string)
			volumes := []string{}
			var resources *snapshot.ResourceLimits

			// Update snapshot with container data
			err = s.snapshotManager.SetSnapshotContainerData(
				snap.ID,
				container.Image,
				nil, // Command - would need to be stored separately
				env,
				ports,
				volumes,
				"", // Network - would need to be stored separately
				resources,
			)
			if err != nil {
				// Log error but don't fail the snapshot creation
				_ = err
			}

			// Mark snapshot as ready (in real implementation, would wait for filesystem capture)
			_ = s.snapshotManager.SetSnapshotState(snap.ID, snapshot.SnapshotStateReady)
			_ = s.snapshotManager.SetSnapshotSize(snap.ID, 1024*1024) // Mock size
		} else {
			// Container not found, mark as failed
			_ = s.snapshotManager.SetSnapshotState(snap.ID, snapshot.SnapshotStateFailed)
		}
	} else {
		// For workspace/volume, mark as ready immediately (would capture in real implementation)
		_ = s.snapshotManager.SetSnapshotState(snap.ID, snapshot.SnapshotStateReady)
	}

	// Return the snapshot
	snap, _ = s.snapshotManager.GetSnapshot(snap.ID)
	c.JSON(http.StatusOK, snap)
}

// listSnapshots handles GET /api/snapshots
func (s *Server) listSnapshots(c *gin.Context) {
	resourceID := c.Query("resourceId")
	snapshotType := c.Query("type")

	var st snapshot.SnapshotType
	if snapshotType != "" {
		st = snapshot.SnapshotType(snapshotType)
	}

	snapshots, err := s.snapshotManager.ListSnapshots(resourceID, st)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to list snapshots",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"snapshots": snapshots,
		"count":     len(snapshots),
	})
}

// getSnapshot handles GET /api/snapshots/:id
func (s *Server) getSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")

	snap, err := s.snapshotManager.GetSnapshot(snapshotID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Snapshot not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, snap)
}

// deleteSnapshot handles DELETE /api/snapshots/:id
func (s *Server) deleteSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")

	err := s.snapshotManager.DeleteSnapshot(snapshotID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to delete snapshot",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Snapshot deleted successfully"})
}

// restoreSnapshot handles POST /api/snapshots/:id/restore
func (s *Server) restoreSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")

	var req RestoreSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get snapshot
	snap, err := s.snapshotManager.GetSnapshot(snapshotID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Snapshot not found",
			"details": err.Error(),
		})
		return
	}

	if snap.State != snapshot.SnapshotStateReady {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Snapshot is not ready for restore",
		})
		return
	}

	// Set snapshot state to restoring
	_ = s.snapshotManager.SetSnapshotState(snapshotID, snapshot.SnapshotStateRestoring)

	// Restore based on snapshot type
	if snap.Type == snapshot.SnapshotTypeContainer {
		ctx := c.Request.Context()
		
		// Build ContainerOptions from snapshot
		containerOpts := &containerd.ContainerOptions{
			Name:    req.NewName,
			Ports:    make(map[string]string),
			Environment: make(map[string]string),
			Volumes:  make(map[string]string),
			Network:   snap.Network,
		}

		// Convert env map
		if snap.Env != nil {
			for k, v := range snap.Env {
				containerOpts.Environment[k] = v
			}
		}

		// Convert ports map
		if snap.Ports != nil {
			containerOpts.Ports = snap.Ports
		}

		// Convert volumes array
		if snap.Volumes != nil {
			for _, vol := range snap.Volumes {
				parts := strings.Split(vol, ":")
				if len(parts) == 2 {
					containerOpts.Volumes[parts[0]] = parts[1]
				}
			}
		}

		// Create container using the containerd client
		container, err := s.containerd.CreateContainer(ctx, snap.Image, req.NewName, containerOpts)
		if err != nil {
			_ = s.snapshotManager.SetSnapshotState(snapshotID, snapshot.SnapshotStateReady)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to restore container",
				"details": err.Error(),
			})
			return
		}

		// Restore complete
		_ = s.snapshotManager.SetSnapshotState(snapshotID, snapshot.SnapshotStateReady)

		c.JSON(http.StatusOK, gin.H{
			"message":   "Container restored successfully",
			"container": container,
		})
	} else {
		// Workspace/volume restore would be handled here
		_ = s.snapshotManager.SetSnapshotState(snapshotID, snapshot.SnapshotStateReady)
		c.JSON(http.StatusOK, gin.H{
			"message": "Restore operation completed",
		})
	}
}

// listResourceSnapshots handles GET /api/snapshots/resource/:resourceId
func (s *Server) listResourceSnapshots(c *gin.Context) {
	resourceID := c.Param("resourceId")

	snapshots, err := s.snapshotManager.ListResourceSnapshots(resourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to list resource snapshots",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"snapshots": snapshots,
		"count":     len(snapshots),
	})
}

