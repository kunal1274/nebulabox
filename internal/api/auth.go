package api

import (
    "crypto/rand"
    "encoding/hex"
    "net/http"
    "github.com/gin-gonic/gin"
)

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func generateToken() string {
    b := make([]byte, 16)
    _, _ = rand.Read(b)
    return hex.EncodeToString(b)
}

func (s *Server) postAuthLogin(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return
    }
    s.authMu.Lock()
    pass, ok := s.users[req.Username]
    if !ok || pass != req.Password {
        s.authMu.Unlock()
        c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"})
        return
    }
    token := generateToken()
    s.sessions[token] = req.Username
    role := s.userRoles[req.Username]
    if role == "" { role = "viewer" }
    s.authMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"username": req.Username, "role": role}})
}

func (s *Server) postAuthLogout(c *gin.Context) {
    token := c.GetHeader("X-App-Auth")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"missing token"}); return
    }
    s.authMu.Lock()
    delete(s.sessions, token)
    s.authMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"status":"logged out"})
}

func (s *Server) getAuthMe(c *gin.Context) {
    token := c.GetHeader("X-App-Auth")
    if token == "" { c.JSON(http.StatusOK, gin.H{"user": nil}); return }
    s.authMu.Lock()
    username, ok := s.sessions[token]
    role := "viewer"
    if ok { role = s.userRoles[username]; if role == "" { role = "viewer" } }
    s.authMu.Unlock()
    if !ok { c.JSON(http.StatusOK, gin.H{"user": nil}); return }
    c.JSON(http.StatusOK, gin.H{"user": gin.H{"username": username, "role": role}})
}

func (s *Server) authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("X-App-Auth")
        if token == "" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"missing auth token"}); return }
        s.authMu.Lock()
        username, ok := s.sessions[token]
        s.authMu.Unlock()
        if !ok { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"invalid session"}); return }
        c.Set("username", username)
        c.Next()
    }
}

func (s *Server) requireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        username, exists := c.Get("username")
        if !exists { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"not authenticated"}); return }
        s.authMu.Lock()
        role := s.userRoles[username.(string)]
        if role == "" { role = "viewer" }
        s.authMu.Unlock()
        for _, r := range allowedRoles { if r == role { c.Next(); return } }
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"insufficient permissions"})
    }
}

type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Role     string `json:"role"`
}

func (s *Server) postCreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return }
    if req.Role == "" { req.Role = "viewer" }
    if req.Role != "admin" && req.Role != "editor" && req.Role != "viewer" {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid role (must be admin/editor/viewer)"}); return
    }
    s.authMu.Lock()
    if _, exists := s.users[req.Username]; exists {
        s.authMu.Unlock()
        c.JSON(http.StatusConflict, gin.H{"error":"user already exists"}); return
    }
    s.users[req.Username] = req.Password
    s.userRoles[req.Username] = req.Role
    s.authMu.Unlock()
    c.JSON(http.StatusCreated, gin.H{"username": req.Username, "role": req.Role})
}


