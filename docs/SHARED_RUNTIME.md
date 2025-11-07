# Shared Runtime Layer - Technical Blueprint

## Overview

The Shared Runtime layer enables secure, multi-user collaborative access to containerized environments. It provides workspace isolation, session multiplexing, secure tunneling, real-time synchronization, and conflict resolution for teams working together on shared container environments.

## Table of Contents

1. [Architecture](#architecture)
2. [Core Components](#core-components)
3. [Data Models](#data-models)
4. [Security Model](#security-model)
5. [Session Management](#session-management)
6. [Synchronization](#synchronization)
7. [API Reference](#api-reference)
8. [Integration Points](#integration-points)
9. [Deployment Considerations](#deployment-considerations)

---

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      Shared Runtime Layer                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Workspace   │  │   Session    │  │    Tunnel    │          │
│  │   Manager    │◄─┤  Manager     │◄─┤   Manager    │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                 │                 │                   │
│         └─────────────────┼─────────────────┘                   │
│                           │                                     │
│  ┌──────────────┐  ┌──────▼───────┐  ┌──────────────┐          │
│  │    Invite    │  │   NebulaSync │  │     CRDT     │          │
│  │   Manager    │  │   Manager    │  │   Manager    │          │
│  └──────────────┘  └──────┬───────┘  └──────┬───────┘          │
│                           │                 │                   │
│  ┌──────────────┐  ┌──────▼───────┐  ┌──────▼───────┐          │
│  │    Audit     │  │   FileSync   │  │   Auto-Sleep │          │
│  │   Logger     │  │   Manager    │  │   Manager    │          │
│  └──────────────┘  └──────────────┘  └───────────────┘          │
│                                                                   │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Container Runtime Layer                       │
│  (containerd / Nebula Runtime)                                   │
└─────────────────────────────────────────────────────────────────┘
```

### Component Interaction Flow

```
┌──────────┐         ┌──────────────┐         ┌─────────────┐
│   User   │────────▶│   Workspace  │────────▶│  Container  │
│  (Owner) │         │   Manager    │         │  Runtime    │
└──────────┘         └──────┬───────┘         └─────────────┘
                             │
                             │ Creates
                             ▼
                    ┌────────────────┐
                    │   Workspace    │
                    │   (Shared)     │
                    └────────┬───────┘
                             │
                ┌────────────┼────────────┐
                │            │            │
                ▼            ▼            ▼
        ┌───────────┐ ┌──────────┐ ┌──────────┐
        │  Invite   │ │ Session  │ │  Tunnel  │
        │  System   │ │  Mux     │ │ Manager  │
        └─────┬─────┘ └────┬─────┘ └────┬─────┘
              │            │            │
              └────────────┼────────────┘
                           │
                           ▼
                  ┌─────────────────┐
                  │   Multiple      │
                  │   Users Join    │
                  └─────────────────┘
```

### Session Multiplexing Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Shared Container                         │
│                 (Single Instance)                          │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│  User A     │ │  User B     │ │  User C    │
│  Session 1  │ │  Session 2  │ │  Session 3  │
│  (Terminal) │ │  (File Edit)│ │  (API Call)│
└──────┬──────┘ └──────┬──────┘ └──────┬──────┘
       │              │                │
       └──────────────┼────────────────┘
                      │
                      ▼
            ┌──────────────────┐
            │ Session          │
            │ Multiplexer      │
            │ (Synchronizes)   │
            └──────────────────┘
```

---

## Core Components

### 1. Workspace Manager

**Purpose**: Manages shared runtime workspaces and their lifecycle.

**Key Responsibilities**:
- Create, update, delete workspaces
- Manage workspace members and roles
- Track workspace status (active, paused, stopped, sleeping)
- Enforce workspace settings and resource limits

**Key Methods**:
```go
CreateWorkspace(name, description, ownerID, containerID string, settings WorkspaceSettings) (*Workspace, error)
GetWorkspace(id string) (*Workspace, bool)
ListWorkspaces(userID string) []*Workspace
UpdateWorkspaceStatus(id, status string) error
AddWorkspaceMember(workspaceID, userID, username, role string) error
RemoveWorkspaceMember(workspaceID, userID string) error
UpdateWorkspaceMemberRole(workspaceID, userID, role string) error
```

### 2. Session Manager

**Purpose**: Handles session multiplexing and synchronization.

**Key Responsibilities**:
- Register and manage multiple user sessions to a single container
- Track session activity and state
- Synchronize session state across users
- Handle session subscriptions and notifications

**Key Methods**:
```go
RegisterSession(workspaceID, sessionID, userID, username, sessionType string) error
UpdateActivity(workspaceID, sessionID string) error
UpdateSessionState(workspaceID, sessionID, state string) error
GetSessionState(workspaceID, sessionID string) (*SessionState, error)
ListActiveSessions(workspaceID string) []*SessionState
SubscribeActivity(workspaceID, userID string) <-chan *SessionState
```

### 3. Tunnel Manager

**Purpose**: Manages secure tunnels for container port access.

**Key Responsibilities**:
- Create secure tunnels to container ports
- Implement port multiplexing for multiple users
- Manage tunnel authentication and authorization
- Track tunnel connections and activity

**Key Methods**:
```go
CreateTunnel(workspaceID, userID, username, containerID string, containerPort int, protocol string) (*Tunnel, error)
GetTunnel(tunnelID string) (*Tunnel, bool)
ListTunnels(workspaceID string) []*Tunnel
CloseTunnel(tunnelID string) error
ValidateTunnelAccess(tunnelID, token string) (*Tunnel, error)
```

### 4. Invite Manager

**Purpose**: Manages workspace invitations and access tokens.

**Key Responsibilities**:
- Generate secure invite tokens
- Manage invite lifecycle (create, accept, reject, revoke)
- Support role-based invites (viewer, editor, admin)
- Cleanup expired invites

**Key Methods**:
```go
CreateInvite(workspaceID, inviterID, email, role string, expiresIn time.Duration) (*Invite, error)
GetInviteByToken(token string) (*Invite, error)
AcceptInvite(token, userID, username string) error
RejectInvite(token string) error
RevokeInvite(inviteID string) error
ListWorkspaceInvites(workspaceID string) []*Invite
```

### 5. NebulaSync Manager

**Purpose**: Synchronizes database changes across workspace instances.

**Key Responsibilities**:
- Track database changes using replication adapters
- Propagate changes to all workspace members
- Provide change history and subscription
- Support multiple replication adapters (in-memory, PostgreSQL, etc.)

**Key Methods**:
```go
RecordChange(workspaceID, resourceType, resourceID, operation string, data interface{}) error
GetChangesSince(workspaceID string, sinceChangeID int64) ([]*SyncEvent, error)
GetLatestChangeID(workspaceID string) (int64, error)
SubscribeToChanges(workspaceID string) <-chan *SyncEvent
```

### 6. FileSync Manager

**Purpose**: Synchronizes filesystem changes across workspace instances.

**Key Responsibilities**:
- Track file changes (create, update, delete)
- Calculate file hashes for integrity verification
- Propagate file changes via events
- Support file synchronization operations

**Key Methods**:
```go
RecordFileChange(workspaceID, containerID, path, operation string, hash string) error
GetFileChanges(workspaceID string, sinceChangeID int64) ([]*FileChange, error)
SyncFile(workspaceID, containerID, path string) error
GetFileHash(workspaceID, containerID, path string) (string, error)
SubscribeToFileChanges(workspaceID string) <-chan *FileChange
```

### 7. CRDT Manager (Conflict Resolution)

**Purpose**: Handles conflict-free replicated data types for workspace state.

**Key Responsibilities**:
- Manage CRDT operations (ORSet, LWWRegister, etc.)
- Detect conflicts between concurrent operations
- Resolve conflicts automatically using CRDT semantics
- Maintain vector clocks for causal ordering

**Key Methods**:
```go
RecordCRDTOperation(workspaceID, resourceType, resourceID, operation string, data interface{}) error
GetCRDTOperations(workspaceID string, sinceTimestamp time.Time) ([]*CRDTOperation, error)
DetectConflicts(workspaceID, resourceType, resourceID string) ([]*Conflict, error)
ResolveConflict(workspaceID, conflictID string, resolution interface{}) error
GetResourceState(workspaceID, resourceType, resourceID string) (interface{}, error)
```

### 8. Auto-Sleep Manager

**Purpose**: Automatically pauses idle workspaces and creates snapshots.

**Key Responsibilities**:
- Monitor workspace activity
- Detect idle workspaces based on configurable timeouts
- Automatically create snapshots before sleeping
- Wake workspaces on activity or manual trigger

**Key Methods**:
```go
SetAutoSleepConfig(workspaceID string, config AutoSleepConfig) error
GetAutoSleepConfig(workspaceID string) (*AutoSleepConfig, error)
RecordActivity(workspaceID string) error
WakeWorkspace(workspaceID string) error
GetIdleWorkspaces() []string
```

### 9. Audit Logger

**Purpose**: Records all significant actions within shared workspaces.

**Key Responsibilities**:
- Log all workspace operations (create, update, delete)
- Track user actions and timestamps
- Provide audit log querying and filtering
- Generate audit statistics

**Key Methods**:
```go
LogAction(workspaceID, userID, action, resourceType, resourceID string, success bool, details interface{}) error
GetAuditLogs(workspaceID string, filters AuditFilters) ([]*AuditLog, error)
GetUserAuditLogs(userID string, filters AuditFilters) ([]*AuditLog, error)
GetAuditStats(workspaceID string) (*AuditStats, error)
```

---

## Data Models

### Workspace

```go
type Workspace struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    OwnerID     string                 `json:"ownerId"`
    ContainerID string                 `json:"containerId"`
    Status      string                 `json:"status"` // active, paused, stopped, sleeping
    Members     []WorkspaceMember      `json:"members"`
    Settings    WorkspaceSettings      `json:"settings"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt   time.Time              `json:"createdAt"`
    UpdatedAt   time.Time              `json:"updatedAt"`
}
```

### WorkspaceMember

```go
type WorkspaceMember struct {
    UserID    string    `json:"userId"`
    Username  string    `json:"username"`
    Role      string    `json:"role"` // owner, admin, editor, viewer
    JoinedAt  time.Time `json:"joinedAt"`
    LastSeen  time.Time `json:"lastSeen,omitempty"`
    IsActive  bool      `json:"isActive"`
}
```

### Session

```go
type Session struct {
    ID          string            `json:"id"`
    WorkspaceID string            `json:"workspaceId"`
    UserID      string            `json:"userId"`
    Username    string            `json:"username"`
    Type        string            `json:"type"` // terminal, file, api, exec
    Connection  SessionConnection `json:"connection"`
    Metadata    map[string]string `json:"metadata,omitempty"`
    StartedAt   time.Time         `json:"startedAt"`
    LastActivity time.Time        `json:"lastActivity"`
}
```

### Tunnel

```go
type Tunnel struct {
    ID           string            `json:"id"`
    WorkspaceID  string            `json:"workspaceId"`
    UserID       string            `json:"userId"`
    Username     string            `json:"username"`
    ContainerID  string            `json:"containerId"`
    ContainerPort int              `json:"containerPort"`
    HostPort     int               `json:"hostPort"`
    TunnelPort   int               `json:"tunnelPort"`
    Protocol     string            `json:"protocol"` // tcp, udp
    Status       string            `json:"status"`   // active, closed, error
    Token        string            `json:"token"`
    AllowedIPs   []string          `json:"allowedIPs,omitempty"`
    CreatedAt    time.Time         `json:"createdAt"`
    LastActivity time.Time         `json:"lastActivity"`
}
```

### Invite

```go
type Invite struct {
    ID          string    `json:"id"`
    WorkspaceID string    `json:"workspaceId"`
    InviterID   string    `json:"inviterId"`
    Email       string    `json:"email,omitempty"`
    Role        string    `json:"role"` // viewer, editor, admin
    Token       string    `json:"token"`
    Status      string    `json:"status"` // pending, accepted, rejected, revoked
    ExpiresAt   time.Time `json:"expiresAt"`
    CreatedAt   time.Time `json:"createdAt"`
}
```

---

## Security Model

### Authentication & Authorization

```
┌─────────────┐
│   User      │
│   Request   │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  Auth Token     │
│  (Bearer JWT)   │
└──────┬──────────┘
       │
       ▼
┌─────────────────┐
│ Permission      │
│ Check           │
└──────┬──────────┘
       │
       ├──► Owner: Full Access
       ├──► Admin: Manage Members, Modify Settings
       ├──► Editor: Modify Resources, Create Sessions
       └──► Viewer: Read-Only Access
```

### Role-Based Permissions

| Action | Owner | Admin | Editor | Viewer |
|--------|-------|-------|--------|--------|
| View Workspace | ✅ | ✅ | ✅ | ✅ |
| Modify Settings | ✅ | ✅ | ❌ | ❌ |
| Add Members | ✅ | ✅ | ❌ | ❌ |
| Remove Members | ✅ | ✅ | ❌ | ❌ |
| Create Session | ✅ | ✅ | ✅ | ❌ |
| Create Tunnel | ✅ | ✅ | ✅ | ❌ |
| Modify Container | ✅ | ✅ | ✅ | ❌ |
| View Audit Logs | ✅ | ✅ | ❌ | ❌ |
| Delete Workspace | ✅ | ❌ | ❌ | ❌ |

### Secure Tunneling

```
┌──────────┐                    ┌──────────────┐
│   User   │                    │   Container  │
│  Client  │                    │    Port      │
└────┬─────┘                    └──────┬───────┘
     │                                  │
     │ 1. Create Tunnel Request        │
     ├─────────────────────────────────▶
     │                                  │
     │ 2. Token Generated               │
     ◄──────────────────────────────────┤
     │                                  │
     │ 3. Connect with Token            │
     ├─────────────────────────────────▶
     │                                  │
     │ 4. Traffic Multiplexed           │
     │◄─────────────────────────────────▶
```

**Security Features**:
- Token-based authentication for tunnels
- IP whitelisting support
- Connection tracking and statistics
- Automatic cleanup of inactive tunnels

---

## Session Management

### Session Lifecycle

```
Create Session
     │
     ▼
┌─────────────┐
│  Register   │
│  Session     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Active    │───► Update Activity ───► Active
└──────┬──────┘                           │
       │                                  │
       │ Idle Timeout                    │
       ▼                                  │
┌─────────────┐                           │
│    Idle     │                           │
└──────┬──────┘                           │
       │                                  │
       │ Activity                         │
       ◄──────────────────────────────────┘
       │
       │ Close
       ▼
┌─────────────┐
│  Closed     │
└─────────────┘
```

### Session Synchronization

```
User A Action
     │
     ▼
┌─────────────┐
│  Session    │
│  State      │
│  Update     │
└──────┬──────┘
       │
       ├──► Broadcast to Subscribers
       │
       ├──► User B Notification
       ├──► User C Notification
       └──► User D Notification
```

---

## Synchronization

### Database Sync (NebulaSync)

```
┌─────────────────────────────────────────────────────────────┐
│                    Database Changes                           │
└────────────────────┬────────────────────────────────────────┘
                     │
         ┌───────────┼───────────┐
         │           │           │
         ▼           ▼           ▼
    ┌────────┐  ┌────────┐  ┌────────┐
    │  User  │  │  User  │  │  User  │
    │   A    │  │   B    │  │   C    │
    └───┬────┘  └───┬────┘  └───┬────┘
        │           │           │
        └───────────┼───────────┘
                    │
                    ▼
            ┌───────────────┐
            │  Replication  │
            │   Adapter     │
            └───────┬───────┘
                    │
                    ▼
            ┌───────────────┐
            │  Sync Manager │
            │  (Orchestrates)│
            └───────┬───────┘
                    │
                    ▼
            ┌───────────────┐
            │   Sync Event   │
            │   Propagation  │
            └────────────────┘
```

### File Sync Flow

```
File Change Detected
        │
        ▼
┌───────────────┐
│  Calculate    │
│  File Hash     │
└───────┬───────┘
        │
        ▼
┌───────────────┐
│ Record Change │
│ Event         │
└───────┬───────┘
        │
        ▼
┌───────────────┐
│ Broadcast to  │
│ Subscribers   │
└───────┬───────┘
        │
        ▼
┌───────────────┐
│  Sync File    │
│  (if needed)  │
└───────────────┘
```

### CRDT Conflict Resolution

```
Concurrent Operations
        │
        ├──► Operation A (User 1)
        │
        └──► Operation B (User 2)
                │
                ▼
        ┌───────────────┐
        │  Vector Clock │
        │  Comparison   │
        └───────┬───────┘
                │
        ┌───────┼───────┐
        │               │
        ▼               ▼
    Causal          Concurrent
    Order           Conflict
        │               │
        │               ▼
        │       ┌───────────────┐
        │       │  CRDT Merge   │
        │       │  (ORSet/LWW)  │
        │       └───────┬───────┘
        │               │
        └───────┬───────┘
                │
                ▼
        ┌───────────────┐
        │  Resolved     │
        │  State        │
        └───────────────┘
```

---

## API Reference

### Workspace Management

#### Create Workspace
```http
POST /api/shareruntime/workspaces
Content-Type: application/json

{
  "name": "My Shared Workspace",
  "description": "Team development environment",
  "containerId": "container-123",
  "settings": {
    "maxMembers": 10,
    "sessionTimeout": 60,
    "auditLogging": true
  }
}
```

#### List Workspaces
```http
GET /api/shareruntime/workspaces
```

#### Get Workspace
```http
GET /api/shareruntime/workspaces/{id}
```

#### Add Workspace Member
```http
POST /api/shareruntime/workspaces/{id}/members
Content-Type: application/json

{
  "userId": "user-456",
  "username": "alice",
  "role": "editor"
}
```

### Invite Management

#### Create Invite
```http
POST /api/shareruntime/workspaces/{id}/invites
Content-Type: application/json

{
  "email": "user@example.com",
  "role": "editor",
  "expiresIn": "24h"
}
```

#### Get Invite Link
```http
GET /api/shareruntime/workspaces/{id}/invites/{token}/link
```

#### Accept Invite
```http
POST /api/shareruntime/invites/{token}/accept
Content-Type: application/json

{
  "username": "newuser"
}
```

### Session Management

#### Create Session
```http
POST /api/shareruntime/workspaces/{id}/sessions
Content-Type: application/json

{
  "type": "terminal",
  "metadata": {}
}
```

#### List Active Sessions
```http
GET /api/shareruntime/workspaces/{id}/sessions
```

#### Update Session Activity
```http
PUT /api/shareruntime/workspaces/{id}/sessions/{sessionId}/activity
```

### Tunnel Management

#### Create Tunnel
```http
POST /api/shareruntime/workspaces/{id}/tunnels
Content-Type: application/json

{
  "containerId": "container-123",
  "containerPort": 8080,
  "protocol": "tcp"
}
```

#### Connect to Tunnel
```http
POST /api/shareruntime/tunnels/{id}/connect
Content-Type: application/json

{
  "token": "tunnel-token-xyz"
}
```

### Synchronization

#### Get Changes Since
```http
GET /api/shareruntime/workspaces/{id}/sync/changes?since={changeId}
```

#### Record File Change
```http
POST /api/shareruntime/workspaces/{id}/filesync/changes
Content-Type: application/json

{
  "containerId": "container-123",
  "path": "/app/config.json",
  "operation": "update",
  "hash": "sha256:abc123..."
}
```

#### Detect Conflicts
```http
GET /api/shareruntime/workspaces/{id}/crdt/conflicts?resourceType={type}&resourceId={id}
```

### Auto-Sleep

#### Set Auto-Sleep Config
```http
PUT /api/shareruntime/workspaces/{id}/autosleep/config
Content-Type: application/json

{
  "enabled": true,
  "idleTimeoutMinutes": 30,
  "createSnapshot": true,
  "autoWakeOnAccess": true
}
```

#### Wake Workspace
```http
POST /api/shareruntime/workspaces/{id}/wake
```

---

## Integration Points

### Container Runtime Integration

```
Shared Runtime Layer
        │
        ├──► Container Lifecycle
        │    (Start, Stop, Pause)
        │
        ├──► Resource Management
        │    (CPU, Memory, Disk)
        │
        └──► Snapshot Management
             (Create, Restore)
```

### Authentication Integration

```
Shared Runtime Layer
        │
        ├──► User Authentication
        │    (Token Validation)
        │
        ├──► Role Verification
        │    (Permission Checks)
        │
        └──► Session Management
             (Token Refresh)
```

### Monitoring Integration

```
Shared Runtime Layer
        │
        ├──► Activity Tracking
        │    (Session Activity)
        │
        ├──► Audit Logging
        │    (All Actions)
        │
        └──► Resource Metrics
             (CPU, Memory Usage)
```

---

## Deployment Considerations

### Scalability

**Horizontal Scaling**:
- Workspace managers can be distributed across multiple instances
- Session state can be stored in distributed cache (Redis)
- Tunnel connections can be load-balanced

**Vertical Scaling**:
- Each workspace manager can handle hundreds of workspaces
- Session multiplexer supports up to 100 concurrent sessions per workspace
- Tunnel manager can handle thousands of concurrent tunnels

### Performance

**Optimization Strategies**:
- Session state cached in memory with periodic persistence
- File sync uses incremental hashing
- CRDT operations batched to reduce network overhead
- Audit logs written asynchronously

**Benchmarks**:
- Workspace creation: < 50ms
- Session registration: < 10ms
- Tunnel creation: < 20ms
- File change detection: < 5ms
- Conflict detection: < 100ms (depends on operation history)

### Reliability

**High Availability**:
- Workspace state can be replicated
- Session state synchronized via message queue
- Tunnel connections automatically reconnected on failure
- Snapshot backup before critical operations

**Disaster Recovery**:
- Workspace snapshots stored in durable storage
- Audit logs persisted to external storage
- Invite tokens can be regenerated
- Session state can be recovered from snapshots

### Security

**Best Practices**:
- All API endpoints require authentication
- Tunnel tokens are cryptographically secure
- Invite tokens expire automatically
- IP whitelisting for sensitive operations
- Audit logging for compliance

**Threat Mitigation**:
- Rate limiting on invite creation
- Connection throttling per user
- Resource limits enforced per workspace
- Automatic cleanup of inactive sessions

---

## Future Enhancements

1. **WebRTC P2P Sync**: Direct peer-to-peer file synchronization
2. **Multi-Region Support**: Cross-region workspace replication
3. **Advanced CRDTs**: More sophisticated conflict resolution strategies
4. **Workspace Templates**: Pre-configured workspace setups
5. **Integration APIs**: Webhooks and external system integration
6. **Mobile Support**: Native mobile apps for workspace access

---

## Conclusion

The Shared Runtime layer provides a comprehensive solution for secure, collaborative container environments. It combines workspace isolation, session multiplexing, secure tunneling, and conflict-free synchronization to enable teams to work together effectively while maintaining security and performance.

For implementation details, see:
- `internal/shareruntime/workspace.go` - Workspace management
- `internal/shareruntime/sessions.go` - Session multiplexing
- `internal/shareruntime/tunnel.go` - Secure tunneling
- `internal/shareruntime/sync.go` - Database synchronization
- `internal/shareruntime/filesync.go` - Filesystem synchronization
- `internal/shareruntime/crdt.go` - Conflict resolution
- `internal/shareruntime/autosleep.go` - Auto-sleep management

