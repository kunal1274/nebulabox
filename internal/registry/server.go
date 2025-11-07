package registry

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

type Server struct {
    router   *gin.Engine
    port     string
    store    Storage            // Storage interface (can be memory or file-based)
    auth     *AuthConfig
    baseDir  string
}

// Storage interface for registry storage
type Storage interface {
    ListRepositories() []string
    ListTags(repo string) []string
    GetBlob(repo, digest string) (*blobData, bool)
    PutBlob(repo, digest string, data []byte) error
    GetManifest(repo, reference string) (*manifestData, bool)
    PutManifest(repo, reference, digest, mediaType string, content []byte) error
    DeleteTag(repo, tag string) bool
    StartUpload(uuid string)
    AppendUpload(uuid string, chunk []byte) error
    FinalizeUpload(repo, uuid, digest, createdBy string, metadata map[string]string, description string) (int64, error)
    GetVersionMetadata(repo string) *VersionMetadata
    SetVersionMetadata(repo string, vm *VersionMetadata) error
}

func NewServer() *Server {
    r := gin.Default()

    port := os.Getenv("NEBULABOX_REGISTRY_PORT")
    if port == "" { port = "5001" }
    
    baseDir := os.Getenv("NEBULABOX_REGISTRY_STORAGE")
    if baseDir == "" { baseDir = "./registry-storage" }
    
    // Use file-based storage if directory is specified, otherwise memory
    var store Storage
    if fileStore, err := NewFileStore(baseDir); err == nil {
        store = fileStore
        log.Printf("📁 Using file-based storage: %s", baseDir)
    } else {
        log.Printf("⚠️  Failed to initialize file storage, using in-memory: %v", err)
        store = newMemoryStore()
    }
    
    auth := NewAuthConfig()
    auth.AddUser("admin", "admin", []string{"admin", "push", "pull"})
    
    s := &Server{
        router:  r,
        port:    port,
        store:   store,
        auth:    auth,
        baseDir: baseDir,
    }
    s.setupRoutes()
    return s
}

func (s *Server) setupRoutes() {
    // Catch-all dispatcher to support multi-segment repositories and root
    s.router.Any("/v2/*path", s.dispatchV2)
    
    // Auth endpoints
    s.router.POST("/auth/login", s.login)
    s.router.POST("/auth/token", s.generateToken)
    
    // Registry API endpoints
    api := s.router.Group("/api/registry")
    {
        api.GET("/repositories", s.listRepositories)
        api.GET("/repositories/:repo/versions", s.listVersions)
        api.GET("/repositories/:repo/versions/:tag", s.getVersion)
        api.DELETE("/repositories/:repo/versions/:tag", s.deleteVersion)
        api.GET("/repositories/:repo/summary", s.getRepositorySummary)
    }
}

func (s *Server) Start() error {
    log.Printf("🚀 Nebula Registry starting on :%s", s.port)
    return s.router.Run(":" + s.port)
}

