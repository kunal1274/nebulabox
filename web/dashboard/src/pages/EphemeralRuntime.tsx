import { useEffect, useState } from 'react'
import { Plus, Play, Pause, Cloud, Clock, Users, Activity, RefreshCw, Power } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { api, type EphemeralRuntime, type Workspace } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function EphemeralRuntime() {
  const [runtimes, setRuntimes] = useState<EphemeralRuntime[]>([])
  const [workspaces, setWorkspaces] = useState<Workspace[]>([])
  const [selectedWorkspaceId, setSelectedWorkspaceId] = useState('')
  const [loading, setLoading] = useState(false)

  // Provision form
  const [name, setName] = useState('')
  const [region, setRegion] = useState('us-east-1')
  const [instanceType, setInstanceType] = useState<'small' | 'medium' | 'large'>('small')
  const [image, setImage] = useState('')
  const [duration, setDuration] = useState(24)
  const [members, setMembers] = useState('')

  useEffect(() => {
    loadWorkspaces()
    loadRuntimes()
  }, [])

  const loadWorkspaces = async () => {
    try {
      const result = await api.listWorkspaces()
      setWorkspaces(result.workspaces)
      if (result.workspaces.length > 0 && !selectedWorkspaceId) {
        setSelectedWorkspaceId(result.workspaces[0].id)
      }
    } catch (error: any) {
      alert(error.message || 'Failed to load workspaces')
    }
  }

  const loadRuntimes = async () => {
    try {
      const result = await api.listEphemeralRuntimes(selectedWorkspaceId || undefined)
      setRuntimes(result.runtimes)
    } catch (error: any) {
      alert(error.message || 'Failed to load runtimes')
    }
  }

  useEffect(() => {
    loadRuntimes()
  }, [selectedWorkspaceId])

  const handleProvision = async () => {
    if (!name || !image || !selectedWorkspaceId) {
      alert('Please fill in all required fields')
      return
    }

    setLoading(true)
    try {
      const memberList = members.split(',').map(m => m.trim()).filter(m => m)
      await api.provisionEphemeralRuntime(selectedWorkspaceId, {
        name,
        region,
        instanceType,
        image,
        duration,
        members: memberList.length > 0 ? memberList : undefined,
      })
      alert('Runtime provisioning initiated')
      setName('')
      setImage('')
      setDuration(24)
      setMembers('')
      loadRuntimes()
    } catch (error: any) {
      alert(error.message || 'Failed to provision runtime')
    } finally {
      setLoading(false)
    }
  }

  const handleTerminate = async (runtimeId: string) => {
    if (!confirm('Are you sure you want to terminate this runtime? This action cannot be undone.')) {
      return
    }

    try {
      await api.terminateEphemeralRuntime(runtimeId)
      alert('Runtime termination initiated')
      loadRuntimes()
    } catch (error: any) {
      alert(error.message || 'Failed to terminate runtime')
    }
  }

  const handleSleep = async (runtimeId: string) => {
    try {
      await api.sleepEphemeralRuntime(runtimeId)
      alert('Runtime put to sleep')
      loadRuntimes()
    } catch (error: any) {
      alert(error.message || 'Failed to sleep runtime')
    }
  }

  const handleWake = async (runtimeId: string) => {
    try {
      await api.wakeEphemeralRuntime(runtimeId)
      alert('Runtime woken up')
      loadRuntimes()
    } catch (error: any) {
      alert(error.message || 'Failed to wake runtime')
    }
  }

  const handleUpdateActivity = async (runtimeId: string) => {
    try {
      await api.updateEphemeralRuntimeActivity(runtimeId)
      loadRuntimes()
    } catch (error: any) {
      alert(error.message || 'Failed to update activity')
    }
  }

  const getStatusBadge = (status: string) => {
    const variants: Record<string, any> = {
      provisioning: 'default',
      active: 'default',
      idle: 'secondary',
      sleeping: 'outline',
      terminating: 'destructive',
      terminated: 'destructive',
    }
    return <Badge variant={variants[status] || 'default'}>{status}</Badge>
  }


  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Cloud className="h-8 w-8" />
            Ephemeral Runtimes
          </h1>
          <p className="text-muted-foreground mt-2">
            Provision temporary remote runtimes for team testing
          </p>
        </div>
        <Button onClick={loadRuntimes} variant="outline">
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Provision Form */}
        <Card>
          <CardHeader>
            <CardTitle>Provision New Runtime</CardTitle>
            <CardDescription>
              Create a temporary remote runtime for testing
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="workspace">Workspace *</Label>
              <Select value={selectedWorkspaceId} onValueChange={setSelectedWorkspaceId}>
                <SelectTrigger>
                  <SelectValue placeholder="Select workspace" />
                </SelectTrigger>
                <SelectContent>
                  {workspaces.map((ws) => (
                    <SelectItem key={ws.id} value={ws.id}>
                      {ws.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="name">Name *</Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="my-test-runtime"
                data-test-id={TEST_IDS.ephemeral.runtimeName}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="region">Region *</Label>
                <Select value={region} onValueChange={setRegion}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="us-east-1">US East (N. Virginia)</SelectItem>
                    <SelectItem value="us-west-2">US West (Oregon)</SelectItem>
                    <SelectItem value="eu-west-1">Europe (Ireland)</SelectItem>
                    <SelectItem value="ap-southeast-1">Asia Pacific (Singapore)</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="instanceType">Instance Type *</Label>
                <Select value={instanceType} onValueChange={(v) => setInstanceType(v as any)}>
                  <SelectTrigger data-test-id={TEST_IDS.ephemeral.instanceType}>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="small">Small (1 CPU, 1GB RAM)</SelectItem>
                    <SelectItem value="medium">Medium (2 CPU, 4GB RAM)</SelectItem>
                    <SelectItem value="large">Large (4 CPU, 8GB RAM)</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="image">Image *</Label>
              <Input
                id="image"
                value={image}
                onChange={(e) => setImage(e.target.value)}
                placeholder="ubuntu:22.04"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="duration">Duration (hours)</Label>
                <Input
                  id="duration"
                  type="number"
                  value={duration}
                  onChange={(e) => setDuration(parseInt(e.target.value) || 24)}
                  min={1}
                  max={168}
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="members">Members (comma-separated user IDs)</Label>
              <Input
                id="members"
                value={members}
                onChange={(e) => setMembers(e.target.value)}
                placeholder="user1, user2, user3"
              />
            </div>

            <Button onClick={handleProvision} disabled={loading || !name || !image || !selectedWorkspaceId} className="w-full" data-test-id={TEST_IDS.ephemeral.provisionRuntime}>
              <Plus className="h-4 w-4 mr-2" />
              {loading ? 'Provisioning...' : 'Provision Runtime'}
            </Button>
          </CardContent>
        </Card>

        {/* Runtime List */}
        <Card>
          <CardHeader>
            <CardTitle>Active Runtimes ({runtimes.length})</CardTitle>
            <CardDescription>
              Manage your ephemeral runtimes
            </CardDescription>
          </CardHeader>
          <CardContent>
            {runtimes.length === 0 ? (
              <p className="text-muted-foreground text-sm text-center py-8">
                No runtimes found. Provision one to get started.
              </p>
            ) : (
              <div className="space-y-4">
                {runtimes.map((runtime) => (
                  <div key={runtime.id} className="p-4 border rounded-lg">
                    <div className="flex items-start justify-between mb-3">
                      <div>
                        <div className="flex items-center gap-2 mb-2">
                          <span className="font-semibold">{runtime.name}</span>
                          {getStatusBadge(runtime.status)}
                        </div>
                        <div className="text-sm text-muted-foreground space-y-1">
                          <div className="flex items-center gap-2">
                            <Clock className="h-3 w-3" />
                            Expires: {new Date(runtime.expiresAt).toLocaleString()}
                          </div>
                          <div className="flex items-center gap-2">
                            <Activity className="h-3 w-3" />
                            Last activity: {new Date(runtime.lastActivityAt).toLocaleString()}
                          </div>
                          <div className="flex items-center gap-2">
                            <Cloud className="h-3 w-3" />
                            {runtime.region} • {runtime.instanceType}
                          </div>
                          {runtime.resources && (
                            <div>
                              {runtime.resources.cpu} CPU • {runtime.resources.memory} RAM • {runtime.resources.disk} Disk
                            </div>
                          )}
                          {runtime.members.length > 0 && (
                            <div className="flex items-center gap-2">
                              <Users className="h-3 w-3" />
                              {runtime.members.length} member(s)
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      {runtime.status === 'active' && (
                        <>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleUpdateActivity(runtime.id)}
                          >
                            <Activity className="h-3 w-3 mr-1" />
                            Update Activity
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleSleep(runtime.id)}
                            data-test-id={TEST_IDS.ephemeral.sleepRuntime}
                          >
                            <Pause className="h-3 w-3 mr-1" />
                            Sleep
                          </Button>
                        </>
                      )}
                      {runtime.status === 'sleeping' && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleWake(runtime.id)}
                          data-test-id={TEST_IDS.ephemeral.wakeRuntime}
                        >
                          <Play className="h-3 w-3 mr-1" />
                          Wake
                        </Button>
                      )}
                      {runtime.status !== 'terminated' && runtime.status !== 'terminating' && (
                        <Button
                          size="sm"
                          variant="destructive"
                          onClick={() => handleTerminate(runtime.id)}
                          data-test-id={TEST_IDS.ephemeral.terminateRuntime}
                        >
                          <Power className="h-3 w-3 mr-1" />
                          Terminate
                        </Button>
                      )}
                    </div>
                    {runtime.accessUrl && (
                      <div className="mt-2 text-xs text-muted-foreground">
                        Access URL: <a href={runtime.accessUrl} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">{runtime.accessUrl}</a>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

