package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/orchestrator"
)

// RegisterNodeRequest represents a node registration request
type RegisterNodeRequest struct {
	ID      string            `json:"id" binding:"required"`
	Name    string            `json:"name,omitempty"`
	Address string            `json:"address" binding:"required"`
	Port    int               `json:"port,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
	Region  string            `json:"region,omitempty"`
	Zone    string            `json:"zone,omitempty"`
}

// CreateDeploymentRequest represents a deployment creation request
type CreateDeploymentRequest struct {
	Name          string                `json:"name" binding:"required"`
	Image         string                `json:"image" binding:"required"`
	Tag           string                `json:"tag,omitempty"`
	Replicas      int                   `json:"replicas" binding:"required"`
	NodeSelector  map[string]string     `json:"nodeSelector,omitempty"`
	Strategy      string                `json:"strategy,omitempty"`
	HealthCheck   *orchestrator.HealthCheckConfig `json:"healthCheck,omitempty"`
	AutoRestart   bool                  `json:"autoRestart,omitempty"`
	MaxRestarts   int                   `json:"maxRestarts,omitempty"`
	RestartPolicy string                `json:"restartPolicy,omitempty"`
	ServiceName   string                `json:"serviceName,omitempty"` // Auto-register as service
	NetworkName   string                `json:"networkName,omitempty"` // Network for deployment
	Ports         []string              `json:"ports,omitempty"` // Port mappings
}

// UpdateDeploymentRequest represents a deployment update request
type UpdateDeploymentRequest struct {
	Replicas    *int                  `json:"replicas,omitempty"`
	Image       string                 `json:"image,omitempty"`
	Tag         string                 `json:"tag,omitempty"`
	HealthCheck *orchestrator.HealthCheckConfig `json:"healthCheck,omitempty"`
	AutoRestart *bool                  `json:"autoRestart,omitempty"`
}

// registerNode handles POST /api/orchestrator/nodes
func (s *Server) registerNode(c *gin.Context) {
	var req RegisterNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	node := &orchestrator.Node{
		ID:      req.ID,
		Name:    req.Name,
		Address: req.Address,
		Port:    req.Port,
		Labels:  req.Labels,
		Region:  req.Region,
		Zone:    req.Zone,
		Resources: orchestrator.NodeResources{
			CPUCores: 4, // Default values, would be reported by node
			MemoryMB: 8192,
			DiskGB:   100,
		},
	}

	if err := s.nodeManager.RegisterNode(node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to register node",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, node)
}

// listNodes handles GET /api/orchestrator/nodes
func (s *Server) listNodes(c *gin.Context) {
	onlineOnly := c.Query("online") == "true"
	
	var nodes []*orchestrator.Node
	if onlineOnly {
		nodes = s.nodeManager.ListOnlineNodes()
	} else {
		nodes = s.nodeManager.ListNodes()
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
		"count": len(nodes),
	})
}

// getNode handles GET /api/orchestrator/nodes/:id
func (s *Server) getNode(c *gin.Context) {
	nodeID := c.Param("id")
	node, exists := s.nodeManager.GetNode(nodeID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Node not found",
		})
		return
	}

	c.JSON(http.StatusOK, node)
}

// updateNodeStatus handles PATCH /api/orchestrator/nodes/:id/status
func (s *Server) updateNodeStatus(c *gin.Context) {
	nodeID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := s.nodeManager.UpdateNodeStatus(nodeID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// createOrchestratorDeployment handles POST /api/orchestrator/deployments
func (s *Server) createOrchestratorDeployment(c *gin.Context) {
	var req CreateDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	deploy := &orchestrator.Deployment{
		Name:         req.Name,
		Image:        req.Image,
		Tag:          req.Tag,
		Replicas:     req.Replicas,
		NodeSelector: req.NodeSelector,
		Strategy:     req.Strategy,
		HealthCheck:  req.HealthCheck,
		AutoRestart:  req.AutoRestart,
		MaxRestarts:  req.MaxRestarts,
		RestartPolicy: req.RestartPolicy,
		ServiceName:   req.ServiceName,
		NetworkName:   req.NetworkName,
		Ports:         req.Ports,
	}

	if err := s.deploymentManager.CreateDeployment(deploy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create deployment",
			"details": err.Error(),
		})
		return
	}

	// Create network if specified and doesn't exist
	if req.NetworkName != "" {
		s.netMu.Lock()
		if _, exists := s.networks[req.NetworkName]; !exists {
		// Create network automatically
		s.networks[req.NetworkName] = &Network{
			ID:      req.NetworkName,
			Name:    req.NetworkName,
			Driver:  "bridge",
			Subnet:  "172.20.0.0/16", // Auto-generated subnet
			Created: time.Now(),
		}
		}
		s.netMu.Unlock()
	}

	// Start health monitoring if enabled
	if deploy.AutoRestart && deploy.HealthCheck != nil {
		interval := time.Duration(deploy.HealthCheck.IntervalSeconds) * time.Second
		if interval == 0 {
			interval = 30 * time.Second
		}
		s.healthMonitor.StartMonitoring(deploy.ID, interval)
	}

	c.JSON(http.StatusOK, deploy)
}

// listOrchestratorDeployments handles GET /api/orchestrator/deployments
func (s *Server) listOrchestratorDeployments(c *gin.Context) {
	deployments := s.deploymentManager.ListDeployments()
	c.JSON(http.StatusOK, gin.H{
		"deployments": deployments,
		"count":       len(deployments),
	})
}

// getOrchestratorDeployment handles GET /api/orchestrator/deployments/:id
func (s *Server) getOrchestratorDeployment(c *gin.Context) {
	deploymentID := c.Param("id")
	deploy, exists := s.deploymentManager.GetDeployment(deploymentID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Deployment not found",
		})
		return
	}

	c.JSON(http.StatusOK, deploy)
}

// updateOrchestratorDeployment handles PATCH /api/orchestrator/deployments/:id
func (s *Server) updateOrchestratorDeployment(c *gin.Context) {
	deploymentID := c.Param("id")
	var req UpdateDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	deploy, exists := s.deploymentManager.GetDeployment(deploymentID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Deployment not found",
		})
		return
	}

	// Update fields
	if req.Replicas != nil {
		deploy.Replicas = *req.Replicas
	}
	if req.Image != "" {
		deploy.Image = req.Image
	}
	if req.Tag != "" {
		deploy.Tag = req.Tag
	}
	if req.HealthCheck != nil {
		deploy.HealthCheck = req.HealthCheck
	}
	if req.AutoRestart != nil {
		deploy.AutoRestart = *req.AutoRestart
		if deploy.AutoRestart && deploy.HealthCheck != nil {
			interval := time.Duration(deploy.HealthCheck.IntervalSeconds) * time.Second
			if interval == 0 {
				interval = 30 * time.Second
			}
			s.healthMonitor.StartMonitoring(deploymentID, interval)
		}
	}

	deploy.UpdatedAt = time.Now()
	
	// Re-register instances if service name exists
	if deploy.ServiceName != "" {
		s.integrateDeploymentWithServiceDiscovery(deploy)
	}
	
	c.JSON(http.StatusOK, deploy)
}

// deleteOrchestratorDeployment handles DELETE /api/orchestrator/deployments/:id
func (s *Server) deleteOrchestratorDeployment(c *gin.Context) {
	deploymentID := c.Param("id")

	// Stop monitoring
	s.healthMonitor.StopMonitoring(deploymentID)

	if err := s.deploymentManager.DeleteDeployment(deploymentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete deployment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment deleted"})
}

// scaleOrchestratorDeployment handles POST /api/orchestrator/deployments/:id/scale
func (s *Server) scaleOrchestratorDeployment(c *gin.Context) {
	deploymentID := c.Param("id")
	var req struct {
		Replicas int `json:"replicas" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	deploy, exists := s.deploymentManager.GetDeployment(deploymentID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Deployment not found",
		})
		return
	}

	if req.Replicas < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Replicas must be >= 0",
		})
		return
	}

	deploy.Replicas = req.Replicas
	deploy.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, gin.H{
		"message": "Deployment scaled",
		"replicas": deploy.Replicas,
	})
}

