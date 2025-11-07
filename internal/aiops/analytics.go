package aiops

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// AnalyticsEngine provides predictive analytics for container operations
type AnalyticsEngine struct {
	metricsHistory map[string][]MetricPoint
	mu             sync.RWMutex
	predictionWindow time.Duration
}

// NewAnalyticsEngine creates a new analytics engine
func NewAnalyticsEngine() *AnalyticsEngine {
	return &AnalyticsEngine{
		metricsHistory:   make(map[string][]MetricPoint),
		predictionWindow: 1 * time.Hour, // Predict 1 hour ahead
	}
}

// MetricPoint represents a single metric data point
type MetricPoint struct {
	Timestamp time.Time
	Value     float64
	CPU       float64
	Memory    float64
	NetworkRx float64
	NetworkTx float64
}

// RecordMetric records a metric point for analysis
func (ae *AnalyticsEngine) RecordMetric(containerID string, point MetricPoint) {
	ae.mu.Lock()
	defer ae.mu.Unlock()

	if ae.metricsHistory[containerID] == nil {
		ae.metricsHistory[containerID] = make([]MetricPoint, 0, 100)
	}

	history := ae.metricsHistory[containerID]
	
	// Keep only recent history (last 24 hours of data points)
	cutoff := time.Now().Add(-24 * time.Hour)
	filtered := history[:0]
	for _, p := range history {
		if p.Timestamp.After(cutoff) {
			filtered = append(filtered, p)
		}
	}
	filtered = append(filtered, point)
	
	// Limit to 1000 points max
	if len(filtered) > 1000 {
		filtered = filtered[len(filtered)-1000:]
	}
	
	ae.metricsHistory[containerID] = filtered
}

// PredictResourceUsage predicts future resource usage
func (ae *AnalyticsEngine) PredictResourceUsage(containerID string, duration time.Duration) (*ResourcePrediction, error) {
	ae.mu.RLock()
	defer ae.mu.RUnlock()

	history, exists := ae.metricsHistory[containerID]
	if !exists || len(history) < 10 {
		return nil, fmt.Errorf("insufficient data for prediction")
	}

	// Simple linear regression for prediction
	cpuPrediction := ae.predictValue(history, func(p MetricPoint) float64 { return p.CPU }, duration)
	memoryPrediction := ae.predictValue(history, func(p MetricPoint) float64 { return p.Memory }, duration)

	// Calculate trends
	cpuTrend := ae.calculateTrend(history, func(p MetricPoint) float64 { return p.CPU })
	memoryTrend := ae.calculateTrend(history, func(p MetricPoint) float64 { return p.Memory })

	// Detect anomalies
	anomalies := ae.detectAnomalies(history)

	return &ResourcePrediction{
		ContainerID:    containerID,
		PredictedAt:    time.Now(),
		Duration:       duration,
		PredictedCPU:   cpuPrediction,
		PredictedMemory: memoryPrediction,
		CPUTrend:       cpuTrend,
		MemoryTrend:    memoryTrend,
		Confidence:     ae.calculateConfidence(history),
		Anomalies:      anomalies,
		Recommendations: ae.generateRecommendations(cpuPrediction, memoryPrediction, cpuTrend, memoryTrend, history),
	}, nil
}

// predictValue uses linear regression to predict future values
func (ae *AnalyticsEngine) predictValue(history []MetricPoint, extractor func(MetricPoint) float64, duration time.Duration) float64 {
	if len(history) < 2 {
		return 0
	}

	// Simple moving average with trend
	recent := history
	if len(history) > 100 {
		recent = history[len(history)-100:]
	}

	// Calculate average and trend
	var sum float64
	var timeSum float64
	var valueSum float64
	var timeValueSum float64
	var timeSquaredSum float64

	for i, point := range recent {
		value := extractor(point)
		time := float64(i)
		
		sum += value
		timeSum += time
		valueSum += value
		timeValueSum += time * value
		timeSquaredSum += time * time
	}

	n := float64(len(recent))
	
	// Linear regression: y = a + b*x
	slope := (n*timeValueSum - timeSum*valueSum) / (n*timeSquaredSum - timeSum*timeSum)
	intercept := (valueSum - slope*timeSum) / n
	
	// Predict future value
	futureTime := float64(len(recent)) + duration.Seconds()/60 // minutes
	predicted := intercept + slope*futureTime
	
	// Ensure non-negative
	if predicted < 0 {
		predicted = 0
	}

	return predicted
}

// calculateTrend calculates the trend (increasing, decreasing, stable)
func (ae *AnalyticsEngine) calculateTrend(history []MetricPoint, extractor func(MetricPoint) float64) string {
	if len(history) < 2 {
		return "unknown"
	}

	recent := history
	if len(history) > 50 {
		recent = history[len(history)-50:]
	}

	first := extractor(recent[0])
	last := extractor(recent[len(recent)-1])
	
	change := ((last - first) / first) * 100
	
	if math.Abs(change) < 5 {
		return "stable"
	} else if change > 0 {
		return "increasing"
	} else {
		return "decreasing"
	}
}

