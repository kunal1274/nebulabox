import { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { api } from '@/lib/api'
import type { Snapshot, Container } from '@/lib/api'
import { Camera, Trash2, RefreshCw, Play, Package, Globe, HardDrive, Clock, User, FileText } from 'lucide-react'
import { TEST_IDS } from '@/lib/test-ids'

export function Snapshots() {
  const [snapshots, setSnapshots] = useState<Snapshot[]>([])
  const [containers, setContainers] = useState<Container[]>([])
  const [selectedSnapshot, setSelectedSnapshot] = useState<Snapshot | null>(null)
  const [filterType, setFilterType] = useState<string>('')
  const [filterResourceId, setFilterResourceId] = useState<string>('')

  // Create snapshot form
  const [snapshotName, setSnapshotName] = useState('')
  const [snapshotDescription, setSnapshotDescription] = useState('')
  const [snapshotType, setSnapshotType] = useState<'container' | 'workspace' | 'volume'>('container')
  const [selectedContainerId, setSelectedContainerId] = useState('')
  const [creating, setCreating] = useState(false)

  // Restore form
  const [restoreName, setRestoreName] = useState('')
  const [restoring, setRestoring] = useState(false)

  useEffect(() => {
    loadSnapshots()
    loadContainers()
  }, [filterType, filterResourceId])

  const loadSnapshots = async () => {
    try {
      const result = await api.listSnapshots(
        filterResourceId || undefined,
        filterType || undefined
      )
      setSnapshots(result.snapshots)
    } catch (error) {
      console.error('Failed to load snapshots:', error)
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

  const handleCreateSnapshot = async () => {
    if (!snapshotName) {
      alert('Please provide a snapshot name')
      return
    }

    if (snapshotType === 'container' && !selectedContainerId) {
      alert('Please select a container')
      return
    }

    setCreating(true)
    try {
      await api.createSnapshot({
        name: snapshotName,
        description: snapshotDescription || undefined,
        type: snapshotType,
        resourceId: selectedContainerId || 'workspace-placeholder',
      })
      setSnapshotName('')
      setSnapshotDescription('')
      setSelectedContainerId('')
      loadSnapshots()
      alert('Snapshot created successfully')
    } catch (error: any) {
      console.error('Failed to create snapshot:', error)
      alert(`Failed to create snapshot: ${error.message}`)
    } finally {
      setCreating(false)
    }
  }

  const handleDeleteSnapshot = async (snapshotId: string) => {
    if (!confirm('Are you sure you want to delete this snapshot?')) return

    try {
      await api.deleteSnapshot(snapshotId)
      loadSnapshots()
      if (selectedSnapshot?.id === snapshotId) {
        setSelectedSnapshot(null)
      }
      alert('Snapshot deleted successfully')
    } catch (error: any) {
      console.error('Failed to delete snapshot:', error)
      alert(`Failed to delete snapshot: ${error.message}`)
    }
  }

  const handleRestoreSnapshot = async () => {
    if (!selectedSnapshot) return

    setRestoring(true)
    try {
      const result = await api.restoreSnapshot(selectedSnapshot.id, restoreName || undefined)
      setRestoreName('')
      setSelectedSnapshot(null)
      loadSnapshots()
      loadContainers()
      alert(result.message || 'Snapshot restored successfully')
    } catch (error: any) {
      console.error('Failed to restore snapshot:', error)
      alert(`Failed to restore snapshot: ${error.message}`)
    } finally {
      setRestoring(false)
    }
  }

  const getSnapshotIcon = (type: string) => {
    switch (type) {
      case 'container':
        return <Package className="h-4 w-4" />
      case 'workspace':
        return <Globe className="h-4 w-4" />
      case 'volume':
        return <HardDrive className="h-4 w-4" />
      default:
        return <FileText className="h-4 w-4" />
    }
  }

  const getStateBadgeVariant = (state: string) => {
    switch (state) {
      case 'ready':
        return 'default'
      case 'creating':
      case 'restoring':
        return 'secondary'
      case 'failed':
        return 'destructive'
      default:
        return 'outline'
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  const formatSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`
  }

  return (
    <div className="container mx-auto p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Snapshots</h1>
        <Button onClick={loadSnapshots} variant="outline">
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh
        </Button>
      </div>

      <Tabs defaultValue="list" className="space-y-6">
        <TabsList>
          <TabsTrigger value="list">All Snapshots</TabsTrigger>
          <TabsTrigger value="create">Create Snapshot</TabsTrigger>
        </TabsList>

        <TabsContent value="list" className="space-y-6">
          {/* Filters */}
          <Card>
            <CardHeader>
              <CardTitle>Filters</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Type</Label>
                  <select
                    className="w-full p-2 border rounded-md mt-1"
                    value={filterType}
                    onChange={(e) => setFilterType(e.target.value)}
                  >
                    <option value="">All Types</option>
                    <option value="container">Container</option>
                    <option value="workspace">Workspace</option>
                    <option value="volume">Volume</option>
                  </select>
                </div>
                <div>
                  <Label>Resource ID</Label>
                  <Input
                    placeholder="Filter by resource ID"
                    value={filterResourceId}
                    onChange={(e) => setFilterResourceId(e.target.value)}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Snapshots List */}
          <Card>
            <CardHeader>
              <CardTitle>Snapshots ({snapshots.length})</CardTitle>
              <CardDescription>
                Manage and restore environment snapshots
              </CardDescription>
            </CardHeader>
            <CardContent>
              {snapshots.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Camera className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>No snapshots found</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {snapshots.map((snapshot) => (
                    <div
                      key={snapshot.id}
                      className={`p-4 border rounded-lg cursor-pointer hover:bg-accent ${
                        selectedSnapshot?.id === snapshot.id ? 'bg-accent' : ''
                      }`}
                      onClick={() => setSelectedSnapshot(snapshot)}
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            {getSnapshotIcon(snapshot.type)}
                            <span className="font-medium">{snapshot.name}</span>
                            <Badge variant={getStateBadgeVariant(snapshot.state)}>
                              {snapshot.state}
                            </Badge>
                            <Badge variant="outline">{snapshot.type}</Badge>
                          </div>
                          {snapshot.description && (
                            <p className="text-sm text-muted-foreground mb-2">
                              {snapshot.description}
                            </p>
                          )}
                          <div className="flex items-center gap-4 text-xs text-muted-foreground">
                            <span className="flex items-center gap-1">
                              <User className="h-3 w-3" />
                              {snapshot.createdBy}
                            </span>
                            <span className="flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              {formatDate(snapshot.createdAt)}
                            </span>
                            <span>Size: {formatSize(snapshot.size)}</span>
                            <span>Resource: {snapshot.resourceId}</span>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          {snapshot.state === 'ready' && (
                            <Button
                              size="sm"
                              onClick={(e) => {
                                e.stopPropagation()
                                setSelectedSnapshot(snapshot)
                              }}
                              data-test-id={TEST_IDS.snapshots.restoreSnapshot}
                            >
                              <Play className="h-4 w-4 mr-1" />
                              Restore
                            </Button>
                          )}
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={(e) => {
                              e.stopPropagation()
                              handleDeleteSnapshot(snapshot.id)
                            }}
                            data-test-id={TEST_IDS.snapshots.deleteSnapshot}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Snapshot Details */}
          {selectedSnapshot && (
            <Card>
              <CardHeader>
                <CardTitle>Snapshot Details</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label>Name</Label>
                    <p className="text-sm">{selectedSnapshot.name}</p>
                  </div>
                  <div>
                    <Label>Type</Label>
                    <p className="text-sm">{selectedSnapshot.type}</p>
                  </div>
                  <div>
                    <Label>State</Label>
                    <p className="text-sm">
                      <Badge variant={getStateBadgeVariant(selectedSnapshot.state)}>
                        {selectedSnapshot.state}
                      </Badge>
                    </p>
                  </div>
                  <div>
                    <Label>Resource ID</Label>
                    <p className="text-sm font-mono">{selectedSnapshot.resourceId}</p>
                  </div>
                  <div>
                    <Label>Size</Label>
                    <p className="text-sm">{formatSize(selectedSnapshot.size)}</p>
                  </div>
                  <div>
                    <Label>Created At</Label>
                    <p className="text-sm">{formatDate(selectedSnapshot.createdAt)}</p>
                  </div>
                  <div>
                    <Label>Created By</Label>
                    <p className="text-sm">{selectedSnapshot.createdBy}</p>
                  </div>
                </div>

                {selectedSnapshot.description && (
                  <div>
                    <Label>Description</Label>
                    <p className="text-sm">{selectedSnapshot.description}</p>
                  </div>
                )}

                {selectedSnapshot.type === 'container' && (
                  <>
                    {selectedSnapshot.image && (
                      <div>
                        <Label>Image</Label>
                        <p className="text-sm font-mono">{selectedSnapshot.image}</p>
                      </div>
                    )}
                    {selectedSnapshot.env && Object.keys(selectedSnapshot.env).length > 0 && (
                      <div>
                        <Label>Environment Variables</Label>
                        <div className="text-sm space-y-1 mt-1">
                          {Object.entries(selectedSnapshot.env).map(([key, value]) => (
                            <div key={key} className="font-mono">
                              {key}={value}
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                    {selectedSnapshot.ports && Object.keys(selectedSnapshot.ports).length > 0 && (
                      <div>
                        <Label>Ports</Label>
                        <div className="text-sm space-y-1 mt-1">
                          {Object.entries(selectedSnapshot.ports).map(([host, container]) => (
                            <div key={host} className="font-mono">
                              {host} → {container}
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  </>
                )}

                {selectedSnapshot.state === 'ready' && (
                  <div className="pt-4 border-t">
                    <Label>Restore Snapshot</Label>
                    <div className="flex gap-2 mt-2">
                      <Input
                        placeholder="New name (optional)"
                        value={restoreName}
                        onChange={(e) => setRestoreName(e.target.value)}
                        data-test-id={TEST_IDS.snapshots.snapshotName}
                      />
                      <Button
                        onClick={handleRestoreSnapshot}
                        disabled={restoring}
                        data-test-id={TEST_IDS.snapshots.restoreSnapshot}
                      >
                        <Play className="h-4 w-4 mr-2" />
                        {restoring ? 'Restoring...' : 'Restore'}
                      </Button>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="create" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Camera className="h-5 w-5" />
                Create Snapshot
              </CardTitle>
              <CardDescription>
                Capture the current state of a container, workspace, or volume
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label>Snapshot Name *</Label>
                <Input
                  value={snapshotName}
                  onChange={(e) => setSnapshotName(e.target.value)}
                  placeholder="my-snapshot"
                  data-test-id={TEST_IDS.snapshots.snapshotName}
                />
              </div>
              <div>
                <Label>Description</Label>
                <Input
                  value={snapshotDescription}
                  onChange={(e) => setSnapshotDescription(e.target.value)}
                  placeholder="Snapshot description"
                />
              </div>
              <div>
                <Label>Type *</Label>
                <select
                  className="w-full p-2 border rounded-md mt-1"
                  value={snapshotType}
                  onChange={(e) => setSnapshotType(e.target.value as 'container' | 'workspace' | 'volume')}
                  data-test-id={TEST_IDS.snapshots.resourceType}
                >
                  <option value="container">Container</option>
                  <option value="workspace">Workspace</option>
                  <option value="volume">Volume</option>
                </select>
              </div>
              {snapshotType === 'container' && (
                <div>
                  <Label>Container *</Label>
                  <select
                    className="w-full p-2 border rounded-md mt-1"
                    value={selectedContainerId}
                    onChange={(e) => setSelectedContainerId(e.target.value)}
                    data-test-id={TEST_IDS.snapshots.resourceType}
                  >
                    <option value="">Select a container...</option>
                    {containers.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name || c.id} ({c.image})
                      </option>
                    ))}
                  </select>
                </div>
              )}
              <Button
                onClick={handleCreateSnapshot}
                disabled={creating || !snapshotName || (snapshotType === 'container' && !selectedContainerId)}
                data-test-id={TEST_IDS.snapshots.createSnapshot}
              >
                <Camera className="h-4 w-4 mr-2" />
                {creating ? 'Creating...' : 'Create Snapshot'}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