// restartOrchestratorDeployment handles POST /api/orchestrator/deployments/:id/restart
func (s *Server) restartOrchestratorDeployment(c *gin.Context) {
	deploymentID := c.Param("id")

	deploy, exists := s.deploymentManager.GetDeployment(deploymentID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Deployment not found",
		})
		return
	}

	// Record restart for all instances
	for _, instance := range deploy.Instances {
		s.deploymentManager.RecordInstanceRestart(deploymentID, instance.ID)
		s.deploymentManager.UpdateInstanceStatus(deploymentID, instance.ID, "pending")
	}

	deploy.UpdatedAt = time.Now()
	c.JSON(http.StatusOK, gin.H{"message": "Deployment restart initiated"})
}

// integrateDeploymentWithServiceDiscovery integrates deployment with service discovery
func (s *Server) integrateDeploymentWithServiceDiscovery(deploy *orchestrator.Deployment) {
	// Create service discovery integration
	integration := orchestrator.NewServiceDiscoveryIntegration(
		// registerService
		func(name, id, address string, port int, network, version string) error {
			s.svcMu.Lock()
			defer s.svcMu.Unlock()
			
			inst := ServiceInstance{
				ID:        id,
				Name:      name,
				Address:   address,
				Port:      port,
				Version:   version,
				Network:   network,
				CreatedAt: time.Now().Unix(),
			}
			s.services[name] = append(s.services[name], inst)
			return nil
		},
		// deregisterService
		func(name, id string) error {
			s.svcMu.Lock()
			defer s.svcMu.Unlock()
			
			arr := s.services[name]
			out := arr[:0]
			for _, x := range arr {
				if x.ID != id {
					out = append(out, x)
				}
			}
			if len(out) == 0 {
				delete(s.services, name)
			} else {
				s.services[name] = out
			}
			return nil
		},
		// addDNSRecord
		func(name string, addresses []string) error {
			s.dnsMu.Lock()
			defer s.dnsMu.Unlock()
			s.dnsRecords[name] = addresses
			return nil
		},
		// deleteDNSRecord
		func(name string) error {
			s.dnsMu.Lock()
			defer s.dnsMu.Unlock()
			delete(s.dnsRecords, name)
			return nil
		},
	)

	// Register all existing instances
	integration.RegisterAllDeploymentInstances(deploy, s.nodeManager)
	
	// Store integration reference for future instance additions
	// (This would be stored in a map in a real implementation)
}

