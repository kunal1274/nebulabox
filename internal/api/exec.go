package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
)

// ExecRequest represents a request to execute a command in a container
type ExecRequest struct {
	Command []string          `json:"command" binding:"required"`
	WorkDir string            `json:"workdir,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	User    string            `json:"user,omitempty"`
	TTY     bool              `json:"tty,omitempty"`
}

// ExecResponse represents the result of container execution
type ExecResponse struct {
	ExitCode int    `json:"exitCode"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
}

// execContainer handles POST /api/containers/:id/exec
func (s *Server) execContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	var req ExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Create exec options
	opts := &containerd.ExecOptions{
		Command: req.Command,
		WorkDir: req.WorkDir,
		Env:     req.Env,
		User:    req.User,
		TTY:     req.TTY,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
	}

	// Execute command in container
	result, err := s.containerd.ExecContainer(ctx, id, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute command in container",
			"details": err.Error(),
		})
		return
	}

	response := ExecResponse{
		ExitCode: result.ExitCode,
		Output:   result.Output,
		Error:    result.Error,
	}

	c.JSON(http.StatusOK, response)
}

// execContainerStream handles POST /api/containers/:id/exec/stream
func (s *Server) execContainerStream(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	var req ExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Create exec options
	opts := &containerd.ExecOptions{
		Command: req.Command,
		WorkDir: req.WorkDir,
		Env:     req.Env,
		User:    req.User,
		TTY:     req.TTY,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
	}

	// Set up streaming response
	c.Header("Content-Type", "text/plain")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)

	// Create a writer that flushes to the client
	writer := &streamingWriter{writer: c.Writer}

	// Execute command with streaming output
	err := s.containerd.ExecContainerWithStream(ctx, id, opts, writer, writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute command in container",
			"details": err.Error(),
		})
		return
	}
}

// streamingWriter implements io.Writer with flushing
type streamingWriter struct {
	writer gin.ResponseWriter
}

func (w *streamingWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	if err != nil {
		return n, err
	}
	
	// Flush the response to send data immediately
	if f, ok := w.writer.(http.Flusher); ok {
		f.Flush()
	}
	
	return n, nil
}

// execContainerShell handles POST /api/containers/:id/exec/shell
func (s *Server) execContainerShell(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	// Get shell from query parameter or use default
	shell := c.Query("shell")
	if shell == "" {
		shell = "/bin/bash"
	}

	// Check if shell exists in container by trying to execute it
	ctx := c.Request.Context()
	
	// Create exec options to check shell availability
	opts := &containerd.ExecOptions{
		Command: []string{shell, "--version"},
		TTY:     false,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
	}

	// Try to execute shell version command
	result, err := s.containerd.ExecContainer(ctx, id, opts)
	if err != nil {
		// If shell doesn't exist, try /bin/sh
		shell = "/bin/sh"
		opts.Command = []string{shell, "--version"}
		result, err = s.containerd.ExecContainer(ctx, id, opts)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No suitable shell found in container",
				"details": err.Error(),
			})
			return
		}
	}

	// Return shell information
	c.JSON(http.StatusOK, gin.H{
		"shell": shell,
		"version": strings.TrimSpace(result.Output),
		"available": true,
	})
}
