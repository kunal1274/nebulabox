package api

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type AlertThresholds struct {
    CPUHigh    float64 `json:"cpuHigh"`
    MemoryHigh float64 `json:"memoryHigh"`
    DiskHigh   float64 `json:"diskHigh"`
}

type AlertEvent struct {
    Type      string  `json:"type"`      // cpu|memory|disk
    Value     float64 `json:"value"`
    Threshold float64 `json:"threshold"`
    Timestamp int64   `json:"timestamp"`
}

// getAlertsConfig returns current thresholds
func (s *Server) getAlertsConfig(c *gin.Context) {
    s.alertsMu.Lock()
    cfg := s.alerts
    s.alertsMu.Unlock()
    c.JSON(http.StatusOK, cfg)
}

// postAlertsConfig updates thresholds
func (s *Server) postAlertsConfig(c *gin.Context) {
    var req AlertThresholds
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body","details": err.Error()})
        return
    }
    s.alertsMu.Lock()
    if req.CPUHigh > 0 { s.alerts.CPUHigh = req.CPUHigh }
    if req.MemoryHigh > 0 { s.alerts.MemoryHigh = req.MemoryHigh }
    if req.DiskHigh > 0 { s.alerts.DiskHigh = req.DiskHigh }
    cfg := s.alerts
    s.alertsMu.Unlock()
    c.JSON(http.StatusOK, cfg)
}

// streamAlerts emits alerts when thresholds exceeded using SSE
func (s *Server) streamAlerts(c *gin.Context) {
    c.Writer.Header().Set("Content-Type", "text/event-stream")
    c.Writer.Header().Set("Cache-Control", "no-cache")
    c.Writer.Header().Set("Connection", "keep-alive")
    c.Writer.Header().Set("X-Accel-Buffering", "no")

    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-c.Request.Context().Done():
            return
        case <-ticker.C:
            // reuse system stats collection
            e := s.currentAlerts(c)
            for _, a := range e {
                c.SSEvent("alert", a)
            }
            c.Writer.Flush()
        }
    }
}

func (s *Server) currentAlerts(c *gin.Context) []AlertEvent {
    ctx := c.Request.Context()
    containers, _ := s.containerd.ListContainers(ctx)
    _ = containers // reserved for future per-container alerts
    // approximate system stats from existing method
    now := time.Now().Unix()
    // We call getSystemStats logic inline (simplified) by using perf/system mocks
    // For simplicity, emit random-like based on time like existing mocks
    cpu := 40.0 + float64(time.Now().UnixNano()%4000)/100.0
    mem := 50.0 + float64(time.Now().UnixNano()%3000)/100.0
    dsk := 30.0 + float64(time.Now().UnixNano()%2000)/100.0

    s.alertsMu.Lock()
    cfg := s.alerts
    s.alertsMu.Unlock()

    var out []AlertEvent
    if cfg.CPUHigh > 0 && cpu >= cfg.CPUHigh {
        out = append(out, AlertEvent{Type:"cpu", Value: cpu, Threshold: cfg.CPUHigh, Timestamp: now})
    }
    if cfg.MemoryHigh > 0 && mem >= cfg.MemoryHigh {
        out = append(out, AlertEvent{Type:"memory", Value: mem, Threshold: cfg.MemoryHigh, Timestamp: now})
    }
    if cfg.DiskHigh > 0 && dsk >= cfg.DiskHigh {
        out = append(out, AlertEvent{Type:"disk", Value: dsk, Threshold: cfg.DiskHigh, Timestamp: now})
    }
    return out
}


