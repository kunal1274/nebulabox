package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
	"github.com/nebulabox/nebulabox/internal/database/repositories"
)

// ContainerRequest represents a request to create a container
type ContainerRequest struct {
	Image   string            `json:"image" binding:"required"`
	Name    string            `json:"name,omitempty"`
	Port    string            `json:"port,omitempty"`
    Ports   []string          `json:"ports,omitempty"`
	Detach  bool              `json:"detach,omitempty"`
	Env     []string          `json:"env,omitempty"`
	Volume  []string          `json:"volume,omitempty"`
    Network string            `json:"network,omitempty"`
    Service string            `json:"service,omitempty"`
    WorkspaceID string        `json:"workspaceId,omitempty"`
    // Health checks
    HealthType        string   `json:"healthType,omitempty"` // http|tcp|cmd
    HealthHTTPPath    string   `json:"healthHttpPath,omitempty"`
    HealthHTTPPort    string   `json:"healthHttpPort,omitempty"`
    HealthTCPPort     string   `json:"healthTcpPort,omitempty"`
    HealthCmd         []string `json:"healthCmd,omitempty"`
    HealthIntervalSec int      `json:"healthIntervalSec,omitempty"`
    HealthTimeoutSec  int      `json:"healthTimeoutSec,omitempty"`
    HealthRetries     int      `json:"healthRetries,omitempty"`
    HealthStartPeriod int      `json:"healthStartPeriodSec,omitempty"`
}

// ContainerResponse represents a container in API responses
type ContainerResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	Created string `json:"created"`
}

