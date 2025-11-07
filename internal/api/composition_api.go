package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/composition"
	"github.com/nebulabox/nebulabox/internal/containerd"
)

// CreateCompositionSpecRequest represents a request to create a composition spec
type CreateCompositionSpecRequest struct {
	Name        string                        `json:"name" binding:"required"`
	Description string                        `json:"description,omitempty"`
	Sources     []composition.SourceContainer `json:"sources" binding:"required"`
	Overrides   *composition.ContainerOverrides `json:"overrides,omitempty"`
	Strategy    string                        `json:"strategy,omitempty"`
}

// ComposeContainerRequest represents a request to create a container from a composition
type ComposeContainerRequest struct {
	SpecName    string `json:"specName" binding:"required"`
	ContainerName string `json:"containerName,omitempty"`
	Start       bool   `json:"start,omitempty"`
}

// createCompositionSpec handles POST /api/composition/specs
func (s *Server) createCompositionSpec(c *gin.Context) {
	var req CreateCompositionSpecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	spec := &composition.CompositionSpec{
		Name:        req.Name,
		Description: req.Description,
		Sources:     req.Sources,
		Overrides:   req.Overrides,
		Strategy:    req.Strategy,
	}

	if err := s.compositionManager.SaveSpec(spec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create composition spec",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, spec)
}

// listCompositionSpecs handles GET /api/composition/specs
func (s *Server) listCompositionSpecs(c *gin.Context) {
	specs := s.compositionManager.ListSpecs()
	c.JSON(http.StatusOK, gin.H{
		"specs": specs,
		"count": len(specs),
	})
}

// getCompositionSpec handles GET /api/composition/specs/:name
func (s *Server) getCompositionSpec(c *gin.Context) {
	name := c.Param("name")
	spec, exists := s.compositionManager.GetSpec(name)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Composition spec not found",
		})
		return
	}

	c.JSON(http.StatusOK, spec)
}

// deleteCompositionSpec handles DELETE /api/composition/specs/:name
func (s *Server) deleteCompositionSpec(c *gin.Context) {
	name := c.Param("name")
	if err := s.compositionManager.DeleteSpec(name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete composition spec",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Composition spec deleted"})
}

