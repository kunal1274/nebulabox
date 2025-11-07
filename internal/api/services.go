package api

import (
    "net/http"
    "time"
    "fmt"
    "github.com/gin-gonic/gin"
)

type ServiceInstance struct {
    ID        string  `json:"id"`
    Name      string  `json:"name"`
    Address   string  `json:"address"`
    Port      int     `json:"port"`
    Version   string  `json:"version,omitempty"`
    Network   string  `json:"network,omitempty"`
    CreatedAt int64   `json:"createdAt"`
}

type RegisterServiceRequest struct {
    Name    string `json:"name" binding:"required"`
    ID      string `json:"id"`
    Address string `json:"address"`
    Port    int    `json:"port"`
    Version string `json:"version"`
    Network string `json:"network"`
}

type DeregisterServiceRequest struct {
    Name string `json:"name" binding:"required"`
    ID   string `json:"id" binding:"required"`
}

func (s *Server) listServices(c *gin.Context) {
    s.svcMu.Lock()
    out := make(map[string][]ServiceInstance, len(s.services))
    for k, v := range s.services { out[k] = append([]ServiceInstance(nil), v...) }
    s.svcMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"services": out})
}

func (s *Server) resolveService(c *gin.Context) {
    name := c.Param("name")
    s.svcMu.Lock()
    list := append([]ServiceInstance(nil), s.services[name]...)
    s.svcMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"name": name, "instances": list})
}

func (s *Server) resolveServiceNext(c *gin.Context) {
    name := c.Param("name")
    s.svcMu.Lock()
    list := s.services[name]
    if len(list) == 0 { s.svcMu.Unlock(); c.JSON(http.StatusOK, gin.H{"name": name, "instance": nil}); return }
    idx := s.svcRR[name] % len(list)
    s.svcRR[name] = (idx + 1) % len(list)
    inst := list[idx]
    s.svcMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"name": name, "instance": inst})
}

func (s *Server) registerService(c *gin.Context) {
    var req RegisterServiceRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body", "details": err.Error()}); return }
    if req.ID == "" { req.ID = fmt.Sprintf("svc-%d", time.Now().UnixNano()) }
    inst := ServiceInstance{ID: req.ID, Name: req.Name, Address: req.Address, Port: req.Port, Version: req.Version, Network: req.Network, CreatedAt: time.Now().Unix()}
    s.svcMu.Lock()
    s.services[req.Name] = append(s.services[req.Name], inst)
    s.svcMu.Unlock()
    c.JSON(http.StatusCreated, inst)
}

func (s *Server) deregisterService(c *gin.Context) {
    var req DeregisterServiceRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body", "details": err.Error()}); return }
    s.svcMu.Lock()
    arr := s.services[req.Name]
    out := arr[:0]
    for _, x := range arr { if x.ID != req.ID { out = append(out, x) } }
    if len(out) == 0 { delete(s.services, req.Name) } else { s.services[req.Name] = out }
    s.svcMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"removed": req.ID})
}


