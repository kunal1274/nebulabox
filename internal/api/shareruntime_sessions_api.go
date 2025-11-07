package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateSessionActivityRequest represents a request to update session activity
type UpdateSessionActivityRequest struct {
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateSessionStateRequest represents a request to update session state
type UpdateSessionStateRequest struct {
	State    string            `json:"state" binding:"required"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// updateSessionActivity handles POST /api/shareruntime/sessions/:id/activity
func (s *Server) updateSessionActivity(c *gin.Context) {
	sessionID := c.Param("id")
	var req UpdateSessionActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Get session to find workspace
	session, exists := s.sharedRuntimeManager.GetSession(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Session not found",
		})
		return
	}

	// Update activity via multiplexer
	mux := s.sharedRuntimeManager.GetSessionManager().GetMultiplexer(session.WorkspaceID)
	if err := mux.UpdateActivity(sessionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update activity",
			"details": err.Error(),
		})
		return
	}

	// Update metadata if provided
	if req.Metadata != nil {
		state, _ := mux.GetSessionState(sessionID)
		if state != nil {
			for k, v := range req.Metadata {
				state.Metadata[k] = v
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity updated"})
}

// updateSessionState handles PATCH /api/shareruntime/sessions/:id/state
func (s *Server) updateSessionState(c *gin.Context) {
	sessionID := c.Param("id")
	var req UpdateSessionStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Get session to find workspace
	session, exists := s.sharedRuntimeManager.GetSession(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Session not found",
		})
		return
	}

	// Update state via multiplexer
	mux := s.sharedRuntimeManager.GetSessionManager().GetMultiplexer(session.WorkspaceID)
	if err := mux.UpdateSessionState(sessionID, req.State); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update state",
			"details": err.Error(),
		})
		return
	}

	// Update metadata if provided
	if req.Metadata != nil {
		state, _ := mux.GetSessionState(sessionID)
		if state != nil {
			for k, v := range req.Metadata {
				state.Metadata[k] = v
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "State updated"})
}

// getSessionState handles GET /api/shareruntime/sessions/:id/state
func (s *Server) getSessionState(c *gin.Context) {
	sessionID := c.Param("id")

	// Get session to find workspace
	session, exists := s.sharedRuntimeManager.GetSession(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Session not found",
		})
		return
	}

	// Get state from multiplexer
	mux := s.sharedRuntimeManager.GetSessionManager().GetMultiplexer(session.WorkspaceID)
	state, exists := mux.GetSessionState(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Session state not found",
		})
		return
	}

	c.JSON(http.StatusOK, state)
}

// listActiveSessions handles GET /api/shareruntime/workspaces/:id/active-sessions
func (s *Server) listActiveSessions(c *gin.Context) {
	workspaceID := c.Param("id")

	mux := s.sharedRuntimeManager.GetSessionManager().GetMultiplexer(workspaceID)
	sessions := mux.ListActiveSessions()

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// subscribeSessionActivity handles GET /api/shareruntime/workspaces/:id/activity-stream
func (s *Server) subscribeSessionActivity(c *gin.Context) {
	workspaceID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	mux := s.sharedRuntimeManager.GetSessionManager().GetMultiplexer(workspaceID)
	activityCh := mux.SubscribeActivity(userID)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Cleanup on disconnect
	defer mux.UnsubscribeActivity(userID)

	// Stream session activity updates
	for {
		select {
		case state, ok := <-activityCh:
			if !ok {
				return
			}
			c.SSEvent("session_activity", state)
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}

