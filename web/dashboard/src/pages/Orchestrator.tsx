import { useEffect, useState } from 'react'
import { Server, Zap, Plus, Trash2, RefreshCw, Square } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { api, NodeInfo, DeploymentInfo, CreateDeploymentRequest } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Orchestrator() {
  const [nodes, setNodes] = useState<NodeInfo[]>([])
  const [deployments, setDeployments] = useState<DeploymentInfo[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('nodes')

  // Node registration form
  const [nodeId, setNodeId] = useState('')
  const [nodeName, setNodeName] = useState('')
  const [nodeAddress, setNodeAddress] = useState('')
  const [nodePort, setNodePort] = useState(8081)
  const [registering, setRegistering] = useState(false)

  // Deployment creation form
  const [deployName, setDeployName] = useState('')
  const [deployImage, setDeployImage] = useState('')
  const [deployTag, setDeployTag] = useState('latest')
  const [deployReplicas, setDeployReplicas] = useState(1)
  const [deployStrategy, setDeployStrategy] = useState('rolling')
  const [deployAutoRestart, setDeployAutoRestart] = useState(true)
  const [deployServiceName, setDeployServiceName] = useState('')
  const [deployNetworkName, setDeployNetworkName] = useState('')
  const [deployPorts, setDeployPorts] = useState('')
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    loadData()
    const interval = setInterval(loadData, 5000) // Refresh every 5 seconds
    return () => clearInterval(interval)
  }, [])

  const loadData = async () => {
    try {
      const [nodesData, deploymentsData] = await Promise.all([
        api.listNodes(),
        api.listDeployments(),
      ])
      setNodes(nodesData.nodes)
      setDeployments(deploymentsData.deployments)
    } catch (error) {
      console.error('Failed to load orchestrator data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleRegisterNode = async () => {
    if (!nodeId.trim() || !nodeAddress.trim()) {
      alert('Please fill in required fields')
      return
    }

    setRegistering(true)
    try {
      await api.registerNode({
        id: nodeId,
        name: nodeName || nodeId,
        address: nodeAddress,
        port: nodePort,
      })
      setNodeId('')
      setNodeName('')
      setNodeAddress('')
      setNodePort(8081)
      loadData()
    } catch (error: any) {
      console.error('Failed to register node:', error)
      alert(`Failed to register node: ${error.message}`)
    } finally {
      setRegistering(false)
    }
  }

  const handleCreateDeployment = async () => {
    if (!deployName.trim() || !deployImage.trim()) {
      alert('Please fill in required fields')
      return
    }

    setCreating(true)
    try {
      const deploy: CreateDeploymentRequest = {
        name: deployName,
        image: deployImage,
        tag: deployTag,
        replicas: deployReplicas,
        strategy: deployStrategy,
        autoRestart: deployAutoRestart,
        serviceName: deployServiceName || undefined,
        networkName: deployNetworkName || undefined,
        ports: deployPorts ? deployPorts.split(',').map(p => p.trim()).filter(p => p) : undefined,
      }
      await api.createDeployment(deploy)
      setDeployName('')
      setDeployImage('')
      setDeployTag('latest')
      setDeployReplicas(1)
      setDeployStrategy('rolling')
      setDeployAutoRestart(true)
      setDeployServiceName('')
      setDeployNetworkName('')
      setDeployPorts('')
      setActiveTab('deployments')
      loadData()
    } catch (error: any) {
      console.error('Failed to create deployment:', error)
      alert(`Failed to create deployment: ${error.message}`)
    } finally {
      setCreating(false)
    }
  }

  const handleScaleDeployment = async (id: string, replicas: number) => {
    try {
      await api.scaleDeployment(id, replicas)
      loadData()
    } catch (error: any) {
      console.error('Failed to scale deployment:', error)
      alert(`Failed to scale: ${error.message}`)
    }
  }

  const handleRestartDeployment = async (id: string) => {
    try {
      await api.restartDeployment(id)
      loadData()
    } catch (error: any) {
      console.error('Failed to restart deployment:', error)
      alert(`Failed to restart: ${error.message}`)
    }
  }

  const handleDeleteDeployment = async (id: string) => {
    if (!confirm('Are you sure you want to delete this deployment?')) {
      return
    }

    try {
      await api.deleteDeployment(id)
      loadData()
    } catch (error: any) {
      console.error('Failed to delete deployment:', error)
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

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'online':
      case 'running':
      case 'healthy':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
      case 'offline':
      case 'failed':
      case 'unhealthy':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
      case 'deploying':
      case 'pending':
      case 'starting':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
    }
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Orchestrator</h1>
        <p className="text-muted-foreground">
          Manage multi-node deployments and cluster nodes
        </p>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList>
          <TabsTrigger value="nodes">
            <Server className="h-4 w-4 mr-2" />
            Nodes
          </TabsTrigger>
          <TabsTrigger value="deployments">
            <Zap className="h-4 w-4 mr-2" />
            Deployments
          </TabsTrigger>
        </TabsList>

        {/* Nodes Tab */}
        <TabsContent value="nodes" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Server className="h-5 w-5" />
                Register New Node
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Node ID *</Label>
                  <Input
                    value={nodeId}
                    onChange={(e) => setNodeId(e.target.value)}
                    placeholder="node-001"
                    data-test-id={TEST_IDS.orchestrator.nodeId}
                  />
                </div>
                <div>
                  <Label>Node Name</Label>
                  <Input
                    value={nodeName}
                    onChange={(e) => setNodeName(e.target.value)}
                    placeholder="Worker Node 1"
                    data-test-id={TEST_IDS.orchestrator.nodeName}
                  />
                </div>
                <div>
                  <Label>Address *</Label>
                  <Input
                    value={nodeAddress}
                    onChange={(e) => setNodeAddress(e.target.value)}
                    placeholder="192.168.1.100"
                    data-test-id={TEST_IDS.orchestrator.nodeAddress}
                  />
                </div>
                <div>
                  <Label>Port</Label>
                  <Input
                    type="number"
                    value={nodePort}
                    onChange={(e) => setNodePort(parseInt(e.target.value) || 8081)}
                  />
                </div>
              </div>
              <Button onClick={handleRegisterNode} disabled={registering} className="mt-4" data-test-id={TEST_IDS.orchestrator.registerNode}>
                <Plus className="h-4 w-4 mr-2" />
                {registering ? 'Registering...' : 'Register Node'}
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Cluster Nodes ({nodes.length})</CardTitle>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-muted-foreground">Loading nodes...</div>
              ) : nodes.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Server className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>No nodes registered. Register a node to get started.</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {nodes.map((node) => (
                    <div key={node.id} className="p-4 border rounded-lg">
                      <div className="flex items-start justify-between mb-2">
                        <div>
                          <div className="flex items-center gap-2 mb-1">
                            <span className="font-semibold text-lg">{node.name}</span>
                            <Badge className={getStatusColor(node.status)}>
                              {node.status}
                            </Badge>
                          </div>
                          <div className="text-sm text-muted-foreground">
                            ID: {node.id} | {node.address}:{node.port}
                          </div>
                          {node.region && (
                            <div className="text-xs text-muted-foreground">
                              {node.region}{node.zone && ` / ${node.zone}`}
                            </div>
                          )}
                        </div>
                      </div>
                      <div className="grid grid-cols-4 gap-4 mt-4 text-sm">
                        <div>
                          <div className="text-muted-foreground">CPU</div>
                          <div className="font-semibold">
                            {node.resources.cpuCores} cores ({node.resources.cpuUsed.toFixed(1)}% used)
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Memory</div>
                          <div className="font-semibold">
                            {node.resources.memoryMB} MB ({node.resources.memoryUsed.toFixed(1)}% used)
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Disk</div>
                          <div className="font-semibold">
                            {node.resources.diskGB} GB ({node.resources.diskUsed.toFixed(1)}% used)
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Containers</div>
                          <div className="font-semibold">{node.containerCount} running</div>
                        </div>
                      </div>
                      <div className="text-xs text-muted-foreground mt-2">
                        Last seen: {formatDate(node.lastSeen)}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Deployments Tab */}
        <TabsContent value="deployments" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Zap className="h-5 w-5" />
                Create New Deployment
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Name *</Label>
                  <Input
                    value={deployName}
                    onChange={(e) => setDeployName(e.target.value)}
                    placeholder="my-app"
                    data-test-id={TEST_IDS.orchestrator.deployName}
                  />
                </div>
                <div>
                  <Label>Image *</Label>
                  <Input
                    value={deployImage}
                    onChange={(e) => setDeployImage(e.target.value)}
                    placeholder="nginx"
                    data-test-id={TEST_IDS.orchestrator.deployImage}
                  />
                </div>
                <div>
                  <Label>Tag</Label>
                  <Input
                    value={deployTag}
                    onChange={(e) => setDeployTag(e.target.value)}
                    placeholder="latest"
                  />
                </div>
                <div>
                  <Label>Replicas</Label>
                  <Input
                    type="number"
                    min="1"
                    value={deployReplicas}
                    onChange={(e) => setDeployReplicas(parseInt(e.target.value) || 1)}
                    data-test-id={TEST_IDS.orchestrator.deployReplicas}
                  />
                </div>
                <div>
                  <Label>Strategy</Label>
                  <select
                    className="w-full p-2 border rounded"
                    value={deployStrategy}
                    onChange={(e) => setDeployStrategy(e.target.value)}
                  >
                    <option value="rolling">Rolling</option>
                    <option value="recreate">Recreate</option>
                    <option value="canary">Canary</option>
                  </select>
                </div>
                <div className="flex items-center gap-2 pt-6">
                  <input
                    type="checkbox"
                    id="autoRestart"
                    checked={deployAutoRestart}
                    onChange={(e) => setDeployAutoRestart(e.target.checked)}
                  />
                  <Label htmlFor="autoRestart">Auto-restart</Label>
                </div>
                <div>
                  <Label>Service Name (for auto-registration)</Label>
                  <Input
                    value={deployServiceName}
                    onChange={(e) => setDeployServiceName(e.target.value)}
                    placeholder="my-api-service"
                  />
                </div>
                <div>
                  <Label>Network Name (auto-create if missing)</Label>
                  <Input
                    value={deployNetworkName}
                    onChange={(e) => setDeployNetworkName(e.target.value)}
                    placeholder="deployment-network"
                  />
                </div>
                <div>
                  <Label>Ports (comma-separated, e.g. "8080:80,8443:443")</Label>
                  <Input
                    value={deployPorts}
                    onChange={(e) => setDeployPorts(e.target.value)}
                    placeholder="8080:80"
                  />
                </div>
              </div>
              <Button onClick={handleCreateDeployment} disabled={creating} className="mt-4" data-test-id={TEST_IDS.orchestrator.createDeployment}>
                <Plus className="h-4 w-4 mr-2" />
                {creating ? 'Creating...' : 'Create Deployment'}
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Deployments ({deployments.length})</CardTitle>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-muted-foreground">Loading deployments...</div>
              ) : deployments.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Zap className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>No deployments. Create one to get started.</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {deployments.map((deploy) => {
                    const healthyCount = deploy.instances.filter(
                      i => i.health === 'healthy' && i.status === 'running'
                    ).length
                    
                    return (
                      <div key={deploy.id} className="p-4 border rounded-lg">
                        <div className="flex items-start justify-between mb-2">
                          <div>
                            <div className="flex items-center gap-2 mb-1">
                              <span className="font-semibold text-lg">{deploy.name}</span>
                              <Badge className={getStatusColor(deploy.status)}>
                                {deploy.status}
                              </Badge>
                              <Badge variant="outline">
                                {healthyCount}/{deploy.replicas} healthy
                              </Badge>
                            </div>
                            <div className="text-sm text-muted-foreground">
                              {deploy.image}:{deploy.tag} | Strategy: {deploy.strategy}
                              {deploy.serviceName && ` | Service: ${deploy.serviceName}`}
                              {deploy.networkName && ` | Network: ${deploy.networkName}`}
                            </div>
                            {deploy.autoRestart && (
                              <div className="text-xs text-muted-foreground mt-1">
                                Auto-restart enabled
                              </div>
                            )}
                          </div>
                          <div className="flex gap-2">
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleScaleDeployment(deploy.id, deploy.replicas + 1)}
                              data-test-id={TEST_IDS.orchestrator.scaleDeployment}
                            >
                              <Plus className="h-4 w-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleScaleDeployment(deploy.id, Math.max(0, deploy.replicas - 1))}
                              data-test-id={TEST_IDS.orchestrator.scaleDeployment}
                            >
                              <Square className="h-4 w-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleRestartDeployment(deploy.id)}
                              data-test-id={TEST_IDS.orchestrator.restartDeployment}
                            >
                              <RefreshCw className="h-4 w-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="destructive"
                              onClick={() => handleDeleteDeployment(deploy.id)}
                              data-test-id={TEST_IDS.orchestrator.deleteDeployment}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                        <div className="mt-4 space-y-2">
                          <div className="text-sm font-semibold">Instances:</div>
                          {deploy.instances.length === 0 ? (
                            <div className="text-sm text-muted-foreground">No instances</div>
                          ) : (
                            deploy.instances.map((instance) => (
                              <div key={instance.id} className="text-sm p-2 bg-muted rounded">
                                <div className="flex items-center gap-2">
                                  <span className="font-mono text-xs">{instance.id.substring(0, 8)}</span>
                                  <Badge className={getStatusColor(instance.status)}>
                                    {instance.status}
                                  </Badge>
                                  <Badge className={getStatusColor(instance.health)}>
                                    {instance.health}
                                  </Badge>
                                  <span className="text-muted-foreground">
                                    Node: {instance.nodeName || instance.nodeId}
                                  </span>
                                  {instance.restarts > 0 && (
                                    <span className="text-muted-foreground">
                                      ({instance.restarts} restarts)
                                    </span>
                                  )}
                                </div>
                              </div>
                            ))
                          )}
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

