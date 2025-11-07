import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Plus, Trash2, Edit, Save, X, Eye, EyeOff, Download, Upload } from 'lucide-react'
import { api } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

interface EnvVar {
  key: string
  value: string
  type: 'string' | 'number' | 'boolean' | 'secret'
}

interface EnvTemplate {
  name: string
  description: string
  variables: EnvVar[]
}

export default function ContainerEnv() {
  const { id } = useParams<{ id: string }>()
  const [envVars, setEnvVars] = useState<EnvVar[]>([])
  const [loading, setLoading] = useState(true)
  const [editing, setEditing] = useState(false)
  const [newVar, setNewVar] = useState<EnvVar>({ key: '', value: '', type: 'string' })
  const [envString, setEnvString] = useState('')
  const [templates, setTemplates] = useState<EnvTemplate[]>([])
  const [showSecrets, setShowSecrets] = useState(false)

  useEffect(() => {
    if (id) {
      loadEnvVars()
      loadTemplates()
    }
  }, [id])

  const loadEnvVars = async () => {
    try {
      setLoading(true)
      const response = await api.getContainerEnvVars(id!)
      if (response.success) {
        setEnvVars(response.variables)
      }
    } catch (error) {
      console.error('Failed to load environment variables:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadTemplates = async () => {
    try {
      const response = await api.getEnvTemplates()
      setTemplates(response.templates)
    } catch (error) {
      console.error('Failed to load templates:', error)
    }
  }

  const addEnvVar = () => {
    if (newVar.key.trim()) {
      setEnvVars([...envVars, { ...newVar }])
      setNewVar({ key: '', value: '', type: 'string' })
    }
  }

  const updateEnvVar = (index: number, field: keyof EnvVar, value: string) => {
    const updated = [...envVars]
    updated[index] = { ...updated[index], [field]: value }
    setEnvVars(updated)
  }

  const removeEnvVar = (index: number) => {
    setEnvVars(envVars.filter((_, i) => i !== index))
  }

  const saveEnvVars = async () => {
    try {
      const response = await api.setContainerEnvVars(id!, envVars)
      if (response.success) {
        setEditing(false)
        await loadEnvVars()
      }
    } catch (error) {
      console.error('Failed to save environment variables:', error)
    }
  }

  const clearEnvVars = async () => {
    try {
      const response = await api.clearContainerEnvVars(id!)
      if (response.success) {
        setEnvVars([])
        setEditing(false)
      }
    } catch (error) {
      console.error('Failed to clear environment variables:', error)
    }
  }

  const applyTemplate = (template: EnvTemplate) => {
    setEnvVars([...envVars, ...template.variables])
  }

  const exportAsString = () => {
    const envString = envVars
      .map(env => `${env.key}=${env.value}`)
      .join('\n')
    setEnvString(envString)
  }

  const importFromString = async () => {
    try {
      const response = await api.setContainerEnvFromString(id!, envString)
      if (response.success) {
        await loadEnvVars()
        setEnvString('')
      }
    } catch (error) {
      console.error('Failed to import environment variables:', error)
    }
  }

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'string': return 'bg-blue-100 text-blue-800'
      case 'number': return 'bg-green-100 text-green-800'
      case 'boolean': return 'bg-yellow-100 text-yellow-800'
      case 'secret': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const maskSecret = (value: string, type: string) => {
    if (type === 'secret' && !showSecrets) {
      return '•'.repeat(8)
    }
    return value
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-gray-600">Loading environment variables...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Environment Variables</h1>
          <p className="text-gray-600">Manage environment variables for container {id}</p>
        </div>
        <div className="flex space-x-2">
          <Button
            variant="outline"
            onClick={() => setShowSecrets(!showSecrets)}
            data-test-id={TEST_IDS.settings.saveSettings}
          >
            {showSecrets ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            {showSecrets ? 'Hide' : 'Show'} Secrets
          </Button>
          <Button
            variant="outline"
            onClick={exportAsString}
            data-test-id={TEST_IDS.settings.saveSettings}
          >
            <Download className="h-4 w-4" />
            Export
          </Button>
          {editing ? (
            <div className="flex space-x-2">
              <Button onClick={saveEnvVars} data-test-id={TEST_IDS.settings.saveSettings}>
                <Save className="h-4 w-4" />
                Save
              </Button>
              <Button variant="outline" onClick={() => setEditing(false)}>
                <X className="h-4 w-4" />
                Cancel
              </Button>
            </div>
          ) : (
            <Button onClick={() => setEditing(true)} data-test-id={TEST_IDS.settings.saveSettings}>
              <Edit className="h-4 w-4" />
              Edit
            </Button>
          )}
        </div>
      </div>

      <Tabs defaultValue="variables" className="space-y-4">
        <TabsList>
          <TabsTrigger value="variables">Variables</TabsTrigger>
          <TabsTrigger value="templates">Templates</TabsTrigger>
          <TabsTrigger value="import">Import/Export</TabsTrigger>
        </TabsList>

        <TabsContent value="variables" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Environment Variables ({envVars.length})</CardTitle>
              <CardDescription>
                Manage environment variables for this container
              </CardDescription>
            </CardHeader>
            <CardContent>
              {editing && (
                <div className="mb-4 p-4 border rounded-lg bg-gray-50">
                  <h3 className="font-medium mb-3">Add New Variable</h3>
                  <div className="grid grid-cols-12 gap-2">
                    <div className="col-span-3">
                      <Label htmlFor="key">Key</Label>
                      <Input
                        id="key"
                        value={newVar.key}
                        onChange={(e) => setNewVar({ ...newVar, key: e.target.value })}
                        placeholder="VARIABLE_NAME"
                      />
                    </div>
                    <div className="col-span-4">
                      <Label htmlFor="value">Value</Label>
                      <Input
                        id="value"
                        value={newVar.value}
                        onChange={(e) => setNewVar({ ...newVar, value: e.target.value })}
                        placeholder="variable_value"
                      />
                    </div>
                    <div className="col-span-3">
                      <Label htmlFor="type">Type</Label>
                      <Select
                        value={newVar.type}
                        onValueChange={(value: any) => setNewVar({ ...newVar, type: value })}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="string">String</SelectItem>
                          <SelectItem value="number">Number</SelectItem>
                          <SelectItem value="boolean">Boolean</SelectItem>
                          <SelectItem value="secret">Secret</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="col-span-2 flex items-end">
                      <Button onClick={addEnvVar} className="w-full">
                        <Plus className="h-4 w-4" />
                        Add
                      </Button>
                    </div>
                  </div>
                </div>
              )}

              <div className="space-y-2">
                {envVars.map((env, index) => (
                  <div key={index} className="flex items-center space-x-2 p-3 border rounded-lg">
                    <div className="flex-1 grid grid-cols-12 gap-2">
                      <div className="col-span-3">
                        <code className="text-sm font-mono">{env.key}</code>
                      </div>
                      <div className="col-span-6">
                        {editing ? (
                          <Input
                            value={env.value}
                            onChange={(e) => updateEnvVar(index, 'value', e.target.value)}
                            type={env.type === 'secret' && !showSecrets ? 'password' : 'text'}
                          />
                        ) : (
                          <code className="text-sm font-mono">
                            {maskSecret(env.value, env.type)}
                          </code>
                        )}
                      </div>
                      <div className="col-span-2">
                        <Badge className={getTypeColor(env.type)}>
                          {env.type}
                        </Badge>
                      </div>
                      <div className="col-span-1">
                        {editing && (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => removeEnvVar(index)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {envVars.length === 0 && (
                <div className="text-center py-8 text-gray-500">
                  No environment variables configured
                </div>
              )}

              {editing && envVars.length > 0 && (
                <div className="mt-4 flex justify-between">
                  <Button variant="destructive" onClick={clearEnvVars}>
                    Clear All
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="templates" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Environment Templates</CardTitle>
              <CardDescription>
                Apply pre-configured environment variable templates
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {templates.map((template, index) => (
                  <Card key={index} className="cursor-pointer hover:shadow-md transition-shadow">
                    <CardHeader>
                      <CardTitle className="text-lg">{template.name}</CardTitle>
                      <CardDescription>{template.description}</CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-2">
                        {template.variables.slice(0, 3).map((env, envIndex) => (
                          <div key={envIndex} className="flex items-center space-x-2">
                            <code className="text-xs font-mono">{env.key}</code>
                            <Badge className={getTypeColor(env.type)}>
                              {env.type}
                            </Badge>
                          </div>
                        ))}
                        {template.variables.length > 3 && (
                          <p className="text-xs text-gray-500">
                            +{template.variables.length - 3} more variables
                          </p>
                        )}
                      </div>
                      <Button
                        className="w-full mt-4"
                        onClick={() => applyTemplate(template)}
                      >
                        Apply Template
                      </Button>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="import" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Import/Export</CardTitle>
              <CardDescription>
                Import environment variables from text or export current variables
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <Label htmlFor="envString">Environment Variables (KEY=VALUE format)</Label>
                  <Textarea
                    id="envString"
                    value={envString}
                    onChange={(e) => setEnvString(e.target.value)}
                    placeholder="NODE_ENV=production&#10;PORT=3000&#10;DEBUG=false"
                    rows={8}
                    className="font-mono text-sm"
                  />
                </div>
                <div className="flex space-x-2">
                  <Button onClick={importFromString}>
                    <Upload className="h-4 w-4" />
                    Import
                  </Button>
                  <Button variant="outline" onClick={exportAsString}>
                    <Download className="h-4 w-4" />
                    Export
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
