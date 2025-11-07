package api

import (
    "net/http"
    "time"
    "fmt"

    "github.com/gin-gonic/gin"
)

type Tenant struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Domain      string    `json:"domain"`
    Created     time.Time `json:"created"`
    CreatedBy   string    `json:"createdBy"`
    Quota       TenantQuota `json:"quota"`
}

type TenantQuota struct {
    MaxContainers int `json:"maxContainers"`
    MaxNetworks   int `json:"maxNetworks"`
    MaxTeams      int `json:"maxTeams"`
    MaxStorageGB  int `json:"maxStorageGB"`
}

type createTenantRequest struct {
    Name   string `json:"name" binding:"required"`
    Domain string `json:"domain"`
    Quota  TenantQuota `json:"quota"`
}

type updateTenantRequest struct {
    Name   string `json:"name"`
    Domain string `json:"domain"`
    Quota  *TenantQuota `json:"quota"`
}

func (s *Server) listTenants(c *gin.Context) {
    username, _ := c.Get("username")
    s.tenantMu.Lock()
    // Check if user is system admin (global role check)
    userRole := s.userRoles[username.(string)]
    isSystemAdmin := userRole == "admin"
    result := []Tenant{}
    if isSystemAdmin {
        // System admin can see all tenants
        for _, t := range s.tenants {
            result = append(result, *t)
        }
    } else {
        // Regular users see only their tenant
        if tenantID, ok := s.userTenants[username.(string)]; ok {
            if t, exists := s.tenants[tenantID]; exists {
                result = append(result, *t)
            }
        }
    }
    s.tenantMu.Unlock()
    c.JSON(http.StatusOK, result)
}

func (s *Server) getTenant(c *gin.Context) {
    id := c.Param("id")
    username, _ := c.Get("username")
    s.tenantMu.Lock()
    t, ok := s.tenants[id]
    if !ok {
        s.tenantMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"tenant not found"})
        return
    }
    // Check access: system admin or tenant member
    userRole := s.userRoles[username.(string)]
    isSystemAdmin := userRole == "admin"
    userTenant, isMember := s.userTenants[username.(string)]
    if !isSystemAdmin && (!isMember || userTenant != id) {
        s.tenantMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"access denied"})
        return
    }
    s.tenantMu.Unlock()
    c.JSON(http.StatusOK, t)
}

func (s *Server) createTenant(c *gin.Context) {
    var req createTenantRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
    }
    username, _ := c.Get("username")
    // Only system admin can create tenants
    userRole := s.userRoles[username.(string)]
    if userRole != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error":"only system admin can create tenants"})
        return
    }
    id := fmt.Sprintf("tenant-%d", time.Now().UnixNano())
    quota := req.Quota
    if quota.MaxContainers == 0 { quota.MaxContainers = 100 }
    if quota.MaxNetworks == 0 { quota.MaxNetworks = 50 }
    if quota.MaxTeams == 0 { quota.MaxTeams = 20 }
    if quota.MaxStorageGB == 0 { quota.MaxStorageGB = 1000 }
    t := &Tenant{
        ID:        id,
        Name:      req.Name,
        Domain:    req.Domain,
        Created:   time.Now(),
        CreatedBy: username.(string),
        Quota:     quota,
    }
    s.tenantMu.Lock()
    if s.tenants == nil { s.tenants = make(map[string]*Tenant) }
    s.tenants[id] = t
    s.tenantMu.Unlock()
    c.JSON(http.StatusCreated, t)
}

func (s *Server) updateTenant(c *gin.Context) {
    id := c.Param("id")
    var req updateTenantRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
    }
    username, _ := c.Get("username")
    s.tenantMu.Lock()
    t, ok := s.tenants[id]
    if !ok {
        s.tenantMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"tenant not found"}); return
    }
    // Only system admin can update tenants
    userRole := s.userRoles[username.(string)]
    if userRole != "admin" {
        s.tenantMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"only system admin can update tenants"}); return
    }
    if req.Name != "" { t.Name = req.Name }
    if req.Domain != "" { t.Domain = req.Domain }
    if req.Quota != nil { t.Quota = *req.Quota }
    s.tenantMu.Unlock()
    c.JSON(http.StatusOK, t)
}

