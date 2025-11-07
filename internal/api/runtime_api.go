package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/runtime"
)

// CreateRuntimeContainerRequest represents a request to create a container using Nebula runtime
type CreateRuntimeContainerRequest struct {
	ID          string            `json:"id" binding:"required"`
	Name        string            `json:"name,omitempty"`
	Image       string            `json:"image" binding:"required"`
	Command     []string          `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	WorkingDir  string            `json:"workingDir,omitempty"`
	User        string            `json:"user,omitempty"`
	NetworkMode string            `json:"networkMode,omitempty"`
	Network     string            `json:"network,omitempty"`
	Ports       map[string]string `json:"ports,omitempty"`
	Volumes     map[string]string `json:"volumes,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Resources   *runtime.ResourceLimits `json:"resources,omitempty"`
	Security    *runtime.SecurityConfig `json:"security,omitempty"`
}

// createRuntimeContainer handles POST /api/runtime/containers
func (s *Server) createRuntimeContainer(c *gin.Context) {
	var req CreateRuntimeContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	spec := &runtime.ContainerSpec{
		ID:          req.ID,
		Name:        req.Name,
		Image:       req.Image,
		Command:     req.Command,
		Args:        req.Args,
		Env:         req.Env,
		WorkingDir:  req.WorkingDir,
		User:        req.User,
		NetworkMode: req.NetworkMode,
		Network:     req.Network,
		Ports:       req.Ports,
		Volumes:     req.Volumes,
		Labels:      req.Labels,
		Resources:   req.Resources,
		Security:    req.Security,
	}

	container, err := s.nebulaRuntime.CreateContainer(c.Request.Context(), spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create container",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, container)
}

// listRuntimeContainers handles GET /api/runtime/containers
func (s *Server) listRuntimeContainers(c *gin.Context) {
	containers, err := s.nebulaRuntime.ListContainers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list containers",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"containers": containers,
		"count":      len(containers),
	})
}

// getRuntimeContainer handles GET /api/runtime/containers/:id
func (s *Server) getRuntimeContainer(c *gin.Context) {
	id := c.Param("id")
	container, err := s.nebulaRuntime.GetContainer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Container not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, container)
}

// startRuntimeContainer handles POST /api/runtime/containers/:id/start
func (s *Server) startRuntimeContainer(c *gin.Context) {
	id := c.Param("id")
	if err := s.nebulaRuntime.StartContainer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to start container",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container started"})
}

// stopRuntimeContainer handles POST /api/runtime/containers/:id/stop
func (s *Server) stopRuntimeContainer(c *gin.Context) {
	id := c.Param("id")
	if err := s.nebulaRuntime.StopContainer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to stop container",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container stopped"})
}

// deleteRuntimeContainer handles DELETE /api/runtime/containers/:id
func (s *Server) deleteRuntimeContainer(c *gin.Context) {
	id := c.Param("id")
	if err := s.nebulaRuntime.DeleteContainer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete container",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container deleted"})
}

// pullRuntimeImage handles POST /api/runtime/images/pull
func (s *Server) pullRuntimeImage(c *gin.Context) {
	var req struct {
		Image string `json:"image" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := s.nebulaRuntime.PullImage(c.Request.Context(), req.Image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to pull image",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image pulled successfully"})
}

// listRuntimeImages handles GET /api/runtime/images
func (s *Server) listRuntimeImages(c *gin.Context) {
	images, err := s.nebulaRuntime.ListImages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list images",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images": images,
		"count":  len(images),
	})
}

// getRuntimeInfo handles GET /api/runtime/info
func (s *Server) getRuntimeInfo(c *gin.Context) {
	info, err := s.nebulaRuntime.Info(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get runtime info",
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// getRuntimeVersion handles GET /api/runtime/version
func (s *Server) getRuntimeVersion(c *gin.Context) {
	version, err := s.nebulaRuntime.Version(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get version",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"version": version})
}

