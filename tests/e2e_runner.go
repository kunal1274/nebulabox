package tests

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

var (
	apiBaseURL   = flag.String("api-url", "http://localhost:8081", "API server base URL")
	startServer  = flag.Bool("start-server", false, "Start API server before tests")
	serverBinary = flag.String("server-binary", "./nebulabox-api", "Path to API server binary")
	waitTimeout  = flag.Duration("wait-timeout", 30*time.Second, "Time to wait for server to start")
)

func TestMain(m *testing.M) {
	flag.Parse()
	
	var apiProcess *exec.Cmd
	var err error
	
	if *startServer {
		log.Println("🚀 Starting API server for E2E tests...")
		
		// Check if binary exists
		if _, err := os.Stat(*serverBinary); os.IsNotExist(err) {
			log.Fatalf("❌ Server binary not found: %s", *serverBinary)
		}
		
		// Start API server
		apiProcess = exec.Command(*serverBinary)
		apiProcess.Stdout = os.Stdout
		apiProcess.Stderr = os.Stderr
		
		if err = apiProcess.Start(); err != nil {
			log.Fatalf("❌ Failed to start API server: %v", err)
		}
		
		log.Printf("⏳ Waiting for API server at %s...", *apiBaseURL)
		
		// Wait for server to be ready
		if !waitForServer(*apiBaseURL, *waitTimeout) {
			log.Fatalf("❌ API server did not become available within %v", *waitTimeout)
		}
		
		log.Println("✅ API server is ready!")
		
		// Setup cleanup
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("\n🛑 Shutting down API server...")
			if apiProcess != nil && apiProcess.Process != nil {
				apiProcess.Process.Kill()
			}
			os.Exit(1)
		}()
		
		// Run tests
		code := m.Run()
		
		// Cleanup
		log.Println("🛑 Stopping API server...")
		if apiProcess != nil && apiProcess.Process != nil {
			if err := apiProcess.Process.Kill(); err != nil {
				log.Printf("⚠️  Error killing server: %v", err)
			}
			apiProcess.Wait()
		}
		
		os.Exit(code)
	} else {
		// Just run tests, assume server is already running
		log.Printf("📡 Using existing API server at %s", *apiBaseURL)
		code := m.Run()
		os.Exit(code)
	}
}

// RunE2ETests runs all E2E tests and returns results
// Note: This function is for CLI use; actual tests run via go test
func RunE2ETests(ctx context.Context, apiURL string) (passed, failed int, err error) {
	// Tests are executed via go test, not programmatically
	// This is a placeholder for potential CLI runner
	return 0, 0, nil
}

// CheckServerHealth checks if the API server is healthy
func CheckServerHealth(baseURL string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL + "/api/auth/me")
	if err != nil {
		return fmt.Errorf("server not reachable: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 500 {
		return fmt.Errorf("server returned error: %d", resp.StatusCode)
	}
	
	return nil
}

