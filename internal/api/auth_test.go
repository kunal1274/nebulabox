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

func TestAuthLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		users:     map[string]string{"testuser": "testpass"},
		userRoles: map[string]string{"testuser": "editor"},
		sessions:  map[string]string{},
	}
	
	router := gin.New()
	router.POST("/auth/login", server.postAuthLogin)
	
	// Test successful login
	body := map[string]string{"username": "testuser", "password": "testpass"}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "editor")
	
	// Test invalid credentials
	body = map[string]string{"username": "testuser", "password": "wrongpass"}
	jsonBody, _ = json.Marshal(body)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		sessions:  map[string]string{"valid-token": "testuser"},
		userRoles: map[string]string{"testuser": "admin"},
	}
	
	router := gin.New()
	router.GET("/auth/me", server.getAuthMe)
	
	// Test with valid token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/me", nil)
	req.Header.Set("X-App-Auth", "valid-token")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "admin")
	
	// Test without token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/auth/me", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "null")
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		sessions:  map[string]string{"admin-token": "admin", "editor-token": "editor"},
		userRoles: map[string]string{"admin": "admin", "editor": "editor"},
	}
	
	router := gin.New()
	router.Use(server.authMiddleware())
	router.GET("/protected", server.requireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	
	// Test admin access
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("X-App-Auth", "admin-token")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Test editor denied
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("X-App-Auth", "editor-token")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}

