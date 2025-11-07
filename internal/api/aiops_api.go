package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/aiops"
)

// RecordMetricRequest represents a request to record metrics
type RecordMetricRequest struct {
	ContainerID string  `json:"containerId" binding:"required"`
	CPU         float64 `json:"cpu"`
	Memory      float64 `json:"memory"`
	NetworkRx   float64 `json:"networkRx"`
	NetworkTx   float64 `json:"networkTx"`
}

// recordMetric handles POST /api/aiops/metrics
func (s *Server) recordMetric(c *gin.Context) {
	var req RecordMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	s.aiOpsAnalytics.RecordMetric(req.ContainerID, aiops.MetricPoint{
		Timestamp: time.Now(),
		CPU:       req.CPU,
		Memory:    req.Memory,
		NetworkRx: req.NetworkRx,
		NetworkTx: req.NetworkTx,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Metric recorded"})
}

// predictResourceUsage handles GET /api/aiops/predict/:containerId
func (s *Server) predictResourceUsage(c *gin.Context) {
	containerID := c.Param("containerId")
	durationStr := c.DefaultQuery("duration", "30m")
	
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		duration = 30 * time.Minute
	}

	prediction, err := s.aiOpsAnalytics.PredictResourceUsage(containerID, duration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate prediction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

// getScalingRecommendation handles GET /api/aiops/scaling/:targetId
func (s *Server) getScalingRecommendation(c *gin.Context) {
	targetID := c.Param("targetId")
	currentReplicas := 1
	
	if replicasStr := c.Query("replicas"); replicasStr != "" {
		// Parse replicas from query (simplified - would parse int in production)
		_ = replicasStr
	}

	recommendation, err := s.aiOpsScaling.GetScalingRecommendation(c.Request.Context(), targetID, currentReplicas)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get scaling recommendation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, recommendation)
}

// setScalingPolicy handles POST /api/aiops/scaling/policy
func (s *Server) setScalingPolicy(c *gin.Context) {
	var policy aiops.ScalingPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	s.aiOpsScaling.SetScalingPolicy(&policy)
	c.JSON(http.StatusOK, gin.H{"message": "Scaling policy set"})
}

// processChatCommand handles POST /api/aiops/chat
func (s *Server) processChatCommand(c *gin.Context) {
	var req struct {
		Command string `json:"command" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	response, err := s.aiOpsChat.ProcessCommand(c.Request.Context(), req.Command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process command",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// getChatHistory handles GET /api/aiops/chat/history
func (s *Server) getChatHistory(c *gin.Context) {
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		// Parse limit (simplified - would parse int in production)
		_ = limitStr
	}

	history := s.aiOpsChat.GetCommandHistory(limit)
	c.JSON(http.StatusOK, gin.H{
		"commands": history,
		"count":    len(history),
	})
}

