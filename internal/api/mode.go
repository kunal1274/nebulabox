package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
)

// ModeRequest represents a mode change request
type ModeRequest struct {
	Mode string `json:"mode" binding:"required"` // "mock", "test", or "live"
}

// ModeResponse represents the current mode
type ModeResponse struct {
	Mode        string `json:"mode"`
	Description string `json:"description"`
	BuiltImages int    `json:"builtImages"`
	TotalImages int    `json:"totalImages"`
}

// setupModeRoutes sets up routes for mode management
func (s *Server) setupModeRoutes() {
	// Get current mode
	s.router.GET("/api/mode", s.getMode)
	
	// Set mode
	s.router.PUT("/api/mode", s.setMode)
}

// getMode returns the current operating mode
func (s *Server) getMode(c *gin.Context) {
	s.modeMu.Lock()
	mode := s.operatingMode
	s.modeMu.Unlock()

	s.builtImagesMu.Lock()
	builtCount := len(s.builtImages)
	s.builtImagesMu.Unlock()

	var description string
	switch mode {
	case "mock":
		description = "Mock mode: Returns only mock/test data, no persistence"
	case "test":
		description = "Test mode: UAT sandbox - Built images and containers are stored in memory"
	case "live":
		description = "Live mode: Production mode with full persistence and registry integration"
	default:
		description = "Unknown mode"
	}

	// Get total images count
	images, _ := c.Get("images")
	totalImages := 0
	if imgList, ok := images.([]ImageResponse); ok {
		totalImages = len(imgList)
	}

	c.JSON(http.StatusOK, ModeResponse{
		Mode:        mode,
		Description: description,
		BuiltImages: builtCount,
		TotalImages: totalImages,
	})
}

// setMode changes the operating mode
func (s *Server) setMode(c *gin.Context) {
	var req ModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate mode
	if req.Mode != "mock" && req.Mode != "test" && req.Mode != "live" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid mode. Must be 'mock', 'test', or 'live'",
		})
		return
	}

	s.modeMu.Lock()
	oldMode := s.operatingMode
	s.operatingMode = req.Mode
	s.modeMu.Unlock()

	// Clear built images and containers when switching to mock mode
	if req.Mode == "mock" {
		s.builtImagesMu.Lock()
		s.builtImages = make(map[string]*ImageResponse)
		s.builtImagesMu.Unlock()
		s.builtContainersMu.Lock()
		s.builtContainers = make(map[string]*containerd.Container)
		s.builtContainersMu.Unlock()
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Mode changed successfully",
		"oldMode":   oldMode,
		"newMode":   req.Mode,
		"mode":      req.Mode,
		"builtImagesCleared": req.Mode == "mock",
	})
}

