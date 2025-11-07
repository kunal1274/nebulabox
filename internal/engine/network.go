package engine

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// NetworkManager manages container networking
type NetworkManager struct {
	networks   map[string]*Network // network name -> Network
	containers map[string]*NetworkConfig // containerID -> NetworkConfig
}

// Network represents a network
type Network struct {
	ID      string
	Name    string
	Driver  string
	Subnet  string
	Gateway string
	Bridge  string
}

// NewNetworkManager creates a new network manager
func NewNetworkManager() (*NetworkManager, error) {
	return &NetworkManager{
		networks:   make(map[string]*Network),
		containers: make(map[string]*NetworkConfig),
	}, nil
}

// SetupContainerNetwork sets up networking for a container
func (nm *NetworkManager) SetupContainerNetwork(containerID, networkName string, portMappings map[string]string) (*NetworkConfig, error) {
	// Get or create network
	network, err := nm.getOrCreateNetwork(networkName)
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %w", err)
	}

	// Create veth pair
	vethHost, _, err := nm.createVethPair(containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create veth pair: %w", err)
	}

	// Add veth to bridge
	if err := nm.addToBridge(vethHost, network.Bridge); err != nil {
		return nil, fmt.Errorf("failed to add to bridge: %w", err)
	}

	// Assign IP address
	ipAddr, err := nm.assignIPAddress(containerID, network)
	if err != nil {
		return nil, fmt.Errorf("failed to assign IP: %w", err)
	}

	// Setup port forwarding
	if err := nm.setupPortForwarding(containerID, portMappings); err != nil {
		return nil, fmt.Errorf("failed to setup port forwarding: %w", err)
	}

	netConfig := &NetworkConfig{
		Mode:        network.Driver,
		Bridge:      network.Bridge,
		IPAddress:   ipAddr.String(),
		Gateway:     network.Gateway,
		PortMappings: portMappings,
	}

	nm.containers[containerID] = netConfig
	return netConfig, nil
}

// CleanupContainerNetwork cleans up networking for a container
func (nm *NetworkManager) CleanupContainerNetwork(containerID string) error {
	netConfig, exists := nm.containers[containerID]
	if !exists {
		return nil // Already cleaned up
	}

	// Remove port forwarding
	if err := nm.removePortForwarding(containerID, netConfig.PortMappings); err != nil {
		return fmt.Errorf("failed to remove port forwarding: %w", err)
	}

	// Remove veth pair
	if err := nm.removeVethPair(containerID); err != nil {
		return fmt.Errorf("failed to remove veth pair: %w", err)
	}

	delete(nm.containers, containerID)
	return nil
}

// getOrCreateNetwork gets an existing network or creates a new one
func (nm *NetworkManager) getOrCreateNetwork(networkName string) (*Network, error) {
	if networkName == "" {
		networkName = "bridge"
	}

	// Check if network exists
	if network, exists := nm.networks[networkName]; exists {
		return network, nil
	}

	// Create new network
	network, err := nm.createNetwork(networkName)
	if err != nil {
		return nil, err
	}

	nm.networks[networkName] = network
	return network, nil
}

// createNetwork creates a new network
func (nm *NetworkManager) createNetwork(networkName string) (*Network, error) {
	// Generate network ID
	networkID := generateNetworkID()

	// Create bridge
	bridgeName := fmt.Sprintf("nb-%s", networkID[:12])
	if err := nm.createBridge(bridgeName); err != nil {
		return nil, fmt.Errorf("failed to create bridge: %w", err)
	}

	// Assign subnet and gateway
	subnet := "172.20.0.0/16"
	gateway := "172.20.0.1"

	// Configure bridge IP
	if err := nm.configureBridgeIP(bridgeName, gateway, subnet); err != nil {
		return nil, fmt.Errorf("failed to configure bridge IP: %w", err)
	}

	return &Network{
		ID:      networkID,
		Name:    networkName,
		Driver:  "bridge",
		Subnet:  subnet,
		Gateway: gateway,
		Bridge:  bridgeName,
	}, nil
}

