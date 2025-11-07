package cli_test

import (
	"net/http"
	"testing"

	"github.com/nebulabox/nebulabox/internal/cli/tests/mocks"
)

// TestSystemStats tests the `nebulabox stats` command
func TestSystemStats(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	mockAPI.RegisterHandler("GET", "/api/system/stats", mocks.DefaultSystemStatsHandler())

	// TODO: Test command execution
	t.Log("System stats test setup complete")
}

// TestSystemMode tests the `nebulabox mode` command
func TestSystemMode(t *testing.T) {
	mockAPI := mocks.NewMockAPIServer()
	defer mockAPI.Close()

	// Test getting mode
	mockAPI.RegisterHandler("GET", "/api/mode", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"mode":"mock"}`))
	})

	// Test setting mode
	mockAPI.RegisterHandler("POST", "/api/mode", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"mode":"test"}`))
	})

	// TODO: Test command execution
	t.Log("System mode test setup complete")
}

