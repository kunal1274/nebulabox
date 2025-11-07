package api

import (
    "bytes"
    "io"
    "net/http"
    "net/url"
    "strings"

    "github.com/gin-gonic/gin"
)
// postRegistryLogin proxies login to registry and returns token
func (s *Server) postRegistryLogin(c *gin.Context) {
    base, _ := url.Parse(s.registryURL)
    base.Path = "/auth/login"
    body, _ := io.ReadAll(c.Request.Body)
    req, _ := http.NewRequest(http.MethodPost, base.String(), bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil { c.JSON(http.StatusBadGateway, gin.H{"error":"registry unavailable","details": err.Error()}); return }
    defer resp.Body.Close()
    c.Status(resp.StatusCode)
    io.Copy(c.Writer, resp.Body)
}

// getRegistryCatalog proxies to /v2/_catalog
func (s *Server) getRegistryCatalog(c *gin.Context) {
    u, _ := url.Parse(s.registryURL)
    u.Path = "/v2/_catalog"
    req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
    if h := c.GetHeader("Authorization"); h != "" { req.Header.Set("Authorization", h) }
    resp, err := http.DefaultClient.Do(req)
    if err != nil { c.JSON(http.StatusBadGateway, gin.H{"error":"registry unavailable","details": err.Error()}); return }
    defer resp.Body.Close()
    c.Status(resp.StatusCode)
    io.Copy(c.Writer, resp.Body)
}

// getRegistryTags proxies to /v2/:repo/tags/list (supports multi-segment repo via wildcard)
func (s *Server) getRegistryTags(c *gin.Context) {
    repo := c.Param("path")
    u, _ := url.Parse(s.registryURL)
    u.Path = "/v2/" + repo + "/tags/list"
    req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
    if h := c.GetHeader("Authorization"); h != "" { req.Header.Set("Authorization", h) }
    resp, err := http.DefaultClient.Do(req)
    if err != nil { c.JSON(http.StatusBadGateway, gin.H{"error":"registry unavailable","details": err.Error()}); return }
    defer resp.Body.Close()
    c.Status(resp.StatusCode)
    io.Copy(c.Writer, resp.Body)
}

// retagRequest: JSON { repo, source, target }
type retagRequest struct {
    Repo   string `json:"repo" binding:"required"`
    Source string `json:"source" binding:"required"`
    Target string `json:"target" binding:"required"`
}

// postRegistryRetag copies manifest from source tag/digest to target tag
func (s *Server) postRegistryRetag(c *gin.Context) {
    var req retagRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body","details": err.Error()}); return }
    base, _ := url.Parse(s.registryURL)
    // GET manifest of source
    getURL := *base; getURL.Path = "/v2/" + req.Repo + "/manifests/" + req.Source
    reqGet, _ := http.NewRequest(http.MethodGet, getURL.String(), nil)
    if h := c.GetHeader("Authorization"); h != "" { reqGet.Header.Set("Authorization", h) }
    resp, err := http.DefaultClient.Do(reqGet)
    if err != nil { c.JSON(http.StatusBadGateway, gin.H{"error":"registry unavailable","details": err.Error()}); return }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        c.Status(resp.StatusCode); io.Copy(c.Writer, resp.Body); return
    }
    body, _ := io.ReadAll(resp.Body)
    contentType := resp.Header.Get("Content-Type")
    if contentType == "" { contentType = "application/vnd.oci.image.manifest.v1+json" }
    // PUT manifest to target
    putURL := *base; putURL.Path = "/v2/" + req.Repo + "/manifests/" + req.Target
    reqPut, _ := http.NewRequest(http.MethodPut, putURL.String(), bytes.NewReader(body))
    reqPut.Header.Set("Content-Type", contentType)
    if h := c.GetHeader("Authorization"); h != "" { reqPut.Header.Set("Authorization", h) }
    resp2, err := http.DefaultClient.Do(reqPut)
    if err != nil { c.JSON(http.StatusBadGateway, gin.H{"error":"registry unavailable","details": err.Error()}); return }
    defer resp2.Body.Close()
    c.Status(resp2.StatusCode)
    io.Copy(c.Writer, resp2.Body)
}

// deleteRegistryTag deletes a tag reference (untag)
// Route: DELETE /api/registry/tags/*path where path is "repo/tag"
func (s *Server) deleteRegistryTag(c *gin.Context) {
    fullPath := c.Param("path")
    // Parse path as "repo/tag" format
    // Remove leading slash if present
    if len(fullPath) > 0 && fullPath[0] == '/' {
        fullPath = fullPath[1:]
    }
    // Split into repo and tag (last segment is tag, rest is repo)
    parts := strings.Split(fullPath, "/")
    if len(parts) < 2 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path format, expected repo/tag"})
        return
    }
    tag := parts[len(parts)-1]
    repo := strings.Join(parts[:len(parts)-1], "/")
    
    base, _ := url.Parse(s.registryURL)
    del := *base; del.Path = "/v2/" + repo + "/manifests/" + tag
    req, _ := http.NewRequest(http.MethodDelete, del.String(), nil)
    if h := c.GetHeader("Authorization"); h != "" { req.Header.Set("Authorization", h) }
    resp, err := http.DefaultClient.Do(req)
    if err != nil { c.JSON(http.StatusBadGateway, gin.H{"error":"registry unavailable","details": err.Error()}); return }
    defer resp.Body.Close()
    c.Status(resp.StatusCode)
    io.Copy(c.Writer, resp.Body)
}


