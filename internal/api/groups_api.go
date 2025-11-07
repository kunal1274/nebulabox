package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/groups"
)

// CreateGroupRequest represents a request to create a container group
type CreateGroupRequest struct {
	Name            string                  `json:"name" binding:"required"`
	Description     string                  `json:"description,omitempty"`
	ParentGroupID   string                  `json:"parentGroupId,omitempty"`
	SharedResources *groups.SharedResources `json:"sharedResources,omitempty"`
}

// UpdateGroupRequest represents a request to update a group
type UpdateGroupRequest struct {
	Name            *string                 `json:"name,omitempty"`
	Description     *string                 `json:"description,omitempty"`
	SharedResources *groups.SharedResources `json:"sharedResources,omitempty"`
}

// CreateRelationshipRequest represents a request to create a container relationship
type CreateRelationshipRequest struct {
	ParentID   string                 `json:"parentId" binding:"required"`
	ChildID    string                 `json:"childId" binding:"required"`
	Type       string                 `json:"type" binding:"required"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// createGroup handles POST /api/groups
func (s *Server) createGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	group, err := s.groupManager.CreateGroup(req.Name, req.Description, req.ParentGroupID, req.SharedResources)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// listGroups handles GET /api/groups
func (s *Server) listGroups(c *gin.Context) {
	var parentID *string
	if pid := c.Query("parentId"); pid != "" {
		parentID = &pid
	}

	groups := s.groupManager.ListGroups(parentID)
	c.JSON(http.StatusOK, gin.H{
		"groups": groups,
		"count":  len(groups),
	})
}

// getGroup handles GET /api/groups/:id
func (s *Server) getGroup(c *gin.Context) {
	groupID := c.Param("id")
	group, exists := s.groupManager.GetGroup(groupID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Group not found",
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// updateGroup handles PATCH /api/groups/:id
func (s *Server) updateGroup(c *gin.Context) {
	groupID := c.Param("id")
	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	group, err := s.groupManager.UpdateGroup(groupID, req.Name, req.Description, req.SharedResources)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// deleteGroup handles DELETE /api/groups/:id
func (s *Server) deleteGroup(c *gin.Context) {
	groupID := c.Param("id")
	force := c.Query("force") == "true"

	if err := s.groupManager.DeleteGroup(groupID, force); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted"})
}

// addContainerToGroup handles POST /api/groups/:id/containers
func (s *Server) addContainerToGroup(c *gin.Context) {
	groupID := c.Param("id")
	var req struct {
		ContainerID string `json:"containerId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := s.groupManager.AddContainerToGroup(groupID, req.ContainerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add container to group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container added to group"})
}

// removeContainerFromGroup handles DELETE /api/groups/:id/containers/:containerId
func (s *Server) removeContainerFromGroup(c *gin.Context) {
	groupID := c.Param("id")
	containerID := c.Param("containerId")

	if err := s.groupManager.RemoveContainerFromGroup(groupID, containerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to remove container from group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container removed from group"})
}

// getGroupHierarchy handles GET /api/groups/:id/hierarchy
func (s *Server) getGroupHierarchy(c *gin.Context) {
	groupID := c.Param("id")
	hierarchy, err := s.groupManager.GetGroupHierarchy(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get hierarchy",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, hierarchy)
}

// createRelationship handles POST /api/groups/relationships
func (s *Server) createRelationship(c *gin.Context) {
	var req CreateRelationshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	relationship, err := s.relationshipManager.CreateRelationship(req.ParentID, req.ChildID, req.Type, req.Properties)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create relationship",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, relationship)
}

// getContainerChildren handles GET /api/groups/relationships/:containerId/children
func (s *Server) getContainerChildren(c *gin.Context) {
	containerID := c.Param("containerId")
	children := s.relationshipManager.GetChildren(containerID)
	c.JSON(http.StatusOK, gin.H{
		"children": children,
		"count":    len(children),
	})
}

// getContainerParent handles GET /api/groups/relationships/:containerId/parent
func (s *Server) getContainerParent(c *gin.Context) {
	containerID := c.Param("containerId")
	parent, exists := s.relationshipManager.GetParent(containerID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No parent relationship found",
		})
		return
	}

	c.JSON(http.StatusOK, parent)
}

// deleteRelationship handles DELETE /api/groups/relationships
func (s *Server) deleteRelationship(c *gin.Context) {
	var req struct {
		ParentID string `json:"parentId" binding:"required"`
		ChildID  string `json:"childId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := s.relationshipManager.DeleteRelationship(req.ParentID, req.ChildID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete relationship",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relationship deleted"})
}

// getContainerTree handles GET /api/groups/relationships/:containerId/tree
func (s *Server) getContainerTree(c *gin.Context) {
	containerID := c.Param("containerId")
	tree := s.relationshipManager.GetContainerTree(containerID)
	c.JSON(http.StatusOK, tree)
}

// getContainerAncestry handles GET /api/groups/relationships/:containerId/ancestry
func (s *Server) getContainerAncestry(c *gin.Context) {
	containerID := c.Param("containerId")
	ancestry := s.relationshipManager.GetAncestry(containerID)
	c.JSON(http.StatusOK, gin.H{
		"ancestry": ancestry,
		"count":    len(ancestry),
	})
}

// startGroup handles POST /api/groups/:id/start
func (s *Server) startGroup(c *gin.Context) {
	groupID := c.Param("id")
	group, exists := s.groupManager.GetGroup(groupID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Group not found",
		})
		return
	}

	// Start all containers in the group
	for _, containerID := range group.ContainerIDs {
		// In production, this would call the containerd client to start containers
		// For now, we just acknowledge the request
		_ = containerID
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group start initiated"})
}

// stopGroup handles POST /api/groups/:id/stop
func (s *Server) stopGroup(c *gin.Context) {
	groupID := c.Param("id")
	group, exists := s.groupManager.GetGroup(groupID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Group not found",
		})
		return
	}

	// Stop all containers in the group
	for _, containerID := range group.ContainerIDs {
		// In production, this would call the containerd client to stop containers
		_ = containerID
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group stop initiated"})
}

