// API client for NebulaBox backend

// Import generated types from schema (re-export for use throughout the app)
// Path is relative to web/dashboard/src/lib/api.ts -> generated/types.ts
import type { Container as GeneratedContainer, Image as GeneratedImage, Workspace as GeneratedWorkspace, WorkspaceMember as GeneratedWorkspaceMember, WorkspaceSettings as GeneratedWorkspaceSettings } from '../../../../generated/types'

// Re-export with original names for backward compatibility
export type Container = GeneratedContainer
export type Image = GeneratedImage
export type Workspace = GeneratedWorkspace
export type WorkspaceMember = GeneratedWorkspaceMember
export type WorkspaceSettings = GeneratedWorkspaceSettings

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/api'

export interface SystemStats {
  cpuUsage: number
  memoryUsage: number
  diskUsage: number
  containersRunning: number
  containersTotal: number
}

export interface SystemHistoryPoint {
  ts: number
  cpu: number
  mem: number
  disk: number
  running: number
  total: number
}

export interface RunContainerOptions {
  image: string
  name?: string
  port?: string
  ports?: string[]
  detach?: boolean
  env?: string[]
  volume?: string[]
  network?: string
  service?: string
  workspaceId?: string
}

export interface EnvVar {
  key: string
  value: string
  type: 'string' | 'number' | 'boolean' | 'secret'
}

export interface EnvTemplate {
  name: string
  description: string
  variables: EnvVar[]
}

export interface RegistryCatalog { repositories: string[] }
export interface RegistryTags { name: string; tags: string[] }

export interface ImageVersion {
  tag: string
  digest: string
  createdAt: string
  createdBy: string
  size: number
  metadata?: Record<string, string>
  description?: string
}

export interface VersionSummary {
  repository: string
  totalVersions: number
  latest: string
  tags: string[]
  updatedAt: string
}

export interface RepositoryVersions {
  repository: string
  versions: ImageVersion[]
  count: number
  latest?: string
}

export interface ImageSignature {
  image: string
  tag: string
  digest: string
  signedBy: string
  signedAt: string
  signature: string
  publicKey: string
  algorithm: string
  keyId: string
}

export interface SigningKeyInfo {
  keyId: string
  createdAt: string
  createdBy: string
  publicKey: string
}

export interface NodeInfo {
  id: string
  name: string
  address: string
  port: number
  status: string
  lastSeen: string
  labels?: Record<string, string>
  resources: {
    cpuCores: number
    memoryMB: number
    diskGB: number
    cpuUsed: number
    memoryUsed: number
    diskUsed: number
    containersRunning: number
  }
  containers: string[]
  containerCount: number
  region?: string
  zone?: string
}

export interface DeploymentInstance {
  id: string
  containerId: string
  nodeId: string
  nodeName: string
  status: string
  health: string
  createdAt: string
  restarts: number
  lastRestart?: string
}

export interface HealthCheckConfig {
  type: string
  httpPath?: string
  httpPort?: number
  tcpPort?: number
  command?: string[]
  intervalSeconds?: number
  timeoutSeconds?: number
  retries?: number
  startPeriodSec?: number
}

export interface CreateDeploymentRequest {
  name: string
  image: string
  tag?: string
  replicas: number
  nodeSelector?: Record<string, string>
  strategy?: string
  healthCheck?: HealthCheckConfig
  autoRestart?: boolean
  maxRestarts?: number
  restartPolicy?: string
  serviceName?: string
  networkName?: string
  ports?: string[]
}

export interface DeploymentInfo {
  id: string
  name: string
  image: string
  tag: string
  replicas: number
  status: string
  createdAt: string
  updatedAt: string
  instances: DeploymentInstance[]
  nodeSelector?: Record<string, string>
  strategy: string
  healthCheck?: HealthCheckConfig
  autoRestart: boolean
  maxRestarts?: number
  restartPolicy: string
  serviceName?: string
  networkName?: string
  ports?: string[]
}

export interface ResourceLimits {
  cpuShares?: number
  memoryLimit?: number
  memorySwap?: number
  cpuQuota?: number
  cpuPeriod?: number
  pidsLimit?: number
  iopsRead?: number
  iopsWrite?: number
  throttleRead?: number
  throttleWrite?: number
}

export interface SecurityConfig {
  readOnlyRootFS?: boolean
  privileged?: boolean
  capAdd?: string[]
  capDrop?: string[]
  securityOpt?: string[]
  appArmorProfile?: string
  selinuxLabel?: string
  noNewPrivileges?: boolean
}

export interface RuntimeContainerSpec {
  id: string
  name?: string
  image: string
  command?: string[]
  args?: string[]
  env?: Record<string, string>
  workingDir?: string
  user?: string
  networkMode?: string
  network?: string
  ports?: Record<string, string>
  volumes?: Record<string, string>
  labels?: Record<string, string>
  resources?: ResourceLimits
  security?: SecurityConfig
}

export interface RuntimeContainer {
  id: string
  name: string
  image: string
  status: string
  createdAt: string
  startedAt?: string
  finishedAt?: string
  exitCode?: number
  pid?: number
  spec?: RuntimeContainerSpec
  stats?: ContainerStats
}

export interface ContainerStats {
  timestamp: string
  cpuTotal: number
  cpuPercent: number
  memoryUsage: number
  memoryLimit: number
  memoryPercent: number
  networkRx: number
  networkTx: number
  blockRead: number
  blockWrite: number
  pidsCurrent: number
}

export interface RuntimeImage {
  id: string
  name: string
  tag: string
  digest: string
  size: number
  createdAt: string
  layers: ImageLayer[]
  config?: ImageConfig
}

export interface ImageLayer {
  digest: string
  size: number
  mediaType: string
  urls?: string[]
  compressed?: boolean
}

export interface ImageConfig {
  architecture: string
  os: string
  config?: Record<string, any>
  history?: HistoryEntry[]
}

export interface HistoryEntry {
  created: string
  createdBy: string
  emptyLayer: boolean
}

export interface RuntimeInfo {
  name: string
  version: string
  apiVersion: string
  arch: string
  os: string
  containers: number
  images: number
  cpuCount: number
  memoryTotal: number
  memoryUsed: number
  storageTotal: number
  storageUsed: number
}

export interface ResourcePrediction {
  containerId: string
  predictedAt: string
  duration: string
  predictedCpu: number
  predictedMemory: number
  cpuTrend: string
  memoryTrend: string
  confidence: number
  anomalies: Anomaly[]
  recommendations: Recommendation[]
}

export interface Anomaly {
  type: string
  severity: string
  timestamp: string
  value: number
  message: string
}

export interface Recommendation {
  type: string
  priority: string
  category: string
  title: string
  description: string
  action: string
}

export interface ScalingRecommendation {
  targetId: string
  action: string
  currentReplicas: number
  recommendedReplicas: number
  reason: string
  message: string
  confidence: number
  predictedCpu: number
  predictedMemory: number
}

export interface ScalingPolicy {
  id?: string
  targetId: string
  type: string
  minReplicas: number
  maxReplicas: number
  targetCpu?: number
  targetMemory?: number
  scaleUpThreshold?: number
  scaleDownThreshold?: number
  cooldownPeriod?: string
}

export interface ChatCommand {
  command: string
  timestamp: string
  processed: boolean
  success: boolean
  response: string
}

export interface ChatResponse {
  success: boolean
  message: string
  data?: Record<string, any>
  suggestions?: string[]
}

export interface SharedResources {
  network?: string
  volumes?: string[]
  envVars?: Record<string, string>
  ports?: string[]
  labels?: Record<string, string>
}

export interface ContainerGroup {
  id: string
  name: string
  description?: string
  parentGroupId?: string
  containerIds: string[]
  labels?: Record<string, string>
  sharedResources?: SharedResources
  createdAt: string
  updatedAt: string
}

export interface CreateGroupRequest {
  name: string
  description?: string
  parentGroupId?: string
  sharedResources?: SharedResources
}

export interface UpdateGroupRequest {
  name?: string
  description?: string
  sharedResources?: SharedResources
}

export interface ContainerRelationship {
  parentId: string
  childId: string
  type: string
  properties?: Record<string, any>
  createdAt: string
}