// listContainers handles GET /api/containers
func (s *Server) listContainers(c *gin.Context) {
	// Parse query parameters
	all := c.Query("all") == "true"
	workspaceID := c.Query("workspaceId")

	s.modeMu.Lock()
	mode := s.operatingMode
	s.modeMu.Unlock()

	ctx := c.Request.Context()
	
	// Create map to track unique containers by ID
	containerMap := make(map[string]*containerd.Container)
	
	// Step 1: Get containers from database if available (test/live mode)
	if (mode == "test" || mode == "live") && s.repos != nil && s.repos.Container != nil {
		dbContainers, err := s.repos.Container.List(all, workspaceID)
		if err == nil && len(dbContainers) > 0 {
			for _, container := range dbContainers {
				containerMap[container.ID] = container
			}
		}
	}
	
	// Step 2: Get containers from containerd (real or mock)
	containers, err := s.containerd.ListContainers(ctx)
	if err != nil {
		// Log error but continue (might be mock mode)
		log.Printf("⚠️  Warning: Failed to list containers from containerd: %v", err)
	} else {
		for _, container := range containers {
			// Only add if not already in map (database takes precedence)
			if _, exists := containerMap[container.ID]; !exists {
				containerMap[container.ID] = container
			}
		}
	}

	// Step 3: In test/live mode, merge with in-memory stored containers (fallback)
	if mode == "test" || mode == "live" {
		s.builtContainersMu.Lock()
		for id, container := range s.builtContainers {
			// Only add if not already in map (database and containerd take precedence)
			if _, exists := containerMap[id]; !exists {
				containerMap[id] = container
			}
		}
		s.builtContainersMu.Unlock()
	}

	// Convert map back to slice
	var allContainers []*containerd.Container
	for _, container := range containerMap {
		allContainers = append(allContainers, container)
	}

	// Filter containers based on 'all' parameter and workspace
	var filteredContainers []*containerd.Container
	s.workspaceMu.Lock()
	for _, container := range allContainers {
		if all || container.Status == "running" {
			// If workspace filter specified, check workspace association
			if workspaceID != "" {
				if wsID, ok := s.containerWorkspaces[container.ID]; !ok || wsID != workspaceID {
					continue
				}
			}
			filteredContainers = append(filteredContainers, container)
		}
	}
	s.workspaceMu.Unlock()

	// Convert to API response format
	responses := make([]ContainerResponse, len(filteredContainers))
	for i, container := range filteredContainers {
		responses[i] = ContainerResponse{
			ID:      container.ID,
			Name:    container.Name,
			Image:   container.Image,
			Status:  container.Status,
			Created: container.Created.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, responses)
}

// getContainer handles GET /api/containers/:id
func (s *Server) getContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	s.modeMu.Lock()
	mode := s.operatingMode
	s.modeMu.Unlock()

	ctx := c.Request.Context()
	containers, err := s.containerd.ListContainers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list containers",
			"details": err.Error(),
		})
		return
	}

	// Find container by ID
	var foundContainer *containerd.Container
	for _, container := range containers {
		if container.ID == id {
			foundContainer = container
			break
		}
	}

	// In test/live mode, also check stored containers
	if foundContainer == nil && (mode == "test" || mode == "live") {
		s.builtContainersMu.Lock()
		if stored, exists := s.builtContainers[id]; exists {
			foundContainer = stored
		}
		s.builtContainersMu.Unlock()
	}

	if foundContainer == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Container not found",
		})
		return
	}

	response := ContainerResponse{
		ID:      foundContainer.ID,
		Name:    foundContainer.Name,
		Image:   foundContainer.Image,
		Status:  foundContainer.Status,
		Created: foundContainer.Created.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// runContainer handles POST /api/containers/run
func (s *Server) runContainer(c *gin.Context) {
	var req ContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Validate workspace membership if workspace ID provided
	username, _ := c.Get("username")
	if req.WorkspaceID != "" {
		s.teamMu.Lock()
		if _, isMember := s.teamMembers[req.WorkspaceID][username.(string)]; !isMember {
			s.teamMu.Unlock()
			c.JSON(http.StatusForbidden, gin.H{"error":"not a member of this workspace"})
			return
		}
		s.teamMu.Unlock()
	}
	// Check tenant quota
	s.tenantMu.Lock()
	tenantID, hasTenant := s.userTenants[username.(string)]
	if hasTenant {
		if !s.checkTenantQuota(tenantID, "container") {
			s.tenantMu.Unlock()
			c.JSON(http.StatusForbidden, gin.H{"error":"tenant container quota exceeded"})
			return
		}
	}
	s.tenantMu.Unlock()

	// Pull image first
	if err := s.containerd.PullImage(ctx, req.Image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to pull image",
			"details": err.Error(),
		})
		return
	}

	// Create container options
	containerOpts := &containerd.ContainerOptions{
		Name:        req.Name,
		Ports:       make(map[string]string),
		Environment: make(map[string]string),
		Volumes:     make(map[string]string),
		Detach:      req.Detach,
        HealthCheck: nil,
        Network:     req.Network,
	}

    // optional: register service after start (mock address/port)

    // Parse port mappings
    if req.Port != "" {
        host, container := parsePortPair(req.Port)
        if host != "" && container != "" {
            containerOpts.Ports[container] = host
        }
    }
    if len(req.Ports) > 0 {
        for _, p := range req.Ports {
            host, container := parsePortPair(p)
            if host != "" && container != "" {
                containerOpts.Ports[container] = host
            }
        }
    }

	// Parse environment variables (format: "KEY=VALUE")
	for _, env := range req.Env {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			containerOpts.Environment[parts[0]] = parts[1]
		} else {
			// If no =, treat as key with empty value
			containerOpts.Environment[env] = ""
		}
	}

	// Parse volume mounts (format: "host:container" or "host:container:ro")
	for _, vol := range req.Volume {
		parts := strings.Split(vol, ":")
		if len(parts) >= 2 {
			hostPath := parts[0]
			containerPath := parts[1]
			containerOpts.Volumes[containerPath] = hostPath
		} else {
			// If no colon, treat as container path only
			containerOpts.Volumes[vol] = vol
		}
	}

    // Health check
    if req.HealthType != "" {
        containerOpts.HealthCheck = &containerd.HealthCheckOptions{
            Type:            req.HealthType,
            HTTPPath:        req.HealthHTTPPath,
            HTTPPort:        req.HealthHTTPPort,
            TCPPort:         req.HealthTCPPort,
            Command:         req.HealthCmd,
            IntervalSeconds: req.HealthIntervalSec,
            TimeoutSeconds:  req.HealthTimeoutSec,
            Retries:         req.HealthRetries,
            StartPeriodSec:  req.HealthStartPeriod,
        }
    }

	// Create container
	container, err := s.containerd.CreateContainer(ctx, req.Image, req.Name, containerOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create container",
			"details": err.Error(),
		})
		return
	}

	// Associate container with workspace if provided
	if req.WorkspaceID != "" {
		s.workspaceMu.Lock()
		s.containerWorkspaces[container.ID] = req.WorkspaceID
		s.workspaceMu.Unlock()
	}

	// Store container in test/live mode
	s.modeMu.Lock()
	mode := s.operatingMode
	s.modeMu.Unlock()
	
	// Prepare storage options for database
	storageOpts := &repositories.ContainerStorageOptions{
		Network:     req.Network,
		WorkspaceID: req.WorkspaceID,
		Command:     "", // Command not directly available from ContainerOptions
	}
	
	// Convert ports, env, volumes to storage format
	if len(containerOpts.Ports) > 0 {
		storageOpts.Ports = make([]string, 0, len(containerOpts.Ports))
		for containerPort, hostPort := range containerOpts.Ports {
			storageOpts.Ports = append(storageOpts.Ports, hostPort+":"+containerPort)
		}
	}
	if len(containerOpts.Environment) > 0 {
		storageOpts.Env = make([]string, 0, len(containerOpts.Environment))
		for key, value := range containerOpts.Environment {
			storageOpts.Env = append(storageOpts.Env, key+"="+value)
		}
	}
	if len(containerOpts.Volumes) > 0 {
		storageOpts.Volumes = make([]string, 0, len(containerOpts.Volumes))
		for containerPath, hostPath := range containerOpts.Volumes {
			storageOpts.Volumes = append(storageOpts.Volumes, hostPath+":"+containerPath)
		}
	}
	
	if mode == "test" || mode == "live" {
		// Save to database if repositories available
		if s.repos != nil && s.repos.Container != nil {
			// Create a copy with updated status
			containerCopy := *container
			containerCopy.Status = "created"
			if err := s.repos.Container.CreateOrUpdate(&containerCopy, storageOpts); err != nil {
				log.Printf("⚠️  Warning: Failed to save container to database: %v", err)
				// Continue with in-memory storage as fallback
			}
		}
		
		// Also store in memory for backward compatibility
		s.builtContainersMu.Lock()
		containerCopy := *container
		containerCopy.Status = "created"
		s.builtContainers[container.ID] = &containerCopy
		s.builtContainersMu.Unlock()
	}

	// Start container
	if err := s.containerd.StartContainer(ctx, container.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start container",
			"details": err.Error(),
		})
		return
	}

	// Update status after start
	if mode == "test" || mode == "live" {
		// Update in database
		if s.repos != nil && s.repos.Container != nil {
			if err := s.repos.Container.UpdateStatus(container.ID, "running"); err != nil {
				log.Printf("⚠️  Warning: Failed to update container status in database: %v", err)
			}
		}
		
		// Update in memory
		s.builtContainersMu.Lock()
		if stored, exists := s.builtContainers[container.ID]; exists {
			stored.Status = "running"
		}
		s.builtContainersMu.Unlock()
	}

    // Reserve host ports for this container (best-effort)
    if len(containerOpts.Ports) > 0 {
        s.portsMu.Lock()
        for _, host := range containerOpts.Ports {
            // host may be "hostPort" or "hostPort/proto"; extract number
            n := 0
            for i := 0; i < len(host); i++ { if host[i] == '/' { host = host[:i]; break } }
            for i := 0; i < len(host); i++ { ch := host[i]; if ch < '0' || ch > '9' { n = 0; break } else { n = n*10 + int(ch-'0') } }
            if n > 0 && n <= 65535 {
                if _, used := s.ports[n]; !used { s.ports[n] = container.ID }
            }
        }
        s.portsMu.Unlock()
    }

    // If service name provided, auto-register an instance (mock address)
    if req.Service != "" {
        inst := ServiceInstance{
            ID:        container.ID,
            Name:      req.Service,
            Address:   "127.0.0.1",
            Port:      0,
            Version:   "",
            Network:   req.Network,
            CreatedAt: time.Now().Unix(),
        }
        s.svcMu.Lock()
        s.services[req.Service] = append(s.services[req.Service], inst)
        s.svcMu.Unlock()
    }

	// Get updated container status
	status := "running"
	if mode == "test" || mode == "live" {
		s.builtContainersMu.Lock()
		if stored, exists := s.builtContainers[container.ID]; exists {
			status = stored.Status
		}
		s.builtContainersMu.Unlock()
	}

	response := ContainerResponse{
		ID:      container.ID,
		Name:    container.Name,
		Image:   container.Image,
		Status:  status,
		Created: container.Created.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, response)
}

