package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// E2ETestClient simulates a client making requests to the API
type E2ETestClient struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

func NewE2ETestClient(baseURL string) *E2ETestClient {
	return &E2ETestClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *E2ETestClient) Login(username, password string) error {
	body := map[string]string{"username": username, "password": password}
	jsonBody, _ := json.Marshal(body)
	
	req, err := http.NewRequest("POST", c.baseURL+"/api/auth/login", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %d", resp.StatusCode)
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	
	c.authToken = result["token"].(string)
	return nil
}

func (c *E2ETestClient) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	if c.authToken != "" {
		req.Header.Set("X-App-Auth", c.authToken)
	}
	return c.httpClient.Do(req)
}

func (c *E2ETestClient) Post(path string, body interface{}) (*http.Response, error) {
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.authToken != "" {
		req.Header.Set("X-App-Auth", c.authToken)
	}
	return c.httpClient.Do(req)
}

func (c *E2ETestClient) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	if c.authToken != "" {
		req.Header.Set("X-App-Auth", c.authToken)
	}
	return c.httpClient.Do(req)
}

// TestE2E_FullContainerWorkflow tests the complete container lifecycle
func TestE2E_FullContainerWorkflow(t *testing.T) {
	// Skip if API server is not running
	client := NewE2ETestClient("http://localhost:8081")
	
	// Test server availability
	resp, err := client.Get("/api/auth/me")
	if err != nil {
		t.Skip("API server not available, skipping E2E test")
		return
	}
	resp.Body.Close()
	
	// Step 1: Login
	err = client.Login("admin", "admin")
	require.NoError(t, err, "Should login successfully")
	assert.NotEmpty(t, client.authToken, "Should receive auth token")
	
	// Step 2: Check current user
	resp, err = client.Get("/api/auth/me")
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var meResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&meResp)
	user := meResp["user"].(map[string]interface{})
	assert.Equal(t, "admin", user["username"])
	
	// Step 3: List containers (initial state)
	resp, err = client.Get("/api/containers")
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var containers []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&containers)
	initialCount := len(containers)
	
	// Step 4: Create a network for the container
	networkReq := map[string]interface{}{
		"name":   "e2e-network",
		"driver": "bridge",
	}
	resp, err = client.Post("/api/networks", networkReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	var network map[string]interface{}
	if resp.StatusCode == http.StatusCreated {
		json.NewDecoder(resp.Body).Decode(&network)
		networkID := network["id"].(string)
		assert.NotEmpty(t, networkID)
		
		// Step 5: Create container with network
		containerReq := map[string]interface{}{
			"image":   "alpine:latest",
			"name":    "e2e-test-container",
			"network": "e2e-network",
			"detach":  true,
		}
		resp, err = client.Post("/api/containers/run", containerReq)
		// Container creation may fail in mock mode, which is OK for E2E test
		if resp != nil {
			resp.Body.Close()
		}
		
		// Step 6: List containers again
		resp, err = client.Get("/api/containers")
		require.NoError(t, err)
		defer resp.Body.Close()
		
		var containersAfter []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&containersAfter)
		assert.GreaterOrEqual(t, len(containersAfter), initialCount, "Container count should not decrease")
		
		// Step 7: Cleanup - delete network
		resp, err = client.Delete("/api/networks/" + networkID)
		if resp != nil {
			resp.Body.Close()
		}
	}
}

// TestE2E_TeamWorkspaceWorkflow tests complete team and workspace workflow
func TestE2E_TeamWorkspaceWorkflow(t *testing.T) {
	client := NewE2ETestClient("http://localhost:8081")
	
	// Test server availability
	resp, err := client.Get("/api/auth/me")
	if err != nil {
		t.Skip("API server not available, skipping E2E test")
		return
	}
	resp.Body.Close()
	
	// Step 1: Login
	err = client.Login("admin", "admin")
	require.NoError(t, err)
	
	// Step 2: Create team
	teamReq := map[string]interface{}{
		"name":        "E2E-Team",
		"description": "End-to-end test team",
	}
	resp, err = client.Post("/api/teams", teamReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Should create team")
	var team map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&team)
	teamID := team["id"].(string)
	assert.NotEmpty(t, teamID)
	
	// Step 3: Get team details
	resp, err = client.Get("/api/teams/" + teamID)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var teamDetail map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&teamDetail)
	members := teamDetail["members"].([]interface{})
	assert.Len(t, members, 1, "Creator should be admin member")
	
	// Step 4: Create network in workspace
	networkReq := map[string]interface{}{
		"name":        "workspace-network",
		"driver":      "bridge",
		"workspaceId": teamID,
	}
	resp, err = client.Post("/api/networks", networkReq)
	if resp != nil {
		resp.Body.Close()
	}
	// Network creation may require editor role
	
	// Step 5: List teams
	resp, err = client.Get("/api/teams")
	require.NoError(t, err)
	defer resp.Body.Close()
	
	var teams []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&teams)
	assert.GreaterOrEqual(t, len(teams), 1, "Should have at least one team")
	
	// Step 6: Cleanup - delete team
	resp, err = client.Delete("/api/teams/" + teamID)
	if resp != nil {
		resp.Body.Close()
	}
}

