package api

import (
    "net/http"
    "time"
    "fmt"

    "github.com/gin-gonic/gin"
)

type Team struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Created     time.Time `json:"created"`
    CreatedBy   string    `json:"createdBy"`
}

type TeamMember struct {
    Username string `json:"username"`
    Role     string `json:"role"` // admin|editor|viewer
    JoinedAt time.Time `json:"joinedAt"`
}

type createTeamRequest struct {
    Name        string `json:"name" binding:"required"`
    Description string `json:"description"`
}

type inviteMemberRequest struct {
    Username string `json:"username" binding:"required"`
    Role     string `json:"role"` // admin|editor|viewer
}

func (s *Server) listTeams(c *gin.Context) {
    username, _ := c.Get("username")
    s.teamMu.Lock()
    // Return teams where user is a member
    result := []Team{}
    for _, t := range s.teams {
        if _, isMember := s.teamMembers[t.ID][username.(string)]; isMember {
            result = append(result, *t)
        }
    }
    s.teamMu.Unlock()
    c.JSON(http.StatusOK, result)
}

func (s *Server) getTeam(c *gin.Context) {
    id := c.Param("id")
    username, _ := c.Get("username")
    s.teamMu.Lock()
    t, ok := s.teams[id]
    if !ok {
        s.teamMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"team not found"})
        return
    }
    // Check membership
    if _, isMember := s.teamMembers[id][username.(string)]; !isMember {
        s.teamMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"not a member of this team"})
        return
    }
    // Get members
    members := []TeamMember{}
    for uname, m := range s.teamMembers[id] {
        members = append(members, TeamMember{Username: uname, Role: m.Role, JoinedAt: m.JoinedAt})
    }
    s.teamMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"team": t, "members": members})
}

func (s *Server) createTeam(c *gin.Context) {
    var req createTeamRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
    }
    username, _ := c.Get("username")
    // Check tenant quota
    s.tenantMu.Lock()
    tenantID, hasTenant := s.userTenants[username.(string)]
    if hasTenant {
        if !s.checkTenantQuota(tenantID, "team") {
            s.tenantMu.Unlock()
            c.JSON(http.StatusForbidden, gin.H{"error":"tenant team quota exceeded"})
            return
        }
    }
    s.tenantMu.Unlock()
    id := fmt.Sprintf("team-%d", time.Now().UnixNano())
    t := &Team{
        ID:          id,
        Name:        req.Name,
        Description: req.Description,
        Created:     time.Now(),
        CreatedBy:   username.(string),
    }
    s.teamMu.Lock()
    if s.teams == nil { s.teams = make(map[string]*Team) }
    if s.teamMembers == nil { s.teamMembers = make(map[string]map[string]*TeamMember) }
    s.teams[id] = t
    if s.teamMembers[id] == nil { s.teamMembers[id] = make(map[string]*TeamMember) }
    // Add creator as admin
    s.teamMembers[id][username.(string)] = &TeamMember{
        Username: username.(string),
        Role:     "admin",
        JoinedAt: time.Now(),
    }
    s.teamMu.Unlock()
    c.JSON(http.StatusCreated, t)
}

func (s *Server) updateTeam(c *gin.Context) {
    id := c.Param("id")
    var req struct {
        Name        string `json:"name"`
        Description string `json:"description"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
    }
    username, _ := c.Get("username")
    s.teamMu.Lock()
    t, ok := s.teams[id]
    if !ok {
        s.teamMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"team not found"}); return
    }
    // Check membership and role (admin or editor can update)
    member, isMember := s.teamMembers[id][username.(string)]
    if !isMember || (member.Role != "admin" && member.Role != "editor") {
        s.teamMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"insufficient permissions"}); return
    }
    if req.Name != "" { t.Name = req.Name }
    if req.Description != "" { t.Description = req.Description }
    s.teamMu.Unlock()
    c.JSON(http.StatusOK, t)
}

func (s *Server) deleteTeam(c *gin.Context) {
    id := c.Param("id")
    username, _ := c.Get("username")
    s.teamMu.Lock()
    _, ok := s.teams[id]
    if !ok {
        s.teamMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"team not found"}); return
    }
    // Only team admin can delete
    member, isMember := s.teamMembers[id][username.(string)]
    if !isMember || member.Role != "admin" {
        s.teamMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"only team admin can delete"}); return
    }
    delete(s.teams, id)
    delete(s.teamMembers, id)
    s.teamMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"deleted": id})
}

func (s *Server) inviteMember(c *gin.Context) {
    id := c.Param("id")
    var req inviteMemberRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
    }
    if req.Role == "" { req.Role = "viewer" }
    if req.Role != "admin" && req.Role != "editor" && req.Role != "viewer" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid role"}); return
    }
    username, _ := c.Get("username")
    s.teamMu.Lock()
    _, ok := s.teams[id]
    if !ok {
        s.teamMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"team not found"}); return
    }
    // Check inviter is admin
    inviter, isMember := s.teamMembers[id][username.(string)]
    if !isMember || inviter.Role != "admin" {
        s.teamMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"only team admin can invite"}); return
    }
    // Check user exists
    s.authMu.Lock()
    _, userExists := s.users[req.Username]
    s.authMu.Unlock()
    if !userExists {
        s.teamMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"user not found"}); return
    }
    // Add member
    if s.teamMembers[id] == nil { s.teamMembers[id] = make(map[string]*TeamMember) }
    s.teamMembers[id][req.Username] = &TeamMember{
        Username: req.Username,
        Role:     req.Role,
        JoinedAt:  time.Now(),
    }
    s.teamMu.Unlock()
    c.JSON(http.StatusCreated, gin.H{"username": req.Username, "role": req.Role})
}

func (s *Server) removeMember(c *gin.Context) {
    id := c.Param("id")
    targetUser := c.Param("username")
    username, _ := c.Get("username")
    s.teamMu.Lock()
    _, ok := s.teams[id]
    if !ok {
        s.teamMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"team not found"}); return
    }
    // Check remover is admin
    remover, isMember := s.teamMembers[id][username.(string)]
    if !isMember || remover.Role != "admin" {
        s.teamMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"only team admin can remove members"}); return
    }
    // Can't remove self
    if targetUser == username.(string) {
        s.teamMu.Unlock()
        c.JSON(http.StatusBadRequest, gin.H{"error":"cannot remove yourself"}); return
    }
    delete(s.teamMembers[id], targetUser)
    s.teamMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"removed": targetUser})
}

func (s *Server) updateMemberRole(c *gin.Context) {
    id := c.Param("id")
    targetUser := c.Param("username")
    var req struct {
        Role string `json:"role" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
    }
    if req.Role != "admin" && req.Role != "editor" && req.Role != "viewer" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid role"}); return
    }
    username, _ := c.Get("username")
    s.teamMu.Lock()
    _, ok := s.teams[id]
    if !ok {
        s.teamMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"team not found"}); return
    }
    // Check updater is admin
    updater, isMember := s.teamMembers[id][username.(string)]
    if !isMember || updater.Role != "admin" {
        s.teamMu.Unlock()
        c.JSON(http.StatusForbidden, gin.H{"error":"only team admin can update roles"}); return
    }
    member, exists := s.teamMembers[id][targetUser]
    if !exists {
        s.teamMu.Unlock()
        c.JSON(http.StatusNotFound, gin.H{"error":"member not found"}); return
    }
    member.Role = req.Role
    s.teamMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"username": targetUser, "role": req.Role})
}

