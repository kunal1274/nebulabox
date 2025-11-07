package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ProvisionEphemeralRuntimeRequest represents a request to provision an ephemeral runtime
type ProvisionEphemeralRuntimeRequest struct {
	Name        string   `json:"name" binding:"required"`
	Region      string   `json:"region" binding:"required"`
	InstanceType string  `json:"instanceType" binding:"required"` // small, medium, large
	Image       string   `json:"image" binding:"required"`
	Duration    int      `json:"duration"` // Duration in hours (default: 24)
	Members     []string `json:"members,omitempty"` // User IDs to grant access
}

// provisionEphemeralRuntime handles POST /api/cloud/ephemeral/runtimes
func (s *Server) provisionEphemeralRuntime(c *gin.Context) {
	username, _ := c.Get("username")
	userID := username.(string)

	var req ProvisionEphemeralRuntimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate instance type
	validTypes := map[string]bool{
		"small":  true,
		"medium": true,
		"large":   true,
	}
	if !validTypes[req.InstanceType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid instance type. Must be one of: small, medium, large",
		})
		return
	}

	// Get workspace ID from query or use default
	workspaceID := c.Query("workspaceId")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "workspaceId query parameter is required",
		})
		return
	}

	duration := time.Duration(req.Duration) * time.Hour
	if req.Duration == 0 {
		duration = 24 * time.Hour // Default: 24 hours
	}

	// Add creator to members if not already included
	members := req.Members
	if members == nil {
		members = []string{}
	}
	found := false
	for _, member := range members {
		if member == userID {
			found = true
			break
		}
	}
	if !found {
		members = append(members, userID)
	}

	runtime, err := s.ephemeralRuntimeManager.ProvisionRuntime(
		req.Name,
		workspaceID,
		req.Region,
		req.InstanceType,
		req.Image,
		userID,
		members,
		duration,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to provision runtime",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, runtime)
}

// listEphemeralRuntimes handles GET /api/cloud/ephemeral/runtimes
func (s *Server) listEphemeralRuntimes(c *gin.Context) {
	workspaceID := c.Query("workspaceId")

	runtimes, err := s.ephemeralRuntimeManager.ListRuntimes(workspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to list runtimes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"runtimes": runtimes,
		"count":    len(runtimes),
	})
}

// getEphemeralRuntime handles GET /api/cloud/ephemeral/runtimes/:id
func (s *Server) getEphemeralRuntime(c *gin.Context) {
	runtimeID := c.Param("id")

	runtime, err := s.ephemeralRuntimeManager.GetRuntime(runtimeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Runtime not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, runtime)
}

// updateEphemeralRuntimeActivity handles POST /api/cloud/ephemeral/runtimes/:id/activity
func (s *Server) updateEphemeralRuntimeActivity(c *gin.Context) {
	runtimeID := c.Param("id")

	err := s.ephemeralRuntimeManager.UpdateActivity(runtimeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to update activity",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity updated"})
}

// sleepEphemeralRuntime handles POST /api/cloud/ephemeral/runtimes/:id/sleep
func (s *Server) sleepEphemeralRuntime(c *gin.Context) {
	runtimeID := c.Param("id")

	var req struct {
		SnapshotID string `json:"snapshotId,omitempty"`
	}
	c.ShouldBindJSON(&req)

	// Create snapshot if not provided
	snapshotID := req.SnapshotID
	if snapshotID == "" {
		snapshotID = fmt.Sprintf("auto-snapshot-%s", runtimeID)
	}

	err := s.ephemeralRuntimeManager.SleepRuntime(runtimeID, snapshotID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to sleep runtime",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Runtime put to sleep",
		"snapshotId": snapshotID,
	})
}

// wakeEphemeralRuntime handles POST /api/cloud/ephemeral/runtimes/:id/wake
func (s *Server) wakeEphemeralRuntime(c *gin.Context) {
	runtimeID := c.Param("id")

	err := s.ephemeralRuntimeManager.WakeRuntime(runtimeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to wake runtime",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Runtime woken up"})
}

// terminateEphemeralRuntime handles DELETE /api/cloud/ephemeral/runtimes/:id
func (s *Server) terminateEphemeralRuntime(c *gin.Context) {
	runtimeID := c.Param("id")

	err := s.ephemeralRuntimeManager.TerminateRuntime(runtimeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to terminate runtime",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Runtime termination initiated"})
}

// addEphemeralRuntimeMember handles POST /api/cloud/ephemeral/runtimes/:id/members
func (s *Server) addEphemeralRuntimeMember(c *gin.Context) {
	runtimeID := c.Param("id")

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err := s.ephemeralRuntimeManager.AddMember(runtimeID, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add member",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

// removeEphemeralRuntimeMember handles DELETE /api/cloud/ephemeral/runtimes/:id/members/:userId
func (s *Server) removeEphemeralRuntimeMember(c *gin.Context) {
	runtimeID := c.Param("id")
	userID := c.Param("userId")

	err := s.ephemeralRuntimeManager.RemoveMember(runtimeID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to remove member",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// checkEphemeralRuntimeHealth handles GET /api/cloud/ephemeral/runtimes/:id/health
func (s *Server) checkEphemeralRuntimeHealth(c *gin.Context) {
	runtimeID := c.Param("id")

	runtime, err := s.ephemeralRuntimeManager.GetRuntime(runtimeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Runtime not found",
			"details": err.Error(),
		})
		return
	}

	// Check if expired
	isExpired := time.Now().After(runtime.ExpiresAt)
	idleDuration := time.Since(runtime.LastActivityAt)

	c.JSON(http.StatusOK, gin.H{
		"runtimeId":     runtimeID,
		"status":        runtime.Status,
		"isExpired":     isExpired,
		"idleDuration":  idleDuration.String(),
		"expiresAt":     runtime.ExpiresAt,
		"lastActivity": runtime.LastActivityAt,
	})
}

