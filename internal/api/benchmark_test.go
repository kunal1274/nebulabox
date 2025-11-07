package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
)

// BenchmarkAPI_ListContainers benchmarks container listing endpoint
func BenchmarkAPI_ListContainers(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	client, _ := containerd.NewClient()
	server := &Server{
		containerd: client,
		networks:   make(map[string]*Network),
	}
	
	router := gin.New()
	router.GET("/containers", server.listContainers)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/containers", nil)
		router.ServeHTTP(w, req)
	}
}

// BenchmarkAPI_ListNetworks benchmarks network listing endpoint
func BenchmarkAPI_ListNetworks(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		networks: make(map[string]*Network),
	}
	
	router := gin.New()
	router.GET("/networks", server.listNetworks)
	
	// Pre-populate with some networks
	for i := 0; i < 100; i++ {
		server.networks[fmt.Sprintf("net-%d", i)] = &Network{
			ID:     fmt.Sprintf("net-%d", i),
			Name:   fmt.Sprintf("network-%d", i),
			Driver: "bridge",
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/networks", nil)
		router.ServeHTTP(w, req)
	}
}

// BenchmarkAPI_ListTeams benchmarks team listing endpoint
func BenchmarkAPI_ListTeams(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		teams:       make(map[string]*Team),
		teamMembers: make(map[string]map[string]*TeamMember),
		users:       map[string]string{"user1": "pass"},
		userRoles:   map[string]string{"user1": "admin"},
		sessions:    map[string]string{"token1": "user1"},
	}
	
	// Pre-populate teams
	for i := 0; i < 50; i++ {
		teamID := fmt.Sprintf("team-%d", i)
		server.teams[teamID] = &Team{
			ID:   teamID,
			Name: fmt.Sprintf("Team %d", i),
		}
		server.teamMembers[teamID] = map[string]*TeamMember{
			"user1": {Username: "user1", Role: "admin"},
		}
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "user1") })
	router.Use(server.authMiddleware())
	router.GET("/teams", server.listTeams)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams", nil)
		req.Header.Set("X-App-Auth", "token1")
		router.ServeHTTP(w, req)
	}
}

// BenchmarkAPI_CreateNetwork benchmarks network creation
func BenchmarkAPI_CreateNetwork(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	client, _ := containerd.NewClient()
	server := &Server{
		containerd:         client,
		networks:          make(map[string]*Network),
		userTenants:       make(map[string]string),
		teams:             make(map[string]*Team),
		teamMembers:       make(map[string]map[string]*TeamMember),
		networkWorkspaces: make(map[string]string),
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "admin") })
	router.POST("/networks", server.authMiddleware(), server.requireRole("admin", "editor"), server.createNetwork)
	
	body := createNetworkRequest{Name: "test-network", Driver: "bridge"}
	jsonBody, _ := json.Marshal(body)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		reqBody := bytes.NewBuffer(jsonBody)
		req, _ := http.NewRequest("POST", "/networks", reqBody)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		// Clean up for next iteration
		server.networks = make(map[string]*Network)
	}
}

// BenchmarkAPI_ListTenants benchmarks tenant listing
func BenchmarkAPI_ListTenants(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		tenants:     make(map[string]*Tenant),
		userTenants: map[string]string{"admin": "tenant-1"},
		userRoles:   map[string]string{"admin": "admin"},
		sessions:    map[string]string{"token1": "admin"},
	}
	
	// Pre-populate tenants
	for i := 0; i < 20; i++ {
		server.tenants[fmt.Sprintf("tenant-%d", i)] = &Tenant{
			ID:   fmt.Sprintf("tenant-%d", i),
			Name: fmt.Sprintf("Tenant %d", i),
			Quota: TenantQuota{
				MaxContainers: 100,
				MaxNetworks:   50,
				MaxTeams:      20,
				MaxStorageGB:  1000,
			},
		}
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "admin") })
	router.Use(server.authMiddleware())
	router.GET("/tenants", server.listTenants)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tenants", nil)
		req.Header.Set("X-App-Auth", "token1")
		router.ServeHTTP(w, req)
	}
}

// BenchmarkAPI_AuthLogin benchmarks login endpoint
func BenchmarkAPI_AuthLogin(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		users:     map[string]string{"testuser": "testpass"},
		userRoles: map[string]string{"testuser": "editor"},
		sessions:  make(map[string]string),
	}
	
	router := gin.New()
	router.POST("/auth/login", server.postAuthLogin)
	
	body := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(body)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
}

