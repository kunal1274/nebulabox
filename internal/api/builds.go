package api

import (
    "net/http"
    "time"
    "fmt"
    "github.com/gin-gonic/gin"
)

type Build struct {
    ID        string `json:"id"`
    Source    string `json:"source"` // github/gitlab/manual
    Repo      string `json:"repo"`
    Ref       string `json:"ref"`
    Status    string `json:"status"` // queued|running|success|failed
    StartedAt int64  `json:"startedAt"`
    EndedAt   int64  `json:"endedAt"`
}

type TriggerBuildRequest struct {
    Source string `json:"source"`
    Repo   string `json:"repo" binding:"required"`
    Ref    string `json:"ref"`
}

func (s *Server) listBuilds(c *gin.Context) {
    s.buildsMu.Lock()
    out := append([]Build(nil), s.builds...)
    s.buildsMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"builds": out})
}

func (s *Server) enqueueBuild(source, repo, ref string) Build {
    b := Build{ID: fmt.Sprintf("b-%d", time.Now().UnixNano()), Source: source, Repo: repo, Ref: ref, Status: "queued", StartedAt: time.Now().Unix()}
    s.buildsMu.Lock()
    s.builds = append([]Build{b}, s.builds...)
    if len(s.builds) > 200 { s.builds = s.builds[:200] }
    s.buildsMu.Unlock()
    go func(id string) {
        // simulate build work
        s.updateBuildStatus(id, "running")
        time.Sleep(2 * time.Second)
        s.updateBuildStatus(id, "success")
    }(b.ID)
    return b
}

func (s *Server) updateBuildStatus(id, status string) {
    s.buildsMu.Lock()
    defer s.buildsMu.Unlock()
    for i := range s.builds {
        if s.builds[i].ID == id {
            s.builds[i].Status = status
            if status == "success" || status == "failed" { s.builds[i].EndedAt = time.Now().Unix() }
            // auto-run tests on success
            if status == "success" {
                repo := s.builds[i].Repo
                ref := s.builds[i].Ref
                go s.enqueueTest(repo, ref, "default")
            }
            break
        }
    }
}

func (s *Server) postTriggerBuild(c *gin.Context) {
    var req TriggerBuildRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return }
    b := s.enqueueBuild(req.Source, req.Repo, req.Ref)
    c.JSON(http.StatusCreated, b)
}


