package api

import (
    "context"
    "log"
    "net/http"
    "sort"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/nebulabox/nebulabox/internal/database/mongodb_repositories"
)

type perfSample struct {
    ts     time.Time
    dur    time.Duration
    status int
}

type perfStats struct {
    mu      sync.Mutex
    samples []perfSample
}

func (p *perfStats) add(d time.Duration, status int) {
    p.mu.Lock()
    p.samples = append(p.samples, perfSample{ts: time.Now(), dur: d, status: status})
    // trim older than 10 minutes
    cutoff := time.Now().Add(-10 * time.Minute)
    i := 0
    for ; i < len(p.samples); i++ {
        if p.samples[i].ts.After(cutoff) { break }
    }
    if i > 0 && i < len(p.samples) { p.samples = p.samples[i:] } else if i >= len(p.samples) { p.samples = p.samples[:0] }
    p.mu.Unlock()
}

func (p *perfStats) snapshot(window time.Duration) (reqps float64, p95ms float64, errRate float64) {
    p.mu.Lock(); defer p.mu.Unlock()
    cutoff := time.Now().Add(-window)
    var arr []time.Duration
    var errs int
    var total int
    for _, s := range p.samples {
        if s.ts.After(cutoff) {
            total++
            arr = append(arr, s.dur)
            if s.status >= 500 { errs++ }
        }
    }
    if total > 0 {
        reqps = float64(total) / window.Seconds()
        if errs > 0 { errRate = float64(errs) / float64(total) }
        sort.Slice(arr, func(i,j int) bool { return arr[i] < arr[j] })
        idx := int(float64(len(arr))*0.95) - 1
        if idx < 0 { idx = 0 }
        if idx >= len(arr) { idx = len(arr)-1 }
        p95ms = float64(arr[idx].Milliseconds())
    }
    return
}

type PerfMetrics struct {
    ReqPerSec1m  float64 `json:"reqPerSec1m"`
    P95Ms1m      float64 `json:"p95Ms1m"`
    ErrorRate1m  float64 `json:"errorRate1m"`
    ReqPerSec5m  float64 `json:"reqPerSec5m"`
    P95Ms5m      float64 `json:"p95Ms5m"`
    ErrorRate5m  float64 `json:"errorRate5m"`
    Timestamp    int64   `json:"timestamp"`
}

func (s *Server) perfMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        d := time.Since(start)
        status := c.Writer.Status()
        s.perf.add(d, status)
        
        // Track per-endpoint metrics
        endpoint := c.Request.Method + " " + c.FullPath()
        if endpoint == "" {
            endpoint = c.Request.Method + " " + c.Request.URL.Path
        }
        
        if s.endpointMetrics != nil {
            metrics := s.endpointMetrics.GetOrCreate(endpoint)
            metrics.RecordRequest(d, status >= 500)
        }
        
        // Save to MongoDB if available (async to avoid blocking)
        if s.mongoRepos != nil && s.mongoRepos.APIMetrics != nil {
            go func() {
                ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
                defer cancel()
                
                // Get user ID if available
                userID := ""
                if username, ok := c.Get("username"); ok {
                    userID = username.(string)
                }
                
                // Get request/response sizes
                requestSize := int64(0)
                if c.Request.Body != nil {
                    // Request size is approximate
                    requestSize = c.Request.ContentLength
                }
                
                responseSize := int64(0)
                if size := c.Writer.Size(); size > 0 {
                    responseSize = int64(size)
                }
                
                entry := &mongodb_repositories.APIMetricEntry{
                    Timestamp:   time.Now(),
                    Endpoint:    endpoint,
                    Method:      c.Request.Method,
                    StatusCode:  status,
                    DurationMs:  d.Milliseconds(),
                    RequestSize: requestSize,
                    ResponseSize: responseSize,
                    Error:       "",
                    UserID:      userID,
                }
                
                if status >= 400 {
                    // Try to extract error message
                    if c.Writer.Status() >= 500 {
                        entry.Error = "server error"
                    } else {
                        entry.Error = "client error"
                    }
                }
                
                if err := s.mongoRepos.APIMetrics.Insert(ctx, entry); err != nil {
                    log.Printf("⚠️  Warning: Failed to save API metric to MongoDB: %v", err)
                }
            }()
        }
    }
}

func (s *Server) getPerfMetrics(c *gin.Context) {
    r1, p1, e1 := s.perf.snapshot(1 * time.Minute)
    r5, p5, e5 := s.perf.snapshot(5 * time.Minute)
    c.JSON(http.StatusOK, PerfMetrics{
        ReqPerSec1m: r1, P95Ms1m: p1, ErrorRate1m: e1,
        ReqPerSec5m: r5, P95Ms5m: p5, ErrorRate5m: e5,
        Timestamp: time.Now().Unix(),
    })
}

func (s *Server) streamPerf(c *gin.Context) {
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
            r1, p1, e1 := s.perf.snapshot(1 * time.Minute)
            r5, p5, e5 := s.perf.snapshot(5 * time.Minute)
            m := PerfMetrics{ReqPerSec1m: r1, P95Ms1m: p1, ErrorRate1m: e1, ReqPerSec5m: r5, P95Ms5m: p5, ErrorRate5m: e5, Timestamp: time.Now().Unix()}
            c.SSEvent("perf", m)
            c.Writer.Flush()
        }
    }
}


