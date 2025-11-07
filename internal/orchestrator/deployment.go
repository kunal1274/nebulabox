package orchestrator

import (
	"fmt"
	"sync"
	"time"
)

// Deployment represents a multi-node deployment
type Deployment struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	Image         string                `json:"image"`
	Tag           string                `json:"tag"`
	Replicas      int                   `json:"replicas"`
	Status        string                `json:"status"` // pending, deploying, running, failed, stopped
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	Instances     []DeploymentInstance  `json:"instances"`
	NodeSelector  map[string]string     `json:"nodeSelector,omitempty"` // Node labels
	Strategy      string                `json:"strategy"` // rolling, recreate, canary
	HealthCheck   *HealthCheckConfig    `json:"healthCheck,omitempty"`
	AutoRestart   bool                  `json:"autoRestart"`
	MaxRestarts   int                   `json:"maxRestarts,omitempty"`
	RestartPolicy string                `json:"restartPolicy"` // always, on-failure, never
	ServiceName   string                `json:"serviceName,omitempty"` // Auto-register instances as a service
	NetworkName   string                `json:"networkName,omitempty"` // Network to create/use for deployment
	Ports         []string              `json:"ports,omitempty"` // Port mappings for service registration
}

// DeploymentInstance represents a single instance of a deployment
type DeploymentInstance struct {
	ID          string    `json:"id"`
	ContainerID string    `json:"containerId"`
	NodeID      string    `json:"nodeId"`
	NodeName    string    `json:"nodeName"`
	Status      string    `json:"status"` // pending, running, failed, stopped
	Health      string    `json:"health"` // healthy, unhealthy, unknown
	CreatedAt   time.Time `json:"createdAt"`
	Restarts    int       `json:"restarts"`
	LastRestart time.Time `json:"lastRestart,omitempty"`
}

// HealthCheckConfig defines deployment-wide health check configuration
type HealthCheckConfig struct {
	Type            string `json:"type"` // http, tcp, cmd
	HTTPPath        string `json:"httpPath,omitempty"`
	HTTPPort        int    `json:"httpPort,omitempty"`
	TCPPort         int    `json:"tcpPort,omitempty"`
	Command         []string `json:"command,omitempty"`
	IntervalSeconds int    `json:"intervalSeconds"`
	TimeoutSeconds  int    `json:"timeoutSeconds"`
	Retries         int    `json:"retries"`
	StartPeriodSec  int    `json:"startPeriodSec"`
}

// DeploymentManager manages multi-node deployments
type DeploymentManager struct {
	deployments map[string]*Deployment
	nodeManager *NodeManager
	mu          sync.RWMutex
}

// NewDeploymentManager creates a new deployment manager
func NewDeploymentManager(nodeManager *NodeManager) *DeploymentManager {
	return &DeploymentManager{
		deployments: make(map[string]*Deployment),
		nodeManager: nodeManager,
	}
}

// CreateDeployment creates a new deployment
func (dm *DeploymentManager) CreateDeployment(deploy *Deployment) error {
	if deploy.Name == "" {
		return fmt.Errorf("deployment name is required")
	}
	if deploy.Image == "" {
		return fmt.Errorf("image is required")
	}
	if deploy.Replicas <= 0 {
		return fmt.Errorf("replicas must be greater than 0")
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Generate ID if not provided
	if deploy.ID == "" {
		deploy.ID = fmt.Sprintf("deploy-%d", time.Now().UnixNano())
	}

	// Set defaults
	if deploy.Tag == "" {
		deploy.Tag = "latest"
	}
	if deploy.Status == "" {
		deploy.Status = "pending"
	}
	if deploy.Strategy == "" {
		deploy.Strategy = "rolling"
	}
	if deploy.RestartPolicy == "" {
		deploy.RestartPolicy = "always"
	}
	if deploy.AutoRestart {
		if deploy.MaxRestarts == 0 {
			deploy.MaxRestarts = 5
		}
	}

	deploy.CreatedAt = time.Now()
	deploy.UpdatedAt = time.Now()
	deploy.Instances = make([]DeploymentInstance, 0)

	dm.deployments[deploy.ID] = deploy
	return nil
}

// GetDeployment retrieves a deployment by ID
func (dm *DeploymentManager) GetDeployment(id string) (*Deployment, bool) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	deploy, exists := dm.deployments[id]
	return deploy, exists
}

