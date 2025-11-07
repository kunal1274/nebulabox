package api

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
)

type Rollback struct {
    ID        string `json:"id"`
    Repo      string `json:"repo"`
    FromRef   string `json:"fromRef"`
    ToRef     string `json:"toRef"`
    Env       string `json:"env"`
    Status    string `json:"status"` // queued|rolling_back|success|failed
    StartedAt int64  `json:"startedAt"`
    EndedAt   int64  `json:"endedAt"`
}

type TriggerRollbackRequest struct {
    Repo string `json:"repo" binding:"required"`
    Env  string `json:"env" binding:"required"`
}

func (s *Server) listRollbacks(c *gin.Context) {
    s.rbMu.Lock()
    out := append([]Rollback(nil), s.rollbacks...)
    s.rbMu.Unlock()
    c.JSON(http.StatusOK, gin.H{"rollbacks": out})
}

// findPreviousSuccessful returns the most recent successful deployment for repo/env excluding the latest one
func (s *Server) findPreviousSuccessful(repo, env string) (current *Deployment, previous *Deployment) {
    s.deployMu.Lock()
    defer s.deployMu.Unlock()
    var list []Deployment
    for _, d := range s.deployments {
        if d.Repo == repo && d.Env == env && d.Status == "success" {
            list = append(list, d)
        }
    }
    if len(list) == 0 { return nil, nil }
    if len(list) == 1 { return &list[0], nil }
    // deployments are stored newest-first; list keeps that order as we appended from s.deployments (already newest-first)
    return &list[0], &list[1]
}

func (s *Server) enqueueRollback(repo, env, toRef, fromRef string) Rollback {
    rb := Rollback{ID: generateID("r"), Repo: repo, Env: env, ToRef: toRef, FromRef: fromRef, Status: "queued", StartedAt: time.Now().Unix()}
    s.rbMu.Lock()
    s.rollbacks = append([]Rollback{rb}, s.rollbacks...)
    if len(s.rollbacks) > 200 { s.rollbacks = s.rollbacks[:200] }
    s.rbMu.Unlock()
    go func(id string) {
        s.updateRollbackStatus(id, "rolling_back")
        time.Sleep(2 * time.Second)
        s.updateRollbackStatus(id, "success")
        // after successful rollback, enqueue a deployment to apply the target ref (simulate)
        s.enqueueDeployment(repo, toRef, env)
    }(rb.ID)
    return rb
}

func (s *Server) updateRollbackStatus(id, status string) {
    s.rbMu.Lock()
    defer s.rbMu.Unlock()
    for i := range s.rollbacks {
        if s.rollbacks[i].ID == id {
            s.rollbacks[i].Status = status
            if status == "success" || status == "failed" { s.rollbacks[i].EndedAt = time.Now().Unix() }
            break
        }
    }
}

func (s *Server) postTriggerRollback(c *gin.Context) {
    var req TriggerRollbackRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"}); return }
    cur, prev := s.findPreviousSuccessful(req.Repo, req.Env)
    if cur == nil || prev == nil {
        c.JSON(http.StatusConflict, gin.H{"error":"no previous successful deployment to roll back to"})
        return
    }
    rb := s.enqueueRollback(req.Repo, req.Env, prev.Ref, cur.Ref)
    c.JSON(http.StatusCreated, rb)
}


