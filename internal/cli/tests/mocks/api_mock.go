package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// MockAPIServer creates a mock HTTP server for testing
type MockAPIServer struct {
	server   *httptest.Server
	handlers map[string]http.HandlerFunc
}

// NewMockAPIServer creates a new mock API server
func NewMockAPIServer() *MockAPIServer {
	mock := &MockAPIServer{
		handlers: make(map[string]http.HandlerFunc),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		key := fmt.Sprintf("%s %s", method, path)
		if handler, ok := mock.handlers[key]; ok {
			handler(w, r)
			return
		}

		// Default handler
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Not found",
		})
	})

	mock.server = httptest.NewServer(mux)
	return mock
}

// URL returns the base URL of the mock server
func (m *MockAPIServer) URL() string {
	return m.server.URL
}

// Close shuts down the mock server
func (m *MockAPIServer) Close() {
	m.server.Close()
}

// RegisterHandler registers a handler for a specific method and path
func (m *MockAPIServer) RegisterHandler(method, path string, handler http.HandlerFunc) {
	key := fmt.Sprintf("%s %s", method, path)
	m.handlers[key] = handler
}

// DefaultContainerListHandler returns a default handler for listing containers
func DefaultContainerListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containers := []map[string]interface{}{
			{
				"id":     "container-1",
				"name":   "test-container",
				"image":  "nginx:latest",
				"status": "running",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(containers)
	}
}

// DefaultImageListHandler returns a default handler for listing images
func DefaultImageListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		images := []map[string]interface{}{
			{
				"id":   "image-1",
				"name": "nginx",
				"tag":  "latest",
				"size": "100MB",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(images)
	}
}

// DefaultSystemStatsHandler returns a default handler for system stats
func DefaultSystemStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := map[string]interface{}{
			"cpuUsage":         45.2,
			"memoryUsage":      62.8,
			"diskUsage":        38.5,
			"containersRunning": 3,
			"containersTotal":   5,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

