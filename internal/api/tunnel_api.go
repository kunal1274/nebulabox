package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTunnelRequest represents a request to create a tunnel
type CreateTunnelRequest struct {
	ContainerID   string   `json:"containerId" binding:"required"`
	ContainerPort int      `json:"containerPort" binding:"required"`
	Protocol      string   `json:"protocol" binding:"required"`
	AllowedIPs    []string `json:"allowedIPs,omitempty"`
}

// createTunnel handles POST /api/shareruntime/workspaces/:id/tunnels
func (s *Server) createTunnel(c *gin.Context) {
	workspaceID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	var req CreateTunnelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Verify workspace exists
	_, exists := s.sharedRuntimeManager.GetWorkspace(workspaceID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Workspace not found",
		})
		return
	}

	// Check permission
	if canCreate, err := s.sharedRuntimeManager.CanCreateTunnel(workspaceID, userID); err != nil || !canCreate {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: cannot create tunnels",
		})
		return
	}

	// Create tunnel
	tunnel, err := s.tunnelManager.CreateTunnel(
		workspaceID,
		userID,
		userID, // Use userID as username for now
		req.ContainerID,
		req.ContainerPort,
		req.Protocol,
		req.AllowedIPs,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create tunnel",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tunnel)
}

// listTunnels handles GET /api/shareruntime/workspaces/:id/tunnels
func (s *Server) listTunnels(c *gin.Context) {
	workspaceID := c.Param("id")
	tunnels := s.tunnelManager.ListTunnels(workspaceID)
	c.JSON(http.StatusOK, gin.H{
		"tunnels": tunnels,
		"count":   len(tunnels),
	})
}

// getUserTunnels handles GET /api/shareruntime/tunnels
func (s *Server) getUserTunnels(c *gin.Context) {
	username, _ := c.Get("username")
	userID := username.(string)
	tunnels := s.tunnelManager.ListUserTunnels(userID)
	c.JSON(http.StatusOK, gin.H{
		"tunnels": tunnels,
		"count":   len(tunnels),
	})
}

// getTunnel handles GET /api/shareruntime/tunnels/:id
func (s *Server) getTunnel(c *gin.Context) {
	tunnelID := c.Param("id")
	tunnel, exists := s.tunnelManager.GetTunnel(tunnelID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tunnel not found",
		})
		return
	}
	c.JSON(http.StatusOK, tunnel)
}

// closeTunnel handles DELETE /api/shareruntime/tunnels/:id
func (s *Server) closeTunnel(c *gin.Context) {
	tunnelID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	// Verify tunnel ownership
	tunnel, exists := s.tunnelManager.GetTunnel(tunnelID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tunnel not found",
		})
		return
	}

	if tunnel.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Not authorized to close this tunnel",
		})
		return
	}

	if err := s.tunnelManager.CloseTunnel(tunnelID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to close tunnel",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tunnel closed"})
}

// listTunnelConnections handles GET /api/shareruntime/tunnels/:id/connections
func (s *Server) listTunnelConnections(c *gin.Context) {
	tunnelID := c.Param("id")
	connections := s.tunnelManager.ListTunnelConnections(tunnelID)
	c.JSON(http.StatusOK, gin.H{
		"connections": connections,
		"count":       len(connections),
	})
}

// validateTunnelAccess handles POST /api/shareruntime/tunnels/:id/validate
func (s *Server) validateTunnelAccess(c *gin.Context) {
	tunnelID := c.Param("id")
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	remoteIP := c.ClientIP()
	if err := s.tunnelManager.ValidateTunnelAccess(tunnelID, req.Token, remoteIP); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// connectTunnel handles POST /api/shareruntime/tunnels/:id/connect
func (s *Server) connectTunnel(c *gin.Context) {
	tunnelID := c.Param("id")
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	remoteIP := c.ClientIP()
	if err := s.tunnelManager.ValidateTunnelAccess(tunnelID, req.Token, remoteIP); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
			"details": err.Error(),
		})
		return
	}

	conn, err := s.tunnelManager.RegisterConnection(tunnelID, remoteIP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to register connection",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, conn)
}

// updateConnectionStats handles PATCH /api/shareruntime/connections/:id/stats
func (s *Server) updateConnectionStats(c *gin.Context) {
	connectionID := c.Param("id")
	var req struct {
		BytesIn  int64 `json:"bytesIn"`
		BytesOut int64 `json:"bytesOut"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := s.tunnelManager.UpdateConnectionStats(connectionID, req.BytesIn, req.BytesOut); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update stats",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stats updated"})
}

// getTunnelByPort handles GET /api/shareruntime/tunnels/by-port/:port
func (s *Server) getTunnelByPort(c *gin.Context) {
	portStr := c.Param("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid port number",
		})
		return
	}

	tunnel, exists := s.tunnelManager.GetTunnelByPort(port)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tunnel not found for port",
		})
		return
	}

	c.JSON(http.StatusOK, tunnel)
}

