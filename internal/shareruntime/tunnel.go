package shareruntime

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"sync"
	"time"
)

// Tunnel represents a secure tunnel for accessing a container port
type Tunnel struct {
	ID          string            `json:"id"`
	WorkspaceID string            `json:"workspaceId"`
	UserID      string            `json:"userId"`
	Username    string            `json:"username"`
	ContainerID string            `json:"containerId"`
	ContainerPort int             `json:"containerPort"`
	HostPort    int               `json:"hostPort"`    // Mapped host port
	TunnelPort  int               `json:"tunnelPort"`  // Tunnel service port
	Protocol    string            `json:"protocol"`    // tcp, udp
	Status      string            `json:"status"`      // active, closed, error
	Token       string            `json:"token"`        // Access token
	AllowedIPs  []string          `json:"allowedIPs,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	LastActivity time.Time        `json:"lastActivity"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// TunnelConnection represents an active connection through a tunnel
type TunnelConnection struct {
	ID        string    `json:"id"`
	TunnelID  string    `json:"tunnelId"`
	RemoteAddr string   `json:"remoteAddr"`
	StartedAt time.Time `json:"startedAt"`
	BytesIn   int64     `json:"bytesIn"`
	BytesOut  int64     `json:"bytesOut"`
	Status    string    `json:"status"` // active, closed
}

// TunnelManager manages secure tunnels for shared containers
type TunnelManager struct {
	tunnels         map[string]*Tunnel              // tunnelID -> Tunnel
	tunnelByToken   map[string]*Tunnel              // token -> Tunnel
	tunnelByPort    map[int]*Tunnel                 // hostPort -> Tunnel
	connections     map[string]*TunnelConnection    // connectionID -> Connection
	tunnelConnections map[string][]string            // tunnelID -> []connectionID
	mu              sync.RWMutex
	portAllocator   *PortAllocator
}

// PortAllocator manages port allocation for tunnels
type PortAllocator struct {
	usedPorts  map[int]bool
	startPort  int
	endPort    int
	mu         sync.Mutex
}

// NewPortAllocator creates a new port allocator
func NewPortAllocator(startPort, endPort int) *PortAllocator {
	return &PortAllocator{
		usedPorts: make(map[int]bool),
		startPort: startPort,
		endPort:  endPort,
	}
}

// AllocatePort allocates an available port
func (pa *PortAllocator) AllocatePort() (int, error) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	for port := pa.startPort; port <= pa.endPort; port++ {
		if !pa.usedPorts[port] {
			// Check if port is actually available
			if pa.isPortAvailable(port) {
				pa.usedPorts[port] = true
				return port, nil
			}
		}
	}

	return 0, fmt.Errorf("no available ports in range %d-%d", pa.startPort, pa.endPort)
}

// ReleasePort releases a port back to the pool
func (pa *PortAllocator) ReleasePort(port int) {
	pa.mu.Lock()
	defer pa.mu.Unlock()
	delete(pa.usedPorts, port)
}

// isPortAvailable checks if a port is actually available
func (pa *PortAllocator) isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// NewTunnelManager creates a new tunnel manager
func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		tunnels:           make(map[string]*Tunnel),
		tunnelByToken:     make(map[string]*Tunnel),
		tunnelByPort:      make(map[int]*Tunnel),
		connections:       make(map[string]*TunnelConnection),
		tunnelConnections: make(map[string][]string),
		portAllocator:     NewPortAllocator(30000, 31000), // Allocate ports in range
	}
}

// CreateTunnel creates a new secure tunnel
func (tm *TunnelManager) CreateTunnel(workspaceID, userID, username, containerID string, containerPort int, protocol string, allowedIPs []string) (*Tunnel, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Allocate a host port
	hostPort, err := tm.portAllocator.AllocatePort()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate port: %w", err)
	}

	// Allocate a tunnel service port
	tunnelPort, err := tm.portAllocator.AllocatePort()
	if err != nil {
		tm.portAllocator.ReleasePort(hostPort)
		return nil, fmt.Errorf("failed to allocate tunnel port: %w", err)
	}

	// Generate secure token
	token, err := generateSecureToken()
	if err != nil {
		tm.portAllocator.ReleasePort(hostPort)
		tm.portAllocator.ReleasePort(tunnelPort)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	tunnel := &Tunnel{
		ID:           fmt.Sprintf("tun-%d", time.Now().UnixNano()),
		WorkspaceID:  workspaceID,
		UserID:       userID,
		Username:     username,
		ContainerID:  containerID,
		ContainerPort: containerPort,
		HostPort:     hostPort,
		TunnelPort:   tunnelPort,
		Protocol:     protocol,
		Status:       "active",
		Token:        token,
		AllowedIPs:   allowedIPs,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
		Metadata:     make(map[string]string),
	}

	tm.tunnels[tunnel.ID] = tunnel
	tm.tunnelByToken[token] = tunnel
	tm.tunnelByPort[hostPort] = tunnel
	tm.tunnelConnections[tunnel.ID] = []string{}

	return tunnel, nil
}

// GetTunnel retrieves a tunnel by ID
func (tm *TunnelManager) GetTunnel(tunnelID string) (*Tunnel, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tunnel, exists := tm.tunnels[tunnelID]
	return tunnel, exists
}

