package api

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

type dnsRecordRequest struct {
    Name string   `json:"name" binding:"required"`
    A    []string `json:"a"`
}

func (s *Server) listDNSRecords(c *gin.Context) {
    s.dnsMu.Lock()
    out := make(map[string][]string, len(s.dnsRecords))
    for k, v := range s.dnsRecords { out[k] = append([]string(nil), v...) }
    s.dnsMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"records": out})
}

func (s *Server) addDNSRecord(c *gin.Context) {
    var req dnsRecordRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body", "details": err.Error()}); return }
    name := strings.ToLower(strings.TrimSpace(req.Name))
    if name == "" { c.JSON(http.StatusBadRequest, gin.H{"error":"name required"}); return }
    s.dnsMu.Lock()
    s.dnsRecords[name] = append([]string(nil), req.A...)
    s.dnsMu.Unlock()
    c.JSON(http.StatusCreated, gin.H{"name": name, "a": req.A})
}

func (s *Server) deleteDNSRecord(c *gin.Context) {
    name := strings.ToLower(strings.TrimSpace(c.Param("name")))
    if name == "" { c.JSON(http.StatusBadRequest, gin.H{"error":"name required"}); return }
    s.dnsMu.Lock()
    delete(s.dnsRecords, name)
    s.dnsMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"deleted": name})
}

// dnsResolve resolves name from explicit records or service registry.
// Names like "api" or "api.svc" are resolved from service registry.
func (s *Server) dnsResolve(c *gin.Context) {
    name := strings.ToLower(strings.TrimSpace(c.Param("name")))
    // explicit A records first
    s.dnsMu.Lock()
    if ips, ok := s.dnsRecords[name]; ok {
        s.dnsMu.Unlock()
        c.JSON(http.StatusOK, gin.H{"name": name, "a": ips})
        return
    }
    s.dnsMu.Unlock()
    // service registry fallback (strip .svc/.local)
    svc := name
    if strings.HasSuffix(svc, ".svc") { svc = strings.TrimSuffix(svc, ".svc") }
    if strings.HasSuffix(svc, ".local") { svc = strings.TrimSuffix(svc, ".local") }
    s.svcMu.Lock()
    list := append([]ServiceInstance(nil), s.services[svc]...)
    s.svcMu.Unlock()
    // mock addresses if not set; default 127.0.0.1
    out := make([]string, 0, len(list))
    for _, inst := range list {
        addr := inst.Address
        if addr == "" { addr = "127.0.0.1" }
        out = append(out, addr)
    }
    c.JSON(http.StatusOK, gin.H{"name": name, "a": out})
}


