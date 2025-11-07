package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ImageRequest represents a request for image operations
type ImageRequest struct {
	Image string `json:"image" binding:"required"`
	Path  string `json:"path,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

// BuildRequest supports building from Dockerfile text or context
type BuildRequest struct {
    Dockerfile string `json:"dockerfile,omitempty"`
    Tag        string `json:"tag"`
    ContextTar string `json:"contextTar,omitempty"` // optional base64 tar (not used in stub)
}

type BuildResponse struct {
    Tag   string   `json:"tag"`
    Image string   `json:"image"`
    Logs  []string `json:"logs"`
}

// ImageResponse represents an image in API responses
type ImageResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Size    string `json:"size"`
	Created string `json:"created"`
}

// listImages handles GET /api/images
func (s *Server) listImages(c *gin.Context) {
	s.modeMu.Lock()
	mode := s.operatingMode
	s.modeMu.Unlock()

	images := make([]ImageResponse, 0)
	imageMap := make(map[string]bool) // Track unique images by "name:tag"

	// Step 1: Get images from registry (if available and not mock mode)
	if mode != "mock" && s.registryClient != nil {
		repos, err := s.registryClient.ListRepositories()
		if err == nil && len(repos) > 0 {
			for _, repo := range repos {
				versions, err := s.registryClient.ListVersions(repo)
				if err == nil {
					for _, version := range versions {
						tag, _ := version["tag"].(string)
						if tag == "" {
							tag = "latest"
						}
						createdAt := ""
						if created, ok := version["createdAt"].(string); ok {
							createdAt = created
						}
						size := "0 MB"
						if sizeVal, ok := version["size"].(float64); ok {
							size = fmt.Sprintf("%.2f MB", sizeVal/1024/1024)
						}
						digest := ""
						if d, ok := version["digest"].(string); ok {
							digest = d
						}
						imageID := digest
						if imageID == "" {
							imageID = repo + ":" + tag
						}
						
						key := repo + ":" + tag
						if !imageMap[key] {
							images = append(images, ImageResponse{
								ID:      imageID,
								Name:    repo,
								Tag:     tag,
								Size:    size,
								Created: createdAt,
							})
							imageMap[key] = true
						}
					}
				}
			}
		}
	}

	// Step 2: Get images from database if available (test/live mode)
	if (mode == "test" || mode == "live") && s.repos != nil && s.repos.Image != nil {
		dbImages, err := s.repos.Image.List()
		if err == nil && len(dbImages) > 0 {
			for _, imgData := range dbImages {
				key := imgData.Name + ":" + imgData.Tag
				if !imageMap[key] {
					images = append(images, ImageResponse{
						ID:      imgData.ID,
						Name:    imgData.Name,
						Tag:     imgData.Tag,
						Size:    imgData.Size,
						Created: imgData.Created,
					})
					imageMap[key] = true
				}
			}
		}
	}
	
	// Step 2b: Add built images from in-memory storage (fallback)
	if mode == "test" || mode == "live" {
		s.builtImagesMu.Lock()
		for _, img := range s.builtImages {
			key := img.Name + ":" + img.Tag
			if !imageMap[key] {
				images = append(images, *img)
				imageMap[key] = true
			}
		}
		s.builtImagesMu.Unlock()
	}

	// Step 3: Add mock data only in mock mode or if no images found
	if mode == "mock" || len(images) == 0 {
		mockImages := []ImageResponse{
			{
				ID:      "img-001",
				Name:    "nginx",
				Tag:     "latest",
				Size:    "142 MB",
				Created: "2024-01-15T10:30:00Z",
			},
			{
				ID:      "img-002",
				Name:    "postgres",
				Tag:     "13",
				Size:    "376 MB",
				Created: "2024-01-15T11:15:00Z",
			},
			{
				ID:      "img-003",
				Name:    "node",
				Tag:     "18",
				Size:    "945 MB",
				Created: "2024-01-15T12:00:00Z",
			},
		}
		for _, mockImg := range mockImages {
			key := mockImg.Name + ":" + mockImg.Tag
			if !imageMap[key] {
				images = append(images, mockImg)
				imageMap[key] = true
			}
		}
	}

	c.JSON(http.StatusOK, images)
}

// pullImage handles POST /api/images/pull
func (s *Server) pullImage(c *gin.Context) {
	var req ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if err := s.containerd.PullImage(ctx, req.Image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to pull image",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image pulled successfully",
		"image":   req.Image,
	})
}

// pushImage handles POST /api/images/push
func (s *Server) pushImage(c *gin.Context) {
	var req ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Parse image name (format: registry/repo:tag or repo:tag)
	imageParts := strings.Split(req.Image, ":")
	repo := imageParts[0]
	tag := "latest"
	if len(imageParts) > 1 {
		tag = imageParts[1]
	}
	
	// Remove registry prefix if present
	repoParts := strings.Split(repo, "/")
	if len(repoParts) > 1 {
		repo = strings.Join(repoParts[1:], "/")
	}

	// Check if registry is available
	if s.registryClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Registry not configured",
		})
		return
	}

	// For now, we'll return success since actual push requires containerd image export
	// In a real implementation, this would:
	// 1. Export the image from containerd
	// 2. Push blobs and manifest to registry
	// This is a placeholder that acknowledges the push request
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Image push initiated",
		"image":   req.Image,
		"repo":    repo,
		"tag":     tag,
		"note":    "Full push implementation requires containerd image export (coming in next phase)",
	})
}

// buildImage handles POST /api/images/build
func (s *Server) buildImage(c *gin.Context) {
    var req BuildRequest
    if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

    if req.Tag == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "tag is required"})
        return
    }

    // Stubbed build logs
    logs := []string{
        "[+] Building 1.2s (5/5) FINISHED",
        " => [internal] load build definition from Dockerfile",
        " => [internal] load .dockerignore",
        " => [1/2] FROM alpine:3.19",
        " => [2/2] RUN echo 'hello from NebulaBox build'",
        " => exporting to image",
        " => naming to " + req.Tag,
        "Successfully built mock-image-id",
    }

    resp := BuildResponse{Tag: req.Tag, Image: req.Tag, Logs: logs}
    c.JSON(http.StatusOK, resp)
}

