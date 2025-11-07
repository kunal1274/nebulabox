package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestServer creates a test API server with fresh state
func setupTestServer(t *testing.T) *Server {
	gin.SetMode(gin.TestMode)
	
	// Create a minimal containerd client (mock mode)
	client, err := containerd.NewClient()
	if err != nil {
		t.Fatalf("Failed to create containerd client: %v", err)
	}
	
	server := &Server{
		containerd:         client,
		users:            map[string]string{"testuser": "testpass", "admin": "admin"},
		userRoles:        map[string]string{"testuser": "editor", "admin": "admin"},
		sessions:         make(map[string]string),
		teams:            make(map[string]*Team),
		teamMembers:     make(map[string]map[string]*TeamMember),
		tenants:          make(map[string]*Tenant),
		userTenants:      make(map[string]string),
		networks:         make(map[string]*Network),
		containerWorkspaces: make(map[string]string),
		networkWorkspaces:   make(map[string]string),
	}
	
	return server
}

// getAuthToken logs in and returns an auth token
func getAuthToken(t *testing.T, server *Server, username, password string) string {
	router := gin.New()
	router.POST("/auth/login", server.postAuthLogin)
	
	body := map[string]string{"username": username, "password": password}
	jsonBody, _ := json.Marshal(body)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	require.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	token, ok := response["token"].(string)
	require.True(t, ok)
	return token
}

func TestIntegration_AuthFlow(t *testing.T) {
	server := setupTestServer(t)
	router := gin.New()
	router.POST("/auth/login", server.postAuthLogin)
	router.GET("/auth/me", server.getAuthMe)
	router.POST("/auth/logout", server.postAuthLogout)
	
	// Step 1: Login
	body := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["token"].(string)
	assert.NotEmpty(t, token)
	
	// Step 2: Check /auth/me
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/auth/me", nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	var meResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &meResp)
	user := meResp["user"].(map[string]interface{})
	assert.Equal(t, "testuser", user["username"])
	assert.Equal(t, "editor", user["role"])
	
	// Step 3: Logout
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Step 4: Verify token no longer works
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/auth/me", nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	var meResp2 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &meResp2)
	assert.Nil(t, meResp2["user"])
}

func TestIntegration_TeamLifecycle(t *testing.T) {
	server := setupTestServer(t)
	router := gin.New()
	router.Use(server.authMiddleware())
	router.POST("/teams", server.createTeam)
	router.GET("/teams", server.listTeams)
	router.GET("/teams/:id", server.getTeam)
	router.POST("/teams/:id/invite", server.inviteMember)
	router.DELETE("/teams/:id", server.deleteTeam)
	
	// Get auth token
	token := getAuthToken(t, server, "testuser", "testpass")
	
	// Step 1: Create team
	body := createTeamRequest{Name: "IntegrationTeam", Description: "Test team"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/teams", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	var team Team
	json.Unmarshal(w.Body.Bytes(), &team)
	teamID := team.ID
	assert.Equal(t, "IntegrationTeam", team.Name)
	
	// Step 2: List teams
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/teams", nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	var teams []Team
	json.Unmarshal(w.Body.Bytes(), &teams)
	assert.Len(t, teams, 1)
	
	// Step 3: Get team details
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/teams/"+teamID, nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	var teamDetail map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &teamDetail)
	members := teamDetail["members"].([]interface{})
	assert.Len(t, members, 1)
	
	// Step 4: Invite member (admin can invite)
	inviteBody := inviteMemberRequest{Username: "admin", Role: "editor"}
	jsonBody, _ = json.Marshal(inviteBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/teams/"+teamID+"/invite", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Step 5: Delete team
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/teams/"+teamID, nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TenantLifecycle(t *testing.T) {
	server := setupTestServer(t)
	router := gin.New()
	router.Use(server.authMiddleware())
	router.Use(server.requireRole("admin"))
	router.POST("/tenants", server.createTenant)
	router.GET("/tenants", server.listTenants)
	router.POST("/tenants/:id/assign", server.assignUserToTenant)
	router.GET("/tenants/:id/usage", server.getTenantUsage)
	
	// Get admin token
	adminToken := getAuthToken(t, server, "admin", "admin")
	
	// Step 1: Create tenant
	body := createTenantRequest{
		Name:   "IntegrationTenant",
		Domain: "test.com",
		Quota: TenantQuota{
			MaxContainers: 50,
			MaxNetworks:   25,
			MaxTeams:      10,
			MaxStorageGB:  500,
		},
	}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Auth", adminToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	var tenant Tenant
	json.Unmarshal(w.Body.Bytes(), &tenant)
	tenantID := tenant.ID
	assert.Equal(t, "IntegrationTenant", tenant.Name)
	
	// Step 2: Assign user to tenant
	assignBody := map[string]string{"username": "testuser"}
	jsonBody, _ = json.Marshal(assignBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/tenants/"+tenantID+"/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Auth", adminToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, tenantID, server.userTenants["testuser"])
	
	// Step 3: Get tenant usage
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/tenants/"+tenantID+"/usage", nil)
	req.Header.Set("X-App-Auth", adminToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	var usageResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &usageResp)
	assert.NotNil(t, usageResp["usage"])
	assert.NotNil(t, usageResp["quota"])
}

func TestIntegration_NetworkOperations(t *testing.T) {
	server := setupTestServer(t)
	router := gin.New()
	router.Use(server.authMiddleware())
	router.Use(server.requireRole("admin", "editor"))
	router.POST("/networks", server.createNetwork)
	router.GET("/networks", server.listNetworks)
	router.DELETE("/networks/:id", server.deleteNetwork)
	
	// Get auth token
	token := getAuthToken(t, server, "testuser", "testpass")
	
	// Step 1: Create network
	body := createNetworkRequest{Name: "test-network", Driver: "bridge", Subnet: "10.0.0.0/24"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/networks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	var network Network
	json.Unmarshal(w.Body.Bytes(), &network)
	networkID := network.ID
	assert.Equal(t, "test-network", network.Name)
	
	// Step 2: List networks
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/networks", nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	var networks []Network
	json.Unmarshal(w.Body.Bytes(), &networks)
	assert.Len(t, networks, 1)
	
	// Step 3: Delete network
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/networks/"+networkID, nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Step 4: Verify deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/networks", nil)
	req.Header.Set("X-App-Auth", token)
	router.ServeHTTP(w, req)
	
	json.Unmarshal(w.Body.Bytes(), &networks)
	assert.Len(t, networks, 0)
}

func TestIntegration_RBAC_AccessControl(t *testing.T) {
	server := setupTestServer(t)
	
	// Create viewer user
	server.users["viewer"] = "viewerpass"
	server.userRoles["viewer"] = "viewer"
	
	router := gin.New()
	router.Use(server.authMiddleware())
	router.POST("/networks", server.requireRole("admin", "editor"), server.createNetwork)
	
	// Test viewer cannot create network
	viewerToken := getAuthToken(t, server, "viewer", "viewerpass")
	body := createNetworkRequest{Name: "forbidden-network"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/networks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Auth", viewerToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	// Test editor can create network
	editorToken := getAuthToken(t, server, "testuser", "testpass")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/networks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Auth", editorToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
}

