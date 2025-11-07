package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/shareruntime"
)

// RecordCRDTOperationRequest represents a request to record a CRDT operation
type RecordCRDTOperationRequest struct {
	Type        string                 `json:"type" binding:"required"` // orset, lwwreg, counter, map
	ResourceID  string                 `json:"resourceId" binding:"required"`
	ResourceType string                `json:"resourceType" binding:"required"`
	Operation   string                 `json:"operation" binding:"required"` // add, remove, update, increment
	Key         string                 `json:"key,omitempty"`
	Value       interface{}            `json:"value,omitempty"`
	VectorClock map[string]int64       `json:"vectorClock,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// recordCRDTOperation handles POST /api/shareruntime/workspaces/:id/crdt/operations
func (s *Server) recordCRDTOperation(c *gin.Context) {
	workspaceID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	var req RecordCRDTOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate CRDT type
	var crdtType shareruntime.CRDTType
	switch req.Type {
	case "orset":
		crdtType = shareruntime.CRDTTypeORSet
	case "lwwreg":
		crdtType = shareruntime.CRDTTypeLWWReg
	case "counter":
		crdtType = shareruntime.CRDTTypeCounter
	case "map":
		crdtType = shareruntime.CRDTTypeMap
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid CRDT type. Must be one of: orset, lwwreg, counter, map",
		})
		return
	}

	op := &shareruntime.CRDTOperation{
		Type:        crdtType,
		WorkspaceID: workspaceID,
		ResourceID:  req.ResourceID,
		ResourceType: req.ResourceType,
		Operation:   req.Operation,
		Key:         req.Key,
		Value:       req.Value,
		UserID:      userID,
		VectorClock: req.VectorClock,
		Metadata:    req.Metadata,
	}

	err := s.crdtManager.RecordOperation(op)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to record operation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, op)
}

// getCRDTOperations handles GET /api/shareruntime/workspaces/:id/crdt/operations
func (s *Server) getCRDTOperations(c *gin.Context) {
	workspaceID := c.Param("id")
	
	sinceStr := c.Query("since")
	var since time.Time
	if sinceStr != "" {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid since timestamp format. Use RFC3339.",
			})
			return
		}
	} else {
		// Default to 1 hour ago
		since = time.Now().Add(-1 * time.Hour)
	}

	operations := s.crdtManager.GetOperationsSince(workspaceID, since)
	
	c.JSON(http.StatusOK, gin.H{
		"operations": operations,
		"count":      len(operations),
		"since":      since,
	})
}

// detectConflicts handles POST /api/shareruntime/workspaces/:id/crdt/conflicts/detect
func (s *Server) detectConflicts(c *gin.Context) {
	workspaceID := c.Param("id")
	
	sinceStr := c.Query("since")
	var since time.Time
	if sinceStr != "" {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid since timestamp format. Use RFC3339.",
			})
			return
		}
	} else {
		// Default to 1 hour ago
		since = time.Now().Add(-1 * time.Hour)
	}

	operations := s.crdtManager.GetOperationsSince(workspaceID, since)
	conflicts := s.crdtManager.DetectConflicts(operations)
	
	c.JSON(http.StatusOK, gin.H{
		"conflicts": conflicts,
		"count":     len(conflicts),
	})
}

// resolveConflict handles POST /api/shareruntime/workspaces/:id/crdt/conflicts/:conflictId/resolve
func (s *Server) resolveConflict(c *gin.Context) {
	workspaceID := c.Param("id")
	conflictID := c.Param("conflictId")

	// In a real implementation, we would look up the conflict by ID
	// For now, we'll use a simple approach
	since := time.Now().Add(-1 * time.Hour)
	operations := s.crdtManager.GetOperationsSince(workspaceID, since)
	conflicts := s.crdtManager.DetectConflicts(operations)

	var targetConflict *shareruntime.Conflict
	for _, conflict := range conflicts {
		if conflict.ResourceID == conflictID {
			targetConflict = &conflict
			break
		}
	}

	if targetConflict == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Conflict not found",
		})
		return
	}

	resolution, err := s.crdtManager.ResolveConflict(*targetConflict)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to resolve conflict",
			"details": err.Error(),
		})
		return
	}

	// Record the resolution
	err = s.crdtManager.RecordOperation(resolution)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to record resolution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conflict":   targetConflict,
		"resolution": resolution,
		"message":    "Conflict resolved",
	})
}

// getResourceState handles GET /api/shareruntime/workspaces/:id/crdt/resources/:resourceId
func (s *Server) getResourceState(c *gin.Context) {
	workspaceID := c.Param("id")
	resourceID := c.Param("resourceId")
	resourceType := c.Query("type")

	// Get all operations for this resource
	since := time.Time{} // Get all operations
	operations := s.crdtManager.GetOperationsSince(workspaceID, since)
	
	var resourceOps []*shareruntime.CRDTOperation
	for _, op := range operations {
		if op.ResourceID == resourceID {
			if resourceType == "" || op.ResourceType == resourceType {
				resourceOps = append(resourceOps, op)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceId": resourceID,
		"resourceType": resourceType,
		"operations": resourceOps,
		"count": len(resourceOps),
	})
}

