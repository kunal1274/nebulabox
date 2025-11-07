package registry

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// listRepositories lists all repositories
// GET /api/registry/repositories
func (s *Server) listRepositories(c *gin.Context) {
	repos := s.store.ListRepositories()
	c.JSON(http.StatusOK, gin.H{
		"repositories": repos,
		"count":        len(repos),
	})
}

// listVersions lists all versions for a repository
// GET /api/registry/repositories/:repo/versions
func (s *Server) listVersions(c *gin.Context) {
	repo := c.Param("repo")
	vm := s.store.GetVersionMetadata(repo)
	
	versions := vm.ListVersions()
	c.JSON(http.StatusOK, gin.H{
		"repository": repo,
		"versions":    versions,
		"count":       len(versions),
		"latest":      vm.Latest,
	})
}

// getVersion gets version information for a specific tag
// GET /api/registry/repositories/:repo/versions/:tag
func (s *Server) getVersion(c *gin.Context) {
	repo := c.Param("repo")
	tag := c.Param("tag")
	
	vm := s.store.GetVersionMetadata(repo)
	version, ok := vm.GetVersion(tag)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}
	
	c.JSON(http.StatusOK, version)
}

// deleteVersion deletes a version tag
// DELETE /api/registry/repositories/:repo/versions/:tag
func (s *Server) deleteVersion(c *gin.Context) {
	if !s.requireAuth(c) {
		return
	}
	
	repo := c.Param("repo")
	tag := c.Param("tag")
	
	if s.store.DeleteTag(repo, tag) {
		c.JSON(http.StatusOK, gin.H{"message": "version deleted"})
		return
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
}

// getRepositorySummary gets a summary of repository information
// GET /api/registry/repositories/:repo/summary
func (s *Server) getRepositorySummary(c *gin.Context) {
	repo := c.Param("repo")
	vm := s.store.GetVersionMetadata(repo)
	
	summary := vm.GetVersionSummary()
	tags := s.store.ListTags(repo)
	summary.Tags = tags
	summary.TotalVersions = len(tags)
	
	c.JSON(http.StatusOK, summary)
}

