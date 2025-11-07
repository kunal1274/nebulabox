package api

import (
    "context"
    "log"
    "net/http"
    "sort"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/nebulabox/nebulabox/internal/database/mongodb_repositories"
)

type LogEntry struct {
    Timestamp int64  `json:"timestamp"`
    Container string `json:"container"`
    Level     string `json:"level"`
    Message   string `json:"message"`
}

// GET /api/logs/search?query=...&containerId=...&since=unix&until=unix&limit=200
func (s *Server) searchLogs(c *gin.Context) {
    queryStr := strings.ToLower(c.Query("query"))
    containerID := c.Query("containerId")
    since, _ := timeFromQuery(c.Query("since"))
    until, _ := timeFromQuery(c.Query("until"))
    limit := 200
    if v := c.Query("limit"); v != "" { if n, err := atoi(v); err == nil { limit = n } }

    ctx := c.Request.Context()
    var out []LogEntry

    // Try to get logs from MongoDB first
    if s.mongoRepos != nil && s.mongoRepos.ContainerLogs != nil {
        filters := mongodb_repositories.ContainerLogFilters{
            ContainerID: containerID,
            Query:       queryStr,
            Since:       since,
            Until:       until,
            Limit:       limit,
        }
        
        dbLogs, err := s.mongoRepos.ContainerLogs.Search(ctx, filters)
        if err == nil && len(dbLogs) > 0 {
            // Convert MongoDB logs to API format
            for _, dbLog := range dbLogs {
                out = append(out, LogEntry{
                    Timestamp: dbLog.Timestamp.Unix(),
                    Container: dbLog.ContainerID,
                    Level:     dbLog.Level,
                    Message:   dbLog.Message,
                })
            }
            c.JSON(http.StatusOK, out)
            return
        }
        // Log error but continue with fallback
        if err != nil {
            log.Printf("⚠️  Warning: Failed to search logs from MongoDB: %v", err)
        }
    }

    // Fallback: synthesize logs from running containers (mock mode)
    containers, _ := s.containerd.ListContainers(ctx)
    now := time.Now().Unix()
    for _, ctr := range containers {
        if containerID != "" && ctr.ID != containerID { continue }
        for i := 0; i < 50; i++ {
            ts := now - int64(50-i)
            lvl := []string{"INFO","WARN","ERROR"}[i%3]
            msg := ctr.Name + " handled request /health" 
            e := LogEntry{Timestamp: ts, Container: ctr.ID, Level: lvl, Message: msg}
            if queryStr != "" && !strings.Contains(strings.ToLower(e.Message), queryStr) { continue }
            if !since.IsZero() && ts < since.Unix() { continue }
            if !until.IsZero() && ts > until.Unix() { continue }
            out = append(out, e)
        }
    }
    sort.Slice(out, func(i,j int) bool { return out[i].Timestamp > out[j].Timestamp })
    if len(out) > limit { out = out[:limit] }
    c.JSON(http.StatusOK, out)
}

// GET /api/logs/stream?containerId=...
func (s *Server) streamLogs(c *gin.Context) {
    containerID := c.Query("containerId")
    c.Writer.Header().Set("Content-Type", "text/event-stream")
    c.Writer.Header().Set("Cache-Control", "no-cache")
    c.Writer.Header().Set("Connection", "keep-alive")
    c.Writer.Header().Set("X-Accel-Buffering", "no")
    ticker := time.NewTicker(1000 * time.Millisecond)
    defer ticker.Stop()
    for {
        select {
        case <-c.Request.Context().Done():
            return
        case t := <-ticker.C:
            // emit a mock line
            msg := "processing request"
            e := LogEntry{Timestamp: t.Unix(), Container: chooseContainerID(s, containerID), Level: "INFO", Message: msg}
            c.SSEvent("log", e)
            c.Writer.Flush()
        }
    }
}

func chooseContainerID(s *Server, prefer string) string {
    if prefer != "" { return prefer }
    ctrs, _ := s.containerd.ListContainers(cancellableBackground())
    if len(ctrs) > 0 { return ctrs[0].ID }
    return ""
}

func cancellableBackground() context.Context { return context.Background() }

func timeFromQuery(v string) (time.Time, error) {
    if v == "" { return time.Time{}, nil }
    if n, err := atoi(v); err == nil { return time.Unix(int64(n), 0), nil }
    return time.Time{}, nil
}

func atoi(s string) (int, error) {
    var n int
    for i := 0; i < len(s); i++ {
        c := s[i] - '0'
        if c > 9 { return 0, fmtError() }
        n = n*10 + int(c)
    }
    return n, nil
}

func fmtError() error { return &parseErr{}
}
type parseErr struct{}
func (*parseErr) Error() string { return "parse error" }