// GetTunnelByToken retrieves a tunnel by access token
func (tm *TunnelManager) GetTunnelByToken(token string) (*Tunnel, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tunnel, exists := tm.tunnelByToken[token]
	return tunnel, exists
}

// GetTunnelByPort retrieves a tunnel by host port
func (tm *TunnelManager) GetTunnelByPort(port int) (*Tunnel, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tunnel, exists := tm.tunnelByPort[port]
	return tunnel, exists
}

// ListTunnels returns tunnels for a workspace
func (tm *TunnelManager) ListTunnels(workspaceID string) []*Tunnel {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var result []*Tunnel
	for _, tunnel := range tm.tunnels {
		if tunnel.WorkspaceID == workspaceID {
			result = append(result, tunnel)
		}
	}
	return result
}

// ListUserTunnels returns tunnels for a user
func (tm *TunnelManager) ListUserTunnels(userID string) []*Tunnel {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var result []*Tunnel
	for _, tunnel := range tm.tunnels {
		if tunnel.UserID == userID {
			result = append(result, tunnel)
		}
	}
	return result
}

// CloseTunnel closes a tunnel
func (tm *TunnelManager) CloseTunnel(tunnelID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tunnel, exists := tm.tunnels[tunnelID]
	if !exists {
		return fmt.Errorf("tunnel not found: %s", tunnelID)
	}

	// Close all connections
	connIDs := tm.tunnelConnections[tunnelID]
	for _, connID := range connIDs {
		if conn, exists := tm.connections[connID]; exists {
			conn.Status = "closed"
		}
	}

	// Release ports
	tm.portAllocator.ReleasePort(tunnel.HostPort)
	tm.portAllocator.ReleasePort(tunnel.TunnelPort)

	// Remove from maps
	delete(tm.tunnels, tunnelID)
	delete(tm.tunnelByToken, tunnel.Token)
	delete(tm.tunnelByPort, tunnel.HostPort)
	delete(tm.tunnelConnections, tunnelID)

	tunnel.Status = "closed"
	return nil
}

// RegisterConnection registers a new connection through a tunnel
func (tm *TunnelManager) RegisterConnection(tunnelID, remoteAddr string) (*TunnelConnection, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tunnel, exists := tm.tunnels[tunnelID]
	if !exists {
		return nil, fmt.Errorf("tunnel not found: %s", tunnelID)
	}

	if tunnel.Status != "active" {
		return nil, fmt.Errorf("tunnel is not active")
	}

	conn := &TunnelConnection{
		ID:         fmt.Sprintf("conn-%d", time.Now().UnixNano()),
		TunnelID:   tunnelID,
		RemoteAddr: remoteAddr,
		StartedAt:  time.Now(),
		Status:     "active",
	}

	tm.connections[conn.ID] = conn
	tm.tunnelConnections[tunnelID] = append(tm.tunnelConnections[tunnelID], conn.ID)
	tunnel.LastActivity = time.Now()

	return conn, nil
}

// CloseConnection closes a tunnel connection
func (tm *TunnelManager) CloseConnection(connectionID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	conn, exists := tm.connections[connectionID]
	if !exists {
		return fmt.Errorf("connection not found: %s", connectionID)
	}

	conn.Status = "closed"

	// Remove from tunnel connections
	if connIDs, exists := tm.tunnelConnections[conn.TunnelID]; exists {
		newIDs := []string{}
		for _, id := range connIDs {
			if id != connectionID {
				newIDs = append(newIDs, id)
			}
		}
		tm.tunnelConnections[conn.TunnelID] = newIDs
	}

	delete(tm.connections, connectionID)
	return nil
}

// ListTunnelConnections returns all connections for a tunnel
func (tm *TunnelManager) ListTunnelConnections(tunnelID string) []*TunnelConnection {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	connIDs, exists := tm.tunnelConnections[tunnelID]
	if !exists {
		return []*TunnelConnection{}
	}

	var result []*TunnelConnection
	for _, connID := range connIDs {
		if conn, exists := tm.connections[connID]; exists {
			result = append(result, conn)
		}
	}
	return result
}

// UpdateConnectionStats updates connection statistics
func (tm *TunnelManager) UpdateConnectionStats(connectionID string, bytesIn, bytesOut int64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	conn, exists := tm.connections[connectionID]
	if !exists {
		return fmt.Errorf("connection not found: %s", connectionID)
	}

	conn.BytesIn += bytesIn
	conn.BytesOut += bytesOut

	// Update tunnel activity
	if tunnel, exists := tm.tunnels[conn.TunnelID]; exists {
		tunnel.LastActivity = time.Now()
	}

	return nil
}

// ValidateTunnelAccess validates if access is allowed for a tunnel
func (tm *TunnelManager) ValidateTunnelAccess(tunnelID, token, remoteIP string) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tunnel, exists := tm.tunnels[tunnelID]
	if !exists {
		return fmt.Errorf("tunnel not found")
	}

	if tunnel.Status != "active" {
		return fmt.Errorf("tunnel is not active")
	}

	if tunnel.Token != token {
		return fmt.Errorf("invalid token")
	}

	// Check IP whitelist if configured
	if len(tunnel.AllowedIPs) > 0 {
		allowed := false
		for _, allowedIP := range tunnel.AllowedIPs {
			if allowedIP == remoteIP || allowedIP == "*" {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("access denied from IP: %s", remoteIP)
		}
	}

	return nil
}

// generateSecureToken generates a secure random token
func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

