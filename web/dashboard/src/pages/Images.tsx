import { useEffect, useState } from 'react'
import { Download, Upload, RefreshCw, Tags, Trash2, Copy } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { api, Image, RegistryTags, RegistryCatalog, ImageScanResult } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Images() {
  const [images, setImages] = useState<Image[]>([])
  const [loading, setLoading] = useState(true)
  const [pullImage, setPullImage] = useState('')
  const [pullLoading, setPullLoading] = useState(false)
  const [repo, setRepo] = useState('nebulabox/nginx')
  const [tags, setTags] = useState<RegistryTags | null>(null)
  const [retagFrom, setRetagFrom] = useState('')
  const [retagTo, setRetagTo] = useState('')
  const [regLoading, setRegLoading] = useState(false)
  const [catalog, setCatalog] = useState<RegistryCatalog | null>(null)
  const [scans, setScans] = useState<Record<string, ImageScanResult>>({})
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('admin')
  const [authMsg, setAuthMsg] = useState('')
  const [dockerfile, setDockerfile] = useState('FROM alpine:3.19\nRUN echo "hello"\n')
  const [buildTag, setBuildTag] = useState('nebulabox/custom:latest')
  const [buildLogs, setBuildLogs] = useState<string[]>([])
  const [buildLoading, setBuildLoading] = useState(false)

  useEffect(() => {
    loadImages()
  }, [])

  const loadImages = async () => {
    setLoading(true)
    try {
      const data = await api.listImages()
      console.log('Loaded images from API:', data)
      setImages(data)
    } catch (error) {
      console.error('Failed to load images:', error)
      // Only set empty array on error, don't show mock data
      // Mock data makes it confusing - we want to see the actual error
      setImages([])
    } finally {
      setLoading(false)
    }
  }

  const handleScan = async (imageRef: string) => {
    try {
      const res = await api.scanImage(imageRef)
      setScans(prev => ({ ...prev, [imageRef]: res }))
    } catch (e) {
      console.error('Failed to scan image:', e)
      alert('Failed to scan image')
    }
  }

  const handleLogin = async () => {
    setAuthMsg('')
    try {
      const res = await api.registryLogin(username, password)
      setAuthMsg(`Token acquired (${res.token_type})`)
    } catch (e) {
      console.error('Login failed:', e)
      setAuthMsg('Login failed')
    }
  }

  const handleLogout = () => {
    api.registryLogout()
    setAuthMsg('Logged out')
  }

  const handleBuild = async () => {
    if (!buildTag.trim() || !dockerfile.trim()) return
    setBuildLoading(true)
    try {
      const res = await api.buildImageFromDockerfile(dockerfile, buildTag)
      setBuildLogs(res.logs)
    } catch (e) {
      console.error('Build failed:', e)
      setBuildLogs([`Build failed: ${String(e)}`])
    } finally {
      setBuildLoading(false)
    }
  }

  const loadTags = async () => {
    if (!repo.trim()) return
    setRegLoading(true)
    try {
      const data = await api.getRegistryTags(repo.trim())
      setTags(data)
    } catch (e) {
      console.error('Failed to load tags:', e)
      setTags(null)
    } finally {
      setRegLoading(false)
    }
  }

  const handleRetag = async () => {
    if (!repo || !retagFrom || !retagTo) return
    setRegLoading(true)
    try {
      await api.retagImage(repo, retagFrom, retagTo)
      setRetagTo('')
      loadTags()
    } catch (e) {
      console.error('Failed to retag:', e)
      alert('Failed to retag')
    } finally {
      setRegLoading(false)
    }
  }

  const handleDeleteTag = async (t: string) => {
    if (!confirm(`Delete tag ${repo}:${t}?`)) return
    setRegLoading(true)
    try {
      await api.deleteTag(repo, t)
      loadTags()
    } catch (e) {
      console.error('Failed to delete tag:', e)
      alert('Failed to delete tag')
    } finally {
      setRegLoading(false)
    }
  }

  const loadCatalog = async () => {
    setRegLoading(true)
    try {
      const data = await api.getRegistryCatalog()
      setCatalog(data)
    } catch (e) {
      console.error('Failed to load catalog:', e)
      setCatalog(null)
    } finally {
      setRegLoading(false)
    }
  }

  const handlePull = async () => {
    if (!pullImage.trim()) return

    setPullLoading(true)
    try {
      await api.pullImage(pullImage)
      setPullImage('')
      loadImages()
    } catch (error) {
      console.error('Failed to pull image:', error)
      alert('Failed to pull image')
    } finally {
      setPullLoading(false)
    }
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold mb-2">Images</h1>
          <p className="text-muted-foreground">
            Manage container images
          </p>
        </div>
        <Button onClick={loadImages} variant="outline" data-test-id={TEST_IDS.images.refresh}>
          <RefreshCw className="mr-2 h-4 w-4" />
          Refresh
        </Button>
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Pull Image</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2">
            <Input
              placeholder="nginx:latest"
              value={pullImage}
              onChange={(e) => setPullImage(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handlePull()}
            />
            <Button onClick={handlePull} disabled={pullLoading || !pullImage.trim()}>
              <Download className="mr-2 h-4 w-4" />
              {pullLoading ? 'Pulling...' : 'Pull'}
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Build Image (Dockerfile)</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-3">
            <div className="col-span-8">
              <Textarea rows={8} value={dockerfile} onChange={(e)=>setDockerfile(e.target.value)} className="font-mono" />
            </div>
            <div className="col-span-4 space-y-2">
              <Label>Tag</Label>
              <Input value={buildTag} onChange={(e)=>setBuildTag(e.target.value)} placeholder="repo/name:tag" />
              <Button onClick={handleBuild} disabled={buildLoading}>
                {buildLoading ? 'Building...' : 'Build'}
              </Button>
              <div className="text-xs text-muted-foreground">Paste a simple Dockerfile and set the tag. This uses a mock builder for now.</div>
            </div>
            {buildLogs.length > 0 && (
              <div className="col-span-12">
                <div className="border rounded p-2 bg-muted max-h-56 overflow-auto text-sm font-mono">
                  {buildLogs.map((l,i)=>(<div key={i}>{l}</div>))}
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Registry Catalog</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 mb-3 items-end">
            <div className="col-span-3">
              <Input placeholder="username" value={username} onChange={(e)=>setUsername(e.target.value)} />
            </div>
            <div className="col-span-3">
              <Input placeholder="password" type="password" value={password} onChange={(e)=>setPassword(e.target.value)} />
            </div>
            <div className="col-span-2">
              <Button variant="secondary" onClick={handleLogin}>Login</Button>
            </div>
            <div className="col-span-2">
              <Button variant="outline" onClick={handleLogout}>Logout</Button>
            </div>
            <div className="col-span-12 text-xs text-muted-foreground">{authMsg}</div>
          </div>
          <div className="flex gap-2 mb-3">
            <Button variant="outline" onClick={loadCatalog} disabled={regLoading}>
              <RefreshCw className="mr-2 h-4 w-4" />
              {regLoading ? 'Loading...' : 'Load Catalog'}
            </Button>
          </div>
          {catalog && (
            <div className="flex flex-wrap gap-2">
              {catalog.repositories?.length ? (
                catalog.repositories.map((r) => (
                  <Button
                    key={r}
                    variant={r === repo ? 'default' : 'secondary'}
                    size="sm"
                    onClick={() => { setRepo(r); setTags(null); setRetagFrom(''); setRetagTo(''); loadTags() }}
                  >
                    {r}
                  </Button>
                ))
              ) : (
                <div className="text-sm text-muted-foreground">No repositories</div>
              )}
            </div>
          )}
        </CardContent>
      </Card>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Registry Tags</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2 mb-3">
            <Input
              placeholder="repository (e.g., nebulabox/nginx)"
              value={repo}
              onChange={(e) => setRepo(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && loadTags()}
            />
            <Button variant="outline" onClick={loadTags} disabled={regLoading}>
              <Tags className="mr-2 h-4 w-4" />
              {regLoading ? 'Loading...' : 'Load Tags'}
            </Button>
          </div>
          {tags && (
            <div className="space-y-3">
              <div className="text-sm text-muted-foreground">{tags.name}</div>
              <div className="flex flex-wrap gap-2">
                {tags.tags?.length ? (
                  tags.tags.map((t) => (
                    <div key={t} className="flex items-center gap-2 border rounded px-2 py-1">
                      <code className="text-sm">{t}</code>
                      <Button size="sm" variant="ghost" onClick={() => setRetagFrom(t)} title="Use as source">
                        <Copy className="h-4 w-4" />
                      </Button>
                      <Button size="sm" variant="ghost" onClick={() => handleDeleteTag(t)} title="Delete tag">
                        <Trash2 className="h-4 w-4 text-red-600" />
                      </Button>
                    </div>
                  ))
                ) : (
                  <div className="text-sm text-muted-foreground">No tags</div>
                )}
              </div>
              <div className="grid grid-cols-12 gap-2 items-end">
                <div className="col-span-5">
                  <Input placeholder="source tag or digest" value={retagFrom} onChange={(e) => setRetagFrom(e.target.value)} />
                </div>
                <div className="col-span-5">
                  <Input placeholder="target tag" value={retagTo} onChange={(e) => setRetagTo(e.target.value)} />
                </div>
                <div className="col-span-2">
                  <Button onClick={handleRetag} disabled={!retagFrom || !retagTo || regLoading}>Retag</Button>
                </div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {loading ? (
        <div className="text-center py-12 text-muted-foreground">
          Loading images...
        </div>
      ) : images.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <p className="text-muted-foreground mb-4">No images found</p>
            <p className="text-sm text-muted-foreground">
              Pull an image to get started
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4" data-test-id={TEST_IDS.images.list}>
          {images.map((image) => (
            <Card key={image.id} data-test-id={TEST_IDS.images.card} data-image-name={image.name} data-image-tag={image.tag}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>
                      {image.name}:{image.tag}
                    </CardTitle>
                    <div className="mt-2 text-sm text-muted-foreground">
                      <div>Size: {image.size}</div>
                      <div>ID: {image.id}</div>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button variant="outline" size="sm" data-test-id={TEST_IDS.images.push}>
                      <Upload className="mr-2 h-4 w-4" />
                      Push
                    </Button>
                    <Button variant="secondary" size="sm" onClick={() => handleScan(`${image.name}:${image.tag}`)} data-test-id={TEST_IDS.images.scan}>
                      Scan
                    </Button>
                  </div>
                </div>
              </CardHeader>
              {scans[`${image.name}:${image.tag}`] && (
                <div className="px-6 pb-6">
                  <div className="text-sm mb-2">
                    Findings: 
                    <span className="ml-2">Critical {scans[`${image.name}:${image.tag}`].criticalCount}</span>
                    <span className="ml-2">High {scans[`${image.name}:${image.tag}`].highCount}</span>
                    <span className="ml-2">Medium {scans[`${image.name}:${image.tag}`].mediumCount}</span>
                    <span className="ml-2">Low {scans[`${image.name}:${image.tag}`].lowCount}</span>
                  </div>
                  <div className="space-y-1 max-h-48 overflow-auto border rounded">
                    {scans[`${image.name}:${image.tag}`].vulnerabilities.map(v => (
                      <div key={v.id+v.package} className="text-sm px-3 py-2 border-b last:border-b-0">
                        <div className="flex justify-between">
                          <span className="font-mono">{v.id}</span>
                          <span className="font-semibold">{v.severity}</span>
                        </div>
                        <div className="text-muted-foreground">{v.package} {v.installed} → {v.fixedVersion || 'n/a'}</div>
                        <div className="text-muted-foreground">{v.title}</div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}

