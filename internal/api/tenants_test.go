package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		tenants:     make(map[string]*Tenant),
		userTenants: make(map[string]string),
		userRoles:   map[string]string{"admin": "admin"},
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "admin") })
	router.POST("/tenants", server.authMiddleware(), server.requireRole("admin"), server.createTenant)
	
	body := createTenantRequest{
		Name:   "acme-corp",
		Domain: "acme.com",
		Quota: TenantQuota{
			MaxContainers: 100,
			MaxNetworks:   50,
			MaxTeams:      20,
			MaxStorageGB:  1000,
		},
	}
	jsonBody, _ := json.Marshal(body)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var tenant Tenant
	err := json.Unmarshal(w.Body.Bytes(), &tenant)
	assert.NoError(t, err)
	assert.Equal(t, "acme-corp", tenant.Name)
	assert.Equal(t, "acme.com", tenant.Domain)
	assert.Equal(t, 100, tenant.Quota.MaxContainers)
}

func TestCheckTenantQuota(t *testing.T) {
	server := &Server{
		tenants: map[string]*Tenant{
			"tenant-1": {
				ID:   "tenant-1",
				Quota: TenantQuota{MaxContainers: 2, MaxNetworks: 1, MaxTeams: 1},
			},
		},
		containerWorkspaces: map[string]string{"c1": "team-1"},
		networkWorkspaces:   map[string]string{"n1": "tenant-1"},
		teams: map[string]*Team{
			"team-1": {CreatedBy: "user1"},
		},
		userTenants: map[string]string{"user1": "tenant-1"},
	}
	
	// Test quota check - should pass
	assert.True(t, server.checkTenantQuota("tenant-1", "container"))
	assert.True(t, server.checkTenantQuota("tenant-1", "network"))
	assert.True(t, server.checkTenantQuota("tenant-1", "team"))
	
	// Add more resources to exceed quota
	server.containerWorkspaces["c2"] = "team-1"
	server.containerWorkspaces["c3"] = "team-1"
	
	// Should fail for containers now
	assert.False(t, server.checkTenantQuota("tenant-1", "container"))
}

func TestListTenants(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		tenants: map[string]*Tenant{
			"tenant-1": {ID: "tenant-1", Name: "TestTenant"},
		},
		userTenants: map[string]string{"user1": "tenant-1"},
		userRoles:   map[string]string{"admin": "admin", "user1": "editor"},
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "user1") })
	router.GET("/tenants", server.authMiddleware(), server.listTenants)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tenants", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var tenants []Tenant
	err := json.Unmarshal(w.Body.Bytes(), &tenants)
	assert.NoError(t, err)
	assert.Len(t, tenants, 1)
	assert.Equal(t, "TestTenant", tenants[0].Name)
}

func TestAssignUserToTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		tenants: map[string]*Tenant{
			"tenant-1": {ID: "tenant-1", Name: "TestTenant"},
		},
		users:      map[string]string{"testuser": "pass"},
		userRoles:  map[string]string{"admin": "admin"},
		userTenants: make(map[string]string),
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "admin") })
	router.POST("/tenants/:id/assign", server.authMiddleware(), server.requireRole("admin"), server.assignUserToTenant)
	
	body := map[string]string{"username": "testuser"}
	jsonBody, _ := json.Marshal(body)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tenants/tenant-1/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "tenant-1", server.userTenants["testuser"])
}