// createBridge creates a network bridge
func (nm *NetworkManager) createBridge(bridgeName string) error {
	// Use ip command to create bridge
	cmd := exec.Command("ip", "link", "add", "name", bridgeName, "type", "bridge")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create bridge: %w", err)
	}

	// Bring bridge up
	cmd = exec.Command("ip", "link", "set", bridgeName, "up")
	return cmd.Run()
}

// configureBridgeIP configures the bridge IP address
func (nm *NetworkManager) configureBridgeIP(bridgeName, ip, subnet string) error {
	// Assign IP to bridge
	cmd := exec.Command("ip", "addr", "add", ip+"/16", "dev", bridgeName)
	return cmd.Run()
}

// createVethPair creates a veth pair for container networking
func (nm *NetworkManager) createVethPair(containerID string) (string, string, error) {
	vethHost := fmt.Sprintf("veth%s", containerID[:12])
	vethContainer := "eth0"

	// Create veth pair
	cmd := exec.Command("ip", "link", "add", vethHost, "type", "veth", "peer", "name", vethContainer)
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to create veth pair: %w", err)
	}

	// Bring host side up
	cmd = exec.Command("ip", "link", "set", vethHost, "up")
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to bring host veth up: %w", err)
	}

	return vethHost, vethContainer, nil
}

// addToBridge adds a veth interface to a bridge
func (nm *NetworkManager) addToBridge(vethName, bridgeName string) error {
	cmd := exec.Command("ip", "link", "set", vethName, "master", bridgeName)
	return cmd.Run()
}

// assignIPAddress assigns an IP address to a container
func (nm *NetworkManager) assignIPAddress(containerID string, network *Network) (net.IP, error) {
	// Generate IP from subnet
	// For simplicity, we'll use a deterministic approach
	ip := net.ParseIP(network.Gateway)
	ip[3]++ // Increment last octet

	return ip, nil
}

// setupPortForwarding sets up port forwarding using iptables
func (nm *NetworkManager) setupPortForwarding(containerID string, portMappings map[string]string) error {
	for containerPort, hostPort := range portMappings {
		// Parse ports
		parts := strings.Split(containerPort, ":")
		if len(parts) > 1 {
			containerPort = parts[1]
		}

		// Setup DNAT rule
		cmd := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING",
			"-p", "tcp", "--dport", hostPort,
			"-j", "DNAT", "--to-destination", "172.20.0.2:"+containerPort)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to setup DNAT: %w", err)
		}
	}

	return nil
}

// removePortForwarding removes port forwarding rules
func (nm *NetworkManager) removePortForwarding(containerID string, portMappings map[string]string) error {
	for containerPort, hostPort := range portMappings {
		parts := strings.Split(containerPort, ":")
		if len(parts) > 1 {
			containerPort = parts[1]
		}

		// Remove DNAT rule
		cmd := exec.Command("iptables", "-t", "nat", "-D", "PREROUTING",
			"-p", "tcp", "--dport", hostPort,
			"-j", "DNAT", "--to-destination", "172.20.0.2:"+containerPort)
		cmd.Run() // Ignore errors
	}

	return nil
}

// removeVethPair removes a veth pair
func (nm *NetworkManager) removeVethPair(containerID string) error {
	vethHost := fmt.Sprintf("veth%s", containerID[:12])
	cmd := exec.Command("ip", "link", "delete", vethHost)
	return cmd.Run()
}

// generateNetworkID generates a unique network ID
func generateNetworkID() string {
	// Simple ID generation
	// In production, use proper UUID
	return "network-" + generateNetworkRandomID(12)
}

// generateNetworkRandomID generates a random ID for networks
func generateNetworkRandomID(length int) string {
	// Simplified ID generation
	// In production, use crypto/rand
	return "net123def456"[:length]
}