// getContainerHealth handles GET /api/containers/:id/health
func (s *Server) getContainerHealth(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Container ID is required",
        })
        return
    }

    ctx := c.Request.Context()
    health, err := s.containerd.GetContainerHealth(ctx, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get container health",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": health.Status,
        "lastChecked": health.LastChecked,
        "failingStreak": health.FailingStreak,
        "message": health.Message,
    })
}

// parsePortPair parses strings like "8080:80" or "8080:80/tcp" into host, container
func parsePortPair(v string) (string, string) {
    // find optional "/proto"
    protoIdx := -1
    for i := 0; i < len(v); i++ {
        if v[i] == '/' {
            protoIdx = i
            break
        }
    }
    core := v
    if protoIdx != -1 {
        core = v[:protoIdx]
    }
    // split by ':'
    sep := -1
    for i := 0; i < len(core); i++ {
        if core[i] == ':' {
            sep = i
            break
        }
    }
    if sep == -1 {
        return "", ""
    }
    host := core[:sep]
    container := core[sep+1:]
    return host, container
}

// stopContainer handles POST /api/containers/:id/stop
func (s *Server) stopContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

    ctx := c.Request.Context()
	if err := s.containerd.StopContainer(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to stop container",
			"details": err.Error(),
		})
		return
	}

    // release any reserved ports owned by this container
    s.portsMu.Lock()
    for p, owner := range s.ports { if owner == id { delete(s.ports, p) } }
    s.portsMu.Unlock()

    c.JSON(http.StatusOK, gin.H{
		"message": "Container stopped successfully",
		"id":      id,
	})
}