// detectAnomalies detects anomalous patterns in metrics
func (ae *AnalyticsEngine) detectAnomalies(history []MetricPoint) []Anomaly {
	if len(history) < 10 {
		return nil
	}

	var anomalies []Anomaly

	// Calculate mean and standard deviation
	var cpuSum, memorySum float64
	for _, p := range history {
		cpuSum += p.CPU
		memorySum += p.Memory
	}
	cpuMean := cpuSum / float64(len(history))
	memoryMean := memorySum / float64(len(history))

	var cpuVar, memoryVar float64
	for _, p := range history {
		cpuVar += math.Pow(p.CPU-cpuMean, 2)
		memoryVar += math.Pow(p.Memory-memoryMean, 2)
	}
	cpuStdDev := math.Sqrt(cpuVar / float64(len(history)))
	memoryStdDev := math.Sqrt(memoryVar / float64(len(history)))

	// Detect outliers (values beyond 2 standard deviations)
	for _, p := range history[len(history)-10:] {
		if math.Abs(p.CPU-cpuMean) > 2*cpuStdDev {
			anomalies = append(anomalies, Anomaly{
				Type:      "cpu_spike",
				Severity:  "medium",
				Timestamp: p.Timestamp,
				Value:     p.CPU,
				Message:   fmt.Sprintf("CPU usage spike detected: %.2f%%", p.CPU),
			})
		}
		if math.Abs(p.Memory-memoryMean) > 2*memoryStdDev {
			anomalies = append(anomalies, Anomaly{
				Type:      "memory_spike",
				Severity:  "medium",
				Timestamp: p.Timestamp,
				Value:     p.Memory,
				Message:   fmt.Sprintf("Memory usage spike detected: %.2f%%", p.Memory),
			})
		}
	}

	return anomalies
}

// calculateConfidence calculates prediction confidence based on data quality
func (ae *AnalyticsEngine) calculateConfidence(history []MetricPoint) float64 {
	if len(history) < 10 {
		return 0.3
	}
	if len(history) < 50 {
		return 0.6
	}
	if len(history) < 100 {
		return 0.8
	}
	return 0.95
}

// generateRecommendations generates recommendations based on predictions
func (ae *AnalyticsEngine) generateRecommendations(cpuPred, memPred float64, cpuTrend, memTrend string, history []MetricPoint) []Recommendation {
	var recommendations []Recommendation

	if cpuPred > 80 && cpuTrend == "increasing" {
		recommendations = append(recommendations, Recommendation{
			Type:        "scale_up",
			Priority:    "high",
			Category:    "cpu",
			Title:       "Scale up due to high CPU usage",
			Description: fmt.Sprintf("Predicted CPU usage: %.2f%%. Consider scaling up resources.", cpuPred),
			Action:      "increase_cpu_limits",
		})
	}

	if memPred > 80 && memTrend == "increasing" {
		recommendations = append(recommendations, Recommendation{
			Type:        "scale_up",
			Priority:    "high",
			Category:    "memory",
			Title:       "Scale up due to high memory usage",
			Description: fmt.Sprintf("Predicted memory usage: %.2f%%. Consider scaling up resources.", memPred),
			Action:      "increase_memory_limits",
		})
	}

	if cpuPred < 20 && cpuTrend == "decreasing" {
		recommendations = append(recommendations, Recommendation{
			Type:        "scale_down",
			Priority:    "low",
			Category:    "cpu",
			Title:       "Consider scaling down",
			Description: fmt.Sprintf("Low CPU usage predicted: %.2f%%. Resources may be over-provisioned.", cpuPred),
			Action:      "decrease_cpu_limits",
		})
	}

	if memPred < 20 && memTrend == "decreasing" {
		recommendations = append(recommendations, Recommendation{
			Type:        "scale_down",
			Priority:    "low",
			Category:    "memory",
			Title:       "Consider scaling down",
			Description: fmt.Sprintf("Low memory usage predicted: %.2f%%. Resources may be over-provisioned.", memPred),
			Action:      "decrease_memory_limits",
		})
	}

	return recommendations
}

// ResourcePrediction represents a resource usage prediction
type ResourcePrediction struct {
	ContainerID     string
	PredictedAt     time.Time
	Duration        time.Duration
	PredictedCPU    float64
	PredictedMemory float64
	CPUTrend        string
	MemoryTrend     string
	Confidence      float64
	Anomalies       []Anomaly
	Recommendations []Recommendation
}

// Anomaly represents a detected anomaly
type Anomaly struct {
	Type      string
	Severity  string // low, medium, high
	Timestamp time.Time
	Value     float64
	Message   string
}

// Recommendation represents an optimization recommendation
type Recommendation struct {
	Type        string // scale_up, scale_down, optimize, alert
	Priority    string // low, medium, high
	Category    string // cpu, memory, network, disk
	Title       string
	Description string
	Action      string
}

