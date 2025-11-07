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

func TestListNetworks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		networks: map[string]*Network{
			"net-1": {ID: "net-1", Name: "testnet", Driver: "bridge"},
		},
	}
	
	router := gin.New()
	router.GET("/networks", server.listNetworks)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/networks", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var networks []Network
	err := json.Unmarshal(w.Body.Bytes(), &networks)
	assert.NoError(t, err)
	assert.Len(t, networks, 1)
	assert.Equal(t, "testnet", networks[0].Name)
}

func TestCreateNetwork(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		networks:         make(map[string]*Network),
		userTenants:     make(map[string]string),
		teams:           make(map[string]*Team),
		teamMembers:     make(map[string]map[string]*TeamMember),
		networkWorkspaces: make(map[string]string),
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "testuser") })
	router.POST("/networks", server.authMiddleware(), server.requireRole("admin", "editor"), server.createNetwork)
	
	body := createNetworkRequest{Name: "mynet", Driver: "bridge"}
	jsonBody, _ := json.Marshal(body)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/networks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var network Network
	err := json.Unmarshal(w.Body.Bytes(), &network)
	assert.NoError(t, err)
	assert.Equal(t, "mynet", network.Name)
	assert.Equal(t, "bridge", network.Driver)
}

func TestDeleteNetwork(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		networks: map[string]*Network{
			"net-1": {ID: "net-1", Name: "testnet"},
		},
	}
	
	router := gin.New()
	router.DELETE("/networks/:id", server.deleteNetwork)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/networks/net-1", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, server.networks, "net-1")
}

