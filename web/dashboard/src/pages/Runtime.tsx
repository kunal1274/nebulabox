import { useEffect, useState } from 'react'
import { Cpu, HardDrive, Activity, Plus, Play, Square, Trash2, Download } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { api, RuntimeContainer, RuntimeImage, RuntimeInfo, RuntimeContainerSpec } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Runtime() {
  const [containers, setContainers] = useState<RuntimeContainer[]>([])
  const [images, setImages] = useState<RuntimeImage[]>([])
  const [info, setInfo] = useState<RuntimeInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('containers')

  // Container creation form
  const [containerId, setContainerId] = useState('')
  const [containerName, setContainerName] = useState('')
  const [containerImage, setContainerImage] = useState('')
  const [creating, setCreating] = useState(false)

  // Image pull form
  const [pullImage, setPullImage] = useState('')
  const [pulling, setPulling] = useState(false)

  useEffect(() => {
    loadData()
    const interval = setInterval(loadData, 5000)
    return () => clearInterval(interval)
  }, [])

  const loadData = async () => {
    try {
      const [containersData, imagesData, infoData] = await Promise.all([
        api.listRuntimeContainers(),
        api.listRuntimeImages(),
        api.getRuntimeInfo(),
      ])
      setContainers(containersData.containers)
      setImages(imagesData.images)
      setInfo(infoData)
    } catch (error) {
      console.error('Failed to load runtime data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateContainer = async () => {
    if (!containerId.trim() || !containerImage.trim()) {
      alert('Please fill in required fields')
      return
    }

    setCreating(true)
    try {
      const spec: RuntimeContainerSpec = {
        id: containerId,
        name: containerName || containerId,
        image: containerImage,
      }
      await api.createRuntimeContainer(spec)
      setContainerId('')
      setContainerName('')
      setContainerImage('')
      loadData()
    } catch (error: any) {
      console.error('Failed to create container:', error)
      alert(`Failed to create container: ${error.message}`)
    } finally {
      setCreating(false)
    }
  }

  const handlePullImage = async () => {
    if (!pullImage.trim()) {
      alert('Please enter an image name')
      return
    }

    setPulling(true)
    try {
      await api.pullRuntimeImage(pullImage)
      setPullImage('')
      loadData()
    } catch (error: any) {
      console.error('Failed to pull image:', error)
      alert(`Failed to pull image: ${error.message}`)
    } finally {
      setPulling(false)
    }
  }

  const handleStartContainer = async (id: string) => {
    try {
      await api.startRuntimeContainer(id)
      loadData()
    } catch (error: any) {
      console.error('Failed to start container:', error)
      alert(`Failed to start: ${error.message}`)
    }
  }

  const handleStopContainer = async (id: string) => {
    try {
      await api.stopRuntimeContainer(id)
      loadData()
    } catch (error: any) {
      console.error('Failed to stop container:', error)
      alert(`Failed to stop: ${error.message}`)
    }
  }

  const handleDeleteContainer = async (id: string) => {
    if (!confirm('Are you sure you want to delete this container?')) {
      return
    }

    try {
      await api.deleteRuntimeContainer(id)
      loadData()
    } catch (error: any) {
      console.error('Failed to delete container:', error)
      alert(`Failed to delete: ${error.message}`)
    }
  }

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleString()
    } catch {
      return dateStr
    }
  }

  const formatBytes = (bytes: number) => {
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    let size = bytes
    let unitIndex = 0
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024
      unitIndex++
    }
    return `${size.toFixed(2)} ${units[unitIndex]}`
  }

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'running':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
      case 'stopped':
      case 'exited':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
      case 'created':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
      default:
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
    }
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Nebula Runtime</h1>
        <p className="text-muted-foreground">
          Custom container runtime with namespaces, cgroups, and OverlayFS
        </p>
      </div>

      {/* Runtime Info Card */}
      {info && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              Runtime Information
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-4 gap-4">
              <div>
                <div className="text-sm text-muted-foreground">Version</div>
                <div className="font-semibold">{info.version}</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Containers</div>
                <div className="font-semibold">{info.containers} running</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Images</div>
                <div className="font-semibold">{info.images} stored</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Memory</div>
                <div className="font-semibold">
                  {formatBytes(info.memoryUsed)} / {formatBytes(info.memoryTotal)}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList>
          <TabsTrigger value="containers">
            <Cpu className="h-4 w-4 mr-2" />
            Containers
          </TabsTrigger>
          <TabsTrigger value="images">
            <HardDrive className="h-4 w-4 mr-2" />
            Images
          </TabsTrigger>
        </TabsList>

        {/* Containers Tab */}
        <TabsContent value="containers" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Plus className="h-5 w-5" />
                Create Container
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Container ID *</Label>
                  <Input
                    value={containerId}
                    onChange={(e) => setContainerId(e.target.value)}
                    placeholder="my-container"
                    data-test-id={TEST_IDS.runtime.containerId}
                  />
                </div>
                <div>
                  <Label>Name</Label>
                  <Input
                    value={containerName}
                    onChange={(e) => setContainerName(e.target.value)}
                    placeholder="My Container"
                  />
                </div>
                <div className="col-span-2">
                  <Label>Image *</Label>
                  <Input
                    value={containerImage}
                    onChange={(e) => setContainerImage(e.target.value)}
                    placeholder="nginx:latest"
                    data-test-id={TEST_IDS.runtime.containerImage}
                  />
                </div>
              </div>
              <Button onClick={handleCreateContainer} disabled={creating} className="mt-4" data-test-id={TEST_IDS.runtime.createContainer}>
                <Plus className="h-4 w-4 mr-2" />
                {creating ? 'Creating...' : 'Create Container'}
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Runtime Containers ({containers.length})</CardTitle>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-muted-foreground">Loading containers...</div>
              ) : containers.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Cpu className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>No containers created. Create one to get started.</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {containers.map((container) => (
                    <div key={container.id} className="p-4 border rounded-lg">
                      <div className="flex items-start justify-between mb-2">
                        <div>
                          <div className="flex items-center gap-2 mb-1">
                            <span className="font-semibold text-lg">{container.name || container.id}</span>
                            <Badge className={getStatusColor(container.status)}>
                              {container.status}
                            </Badge>
                            {container.pid && (
                              <Badge variant="outline">PID: {container.pid}</Badge>
                            )}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            ID: {container.id} | Image: {container.image}
                          </div>
                          <div className="text-xs text-muted-foreground mt-1">
                            Created: {formatDate(container.createdAt)}
                            {container.startedAt && ` | Started: ${formatDate(container.startedAt)}`}
                          </div>
                        </div>
                        <div className="flex gap-2">
                          {container.status !== 'running' && (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleStartContainer(container.id)}
                              data-test-id={TEST_IDS.runtime.startContainer}
                            >
                              <Play className="h-4 w-4" />
                            </Button>
                          )}
                          {container.status === 'running' && (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleStopContainer(container.id)}
                              data-test-id={TEST_IDS.runtime.stopContainer}
                            >
                              <Square className="h-4 w-4" />
                            </Button>
                          )}
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => handleDeleteContainer(container.id)}
                            data-test-id={TEST_IDS.runtime.deleteContainer}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                      {container.stats && (
                        <div className="mt-4 grid grid-cols-4 gap-4 text-sm">
                          <div>
                            <div className="text-muted-foreground">CPU</div>
                            <div className="font-semibold">{container.stats.cpuPercent.toFixed(1)}%</div>
                          </div>
                          <div>
                            <div className="text-muted-foreground">Memory</div>
                            <div className="font-semibold">
                              {formatBytes(container.stats.memoryUsage)}
                            </div>
                          </div>
                          <div>
                            <div className="text-muted-foreground">Network</div>
                            <div className="font-semibold">
                              {formatBytes(container.stats.networkRx)} ↓ / {formatBytes(container.stats.networkTx)} ↑
                            </div>
                          </div>
                          <div>
                            <div className="text-muted-foreground">PIDs</div>
                            <div className="font-semibold">{container.stats.pidsCurrent}</div>
                          </div>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Images Tab */}
        <TabsContent value="images" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Download className="h-5 w-5" />
                Pull Image
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex gap-2">
                <Input
                  placeholder="nginx:latest"
                  value={pullImage}
                  onChange={(e) => setPullImage(e.target.value)}
                  data-test-id={TEST_IDS.runtime.imageName}
                />
                <Button onClick={handlePullImage} disabled={pulling} data-test-id={TEST_IDS.runtime.pullImage}>
                  <Download className="h-4 w-4 mr-2" />
                  {pulling ? 'Pulling...' : 'Pull Image'}
                </Button>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Runtime Images ({images.length})</CardTitle>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-muted-foreground">Loading images...</div>
              ) : images.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <HardDrive className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>No images. Pull an image to get started.</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {images.map((image) => (
                    <div key={image.id} className="p-4 border rounded-lg">
                      <div className="flex items-start justify-between mb-2">
                        <div>
                          <div className="font-semibold text-lg">
                            {image.name}:{image.tag}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            ID: {image.id} | Size: {formatBytes(image.size)}
                          </div>
                          {image.digest && (
                            <div className="text-xs text-muted-foreground font-mono">
                              Digest: {image.digest.substring(0, 32)}...
                            </div>
                          )}
                          <div className="text-xs text-muted-foreground mt-1">
                            Created: {formatDate(image.createdAt)}
                          </div>
                        </div>
                      </div>
                      {image.layers && image.layers.length > 0 && (
                        <details className="mt-2">
                          <summary className="cursor-pointer text-sm text-muted-foreground">
                            {image.layers.length} layers
                          </summary>
                          <div className="mt-2 space-y-1">
                            {image.layers.map((layer, idx) => (
                              <div key={idx} className="text-xs text-muted-foreground pl-4">
                                Layer {idx + 1}: {formatBytes(layer.size)} ({layer.digest.substring(0, 16)}...)
                              </div>
                            ))}
                          </div>
                        </details>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

