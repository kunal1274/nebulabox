package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
)

// LoadTestResult holds results from a load test
type LoadTestResult struct {
	TotalRequests    int64
	SuccessfulReqs   int64
	FailedReqs       int64
	TotalDuration    time.Duration
	AvgResponseTime  time.Duration
	MinResponseTime  time.Duration
	MaxResponseTime  time.Duration
	RequestsPerSec   float64
	ErrorRate        float64
}

// runLoadTest runs concurrent requests against an endpoint
func runLoadTest(handler gin.HandlerFunc, method, path string, body interface{}, concurrency, requests int) LoadTestResult {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(method, path, handler)
	
	var totalReqs int64
	var successReqs int64
	var failedReqs int64
	var totalTime int64
	var minTime int64 = 999999999
	var maxTime int64
	
	var wg sync.WaitGroup
	startTime := time.Now()
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for j := 0; j < requests/concurrency; j++ {
				reqStart := time.Now()
				
				var req *http.Request
				var err error
				if body != nil {
					jsonBody, _ := json.Marshal(body)
					req, err = http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
					req.Header.Set("Content-Type", "application/json")
				} else {
					req, err = http.NewRequest(method, path, nil)
				}
				if err != nil {
					atomic.AddInt64(&failedReqs, 1)
					atomic.AddInt64(&totalReqs, 1)
					continue
				}
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				
				duration := time.Since(reqStart).Nanoseconds()
				atomic.AddInt64(&totalReqs, 1)
				atomic.AddInt64(&totalTime, duration)
				
				if w.Code >= 200 && w.Code < 300 {
					atomic.AddInt64(&successReqs, 1)
				} else {
					atomic.AddInt64(&failedReqs, 1)
				}
				
				// Update min/max
				for {
					currentMin := atomic.LoadInt64(&minTime)
					if duration >= currentMin {
						break
					}
					if atomic.CompareAndSwapInt64(&minTime, currentMin, duration) {
						break
					}
				}
				for {
					currentMax := atomic.LoadInt64(&maxTime)
					if duration <= currentMax {
						break
					}
					if atomic.CompareAndSwapInt64(&maxTime, currentMax, duration) {
						break
					}
				}
			}
		}()
	}
	
	wg.Wait()
	totalDuration := time.Since(startTime)
	
	totalReqsVal := atomic.LoadInt64(&totalReqs)
	successReqsVal := atomic.LoadInt64(&successReqs)
	failedReqsVal := atomic.LoadInt64(&failedReqs)
	totalTimeVal := atomic.LoadInt64(&totalTime)
	minTimeVal := atomic.LoadInt64(&minTime)
	maxTimeVal := atomic.LoadInt64(&maxTime)
	
	avgTime := time.Duration(totalTimeVal / totalReqsVal)
	if totalReqsVal == 0 {
		avgTime = 0
	}
	
	reqsPerSec := float64(totalReqsVal) / totalDuration.Seconds()
	errorRate := float64(failedReqsVal) / float64(totalReqsVal) * 100
	if totalReqsVal == 0 {
		errorRate = 0
	}
	
	return LoadTestResult{
		TotalRequests:   totalReqsVal,
		SuccessfulReqs: successReqsVal,
		FailedReqs:     failedReqsVal,
		TotalDuration:  totalDuration,
		AvgResponseTime: avgTime,
		MinResponseTime: time.Duration(minTimeVal),
		MaxResponseTime: time.Duration(maxTimeVal),
		RequestsPerSec:  reqsPerSec,
		ErrorRate:       errorRate,
	}
}

