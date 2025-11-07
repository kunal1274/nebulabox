package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// listRegistryRepositories lists all repositories from the registry
// GET /api/registry/repositories
func (s *Server) listRegistryRepositories(c *gin.Context) {
	if s.registryClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Registry not configured",
		})
		return
	}

	repos, err := s.registryClient.ListRepositories()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "Failed to list repositories",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repositories": repos,
		"count":        len(repos),
	})
}

// listRegistryVersions lists all versions for a repository
// GET /api/registry/repositories/:repo/versions
func (s *Server) listRegistryVersions(c *gin.Context) {
	if s.registryClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Registry not configured",
		})
		return
	}

	repo := c.Param("repo")
	versions, err := s.registryClient.ListVersions(repo)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "Failed to list versions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repository": repo,
		"versions":   versions,
		"count":      len(versions),
	})
}

// getRegistrySummary gets a summary of repository information
// GET /api/registry/repositories/:repo/summary
func (s *Server) getRegistrySummary(c *gin.Context) {
	if s.registryClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Registry not configured",
		})
		return
	}

	repo := c.Param("repo")
	summary, err := s.registryClient.GetRepositorySummary(repo)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "Failed to get repository summary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, summary)
}

