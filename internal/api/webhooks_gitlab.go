package api

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
    "github.com/gin-gonic/gin"
)

type GitLabEvent struct {
    ID        string `json:"id"`
    Event     string `json:"event"`
    Project   string `json:"project"`
    Ref       string `json:"ref"`
    User      string `json:"user"`
    Timestamp int64  `json:"timestamp"`
}

func (s *Server) postGitLabWebhook(c *gin.Context) {
    token := c.GetHeader("X-Gitlab-Token")
    if s.gitlabSecret != "" && token != s.gitlabSecret {
        c.JSON(http.StatusUnauthorized, gin.H{"error":"bad token"}); return
    }
    event := c.GetHeader("X-Gitlab-Event")
    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"read body"}); return }
    var payload map[string]any
    _ = json.Unmarshal(body, &payload)
    proj := ""
    if p, ok := payload["project"].(map[string]any); ok {
        if n, ok := p["path_with_namespace"].(string); ok { proj = n }
    }
    ref := ""
    if v, ok := payload["ref"].(string); ok { ref = v }
    user := ""
    if v, ok := payload["user_username"].(string); ok { user = v }
    if user == "" {
        if u, ok := payload["user"].(map[string]any); ok { if n, ok := u["username"].(string); ok { user = n } }
    }
    id := ""
    if v, ok := payload["object_kind"].(string); ok { id = v }
    if v, ok := payload["event_id"].(string); ok { id = v }
    if id == "" { id = strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "-") }
    ev := GitLabEvent{ID: id, Event: event, Project: proj, Ref: ref, User: user, Timestamp: time.Now().Unix()}
    s.glMu.Lock()
    s.glEvents = append([]GitLabEvent{ev}, s.glEvents...)
    if len(s.glEvents) > 100 { s.glEvents = s.glEvents[:100] }
    s.glMu.Unlock()
    if strings.Contains(strings.ToLower(event), "push") && proj != "" {
        s.enqueueBuild("gitlab", proj, ref)
    }
    c.JSON(http.StatusOK, gin.H{"status":"ok"})
}

func (s *Server) getGitLabEvents(c *gin.Context) {
    s.glMu.Lock()
    out := append([]GitLabEvent(nil), s.glEvents...)
    s.glMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"events": out})
}


