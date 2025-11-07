package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/shareruntime"
)

// RecordFileChangeRequest represents a request to record a file change
type RecordFileChangeRequest struct {
	ContainerID string            `json:"containerId" binding:"required"`
	Path        string            `json:"path" binding:"required"`
	ChangeType  string            `json:"changeType" binding:"required"` // created, modified, deleted, renamed
	OldPath     string            `json:"oldPath,omitempty"`            // For rename operations
	IsDirectory bool              `json:"isDirectory"`
	Size        int64             `json:"size,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// recordFileChange handles POST /api/shareruntime/workspaces/:id/filesync/changes
func (s *Server) recordFileChange(c *gin.Context) {
	workspaceID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	var req RecordFileChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate change type
	validChangeTypes := map[string]bool{
		"created":  true,
		"modified": true,
		"deleted":   true,
		"renamed":   true,
	}
	if !validChangeTypes[req.ChangeType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid change type. Must be one of: created, modified, deleted, renamed",
		})
		return
	}

	err := s.fileSyncManager.RecordFileChange(
		workspaceID,
		req.ContainerID,
		req.Path,
		shareruntime.FileChangeType(req.ChangeType),
		userID,
		req.IsDirectory,
		req.Size,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to record file change",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File change recorded"})
}

// getFileChanges handles GET /api/shareruntime/workspaces/:id/filesync/changes
func (s *Server) getFileChanges(c *gin.Context) {
	workspaceID := c.Param("id")
	containerID := c.Query("containerId")
	
	sinceStr := c.Query("since")
	var since time.Time
	if sinceStr != "" {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid since timestamp format. Use RFC3339.",
			})
			return
		}
	} else {
		// Default to 1 hour ago
		since = time.Now().Add(-1 * time.Hour)
	}

	changes, err := s.fileSyncManager.GetFileChangesSince(workspaceID, containerID, since)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get file changes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"changes": changes,
		"count":   len(changes),
		"since":   since,
	})
}

// subscribeToFileChanges handles GET /api/shareruntime/workspaces/:id/filesync/subscribe
func (s *Server) subscribeToFileChanges(c *gin.Context) {
	workspaceID := c.Param("id")
	containerID := c.Query("containerId")

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	subID, err := s.fileSyncManager.SubscribeToFileChanges(workspaceID, containerID, func(change *shareruntime.FileChange) {
		c.SSEvent("file_change", change)
		c.Writer.Flush()
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to subscribe",
			"details": err.Error(),
		})
		return
	}

	// Cleanup on disconnect
	defer s.fileSyncManager.UnsubscribeFromFileChanges(subID)

	// Keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.SSEvent("ping", gin.H{"time": time.Now()})
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}

// syncFile handles POST /api/shareruntime/workspaces/:id/filesync/sync
func (s *Server) syncFile(c *gin.Context) {
	workspaceID := c.Param("id")
	
	var req struct {
		ContainerID string `json:"containerId" binding:"required"`
		FilePath    string `json:"filePath" binding:"required"`
		TargetPath  string `json:"targetPath" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err := s.fileSyncManager.SyncFile(workspaceID, req.ContainerID, req.FilePath, req.TargetPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to sync file",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File synced successfully"})
}

// getFileHash handles GET /api/shareruntime/filesync/hash
func (s *Server) getFileHash(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "path parameter is required",
		})
		return
	}

	hash, err := s.fileSyncManager.GetFileHash(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to calculate file hash",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": filePath,
		"hash": hash,
	})
}

