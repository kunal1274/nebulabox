import { useState } from 'react'
import { Code, Play, CheckCircle, XCircle, FileCode, RefreshCw } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { api } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function BuildSpec() {
  const [spec, setSpec] = useState<any>({
    version: '1.0',
    name: 'my-app',
    tag: 'my-app:latest',
    base: {
      image: 'alpine',
      tag: '3.19',
    },
    steps: [
      { type: 'run', command: 'apk add --no-cache nodejs npm' },
    ],
    env: {},
    workdir: '/app',
  })
  const [specJson, setSpecJson] = useState(JSON.stringify(spec, null, 2))
  const [tag, setTag] = useState('my-app:latest')
  const [validation, setValidation] = useState<{ valid: boolean; dockerfile?: string; errors?: string[]; message?: string } | null>(null)
  const [dockerfile, setDockerfile] = useState('')
  const [buildLogs, setBuildLogs] = useState<string[]>([])
  const [loading, setLoading] = useState(false)

  const updateSpec = () => {
    try {
      const parsed = JSON.parse(specJson)
      setSpec(parsed)
      if (parsed.tag) setTag(parsed.tag)
    } catch (e) {
      console.error('Invalid JSON:', e)
    }
  }

  const handleValidate = async () => {
    setLoading(true)
    try {
      const parsed = JSON.parse(specJson)
      const result = await api.validateBuildSpec(parsed)
      setValidation(result)
      if (result.valid && result.dockerfile) {
        setDockerfile(result.dockerfile)
      }
    } catch (error: any) {
      setValidation({
        valid: false,
        errors: [error.message || 'Validation failed'],
        message: 'Failed to validate specification',
      })
    } finally {
      setLoading(false)
    }
  }

  const handleConvert = async () => {
    setLoading(true)
    try {
      const parsed = JSON.parse(specJson)
      const result = await api.convertBuildSpec(parsed)
      if (result.dockerfile) {
        setDockerfile(result.dockerfile)
        setValidation({ valid: true, message: result.message })
      }
    } catch (error: any) {
      setValidation({
        valid: false,
        errors: [error.message || 'Conversion failed'],
        message: 'Failed to convert specification',
      })
    } finally {
      setLoading(false)
    }
  }

  const handleBuild = async () => {
    if (!tag.trim()) {
      alert('Please enter a tag')
      return
    }
    
    setLoading(true)
    setBuildLogs([])
    try {
      const parsed = JSON.parse(specJson)
      const result = await api.buildFromSpec(parsed, tag)
      if (result.logs) {
        setBuildLogs(result.logs)
      }
      if (result.dockerfile) {
        setDockerfile(result.dockerfile)
      }
      setValidation({ valid: true, message: result.message || 'Build completed' })
    } catch (error: any) {
      setBuildLogs([`Build failed: ${error.message || 'Unknown error'}`])
      setValidation({
        valid: false,
        errors: [error.message || 'Build failed'],
        message: 'Build failed',
      })
    } finally {
      setLoading(false)
    }
  }

  const exampleSpec = {
    version: '1.0',
    name: 'node-app',
    tag: 'node-app:latest',
    base: {
      image: 'node',
      tag: '18-alpine',
    },
    steps: [
      { type: 'run', command: 'npm install -g npm@latest', comment: 'Update npm' },
      { type: 'copy', source: 'package.json', dest: '/app/package.json', comment: 'Copy package file' },
      { type: 'run', command: 'npm install', comment: 'Install dependencies', workdir: '/app' },
      { type: 'copy', source: '.', dest: '/app', comment: 'Copy application code' },
      { type: 'cmd', command: '["node", "index.js"]', comment: 'Start command' },
    ],
    env: {
      NODE_ENV: 'production',
    },
    workdir: '/app',
    expose: [3000],
    labels: {
      'maintainer': 'NebulaBox',
      'version': '1.0',
    },
    health: {
      type: 'http',
      path: '/health',
      port: 3000,
      interval: 30,
      timeout: 10,
      retries: 3,
    },
  }

  const loadExample = () => {
    setSpec(exampleSpec)
    setSpecJson(JSON.stringify(exampleSpec, null, 2))
    setTag(exampleSpec.tag)
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Build Specification</h1>
        <p className="text-muted-foreground">
          Create and build container images using NebulaBox build specification format
        </p>
      </div>

      <div className="grid lg:grid-cols-2 gap-6 mb-6">
        {/* Build Specification Editor */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <span className="flex items-center gap-2">
                <FileCode className="h-5 w-5" />
                Build Specification
              </span>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={loadExample} data-test-id={TEST_IDS.buildspec.loadExample}>
                  Load Example
                </Button>
                <Button variant="outline" size="sm" onClick={updateSpec}>
                  <RefreshCw className="h-4 w-4" />
                </Button>
              </div>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <Label>Tag</Label>
                <Input
                  value={tag}
                  onChange={(e) => setTag(e.target.value)}
                  placeholder="my-app:latest"
                  data-test-id={TEST_IDS.buildspec.tagInput}
                />
              </div>
              <div>
                <Label>Specification JSON</Label>
                <Textarea
                  rows={20}
                  value={specJson}
                  onChange={(e) => setSpecJson(e.target.value)}
                  className="font-mono text-sm"
                  placeholder='{"version": "1.0", "name": "...", ...}'
                  data-test-id={TEST_IDS.buildspec.specEditor}
                />
              </div>
              <div className="flex gap-2">
                <Button onClick={handleValidate} disabled={loading} variant="outline" data-test-id={TEST_IDS.buildspec.validate}>
                  <CheckCircle className="mr-2 h-4 w-4" />
                  Validate
                </Button>
                <Button onClick={handleConvert} disabled={loading} variant="outline" data-test-id={TEST_IDS.buildspec.convert}>
                  <Code className="mr-2 h-4 w-4" />
                  Convert to Dockerfile
                </Button>
                <Button onClick={handleBuild} disabled={loading} data-test-id={TEST_IDS.buildspec.build}>
                  <Play className="mr-2 h-4 w-4" />
                  Build
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Validation Result & Dockerfile */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Code className="h-5 w-5" />
              Output
            </CardTitle>
          </CardHeader>
          <CardContent>
            <Tabs defaultValue="dockerfile" className="w-full">
              <TabsList>
                <TabsTrigger value="dockerfile" data-test-id={TEST_IDS.buildspec.dockerfileTab}>Dockerfile</TabsTrigger>
                <TabsTrigger value="validation" data-test-id={TEST_IDS.buildspec.validationTab}>Validation</TabsTrigger>
                <TabsTrigger value="logs" data-test-id={TEST_IDS.buildspec.logsTab}>Build Logs</TabsTrigger>
              </TabsList>
              
              <TabsContent value="dockerfile" className="space-y-2">
                <Textarea
                  rows={20}
                  value={dockerfile}
                  readOnly
                  className="font-mono text-xs"
                  placeholder="Generated Dockerfile will appear here..."
                />
                {dockerfile && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => navigator.clipboard.writeText(dockerfile)}
                  >
                    Copy Dockerfile
                  </Button>
                )}
              </TabsContent>
              
              <TabsContent value="validation" className="space-y-2">
                {validation && (
                  <div className={`p-4 rounded-lg border ${
                    validation.valid
                      ? 'bg-green-50 border-green-200 dark:bg-green-950 dark:border-green-800'
                      : 'bg-red-50 border-red-200 dark:bg-red-950 dark:border-red-800'
                  }`}>
                    <div className="flex items-center gap-2 mb-2">
                      {validation.valid ? (
                        <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400" />
                      ) : (
                        <XCircle className="h-5 w-5 text-red-600 dark:text-red-400" />
                      )}
                      <span className="font-semibold">
                        {validation.valid ? 'Valid' : 'Invalid'}
                      </span>
                    </div>
                    {validation.message && (
                      <p className="text-sm mb-2">{validation.message}</p>
                    )}
                    {validation.errors && validation.errors.length > 0 && (
                      <div className="space-y-1">
                        <p className="text-sm font-medium">Errors:</p>
                        <ul className="list-disc list-inside text-sm">
                          {validation.errors.map((err, i) => (
                            <li key={i}>{err}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                )}
                {!validation && (
                  <p className="text-sm text-muted-foreground">
                    Click "Validate" to check your specification
                  </p>
                )}
              </TabsContent>
              
              <TabsContent value="logs" className="space-y-2">
                {buildLogs.length > 0 ? (
                  <div className="border rounded p-3 bg-muted max-h-96 overflow-auto">
                    <pre className="text-xs font-mono whitespace-pre-wrap">
                      {buildLogs.join('\n')}
                    </pre>
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">
                    Build logs will appear here after building
                  </p>
                )}
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>
      </div>

      {/* Documentation */}
      <Card>
        <CardHeader>
          <CardTitle>Build Specification Format</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4 text-sm">
            <div>
              <h3 className="font-semibold mb-2">Required Fields:</h3>
              <ul className="list-disc list-inside space-y-1 text-muted-foreground">
                <li><code className="bg-muted px-1 rounded">name</code> - Application name</li>
                <li><code className="bg-muted px-1 rounded">base</code> - Base image (image, tag)</li>
                <li><code className="bg-muted px-1 rounded">steps</code> - Build steps array</li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold mb-2">Step Types:</h3>
              <ul className="list-disc list-inside space-y-1 text-muted-foreground">
                <li><code className="bg-muted px-1 rounded">run</code> - Execute a command</li>
                <li><code className="bg-muted px-1 rounded">copy</code> - Copy files (source, dest)</li>
                <li><code className="bg-muted px-1 rounded">add</code> - Add files with URL support</li>
                <li><code className="bg-muted px-1 rounded">cmd</code> - Set default command</li>
                <li><code className="bg-muted px-1 rounded">arg</code> - Define build arguments</li>
                <li><code className="bg-muted px-1 rounded">volume</code> - Create volume mount</li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold mb-2">Optional Fields:</h3>
              <ul className="list-disc list-inside space-y-1 text-muted-foreground">
                <li><code className="bg-muted px-1 rounded">env</code> - Environment variables</li>
                <li><code className="bg-muted px-1 rounded">workdir</code> - Working directory</li>
                <li><code className="bg-muted px-1 rounded">expose</code> - Exposed ports array</li>
                <li><code className="bg-muted px-1 rounded">labels</code> - Image labels</li>
                <li><code className="bg-muted px-1 rounded">health</code> - Health check configuration</li>
                <li><code className="bg-muted px-1 rounded">user</code> - Run as user</li>
              </ul>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

