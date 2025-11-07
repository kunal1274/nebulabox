package api

import (
    "net/http"
    "sort"
    "strconv"
    "github.com/gin-gonic/gin"
)

type reservePortRequest struct {
    Port int    `json:"port"`
    ID   string `json:"id"`
}

type releasePortRequest struct {
    Port int    `json:"port"`
    ID   string `json:"id"`
}

func (s *Server) listPorts(c *gin.Context) {
    s.portsMu.Lock()
    out := make([]gin.H, 0, len(s.ports))
    for p, id := range s.ports { out = append(out, gin.H{"port": p, "id": id}) }
    s.portsMu.Unlock()
    sort.Slice(out, func(i, j int) bool { return out[i]["port"].(int) < out[j]["port"].(int) })
    c.JSON(http.StatusOK, gin.H{"ports": out})
}

func (s *Server) reservePort(c *gin.Context) {
    var req reservePortRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return }
    if req.Port <= 0 || req.Port > 65535 || req.ID == "" { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid input"}); return }
    s.portsMu.Lock()
    defer s.portsMu.Unlock()
    if _, used := s.ports[req.Port]; used { c.JSON(http.StatusConflict, gin.H{"error":"port in use"}); return }
    s.ports[req.Port] = req.ID
    c.JSON(http.StatusCreated, gin.H{"port": req.Port, "id": req.ID})
}

func (s *Server) releasePort(c *gin.Context) {
    var req releasePortRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return }
    s.portsMu.Lock()
    if owner, ok := s.ports[req.Port]; ok {
        if req.ID == "" || owner == req.ID { delete(s.ports, req.Port) }
    }
    s.portsMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"released": req.Port})
}

func (s *Server) suggestPort(c *gin.Context) {
    base := 30000
    if v := c.Query("from"); v != "" { if n, err := strconv.Atoi(v); err == nil && n > 0 { base = n } }
    if base < 1024 { base = 1024 }
    s.portsMu.Lock()
    p := base
    for {
        if _, used := s.ports[p]; !used { break }
        p++
        if p > 65535 { p = 1024 }
        if p == base { break }
    }
    s.portsMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"port": p})
}


