package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthMonitor monitors container health and manages auto-restart
type HealthMonitor struct {
	nodeManager      *NodeManager
	deploymentManager *DeploymentManager
	monitoring       map[string]*monitorContext // deploymentID -> context
	mu               sync.RWMutex
	stopCh           chan struct{}
	wg               sync.WaitGroup
}

// monitorContext tracks monitoring state for a deployment
type monitorContext struct {
	deploymentID string
	cancel       context.CancelFunc
	interval     time.Duration
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor(nodeManager *NodeManager, deploymentManager *DeploymentManager) *HealthMonitor {
	return &HealthMonitor{
		nodeManager:       nodeManager,
		deploymentManager: deploymentManager,
		monitoring:         make(map[string]*monitorContext),
		stopCh:             make(chan struct{}),
	}
}

// StartMonitoring starts health monitoring for a deployment
func (hm *HealthMonitor) StartMonitoring(deploymentID string, interval time.Duration) error {
	if interval <= 0 {
		interval = 30 * time.Second // Default
	}

	hm.mu.Lock()
	defer hm.mu.Unlock()

	// Check if already monitoring
	if _, exists := hm.monitoring[deploymentID]; exists {
		return fmt.Errorf("already monitoring deployment: %s", deploymentID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	hm.monitoring[deploymentID] = &monitorContext{
		deploymentID: deploymentID,
		cancel:       cancel,
		interval:     interval,
	}

	hm.wg.Add(1)
	go hm.monitorDeployment(ctx, deploymentID, interval)

	return nil
}

// StopMonitoring stops health monitoring for a deployment
func (hm *HealthMonitor) StopMonitoring(deploymentID string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	monitor, exists := hm.monitoring[deploymentID]
	if !exists {
		return
	}

	monitor.cancel()
	delete(hm.monitoring, deploymentID)
}

// Stop stops all monitoring
func (hm *HealthMonitor) Stop() {
	close(hm.stopCh)

	hm.mu.Lock()
	for _, monitor := range hm.monitoring {
		monitor.cancel()
	}
	hm.monitoring = make(map[string]*monitorContext)
	hm.mu.Unlock()

	hm.wg.Wait()
}

// monitorDeployment continuously monitors a deployment's health
func (hm *HealthMonitor) monitorDeployment(ctx context.Context, deploymentID string, interval time.Duration) {
	defer hm.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-hm.stopCh:
			return
		case <-ticker.C:
			hm.checkDeploymentHealth(deploymentID)
		}
	}
}

// checkDeploymentHealth checks health of all instances in a deployment
func (hm *HealthMonitor) checkDeploymentHealth(deploymentID string) {
	deploy, exists := hm.deploymentManager.GetDeployment(deploymentID)
	if !exists {
		return
	}

	for _, instance := range deploy.Instances {
		// Mock health check - in production, this would query the container health
		// For now, simulate health checking
		health := hm.checkInstanceHealth(instance)

		if health != instance.Health {
			hm.deploymentManager.UpdateInstanceHealth(deploymentID, instance.ID, health)

			// If unhealthy and auto-restart enabled, restart the instance
			if health == "unhealthy" && deploy.AutoRestart {
				hm.handleUnhealthyInstance(deploy, instance)
			}
		}
	}
}

// checkInstanceHealth checks the health of a single instance
func (hm *HealthMonitor) checkInstanceHealth(instance DeploymentInstance) string {
	// Mock implementation - in production, this would:
	// 1. Query the container health via containerd API
	// 2. Check HTTP/TCP/CMD health probes
	// 3. Return healthy/unhealthy/unknown

	// Simple heuristic: if status is running, assume healthy
	if instance.Status == "running" {
		return "healthy"
	}
	if instance.Status == "failed" || instance.Status == "stopped" {
		return "unhealthy"
	}
	return "unknown"
}

// handleUnhealthyInstance handles an unhealthy instance (restart if needed)
func (hm *HealthMonitor) handleUnhealthyInstance(deploy *Deployment, instance DeploymentInstance) {
	// Check restart policy
	if deploy.RestartPolicy == "never" {
		return
	}

	// Check max restarts
	if deploy.MaxRestarts > 0 && instance.Restarts >= deploy.MaxRestarts {
		// Mark instance as failed permanently
		hm.deploymentManager.UpdateInstanceStatus(deploy.ID, instance.ID, "failed")
		return
	}

	// Check restart policy
	if deploy.RestartPolicy == "on-failure" && instance.Status != "failed" {
		return
	}

	// Record restart
	hm.deploymentManager.RecordInstanceRestart(deploy.ID, instance.ID)

	// Restart the instance (this would trigger actual container restart in production)
	// For now, just update status
	hm.deploymentManager.UpdateInstanceStatus(deploy.ID, instance.ID, "pending")
	
	// In production, this would:
	// 1. Stop the unhealthy container
	// 2. Create a new container instance
	// 3. Deploy to a node (possibly different node if current is down)
}