// startContainer handles POST /api/containers/:id/start
func (s *Server) startContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	ctx := c.Request.Context()
	if err := s.containerd.StartContainer(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start container",
			"details": err.Error(),
		})
		return
	}

	// Update stored container status
	s.modeMu.Lock()
	mode := s.operatingMode
	s.modeMu.Unlock()
	
	if mode == "test" || mode == "live" {
		s.builtContainersMu.Lock()
		if container, exists := s.builtContainers[id]; exists {
			container.Status = "running"
		}
		s.builtContainersMu.Unlock()
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Container started successfully",
		"id":      id,
	})
}

// getContainerLogs handles GET /api/containers/:id/logs
func (s *Server) getContainerLogs(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	ctx := c.Request.Context()
	
	// Try to get logs from MongoDB first
	if s.mongoRepos != nil && s.mongoRepos.ContainerLogs != nil {
		dbLogs, err := s.mongoRepos.ContainerLogs.GetLatest(ctx, id, 1000)
		if err == nil && len(dbLogs) > 0 {
			// Convert to API format
			logs := make([]LogEntry, len(dbLogs))
			for i, dbLog := range dbLogs {
				logs[i] = LogEntry{
					Timestamp: dbLog.Timestamp.Unix(),
					Container: dbLog.ContainerID,
					Level:     dbLog.Level,
					Message:   dbLog.Message,
				}
			}
			c.JSON(http.StatusOK, logs)
			return
		}
		// Log error but continue with fallback
		if err != nil {
			log.Printf("⚠️  Warning: Failed to get logs from MongoDB: %v", err)
		}
	}

	// Fallback: get logs from containerd
	logs, err := s.containerd.GetContainerLogs(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get container logs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}
