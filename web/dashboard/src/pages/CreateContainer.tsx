import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { api, RunContainerOptions } from '@/lib/api'
import { ArrowLeft } from 'lucide-react'
import { Link } from 'react-router-dom'
import { TEST_IDS } from '@/lib/test-ids'

export function CreateContainer() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<RunContainerOptions>({
    image: '',
    name: '',
    port: '',
    detach: true,
    env: [],
    volume: [],
  })
  const [networks, setNetworks] = useState<Array<{ id:string; name:string }>>([])
  const [teams, setTeams] = useState<Array<{ id:string; name:string }>>([])
  const [selectedWorkspace, setSelectedWorkspace] = useState<string>('')
  const [volumes, setVolumes] = useState<Array<{ host: string; container: string; ro: boolean }>>([
    { host: '', container: '', ro: false },
  ])
  const [ports, setPorts] = useState<Array<{ host: string; container: string; proto: 'tcp' | 'udp' }>>([
    { host: '', container: '', proto: 'tcp' },
  ])
  const [health, setHealth] = useState<{
    type: '' | 'http' | 'tcp' | 'cmd'
    httpPath: string
    httpPort: string
    tcpPort: string
    cmd: string
    interval: number
    timeout: number
    retries: number
    startPeriod: number
  }>({
    type: '',
    httpPath: '/health',
    httpPort: '',
    tcpPort: '',
    cmd: '',
    interval: 10,
    timeout: 2,
    retries: 3,
    startPeriod: 0,
  })

  useEffect(() => {
    api.listNetworks().then(list => {
      const simple = list.map(n=>({ id: n.id, name: n.name }))
      setNetworks(simple)
    }).catch(()=>{})
    api.listTeams().then(list => {
      const simple = list.map(t=>({ id: t.id, name: t.name }))
      setTeams(simple)
    }).catch(()=>{})
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      // convert volumes to API format: host:container[:ro]
      const volumeSpecs = volumes
        .filter(v => v.host.trim() && v.container.trim())
        .map(v => `${v.host.trim()}:${v.container.trim()}${v.ro ? ':ro' : ''}`)

      const portSpecs = ports
        .filter(p => p.host.trim() && p.container.trim())
        .map(p => `${p.host.trim()}:${p.container.trim()}${p.proto ? `/${p.proto}` : ''}`)

      await api.runContainer({
        ...formData,
        volume: volumeSpecs,
        ports: portSpecs,
        workspaceId: selectedWorkspace || undefined,
        // health mapping
        healthType: health.type || undefined,
        healthHttpPath: health.type === 'http' ? health.httpPath : undefined,
        healthHttpPort: health.type === 'http' ? health.httpPort : undefined,
        healthTcpPort: health.type === 'tcp' ? health.tcpPort : undefined,
        healthCmd: health.type === 'cmd' && health.cmd.trim() ? health.cmd.trim().split(/\s+/) : undefined,
        healthIntervalSec: health.type ? health.interval : undefined,
        healthTimeoutSec: health.type ? health.timeout : undefined,
        healthRetries: health.type ? health.retries : undefined,
        healthStartPeriodSec: health.type ? health.startPeriod : undefined,
      } as any)
      navigate('/containers')
    } catch (error) {
      console.error('Failed to create container:', error)
      alert('Failed to create container. Check console for details.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="p-8">
      <Link to="/containers">
        <Button variant="ghost" className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Containers
        </Button>
      </Link>

      <Card className="max-w-2xl">
        <CardHeader>
          <CardTitle>Create New Container</CardTitle>
          <CardDescription>
            Run a container from an image
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="image">Image *</Label>
              <Input
                id="image"
                placeholder="nginx:latest"
                value={formData.image}
                onChange={(e) =>
                  setFormData({ ...formData, image: e.target.value })
                }
                required
                data-test-id={TEST_IDS.createContainer.imageInput}
              />
              <p className="text-xs text-muted-foreground">
                Enter the image name (e.g., nginx, postgres:13, node:18)
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="name">Container Name</Label>
              <Input
                id="name"
                placeholder="my-container"
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                data-test-id={TEST_IDS.createContainer.nameInput}
              />
              <p className="text-xs text-muted-foreground">
                Optional: Assign a name to the container
              </p>
            </div>

            <div className="space-y-2">
              <Label>Workspace (optional)</Label>
              <select
                className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                value={selectedWorkspace}
                onChange={(e) => setSelectedWorkspace(e.target.value)}
              >
                <option value="">No workspace (personal)</option>
                {teams.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
              </select>
              <p className="text-xs text-muted-foreground">Associate this container with a team workspace for shared access.</p>
            </div>

            <div className="space-y-2">
              <Label>Service Name (optional)</Label>
              <Input
                placeholder="api"
                value={(formData as any).service || ''}
                onChange={(e)=>setFormData({ ...formData, service: e.target.value || undefined } as any)}
              />
              <p className="text-xs text-muted-foreground">Registers this container instance under the service name for discovery.</p>
            </div>

            <div className="space-y-2">
              <Label>Network</Label>
              <select
                className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                value={formData.network || ''}
                onChange={(e)=>setFormData({ ...formData, network: e.target.value || undefined })}
              >
                <option value="">Default</option>
                {networks.map(n=>(<option key={n.id} value={n.name}>{n.name}</option>))}
              </select>
              <p className="text-xs text-muted-foreground">Choose a custom network (bridge driver) or leave Default.</p>
            </div>

            <div className="space-y-3">
              <Label>Port Mappings</Label>
              <div className="space-y-2">
                {ports.map((p, idx) => (
                  <div key={idx} className="grid grid-cols-12 gap-2 items-end">
                    <div className="col-span-4">
                      <Label className="text-xs">Host Port</Label>
                      <Input
                        placeholder="8080"
                        value={p.host}
                        onChange={(e) => {
                          const next = [...ports]
                          next[idx].host = e.target.value
                          setPorts(next)
                        }}
                      />
                      <div className="mt-1 flex gap-2">
                        <Button type="button" variant="outline" onClick={async ()=>{
                          const from = Number(p.host) > 0 ? Number(p.host) : undefined
                          try { const s = await api.suggestPort(from); const next = [...ports]; next[idx].host = String(s.port); setPorts(next) } catch {}
                        }}>Auto-assign</Button>
                        <Button type="button" variant="outline" onClick={async ()=>{
                          const n = Number(p.host); if (!n) return
                          try { const list = await api.listPorts(); const used = (list.ports||[]).some(x=>x.port===n)
                            if (used) { alert(`Port ${n} is already reserved`) } else { alert(`Port ${n} appears available`) }
                          } catch {}
                        }}>Check</Button>
                      </div>
                    </div>
                    <div className="col-span-4">
                      <Label className="text-xs">Container Port</Label>
                      <Input
                        placeholder="80"
                        value={p.container}
                        onChange={(e) => {
                          const next = [...ports]
                          next[idx].container = e.target.value
                          setPorts(next)
                        }}
                      />
                    </div>
                    <div className="col-span-2">
                      <Label className="text-xs">Protocol</Label>
                      <select
                        className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                        value={p.proto}
                        onChange={(e) => {
                          const next = [...ports]
                          next[idx].proto = (e.target.value as 'tcp' | 'udp') || 'tcp'
                          setPorts(next)
                        }}
                      >
                        <option value="tcp">tcp</option>
                        <option value="udp">udp</option>
                      </select>
                    </div>
                    {ports.length > 1 && (
                      <div className="col-span-2 flex justify-end">
                        <Button
                          type="button"
                          variant="outline"
                          onClick={() => setPorts(ports.filter((_, i) => i !== idx))}
                        >
                          Remove
                        </Button>
                      </div>
                    )}
                  </div>
                ))}
              </div>
              <div>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => setPorts([...ports, { host: '', container: '', proto: 'tcp' }])}
                >
                  Add Port
                </Button>
              </div>
              <p className="text-xs text-muted-foreground">
                Each entry becomes hostPort:containerPort/proto. Backend currently records mappings; real NAT publishing is planned under networking tasks.
              </p>
            </div>

            <div className="space-y-3">
              <Label>Volumes</Label>
              <div className="space-y-2">
                {volumes.map((vol, idx) => (
                  <div key={idx} className="grid grid-cols-12 gap-2 items-end">
                    <div className="col-span-5">
                      <Label className="text-xs">Host Path</Label>
                      <Input
                        placeholder="/host/path"
                        value={vol.host}
                        onChange={(e) => {
                          const next = [...volumes]
                          next[idx].host = e.target.value
                          setVolumes(next)
                        }}
                      />
                    </div>
                    <div className="col-span-5">
                      <Label className="text-xs">Container Path</Label>
                      <Input
                        placeholder="/container/path"
                        value={vol.container}
                        onChange={(e) => {
                          const next = [...volumes]
                          next[idx].container = e.target.value
                          setVolumes(next)
                        }}
                      />
                    </div>
                    <div className="col-span-2 flex items-center space-x-2">
                      <input
                        id={`ro-${idx}`}
                        type="checkbox"
                        className="rounded border-gray-300"
                        checked={vol.ro}
                        onChange={(e) => {
                          const next = [...volumes]
                          next[idx].ro = e.target.checked
                          setVolumes(next)
                        }}
                      />
                      <Label htmlFor={`ro-${idx}`} className="text-sm">Read-only</Label>
                    </div>
                    {volumes.length > 1 && (
                      <div className="col-span-12 flex justify-end">
                        <Button
                          type="button"
                          variant="outline"
                          onClick={() => setVolumes(volumes.filter((_, i) => i !== idx))}
                        >
                          Remove
                        </Button>
                      </div>
                    )}
                  </div>
                ))}
              </div>
              <div>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => setVolumes([...volumes, { host: '', container: '', ro: false }])}
                >
                  Add Volume
                </Button>
              </div>
              <p className="text-xs text-muted-foreground">
                Bind-mount host directories/files into the container. Each entry becomes
                host:container or host:container:ro.
              </p>
            </div>

            <div className="space-y-3">
              <Label>Health Check (optional)</Label>
              <div className="grid grid-cols-12 gap-2 items-end">
                <div className="col-span-3">
                  <Label className="text-xs">Type</Label>
                  <select
                    className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                    value={health.type}
                    onChange={(e) => setHealth({ ...health, type: e.target.value as any })}
                  >
                    <option value="">None</option>
                    <option value="http">HTTP</option>
                    <option value="tcp">TCP</option>
                    <option value="cmd">Command</option>
                  </select>
                </div>
                {health.type === 'http' && (
                  <>
                    <div className="col-span-5">
                      <Label className="text-xs">HTTP Path</Label>
                      <Input
                        placeholder="/health"
                        value={health.httpPath}
                        onChange={(e) => setHealth({ ...health, httpPath: e.target.value })}
                      />
                    </div>
                    <div className="col-span-4">
                      <Label className="text-xs">HTTP Port</Label>
                      <Input
                        placeholder="80"
                        value={health.httpPort}
                        onChange={(e) => setHealth({ ...health, httpPort: e.target.value })}
                      />
                    </div>
                  </>
                )}
                {health.type === 'tcp' && (
                  <div className="col-span-4">
                    <Label className="text-xs">TCP Port</Label>
                    <Input
                      placeholder="80"
                      value={health.tcpPort}
                      onChange={(e) => setHealth({ ...health, tcpPort: e.target.value })}
                    />
                  </div>
                )}
                {health.type === 'cmd' && (
                  <div className="col-span-9">
                    <Label className="text-xs">Command</Label>
                    <Input
                      placeholder="curl -f http://localhost:80/health || exit 1"
                      value={health.cmd}
                      onChange={(e) => setHealth({ ...health, cmd: e.target.value })}
                    />
                  </div>
                )}
              </div>
              {health.type && (
                <div className="grid grid-cols-12 gap-2">
                  <div className="col-span-3">
                    <Label className="text-xs">Interval (s)</Label>
                    <Input
                      type="number"
                      min={1}
                      value={health.interval}
                      onChange={(e) => setHealth({ ...health, interval: Number(e.target.value) })}
                    />
                  </div>
                  <div className="col-span-3">
                    <Label className="text-xs">Timeout (s)</Label>
                    <Input
                      type="number"
                      min={1}
                      value={health.timeout}
                      onChange={(e) => setHealth({ ...health, timeout: Number(e.target.value) })}
                    />
                  </div>
                  <div className="col-span-3">
                    <Label className="text-xs">Retries</Label>
                    <Input
                      type="number"
                      min={0}
                      value={health.retries}
                      onChange={(e) => setHealth({ ...health, retries: Number(e.target.value) })}
                    />
                  </div>
                  <div className="col-span-3">
                    <Label className="text-xs">Start Period (s)</Label>
                    <Input
                      type="number"
                      min={0}
                      value={health.startPeriod}
                      onChange={(e) => setHealth({ ...health, startPeriod: Number(e.target.value) })}
                    />
                  </div>
                </div>
              )}
              <p className="text-xs text-muted-foreground">Configure liveness-style checks. Real probe execution will be enabled in real mode wiring.</p>
            </div>

            <div className="flex items-center space-x-2">
              <input
                type="checkbox"
                id="detach"
                checked={formData.detach}
                onChange={(e) =>
                  setFormData({ ...formData, detach: e.target.checked })
                }
                className="rounded border-gray-300"
              />
              <Label htmlFor="detach">Run in background (detached)</Label>
            </div>

            <div className="flex gap-2 pt-4">
              <Button type="submit" disabled={loading} data-test-id={TEST_IDS.createContainer.createButton}>
                {loading ? 'Creating...' : 'Create Container'}
              </Button>
              <Link to="/containers">
                <Button type="button" variant="outline" data-test-id={TEST_IDS.createContainer.cancelButton}>
                  Cancel
                </Button>
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}

