package aiops

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// ScalingAdvisor provides auto-scaling recommendations
type ScalingAdvisor struct {
	analytics *AnalyticsEngine
	policies  map[string]*ScalingPolicy
	mu        sync.RWMutex
}

// NewScalingAdvisor creates a new scaling advisor
func NewScalingAdvisor(analytics *AnalyticsEngine) *ScalingAdvisor {
	return &ScalingAdvisor{
		analytics: analytics,
		policies:  make(map[string]*ScalingPolicy),
	}
}

// ScalingPolicy defines scaling rules for a container or deployment
type ScalingPolicy struct {
	ID               string
	TargetID         string // container ID or deployment ID
	Type             string // "container" or "deployment"
	MinReplicas      int
	MaxReplicas      int
	TargetCPU        float64 // Target CPU percentage
	TargetMemory     float64 // Target memory percentage
	ScaleUpThreshold float64 // Scale up when usage exceeds this
	ScaleDownThreshold float64 // Scale down when usage is below this
	CooldownPeriod   time.Duration
	LastScaling      time.Time
}

// GetScalingRecommendation provides scaling recommendations
func (sa *ScalingAdvisor) GetScalingRecommendation(ctx context.Context, targetID string, currentReplicas int) (*ScalingRecommendation, error) {
	sa.mu.RLock()
	policy, hasPolicy := sa.policies[targetID]
	sa.mu.RUnlock()

	if !hasPolicy {
		// Return default recommendation based on current usage
		prediction, err := sa.analytics.PredictResourceUsage(targetID, 30*time.Minute)
		if err != nil {
			return nil, err
		}

		rec, err := sa.generateRecommendation(targetID, currentReplicas, prediction, nil)
		return rec, err
	}

	// Check cooldown period
	if time.Since(policy.LastScaling) < policy.CooldownPeriod {
		rec := &ScalingRecommendation{
			TargetID:    targetID,
			Action:      "none",
			CurrentReplicas: currentReplicas,
			RecommendedReplicas: currentReplicas,
			Reason:      "cooldown_period",
			Message:     fmt.Sprintf("Scaling on cooldown. Next scaling available in %v", policy.CooldownPeriod-time.Since(policy.LastScaling)),
			Confidence:  1.0,
		}
		return rec, nil
	}

	// Get prediction
	prediction, err := sa.analytics.PredictResourceUsage(targetID, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	rec, err := sa.generateRecommendation(targetID, currentReplicas, prediction, policy)
	return rec, err
}

// generateRecommendation generates a scaling recommendation
func (sa *ScalingAdvisor) generateRecommendation(targetID string, currentReplicas int, prediction *ResourcePrediction, policy *ScalingPolicy) (*ScalingRecommendation, error) {
	var action string
	var newReplicas int
	var reason string
	var message string

	// Use policy thresholds if available, otherwise use defaults
	cpuThresholdUp := 80.0
	cpuThresholdDown := 20.0
	memThresholdUp := 80.0
	memThresholdDown := 20.0
	minReplicas := 1
	maxReplicas := 10

	if policy != nil {
		cpuThresholdUp = policy.ScaleUpThreshold
		cpuThresholdDown = policy.ScaleDownThreshold
		memThresholdUp = policy.ScaleUpThreshold
		memThresholdDown = policy.ScaleDownThreshold
		minReplicas = policy.MinReplicas
		maxReplicas = policy.MaxReplicas
	}

	// Check if scale up is needed
	if prediction.PredictedCPU > cpuThresholdUp || prediction.PredictedMemory > memThresholdUp {
		if currentReplicas < maxReplicas {
			action = "scale_up"
			// Calculate new replica count based on predicted usage
			maxUsage := prediction.PredictedCPU
			if prediction.PredictedMemory > maxUsage {
				maxUsage = prediction.PredictedMemory
			}
			// Scale to bring usage below threshold (with 20% headroom)
			targetUsage := cpuThresholdUp * 0.8
			newReplicas = int(math.Ceil(float64(currentReplicas) * maxUsage / targetUsage))
			if newReplicas > maxReplicas {
				newReplicas = maxReplicas
			}
		reason = "high_resource_usage"
		message = fmt.Sprintf("High resource usage predicted (CPU: %.2f%%, Memory: %.2f%%). Scaling from %d to %d replicas.", 
			prediction.PredictedCPU, prediction.PredictedMemory, currentReplicas, newReplicas)
		newReplicas = newReplicas
		} else {
		action = "none"
		reason = "max_replicas_reached"
		message = fmt.Sprintf("High resource usage but already at max replicas (%d)", maxReplicas)
		newReplicas = currentReplicas
	}
	} else if prediction.PredictedCPU < cpuThresholdDown && prediction.PredictedMemory < memThresholdDown {
		// Check if scale down is possible
		if currentReplicas > minReplicas {
			action = "scale_down"
			// Calculate new replica count
			avgUsage := (prediction.PredictedCPU + prediction.PredictedMemory) / 2
			targetUsage := cpuThresholdDown * 1.2 // 20% above threshold for safety
			newReplicas = int(math.Floor(float64(currentReplicas) * avgUsage / targetUsage))
			if newReplicas < minReplicas {
				newReplicas = minReplicas
			}
			reason = "low_resource_usage"
			message = fmt.Sprintf("Low resource usage predicted (CPU: %.2f%%, Memory: %.2f%%). Scaling from %d to %d replicas.", 
				prediction.PredictedCPU, prediction.PredictedMemory, currentReplicas, newReplicas)
		} else {
			action = "none"
			reason = "min_replicas_reached"
			message = fmt.Sprintf("Low resource usage but already at min replicas (%d)", minReplicas)
		}
	} else {
		action = "none"
		reason = "optimal"
		message = "Resource usage is within optimal range. No scaling needed."
		newReplicas = currentReplicas
	}

	rec := &ScalingRecommendation{
		TargetID:        targetID,
		Action:          action,
		CurrentReplicas: currentReplicas,
		RecommendedReplicas: newReplicas,
		Reason:          reason,
		Message:        message,
		Confidence:      prediction.Confidence,
		PredictedCPU:    prediction.PredictedCPU,
		PredictedMemory: prediction.PredictedMemory,
	}
	return rec, nil
}

// SetScalingPolicy sets a scaling policy for a target
func (sa *ScalingAdvisor) SetScalingPolicy(policy *ScalingPolicy) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.policies[policy.TargetID] = policy
}

// GetScalingPolicy retrieves a scaling policy
func (sa *ScalingAdvisor) GetScalingPolicy(targetID string) (*ScalingPolicy, bool) {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	policy, exists := sa.policies[targetID]
	return policy, exists
}

// ScalingRecommendation represents a scaling recommendation
type ScalingRecommendation struct {
	TargetID           string
	Action             string // scale_up, scale_down, none
	CurrentReplicas    int
	RecommendedReplicas int
	Reason             string
	Message            string
	Confidence         float64
	PredictedCPU       float64
	PredictedMemory    float64
}

// RecordScaling records a scaling action for cooldown tracking
func (sa *ScalingAdvisor) RecordScaling(targetID string) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	
	if policy, exists := sa.policies[targetID]; exists {
		policy.LastScaling = time.Now()
	}
}

