package mocks

import (
	"encoding/json"
	"net/http"
)

// DefaultWorkspaceListHandler returns a default handler for listing workspaces
func DefaultWorkspaceListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workspaces := []map[string]interface{}{
			{
				"id":     "workspace-1",
				"name":   "test-workspace",
				"status": "active",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"workspaces": workspaces,
		})
	}
}

// DefaultErrorHandler returns a handler that returns an error
func DefaultErrorHandler(statusCode int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": message,
		})
	}
}

// DefaultSuccessHandler returns a handler that returns success
func DefaultSuccessHandler(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": message,
		})
	}
}

