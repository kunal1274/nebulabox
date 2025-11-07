package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getEndpointMetrics returns performance metrics for all endpoints
// GET /api/performance/endpoints
func (s *Server) getEndpointMetrics(c *gin.Context) {
	if s.endpointMetrics == nil {
		c.JSON(http.StatusOK, gin.H{"endpoints": map[string]interface{}{}})
		return
	}
	
	stats := s.endpointMetrics.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"endpoints": stats,
		"timestamp": gin.H{"unix": gin.H{"seconds": gin.H{}}},
	})
}