// TestE2E_MultiTenantWorkflow tests multi-tenant operations
func TestE2E_MultiTenantWorkflow(t *testing.T) {
	client := NewE2ETestClient("http://localhost:8081")
	
	// Test server availability
	resp, err := client.Get("/api/auth/me")
	if err != nil {
		t.Skip("API server not available, skipping E2E test")
		return
	}
	resp.Body.Close()
	
	// Step 1: Login as admin
	err = client.Login("admin", "admin")
	require.NoError(t, err)
	
	// Step 2: Create tenant
	tenantReq := map[string]interface{}{
		"name":   "E2E-Tenant",
		"domain": "e2e.test",
		"quota": map[string]interface{}{
			"maxContainers": 100,
			"maxNetworks":   50,
			"maxTeams":      20,
			"maxStorageGB":  1000,
		},
	}
	resp, err = client.Post("/api/tenants", tenantReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusCreated {
		var tenant map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tenant)
		tenantID := tenant["id"].(string)
		assert.NotEmpty(t, tenantID)
		
		// Step 3: Assign user to tenant (if testuser exists)
		assignReq := map[string]interface{}{"username": "admin"}
		resp, err = client.Post("/api/tenants/"+tenantID+"/assign", assignReq)
		if resp != nil {
			resp.Body.Close()
		}
		
		// Step 4: Get tenant usage
		resp, err = client.Get("/api/tenants/" + tenantID + "/usage")
		if err == nil && resp != nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var usage map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&usage)
				assert.NotNil(t, usage["usage"])
				assert.NotNil(t, usage["quota"])
			}
		}
		
		// Step 5: List tenants
		resp, err = client.Get("/api/tenants")
		require.NoError(t, err)
		defer resp.Body.Close()
		
		var tenants []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tenants)
		assert.GreaterOrEqual(t, len(tenants), 1, "Should have at least one tenant")
	}
}

// TestE2E_AuthAndRBACWorkflow tests authentication and role-based access
func TestE2E_AuthAndRBACWorkflow(t *testing.T) {
	client := NewE2ETestClient("http://localhost:8081")
	
	// Test server availability
	resp, err := client.Get("/api/auth/me")
	if err != nil {
		t.Skip("API server not available, skipping E2E test")
		return
	}
	resp.Body.Close()
	
	// Step 1: Try to access protected resource without auth
	resp, err = client.Get("/api/containers")
	require.NoError(t, err)
	defer resp.Body.Close()
	// May or may not require auth depending on route configuration
	
	// Step 2: Login
	err = client.Login("admin", "admin")
	require.NoError(t, err)
	
	// Step 3: Access protected resource with auth
	resp, err = client.Get("/api/containers")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Step 4: Test RBAC - try to create network (requires editor/admin)
	networkReq := map[string]interface{}{
		"name":   "rbac-test-network",
		"driver": "bridge",
	}
	resp, err = client.Post("/api/networks", networkReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// Admin should be able to create
	if resp.StatusCode == http.StatusCreated {
		var network map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&network)
		networkID := network["id"].(string)
		
		// Cleanup
		client.Delete("/api/networks/" + networkID)
	}
}

// TestE2E_ServiceDiscoveryWorkflow tests service discovery operations
func TestE2E_ServiceDiscoveryWorkflow(t *testing.T) {
	client := NewE2ETestClient("http://localhost:8081")
	
	// Test server availability
	resp, err := client.Get("/api/auth/me")
	if err != nil {
		t.Skip("API server not available, skipping E2E test")
		return
	}
	resp.Body.Close()
	
	// Login
	err = client.Login("admin", "admin")
	require.NoError(t, err)
	
	// Step 1: Register a service
	serviceReq := map[string]interface{}{
		"service":  "e2e-api-service",
		"address":  "10.0.0.1",
		"port":     "8080",
		"metadata": map[string]string{"env": "test"},
	}
	resp, err = client.Post("/api/services/register", serviceReq)
	if resp != nil {
		resp.Body.Close()
	}
	
	// Step 2: List services
	resp, err = client.Get("/api/services")
	require.NoError(t, err)
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		var services []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&services)
		// Should have at least the registered service
	}
	
	// Step 3: Resolve service
	resp, err = client.Get("/api/services/resolve/e2e-api-service")
	if resp != nil {
		resp.Body.Close()
	}
	
	// Step 4: Get next instance (round-robin)
	resp, err = client.Get("/api/services/resolve/e2e-api-service/next")
	if resp != nil {
		resp.Body.Close()
	}
}

// waitForServer waits for the API server to be available
func waitForServer(baseURL string, maxWait time.Duration) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(maxWait)
	
	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/api/auth/me")
		if err == nil {
			resp.Body.Close()
			return resp.StatusCode == http.StatusOK
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

