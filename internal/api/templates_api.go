package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/templates"
)

// listTemplates handles GET /api/templates
func (s *Server) listTemplates(c *gin.Context) {
	category := c.Query("category")
	tag := c.Query("tag")

	templateList := s.templateManager.ListTemplates(category, tag)
	c.JSON(http.StatusOK, gin.H{
		"templates": templateList,
		"count":     len(templateList),
	})
}

// getTemplate handles GET /api/templates/:id
func (s *Server) getTemplate(c *gin.Context) {
	id := c.Param("id")
	template, exists := s.templateManager.GetTemplate(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Template not found",
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// saveTemplate handles POST /api/templates
func (s *Server) saveTemplate(c *gin.Context) {
	var template templates.StackTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := s.templateManager.SaveTemplate(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to save template",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// deleteTemplate handles DELETE /api/templates/:id
func (s *Server) deleteTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := s.templateManager.DeleteTemplate(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete template",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted"})
}

// deployTemplate handles POST /api/templates/:id/deploy
func (s *Server) deployTemplate(c *gin.Context) {
	id := c.Param("id")
	template, exists := s.templateManager.GetTemplate(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Template not found",
		})
		return
	}

	var req struct {
		Prefix  string            `json:"prefix,omitempty"`  // Prefix for container names
		EnvVars map[string]string `json:"envVars,omitempty"` // Override env vars
		Start   bool              `json:"start,omitempty"`   // Start containers immediately
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// In a full implementation, this would:
	// 1. Create networks and volumes
	// 2. Create containers in dependency order
	// 3. Start containers if requested
	// For now, return the template configuration

	deployment := gin.H{
		"template":  template,
		"prefix":    req.Prefix,
		"containers": len(template.Containers),
		"status":    "pending",
		"message":   "Template deployment initiated. In production, this would create and start all containers.",
	}

	c.JSON(http.StatusOK, deployment)
}

