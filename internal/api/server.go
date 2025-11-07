package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
    "sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/aiops"
	"github.com/nebulabox/nebulabox/internal/cloud"
	"github.com/nebulabox/nebulabox/internal/composition"
	"github.com/nebulabox/nebulabox/internal/containerd"
	"github.com/nebulabox/nebulabox/internal/database"
	"github.com/nebulabox/nebulabox/internal/database/mongodb_repositories"
	"github.com/nebulabox/nebulabox/internal/database/repositories"
	"github.com/nebulabox/nebulabox/internal/groups"
	"github.com/nebulabox/nebulabox/internal/orchestrator"
	"github.com/nebulabox/nebulabox/internal/runtime"
	"github.com/nebulabox/nebulabox/internal/security"
	"github.com/nebulabox/nebulabox/internal/shareruntime"
	"github.com/nebulabox/nebulabox/internal/snapshot"
	"github.com/nebulabox/nebulabox/internal/templates"
)

// Server represents the NebulaBox API server
type Server struct {
	router     *gin.Engine
	containerd *containerd.Client
	port       string
    registryURL string
    registryClient *RegistryClient
    perf       *perfStats
    alerts     AlertThresholds
    alertsMu   sync.Mutex
    sysHist    *systemHistory
    networks   map[string]*Network
    netMu      sync.Mutex
    services   map[string][]ServiceInstance
    svcMu      sync.Mutex
    svcRR      map[string]int
    dnsRecords map[string][]string
    dnsMu      sync.Mutex
    ports      map[int]string // hostPort -> containerID
    portsMu    sync.Mutex
    githubSecret string
    ghEvents     []GitHubEvent
    ghMu         sync.Mutex
    gitlabSecret string
    glEvents     []GitLabEvent
    glMu         sync.Mutex
    builds       []Build
    buildsMu     sync.Mutex
    tests        []TestRun
    testsMu      sync.Mutex
    deployments  []Deployment
    deployMu     sync.Mutex
    rollbacks    []Rollback
    rbMu         sync.Mutex
    users        map[string]string // username -> password (mock)
    userRoles    map[string]string // username -> role
    sessions     map[string]string // token -> username
    authMu       sync.Mutex
    teams        map[string]*Team // teamID -> Team
    teamMembers  map[string]map[string]*TeamMember // teamID -> username -> TeamMember
    teamMu       sync.Mutex
    containerWorkspaces map[string]string // containerID -> workspaceID
    networkWorkspaces   map[string]string // networkID -> workspaceID
    workspaceMu        sync.Mutex
    tenants      map[string]*Tenant // tenantID -> Tenant
    userTenants  map[string]string // username -> tenantID
    tenantMu     sync.Mutex
    endpointMetrics *EndpointMetrics
    keyManager   *security.KeyManager
    nodeManager  *orchestrator.NodeManager
    deploymentManager *orchestrator.DeploymentManager
    healthMonitor *orchestrator.HealthMonitor
    nebulaRuntime runtime.Runtime
    aiOpsAnalytics *aiops.AnalyticsEngine
    aiOpsScaling *aiops.ScalingAdvisor
    aiOpsChat *aiops.ChatOpsEngine
    groupManager *groups.GroupManager
    relationshipManager *groups.RelationshipManager
    compositionManager *composition.CompositionManager
    templateManager *templates.TemplateManager
    sharedRuntimeManager *shareruntime.WorkspaceManager
    sharedRuntimeInviteManager *shareruntime.InviteManager
    tunnelManager *shareruntime.TunnelManager
    auditLogger *shareruntime.AuditLogger
    syncManager *shareruntime.SyncManager
    fileSyncManager *shareruntime.FileSyncManager
    snapshotManager *snapshot.SnapshotManager
    crdtManager *shareruntime.ConflictResolutionManager
    ephemeralRuntimeManager *cloud.EphemeralRuntimeManager
    autoSleepManager *shareruntime.AutoSleepManager
    // Image storage for built images (test mode)
    builtImages map[string]*ImageResponse // tag -> image
    builtImagesMu sync.Mutex
    // Container storage for created containers (test mode)
    builtContainers map[string]*containerd.Container // id -> container
    builtContainersMu sync.Mutex
    // Operating mode: "mock", "test", "live"
    operatingMode string
    modeMu sync.Mutex
    // Database repositories (nil if database not available)
    repos *repositories.Repositories
    // MongoDB repositories (nil if MongoDB not available)
    mongoRepos *mongodb_repositories.Repositories
}