func (s *Server) dispatchV2(c *gin.Context) {
    p := strings.TrimPrefix(c.Param("path"), "/")
    method := c.Request.Method

    // Root
    if p == "" {
        c.Status(http.StatusOK)
        return
    }

    // Catalog
    if p == "_catalog" && method == http.MethodGet {
        c.JSON(http.StatusOK, gin.H{"repositories": s.store.ListRepositories()})
        return
    }

    // Tags list: <repo>/tags/list
    if strings.HasSuffix(p, "/tags/list") && method == http.MethodGet {
        repo := strings.TrimSuffix(p, "/tags/list")
        repo = strings.TrimSuffix(repo, "/")
        c.JSON(http.StatusOK, gin.H{"name": repo, "tags": s.store.ListTags(repo)})
        return
    }

    // Blobs uploads
    if strings.HasSuffix(p, "/blobs/uploads/") && method == http.MethodPost {
        if !s.requireAuth(c) { return }
        repo := strings.TrimSuffix(p, "/blobs/uploads/")
        uuid := genUUID()
        s.store.StartUpload(uuid)
        c.Header("Location", "/v2/"+repo+"/blobs/uploads/"+uuid)
        c.Header("Docker-Upload-UUID", uuid)
        c.Status(http.StatusAccepted)
        return
    }
    if strings.Contains(p, "/blobs/uploads/") && method == http.MethodPatch {
        if !s.requireAuth(c) { return }
        idx := strings.LastIndex(p, "/blobs/uploads/")
        uuid := p[idx+len("/blobs/uploads/"):]
        b, _ := io.ReadAll(c.Request.Body)
        _ = s.store.AppendUpload(uuid, b)
        c.Status(http.StatusNoContent)
        return
    }
    if strings.Contains(p, "/blobs/uploads/") && method == http.MethodPut {
        if !s.requireAuth(c) { return }
        idx := strings.LastIndex(p, "/blobs/uploads/")
        repo := p[:idx]
        uuid := p[idx+len("/blobs/uploads/"):]
        digest := c.Query("digest")
        if digest == "" { c.Status(http.StatusBadRequest); return }
        
        username := s.getUsernameFromAuth(c)
        metadata := make(map[string]string)
        description := c.Query("description")
        
        size, err := s.store.FinalizeUpload(repo, uuid, digest, username, metadata, description)
        if err != nil {
            c.Status(http.StatusNotFound); return
        }
        c.Header("Location", "/v2/"+repo+"/blobs/"+digest)
        c.Header("Content-Length", fmt.Sprintf("%d", size))
        c.Status(http.StatusCreated)
        return
    }

    // Blobs GET: <repo>/blobs/<digest>
    if strings.Contains(p, "/blobs/") && method == http.MethodGet {
        idx := strings.LastIndex(p, "/blobs/")
        repo := p[:idx]
        digest := p[idx+len("/blobs/"):]
        if b, ok := s.store.GetBlob(repo, digest); ok {
            c.Data(http.StatusOK, "application/octet-stream", b.content)
            return
        }
        c.Status(http.StatusNotFound)
        return
    }

    // Manifests GET/PUT: <repo>/manifests/<reference>
    if strings.Contains(p, "/manifests/") {
        idx := strings.LastIndex(p, "/manifests/")
        repo := p[:idx]
        ref := p[idx+len("/manifests/"):]
        if method == http.MethodGet {
            if mf, ok := s.store.GetManifest(repo, ref); ok {
                mt := mf.mediaType
                if mt == "" { mt = "application/vnd.oci.image.manifest.v1+json" }
                c.Data(http.StatusOK, mt, mf.content)
                return
            }
            c.Status(http.StatusNotFound)
            return
        }
        if method == http.MethodPut {
            if !s.requireAuth(c) { return }
            b, _ := io.ReadAll(c.Request.Body)
            digest := c.Query("digest")
            if digest == "" { digest = "sha256:dummy" }
            mediaType := c.GetHeader("Content-Type")
            
            // Get username from auth
            username := s.getUsernameFromAuth(c)
            
            // Store manifest
            s.store.PutManifest(repo, ref, digest, mediaType, b)
            
            // Update version metadata
            vm := s.store.GetVersionMetadata(repo)
            metadata := map[string]string{
                "mediaType": mediaType,
                "digest":    digest,
            }
            vm.AddVersion(ref, digest, username, int64(len(b)), metadata, "")
            s.store.SetVersionMetadata(repo, vm)
            
            c.Status(http.StatusCreated)
            return
        }
        if method == http.MethodDelete {
            if !s.requireAuth(c) { return }
            // Treat DELETE on tag as untag, not deleting manifest content
            if len(ref) > 0 && (len(ref) < 7 || ref[:7] != "sha256:") {
                if s.store.DeleteTag(repo, ref) {
                    c.Status(http.StatusAccepted)
                    return
                }
                c.Status(http.StatusNotFound)
                return
            }
            c.Status(http.StatusMethodNotAllowed)
            return
        }
    }

    c.Status(http.StatusNotFound)
}

func genUUID() string { return time.Now().Format("20060102150405.000000000") }

func (s *Server) login(c *gin.Context) {
    var body struct{
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"})
        return
    }
    
    user, err := s.auth.AuthenticateUser(body.Username, body.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"})
        return
    }
    
    scopes := []string{"pull", "push"}
    if user.HasRole("admin") {
        scopes = append(scopes, "*")
    }
    
    token := s.auth.GenerateToken(body.Username, scopes)
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "token_type": "Bearer",
        "expires_in": int(s.auth.TokenExpiry.Seconds()),
        "scope": strings.Join(scopes, " "),
    })
}

func (s *Server) generateToken(c *gin.Context) {
    var body struct{
        Username string   `json:"username"`
        Password string   `json:"password"`
        Scopes   []string `json:"scopes,omitempty"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"})
        return
    }
    
    user, err := s.auth.AuthenticateUser(body.Username, body.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"})
        return
    }
    
    scopes := body.Scopes
    if len(scopes) == 0 {
        scopes = []string{"pull", "push"}
        if user.HasRole("admin") {
            scopes = append(scopes, "*")
        }
    }
    
    token := s.auth.GenerateToken(body.Username, scopes)
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "token_type": "Bearer",
        "expires_in": int(s.auth.TokenExpiry.Seconds()),
        "scope": strings.Join(scopes, " "),
    })
}

func (s *Server) requireAuth(c *gin.Context) bool {
    auth := c.GetHeader("Authorization")
    
    // Try Bearer token first
    if token, ok := ParseBearerAuth(auth); ok {
        if _, err := s.auth.ValidateToken(token); err == nil {
            return true
        }
    }
    
    // Try Basic auth
    if username, password, ok := ParseBasicAuth(auth); ok {
        if _, err := s.auth.AuthenticateUser(username, password); err == nil {
            return true
        }
    }
    
    c.Header("WWW-Authenticate", "Bearer realm=\"registry\",service=\"nebulabox-registry\"")
    c.Status(http.StatusUnauthorized)
    return false
}

func (s *Server) getUsernameFromAuth(c *gin.Context) string {
    auth := c.GetHeader("Authorization")
    
    // Try Bearer token
    if token, ok := ParseBearerAuth(auth); ok {
        if info, err := s.auth.ValidateToken(token); err == nil {
            return info.Username
        }
    }
    
    // Try Basic auth
    if username, _, ok := ParseBasicAuth(auth); ok {
        return username
    }
    
    return "unknown"
}


