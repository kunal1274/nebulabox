package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/shareruntime"
)

// getChangesSince handles GET /api/shareruntime/workspaces/:id/sync/changes
func (s *Server) getChangesSince(c *gin.Context) {
	workspaceID := c.Param("id")
	
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

	changes, err := s.syncManager.GetChangesSince(workspaceID, since)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get changes",
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

// getLatestChangeID handles GET /api/shareruntime/workspaces/:id/sync/latest
func (s *Server) getLatestChangeID(c *gin.Context) {
	workspaceID := c.Param("id")
	
	changeID, err := s.syncManager.GetLatestChangeID(workspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get latest change ID",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"changeId": changeID,
	})
}

// subscribeToChanges handles GET /api/shareruntime/workspaces/:id/sync/subscribe
func (s *Server) subscribeToChanges(c *gin.Context) {
	workspaceID := c.Param("id")

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	subID, err := s.syncManager.SubscribeToChanges(workspaceID, func(event *shareruntime.SyncEvent) {
		c.SSEvent("change", event)
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
	defer s.syncManager.UnsubscribeFromChanges(subID)

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

// syncWorkspace handles POST /api/shareruntime/workspaces/:id/sync/apply
func (s *Server) syncWorkspace(c *gin.Context) {
	workspaceID := c.Param("id")
	
	var req struct {
		Changes []*shareruntime.SyncEvent `json:"changes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Apply changes (for now, just verify they can be processed)
	// In a real implementation, this would apply the changes to local state
	applied := 0
	errors := []string{}

	for _, change := range req.Changes {
		if change.WorkspaceID != workspaceID {
			errors = append(errors, fmt.Sprintf("Change %s belongs to different workspace", change.ID))
			continue
		}

		// Here we would apply the change to local state
		// For now, we just validate
		applied++
	}

	c.JSON(http.StatusOK, gin.H{
		"applied": applied,
		"total":   len(req.Changes),
		"errors":  errors,
	})
}