// TestLoad_ListContainers tests concurrent container listing
func TestLoad_ListContainers(t *testing.T) {
	client, _ := containerd.NewClient()
	server := &Server{
		containerd: client,
		networks:   make(map[string]*Network),
	}
	
	result := runLoadTest(server.listContainers, "GET", "/containers", nil, 10, 100)
	
	fmt.Printf("\n📊 Load Test Results - List Containers\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("Total Requests:    %d\n", result.TotalRequests)
	fmt.Printf("Successful:        %d\n", result.SuccessfulReqs)
	fmt.Printf("Failed:            %d\n", result.FailedReqs)
	fmt.Printf("Total Duration:    %v\n", result.TotalDuration)
	fmt.Printf("Avg Response Time: %v\n", result.AvgResponseTime)
	fmt.Printf("Min Response Time: %v\n", result.MinResponseTime)
	fmt.Printf("Max Response Time: %v\n", result.MaxResponseTime)
	fmt.Printf("Requests/sec:      %.2f\n", result.RequestsPerSec)
	fmt.Printf("Error Rate:        %.2f%%\n", result.ErrorRate)
	fmt.Printf("=====================================\n\n")
	
	if result.ErrorRate > 5.0 {
		t.Errorf("Error rate too high: %.2f%%", result.ErrorRate)
	}
}

// TestLoad_ListNetworks tests concurrent network listing
func TestLoad_ListNetworks(t *testing.T) {
	server := &Server{
		networks: make(map[string]*Network),
	}
	
	// Pre-populate networks
	for i := 0; i < 50; i++ {
		server.networks[fmt.Sprintf("net-%d", i)] = &Network{
			ID:     fmt.Sprintf("net-%d", i),
			Name:   fmt.Sprintf("network-%d", i),
			Driver: "bridge",
		}
	}
	
	result := runLoadTest(server.listNetworks, "GET", "/networks", nil, 20, 200)
	
	fmt.Printf("\n📊 Load Test Results - List Networks\n")
	fmt.Printf("====================================\n")
	fmt.Printf("Total Requests:    %d\n", result.TotalRequests)
	fmt.Printf("Successful:        %d\n", result.SuccessfulReqs)
	fmt.Printf("Failed:            %d\n", result.FailedReqs)
	fmt.Printf("Total Duration:    %v\n", result.TotalDuration)
	fmt.Printf("Avg Response Time: %v\n", result.AvgResponseTime)
	fmt.Printf("Min Response Time: %v\n", result.MinResponseTime)
	fmt.Printf("Max Response Time: %v\n", result.MaxResponseTime)
	fmt.Printf("Requests/sec:      %.2f\n", result.RequestsPerSec)
	fmt.Printf("Error Rate:        %.2f%%\n", result.ErrorRate)
	fmt.Printf("====================================\n\n")
	
	if result.ErrorRate > 5.0 {
		t.Errorf("Error rate too high: %.2f%%", result.ErrorRate)
	}
}

// TestLoad_AuthLogin tests concurrent login requests
func TestLoad_AuthLogin(t *testing.T) {
	server := &Server{
		users:     map[string]string{"testuser": "testpass"},
		userRoles: map[string]string{"testuser": "editor"},
		sessions:  make(map[string]string),
	}
	
	body := map[string]string{"username": "testuser", "password": "testpass"}
	
	result := runLoadTest(server.postAuthLogin, "POST", "/auth/login", body, 15, 150)
	
	fmt.Printf("\n📊 Load Test Results - Auth Login\n")
	fmt.Printf("=================================\n")
	fmt.Printf("Total Requests:    %d\n", result.TotalRequests)
	fmt.Printf("Successful:        %d\n", result.SuccessfulReqs)
	fmt.Printf("Failed:            %d\n", result.FailedReqs)
	fmt.Printf("Total Duration:    %v\n", result.TotalDuration)
	fmt.Printf("Avg Response Time: %v\n", result.AvgResponseTime)
	fmt.Printf("Min Response Time: %v\n", result.MinResponseTime)
	fmt.Printf("Max Response Time: %v\n", result.MaxResponseTime)
	fmt.Printf("Requests/sec:      %.2f\n", result.RequestsPerSec)
	fmt.Printf("Error Rate:        %.2f%%\n", result.ErrorRate)
	fmt.Printf("=================================\n\n")
	
	if result.ErrorRate > 5.0 {
		t.Errorf("Error rate too high: %.2f%%", result.ErrorRate)
	}
}

