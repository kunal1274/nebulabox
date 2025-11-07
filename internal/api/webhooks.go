package api

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
    "github.com/gin-gonic/gin"
)

type GitHubEvent struct {
    ID        string `json:"id"`
    Event     string `json:"event"`
    Repo      string `json:"repo"`
    Ref       string `json:"ref"`
    Action    string `json:"action,omitempty"`
    Sender    string `json:"sender,omitempty"`
    Timestamp int64  `json:"timestamp"`
}

func (s *Server) verifyGitHubSignature(body []byte, sigHeader string) bool {
    if s.githubSecret == "" { return true }
    if !strings.HasPrefix(sigHeader, "sha256=") { return false }
    sig := strings.TrimPrefix(sigHeader, "sha256=")
    mac := hmac.New(sha256.New, []byte(s.githubSecret))
    mac.Write(body)
    expected := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(sig), []byte(expected))
}

func (s *Server) postGitHubWebhook(c *gin.Context) {
    event := c.GetHeader("X-GitHub-Event")
    delivery := c.GetHeader("X-GitHub-Delivery")
    sig := c.GetHeader("X-Hub-Signature-256")

    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"read body"}); return }

    if !s.verifyGitHubSignature(body, sig) {
        c.JSON(http.StatusUnauthorized, gin.H{"error":"bad signature"}); return
    }

    // Parse minimal fields
    var payload map[string]any
    _ = json.Unmarshal(body, &payload)
    repo := ""
    if r, ok := payload["repository"].(map[string]any); ok {
        if n, ok := r["full_name"].(string); ok { repo = n }
    }
    ref := ""
    if v, ok := payload["ref"].(string); ok { ref = v }
    action := ""
    if v, ok := payload["action"].(string); ok { action = v }
    sender := ""
    if sdr, ok := payload["sender"].(map[string]any); ok {
        if n, ok := sdr["login"].(string); ok { sender = n }
    }

    ev := GitHubEvent{ID: delivery, Event: event, Repo: repo, Ref: ref, Action: action, Sender: sender, Timestamp: time.Now().Unix()}
    s.ghMu.Lock()
    s.ghEvents = append([]GitHubEvent{ev}, s.ghEvents...)
    if len(s.ghEvents) > 100 { s.ghEvents = s.ghEvents[:100] }
    s.ghMu.Unlock()

    // enqueue build on push events
    if event == "push" && repo != "" {
        s.enqueueBuild("github", repo, ref)
    }

    // TODO: trigger CI/CD pipeline, build, etc.
    c.JSON(http.StatusOK, gin.H{"status":"ok"})
}

func (s *Server) getGitHubEvents(c *gin.Context) {
    s.ghMu.Lock()
    out := append([]GitHubEvent(nil), s.ghEvents...)
    s.ghMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"events": out})
}