// NewServer creates a new API server instance
func NewServer() (*Server, error) {
	// Initialize containerd client
	client, err := containerd.NewClient()
	if err != nil {
		return nil, err
	}

	// Create Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

    registryURL := os.Getenv("NEBULABOX_REGISTRY_URL")
    if registryURL == "" {
        registryURL = "http://localhost:5001"
    }

    registryClient := NewRegistryClient(registryURL)
    // Auto-login with admin credentials
    if adminUser := os.Getenv("NEBULABOX_ADMIN_USER"); adminUser != "" {
        adminPass := os.Getenv("NEBULABOX_ADMIN_PASS")
        if adminPass == "" {
            adminPass = "admin"
        }
        if token, err := registryClient.Login(adminUser, adminPass); err == nil {
            log.Printf("✅ Authenticated with registry as %s", adminUser)
            _ = token // Token is stored in client
        } else {
            log.Printf("⚠️  Failed to authenticate with registry: %v", err)
        }
    }

    // Determine operating mode
    mode := os.Getenv("NEBULABOX_MODE")
    if mode == "" {
        mode = "test" // Default to test mode for UAT sandbox
    }
    if mode != "mock" && mode != "test" && mode != "live" {
        mode = "test"
    }

    server := &Server{
        router:         router,
        containerd:     client,
        port:           "8081",
        registryURL:    registryURL,
        registryClient: registryClient,
        builtImages:    make(map[string]*ImageResponse),
        builtContainers: make(map[string]*containerd.Container),
        operatingMode:  mode,
        perf:           &perfStats{},
        alerts:      AlertThresholds{CPUHigh: 85, MemoryHigh: 80, DiskHigh: 90},
        sysHist:     newSystemHistory(24*60*60), // keep up to 24h at 1s resolution
        networks:    make(map[string]*Network),
        services:    make(map[string][]ServiceInstance),
        svcRR:       make(map[string]int),
        dnsRecords:  make(map[string][]string),
        ports:       make(map[int]string),
        githubSecret: os.Getenv("NEBULABOX_GITHUB_SECRET"),
        gitlabSecret: os.Getenv("NEBULABOX_GITLAB_SECRET"),
        users:        map[string]string{},
        userRoles:    map[string]string{},
        sessions:     map[string]string{},
        teams:        map[string]*Team{},
        teamMembers:  map[string]map[string]*TeamMember{},
        containerWorkspaces: map[string]string{},
        networkWorkspaces:   map[string]string{},
        tenants:      map[string]*Tenant{},
        userTenants:  map[string]string{},
        endpointMetrics: NewEndpointMetrics(),
        keyManager:   security.NewKeyManager(),
        nodeManager:  orchestrator.NewNodeManager(),
    }

    // Initialize deployment manager and health monitor
    server.deploymentManager = orchestrator.NewDeploymentManager(server.nodeManager)
    server.healthMonitor = orchestrator.NewHealthMonitor(server.nodeManager, server.deploymentManager)
    
    // Initialize Nebula runtime
    runtimeBasePath := os.Getenv("NEBULA_RUNTIME_PATH")
    if runtimeBasePath == "" {
        runtimeBasePath = "/var/lib/nebula-runtime"
    }
    server.nebulaRuntime = runtime.NewRuntime(runtimeBasePath)

    // Initialize database repositories (if available)
    repos, err := repositories.InitRepositories()
    if err != nil {
        log.Printf("⚠️  Warning: Failed to initialize repositories: %v", err)
        log.Println("   Continuing with in-memory storage only")
    } else if repos != nil {
        log.Println("✅ Database repositories initialized")
        server.repos = repos
    } else {
        log.Println("ℹ️  No database repositories available, using in-memory storage")
    }

    // Initialize MongoDB repositories (if available)
    mongoDB := database.GetMongoDB()
    if mongoDB != nil && mongoDB.DB != nil {
        mongoRepos := mongodb_repositories.InitMongoRepositories(mongoDB.DB)
        if mongoRepos != nil {
            log.Println("✅ MongoDB repositories initialized")
            server.mongoRepos = mongoRepos
        } else {
            log.Println("ℹ️  No MongoDB repositories available, logs/metrics won't be persisted")
        }
    } else {
        log.Println("ℹ️  MongoDB not available, logs/metrics won't be persisted")
    }

    // Initialize AI Ops components
    server.aiOpsAnalytics = aiops.NewAnalyticsEngine()
    server.aiOpsScaling = aiops.NewScalingAdvisor(server.aiOpsAnalytics)
    server.aiOpsChat = aiops.NewChatOpsEngine(server.aiOpsAnalytics, server.aiOpsScaling)

    // Initialize container groups and relationships
    server.groupManager = groups.NewGroupManager()
    server.relationshipManager = groups.NewRelationshipManager()

    // Initialize container composition
    server.compositionManager = composition.NewCompositionManager()

    // Initialize stack templates
    server.templateManager = templates.NewTemplateManager()

    // Initialize Shared Runtime
    server.sharedRuntimeManager = shareruntime.NewWorkspaceManager()
    server.sharedRuntimeInviteManager = shareruntime.NewInviteManager()
    server.tunnelManager = shareruntime.NewTunnelManager()
    server.auditLogger = shareruntime.NewAuditLogger(10000)
    
    // Initialize NebulaSync
    replicationAdapter := shareruntime.NewInMemoryReplicationAdapter()
    server.syncManager = shareruntime.NewSyncManager(replicationAdapter)
    
    // Initialize FileSync
    fileSyncAdapter := shareruntime.NewInMemoryFileSyncAdapter()
    server.fileSyncManager = shareruntime.NewFileSyncManager(fileSyncAdapter)
    
    // Initialize Snapshot Manager
    server.snapshotManager = snapshot.NewSnapshotManager()
    
    // Initialize Auto-Sleep Manager (after workspace and snapshot managers)
    snapshotAdapter := &SnapshotManagerAdapter{snapshotManager: server.snapshotManager}
    server.autoSleepManager = shareruntime.NewAutoSleepManager(server.sharedRuntimeManager, snapshotAdapter)
    server.autoSleepManager.Start()
    
    // Initialize CRDT/Conflict Resolution Manager
    nodeID := os.Getenv("NEBULABOX_NODE_ID")
    if nodeID == "" {
        nodeID = "node-1"
    }
    server.crdtManager = shareruntime.NewConflictResolutionManager(nodeID)
    
    // Initialize Ephemeral Runtime Manager
    idleTimeout := 30 * time.Minute // Default idle timeout
    server.ephemeralRuntimeManager = cloud.NewEphemeralRuntimeManager(idleTimeout)
    
    // Start background tasks for idle/expired runtime checks
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()
        for range ticker.C {
            server.ephemeralRuntimeManager.CheckIdleRuntimes()
            expired := server.ephemeralRuntimeManager.CheckExpiredRuntimes()
            for _, id := range expired {
                server.ephemeralRuntimeManager.TerminateRuntime(id)
            }
        }
    }()

    // init default user (mock)
    adminUser := os.Getenv("NEBULABOX_ADMIN_USER")
    if adminUser == "" { adminUser = "admin" }
    adminPass := os.Getenv("NEBULABOX_ADMIN_PASS")
    if adminPass == "" { adminPass = "admin" }
    server.users[adminUser] = adminPass
    server.userRoles[adminUser] = "admin"

	// Setup routes
	server.setupRoutes()
	server.setupModeRoutes()

	return server, nil
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
    // perf middleware
    s.router.Use(s.perfMiddleware())

    api := s.router.Group("/api")
	{
		// Container routes
		containers := api.Group("/containers")
		{
			containers.GET("", s.listContainers)
			containers.POST("/run", s.runContainer)
			containers.GET("/:id", s.getContainer)
			containers.POST("/:id/start", s.startContainer)
			containers.POST("/:id/stop", s.stopContainer)
			containers.GET("/:id/logs", s.getContainerLogs)
			containers.GET("/:id/health", s.getContainerHealth)
			containers.POST("/:id/exec", s.execContainer)
			containers.POST("/:id/exec/stream", s.execContainerStream)
			containers.POST("/:id/exec/shell", s.execContainerShell)
			containers.POST("/:id/files/upload", s.uploadFile)
			containers.POST("/:id/files/download", s.downloadFile)
			containers.POST("/:id/files/list", s.listFiles)
			containers.POST("/:id/files/upload/stream", s.uploadFileStream)
			containers.GET("/:id/files/download/stream", s.downloadFileStream)
			containers.POST("/:id/env", s.setEnvVars)
			containers.GET("/:id/env", s.getEnvVars)
			containers.PUT("/:id/env", s.updateEnvVar)
			containers.DELETE("/:id/env", s.clearEnvVars)
			containers.POST("/:id/env/string", s.setEnvFromString)
			containers.GET("/:id/env/string", s.getEnvAsString)
			containers.POST("/:id/env/validate", s.validateEnvVar)
			containers.POST("/:id/env/parse", s.parseEnvString)
		}

		// Image routes
		images := api.Group("/images")
		{
			images.GET("", s.listImages)
			images.POST("/pull", s.pullImage)
			images.POST("/push", s.pushImage)
			images.POST("/build", s.buildImage)
			images.POST("/scan", s.postScanImage)
		}

		// Build specification endpoints
		buildspec := api.Group("/buildspec")
		{
			buildspec.POST("/validate", s.validateBuildSpec)
			buildspec.POST("/convert", s.convertBuildSpec)
			buildspec.POST("/build", s.buildFromSpec)
		}

		// Security endpoints (signing, scanning)
		security := api.Group("/security")
		{
			security.POST("/sign", s.signImage)
			security.POST("/verify", s.verifySignature)
			security.GET("/keys", s.listKeys)
			security.GET("/keys/:keyId", s.getKey)
			security.POST("/keys/generate", s.generateKey)
		}

		// Orchestrator endpoints (multi-node deployment)
		orch := api.Group("/orchestrator")
		{
			// Node management
			orch.POST("/nodes", s.registerNode)
			orch.GET("/nodes", s.listNodes)
			orch.GET("/nodes/:id", s.getNode)
			orch.PATCH("/nodes/:id/status", s.updateNodeStatus)
			
			// Deployment management
			orch.POST("/deployments", s.createOrchestratorDeployment)
			orch.GET("/deployments", s.listOrchestratorDeployments)
			orch.GET("/deployments/:id", s.getOrchestratorDeployment)
			orch.PATCH("/deployments/:id", s.updateOrchestratorDeployment)
			orch.DELETE("/deployments/:id", s.deleteOrchestratorDeployment)
			orch.POST("/deployments/:id/scale", s.scaleOrchestratorDeployment)
			orch.POST("/deployments/:id/restart", s.restartOrchestratorDeployment)
		}

		// Nebula Runtime endpoints (custom container runtime)
		rt := api.Group("/runtime")
		{
			// Container management
			rt.POST("/containers", s.createRuntimeContainer)
			rt.GET("/containers", s.listRuntimeContainers)
			rt.GET("/containers/:id", s.getRuntimeContainer)
			rt.POST("/containers/:id/start", s.startRuntimeContainer)
			rt.POST("/containers/:id/stop", s.stopRuntimeContainer)
			rt.DELETE("/containers/:id", s.deleteRuntimeContainer)
			
			// Image management
			rt.POST("/images/pull", s.pullRuntimeImage)
			rt.GET("/images", s.listRuntimeImages)
			
			// Runtime info
			rt.GET("/info", s.getRuntimeInfo)
			rt.GET("/version", s.getRuntimeVersion)
		}

		// AI Ops endpoints (predictive analytics and auto-scaling)
		aiops := api.Group("/aiops")
		{
			// Metrics recording
			aiops.POST("/metrics", s.recordMetric)
			
			// Predictions
			aiops.GET("/predict/:containerId", s.predictResourceUsage)
			
			// Scaling recommendations
			aiops.GET("/scaling/:targetId", s.getScalingRecommendation)
			aiops.POST("/scaling/policy", s.setScalingPolicy)
			
			// ChatOps
			aiops.POST("/chat", s.processChatCommand)
			aiops.GET("/chat/history", s.getChatHistory)
		}

		// Container Groups and Hierarchy endpoints
		g := api.Group("/groups")
		{
			// Group management
			g.POST("", s.createGroup)
			g.GET("", s.listGroups)
			g.GET("/:id", s.getGroup)
			g.PATCH("/:id", s.updateGroup)
			g.DELETE("/:id", s.deleteGroup)
			g.GET("/:id/hierarchy", s.getGroupHierarchy)
			
			// Container membership
			g.POST("/:id/containers", s.addContainerToGroup)
			g.DELETE("/:id/containers/:containerId", s.removeContainerFromGroup)
			
			// Group operations
			g.POST("/:id/start", s.startGroup)
			g.POST("/:id/stop", s.stopGroup)
			
			// Relationships
			g.POST("/relationships", s.createRelationship)
			g.DELETE("/relationships", s.deleteRelationship)
			g.GET("/relationships/:containerId/children", s.getContainerChildren)
			g.GET("/relationships/:containerId/parent", s.getContainerParent)
			g.GET("/relationships/:containerId/tree", s.getContainerTree)
			g.GET("/relationships/:containerId/ancestry", s.getContainerAncestry)
		}

		// Container Composition endpoints
		comp := api.Group("/composition")
		{
			// Spec management
			comp.POST("/specs", s.createCompositionSpec)
			comp.GET("/specs", s.listCompositionSpecs)
			comp.GET("/specs/:name", s.getCompositionSpec)
			comp.DELETE("/specs/:name", s.deleteCompositionSpec)
			
			// Composition operations
			comp.POST("/preview", s.previewComposition)
			comp.POST("/compose", s.composeContainerFromSpec)
		}

		// Snapshot endpoints
		snap := api.Group("/snapshots")
		{
			snap.POST("", s.createSnapshot)
			snap.GET("", s.listSnapshots)
			snap.GET("/:id", s.getSnapshot)
			snap.DELETE("/:id", s.deleteSnapshot)
			snap.POST("/:id/restore", s.restoreSnapshot)
			snap.GET("/resource/:resourceId", s.listResourceSnapshots)
		}

		// Stack Templates endpoints
		tmpl := api.Group("/templates")
		{
			tmpl.GET("", s.listTemplates)
			tmpl.GET("/:id", s.getTemplate)
			tmpl.POST("", s.saveTemplate)
			tmpl.DELETE("/:id", s.deleteTemplate)
			tmpl.POST("/:id/deploy", s.deployTemplate)
		}

		// Shared Runtime endpoints
		sr := api.Group("/shareruntime")
		{
			// Workspace management
			sr.POST("/workspaces", s.createWorkspace)
			sr.GET("/workspaces", s.listWorkspaces)
			sr.GET("/workspaces/:id", s.getWorkspace)
			sr.DELETE("/workspaces/:id", s.deleteWorkspace)
			sr.PATCH("/workspaces/:id/status", s.updateWorkspaceStatus)
			
			// Member management
			sr.POST("/workspaces/:id/members", s.addWorkspaceMember)
			sr.DELETE("/workspaces/:id/members/:userId", s.removeWorkspaceMember)
			sr.PATCH("/workspaces/:id/members/:userId/role", s.updateWorkspaceMemberRole)
			
			// Invites
			sr.POST("/workspaces/:id/invites", s.createInvite)
			sr.GET("/workspaces/:id/invites", s.listInvites)
			sr.GET("/workspaces/:id/invites/:token/link", s.getWorkspaceInviteLink)
			sr.GET("/invites/:token", s.getInviteLink)
			sr.GET("/invites/:token/info", s.getInviteInfo)
			sr.POST("/invites/:token/accept", s.acceptInvite)
			
			// Permissions
			sr.GET("/workspaces/:id/permissions", s.getUserPermissions)
			
			// Audit Logs
			sr.GET("/workspaces/:id/audit-logs", s.getAuditLogs)
			sr.GET("/workspaces/:id/audit-stats", s.getAuditStats)
			sr.GET("/audit-logs", s.getUserAuditLogs)
			
			// NebulaSync
			sr.GET("/workspaces/:id/sync/changes", s.getChangesSince)
			sr.GET("/workspaces/:id/sync/latest", s.getLatestChangeID)
			sr.GET("/workspaces/:id/sync/subscribe", s.subscribeToChanges)
			sr.POST("/workspaces/:id/sync/apply", s.syncWorkspace)
			
			// CRDT/Conflict Resolution
			sr.POST("/workspaces/:id/crdt/operations", s.recordCRDTOperation)
			sr.GET("/workspaces/:id/crdt/operations", s.getCRDTOperations)
			sr.POST("/workspaces/:id/crdt/conflicts/detect", s.detectConflicts)
			sr.POST("/workspaces/:id/crdt/conflicts/:conflictId/resolve", s.resolveConflict)
			sr.GET("/workspaces/:id/crdt/resources/:resourceId", s.getResourceState)
			
			// Auto-Sleep
			sr.PUT("/workspaces/:id/autosleep/config", s.setAutoSleepConfig)
			sr.GET("/workspaces/:id/autosleep/config", s.getAutoSleepConfig)
			sr.POST("/workspaces/:id/activity", s.recordWorkspaceActivity)
			sr.POST("/workspaces/:id/wake", s.wakeWorkspace)
			sr.GET("/autosleep/idle", s.getIdleWorkspaces)
			
			// Sessions
			sr.POST("/workspaces/:id/sessions", s.createSession)
			sr.GET("/workspaces/:id/sessions", s.listWorkspaceSessions)
			sr.DELETE("/sessions/:id", s.closeSession)
			sr.POST("/sessions/:id/activity", s.updateSessionActivity)
			sr.PATCH("/sessions/:id/state", s.updateSessionState)
			sr.GET("/sessions/:id/state", s.getSessionState)
			sr.GET("/workspaces/:id/active-sessions", s.listActiveSessions)
			sr.GET("/workspaces/:id/activity-stream", s.subscribeSessionActivity)
			
			// Tunnels
			sr.POST("/workspaces/:id/tunnels", s.createTunnel)
			sr.GET("/workspaces/:id/tunnels", s.listTunnels)
			sr.GET("/tunnels", s.getUserTunnels)
			sr.GET("/tunnels/:id", s.getTunnel)
			sr.DELETE("/tunnels/:id", s.closeTunnel)
			sr.GET("/tunnels/:id/connections", s.listTunnelConnections)
			sr.POST("/tunnels/:id/validate", s.validateTunnelAccess)
			sr.POST("/tunnels/:id/connect", s.connectTunnel)
			sr.GET("/tunnels/by-port/:port", s.getTunnelByPort)
			sr.PATCH("/connections/:id/stats", s.updateConnectionStats)
		}

		// Registry routes (proxy and management)
		registry := api.Group("/registry")
		{
			registry.GET("/catalog", s.getRegistryCatalog)
			registry.GET("/repositories", s.listRegistryRepositories)
			registry.GET("/repositories/:repo/versions", s.listRegistryVersions)
			registry.GET("/repositories/:repo/summary", s.getRegistrySummary)
			registry.DELETE("/tags/*path", s.deleteRegistryTag) // Must be before GET with *path
			registry.GET("/tags/*path", s.getRegistryTags)
			registry.POST("/retag", s.postRegistryRetag)
			registry.POST("/login", s.postRegistryLogin)
		}

		// Cloud/Ephemeral Runtime endpoints
		cloud := api.Group("/cloud")
		{
			ephemeral := cloud.Group("/ephemeral")
			{
				ephemeral.POST("/runtimes", s.provisionEphemeralRuntime)
				ephemeral.GET("/runtimes", s.listEphemeralRuntimes)
				ephemeral.GET("/runtimes/:id", s.getEphemeralRuntime)
				ephemeral.DELETE("/runtimes/:id", s.terminateEphemeralRuntime)
				ephemeral.POST("/runtimes/:id/activity", s.updateEphemeralRuntimeActivity)
				ephemeral.POST("/runtimes/:id/sleep", s.sleepEphemeralRuntime)
				ephemeral.POST("/runtimes/:id/wake", s.wakeEphemeralRuntime)
				ephemeral.GET("/runtimes/:id/health", s.checkEphemeralRuntimeHealth)
				ephemeral.POST("/runtimes/:id/members", s.addEphemeralRuntimeMember)
				ephemeral.DELETE("/runtimes/:id/members/:userId", s.removeEphemeralRuntimeMember)
			}
		}

		// System routes
		system := api.Group("/system")
		{
			system.GET("/stats", s.getSystemStats)
			system.GET("/stream", s.getSystemStream)
			system.GET("/history", s.getSystemHistory)
		}

        // Network routes
        nets := api.Group("/networks")
        {
            nets.GET("", s.listNetworks)
            nets.POST("", s.authMiddleware(), s.requireRole("admin", "editor"), s.createNetwork)
            nets.DELETE(":id", s.authMiddleware(), s.requireRole("admin", "editor"), s.deleteNetwork)
        }

        // Service discovery routes
        svc := api.Group("/services")
        {
            svc.GET("", s.listServices)
            svc.GET("/resolve/:name", s.resolveService)
            svc.GET("/next/:name", s.resolveServiceNext)
            svc.POST("/register", s.registerService)
            svc.POST("/deregister", s.deregisterService)
        }

        // DNS-style resolution routes
        dns := api.Group("/dns")
        {
            dns.GET("/records", s.listDNSRecords)
            dns.POST("/records", s.addDNSRecord)
            dns.DELETE("/records/:name", s.deleteDNSRecord)
            dns.GET("/resolve/:name", s.dnsResolve)
        }

        // Ports management routes
        ports := api.Group("/ports")
        {
            ports.GET("", s.listPorts)
            ports.POST("/reserve", s.reservePort)
            ports.POST("/release", s.releasePort)
            ports.GET("/suggest", s.suggestPort)
        }

		// Logs routes
		logs := api.Group("/logs")
		{
			logs.GET("/search", s.searchLogs)
			logs.GET("/stream", s.streamLogs)
		}

		// Environment variable routes
		env := api.Group("/env")
		{
			env.GET("/templates", s.getEnvTemplates)
		}

		// Health check
		api.GET("/health", s.healthCheck)

		// Perf routes
		api.GET("/perf/metrics", s.getPerfMetrics)
		api.GET("/perf/stream", s.streamPerf)
		api.GET("/perf/endpoints", s.getEndpointMetrics)

		// Alerts routes
		api.GET("/alerts", s.getAlertsConfig)
		api.POST("/alerts", s.postAlertsConfig)
		api.GET("/alerts/stream", s.streamAlerts)

        // Webhooks
        api.POST("/webhooks/github", s.postGitHubWebhook)
        api.GET("/webhooks/github/events", s.getGitHubEvents)
        api.POST("/webhooks/gitlab", s.postGitLabWebhook)
        api.GET("/webhooks/gitlab/events", s.getGitLabEvents)

        // Builds
        api.GET("/builds", s.listBuilds)
        api.POST("/builds/trigger", s.postTriggerBuild)

        // Tests
        api.GET("/tests", s.listTests)
        api.POST("/tests/run", s.postRunTests)

        // Deployments
        api.GET("/deployments", s.listDeployments)
        api.POST("/deployments/trigger", s.postTriggerDeployment)
        api.GET("/deployments/rollbacks", s.listRollbacks)
        api.POST("/deployments/rollback", s.postTriggerRollback)

        // Teams/Workspaces
        teams := api.Group("/teams")
        teams.Use(s.authMiddleware())
        {
            teams.GET("", s.listTeams)
            teams.POST("", s.createTeam)
            teams.GET("/:id", s.getTeam)
            teams.PUT("/:id", s.updateTeam)
            teams.DELETE("/:id", s.deleteTeam)
            teams.POST("/:id/invite", s.inviteMember)
            teams.DELETE("/:id/members/:username", s.removeMember)
            teams.PUT("/:id/members/:username/role", s.updateMemberRole)
        }

        // Tenants
        tenants := api.Group("/tenants")
        tenants.Use(s.authMiddleware())
        {
            tenants.GET("", s.listTenants)
            tenants.POST("", s.requireRole("admin"), s.createTenant)
            tenants.GET("/:id", s.getTenant)
            tenants.PUT("/:id", s.requireRole("admin"), s.updateTenant)
            tenants.DELETE("/:id", s.requireRole("admin"), s.deleteTenant)
            tenants.POST("/:id/assign", s.requireRole("admin"), s.assignUserToTenant)
            tenants.GET("/:id/usage", s.getTenantUsage)
        }

        // Auth
        api.POST("/auth/login", s.postAuthLogin)
        api.POST("/auth/logout", s.postAuthLogout)
        api.GET("/auth/me", s.getAuthMe)
        api.POST("/auth/users", s.authMiddleware(), s.requireRole("admin"), s.postCreateUser)
	}
}

// Start starts the API server
func (s *Server) Start() error {
	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 NebulaBox API server starting on port %s", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited")
	return nil
}

// Close closes the server and its dependencies
func (s *Server) Close() error {
	return s.containerd.Close()
}

// healthCheck returns server health status
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "nebulabox-api",
		"version":   "0.1.0-alpha",
	})
}