// autoRegisterInstance registers a new deployment instance with service discovery
func (s *Server) autoRegisterInstance(deploy *orchestrator.Deployment, instance orchestrator.DeploymentInstance) {
	if deploy.ServiceName == "" {
		return
	}

	node, exists := s.nodeManager.GetNode(instance.NodeID)
	if !exists {
		return
	}

	integration := orchestrator.NewServiceDiscoveryIntegration(
		func(name, id, address string, port int, network, version string) error {
			s.svcMu.Lock()
			defer s.svcMu.Unlock()
			inst := ServiceInstance{
				ID:        id,
				Name:      name,
				Address:   address,
				Port:      port,
				Version:   version,
				Network:   network,
				CreatedAt: time.Now().Unix(),
			}
			s.services[name] = append(s.services[name], inst)
			return nil
		},
		func(name, id string) error {
			s.svcMu.Lock()
			defer s.svcMu.Unlock()
			arr := s.services[name]
			out := arr[:0]
			for _, x := range arr {
				if x.ID != id {
					out = append(out, x)
				}
			}
			if len(out) == 0 {
				delete(s.services, name)
			} else {
				s.services[name] = out
			}
			return nil
		},
		func(name string, addresses []string) error {
			s.dnsMu.Lock()
			defer s.dnsMu.Unlock()
			s.dnsRecords[name] = addresses
			return nil
		},
		func(name string) error {
			s.dnsMu.Lock()
			defer s.dnsMu.Unlock()
			delete(s.dnsRecords, name)
			return nil
		},
	)

	integration.RegisterDeploymentInstance(deploy, instance, node.Address)
}

