package api

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
)

type Deployment struct {
    ID        string `json:"id"`
    Repo      string `json:"repo"`
    Ref       string `json:"ref"`
    Env       string `json:"env"`
    Status    string `json:"status"` // queued|deploying|success|failed
    StartedAt int64  `json:"startedAt"`
    EndedAt   int64  `json:"endedAt"`
}

type TriggerDeploymentRequest struct {
    Repo string `json:"repo" binding:"required"`
    Ref  string `json:"ref"`
    Env  string `json:"env"`
}

func (s *Server) listDeployments(c *gin.Context) {
    s.deployMu.Lock()
    out := append([]Deployment(nil), s.deployments...)
    s.deployMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"deployments": out})
}

func (s *Server) enqueueDeployment(repo, ref, env string) Deployment {
    if env == "" { env = "staging" }
    d := Deployment{ID: generateID("d"), Repo: repo, Ref: ref, Env: env, Status: "queued", StartedAt: time.Now().Unix()}
    s.deployMu.Lock()
    s.deployments = append([]Deployment{d}, s.deployments...)
    if len(s.deployments) > 200 { s.deployments = s.deployments[:200] }
    s.deployMu.Unlock()
    go func(id string) {
        s.updateDeploymentStatus(id, "deploying")
        time.Sleep(2 * time.Second)
        s.updateDeploymentStatus(id, "success")
    }(d.ID)
    return d
}

func (s *Server) updateDeploymentStatus(id, status string) {
    s.deployMu.Lock()
    defer s.deployMu.Unlock()
    for i := range s.deployments {
        if s.deployments[i].ID == id {
            s.deployments[i].Status = status
            if status == "success" || status == "failed" { s.deployments[i].EndedAt = time.Now().Unix() }
            break
        }
    }
}

func (s *Server) postTriggerDeployment(c *gin.Context) {
    var req TriggerDeploymentRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return }
    d := s.enqueueDeployment(req.Repo, req.Ref, req.Env)
    c.JSON(http.StatusCreated, d)
}


