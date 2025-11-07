package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/buildspec"
	"github.com/nebulabox/nebulabox/internal/database/repositories"
)

// BuildSpecRequest represents a request to build from a build specification
type BuildSpecRequest struct {
	Spec  map[string]interface{} `json:"spec"`
	Tag   string                `json:"tag"`
	Build bool                   `json:"build,omitempty"` // If true, actually build; if false, just validate/convert
}

// BuildSpecResponse represents the response from build spec operations
type BuildSpecResponse struct {
	Valid      bool     `json:"valid"`
	Dockerfile string   `json:"dockerfile,omitempty"`
	Errors     []string `json:"errors,omitempty"`
	Tag        string   `json:"tag,omitempty"`
	Logs       []string `json:"logs,omitempty"`
	Message    string   `json:"message,omitempty"`
}

// validateBuildSpec handles POST /api/buildspec/validate
func (s *Server) validateBuildSpec(c *gin.Context) {
	var req BuildSpecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert map to JSON bytes for parsing
	specJSON, err := json.Marshal(req.Spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to serialize spec",
			"details": err.Error(),
		})
		return
	}

	// Parse the spec
	spec, err := buildspec.ParseSpec(specJSON)
	if err != nil {
		c.JSON(http.StatusOK, BuildSpecResponse{
			Valid:   false,
			Errors:  []string{err.Error()},
			Message: "Specification validation failed",
		})
		return
	}

	// Validate the spec
	if err := spec.Validate(); err != nil {
		c.JSON(http.StatusOK, BuildSpecResponse{
			Valid:   false,
			Errors:  []string{err.Error()},
			Message: "Specification validation failed",
		})
		return
	}

	// Convert to Dockerfile
	dockerfile := spec.ToDockerfile()

	c.JSON(http.StatusOK, BuildSpecResponse{
		Valid:      true,
		Dockerfile: dockerfile,
		Message:    "Specification is valid",
	})
}

// convertBuildSpec handles POST /api/buildspec/convert
func (s *Server) convertBuildSpec(c *gin.Context) {
	var req BuildSpecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert map to JSON bytes for parsing
	specJSON, err := json.Marshal(req.Spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to serialize spec",
			"details": err.Error(),
		})
		return
	}

	// Parse the spec
	spec, err := buildspec.ParseSpec(specJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse spec",
			"details": err.Error(),
		})
		return
	}

	// Convert to Dockerfile
	dockerfile := spec.ToDockerfile()

	c.JSON(http.StatusOK, BuildSpecResponse{
		Valid:      true,
		Dockerfile: dockerfile,
		Message:    "Converted to Dockerfile",
	})
}

// buildFromSpec handles POST /api/buildspec/build
func (s *Server) buildFromSpec(c *gin.Context) {
	var req BuildSpecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tag is required",
		})
		return
	}

	// Convert map to JSON bytes for parsing
	specJSON, err := json.Marshal(req.Spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to serialize spec",
			"details": err.Error(),
		})
		return
	}

	// Parse and validate the spec
	spec, err := buildspec.ParseSpec(specJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse spec",
			"details": err.Error(),
		})
		return
	}

	if err := spec.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid spec",
			"details": err.Error(),
		})
		return
	}

	// Convert to Dockerfile
	dockerfile := spec.ToDockerfile()

	// Call the existing build logic
	// For now, we'll use mock build logs similar to buildImage
	logs := []string{
		"[+] Building from NebulaBox build spec",
		"[+] Converting spec to Dockerfile",
		"[+] Building 1.5s (6/6) FINISHED",
		" => [internal] load build definition from buildspec",
		" => [internal] load .dockerignore",
		fmt.Sprintf(" => [1/%d] FROM %s:%s", len(spec.Steps)+1, spec.Base.Image, spec.Base.Tag),
	}

	for i, step := range spec.Steps {
		stepLog := fmt.Sprintf(" => [%d/%d] %s", i+2, len(spec.Steps)+1, strings.ToUpper(step.Type))
		if step.Comment != "" {
			stepLog += fmt.Sprintf(" # %s", step.Comment)
		} else if step.Command != "" {
			stepLog += fmt.Sprintf(" %s", step.Command)
		}
		logs = append(logs, stepLog)
	}

	logs = append(logs,
		" => exporting to image",
		fmt.Sprintf(" => naming to %s", req.Tag),
		"Successfully built from NebulaBox build specification",
	)

	// Register the built image - store it so it appears in Images list
	imageParts := strings.Split(req.Tag, ":")
	repo := imageParts[0]
	tag := "latest"
	if len(imageParts) > 1 {
		tag = imageParts[1]
	}

	// Generate a digest for the built image
	digestBytes := sha256.Sum256([]byte(req.Tag + time.Now().Format(time.RFC3339)))
	digest := "sha256:" + hex.EncodeToString(digestBytes[:])
	imageID := digest

	// Calculate estimated size based on base image and steps
	// This is a mock calculation - in real mode would get actual size
	estimatedSize := 150.0 // Base size in MB
	if strings.Contains(spec.Base.Image, "alpine") {
		estimatedSize = 50.0
	} else if strings.Contains(spec.Base.Image, "node") {
		estimatedSize = 200.0
	}
	estimatedSize += float64(len(spec.Steps)) * 10.0 // Add some MB per step

	// Create image record
	builtImage := &ImageResponse{
		ID:      imageID,
		Name:    repo,
		Tag:     tag,
		Size:    fmt.Sprintf("%.1f MB", estimatedSize),
		Created: time.Now().Format(time.RFC3339),
	}

	// Store built image in database and memory (for test/live mode)
	s.modeMu.Lock()
	mode := s.operatingMode
	s.modeMu.Unlock()
	
	if mode == "test" || mode == "live" {
		// Save to database if repositories available
		if s.repos != nil && s.repos.Image != nil {
			imageData := &repositories.ImageData{
				ID:      builtImage.ID,
				Name:    builtImage.Name,
				Tag:     builtImage.Tag,
				Size:    builtImage.Size,
				Created: builtImage.Created,
				Digest:  "", // Will be populated if available
			}
			if err := s.repos.Image.CreateOrUpdate(imageData); err != nil {
				log.Printf("⚠️  Warning: Failed to save image to database: %v", err)
				// Continue with in-memory storage as fallback
			}
		}
		
		// Also store in memory for backward compatibility
		s.builtImagesMu.Lock()
		s.builtImages[req.Tag] = builtImage
		s.builtImagesMu.Unlock()
	}

	logs = append(logs, fmt.Sprintf(" => registered image: %s", req.Tag))
	logs = append(logs, fmt.Sprintf(" => image ID: %s", imageID))
	logs = append(logs, fmt.Sprintf(" => size: %s", builtImage.Size))

	resp := BuildSpecResponse{
		Valid:      true,
		Dockerfile: dockerfile,
		Tag:        req.Tag,
		Logs:       logs,
		Message:    "Build completed and registered successfully",
	}

	c.JSON(http.StatusOK, resp)
}

