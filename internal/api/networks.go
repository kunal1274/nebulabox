package api

import (
    "net/http"
    "time"
    "fmt"

    "github.com/gin-gonic/gin"
)

type Network struct {
    ID      string    `json:"id"`
    Name    string    `json:"name"`
    Driver  string    `json:"driver"`
    Subnet  string    `json:"subnet"`
    Created time.Time `json:"created"`
}

type createNetworkRequest struct {
    Name        string `json:"name" binding:"required"`
    Driver      string `json:"driver"`
    Subnet      string `json:"subnet"`
    WorkspaceID string `json:"workspaceId,omitempty"`
}

func (s *Server) listNetworks(c *gin.Context) {
    workspaceID := c.Query("workspaceId")
    s.netMu.Lock()
    list := make([]Network, 0, len(s.networks))
    s.workspaceMu.Lock()
    for _, n := range s.networks {
        // Filter by workspace if specified
        if workspaceID != "" {
            if wsID, ok := s.networkWorkspaces[n.ID]; !ok || wsID != workspaceID {
                continue
            }
        }
        list = append(list, *n)
    }
    s.workspaceMu.Unlock()
    s.netMu.Unlock()
    c.JSON(http.StatusOK, list)
}

func (s *Server) createNetwork(c *gin.Context) {
    var req createNetworkRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request body", "details": err.Error()})
        return
    }
    // Validate workspace membership if workspace ID provided
    username, _ := c.Get("username")
    if req.WorkspaceID != "" {
        s.teamMu.Lock()
        if _, isMember := s.teamMembers[req.WorkspaceID][username.(string)]; !isMember {
            s.teamMu.Unlock()
            c.JSON(http.StatusForbidden, gin.H{"error":"not a member of this workspace"})
            return
        }
        s.teamMu.Unlock()
    }
    // Check tenant quota
    s.tenantMu.Lock()
    tenantID, hasTenant := s.userTenants[username.(string)]
    if hasTenant {
        if !s.checkTenantQuota(tenantID, "network") {
            s.tenantMu.Unlock()
            c.JSON(http.StatusForbidden, gin.H{"error":"tenant network quota exceeded"})
            return
        }
    }
    s.tenantMu.Unlock()
    if req.Driver == "" { req.Driver = "bridge" }
    id := fmt.Sprintf("net-%d", time.Now().UnixNano())
    n := &Network{ID: id, Name: req.Name, Driver: req.Driver, Subnet: req.Subnet, Created: time.Now()}
    s.netMu.Lock()
    if s.networks == nil { s.networks = make(map[string]*Network) }
    s.networks[id] = n
    s.netMu.Unlock()
    // Associate network with workspace if provided
    if req.WorkspaceID != "" {
        s.workspaceMu.Lock()
        s.networkWorkspaces[id] = req.WorkspaceID
        s.workspaceMu.Unlock()
    }
    c.JSON(http.StatusCreated, n)
}

func (s *Server) deleteNetwork(c *gin.Context) {
    id := c.Param("id")
    if id == "" { c.JSON(http.StatusBadRequest, gin.H{"error":"id required"}); return }
    s.netMu.Lock()
    _, ok := s.networks[id]
    if ok { delete(s.networks, id) }
    s.netMu.Unlock()
    if !ok { c.JSON(http.StatusNotFound, gin.H{"error":"network not found"}); return }
    c.JSON(http.StatusOK, gin.H{"deleted": id})
}


