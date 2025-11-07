import { useEffect, useState } from 'react'
import { Plus, Trash2, Play, Pause, Square, Mail, Activity, Globe, Eye, FileText, BarChart3, FolderSync, RefreshCw, AlertTriangle, CheckCircle, Copy, ExternalLink, UserPlus } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { api, Workspace, Session, SessionState, Invite, Container, AuditLog, AuditStats, FileChange, CRDTOperation, Conflict, AutoSleepConfig, WorkspaceActivity } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function SharedRuntime() {
  const [activeTab, setActiveTab] = useState('workspaces')
  const [workspaces, setWorkspaces] = useState<Workspace[]>([])
  const [containers, setContainers] = useState<Container[]>([])
  const [selectedWorkspace, setSelectedWorkspace] = useState<Workspace | null>(null)
  const [sessions, setSessions] = useState<Session[]>([])
  const [activeSessions, setActiveSessions] = useState<SessionState[]>([])
  const [invites, setInvites] = useState<Invite[]>([])
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([])
  const [auditStats, setAuditStats] = useState<AuditStats | null>(null)
  const [fileChanges, setFileChanges] = useState<FileChange[]>([])
  const [selectedContainerForSync, setSelectedContainerForSync] = useState('')
  const [syncFilePath, setSyncFilePath] = useState('')
  const [syncTargetPath, setSyncTargetPath] = useState('')
  const [crdtOperations, setCrdtOperations] = useState<CRDTOperation[]>([])
  const [conflicts, setConflicts] = useState<Conflict[]>([])
  const [detectingConflicts, setDetectingConflicts] = useState(false)
  const [joinToken, setJoinToken] = useState('')
  const [joinError, setJoinError] = useState('')

  // Create workspace form
  const [workspaceName, setWorkspaceName] = useState('')
  const [workspaceDescription, setWorkspaceDescription] = useState('')
  const [selectedContainerId, setSelectedContainerId] = useState('')
  const [creating, setCreating] = useState(false)

  // Invite form
  const [inviteEmail, setInviteEmail] = useState('')
  const [inviteRole, setInviteRole] = useState('viewer')
  const [creatingInvite, setCreatingInvite] = useState(false)

  // Session creation
  const [sessionType, setSessionType] = useState('terminal')
  const [creatingSession, setCreatingSession] = useState(false)

  useEffect(() => {
    loadWorkspaces()
    loadContainers()
  }, [])

  useEffect(() => {
    if (selectedWorkspace) {
      loadWorkspaceData()
      loadAuditData()
      const interval = setInterval(() => {
        loadWorkspaceData()
        loadAuditData()
      }, 5000)
      return () => clearInterval(interval)
    }
  }, [selectedWorkspace])

  const loadWorkspaces = async () => {
    try {
      const data = await api.listWorkspaces()
      setWorkspaces(data.workspaces)
    } catch (error) {
      console.error('Failed to load workspaces:', error)
    }
  }

  const loadContainers = async () => {
    try {
      const data = await api.listContainers(true)
      setContainers(data)
    } catch (error) {
      console.error('Failed to load containers:', error)
    }
  }

  const loadWorkspaceData = async () => {
    if (!selectedWorkspace) return

    try {
      const [sessionsData, activeData, invitesData] = await Promise.all([
        api.listWorkspaceSessions(selectedWorkspace.id),
        api.listActiveSessions(selectedWorkspace.id),
        api.listInvites(selectedWorkspace.id),
      ])
      setSessions(sessionsData.sessions)
      setActiveSessions(activeData.sessions)
      // Fetch invite links for each invite
      const invitesWithLinks = await Promise.all(
        invitesData.invites.map(async (invite) => {
          try {
            const linkResult = await api.getInviteLink(selectedWorkspace.id, invite.token)
            return { ...invite, link: linkResult.link }
          } catch {
            return invite
          }
        })
      )
      setInvites(invitesWithLinks)
    } catch (error) {
      console.error('Failed to load workspace data:', error)
    }
  }

  const loadAuditData = async () => {
    if (!selectedWorkspace) return

    try {
      const [logsData, statsData] = await Promise.all([
        api.getAuditLogs(selectedWorkspace.id, { limit: 50 }),
        api.getAuditStats(selectedWorkspace.id),
      ])
      setAuditLogs(logsData.logs)
      setAuditStats(statsData)
    } catch (error) {
      console.error('Failed to load audit data:', error)
    }
  }

  const handleCreateWorkspace = async () => {
    if (!workspaceName || !selectedContainerId) {
      alert('Please provide workspace name and select a container')
      return
    }

    setCreating(true)
    try {
      await api.createWorkspace({
        name: workspaceName,
        description: workspaceDescription || undefined,
        containerId: selectedContainerId,
        settings: {
          maxMembers: 10,
          sessionTimeout: 60,
        },
      })
      setWorkspaceName('')
      setWorkspaceDescription('')
      setSelectedContainerId('')
      loadWorkspaces()
    } catch (error: any) {
      console.error('Failed to create workspace:', error)
      alert(`Failed to create workspace: ${error.message}`)
    } finally {
      setCreating(false)
    }
  }

  const handleDeleteWorkspace = async (workspaceId: string) => {
    if (!confirm('Are you sure you want to delete this workspace?')) return

    try {
      await api.deleteWorkspace(workspaceId)
      loadWorkspaces()
      if (selectedWorkspace?.id === workspaceId) {
        setSelectedWorkspace(null)
      }
    } catch (error: any) {
      console.error('Failed to delete workspace:', error)
      alert(`Failed to delete workspace: ${error.message}`)
    }
  }

  const handleShareWorkspace = async (workspaceId: string) => {
    try {
      // Try to get existing invite link first
      const invites = await api.listInvites(workspaceId)
      if (invites.invites.length > 0) {
        const invite = invites.invites[0]
        const linkResult = await api.getInviteLink(workspaceId, invite.token)
        await copyToClipboard(linkResult.link)
        alert('Invite link copied to clipboard!')
        return
      }

      // Create new invite if none exists
      const invite = await api.createInvite(workspaceId, { role: 'editor', expiresInHours: 0 })
      const linkResult = await api.getInviteLink(workspaceId, invite.token)
      await copyToClipboard(linkResult.link)
      alert('Invite link created and copied to clipboard!')
    } catch (error: any) {
      alert(error.message || 'Failed to share workspace')
    }
  }

  const handleJoinWorkspace = async () => {
    if (!joinToken.trim()) {
      setJoinError('Please enter an invite token')
      return
    }

    setJoinError('')
    try {
      await api.acceptInvite(joinToken.trim())
      alert(`Successfully joined workspace!`)
      setJoinToken('')
      loadWorkspaces()
    } catch (error: any) {
      setJoinError(error.message || 'Failed to join workspace')
    }
  }

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
    } catch (err) {
      // Fallback for older browsers
      const textArea = document.createElement('textarea')
      textArea.value = text
      textArea.style.position = 'fixed'
      textArea.style.opacity = '0'
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
    }
  }

  const handleUpdateStatus = async (workspaceId: string, status: string) => {
    try {
      await api.updateWorkspaceStatus(workspaceId, status)
      loadWorkspaces()
      if (selectedWorkspace?.id === workspaceId) {
        const updated = await api.getWorkspace(workspaceId)
        setSelectedWorkspace(updated)
      }
    } catch (error: any) {
      console.error('Failed to update status:', error)
      alert(`Failed to update status: ${error.message}`)
    }
  }

  const handleCreateInvite = async () => {
    if (!selectedWorkspace) return

    setCreatingInvite(true)
    try {
      await api.createInvite(selectedWorkspace.id, {
        email: inviteEmail || undefined,
        role: inviteRole,
        expiresInHours: 24,
      })
      setInviteEmail('')
      loadWorkspaceData()
    } catch (error: any) {
      console.error('Failed to create invite:', error)
      alert(`Failed to create invite: ${error.message}`)
    } finally {
      setCreatingInvite(false)
    }
  }

  const handleCreateSession = async () => {
    if (!selectedWorkspace) return

    setCreatingSession(true)
    try {
      await api.createSession(selectedWorkspace.id, {
        type: sessionType,
      })
      loadWorkspaceData()
    } catch (error: any) {
      console.error('Failed to create session:', error)
      alert(`Failed to create session: ${error.message}`)
    } finally {
      setCreatingSession(false)
    }
  }

  const handleRemoveMember = async (userId: string) => {
    if (!selectedWorkspace) return
    if (!confirm('Remove this member from the workspace?')) return

    try {
      await api.removeWorkspaceMember(selectedWorkspace.id, userId)
      loadWorkspaces()
      if (selectedWorkspace) {
        const updated = await api.getWorkspace(selectedWorkspace.id)
        setSelectedWorkspace(updated)
      }
    } catch (error: any) {
      console.error('Failed to remove member:', error)
      alert(`Failed to remove member: ${error.message}`)
    }
  }

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleString()
    } catch {
      return dateStr
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'success'
      case 'paused':
        return 'secondary'
      case 'stopped':
        return 'default'
      default:
        return 'default'
    }
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold mb-2">Shared Runtime</h1>
          <p className="text-muted-foreground">
            Share container access with team members for collaborative development
          </p>
        </div>
        <Card className="p-4">
          <div className="flex items-center gap-2">
            <Input
              placeholder="Enter invite token to join..."
              value={joinToken}
              onChange={(e) => {
                setJoinToken(e.target.value)
                setJoinError('')
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleJoinWorkspace()
                }
              }}
              className="w-64"
            />
            <Button onClick={handleJoinWorkspace} data-test-id={TEST_IDS.shareruntime.joinWorkspace}>
              <UserPlus className="h-4 w-4 mr-2" />
              Join
            </Button>
          </div>
          {joinError && (
            <p className="text-sm text-red-600 mt-2">{joinError}</p>
          )}
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList>
          <TabsTrigger value="workspaces">
            <Globe className="h-4 w-4 mr-2" />
            Workspaces
          </TabsTrigger>
          <TabsTrigger value="sessions" disabled={!selectedWorkspace}>
            <Activity className="h-4 w-4 mr-2" />
            Sessions
          </TabsTrigger>
          <TabsTrigger value="audit" disabled={!selectedWorkspace}>
            <FileText className="h-4 w-4 mr-2" />
            Audit Logs
          </TabsTrigger>
          <TabsTrigger value="filesync" disabled={!selectedWorkspace}>
            <FolderSync className="h-4 w-4 mr-2" />
            File Sync
          </TabsTrigger>
        </TabsList>

        {/* Workspaces Tab */}
        <TabsContent value="workspaces" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Plus className="h-5 w-5" />
                Create Shared Workspace
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Workspace Name *</Label>
                  <Input
                    value={workspaceName}
                    onChange={(e) => setWorkspaceName(e.target.value)}
                    placeholder="my-shared-workspace"
                  />
                </div>
                <div>
                  <Label>Container *</Label>
                  <select
                    value={selectedContainerId}
                    onChange={(e) => setSelectedContainerId(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md"
                  >
                    <option value="">Select a container...</option>
                    {containers.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name || c.id} ({c.image})
                      </option>
                    ))}
                  </select>
                </div>
                <div className="col-span-2">
                  <Label>Description</Label>
                  <Input
                    value={workspaceDescription}
                    onChange={(e) => setWorkspaceDescription(e.target.value)}
                    placeholder="Workspace description"
                  />
                </div>
              </div>
              <Button onClick={handleCreateWorkspace} disabled={creating} className="mt-4" data-test-id={TEST_IDS.shareruntime.createWorkspace}>
                <Plus className="h-4 w-4 mr-2" />
                {creating ? 'Creating...' : 'Create Workspace'}
              </Button>
            </CardContent>
          </Card>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {workspaces.map((workspace) => (
              <Card
                key={workspace.id}
                className={`cursor-pointer transition-all hover:shadow-lg ${
                  selectedWorkspace?.id === workspace.id ? 'ring-2 ring-primary' : ''
                }`}
                onClick={() => {
                  setSelectedWorkspace(workspace)
                  setActiveTab('sessions')
                }}
              >
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <CardTitle className="flex items-center gap-2 mb-2">
                        {workspace.name}
                        <Badge variant={getStatusColor(workspace.status) as any}>
                          {workspace.status}
                        </Badge>
                      </CardTitle>
                      <p className="text-sm text-muted-foreground">{workspace.description || 'No description'}</p>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    <div className="text-xs text-muted-foreground">
                      Container: {workspace.containerId ? workspace.containerId.substring(0, 12) : 'N/A'}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      Members: {(workspace.members || []).length} | Created: {formatDate(workspace.createdAt)}
                    </div>
                    <div className="flex gap-2 mt-4">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={(e) => {
                          e.stopPropagation()
                          setSelectedWorkspace(workspace)
                          setActiveTab('sessions')
                        }}
                      >
                        <Eye className="h-4 w-4 mr-1" />
                        View
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={async (e) => {
                          e.stopPropagation()
                          await handleShareWorkspace(workspace.id)
                        }}
                      >
                        <ExternalLink className="h-4 w-4 mr-1" />
                        Share
                      </Button>
                      {workspace.status === 'active' && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={(e) => {
                            e.stopPropagation()
                            handleUpdateStatus(workspace.id, 'paused')
                          }}
                        >
                          <Pause className="h-4 w-4 mr-1" />
                          Pause
                        </Button>
                      )}
                      {workspace.status === 'paused' && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={(e) => {
                            e.stopPropagation()
                            handleUpdateStatus(workspace.id, 'active')
                          }}
                        >
                          <Play className="h-4 w-4 mr-1" />
                          Resume
                        </Button>
                      )}
                      <Button
                        size="sm"
                        variant="destructive"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleDeleteWorkspace(workspace.id)
                        }}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        {/* Sessions Tab */}
        <TabsContent value="sessions" className="space-y-6">
          {selectedWorkspace ? (
            <>
              <Card>
                <CardHeader>
                  <CardTitle>{selectedWorkspace.name}</CardTitle>
                  <div className="text-sm text-muted-foreground mt-2">
                    Container: {selectedWorkspace.containerId ? selectedWorkspace.containerId.substring(0, 12) : 'N/A'} | 
                    Members: {(selectedWorkspace.members || []).length} | 
                    Status: {selectedWorkspace.status}
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div>
                      <h3 className="font-semibold mb-3">Members</h3>
                      <div className="space-y-2">
                        {(selectedWorkspace.members || []).map((member: any) => (
                          <div key={member.userId} className="flex items-center justify-between p-2 border rounded">
                            <div>
                              <div className="font-medium">{member.username}</div>
                              <div className="text-xs text-muted-foreground">{member.role}</div>
                            </div>
                            {member.userId !== selectedWorkspace.ownerId && (
                              <Button
                                size="sm"
                                variant="ghost"
                                onClick={() => handleRemoveMember(member.userId)}
                              >
                                <Trash2 className="h-3 w-3" />
                              </Button>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>

                    <div>
                      <h3 className="font-semibold mb-3">Create Invite</h3>
                      <div className="space-y-2">
                        <Input
                          placeholder="Email (optional)"
                          value={inviteEmail}
                          onChange={(e) => setInviteEmail(e.target.value)}
                        />
                        <select
                          value={inviteRole}
                          onChange={(e) => setInviteRole(e.target.value)}
                          className="w-full px-3 py-2 border rounded-md"
                        >
                          <option value="viewer">Viewer</option>
                          <option value="editor">Editor</option>
                          <option value="admin">Admin</option>
                        </select>
                        <Button
                          onClick={handleCreateInvite}
                          disabled={creatingInvite}
                          className="w-full"
                        >
                          <Mail className="h-4 w-4 mr-2" />
                          {creatingInvite ? 'Creating...' : 'Create Invite'}
                        </Button>
                      </div>
                      {invites.length > 0 && (
                        <div className="mt-4">
                          <div className="font-semibold text-sm mb-2">Pending Invites</div>
                          <div className="space-y-2">
                            {invites.filter(i => i.status === 'pending').map((invite) => {
                              const inviteUrl = invite.link || `${window.location.origin}/shareruntime/accept/${invite.token}`
                              return (
                                <div key={invite.id} className="text-xs p-2 border rounded">
                                  <div className="flex items-center justify-between mb-1">
                                    <div>
                                      <Badge variant="secondary">{invite.role}</Badge>
                                      <span className="ml-2">{invite.email || 'No email'}</span>
                                    </div>
                                  </div>
                                  <div className="mt-2 flex gap-2">
                                    <Input
                                      readOnly
                                      value={inviteUrl}
                                      className="text-xs flex-1"
                                      onClick={(e) => e.currentTarget.select()}
                                    />
                                    <Button
                                      size="sm"
                                      variant="outline"
                                      onClick={async () => {
                                        await copyToClipboard(inviteUrl)
                                        alert('Invite link copied to clipboard!')
                                      }}
                                    >
                                      <Copy className="h-4 w-4" />
                                    </Button>
                                    <Button
                                      size="sm"
                                      variant="outline"
                                      onClick={() => window.open(inviteUrl, '_blank')}
                                    >
                                      <ExternalLink className="h-4 w-4" />
                                    </Button>
                                  </div>
                                  <div className="text-muted-foreground mt-1">
                                    Expires: {formatDate(invite.expiresAt)}
                                  </div>
                                </div>
                              )
                            })}
                          </div>
                        </div>
                      )}
                    </div>

                    <div>
                      <h3 className="font-semibold mb-3">Create Session</h3>
                      <div className="space-y-2">
                        <select
                          value={sessionType}
                          onChange={(e) => setSessionType(e.target.value)}
                          className="w-full px-3 py-2 border rounded-md"
                        >
                          <option value="terminal">Terminal</option>
                          <option value="file">File Access</option>
                          <option value="api">API</option>
                          <option value="exec">Exec</option>
                        </select>
                        <Button
                          onClick={handleCreateSession}
                          disabled={creatingSession}
                          className="w-full"
                        >
                          <Plus className="h-4 w-4 mr-2" />
                          {creatingSession ? 'Creating...' : 'Create Session'}
                        </Button>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Active Sessions ({activeSessions.length})</CardTitle>
                </CardHeader>
                <CardContent>
                  {activeSessions.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground">
                      <Activity className="h-12 w-12 mx-auto mb-2 opacity-50" />
                      <p>No active sessions</p>
                    </div>
                  ) : (
                    <div className="space-y-3">
                      {activeSessions.map((session) => (
                        <div key={session.sessionId} className="p-3 border rounded-lg">
                          <div className="flex items-center justify-between">
                            <div>
                              <div className="font-medium">{session.username}</div>
                              <div className="text-sm text-muted-foreground">
                                {session.type} | {session.state} | Last activity: {formatDate(session.lastActivity)}
                              </div>
                            </div>
                            <Badge variant={session.state === 'active' ? 'success' : 'secondary'}>
                              {session.state}
                            </Badge>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>All Sessions ({sessions.length})</CardTitle>
                </CardHeader>
                <CardContent>
                  {sessions.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground">
                      <Activity className="h-12 w-12 mx-auto mb-2 opacity-50" />
                      <p>No sessions</p>
                    </div>
                  ) : (
                    <div className="space-y-3">
                      {sessions.map((session) => (
                        <div key={session.id} className="p-3 border rounded-lg">
                          <div className="flex items-center justify-between">
                            <div>
                              <div className="font-medium">{session.username}</div>
                              <div className="text-sm text-muted-foreground">
                                {session.type} | Started: {formatDate(session.startedAt)}
                              </div>
                            </div>
                            <Button
                              size="sm"
                              variant="destructive"
                              onClick={() => {
                                api.closeSession(session.id).then(() => loadWorkspaceData())
                              }}
                            >
                              <Square className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </>
          ) : (
            <Card>
              <CardContent className="flex flex-col items-center justify-center py-12">
                <Globe className="h-12 w-12 mb-4 opacity-50" />
                <p className="text-muted-foreground">Select a workspace to view sessions</p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* Audit Logs Tab */}
        <TabsContent value="audit" className="space-y-6">
          {selectedWorkspace ? (
            <>
              {auditStats && (
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <BarChart3 className="h-5 w-5" />
                      Audit Statistics (Last 24 Hours)
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-3 gap-4">
                      <div>
                        <div className="text-sm text-muted-foreground">Total Actions</div>
                        <div className="text-2xl font-bold">{auditStats.totalActions}</div>
                      </div>
                      <div>
                        <div className="text-sm text-muted-foreground">Successful</div>
                        <div className="text-2xl font-bold text-green-600">{auditStats.successfulActions}</div>
                      </div>
                      <div>
                        <div className="text-sm text-muted-foreground">Failed</div>
                        <div className="text-2xl font-bold text-red-600">{auditStats.failedActions}</div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              )}

              <Card>
                <CardHeader>
                  <CardTitle>Recent Audit Logs ({auditLogs.length})</CardTitle>
                </CardHeader>
                <CardContent>
                  {auditLogs.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground">
                      <FileText className="h-12 w-12 mx-auto mb-2 opacity-50" />
                      <p>No audit logs</p>
                    </div>
                  ) : (
                    <div className="space-y-2">
                      {auditLogs.map((log) => (
                        <div key={log.id} className={`p-3 border rounded-lg ${log.success ? '' : 'bg-red-50 border-red-200'}`}>
                          <div className="flex items-start justify-between">
                            <div className="flex-1">
                              <div className="flex items-center gap-2 mb-1">
                                <span className="font-medium">{log.username}</span>
                                <Badge variant={log.success ? 'success' : 'destructive'}>
                                  {log.success ? 'Success' : 'Failed'}
                                </Badge>
                                <span className="text-sm text-muted-foreground">{log.action}</span>
                              </div>
                              {log.details && Object.keys(log.details).length > 0 && (
                                <div className="text-xs text-muted-foreground mt-1">
                                  {Object.entries(log.details).map(([key, value]) => (
                                    <span key={key} className="mr-3">
                                      <strong>{key}:</strong> {String(value)}
                                    </span>
                                  ))}
                                </div>
                              )}
                              {log.errorMessage && (
                                <div className="text-xs text-red-600 mt-1">
                                  Error: {log.errorMessage}
                                </div>
                              )}
                              <div className="text-xs text-muted-foreground mt-1">
                                {formatDate(log.timestamp)}
                              </div>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </>
          ) : (
            <Card>
              <CardContent className="flex flex-col items-center justify-center py-12">
                <FileText className="h-12 w-12 mb-4 opacity-50" />
                <p className="text-muted-foreground">Select a workspace to view audit logs</p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="filesync">
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FolderSync className="h-5 w-5" />
                  File System Synchronization
                </CardTitle>
                <CardDescription>
                  Track and synchronize file changes across shared containers
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium mb-2 block">Workspace</label>
                    <select
                      className="w-full p-2 border rounded"
                      value={selectedWorkspace?.id || ''}
                      onChange={(e) => {
                        const ws = workspaces.find(w => w.id === e.target.value)
                        setSelectedWorkspace(ws || null)
                      }}
                    >
                      <option value="">Select workspace...</option>
                      {workspaces.map((w) => (
                        <option key={w.id} value={w.id}>
                          {w.name}
                        </option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="text-sm font-medium mb-2 block">Container ID</label>
                    <Input
                      placeholder="Container ID (optional)"
                      value={selectedContainerForSync}
                      onChange={(e) => setSelectedContainerForSync(e.target.value)}
                    />
                  </div>
                </div>
                <Button
                  onClick={async () => {
                    if (!selectedWorkspace) return
                    try {
                      const result = await api.getFileChanges(
                        selectedWorkspace.id,
                        selectedContainerForSync || undefined,
                        new Date(Date.now() - 3600000).toISOString()
                      )
                      setFileChanges(result.changes)
                    } catch (error: any) {
                      alert(error.message || 'Failed to fetch file changes')
                    }
                  }}
                  disabled={!selectedWorkspace}
                >
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Refresh Changes
                </Button>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Recent File Changes</CardTitle>
                <CardDescription>
                  Last {fileChanges.length} file changes in the workspace
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {fileChanges.length === 0 ? (
                    <p className="text-muted-foreground text-sm">No file changes found</p>
                  ) : (
                    fileChanges.slice(0, 50).map((change) => (
                      <div
                        key={change.id}
                        className="flex items-center justify-between p-3 border rounded"
                      >
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <Badge variant={
                              change.changeType === 'created' ? 'default' :
                              change.changeType === 'modified' ? 'secondary' :
                              change.changeType === 'deleted' ? 'destructive' : 'outline'
                            }>
                              {change.changeType}
                            </Badge>
                            <span className="font-mono text-sm">{change.path}</span>
                            {change.isDirectory && (
                              <Badge variant="outline">Directory</Badge>
                            )}
                          </div>
                          <div className="text-xs text-muted-foreground mt-1">
                            {change.containerId} • {new Date(change.timestamp).toLocaleString()} • {change.size > 0 ? `${(change.size / 1024).toFixed(2)} KB` : '0 B'}
                            {change.hash && (
                              <span className="ml-2 font-mono">Hash: {change.hash.substring(0, 16)}...</span>
                            )}
                          </div>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Sync File</CardTitle>
                <CardDescription>
                  Synchronize a file between containers using rsync-like operations
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium mb-2 block">Source File Path</label>
                    <Input
                      placeholder="/path/to/source/file"
                      value={syncFilePath}
                      onChange={(e) => setSyncFilePath(e.target.value)}
                    />
                  </div>
                  <div>
                    <label className="text-sm font-medium mb-2 block">Target File Path</label>
                    <Input
                      placeholder="/path/to/target/file"
                      value={syncTargetPath}
                      onChange={(e) => setSyncTargetPath(e.target.value)}
                    />
                  </div>
                </div>
                <Button
                  onClick={async () => {
                    if (!selectedWorkspace || !selectedContainerForSync || !syncFilePath || !syncTargetPath) {
                      alert('Please fill in all fields')
                      return
                    }
                    try {
                      await api.syncFile(selectedWorkspace.id, {
                        containerId: selectedContainerForSync,
                        filePath: syncFilePath,
                        targetPath: syncTargetPath,
                      })
                      alert('File sync initiated successfully')
                      setSyncFilePath('')
                      setSyncTargetPath('')
                    } catch (error: any) {
                      alert(error.message || 'Failed to sync file')
                    }
                  }}
                  disabled={!selectedWorkspace || !selectedContainerForSync || !syncFilePath || !syncTargetPath}
                >
                  <FolderSync className="h-4 w-4 mr-2" />
                  Sync File
                </Button>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="conflicts">
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <AlertTriangle className="h-5 w-5" />
                  Conflict Resolution (CRDT)
                </CardTitle>
                <CardDescription>
                  Detect and resolve conflicts in shared workspaces using CRDT algorithms
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex gap-4">
                  <Button
                    onClick={async () => {
                      if (!selectedWorkspace) return
                      try {
                        const result = await api.getCRDTOperations(
                          selectedWorkspace.id,
                          new Date(Date.now() - 3600000).toISOString()
                        )
                        setCrdtOperations(result.operations)
                      } catch (error: any) {
                        alert(error.message || 'Failed to fetch CRDT operations')
                      }
                    }}
                    disabled={!selectedWorkspace}
                  >
                    <RefreshCw className="h-4 w-4 mr-2" />
                    Load Operations
                  </Button>
                  <Button
                    onClick={async () => {
                      if (!selectedWorkspace) return
                      setDetectingConflicts(true)
                      try {
                        const result = await api.detectConflicts(
                          selectedWorkspace.id,
                          new Date(Date.now() - 3600000).toISOString()
                        )
                        setConflicts(result.conflicts)
                      } catch (error: any) {
                        alert(error.message || 'Failed to detect conflicts')
                      } finally {
                        setDetectingConflicts(false)
                      }
                    }}
                    disabled={!selectedWorkspace || detectingConflicts}
                    variant="outline"
                  >
                    <AlertTriangle className="h-4 w-4 mr-2" />
                    {detectingConflicts ? 'Detecting...' : 'Detect Conflicts'}
                  </Button>
                </div>
              </CardContent>
            </Card>

            {conflicts.length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle>Detected Conflicts ({conflicts.length})</CardTitle>
                  <CardDescription>
                    Conflicts requiring resolution
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {conflicts.map((conflict, idx) => (
                      <div key={idx} className="p-4 border rounded-lg bg-yellow-50">
                        <div className="flex items-start justify-between mb-3">
                          <div>
                            <div className="flex items-center gap-2 mb-2">
                              <Badge variant="destructive">{conflict.type}</Badge>
                              <span className="font-medium">{conflict.resourceType}: {conflict.resourceId}</span>
                              {conflict.resolved && (
                                <Badge variant="default">
                                  <CheckCircle className="h-3 w-3 mr-1" />
                                  Resolved
                                </Badge>
                              )}
                            </div>
                            <div className="text-sm text-muted-foreground">
                              {conflict.operations.length} concurrent operations
                            </div>
                          </div>
                          {!conflict.resolved && (
                            <Button
                              size="sm"
                              onClick={async () => {
                                if (!selectedWorkspace) return
                                try {
                                  await api.resolveConflict(selectedWorkspace.id, conflict.resourceId)
                                  alert('Conflict resolved successfully')
                                  // Refresh conflicts
                                  const result = await api.detectConflicts(
                                    selectedWorkspace.id,
                                    new Date(Date.now() - 3600000).toISOString()
                                  )
                                  setConflicts(result.conflicts)
                                } catch (error: any) {
                                  alert(error.message || 'Failed to resolve conflict')
                                }
                              }}
                            >
                              <CheckCircle className="h-4 w-4 mr-1" />
                              Resolve
                            </Button>
                          )}
                        </div>
                        <div className="space-y-2 mt-3">
                          {conflict.operations.map((op, opIdx) => (
                            <div key={opIdx} className="p-2 bg-white rounded border text-sm">
                              <div className="flex items-center gap-2">
                                <Badge variant="outline">{op.type}</Badge>
                                <span className="font-mono">{op.operation}</span>
                                <span className="text-muted-foreground">by {op.userId}</span>
                                <span className="text-muted-foreground ml-auto">
                                  {new Date(op.timestamp).toLocaleString()}
                                </span>
                              </div>
                              {op.value && (
                                <div className="mt-1 text-xs text-muted-foreground">
                                  Value: {typeof op.value === 'string' ? op.value : JSON.stringify(op.value)}
                                </div>
                              )}
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )}

            <Card>
              <CardHeader>
                <CardTitle>CRDT Operations ({crdtOperations.length})</CardTitle>
                <CardDescription>
                  Recent CRDT operations for conflict tracking
                </CardDescription>
              </CardHeader>
              <CardContent>
                {crdtOperations.length === 0 ? (
                  <p className="text-muted-foreground text-sm">No operations found</p>
                ) : (
                  <div className="space-y-2">
                    {crdtOperations.slice(0, 50).map((op) => (
                      <div key={op.id} className="p-3 border rounded">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            <Badge variant="outline">{op.type}</Badge>
                            <span className="font-mono text-sm">{op.operation}</span>
                            <span className="text-xs text-muted-foreground">
                              {op.resourceType}: {op.resourceId}
                            </span>
                          </div>
                          <div className="flex items-center gap-2 text-xs text-muted-foreground">
                            <span>{op.userId}</span>
                            <span>•</span>
                            <span>{new Date(op.timestamp).toLocaleString()}</span>
                          </div>
                        </div>
                        {op.value && (
                          <div className="mt-1 text-xs text-muted-foreground">
                            {typeof op.value === 'string' ? op.value : JSON.stringify(op.value)}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="autosleep">
          <AutoSleepTab workspace={selectedWorkspace} />
        </TabsContent>
      </Tabs>
    </div>
  )
}

function AutoSleepTab({ workspace }: { workspace: Workspace | null }) {
  const [config, setConfig] = useState<AutoSleepConfig | null>(null)
  const [idleWorkspaces, setIdleWorkspaces] = useState<WorkspaceActivity[]>([])
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (workspace) {
      loadConfig()
      loadIdleWorkspaces()
    }
  }, [workspace])

  const loadConfig = async () => {
    if (!workspace) return
    try {
      const result = await api.getAutoSleepConfig(workspace.id)
      setConfig(result.config)
    } catch (error: any) {
      alert(error.message || 'Failed to load auto-sleep config')
    }
  }

  const loadIdleWorkspaces = async () => {
    try {
      const result = await api.getIdleWorkspaces()
      setIdleWorkspaces(result.workspaces)
    } catch (error: any) {
      alert(error.message || 'Failed to load idle workspaces')
    }
  }

  const handleSaveConfig = async () => {
    if (!workspace || !config) return
    setSaving(true)
    try {
      await api.setAutoSleepConfig(workspace.id, config)
      alert('Auto-sleep configuration saved')
    } catch (error: any) {
      alert(error.message || 'Failed to save config')
    } finally {
      setSaving(false)
    }
  }

  const handleWakeWorkspace = async (workspaceId: string) => {
    try {
      await api.wakeWorkspace(workspaceId)
      alert('Workspace woken up')
      loadIdleWorkspaces()
    } catch (error: any) {
      alert(error.message || 'Failed to wake workspace')
    }
  }

  const handleRecordActivity = async () => {
    if (!workspace) return
    try {
      await api.recordWorkspaceActivity(workspace.id)
      alert('Activity recorded')
      loadIdleWorkspaces()
    } catch (error: any) {
      alert(error.message || 'Failed to record activity')
    }
  }

  if (!workspace) {
    return (
      <Card>
        <CardContent className="pt-6">
          <p className="text-muted-foreground text-sm text-center">
            Select a workspace to configure auto-sleep
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Pause className="h-5 w-5" />
            Auto-Sleep Configuration
          </CardTitle>
          <CardDescription>
            Configure automatic sleep and snapshot for idle workspaces
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {config ? (
            <>
              <div className="flex items-center justify-between">
                <Label htmlFor="enabled">Enable Auto-Sleep</Label>
                <input
                  type="checkbox"
                  id="enabled"
                  checked={config.enabled}
                  onChange={(e) => setConfig({ ...config, enabled: e.target.checked })}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="idleTimeout">Idle Timeout (minutes)</Label>
                <Input
                  id="idleTimeout"
                  type="number"
                  value={config.idleTimeout}
                  onChange={(e) => setConfig({ ...config, idleTimeout: parseInt(e.target.value) || 30 })}
                  min={1}
                  disabled={!config.enabled}
                />
                <p className="text-xs text-muted-foreground">
                  Workspace is considered idle after this duration of inactivity
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="sleepTimeout">Sleep Timeout (minutes)</Label>
                <Input
                  id="sleepTimeout"
                  type="number"
                  value={config.sleepTimeout}
                  onChange={(e) => setConfig({ ...config, sleepTimeout: parseInt(e.target.value) || 15 })}
                  min={1}
                  disabled={!config.enabled}
                />
                <p className="text-xs text-muted-foreground">
                  Workspace will be put to sleep after this additional duration
                </p>
              </div>

              <div className="flex items-center justify-between">
                <Label htmlFor="createSnapshot">Create Snapshot Before Sleep</Label>
                <input
                  type="checkbox"
                  id="createSnapshot"
                  checked={config.createSnapshot}
                  onChange={(e) => setConfig({ ...config, createSnapshot: e.target.checked })}
                  disabled={!config.enabled}
                />
              </div>

              <div className="flex items-center justify-between">
                <Label htmlFor="autoWake">Auto-Wake On Access</Label>
                <input
                  type="checkbox"
                  id="autoWake"
                  checked={config.autoWakeOnAccess}
                  onChange={(e) => setConfig({ ...config, autoWakeOnAccess: e.target.checked })}
                  disabled={!config.enabled}
                />
              </div>

              <div className="flex gap-4">
                <Button onClick={handleSaveConfig} disabled={saving}>
                  {saving ? 'Saving...' : 'Save Configuration'}
                </Button>
                <Button onClick={handleRecordActivity} variant="outline">
                  <Activity className="h-4 w-4 mr-2" />
                  Record Activity
                </Button>
              </div>
            </>
          ) : (
            <p className="text-muted-foreground text-sm">Loading configuration...</p>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Idle Workspaces</CardTitle>
          <CardDescription>
            Workspaces that are currently idle and may be put to sleep
          </CardDescription>
        </CardHeader>
        <CardContent>
          {idleWorkspaces.length === 0 ? (
            <p className="text-muted-foreground text-sm">No idle workspaces found</p>
          ) : (
            <div className="space-y-4">
              {idleWorkspaces.map((activity) => (
                <div key={activity.workspaceId} className="p-4 border rounded-lg">
                  <div className="flex items-start justify-between">
                    <div>
                      <div className="font-semibold mb-2">{activity.workspaceName}</div>
                      <div className="text-sm text-muted-foreground space-y-1">
                        <div>Last Activity: {new Date(activity.lastActivity).toLocaleString()}</div>
                        <div>Idle Duration: {activity.idleDuration}</div>
                        <div>Status: <Badge variant="outline">{activity.status}</Badge></div>
                      </div>
                    </div>
                    {activity.status === 'sleeping' && (
                      <Button
                        size="sm"
                        onClick={() => handleWakeWorkspace(activity.workspaceId)}
                      >
                        <Play className="h-4 w-4 mr-1" />
                        Wake
                      </Button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

