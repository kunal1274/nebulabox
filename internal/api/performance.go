package api

import (
	"sync"
	"sync/atomic"
	"time"
)

// PerformanceMetrics tracks API performance metrics
type PerformanceMetrics struct {
	RequestCount     int64
	TotalLatency    int64 // nanoseconds
	MinLatency      int64
	MaxLatency      int64
	ErrorCount      int64
	LastRequestTime time.Time
	mu              sync.RWMutex
}

// RecordRequest records a request's latency
func (pm *PerformanceMetrics) RecordRequest(latency time.Duration, isError bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	atomic.AddInt64(&pm.RequestCount, 1)
	atomic.AddInt64(&pm.TotalLatency, int64(latency))
	
	if isError {
		atomic.AddInt64(&pm.ErrorCount, 1)
	}
	
	latencyNs := int64(latency)
	
	// Update min
	for {
		currentMin := atomic.LoadInt64(&pm.MinLatency)
		if currentMin == 0 || latencyNs < currentMin {
			if atomic.CompareAndSwapInt64(&pm.MinLatency, currentMin, latencyNs) {
				break
			}
		} else {
			break
		}
	}
	
	// Update max
	for {
		currentMax := atomic.LoadInt64(&pm.MaxLatency)
		if latencyNs > currentMax {
			if atomic.CompareAndSwapInt64(&pm.MaxLatency, currentMax, latencyNs) {
				break
			}
		} else {
			break
		}
	}
	
	pm.LastRequestTime = time.Now()
}

// GetStats returns current performance statistics
func (pm *PerformanceMetrics) GetStats() map[string]interface{} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	count := atomic.LoadInt64(&pm.RequestCount)
	totalLatency := atomic.LoadInt64(&pm.TotalLatency)
	minLatency := atomic.LoadInt64(&pm.MinLatency)
	maxLatency := atomic.LoadInt64(&pm.MaxLatency)
	errorCount := atomic.LoadInt64(&pm.ErrorCount)
	
	var avgLatency time.Duration
	if count > 0 {
		avgLatency = time.Duration(totalLatency / count)
	}
	
	var errorRate float64
	if count > 0 {
		errorRate = float64(errorCount) / float64(count) * 100
	}
	
	return map[string]interface{}{
		"requestCount":   count,
		"avgLatencyMs":   avgLatency.Milliseconds(),
		"minLatencyMs":   time.Duration(minLatency).Milliseconds(),
		"maxLatencyMs":   time.Duration(maxLatency).Milliseconds(),
		"errorCount":     errorCount,
		"errorRate":      errorRate,
		"lastRequestTime": pm.LastRequestTime,
	}
}

// Reset clears all metrics
func (pm *PerformanceMetrics) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	atomic.StoreInt64(&pm.RequestCount, 0)
	atomic.StoreInt64(&pm.TotalLatency, 0)
	atomic.StoreInt64(&pm.MinLatency, 0)
	atomic.StoreInt64(&pm.MaxLatency, 0)
	atomic.StoreInt64(&pm.ErrorCount, 0)
	pm.LastRequestTime = time.Time{}
}

// EndpointMetrics tracks performance per endpoint
type EndpointMetrics struct {
	metrics map[string]*PerformanceMetrics
	mu      sync.RWMutex
}

// NewEndpointMetrics creates a new endpoint metrics tracker
func NewEndpointMetrics() *EndpointMetrics {
	return &EndpointMetrics{
		metrics: make(map[string]*PerformanceMetrics),
	}
}

// GetOrCreate gets or creates metrics for an endpoint
func (em *EndpointMetrics) GetOrCreate(endpoint string) *PerformanceMetrics {
	em.mu.RLock()
	pm, exists := em.metrics[endpoint]
	em.mu.RUnlock()
	
	if exists {
		return pm
	}
	
	em.mu.Lock()
	defer em.mu.Unlock()
	
	// Double-check after acquiring write lock
	pm, exists = em.metrics[endpoint]
	if exists {
		return pm
	}
	
	pm = &PerformanceMetrics{}
	em.metrics[endpoint] = pm
	return pm
}

// GetStats returns stats for all endpoints
func (em *EndpointMetrics) GetStats() map[string]map[string]interface{} {
	em.mu.RLock()
	defer em.mu.RUnlock()
	
	stats := make(map[string]map[string]interface{})
	for endpoint, metrics := range em.metrics {
		stats[endpoint] = metrics.GetStats()
	}
	return stats
}

// Reset clears all endpoint metrics
func (em *EndpointMetrics) Reset() {
	em.mu.Lock()
	defer em.mu.Unlock()
	
	for _, metrics := range em.metrics {
		metrics.Reset()
	}
}

