package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/shareruntime"
)

// AutoSleepConfigRequest represents a request to set auto-sleep configuration
type AutoSleepConfigRequest struct {
	Enabled          bool   `json:"enabled"`
	IdleTimeout      int    `json:"idleTimeout"`      // Minutes
	SleepTimeout     int    `json:"sleepTimeout"`     // Minutes
	CreateSnapshot   bool   `json:"createSnapshot"`
	AutoWakeOnAccess bool   `json:"autoWakeOnAccess"`
}

// setAutoSleepConfig handles PUT /api/shareruntime/workspaces/:id/autosleep/config
func (s *Server) setAutoSleepConfig(c *gin.Context) {
	workspaceID := c.Param("id")

	var req AutoSleepConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	config := shareruntime.AutoSleepConfig{
		Enabled:          req.Enabled,
		IdleTimeout:      time.Duration(req.IdleTimeout) * time.Minute,
		SleepTimeout:     time.Duration(req.SleepTimeout) * time.Minute,
		CreateSnapshot:   req.CreateSnapshot,
		AutoWakeOnAccess: req.AutoWakeOnAccess,
	}

	s.autoSleepManager.SetConfig(workspaceID, config)

	c.JSON(http.StatusOK, gin.H{
		"message": "Auto-sleep configuration updated",
		"config": gin.H{
			"enabled":          config.Enabled,
			"idleTimeout":      int(config.IdleTimeout.Minutes()),
			"sleepTimeout":     int(config.SleepTimeout.Minutes()),
			"createSnapshot":   config.CreateSnapshot,
			"autoWakeOnAccess": config.AutoWakeOnAccess,
		},
	})
}

// getAutoSleepConfig handles GET /api/shareruntime/workspaces/:id/autosleep/config
func (s *Server) getAutoSleepConfig(c *gin.Context) {
	workspaceID := c.Param("id")

	config := s.autoSleepManager.GetConfig(workspaceID)

	c.JSON(http.StatusOK, gin.H{
		"config": gin.H{
			"enabled":          config.Enabled,
			"idleTimeout":      int(config.IdleTimeout.Minutes()),
			"sleepTimeout":     int(config.SleepTimeout.Minutes()),
			"createSnapshot":   config.CreateSnapshot,
			"autoWakeOnAccess": config.AutoWakeOnAccess,
		},
	})
}

// recordWorkspaceActivity handles POST /api/shareruntime/workspaces/:id/activity
func (s *Server) recordWorkspaceActivity(c *gin.Context) {
	workspaceID := c.Param("id")

	s.autoSleepManager.RecordActivity(workspaceID)

	c.JSON(http.StatusOK, gin.H{"message": "Activity recorded"})
}

// getIdleWorkspaces handles GET /api/shareruntime/autosleep/idle
func (s *Server) getIdleWorkspaces(c *gin.Context) {
	activities := s.autoSleepManager.GetIdleWorkspaces()

	var result []gin.H
	for _, activity := range activities {
		result = append(result, gin.H{
			"workspaceId":   activity.WorkspaceID,
			"workspaceName": activity.WorkspaceName,
			"lastActivity":  activity.LastActivity,
			"idleDuration":  activity.IdleDuration.String(),
			"status":        activity.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"workspaces": result,
		"count":      len(result),
	})
}

// wakeWorkspace handles POST /api/shareruntime/workspaces/:id/wake
func (s *Server) wakeWorkspace(c *gin.Context) {
	workspaceID := c.Param("id")

	err := s.autoSleepManager.WakeWorkspace(workspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to wake workspace",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace woken up"})
}