// previewComposition handles POST /api/composition/preview
func (s *Server) previewComposition(c *gin.Context) {
	var spec composition.CompositionSpec
	if err := c.ShouldBindJSON(&spec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	composed, err := s.compositionManager.ComposeContainer(
		&spec,
		s.getContainerSourceData,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to compose container",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, composed)
}

// composeContainerFromSpec handles POST /api/composition/compose
func (s *Server) composeContainerFromSpec(c *gin.Context) {
	var req ComposeContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	spec, exists := s.compositionManager.GetSpec(req.SpecName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Composition spec not found",
		})
		return
	}

	composed, err := s.compositionManager.ComposeContainer(
		spec,
		s.getContainerSourceData,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to compose container",
			"details": err.Error(),
		})
		return
	}

	// Create container from composed spec
	containerReq := ContainerRequest{
		Image:  composed.Image,
		Name:   req.ContainerName,
		Detach: true,
		Network: composed.Network,
		Service: composed.Service,
	}

	// Convert env vars
	for k, v := range composed.EnvVars {
		containerReq.Env = append(containerReq.Env, k+"="+v)
	}

	// Convert ports
	containerReq.Ports = make([]string, 0)
	for containerPort, hostPort := range composed.Ports {
		containerReq.Ports = append(containerReq.Ports, hostPort+":"+containerPort)
	}

	// Convert volumes
	for source, dest := range composed.Volumes {
		containerReq.Volume = append(containerReq.Volume, source+":"+dest)
	}

	// Health check
	if composed.HealthCheck != nil {
		containerReq.HealthType = composed.HealthCheck.Type
		containerReq.HealthHTTPPath = composed.HealthCheck.HTTPPath
		containerReq.HealthHTTPPort = composed.HealthCheck.HTTPPort
		containerReq.HealthTCPPort = composed.HealthCheck.TCPPort
		containerReq.HealthCmd = composed.HealthCheck.Command
		containerReq.HealthIntervalSec = composed.HealthCheck.IntervalSeconds
		containerReq.HealthTimeoutSec = composed.HealthCheck.TimeoutSeconds
		containerReq.HealthRetries = composed.HealthCheck.Retries
		containerReq.HealthStartPeriod = composed.HealthCheck.StartPeriodSec
	}

	// Create container using existing runContainer logic (simplified)
	ctx := c.Request.Context()
	
	// Pull image
	if err := s.containerd.PullImage(ctx, containerReq.Image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to pull image",
			"details": err.Error(),
		})
		return
	}

	// Create container options
	containerOpts := &containerd.ContainerOptions{
		Name:        containerReq.Name,
		Ports:       make(map[string]string),
		Environment: make(map[string]string),
		Volumes:     make(map[string]string),
		Detach:      containerReq.Detach,
		Network:     containerReq.Network,
	}

	// Ports
	for _, p := range containerReq.Ports {
		host, container := parsePortPair(p)
		if host != "" && container != "" {
			containerOpts.Ports[container] = host
		}
	}

	// Environment
	for _, env := range containerReq.Env {
		parts := splitEnv(env)
		if len(parts) == 2 {
			containerOpts.Environment[parts[0]] = parts[1]
		}
	}

	// Volumes
	for _, vol := range containerReq.Volume {
		parts := splitVolume(vol)
		if len(parts) == 2 {
			containerOpts.Volumes[parts[0]] = parts[1]
		}
	}

	// Health check
	if composed.HealthCheck != nil {
		containerOpts.HealthCheck = &containerd.HealthCheckOptions{
			Type:            composed.HealthCheck.Type,
			HTTPPath:        composed.HealthCheck.HTTPPath,
			HTTPPort:        composed.HealthCheck.HTTPPort,
			TCPPort:         composed.HealthCheck.TCPPort,
			Command:         composed.HealthCheck.Command,
			IntervalSeconds: composed.HealthCheck.IntervalSeconds,
			TimeoutSeconds:  composed.HealthCheck.TimeoutSeconds,
			Retries:         composed.HealthCheck.Retries,
			StartPeriodSec:  composed.HealthCheck.StartPeriodSec,
		}
	}

	// Create container
	container, err := s.containerd.CreateContainer(ctx, containerReq.Image, containerReq.Name, containerOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create container",
			"details": err.Error(),
		})
		return
	}

	// Start if requested
	if req.Start {
		if err := s.containerd.StartContainer(ctx, container.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to start container",
				"details": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"container": ContainerResponse{
			ID:      container.ID,
			Name:    container.Name,
			Image:   container.Image,
			Status:  container.Status,
			Created: container.Created.Format("2006-01-02T15:04:05Z07:00"),
		},
		"composition": composed,
	})
}

// getContainerSourceData extracts source data from a container
func (s *Server) getContainerSourceData(containerID string) (*composition.ContainerSourceData, error) {
	ctx := context.Background()
	
	// Get container details (using mock data structure for now)
	// In production, this would query containerd for actual container configuration
	containers, err := s.containerd.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	var container *containerd.Container
	for _, c := range containers {
		if c.ID == containerID {
			container = c
			break
		}
	}

	if container == nil {
		return nil, &compositionError{message: "container not found: " + containerID}
	}

	// Extract data from container
	// Note: This is simplified - in production, we'd need to extract actual config
	data := &composition.ContainerSourceData{
		ContainerID: containerID,
		Image:       container.Image,
		EnvVars:     make(map[string]string),
		Ports:      make(map[string]string),
		Volumes:    make(map[string]string),
		Labels:     make(map[string]string),
	}

	// In a real implementation, we'd parse the container configuration
	// For now, return mock data structure

	return data, nil
}

type compositionError struct {
	message string
}

func (e *compositionError) Error() string {
	return e.message
}

// Helper functions
func splitEnv(env string) []string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return []string{env[:i], env[i+1:]}
		}
	}
	return []string{env}
}

func splitVolume(vol string) []string {
	for i := 0; i < len(vol); i++ {
		if vol[i] == ':' {
			return []string{vol[:i], vol[i+1:]}
		}
	}
	return []string{vol}
}

