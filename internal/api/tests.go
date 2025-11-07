package api

import (
    "net/http"
    "time"
    "math/rand"
    "github.com/gin-gonic/gin"
)

type TestRun struct {
    ID        string `json:"id"`
    Repo      string `json:"repo"`
    Ref       string `json:"ref"`
    Suite     string `json:"suite"`
    Status    string `json:"status"` // queued|running|passed|failed
    StartedAt int64  `json:"startedAt"`
    EndedAt   int64  `json:"endedAt"`
}

type RunTestsRequest struct {
    Repo  string `json:"repo" binding:"required"`
    Ref   string `json:"ref"`
    Suite string `json:"suite"`
}

func (s *Server) listTests(c *gin.Context) {
    s.testsMu.Lock()
    out := append([]TestRun(nil), s.tests...)
    s.testsMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"tests": out})
}

func (s *Server) enqueueTest(repo, ref, suite string) TestRun {
    if suite == "" { suite = "default" }
    tr := TestRun{ID: generateID("t"), Repo: repo, Ref: ref, Suite: suite, Status: "queued", StartedAt: time.Now().Unix()}
    s.testsMu.Lock()
    s.tests = append([]TestRun{tr}, s.tests...)
    if len(s.tests) > 200 { s.tests = s.tests[:200] }
    s.testsMu.Unlock()
    go func(id string) {
        s.updateTestStatus(id, "running")
        time.Sleep(1500 * time.Millisecond)
        // randomly pass/fail but bias to pass
        if rand.Intn(100) < 85 { s.updateTestStatus(id, "passed") } else { s.updateTestStatus(id, "failed") }
    }(tr.ID)
    return tr
}

func (s *Server) updateTestStatus(id, status string) {
    s.testsMu.Lock()
    defer s.testsMu.Unlock()
    for i := range s.tests {
        if s.tests[i].ID == id {
            s.tests[i].Status = status
            if status == "passed" || status == "failed" { s.tests[i].EndedAt = time.Now().Unix() }
            if status == "passed" {
                repo := s.tests[i].Repo
                ref := s.tests[i].Ref
                go s.enqueueDeployment(repo, ref, "staging")
            }
            break
        }
    }
}

func generateID(prefix string) string {
    return prefix + "-" + time.Now().Format("20060102T150405.000000000")
}

func (s *Server) postRunTests(c *gin.Context) {
    var req RunTestsRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return }
    tr := s.enqueueTest(req.Repo, req.Ref, req.Suite)
    c.JSON(http.StatusCreated, tr)
}