export interface ContainerTree {
  containerId: string
  children: ContainerTree[]
  parentId?: string
  type?: string
}

export interface ContainerElements {
  image?: boolean
  envVars?: boolean
  selectedEnv?: string[]
  ports?: boolean
  selectedPorts?: string[]
  volumes?: boolean
  selectedVolumes?: string[]
  network?: boolean
  service?: boolean
  healthCheck?: boolean
  labels?: boolean
  selectedLabels?: string[]
  command?: boolean
  workingDir?: boolean
  resources?: boolean
}

export interface HealthCheckOverride {
  type?: string
  httpPath?: string
  httpPort?: string
  tcpPort?: string
  command?: string[]
  intervalSeconds?: number
  timeoutSeconds?: number
  retries?: number
  startPeriodSec?: number
}

export interface ContainerOverrides {
  image?: string
  envVars?: Record<string, string>
  ports?: Record<string, string>
  volumes?: Record<string, string>
  network?: string
  service?: string
  healthCheck?: HealthCheckOverride
  labels?: Record<string, string>
  command?: string[]
  workingDir?: string
}

export interface SourceContainer {
  containerId: string
  elements: ContainerElements
  priority?: number
  description?: string
}

export interface CompositionSpec {
  name: string
  description?: string
  sources: SourceContainer[]
  overrides?: ContainerOverrides
  strategy?: string
  createdAt?: string
}

export interface ConflictResolution {
  type: string
  source1: string
  source2: string
  value1: any
  value2: any
  resolution: string
  finalValue: any
  message?: string
}

export interface ComposedContainerSpec {
  name: string
  image: string
  envVars: Record<string, string>
  ports: Record<string, string>
  volumes: Record<string, string>
  network?: string
  service?: string
  healthCheck?: HealthCheckOverride
  labels: Record<string, string>
  command?: string[]
  workingDir?: string
  conflicts?: ConflictResolution[]
  sources: string[]
}

export interface GroupHierarchy {
  group: ContainerGroup
  children: GroupHierarchy[]
  containers: string[]
}

export interface StackHealthCheckConfig {
  type: string
  httpPath?: string
  httpPort?: string
  tcpPort?: string
  command?: string[]
  intervalSeconds?: number
  timeoutSeconds?: number
  retries?: number
  startPeriodSec?: number
}

export interface ResourceConfig {
  cpuShares?: number
  memoryLimit?: number
  memorySwap?: number
  cpuQuota?: number
  cpuPeriod?: number
  pidsLimit?: number
  iopsRead?: number
  iopsWrite?: number
}

export interface ContainerConfig {
  name: string
  image: string
  ports?: Record<string, string>
  envVars?: Record<string, string>
  volumes?: string[]
  network?: string
  service?: string
  dependsOn?: string[]
  healthCheck?: StackHealthCheckConfig
  labels?: Record<string, string>
  command?: string[]
  resources?: ResourceConfig
}

export interface NetworkConfig {
  name: string
  driver?: string
  labels?: Record<string, string>
  options?: Record<string, string>
}

export interface VolumeConfig {
  name: string
  driver?: string
  labels?: Record<string, string>
  options?: Record<string, string>
}

export interface StackTemplate {
  id: string
  name: string
  description: string
  category: string
  containers: ContainerConfig[]
  networks?: NetworkConfig[]
  volumes?: VolumeConfig[]
  envVars?: Record<string, string>
  labels?: Record<string, string>
  tags?: string[]
  createdAt?: string
}

export interface ResourceLimits {
  maxCpu?: number
  maxMemory?: number
  maxDisk?: number
}

// WorkspaceSettings is now imported from generated/types
// Extended WorkspaceSettings for additional fields
export interface WorkspaceSettingsExtended extends WorkspaceSettings {
  allowedIPs?: string[]
  permissions?: Record<string, string>
  resourceLimits?: ResourceLimits
}

// WorkspaceMember and Workspace are now imported from generated/types
// Extended types for backward compatibility

export interface SessionConnection {
  protocol: string
  endpoint: string
  token?: string
  port?: number
}

// Extended Workspace for additional fields not in schema
export interface WorkspaceExtended extends Workspace {
  containerId: string
  members: WorkspaceMember[]
  settings: WorkspaceSettingsExtended
}

export interface Session {
  id: string
  workspaceId: string
  userId: string
  username: string
  type: string
  connection: SessionConnection
  metadata?: Record<string, string>
  startedAt: string
  lastActivity: string
}

export interface SessionState {
  sessionId: string
  userId: string
  username: string
  type: string
  state: string
  lastActivity: string
  metadata?: Record<string, string>
}

export interface Invite {
  id: string
  workspaceId: string
  inviterId: string
  inviterName: string
  email?: string
  role: string
  token: string
  status: string
  expiresAt: string
  createdAt: string
  acceptedAt?: string
  metadata?: Record<string, string>
  link?: string
}

export interface AuditLog {
  id: string
  workspaceId?: string
  userId: string
  username: string
  action: string
  resourceType?: string
  resourceId?: string
  details?: Record<string, string>
  ipAddress?: string
  userAgent?: string
  success: boolean
  errorMessage?: string
  timestamp: string
}

export interface AuditStats {
  totalActions: number
  successfulActions: number
  failedActions: number
  actionsByType: Record<string, number>
  usersByActivity: Record<string, number>
}

export interface SyncEvent {
  id: string
  workspaceId: string
  resourceType: string
  resourceId: string
  changeType: 'create' | 'update' | 'delete'
  data?: any
  timestamp: string
  userId: string
  metadata?: Record<string, string>
}

export interface FileChange {
  id: string
  workspaceId: string
  containerId: string
  path: string
  changeType: 'created' | 'modified' | 'deleted' | 'renamed'
  oldPath?: string
  hash?: string
  size: number
  isDirectory: boolean
  userId: string
  timestamp: string
  metadata?: Record<string, string>
}

export interface Snapshot {
  id: string
  name: string
  description?: string
  type: 'container' | 'workspace' | 'volume'
  resourceId: string
  state: 'creating' | 'ready' | 'failed' | 'restoring'
  size: number
  createdAt: string
  createdBy: string
  metadata?: Record<string, any>
  image?: string
  command?: string[]
  env?: Record<string, string>
  ports?: Record<string, string>
  volumes?: string[]
  network?: string
  resources?: {
    cpu?: string
    memory?: string
    iops?: number
    pids?: number
  }
  filesystemHash?: string
  workspaceSettings?: Record<string, any>
  members?: string[]
}

export interface CRDTOperation {
  id: string
  type: 'orset' | 'lwwreg' | 'counter' | 'map'
  workspaceId: string
  resourceId: string
  resourceType: string
  operation: 'add' | 'remove' | 'update' | 'increment' | 'merge'
  key?: string
  value?: any
  timestamp: string
  userId: string
  vectorClock: Record<string, number>
  metadata?: Record<string, any>
}

export interface Conflict {
  resourceId: string
  resourceType: string
  operations: CRDTOperation[]
  type: 'concurrent' | 'divergent' | 'lost_update'
  resolved: boolean
  resolution?: CRDTOperation
}

export interface EphemeralRuntime {
  id: string
  name: string
  workspaceId: string
  region: string
  status: 'provisioning' | 'active' | 'idle' | 'sleeping' | 'terminating' | 'terminated'
  instanceType: string
  image: string
  createdAt: string
  expiresAt: string
  lastActivityAt: string
  createdBy: string
  members: string[]
  accessUrl: string
  sshKey?: string
  resources?: {
    cpu: number
    memory: string
    disk: string
    network?: string
    publicIp?: string
    privateIp?: string
  }
  metadata?: Record<string, any>
  snapshotId?: string
}

export interface AutoSleepConfig {
  enabled: boolean
  idleTimeout: number // minutes
  sleepTimeout: number // minutes
  createSnapshot: boolean
  autoWakeOnAccess: boolean
}

export interface WorkspaceActivity {
  workspaceId: string
  workspaceName: string
  lastActivity: string
  idleDuration: string
  status: string
}

