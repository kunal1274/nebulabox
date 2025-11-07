package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/database/mongodb_repositories"
	"github.com/nebulabox/nebulabox/internal/shareruntime"
)

// getAuditLogs handles GET /api/shareruntime/workspaces/:id/audit-logs
func (s *Server) getAuditLogs(c *gin.Context) {
	workspaceID := c.Param("id")

	// Parse filters
	filters := shareruntime.AuditLogFilters{
		WorkspaceID: workspaceID,
	}

	if action := c.Query("action"); action != "" {
		filters.Action = shareruntime.AuditAction(action)
	}
	if userID := c.Query("userId"); userID != "" {
		filters.UserID = userID
	}
	if resourceType := c.Query("resourceType"); resourceType != "" {
		filters.ResourceType = resourceType
	}
	if resourceID := c.Query("resourceId"); resourceID != "" {
		filters.ResourceID = resourceID
	}
	if successStr := c.Query("success"); successStr != "" {
		success, err := strconv.ParseBool(successStr)
		if err == nil {
			filters.Success = &success
		}
	}
	if startTimeStr := c.Query("startTime"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filters.StartTime = &startTime
		}
	}
	if endTimeStr := c.Query("endTime"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filters.EndTime = &endTime
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}
	if filters.Limit == 0 {
		filters.Limit = 100 // Default limit
	}

	// Try to get logs from MongoDB first
	if s.mongoRepos != nil && s.mongoRepos.AuditLogs != nil {
		ctx := c.Request.Context()
		mongoFilters := mongodb_repositories.AuditLogFilters{
			WorkspaceID:  filters.WorkspaceID,
			UserID:       filters.UserID,
			Action:       string(filters.Action),
			ResourceType: filters.ResourceType,
			ResourceID:   filters.ResourceID,
			Success:      filters.Success,
			Limit:        filters.Limit,
		}
		if filters.StartTime != nil {
			mongoFilters.StartTime = *filters.StartTime
		}
		if filters.EndTime != nil {
			mongoFilters.EndTime = *filters.EndTime
		}
		
		dbLogs, err := s.mongoRepos.AuditLogs.GetLogs(ctx, mongoFilters)
		if err == nil && len(dbLogs) > 0 {
			// Convert MongoDB logs to shareruntime format
			logs := make([]*shareruntime.AuditLog, len(dbLogs))
			for i, dbLog := range dbLogs {
				details := make(map[string]string)
				if dbLog.Details != nil {
					for k, v := range dbLog.Details {
						if str, ok := v.(string); ok {
							details[k] = str
						}
					}
				}
				logs[i] = &shareruntime.AuditLog{
					ID:           dbLog.ID.Hex(),
					WorkspaceID:  dbLog.WorkspaceID,
					UserID:       dbLog.UserID,
					Username:     dbLog.Username,
					Action:       shareruntime.AuditAction(dbLog.Action),
					ResourceType: dbLog.ResourceType,
					ResourceID:   dbLog.ResourceID,
					Details:      details,
					IPAddress:    dbLog.IPAddress,
					UserAgent:    dbLog.UserAgent,
					Success:      dbLog.Success,
					ErrorMessage: dbLog.Message,
					Timestamp:    dbLog.Timestamp,
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"logs": logs,
				"count": len(logs),
			})
			return
		}
		// Log error but continue with fallback
		if err != nil {
			log.Printf("⚠️  Warning: Failed to get audit logs from MongoDB: %v", err)
		}
	}

	// Fallback: get logs from in-memory audit logger
	logs := s.auditLogger.GetLogs(filters)

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"count": len(logs),
	})
}

// getUserAuditLogs handles GET /api/shareruntime/audit-logs
func (s *Server) getUserAuditLogs(c *gin.Context) {
	username, _ := c.Get("username")
	userID := username.(string)

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs := s.auditLogger.GetUserLogs(userID, limit)

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"count": len(logs),
	})
}

// getAuditStats handles GET /api/shareruntime/workspaces/:id/audit-stats
func (s *Server) getAuditStats(c *gin.Context) {
	workspaceID := c.Param("id")

	startTime := time.Now().Add(-24 * time.Hour) // Default: last 24 hours
	endTime := time.Now()

	if startTimeStr := c.Query("startTime"); startTimeStr != "" {
		if parsed, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			startTime = parsed
		}
	}
	if endTimeStr := c.Query("endTime"); endTimeStr != "" {
		if parsed, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			endTime = parsed
		}
	}

	stats := s.auditLogger.GetStats(workspaceID, startTime, endTime)

	c.JSON(http.StatusOK, stats)
}