// ListDeployments returns all deployments
func (dm *DeploymentManager) ListDeployments() []*Deployment {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	deployments := make([]*Deployment, 0, len(dm.deployments))
	for _, deploy := range dm.deployments {
		deployCopy := *deploy
		deployments = append(deployments, &deployCopy)
	}
	return deployments
}

// UpdateDeploymentStatus updates deployment status
func (dm *DeploymentManager) UpdateDeploymentStatus(id, status string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deploy, exists := dm.deployments[id]
	if !exists {
		return fmt.Errorf("deployment not found: %s", id)
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"deploying": true,
		"running":   true,
		"failed":    true,
		"stopped":   true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	deploy.Status = status
	deploy.UpdatedAt = time.Now()
	return nil
}

// AddInstance adds an instance to a deployment
func (dm *DeploymentManager) AddInstance(deploymentID string, instance DeploymentInstance) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deploy, exists := dm.deployments[deploymentID]
	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	// Check if instance already exists
	for _, inst := range deploy.Instances {
		if inst.ID == instance.ID {
			return nil // Already exists
		}
	}

	deploy.Instances = append(deploy.Instances, instance)
	deploy.UpdatedAt = time.Now()

	// Update node container count
	if instance.NodeID != "" {
		dm.nodeManager.AddContainerToNode(instance.NodeID, instance.ContainerID)
	}

	return nil
}

// GetDeploymentByServiceName returns a deployment by its service name
func (dm *DeploymentManager) GetDeploymentByServiceName(serviceName string) (*Deployment, bool) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	for _, deploy := range dm.deployments {
		if deploy.ServiceName == serviceName {
			return deploy, true
		}
	}
	return nil, false
}

// UpdateInstanceStatus updates instance status
func (dm *DeploymentManager) UpdateInstanceStatus(deploymentID, instanceID, status string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deploy, exists := dm.deployments[deploymentID]
	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	for i := range deploy.Instances {
		if deploy.Instances[i].ID == instanceID {
			deploy.Instances[i].Status = status
			deploy.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("instance not found: %s", instanceID)
}

// UpdateInstanceHealth updates instance health status
func (dm *DeploymentManager) UpdateInstanceHealth(deploymentID, instanceID, health string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deploy, exists := dm.deployments[deploymentID]
	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	for i := range deploy.Instances {
		if deploy.Instances[i].ID == instanceID {
			deploy.Instances[i].Health = health
			deploy.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("instance not found: %s", instanceID)
}

// RecordInstanceRestart records an instance restart
func (dm *DeploymentManager) RecordInstanceRestart(deploymentID, instanceID string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deploy, exists := dm.deployments[deploymentID]
	if !exists {
		return fmt.Errorf("deployment not found: %s", deploymentID)
	}

	for i := range deploy.Instances {
		if deploy.Instances[i].ID == instanceID {
			deploy.Instances[i].Restarts++
			deploy.Instances[i].LastRestart = time.Now()
			deploy.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("instance not found: %s", instanceID)
}

// DeleteDeployment deletes a deployment
func (dm *DeploymentManager) DeleteDeployment(id string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	deploy, exists := dm.deployments[id]
	if !exists {
		return fmt.Errorf("deployment not found: %s", id)
	}

	if deploy.Status == "deploying" || deploy.Status == "running" {
		return fmt.Errorf("cannot delete running deployment, stop it first")
	}

	delete(dm.deployments, id)
	return nil
}

// GetDeploymentInstances returns all instances for a deployment
func (dm *DeploymentManager) GetDeploymentInstances(deploymentID string) ([]DeploymentInstance, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	deploy, exists := dm.deployments[deploymentID]
	if !exists {
		return nil, fmt.Errorf("deployment not found: %s", deploymentID)
	}

	instances := make([]DeploymentInstance, len(deploy.Instances))
	copy(instances, deploy.Instances)
	return instances, nil
}

// GetHealthyInstances returns healthy instances count
func (dm *DeploymentManager) GetHealthyInstances(deploymentID string) int {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	deploy, exists := dm.deployments[deploymentID]
	if !exists {
		return 0
	}

	count := 0
	for _, instance := range deploy.Instances {
		if instance.Health == "healthy" && instance.Status == "running" {
			count++
		}
	}
	return count
}