export interface Vulnerability { id: string; package: string; installed: string; fixedVersion?: string; severity: 'CRITICAL'|'HIGH'|'MEDIUM'|'LOW'|'UNKNOWN'; title: string; description?: string; source: string }
export interface ImageScanResult { image: string; scannedAt: string; criticalCount: number; highCount: number; mediumCount: number; lowCount: number; unknownCount: number; vulnerabilities: Vulnerability[] }
export interface LogEntry { timestamp: number; container: string; level: string; message: string }
export interface AlertThresholds { cpuHigh: number; memoryHigh: number; diskHigh: number }
export interface AlertEvent { type: 'cpu'|'memory'|'disk'; value: number; threshold: number; timestamp: number }

export class ApiClient {
  private authToken: string | null
  private appToken: string | null

  constructor() {
    this.authToken = localStorage.getItem('nebula_registry_token')
    this.appToken = localStorage.getItem('nebula_app_token')
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
      // Import API tracker dynamically to avoid circular dependencies
      const { apiTracker } = await import('./api-test-aware')
      const method = options?.method || 'GET'
      const startTime = Date.now()

      try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
          ...options,
          headers: {
            'Content-Type': 'application/json',
            ...(this.authToken ? { Authorization: `Bearer ${this.authToken}` } : {}),
            ...(this.appToken ? { 'X-App-Auth': this.appToken } : {}),
            ...(options?.headers as Record<string, string>),
          },
        })

        if (!response.ok) {
          const errorText = await response.text()
          const error = errorText ? JSON.parse(errorText) : { message: `HTTP ${response.status}` }
          apiTracker.recordCall(method, endpoint, options?.body, undefined, error.message || error.error || `HTTP ${response.status}`, Date.now() - startTime)
          throw new Error(error.message || error.error || `API Error: ${response.statusText}`)
        }

        const data = await response.json()
        apiTracker.recordCall(method, endpoint, options?.body, data, undefined, Date.now() - startTime)
        return data
      } catch (error: any) {
        apiTracker.recordCall(method, endpoint, options?.body, undefined, error.message || 'Network error', Date.now() - startTime)
        throw error
      }
  }

  // Container operations
  async listContainers(all?: boolean): Promise<Container[]> {
    const url = all ? '/containers?all=true' : '/containers'
    return this.request<Container[]>(url)
  }

  async getContainer(id: string): Promise<Container> {
    return this.request<Container>(`/containers/${id}`)
  }

  async runContainer(options: RunContainerOptions): Promise<Container> {
    return this.request<Container>('/containers/run', {
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  async stopContainer(id: string): Promise<void> {
    return this.request<void>(`/containers/${id}/stop`, {
      method: 'POST',
    })
  }

  async getContainerLogs(id: string): Promise<string[]> {
    return this.request<string[]>(`/containers/${id}/logs`)
  }

  // Image operations
  async listImages(): Promise<Image[]> {
    return this.request<Image[]>('/images')
  }

  async pullImage(image: string): Promise<void> {
    return this.request<void>(`/images/pull`, {
      method: 'POST',
      body: JSON.stringify({ image }),
    })
  }

  async pushImage(image: string): Promise<void> {
    return this.request<void>(`/images/push`, {
      method: 'POST',
      body: JSON.stringify({ image }),
    })
  }

  async buildImage(path: string, tag?: string): Promise<void> {
    return this.request<void>('/images/build', {
      method: 'POST',
      body: JSON.stringify({ path, tag }),
    })
  }

  async buildImageFromDockerfile(dockerfile: string, tag: string): Promise<{ tag: string; image: string; logs: string[] }> {
    return this.request<{ tag: string; image: string; logs: string[] }>('/images/build', {
      method: 'POST',
      body: JSON.stringify({ dockerfile, tag })
    })
  }

  // Build Specification APIs
  async validateBuildSpec(spec: Record<string, any>): Promise<{ valid: boolean; dockerfile?: string; errors?: string[]; message?: string }> {
    return this.request<{ valid: boolean; dockerfile?: string; errors?: string[]; message?: string }>('/buildspec/validate', {
      method: 'POST',
      body: JSON.stringify({ spec }),
    })
  }

  async convertBuildSpec(spec: Record<string, any>): Promise<{ valid: boolean; dockerfile: string; message?: string }> {
    return this.request<{ valid: boolean; dockerfile: string; message?: string }>('/buildspec/convert', {
      method: 'POST',
      body: JSON.stringify({ spec }),
    })
  }

  async buildFromSpec(spec: Record<string, any>, tag: string): Promise<{ valid: boolean; dockerfile?: string; tag: string; logs: string[]; message?: string }> {
    return this.request<{ valid: boolean; dockerfile?: string; tag: string; logs: string[]; message?: string }>('/buildspec/build', {
      method: 'POST',
      body: JSON.stringify({ spec, tag }),
    })
  }

  // Security APIs (signing, scanning)
  async signImage(image: string, tag: string, digest: string, keyId: string): Promise<ImageSignature> {
    return this.request<ImageSignature>('/security/sign', {
      method: 'POST',
      body: JSON.stringify({ image, tag, digest, keyId }),
    })
  }

  async verifySignature(image: string, tag: string, digest: string, signature: ImageSignature): Promise<{ valid: boolean; image: string; tag: string; digest: string }> {
    return this.request<{ valid: boolean; image: string; tag: string; digest: string }>('/security/verify', {
      method: 'POST',
      body: JSON.stringify({ image, tag, digest, signature }),
    })
  }

  async generateSigningKey(keyId: string): Promise<SigningKeyInfo> {
    return this.request<SigningKeyInfo>('/security/keys/generate', {
      method: 'POST',
      body: JSON.stringify({ keyId }),
    })
  }

  async listSigningKeys(): Promise<{ keys: SigningKeyInfo[]; count: number }> {
    return this.request<{ keys: SigningKeyInfo[]; count: number }>('/security/keys')
  }

  async getSigningKey(keyId: string): Promise<SigningKeyInfo> {
    return this.request<SigningKeyInfo>(`/security/keys/${encodeURIComponent(keyId)}`)
  }

  // Orchestrator APIs (multi-node deployment)
  async registerNode(node: { id: string; name?: string; address: string; port?: number; labels?: Record<string, string>; region?: string; zone?: string }): Promise<NodeInfo> {
    return this.request<NodeInfo>('/orchestrator/nodes', {
      method: 'POST',
      body: JSON.stringify(node),
    })
  }

  async listNodes(onlineOnly?: boolean): Promise<{ nodes: NodeInfo[]; count: number }> {
    const url = onlineOnly ? '/orchestrator/nodes?online=true' : '/orchestrator/nodes'
    return this.request<{ nodes: NodeInfo[]; count: number }>(url)
  }

  async getNode(nodeId: string): Promise<NodeInfo> {
    return this.request<NodeInfo>(`/orchestrator/nodes/${encodeURIComponent(nodeId)}`)
  }

  async updateNodeStatus(nodeId: string, status: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/orchestrator/nodes/${encodeURIComponent(nodeId)}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })
  }

  async createDeployment(deploy: CreateDeploymentRequest): Promise<DeploymentInfo> {
    return this.request<DeploymentInfo>('/orchestrator/deployments', {
      method: 'POST',
      body: JSON.stringify(deploy),
    })
  }

  async listDeployments(): Promise<{ deployments: DeploymentInfo[]; count: number }> {
    return this.request<{ deployments: DeploymentInfo[]; count: number }>('/orchestrator/deployments')
  }

  async getDeployment(deploymentId: string): Promise<DeploymentInfo> {
    return this.request<DeploymentInfo>(`/orchestrator/deployments/${encodeURIComponent(deploymentId)}`)
  }

  async updateDeployment(deploymentId: string, updates: Partial<CreateDeploymentRequest>): Promise<DeploymentInfo> {
    return this.request<DeploymentInfo>(`/orchestrator/deployments/${encodeURIComponent(deploymentId)}`, {
      method: 'PATCH',
      body: JSON.stringify(updates),
    })
  }

  async deleteDeployment(deploymentId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/orchestrator/deployments/${encodeURIComponent(deploymentId)}`, {
      method: 'DELETE',
    })
  }

  async scaleDeployment(deploymentId: string, replicas: number): Promise<{ message: string; replicas: number }> {
    return this.request<{ message: string; replicas: number }>(`/orchestrator/deployments/${encodeURIComponent(deploymentId)}/scale`, {
      method: 'POST',
      body: JSON.stringify({ replicas }),
    })
  }

  async restartDeployment(deploymentId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/orchestrator/deployments/${encodeURIComponent(deploymentId)}/restart`, {
      method: 'POST',
    })
  }

  // Nebula Runtime APIs (custom container runtime)
  async createRuntimeContainer(container: RuntimeContainerSpec): Promise<RuntimeContainer> {
    return this.request<RuntimeContainer>('/runtime/containers', {
      method: 'POST',
      body: JSON.stringify(container),
    })
  }

  async listRuntimeContainers(): Promise<{ containers: RuntimeContainer[]; count: number }> {
    return this.request<{ containers: RuntimeContainer[]; count: number }>('/runtime/containers')
  }

  async getRuntimeContainer(containerId: string): Promise<RuntimeContainer> {
    return this.request<RuntimeContainer>(`/runtime/containers/${encodeURIComponent(containerId)}`)
  }

  async startRuntimeContainer(containerId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/runtime/containers/${encodeURIComponent(containerId)}/start`, {
      method: 'POST',
    })
  }

  async stopRuntimeContainer(containerId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/runtime/containers/${encodeURIComponent(containerId)}/stop`, {
      method: 'POST',
    })
  }

  async deleteRuntimeContainer(containerId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/runtime/containers/${encodeURIComponent(containerId)}`, {
      method: 'DELETE',
    })
  }

  async pullRuntimeImage(image: string): Promise<{ message: string }> {
    return this.request<{ message: string }>('/runtime/images/pull', {
      method: 'POST',
      body: JSON.stringify({ image }),
    })
  }

  async listRuntimeImages(): Promise<{ images: RuntimeImage[]; count: number }> {
    return this.request<{ images: RuntimeImage[]; count: number }>('/runtime/images')
  }

  async getRuntimeInfo(): Promise<RuntimeInfo> {
    return this.request<RuntimeInfo>('/runtime/info')
  }

  async getRuntimeVersion(): Promise<{ version: string }> {
    return this.request<{ version: string }>('/runtime/version')
  }

  // AI Ops APIs (predictive analytics and auto-scaling)
  async recordMetric(containerId: string, cpu: number, memory: number, networkRx?: number, networkTx?: number): Promise<{ message: string }> {
    return this.request<{ message: string }>('/aiops/metrics', {
      method: 'POST',
      body: JSON.stringify({ containerId, cpu, memory, networkRx, networkTx }),
    })
  }

  async predictResourceUsage(containerId: string, duration?: string): Promise<ResourcePrediction> {
    const params = duration ? `?duration=${encodeURIComponent(duration)}` : ''
    return this.request<ResourcePrediction>(`/aiops/predict/${encodeURIComponent(containerId)}${params}`)
  }

  async getScalingRecommendation(targetId: string, replicas?: number): Promise<ScalingRecommendation> {
    const params = replicas ? `?replicas=${replicas}` : ''
    return this.request<ScalingRecommendation>(`/aiops/scaling/${encodeURIComponent(targetId)}${params}`)
  }

  async setScalingPolicy(policy: ScalingPolicy): Promise<{ message: string }> {
    return this.request<{ message: string }>('/aiops/scaling/policy', {
      method: 'POST',
      body: JSON.stringify(policy),
    })
  }

  async processChatCommand(command: string): Promise<ChatResponse> {
    return this.request<ChatResponse>('/aiops/chat', {
      method: 'POST',
      body: JSON.stringify({ command }),
    })
  }

  async getChatHistory(limit?: number): Promise<{ commands: ChatCommand[]; count: number }> {
    const params = limit ? `?limit=${limit}` : ''
    return this.request<{ commands: ChatCommand[]; count: number }>(`/aiops/chat/history${params}`)
  }

  // Container Groups and Hierarchy APIs
  async createGroup(group: CreateGroupRequest): Promise<ContainerGroup> {
    return this.request<ContainerGroup>('/groups', {
      method: 'POST',
      body: JSON.stringify(group),
    })
  }

  async listGroups(parentId?: string): Promise<{ groups: ContainerGroup[]; count: number }> {
    const params = parentId ? `?parentId=${encodeURIComponent(parentId)}` : ''
    return this.request<{ groups: ContainerGroup[]; count: number }>(`/groups${params}`)
  }

  async getGroup(groupId: string): Promise<ContainerGroup> {
    return this.request<ContainerGroup>(`/groups/${encodeURIComponent(groupId)}`)
  }

  async updateGroup(groupId: string, group: UpdateGroupRequest): Promise<ContainerGroup> {
    return this.request<ContainerGroup>(`/groups/${encodeURIComponent(groupId)}`, {
      method: 'PATCH',
      body: JSON.stringify(group),
    })
  }

  async deleteGroup(groupId: string, force?: boolean): Promise<{ message: string }> {
    const params = force ? '?force=true' : ''
    return this.request<{ message: string }>(`/groups/${encodeURIComponent(groupId)}${params}`, {
      method: 'DELETE',
    })
  }

  async addContainerToGroup(groupId: string, containerId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/groups/${encodeURIComponent(groupId)}/containers`, {
      method: 'POST',
      body: JSON.stringify({ containerId }),
    })
  }

  async removeContainerFromGroup(groupId: string, containerId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/groups/${encodeURIComponent(groupId)}/containers/${encodeURIComponent(containerId)}`, {
      method: 'DELETE',
    })
  }

  async getGroupHierarchy(groupId: string): Promise<GroupHierarchy> {
    return this.request<GroupHierarchy>(`/groups/${encodeURIComponent(groupId)}/hierarchy`)
  }

  async startGroup(groupId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/groups/${encodeURIComponent(groupId)}/start`, {
      method: 'POST',
    })
  }

  async stopGroup(groupId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/groups/${encodeURIComponent(groupId)}/stop`, {
      method: 'POST',
    })
  }

  async createRelationship(parentId: string, childId: string, type: string, properties?: Record<string, any>): Promise<ContainerRelationship> {
    return this.request<ContainerRelationship>('/groups/relationships', {
      method: 'POST',
      body: JSON.stringify({ parentId, childId, type, properties }),
    })
  }

  async deleteRelationship(parentId: string, childId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>('/groups/relationships', {
      method: 'DELETE',
      body: JSON.stringify({ parentId, childId }),
    })
  }

  async getContainerChildren(containerId: string): Promise<{ children: ContainerRelationship[]; count: number }> {
    return this.request<{ children: ContainerRelationship[]; count: number }>(`/groups/relationships/${encodeURIComponent(containerId)}/children`)
  }

  async getContainerParent(containerId: string): Promise<ContainerRelationship> {
    return this.request<ContainerRelationship>(`/groups/relationships/${encodeURIComponent(containerId)}/parent`)
  }

  async getContainerTree(containerId: string): Promise<ContainerTree> {
    return this.request<ContainerTree>(`/groups/relationships/${encodeURIComponent(containerId)}/tree`)
  }

  async getContainerAncestry(containerId: string): Promise<{ ancestry: ContainerRelationship[]; count: number }> {
    return this.request<{ ancestry: ContainerRelationship[]; count: number }>(`/groups/relationships/${encodeURIComponent(containerId)}/ancestry`)
  }

  // Container Composition APIs
  async createCompositionSpec(spec: CompositionSpec): Promise<CompositionSpec> {
    return this.request<CompositionSpec>('/composition/specs', {
      method: 'POST',
      body: JSON.stringify(spec),
    })
  }

  async listCompositionSpecs(): Promise<{ specs: CompositionSpec[]; count: number }> {
    return this.request<{ specs: CompositionSpec[]; count: number }>('/composition/specs')
  }

  async getCompositionSpec(name: string): Promise<CompositionSpec> {
    return this.request<CompositionSpec>(`/composition/specs/${encodeURIComponent(name)}`)
  }

  async deleteCompositionSpec(name: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/composition/specs/${encodeURIComponent(name)}`, {
      method: 'DELETE',
    })
  }

  async previewComposition(spec: CompositionSpec): Promise<ComposedContainerSpec> {
    return this.request<ComposedContainerSpec>('/composition/preview', {
      method: 'POST',
      body: JSON.stringify(spec),
    })
  }

  async composeContainerFromSpec(specName: string, containerName?: string, start?: boolean): Promise<{ container: Container; composition: ComposedContainerSpec }> {
    return this.request<{ container: Container; composition: ComposedContainerSpec }>('/composition/compose', {
      method: 'POST',
      body: JSON.stringify({ specName, containerName, start }),
    })
  }

  // Stack Templates APIs
  async listTemplates(category?: string, tag?: string): Promise<{ templates: StackTemplate[]; count: number }> {
    const params = new URLSearchParams()
    if (category) params.append('category', category)
    if (tag) params.append('tag', tag)
    const query = params.toString()
    return this.request<{ templates: StackTemplate[]; count: number }>(`/templates${query ? `?${query}` : ''}`)
  }

  async getTemplate(id: string): Promise<StackTemplate> {
    return this.request<StackTemplate>(`/templates/${encodeURIComponent(id)}`)
  }

  async saveTemplate(template: StackTemplate): Promise<StackTemplate> {
    return this.request<StackTemplate>('/templates', {
      method: 'POST',
      body: JSON.stringify(template),
    })
  }

  async deleteTemplate(id: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/templates/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    })
  }

  async deployTemplate(id: string, prefix?: string, envVars?: Record<string, string>, start?: boolean): Promise<any> {
    return this.request<any>(`/templates/${encodeURIComponent(id)}/deploy`, {
      method: 'POST',
      body: JSON.stringify({ prefix, envVars, start }),
    })
  }

  // Shared Runtime APIs
  async createWorkspace(workspace: { name: string; description?: string; containerId: string; settings?: WorkspaceSettings }): Promise<Workspace> {
    return this.request<Workspace>('/shareruntime/workspaces', {
      method: 'POST',
      body: JSON.stringify(workspace),
    })
  }

  async listWorkspaces(): Promise<{ workspaces: Workspace[]; count: number }> {
    return this.request<{ workspaces: Workspace[]; count: number }>('/shareruntime/workspaces')
  }

  async getWorkspace(id: string): Promise<Workspace> {
    return this.request<Workspace>(`/shareruntime/workspaces/${encodeURIComponent(id)}`)
  }

  async deleteWorkspace(id: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    })
  }

  async updateWorkspaceStatus(id: string, status: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(id)}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })
  }

  async addWorkspaceMember(workspaceId: string, member: { userId: string; username: string; role: string }): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/members`, {
      method: 'POST',
      body: JSON.stringify(member),
    })
  }

  async removeWorkspaceMember(workspaceId: string, userId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/members/${encodeURIComponent(userId)}`, {
      method: 'DELETE',
    })
  }

  async updateWorkspaceMemberRole(workspaceId: string, userId: string, role: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/members/${encodeURIComponent(userId)}/role`, {
      method: 'PATCH',
      body: JSON.stringify({ role }),
    })
  }

  async createInvite(workspaceId: string, invite: { email?: string; role: string; expiresInHours?: number }): Promise<Invite> {
    return this.request<Invite>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/invites`, {
      method: 'POST',
      body: JSON.stringify(invite),
    })
  }

  async listInvites(workspaceId: string): Promise<{ invites: Invite[]; count: number }> {
    return this.request<{ invites: Invite[]; count: number }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/invites`)
  }

  async acceptInvite(token: string): Promise<{ invite: Invite; message: string }> {
    return this.request<{ invite: Invite; message: string }>(`/shareruntime/invites/${encodeURIComponent(token)}/accept`, {
      method: 'POST',
    })
  }

  async getInviteLink(workspaceId: string, token: string): Promise<{ link: string }> {
    return this.request<{ link: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/invites/${encodeURIComponent(token)}/link`)
  }

  async getInviteInfo(token: string): Promise<{ workspaceId: string; role: string; inviterName: string; expiresAt: string; status: string }> {
    return this.request<{ workspaceId: string; role: string; inviterName: string; expiresAt: string; status: string }>(`/shareruntime/invites/${encodeURIComponent(token)}/info`)
  }

  async getUserPermissions(workspaceId: string): Promise<{ role: string; permissions: Record<string, boolean> }> {
    return this.request<{ role: string; permissions: Record<string, boolean> }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/permissions`)
  }

  async getAuditLogs(workspaceId: string, filters?: {
    action?: string
    userId?: string
    resourceType?: string
    resourceId?: string
    success?: boolean
    startTime?: string
    endTime?: string
    limit?: number
  }): Promise<{ logs: AuditLog[]; count: number }> {
    const params = new URLSearchParams()
    if (filters) {
      if (filters.action) params.append('action', filters.action)
      if (filters.userId) params.append('userId', filters.userId)
      if (filters.resourceType) params.append('resourceType', filters.resourceType)
      if (filters.resourceId) params.append('resourceId', filters.resourceId)
      if (filters.success !== undefined) params.append('success', String(filters.success))
      if (filters.startTime) params.append('startTime', filters.startTime)
      if (filters.endTime) params.append('endTime', filters.endTime)
      if (filters.limit) params.append('limit', String(filters.limit))
    }
    const query = params.toString()
    return this.request<{ logs: AuditLog[]; count: number }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/audit-logs${query ? '?' + query : ''}`)
  }

  async getUserAuditLogs(limit?: number): Promise<{ logs: AuditLog[]; count: number }> {
    const params = limit ? `?limit=${limit}` : ''
    return this.request<{ logs: AuditLog[]; count: number }>(`/shareruntime/audit-logs${params}`)
  }

  async getAuditStats(workspaceId: string, startTime?: string, endTime?: string): Promise<AuditStats> {
    const params = new URLSearchParams()
    if (startTime) params.append('startTime', startTime)
    if (endTime) params.append('endTime', endTime)
    const query = params.toString()
    return this.request<AuditStats>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/audit-stats${query ? '?' + query : ''}`)
  }

  // NebulaSync APIs
  async getChangesSince(workspaceId: string, since?: string): Promise<{ changes: SyncEvent[]; count: number; since: string }> {
    const params = since ? `?since=${encodeURIComponent(since)}` : ''
    return this.request<{ changes: SyncEvent[]; count: number; since: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/sync/changes${params}`)
  }

  async getLatestChangeID(workspaceId: string): Promise<{ changeId: string }> {
    return this.request<{ changeId: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/sync/latest`)
  }

  async syncWorkspace(workspaceId: string, changes: SyncEvent[]): Promise<{ applied: number; total: number; errors: string[] }> {
    return this.request<{ applied: number; total: number; errors: string[] }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/sync/apply`, {
      method: 'POST',
      body: JSON.stringify({ changes }),
    })
  }

  // FileSync APIs
  async recordFileChange(workspaceId: string, change: {
    containerId: string
    path: string
    changeType: 'created' | 'modified' | 'deleted' | 'renamed'
    oldPath?: string
    isDirectory?: boolean
    size?: number
    metadata?: Record<string, string>
  }): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/filesync/changes`, {
      method: 'POST',
      body: JSON.stringify(change),
    })
  }

  async getFileChanges(workspaceId: string, containerId?: string, since?: string): Promise<{ changes: FileChange[]; count: number; since: string }> {
    const params = new URLSearchParams()
    if (containerId) params.append('containerId', containerId)
    if (since) params.append('since', since)
    const query = params.toString()
    return this.request<{ changes: FileChange[]; count: number; since: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/filesync/changes${query ? '?' + query : ''}`)
  }

  async syncFile(workspaceId: string, sync: {
    containerId: string
    filePath: string
    targetPath: string
  }): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/filesync/sync`, {
      method: 'POST',
      body: JSON.stringify(sync),
    })
  }

  async getFileHash(filePath: string): Promise<{ path: string; hash: string }> {
    return this.request<{ path: string; hash: string }>(`/shareruntime/filesync/hash?path=${encodeURIComponent(filePath)}`)
  }

  // Snapshot APIs
  async createSnapshot(snapshot: {
    name: string
    description?: string
    type: 'container' | 'workspace' | 'volume'
    resourceId: string
    metadata?: Record<string, any>
  }): Promise<Snapshot> {
    return this.request<Snapshot>(`/snapshots`, {
      method: 'POST',
      body: JSON.stringify(snapshot),
    })
  }

  async listSnapshots(resourceId?: string, type?: string): Promise<{ snapshots: Snapshot[]; count: number }> {
    const params = new URLSearchParams()
    if (resourceId) params.append('resourceId', resourceId)
    if (type) params.append('type', type)
    const query = params.toString()
    return this.request<{ snapshots: Snapshot[]; count: number }>(`/snapshots${query ? '?' + query : ''}`)
  }

  async getSnapshot(snapshotId: string): Promise<Snapshot> {
    return this.request<Snapshot>(`/snapshots/${encodeURIComponent(snapshotId)}`)
  }

  async deleteSnapshot(snapshotId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/snapshots/${encodeURIComponent(snapshotId)}`, {
      method: 'DELETE',
    })
  }

  async restoreSnapshot(snapshotId: string, newName?: string): Promise<{ message: string; container?: any }> {
    return this.request<{ message: string; container?: any }>(`/snapshots/${encodeURIComponent(snapshotId)}/restore`, {
      method: 'POST',
      body: JSON.stringify({ newName }),
    })
  }

  async listResourceSnapshots(resourceId: string): Promise<{ snapshots: Snapshot[]; count: number }> {
    return this.request<{ snapshots: Snapshot[]; count: number }>(`/snapshots/resource/${encodeURIComponent(resourceId)}`)
  }

  // CRDT/Conflict Resolution APIs
  async recordCRDTOperation(workspaceId: string, operation: {
    type: 'orset' | 'lwwreg' | 'counter' | 'map'
    resourceId: string
    resourceType: string
    operation: 'add' | 'remove' | 'update' | 'increment'
    key?: string
    value?: any
    vectorClock?: Record<string, number>
    metadata?: Record<string, any>
  }): Promise<CRDTOperation> {
    return this.request<CRDTOperation>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/crdt/operations`, {
      method: 'POST',
      body: JSON.stringify(operation),
    })
  }

  async getCRDTOperations(workspaceId: string, since?: string): Promise<{ operations: CRDTOperation[]; count: number; since: string }> {
    const params = new URLSearchParams()
    if (since) params.append('since', since)
    const query = params.toString()
    return this.request<{ operations: CRDTOperation[]; count: number; since: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/crdt/operations${query ? '?' + query : ''}`)
  }

  async detectConflicts(workspaceId: string, since?: string): Promise<{ conflicts: Conflict[]; count: number }> {
    const params = new URLSearchParams()
    if (since) params.append('since', since)
    const query = params.toString()
    return this.request<{ conflicts: Conflict[]; count: number }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/crdt/conflicts/detect${query ? '?' + query : ''}`, {
      method: 'POST',
    })
  }

  async resolveConflict(workspaceId: string, conflictId: string): Promise<{ conflict: Conflict; resolution: CRDTOperation; message: string }> {
    return this.request<{ conflict: Conflict; resolution: CRDTOperation; message: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/crdt/conflicts/${encodeURIComponent(conflictId)}/resolve`, {
      method: 'POST',
    })
  }

  async getResourceState(workspaceId: string, resourceId: string, resourceType?: string): Promise<{ resourceId: string; resourceType: string; operations: CRDTOperation[]; count: number }> {
    const params = new URLSearchParams()
    if (resourceType) params.append('type', resourceType)
    const query = params.toString()
    return this.request<{ resourceId: string; resourceType: string; operations: CRDTOperation[]; count: number }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/crdt/resources/${encodeURIComponent(resourceId)}${query ? '?' + query : ''}`)
  }

  // Ephemeral Runtime APIs
  async provisionEphemeralRuntime(workspaceId: string, runtime: {
    name: string
    region: string
    instanceType: 'small' | 'medium' | 'large'
    image: string
    duration?: number
    members?: string[]
  }): Promise<EphemeralRuntime> {
    const params = new URLSearchParams()
    params.append('workspaceId', workspaceId)
    return this.request<EphemeralRuntime>(`/cloud/ephemeral/runtimes?${params.toString()}`, {
      method: 'POST',
      body: JSON.stringify(runtime),
    })
  }

  async listEphemeralRuntimes(workspaceId?: string): Promise<{ runtimes: EphemeralRuntime[]; count: number }> {
    const params = new URLSearchParams()
    if (workspaceId) params.append('workspaceId', workspaceId)
    const query = params.toString()
    return this.request<{ runtimes: EphemeralRuntime[]; count: number }>(`/cloud/ephemeral/runtimes${query ? '?' + query : ''}`)
  }

  async getEphemeralRuntime(runtimeId: string): Promise<EphemeralRuntime> {
    return this.request<EphemeralRuntime>(`/cloud/ephemeral/runtimes/${encodeURIComponent(runtimeId)}`)
  }

  async terminateEphemeralRuntime(runtimeId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/cloud/ephemeral/runtimes/${encodeURIComponent(runtimeId)}`, {
      method: 'DELETE',
    })
  }

  async updateEphemeralRuntimeActivity(runtimeId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/cloud/ephemeral/runtimes/${encodeURIComponent(runtimeId)}/activity`, {
      method: 'POST',
    })
  }

  async sleepEphemeralRuntime(runtimeId: string, snapshotId?: string): Promise<{ message: string; snapshotId: string }> {
    return this.request<{ message: string; snapshotId: string }>(`/cloud/ephemeral/runtimes/${encodeURIComponent(runtimeId)}/sleep`, {
      method: 'POST',
      body: JSON.stringify({ snapshotId }),
    })
  }

  async wakeEphemeralRuntime(runtimeId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/cloud/ephemeral/runtimes/${encodeURIComponent(runtimeId)}/wake`, {
      method: 'POST',
    })
  }

  async checkEphemeralRuntimeHealth(runtimeId: string): Promise<{
    runtimeId: string
    status: string
    isExpired: boolean
    idleDuration: string
    expiresAt: string
    lastActivity: string
  }> {
    return this.request<{
      runtimeId: string
      status: string
      isExpired: boolean
      idleDuration: string
      expiresAt: string
      lastActivity: string
    }>(`/cloud/ephemeral/runtimes/${encodeURIComponent(runtimeId)}/health`)
  }

  async addEphemeralRuntimeMember(runtimeId: string, userId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/cloud/ephemeral/runtimes/${encodeURIComponent(runtimeId)}/members`, {
      method: 'POST',
      body: JSON.stringify({ userId }),
    })
  }

  async removeEphemeralRuntimeMember(runtimeId: string, userId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/cloud/ephemeral/runtimes/${encodeURIComponent(runtimeId)}/members/${encodeURIComponent(userId)}`, {
      method: 'DELETE',
    })
  }

  // Auto-Sleep APIs
  async setAutoSleepConfig(workspaceId: string, config: {
    enabled: boolean
    idleTimeout: number // minutes
    sleepTimeout: number // minutes
    createSnapshot: boolean
    autoWakeOnAccess: boolean
  }): Promise<{ message: string; config: any }> {
    return this.request<{ message: string; config: any }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/autosleep/config`, {
      method: 'PUT',
      body: JSON.stringify(config),
    })
  }

  async getAutoSleepConfig(workspaceId: string): Promise<{ config: AutoSleepConfig }> {
    return this.request<{ config: AutoSleepConfig }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/autosleep/config`)
  }

  async recordWorkspaceActivity(workspaceId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/activity`, {
      method: 'POST',
    })
  }

  async wakeWorkspace(workspaceId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/wake`, {
      method: 'POST',
    })
  }

  async getIdleWorkspaces(): Promise<{ workspaces: WorkspaceActivity[]; count: number }> {
    return this.request<{ workspaces: WorkspaceActivity[]; count: number }>(`/shareruntime/autosleep/idle`)
  }

  async createSession(workspaceId: string, session: { type: string; connection?: SessionConnection }): Promise<Session> {
    return this.request<Session>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/sessions`, {
      method: 'POST',
      body: JSON.stringify(session),
    })
  }

  async listWorkspaceSessions(workspaceId: string): Promise<{ sessions: Session[]; count: number }> {
    return this.request<{ sessions: Session[]; count: number }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/sessions`)
  }

  async closeSession(sessionId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/sessions/${encodeURIComponent(sessionId)}`, {
      method: 'DELETE',
    })
  }

  async updateSessionActivity(sessionId: string, metadata?: Record<string, string>): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/sessions/${encodeURIComponent(sessionId)}/activity`, {
      method: 'POST',
      body: JSON.stringify({ metadata }),
    })
  }

  async updateSessionState(sessionId: string, state: string, metadata?: Record<string, string>): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/shareruntime/sessions/${encodeURIComponent(sessionId)}/state`, {
      method: 'PATCH',
      body: JSON.stringify({ state, metadata }),
    })
  }

  async getSessionState(sessionId: string): Promise<SessionState> {
    return this.request<SessionState>(`/shareruntime/sessions/${encodeURIComponent(sessionId)}/state`)
  }

  async listActiveSessions(workspaceId: string): Promise<{ sessions: SessionState[]; count: number }> {
    return this.request<{ sessions: SessionState[]; count: number }>(`/shareruntime/workspaces/${encodeURIComponent(workspaceId)}/active-sessions`)
  }

  // System operations
  async getSystemStats(): Promise<SystemStats> {
    return this.request<SystemStats>('/system/stats')
  }

  // Networks
  async listNetworks(): Promise<Array<{ id:string; name:string; driver:string; subnet:string; created:string }>> {
    return this.request('/networks')
  }

  // Services
  async listServices(): Promise<{ services: Record<string, Array<{ id:string; name:string; address:string; port:number; version?:string; network?:string; createdAt:number }>> }> {
    return this.request('/services')
  }
  async resolveService(name: string): Promise<{ name: string; instances: Array<{ id:string; name:string; address:string; port:number; version?:string; network?:string; createdAt:number }> }> {
    return this.request(`/services/resolve/${encodeURIComponent(name)}`)
  }
  async resolveServiceNext(name: string): Promise<{ name: string; instance: { id:string; name:string; address:string; port:number; version?:string; network?:string; createdAt:number } | null }> {
    return this.request(`/services/next/${encodeURIComponent(name)}`)
  }
  async registerService(payload: { name: string; id?: string; address?: string; port?: number; version?: string; network?: string }): Promise<any> {
    return this.request('/services/register', { method:'POST', body: JSON.stringify(payload) })
  }
  async deregisterService(payload: { name: string; id: string }): Promise<any> {
    return this.request('/services/deregister', { method:'POST', body: JSON.stringify(payload) })
  }

  // DNS
  async listDNSRecords(): Promise<{ records: Record<string, string[]> }> {
    return this.request('/dns/records')
  }
  async addDNSRecord(name: string, a: string[]): Promise<{ name: string; a: string[] }> {
    return this.request('/dns/records', { method:'POST', body: JSON.stringify({ name, a }) })
  }
  async deleteDNSRecord(name: string): Promise<{ deleted: string }> {
    return this.request(`/dns/records/${encodeURIComponent(name)}`, { method:'DELETE' })
  }
  async dnsResolve(name: string): Promise<{ name: string; a: string[] }> {
    return this.request(`/dns/resolve/${encodeURIComponent(name)}`)
  }

  // Ports
  async listPorts(): Promise<{ ports: Array<{ port: number; id: string }> }> {
    return this.request('/ports')
  }
  async reservePort(port: number, id: string): Promise<{ port: number; id: string }> {
    return this.request('/ports/reserve', { method: 'POST', body: JSON.stringify({ port, id }) })
  }
  async releasePort(port: number, id?: string): Promise<{ released: number }> {
    return this.request('/ports/release', { method: 'POST', body: JSON.stringify({ port, id }) })
  }
  async suggestPort(from?: number): Promise<{ port: number }> {
    const qs = from ? `?from=${from}` : ''
    return this.request(`/ports/suggest${qs}`)
  }

  // Webhooks
  async getGitHubEvents(): Promise<{ events: Array<{ id:string; event:string; repo:string; ref:string; action?:string; sender?:string; timestamp:number }> }> {
    return this.request('/webhooks/github/events')
  }
  async getGitLabEvents(): Promise<{ events: Array<{ id:string; event:string; project:string; ref:string; user:string; timestamp:number }> }> {
    return this.request('/webhooks/gitlab/events')
  }

  // Builds
  async getBuilds(): Promise<{ builds: Array<{ id:string; source:string; repo:string; ref:string; status:string; startedAt:number; endedAt:number }> }> {
    return this.request('/builds')
  }
  async triggerBuild(payload: { source?: string; repo: string; ref?: string }): Promise<{ id:string }> {
    return this.request('/builds/trigger', { method:'POST', body: JSON.stringify(payload) })
  }

  // Tests
  async getTests(): Promise<{ tests: Array<{ id:string; repo:string; ref:string; suite:string; status:string; startedAt:number; endedAt:number }> }> {
    return this.request('/tests')
  }
  async runTests(payload: { repo: string; ref?: string; suite?: string }): Promise<{ id:string }> {
    return this.request('/tests/run', { method:'POST', body: JSON.stringify(payload) })
  }

  // Deployments
  async getDeployments(): Promise<{ deployments: Array<{ id:string; repo:string; ref:string; env:string; status:string; startedAt:number; endedAt:number }> }> {
    return this.request('/deployments')
  }
  async triggerDeployment(payload: { repo: string; ref?: string; env?: string }): Promise<{ id:string }> {
    return this.request('/deployments/trigger', { method:'POST', body: JSON.stringify(payload) })
  }

  // Rollbacks
  async getRollbacks(): Promise<{ rollbacks: Array<{ id:string; repo:string; fromRef:string; toRef:string; env:string; status:string; startedAt:number; endedAt:number }> }> {
    return this.request('/deployments/rollbacks')
  }
  async triggerRollback(payload: { repo: string; env: string }): Promise<{ id:string }> {
    return this.request('/deployments/rollback', { method:'POST', body: JSON.stringify(payload) })
  }
  async createNetwork(payload: { name: string; driver?: string; subnet?: string }): Promise<{ id:string; name:string; driver:string; subnet:string; created:string }> {
    return this.request('/networks', { method:'POST', body: JSON.stringify(payload) })
  }
  async deleteNetwork(id: string): Promise<{ deleted: string }> {
    return this.request(`/networks/${id}`, { method:'DELETE' })
  }

  async getSystemHistory(params: { range?: number; step?: number } = {}): Promise<{ points: SystemHistoryPoint[] }> {
    const qs = new URLSearchParams()
    if (params.range) qs.set('range', String(params.range))
    if (params.step) qs.set('step', String(params.step))
    const q = qs.toString()
    const path = q ? `/system/history?${q}` : '/system/history'
    return this.request<{ points: SystemHistoryPoint[] }>(path)
  }

  // Environment variables
  async getContainerEnvVars(containerId: string): Promise<{ success: boolean; message: string; variables: EnvVar[]; error?: string }> {
    return this.request<{ success: boolean; message: string; variables: EnvVar[]; error?: string }>(`/containers/${containerId}/env`)
  }

  async setContainerEnvVars(containerId: string, variables: EnvVar[]): Promise<{ success: boolean; message: string; variables: EnvVar[]; error?: string }> {
    return this.request<{ success: boolean; message: string; variables: EnvVar[]; error?: string }>(`/containers/${containerId}/env`, {
      method: 'POST',
      body: JSON.stringify({ variables }),
    })
  }

  async clearContainerEnvVars(containerId: string): Promise<{ success: boolean; message: string; error?: string }> {
    return this.request<{ success: boolean; message: string; error?: string }>(`/containers/${containerId}/env`, {
      method: 'DELETE',
    })
  }

  async setContainerEnvFromString(containerId: string, envString: string): Promise<{ success: boolean; message: string; error?: string }> {
    return this.request<{ success: boolean; message: string; error?: string }>(`/containers/${containerId}/env/string`, {
      method: 'POST',
      body: JSON.stringify({ envString }),
    })
  }

  async getEnvTemplates(): Promise<{ templates: EnvTemplate[]; count: number }> {
    return this.request<{ templates: EnvTemplate[]; count: number }>('/env/templates')
  }

  // Registry
  async getRegistryCatalog(): Promise<RegistryCatalog> {
    return this.request<RegistryCatalog>('/registry/catalog')
  }
  
  async getRegistryTags(repo: string): Promise<RegistryTags> {
    return this.request<RegistryTags>(`/registry/tags/${repo}`)
  }
  
  async listRegistryRepositories(): Promise<{ repositories: string[]; count: number }> {
    return this.request<{ repositories: string[]; count: number }>('/registry/repositories')
  }
  
  async listRepositoryVersions(repo: string): Promise<RepositoryVersions> {
    return this.request<RepositoryVersions>(`/registry/repositories/${encodeURIComponent(repo)}/versions`)
  }
  
  async getRepositorySummary(repo: string): Promise<VersionSummary> {
    return this.request<VersionSummary>(`/registry/repositories/${encodeURIComponent(repo)}/summary`)
  }
  
  async getVersion(repo: string, tag: string): Promise<ImageVersion> {
    return this.request<ImageVersion>(`/registry/repositories/${encodeURIComponent(repo)}/versions/${encodeURIComponent(tag)}`)
  }
  
  async deleteVersion(repo: string, tag: string): Promise<void> {
    return this.request<void>(`/registry/repositories/${encodeURIComponent(repo)}/versions/${encodeURIComponent(tag)}`, {
      method: 'DELETE',
    })
  }
  async retagImage(repo: string, source: string, target: string): Promise<void> {
    await this.request<void>('/registry/retag', {
      method: 'POST',
      body: JSON.stringify({ repo, source, target })
    })
  }
  async deleteTag(repo: string, tag: string): Promise<void> {
    await this.request<void>(`/registry/tags/${repo}/${tag}`, { method: 'DELETE' })
  }

  async scanImage(image: string): Promise<ImageScanResult> {
    return this.request<ImageScanResult>('/images/scan', {
      method: 'POST',
      body: JSON.stringify({ image })
    })
  }

  async searchLogs(params: { query?: string; containerId?: string; since?: number; until?: number; limit?: number }): Promise<LogEntry[]> {
    const qs = new URLSearchParams()
    if (params.query) qs.set('query', params.query)
    if (params.containerId) qs.set('containerId', params.containerId)
    if (params.since) qs.set('since', String(params.since))
    if (params.until) qs.set('until', String(params.until))
    if (params.limit) qs.set('limit', String(params.limit))
    return this.request<LogEntry[]>(`/logs/search?${qs.toString()}`)
  }

  async getAlerts(): Promise<AlertThresholds> {
    return this.request<AlertThresholds>('/alerts')
  }
  async setAlerts(cfg: AlertThresholds): Promise<AlertThresholds> {
    return this.request<AlertThresholds>('/alerts', { method:'POST', body: JSON.stringify(cfg) })
  }

  async registryLogin(username: string, password: string): Promise<{ token: string; token_type: string }> {
    const res = await this.request<{ token: string; token_type: string }>('/registry/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    })
    this.authToken = res.token
    localStorage.setItem('nebula_registry_token', res.token)
    return res
  }

  registryLogout() {
    this.authToken = null
    localStorage.removeItem('nebula_registry_token')
  }

  // App auth
  async login(username: string, password: string): Promise<{ token: string; user: { username: string; role: string } }> {
    const res = await this.request<{ token: string; user: { username: string; role: string } }>(`/auth/login`, {
      method: 'POST',
      body: JSON.stringify({ username, password })
    })
    this.appToken = res.token
    localStorage.setItem('nebula_app_token', res.token)
    return res
  }
  async logout(): Promise<void> {
    await this.request(`/auth/logout`, { method: 'POST' })
    this.appToken = null
    localStorage.removeItem('nebula_app_token')
  }
  async me(): Promise<{ user: { username: string; role: string } | null }> {
    return this.request('/auth/me')
  }

  // Teams/Workspaces
  async listTeams(): Promise<Array<{ id: string; name: string; description: string; created: string; createdBy: string }>> {
    return this.request('/teams')
  }
  async getTeam(id: string): Promise<{ team: { id: string; name: string; description: string; created: string; createdBy: string }; members: Array<{ username: string; role: string; joinedAt: string }> }> {
    return this.request(`/teams/${id}`)
  }
  async createTeam(name: string, description?: string): Promise<{ id: string; name: string; description: string; created: string; createdBy: string }> {
    return this.request('/teams', { method: 'POST', body: JSON.stringify({ name, description }) })
  }
  async updateTeam(id: string, name?: string, description?: string): Promise<{ id: string; name: string; description: string }> {
    return this.request(`/teams/${id}`, { method: 'PUT', body: JSON.stringify({ name, description }) })
  }
  async deleteTeam(id: string): Promise<{ deleted: string }> {
    return this.request(`/teams/${id}`, { method: 'DELETE' })
  }
  async inviteMember(teamId: string, username: string, role: 'admin' | 'editor' | 'viewer' = 'viewer'): Promise<{ username: string; role: string }> {
    return this.request(`/teams/${teamId}/invite`, { method: 'POST', body: JSON.stringify({ username, role }) })
  }
  async removeMember(teamId: string, username: string): Promise<{ removed: string }> {
    return this.request(`/teams/${teamId}/members/${username}`, { method: 'DELETE' })
  }
  async updateMemberRole(teamId: string, username: string, role: 'admin' | 'editor' | 'viewer'): Promise<{ username: string; role: string }> {
    return this.request(`/teams/${teamId}/members/${username}/role`, { method: 'PUT', body: JSON.stringify({ role }) })
  }

  // Tenants
  async listTenants(): Promise<Array<{ id: string; name: string; domain: string; created: string; createdBy: string; quota: { maxContainers: number; maxNetworks: number; maxTeams: number; maxStorageGB: number } }>> {
    return this.request('/tenants')
  }
  async getTenant(id: string): Promise<{ id: string; name: string; domain: string; created: string; createdBy: string; quota: { maxContainers: number; maxNetworks: number; maxTeams: number; maxStorageGB: number } }> {
    return this.request(`/tenants/${id}`)
  }
  async createTenant(name: string, domain?: string, quota?: { maxContainers?: number; maxNetworks?: number; maxTeams?: number; maxStorageGB?: number }): Promise<{ id: string; name: string; domain: string; quota: any }> {
    return this.request('/tenants', { method: 'POST', body: JSON.stringify({ name, domain, quota }) })
  }
  async updateTenant(id: string, name?: string, domain?: string, quota?: { maxContainers?: number; maxNetworks?: number; maxTeams?: number; maxStorageGB?: number }): Promise<{ id: string; name: string; domain: string; quota: any }> {
    return this.request(`/tenants/${id}`, { method: 'PUT', body: JSON.stringify({ name, domain, quota }) })
  }
  async deleteTenant(id: string): Promise<{ deleted: string }> {
    return this.request(`/tenants/${id}`, { method: 'DELETE' })
  }
  async assignUserToTenant(tenantId: string, username: string): Promise<{ username: string; tenantId: string }> {
    return this.request(`/tenants/${tenantId}/assign`, { method: 'POST', body: JSON.stringify({ username }) })
  }
async getTenantUsage(tenantId: string): Promise<{ tenant: any; usage: { containers: number; networks: number; teams: number }; quota: any }> {
return this.request(`/tenants/${tenantId}/usage`)
  }

  // Mode management
  async getMode(): Promise<{ mode: string; description: string; builtImages: number; totalImages: number }> {
    return this.request('/mode')
  }

  async setMode(mode: 'mock' | 'test' | 'live'): Promise<{ message: string; oldMode: string; newMode: string; mode: string; builtImagesCleared: boolean }> {
    return this.request('/mode', {
      method: 'PUT',
      body: JSON.stringify({ mode }),
    })
  }
}

export const api = new ApiClient()