func (s *Server) deleteTenant(c *gin.Context) {
    id := c.Param("id")
    username, _ := c.Get("username")
    s.tenantMu.Lock()
    _, ok := s.tenants[id]
    if !ok {
        s.tenantMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"tenant not found"}); return
    }
    // Only system admin can delete tenants
    userRole := s.userRoles[username.(string)]
    if userRole != "admin" {
        s.tenantMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"only system admin can delete tenants"}); return
    }
    delete(s.tenants, id)
    // Clean up user associations
    for user, tenantID := range s.userTenants {
        if tenantID == id {
            delete(s.userTenants, user)
        }
    }
    s.tenantMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"deleted": id})
}

func (s *Server) assignUserToTenant(c *gin.Context) {
    tenantID := c.Param("id")
    var req struct {
        Username string `json:"username" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
    }
    username, _ := c.Get("username")
    // Only system admin can assign users
    userRole := s.userRoles[username.(string)]
    if userRole != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error":"only system admin can assign users"})
        return
    }
    s.tenantMu.Lock()
    _, tenantExists := s.tenants[tenantID]
    if !tenantExists {
        s.tenantMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"tenant not found"}); return
    }
    s.authMu.Lock()
    _, userExists := s.users[req.Username]
    s.authMu.Unlock()
    if !userExists {
        s.tenantMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"user not found"}); return
    }
    if s.userTenants == nil { s.userTenants = make(map[string]string) }
    s.userTenants[req.Username] = tenantID
    s.tenantMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"username": req.Username, "tenantId": tenantID})
}

func (s *Server) checkTenantQuota(tenantID string, resourceType string) bool {
    s.tenantMu.Lock()
    defer s.tenantMu.Unlock()
    t, ok := s.tenants[tenantID]
    if !ok { return false }
    switch resourceType {
    case "container":
        // Count containers for this tenant
        count := 0
        s.workspaceMu.Lock()
        for _, wsID := range s.containerWorkspaces {
            // Check if workspace belongs to tenant's teams
            s.teamMu.Lock()
            for _, team := range s.teams {
                if team.ID == wsID {
                    // For simplicity, assume teams inherit tenant from creator
                    // In production, you'd track team->tenant mapping
                    count++
                    break
                }
            }
            s.teamMu.Unlock()
        }
        s.workspaceMu.Unlock()
        return count < t.Quota.MaxContainers
    case "network":
        count := 0
        s.workspaceMu.Lock()
        for _, wsID := range s.networkWorkspaces {
            if wsID == tenantID { count++ }
        }
        s.workspaceMu.Unlock()
        return count < t.Quota.MaxNetworks
    case "team":
        count := 0
        s.teamMu.Lock()
        for _, team := range s.teams {
            // Simplified: check if team creator is in tenant
            if _, ok := s.userTenants[team.CreatedBy]; ok {
                count++
            }
        }
        s.teamMu.Unlock()
        return count < t.Quota.MaxTeams
    }
    return true
}

func (s *Server) getTenantUsage(c *gin.Context) {
    tenantID := c.Param("id")
    username, _ := c.Get("username")
    s.tenantMu.Lock()
    t, ok := s.tenants[tenantID]
    if !ok {
        s.tenantMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"tenant not found"})
        return
    }
    // Check access: system admin or tenant member
    userRole := s.userRoles[username.(string)]
    isSystemAdmin := userRole == "admin"
    userTenant, isMember := s.userTenants[username.(string)]
    if !isSystemAdmin && (!isMember || userTenant != tenantID) {
        s.tenantMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"access denied"})
        return
    }
    usage := map[string]int{
        "containers": 0,
        "networks":    0,
        "teams":       0,
    }
    s.workspaceMu.Lock()
    // Count containers (simplified - would need proper tenant->workspace mapping)
    for range s.containerWorkspaces {
        usage["containers"]++
    }
    for _, wsID := range s.networkWorkspaces {
        if wsID == tenantID { usage["networks"]++ }
    }
    s.workspaceMu.Unlock()
    s.teamMu.Lock()
    for _, team := range s.teams {
        if creatorTenant, ok := s.userTenants[team.CreatedBy]; ok && creatorTenant == tenantID {
            usage["teams"]++
        }
    }
    s.teamMu.Unlock()
    s.tenantMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"tenant": t, "usage": usage, "quota": t.Quota})
}

