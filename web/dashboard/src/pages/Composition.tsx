import { useEffect, useState } from 'react'
import { Shuffle, Plus, Trash2, Eye, Save, Play, AlertTriangle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Checkbox } from '@/components/ui/checkbox'
import { api, Container, CompositionSpec, ComposedContainerSpec, SourceContainer, ContainerElements } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Composition() {
  const [containers, setContainers] = useState<Container[]>([])
  const [specs, setSpecs] = useState<CompositionSpec[]>([])
  const [activeTab, setActiveTab] = useState('builder')
  
  // Builder state
  const [specName, setSpecName] = useState('')
  const [specDescription, setSpecDescription] = useState('')
  const [sources, setSources] = useState<SourceContainer[]>([])
  const [preview, setPreview] = useState<ComposedContainerSpec | null>(null)
  const [previewLoading, setPreviewLoading] = useState(false)
  const [saving, setSaving] = useState(false)

  // Source selection
  const [selectedContainerId, setSelectedContainerId] = useState('')
  const [sourceElements, setSourceElements] = useState<ContainerElements>({})
  const [sourcePriority, setSourcePriority] = useState(1)
  const [sourceDescription, setSourceDescription] = useState('')

  useEffect(() => {
    loadContainers()
    loadSpecs()
  }, [])

  const loadContainers = async () => {
    try {
      const data = await api.listContainers(true)
      setContainers(data)
    } catch (error) {
      console.error('Failed to load containers:', error)
    }
  }

  const loadSpecs = async () => {
    try {
      const data = await api.listCompositionSpecs()
      setSpecs(data.specs)
    } catch (error) {
      console.error('Failed to load specs:', error)
    }
  }

  const handleAddSource = () => {
    if (!selectedContainerId) {
      alert('Please select a container')
      return
    }

    const newSource: SourceContainer = {
      containerId: selectedContainerId,
      elements: { ...sourceElements },
      priority: sourcePriority,
      description: sourceDescription || undefined,
    }

    setSources([...sources, newSource])
    setSelectedContainerId('')
    setSourceElements({})
    setSourcePriority(1)
    setSourceDescription('')
  }

  const handleRemoveSource = (index: number) => {
    setSources(sources.filter((_, i) => i !== index))
  }

  const handlePreview = async () => {
    if (!specName || sources.length === 0) {
      alert('Please provide a spec name and add at least one source')
      return
    }

    const spec: CompositionSpec = {
      name: specName,
      description: specDescription || undefined,
      sources,
    }

    setPreviewLoading(true)
    try {
      const composed = await api.previewComposition(spec)
      setPreview(composed)
    } catch (error: any) {
      console.error('Failed to preview:', error)
      alert(`Failed to preview composition: ${error.message}`)
    } finally {
      setPreviewLoading(false)
    }
  }

  const handleSave = async () => {
    if (!specName || sources.length === 0) {
      alert('Please provide a spec name and add at least one source')
      return
    }

    setSaving(true)
    try {
      const spec: CompositionSpec = {
        name: specName,
        description: specDescription || undefined,
        sources,
      }
      await api.createCompositionSpec(spec)
      loadSpecs()
      alert('Composition spec saved!')
    } catch (error: any) {
      console.error('Failed to save:', error)
      alert(`Failed to save composition spec: ${error.message}`)
    } finally {
      setSaving(false)
    }
  }

  const handleCompose = async (specName: string) => {
    const containerName = prompt('Enter container name (optional):')
    if (containerName === null) return

    try {
      const result = await api.composeContainerFromSpec(specName, containerName || undefined, true)
      alert(`Container ${result.container.id} created successfully!`)
      loadContainers()
    } catch (error: any) {
      console.error('Failed to compose:', error)
      alert(`Failed to create container: ${error.message}`)
    }
  }

  const handleDeleteSpec = async (name: string) => {
    if (!confirm(`Delete composition spec "${name}"?`)) return

    try {
      await api.deleteCompositionSpec(name)
      loadSpecs()
    } catch (error: any) {
      console.error('Failed to delete:', error)
      alert(`Failed to delete spec: ${error.message}`)
    }
  }

  const getContainerName = (id: string) => {
    const container = containers.find(c => c.id === id)
    return container?.name || container?.id || id
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Container Composition Builder</h1>
        <p className="text-muted-foreground">
          Create new containers by mixing elements from existing containers
        </p>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList>
          <TabsTrigger value="builder">
            <Shuffle className="h-4 w-4 mr-2" />
            Builder
          </TabsTrigger>
          <TabsTrigger value="specs">
            <Save className="h-4 w-4 mr-2" />
            Saved Specs
          </TabsTrigger>
        </TabsList>

        {/* Builder Tab */}
        <TabsContent value="builder" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Composition Specification</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Spec Name *</Label>
                  <Input
                    value={specName}
                    onChange={(e) => setSpecName(e.target.value)}
                    placeholder="my-composed-container"
                    data-test-id={TEST_IDS.composition.specName}
                  />
                </div>
                <div>
                  <Label>Description</Label>
                  <Input
                    value={specDescription}
                    onChange={(e) => setSpecDescription(e.target.value)}
                    placeholder="Description of this composition"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Add Source Container</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label>Container</Label>
                <select
                  value={selectedContainerId}
                  onChange={(e) => setSelectedContainerId(e.target.value)}
                  className="w-full px-3 py-2 border rounded-md"
                  data-test-id={TEST_IDS.composition.addSource}
                >
                  <option value="">Select a container...</option>
                  {containers.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.name || c.id} ({c.image}) - {c.status}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <Label>Priority (higher = wins conflicts)</Label>
                <Input
                  type="number"
                  min="1"
                  value={sourcePriority}
                  onChange={(e) => setSourcePriority(parseInt(e.target.value) || 1)}
                />
              </div>

              <div>
                <Label>Description (optional)</Label>
                <Input
                  value={sourceDescription}
                  onChange={(e) => setSourceDescription(e.target.value)}
                  placeholder="What to use from this container"
                />
              </div>

              <div>
                <Label className="mb-3 block">Elements to Extract</Label>
                <div className="grid grid-cols-3 gap-3">
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-image"
                      checked={sourceElements.image || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, image: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-image" className="text-sm">Image</label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-env"
                      checked={sourceElements.envVars || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, envVars: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-env" className="text-sm">Env Vars</label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-ports"
                      checked={sourceElements.ports || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, ports: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-ports" className="text-sm">Ports</label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-volumes"
                      checked={sourceElements.volumes || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, volumes: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-volumes" className="text-sm">Volumes</label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-network"
                      checked={sourceElements.network || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, network: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-network" className="text-sm">Network</label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-service"
                      checked={sourceElements.service || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, service: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-service" className="text-sm">Service</label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-health"
                      checked={sourceElements.healthCheck || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, healthCheck: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-health" className="text-sm">Health Check</label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-labels"
                      checked={sourceElements.labels || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, labels: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-labels" className="text-sm">Labels</label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="elem-command"
                      checked={sourceElements.command || false}
                      onCheckedChange={(checked) =>
                        setSourceElements({ ...sourceElements, command: checked as boolean })
                      }
                    />
                    <label htmlFor="elem-command" className="text-sm">Command</label>
                  </div>
                </div>
              </div>

              <Button onClick={handleAddSource} className="w-full" data-test-id={TEST_IDS.composition.addSource}>
                <Plus className="h-4 w-4 mr-2" />
                Add Source
              </Button>
            </CardContent>
          </Card>

          {sources.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Source Containers ({sources.length})</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {sources.map((source, index) => (
                    <div key={index} className="p-3 border rounded-lg flex items-start justify-between">
                      <div className="flex-1">
                        <div className="font-semibold">
                          {getContainerName(source.containerId)}
                          {source.priority && source.priority > 1 && (
                            <Badge variant="outline" className="ml-2">Priority {source.priority}</Badge>
                          )}
                        </div>
                        {source.description && (
                          <div className="text-sm text-muted-foreground mt-1">{source.description}</div>
                        )}
                        <div className="text-xs text-muted-foreground mt-2">
                          Elements: {Object.entries(source.elements)
                            .filter(([_, v]) => v === true)
                            .map(([k]) => k)
                            .join(', ') || 'none'}
                        </div>
                      </div>
                      <Button
                        size="sm"
                        variant="destructive"
                        onClick={() => handleRemoveSource(index)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          <div className="flex gap-2">
            <Button onClick={handlePreview} disabled={previewLoading || sources.length === 0} data-test-id={TEST_IDS.composition.previewComposition}>
              <Eye className="h-4 w-4 mr-2" />
              {previewLoading ? 'Previewing...' : 'Preview Composition'}
            </Button>
            <Button onClick={handleSave} disabled={saving || sources.length === 0} variant="outline" data-test-id={TEST_IDS.composition.saveSpec}>
              <Save className="h-4 w-4 mr-2" />
              {saving ? 'Saving...' : 'Save Spec'}
            </Button>
          </div>

          {preview && (
            <Card>
              <CardHeader>
                <CardTitle>Preview: {preview.name}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label className="font-semibold">Image</Label>
                  <div className="p-2 bg-gray-50 dark:bg-gray-900 rounded">{preview.image}</div>
                </div>

                {Object.keys(preview.envVars).length > 0 && (
                  <div>
                    <Label className="font-semibold">Environment Variables</Label>
                    <div className="p-2 bg-gray-50 dark:bg-gray-900 rounded max-h-32 overflow-y-auto">
                      {Object.entries(preview.envVars).map(([k, v]) => (
                        <div key={k} className="text-sm font-mono">
                          {k}={v}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {Object.keys(preview.ports).length > 0 && (
                  <div>
                    <Label className="font-semibold">Ports</Label>
                    <div className="p-2 bg-gray-50 dark:bg-gray-900 rounded">
                      {Object.entries(preview.ports).map(([container, host]) => (
                        <div key={container} className="text-sm">
                          {host}:{container}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {preview.network && (
                  <div>
                    <Label className="font-semibold">Network</Label>
                    <div className="p-2 bg-gray-50 dark:bg-gray-900 rounded">{preview.network}</div>
                  </div>
                )}

                {preview.service && (
                  <div>
                    <Label className="font-semibold">Service</Label>
                    <div className="p-2 bg-gray-50 dark:bg-gray-900 rounded">{preview.service}</div>
                  </div>
                )}

                {preview.conflicts && preview.conflicts.length > 0 && (
                  <div>
                    <Label className="font-semibold text-orange-600 dark:text-orange-400 flex items-center gap-2">
                      <AlertTriangle className="h-4 w-4" />
                      Conflicts Resolved ({preview.conflicts.length})
                    </Label>
                    <div className="space-y-2 mt-2">
                      {preview.conflicts.map((conflict, idx) => (
                        <div key={idx} className="p-2 bg-orange-50 dark:bg-orange-900/20 rounded text-sm">
                          <div className="font-semibold">{conflict.type.toUpperCase()}</div>
                          <div className="text-muted-foreground">
                            {conflict.message || `${conflict.source1} vs ${conflict.source2}`}
                          </div>
                          <div className="text-xs mt-1">
                            Resolution: {conflict.resolution} → {String(conflict.finalValue)}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                <div>
                  <Label className="font-semibold">Sources Used</Label>
                  <div className="flex flex-wrap gap-1 mt-2">
                    {preview.sources.map((sourceId) => (
                      <Badge key={sourceId} variant="secondary">
                        {getContainerName(sourceId)}
                      </Badge>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* Saved Specs Tab */}
        <TabsContent value="specs" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Saved Composition Specs ({specs.length})</CardTitle>
            </CardHeader>
            <CardContent>
              {specs.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Save className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>No saved composition specs. Create one in the Builder tab.</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {specs.map((spec) => (
                    <div key={spec.name} className="p-4 border rounded-lg">
                      <div className="flex items-start justify-between mb-2">
                        <div className="flex-1">
                          <div className="font-semibold text-lg">{spec.name}</div>
                          {spec.description && (
                            <div className="text-sm text-muted-foreground mt-1">{spec.description}</div>
                          )}
                          <div className="text-xs text-muted-foreground mt-2">
                            {spec.sources.length} source(s) | Created: {spec.createdAt ? new Date(spec.createdAt).toLocaleString() : 'N/A'}
                          </div>
                        </div>
                        <div className="flex gap-2">
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleCompose(spec.name)}
                          >
                            <Play className="h-4 w-4 mr-1" />
                            Compose
                          </Button>
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => handleDeleteSpec(spec.name)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                      <div className="mt-2">
                        <div className="text-sm font-semibold mb-1">Sources:</div>
                        <div className="space-y-1">
                          {spec.sources.map((source, idx) => (
                            <div key={idx} className="text-xs text-muted-foreground pl-2">
                              • {getContainerName(source.containerId)} (Priority: {source.priority || 1})
                            </div>
                          ))}
                        </div>
                      </div>
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

