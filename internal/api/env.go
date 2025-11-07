package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
)

// EnvVarRequest represents a request to manage environment variables
type EnvVarRequest struct {
	Variables []containerd.EnvVar `json:"variables"`
}

// EnvOperationRequest represents a request for environment variable operations
type EnvOperationRequest struct {
	Action string `json:"action" binding:"required"` // "set", "unset", "update", "clear"
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Type   string `json:"type,omitempty"`
}

// EnvStringRequest represents a request with environment variables as string
type EnvStringRequest struct {
	EnvString string `json:"envString" binding:"required"`
}

// EnvResponse represents the response for environment variable operations
type EnvResponse struct {
	Success   bool                `json:"success"`
	Message   string              `json:"message"`
	Variables []containerd.EnvVar `json:"variables,omitempty"`
	Error     string              `json:"error,omitempty"`
}

// setEnvVars handles POST /api/containers/:id/env
func (s *Server) setEnvVars(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	var req EnvVarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate environment variables
	for _, env := range req.Variables {
		if err := containerd.ValidateEnvVar(env); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid environment variable",
				"details": err.Error(),
			})
			return
		}
	}

	ctx := c.Request.Context()

	// Set environment variables
	result, err := s.containerd.SetEnvVars(ctx, id, req.Variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to set environment variables",
			"details": err.Error(),
		})
		return
	}

	response := EnvResponse{
		Success:   result.Success,
		Message:   result.Message,
		Variables: result.Variables,
		Error:     result.Error,
	}

	c.JSON(http.StatusOK, response)
}

// getEnvVars handles GET /api/containers/:id/env
func (s *Server) getEnvVars(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	ctx := c.Request.Context()

	// Get environment variables
	result, err := s.containerd.GetEnvVars(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get environment variables",
			"details": err.Error(),
		})
		return
	}

	response := EnvResponse{
		Success:   result.Success,
		Message:   result.Message,
		Variables: result.Variables,
		Error:     result.Error,
	}

	c.JSON(http.StatusOK, response)
}

// updateEnvVar handles PUT /api/containers/:id/env
func (s *Server) updateEnvVar(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	var req EnvOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate operation
	if req.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Action is required",
		})
		return
	}

	ctx := c.Request.Context()

	// Create operation
	operation := &containerd.EnvOperation{
		Action: req.Action,
		Key:    req.Key,
		Value:  req.Value,
		Type:   req.Type,
	}

	// Update environment variable
	result, err := s.containerd.UpdateEnvVar(ctx, id, operation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update environment variable",
			"details": err.Error(),
		})
		return
	}

	response := EnvResponse{
		Success:   result.Success,
		Message:   result.Message,
		Variables: result.Variables,
		Error:     result.Error,
	}

	c.JSON(http.StatusOK, response)
}

// clearEnvVars handles DELETE /api/containers/:id/env
func (s *Server) clearEnvVars(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	ctx := c.Request.Context()

	// Clear environment variables
	result, err := s.containerd.ClearEnvVars(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear environment variables",
			"details": err.Error(),
		})
		return
	}

	response := EnvResponse{
		Success:   result.Success,
		Message:   result.Message,
		Variables: result.Variables,
		Error:     result.Error,
	}

	c.JSON(http.StatusOK, response)
}

// setEnvFromString handles POST /api/containers/:id/env/string
func (s *Server) setEnvFromString(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	var req EnvStringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Parse environment string
	variables, err := containerd.ParseEnvString(req.EnvString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid environment string format",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Set environment variables
	result, err := s.containerd.SetEnvVars(ctx, id, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to set environment variables from string",
			"details": err.Error(),
		})
		return
	}

	response := EnvResponse{
		Success:   result.Success,
		Message:   result.Message,
		Variables: result.Variables,
		Error:     result.Error,
	}

	c.JSON(http.StatusOK, response)
}

// getEnvAsString handles GET /api/containers/:id/env/string
func (s *Server) getEnvAsString(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	ctx := c.Request.Context()

	// Get environment variables
	result, err := s.containerd.GetEnvVars(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get environment variables",
			"details": err.Error(),
		})
		return
	}

	// Format as string
	envString := containerd.FormatEnvString(result.Variables)

	c.JSON(http.StatusOK, gin.H{
		"success":   result.Success,
		"message":   result.Message,
		"envString": envString,
		"count":     len(result.Variables),
	})
}

// validateEnvVar handles POST /api/containers/:id/env/validate
func (s *Server) validateEnvVar(c *gin.Context) {
	var env containerd.EnvVar
	if err := c.ShouldBindJSON(&env); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate environment variable
	err := containerd.ValidateEnvVar(env)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"message": "Environment variable is valid",
	})
}

// parseEnvString handles POST /api/containers/:id/env/parse
func (s *Server) parseEnvString(c *gin.Context) {
	var req EnvStringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Parse environment string
	variables, err := containerd.ParseEnvString(req.EnvString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid environment string format",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"variables": variables,
		"count":     len(variables),
	})
}

// getEnvTemplates handles GET /api/env/templates
func (s *Server) getEnvTemplates(c *gin.Context) {
	templates := []gin.H{
		{
			"name":        "Node.js Application",
			"description": "Common environment variables for Node.js applications",
			"variables": []containerd.EnvVar{
				{Key: "NODE_ENV", Value: "production", Type: "string"},
				{Key: "PORT", Value: "3000", Type: "number"},
				{Key: "DEBUG", Value: "false", Type: "boolean"},
				{Key: "API_URL", Value: "https://api.example.com", Type: "string"},
				{Key: "DATABASE_URL", Value: "***", Type: "secret"},
			},
		},
		{
			"name":        "Python Application",
			"description": "Common environment variables for Python applications",
			"variables": []containerd.EnvVar{
				{Key: "PYTHONPATH", Value: "/app", Type: "string"},
				{Key: "FLASK_ENV", Value: "production", Type: "string"},
				{Key: "FLASK_DEBUG", Value: "false", Type: "boolean"},
				{Key: "SECRET_KEY", Value: "***", Type: "secret"},
				{Key: "DATABASE_URL", Value: "***", Type: "secret"},
			},
		},
		{
			"name":        "Docker Container",
			"description": "Basic environment variables for Docker containers",
			"variables": []containerd.EnvVar{
				{Key: "PATH", Value: "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", Type: "string"},
				{Key: "HOME", Value: "/root", Type: "string"},
				{Key: "USER", Value: "root", Type: "string"},
				{Key: "SHELL", Value: "/bin/bash", Type: "string"},
				{Key: "LANG", Value: "en_US.UTF-8", Type: "string"},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
		"count":     len(templates),
	})
}
