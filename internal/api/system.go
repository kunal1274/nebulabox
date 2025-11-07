package api

import (
	"context"
	"log"
    "net/http"
    "runtime"
    "time"
    "strconv"

    "github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/database/mongodb_repositories"
)

// internal sample for history
type systemSample struct {
    ts       int64
    cpu      float64
    mem      float64
    disk     float64
    running  int
    total    int
}

type systemHistory struct {
    buf   []systemSample
    head  int
    size  int
}

func newSystemHistory(capacity int) *systemHistory {
    return &systemHistory{buf: make([]systemSample, capacity)}
}

// add stores the latest sample into a ring buffer
func (h *systemHistory) add(s systemSample) {
    if len(h.buf) == 0 { return }
    h.buf[h.head%len(h.buf)] = s
    h.head++
    if h.size < len(h.buf) { h.size++ }
}

// query returns samples within the last rangeSec seconds, downsampled by stepSec
func (h *systemHistory) query(rangeSec, stepSec int) []systemSample {
    if h.size == 0 { return nil }
    if stepSec <= 0 { stepSec = 1 }
    now := time.Now().Unix()
    from := now - int64(rangeSec)
    out := make([]systemSample, 0, h.size)
    lastEmitted := int64(0)
    // walk newest->oldest through ring
    for i := 0; i < h.size; i++ {
        idx := (h.head - 1 - i)
        if idx < 0 { idx = idx + ((-idx/len(h.buf))+1)*len(h.buf) }
        s := h.buf[idx%len(h.buf)]
        if s.ts < from { break }
        if lastEmitted == 0 || s.ts <= now && (lastEmitted-s.ts) >= int64(stepSec) {
            out = append(out, s)
            lastEmitted = s.ts
        }
    }
    // reverse to chronological
    for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
        out[i], out[j] = out[j], out[i]
    }
    return out
}

// SystemStatsResponse represents system statistics
type SystemStatsResponse struct {
	CPUUsage          float64 `json:"cpuUsage"`
	MemoryUsage       float64 `json:"memoryUsage"`
	DiskUsage         float64 `json:"diskUsage"`
	ContainersRunning int     `json:"containersRunning"`
	ContainersTotal   int     `json:"containersTotal"`
	Timestamp         int64   `json:"timestamp"`
}

// getSystemStats handles GET /api/system/stats
func (s *Server) getSystemStats(c *gin.Context) {
	// Get container statistics
	ctx := c.Request.Context()
	containers, err := s.containerd.ListContainers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get container stats",
			"details": err.Error(),
		})
		return
	}

	runningCount := 0
	for _, container := range containers {
		if container.Status == "running" {
			runningCount++
		}
	}

	// Get system memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate memory usage percentage (simplified)
	memoryUsage := float64(m.Alloc) / float64(m.Sys) * 100
	if memoryUsage > 100 {
		memoryUsage = 100
	}

	// Mock CPU and disk usage for now
	// In a real implementation, you would use system monitoring libraries
	cpuUsage := 45.2 + float64(time.Now().Unix()%10) // Simulate varying CPU usage
	diskUsage := 38.5 + float64(time.Now().Unix()%5) // Simulate varying disk usage

	stats := SystemStatsResponse{
		CPUUsage:          cpuUsage,
		MemoryUsage:       memoryUsage,
		DiskUsage:         diskUsage,
		ContainersRunning: runningCount,
		ContainersTotal:   len(containers),
		Timestamp:         time.Now().Unix(),
	}

	// store in history (per second resolution)
	if s.sysHist != nil {
		ss := systemSample{ts: stats.Timestamp, cpu: stats.CPUUsage, mem: stats.MemoryUsage, disk: stats.DiskUsage, running: stats.ContainersRunning, total: stats.ContainersTotal}
		s.sysHist.add(ss)
	}

	// Save to MongoDB if available (async to avoid blocking)
	if s.mongoRepos != nil && s.mongoRepos.SystemMetrics != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			entry := &mongodb_repositories.SystemMetricEntry{
				Timestamp:         time.Unix(stats.Timestamp, 0),
				CPUUsage:          stats.CPUUsage,
				MemoryUsage:       stats.MemoryUsage,
				DiskUsage:         stats.DiskUsage,
				ContainersRunning: stats.ContainersRunning,
				ContainersTotal:   stats.ContainersTotal,
			}
			
			if err := s.mongoRepos.SystemMetrics.Insert(ctx, entry); err != nil {
				log.Printf("⚠️  Warning: Failed to save system metric to MongoDB: %v", err)
			}
		}()
	}

	c.JSON(http.StatusOK, stats)
}

// getSystemStream streams system stats via SSE: GET /api/system/stream
func (s *Server) getSystemStream(c *gin.Context) {
    c.Writer.Header().Set("Content-Type", "text/event-stream")
    c.Writer.Header().Set("Cache-Control", "no-cache")
    c.Writer.Header().Set("Connection", "keep-alive")
    c.Writer.Header().Set("X-Accel-Buffering", "no")

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    // send initial
    send := func() bool {
        // reuse the logic from getSystemStats
        ctx := c.Request.Context()
        containers, err := s.containerd.ListContainers(ctx)
        if err != nil {
            return false
        }
        runningCount := 0
        for _, container := range containers { if container.Status == "running" { runningCount++ } }
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        memoryUsage := float64(m.Alloc) / float64(m.Sys) * 100
        if memoryUsage > 100 { memoryUsage = 100 }
        cpuUsage := 40.0 + float64(time.Now().UnixNano()%2000)/100.0
        diskUsage := 35.0 + float64(time.Now().UnixNano()%1500)/100.0
		stats := SystemStatsResponse{
            CPUUsage: cpuUsage,
            MemoryUsage: memoryUsage,
            DiskUsage: diskUsage,
            ContainersRunning: runningCount,
            ContainersTotal: len(containers),
            Timestamp: time.Now().Unix(),
        }
		// store in history as well
		if s.sysHist != nil {
			ss := systemSample{ts: stats.Timestamp, cpu: stats.CPUUsage, mem: stats.MemoryUsage, disk: stats.DiskUsage, running: stats.ContainersRunning, total: stats.ContainersTotal}
			s.sysHist.add(ss)
		}
        c.SSEvent("stats", stats)
        c.Writer.Flush()
        return true
    }

    // initial push
    if !send() { return }

    for {
        select {
        case <-c.Request.Context().Done():
            return
        case <-ticker.C:
            if !send() { return }
        }
    }
}

// getSystemHistory: GET /api/system/history?range=3600&step=10
func (s *Server) getSystemHistory(c *gin.Context) {
    // defaults: last 1 hour, 10s step
    rangeSec := 3600
    stepSec := 10
    if v := c.Query("range"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 { rangeSec = n }
    }
    if v := c.Query("step"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 { stepSec = n }
    }
    if s.sysHist == nil {
        c.JSON(http.StatusOK, gin.H{"points": []any{}})
        return
    }
    samples := s.sysHist.query(rangeSec, stepSec)
    type point struct {
        Ts int64 `json:"ts"`
        Cpu float64 `json:"cpu"`
        Mem float64 `json:"mem"`
        Disk float64 `json:"disk"`
        Running int `json:"running"`
        Total int `json:"total"`
    }
    out := make([]point, 0, len(samples))
    for _, s := range samples {
        out = append(out, point{Ts: s.ts, Cpu: s.cpu, Mem: s.mem, Disk: s.disk, Running: s.running, Total: s.total})
    }
    c.JSON(http.StatusOK, gin.H{"points": out})
}
