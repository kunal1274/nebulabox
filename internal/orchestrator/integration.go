package orchestrator

import (
	"fmt"
	"strconv"
	"strings"
)

// ServiceDiscoveryIntegration handles integration with service discovery
type ServiceDiscoveryIntegration struct {
	registerService    func(name, id, address string, port int, network, version string) error
	deregisterService  func(name, id string) error
	addDNSRecord       func(name string, addresses []string) error
	deleteDNSRecord    func(name string) error
}

// NewServiceDiscoveryIntegration creates a new integration handler
func NewServiceDiscoveryIntegration(
	registerService func(name, id, address string, port int, network, version string) error,
	deregisterService func(name, id string) error,
	addDNSRecord func(name string, addresses []string) error,
	deleteDNSRecord func(name string) error,
) *ServiceDiscoveryIntegration {
	return &ServiceDiscoveryIntegration{
		registerService:   registerService,
		deregisterService: deregisterService,
		addDNSRecord:      addDNSRecord,
		deleteDNSRecord:   deleteDNSRecord,
	}
}

// RegisterDeploymentInstance registers a deployment instance as a service
func (sdi *ServiceDiscoveryIntegration) RegisterDeploymentInstance(
	deploy *Deployment,
	instance DeploymentInstance,
	nodeAddress string,
) error {
	if deploy.ServiceName == "" {
		return nil // No service registration needed
	}

	// Extract port from deployment ports or use default
	port := 8080 // Default port
	if len(deploy.Ports) > 0 {
		// Parse first port mapping (format: "host:container" or "container")
		firstPort := deploy.Ports[0]
		if strings.Contains(firstPort, ":") {
			parts := strings.Split(firstPort, ":")
			if len(parts) >= 2 {
				if p, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
					port = p
				}
			}
		} else {
			if p, err := strconv.Atoi(strings.TrimSpace(firstPort)); err == nil {
				port = p
			}
		}
	}

	// Use node address or instance address
	address := nodeAddress
	if address == "" {
		address = instance.ContainerID // Fallback to container ID
	}

	// Register as service instance
	if sdi.registerService != nil {
		err := sdi.registerService(
			deploy.ServiceName,
			instance.ID,
			address,
			port,
			deploy.NetworkName,
			deploy.Tag,
		)
		if err != nil {
			return fmt.Errorf("failed to register service: %w", err)
		}
	}

	// Add DNS record for service name
	if sdi.addDNSRecord != nil {
		err := sdi.addDNSRecord(
			deploy.ServiceName,
			[]string{address},
		)
		if err != nil {
			// DNS record creation is best-effort
			// Log but don't fail registration
		}
	}

	return nil
}

// DeregisterDeploymentInstance removes a deployment instance from service discovery
func (sdi *ServiceDiscoveryIntegration) DeregisterDeploymentInstance(
	deploy *Deployment,
	instance DeploymentInstance,
) error {
	if deploy.ServiceName == "" {
		return nil // No service to deregister
	}

	// Deregister service instance
	if sdi.deregisterService != nil {
		if err := sdi.deregisterService(deploy.ServiceName, instance.ID); err != nil {
			return fmt.Errorf("failed to deregister service: %w", err)
		}
	}

	return nil
}

// RegisterAllDeploymentInstances registers all instances of a deployment
func (sdi *ServiceDiscoveryIntegration) RegisterAllDeploymentInstances(
	deploy *Deployment,
	nodeManager *NodeManager,
) error {
	for _, instance := range deploy.Instances {
		if instance.Status != "running" {
			continue // Only register running instances
		}

		// Get node address
		var nodeAddress string
		if instance.NodeID != "" {
			node, exists := nodeManager.GetNode(instance.NodeID)
			if exists {
				nodeAddress = node.Address
			}
		}

		if err := sdi.RegisterDeploymentInstance(deploy, instance, nodeAddress); err != nil {
			// Continue with other instances even if one fails
			continue
		}
	}

	return nil
}

