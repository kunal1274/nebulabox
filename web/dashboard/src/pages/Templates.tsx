import { useEffect, useState } from 'react'
import { Layers, Play, Trash2, Eye, Filter, Search } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { api, StackTemplate } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Templates() {
  const [templates, setTemplates] = useState<StackTemplate[]>([])
  const [filteredTemplates, setFilteredTemplates] = useState<StackTemplate[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedTemplate, setSelectedTemplate] = useState<StackTemplate | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [tagFilter, setTagFilter] = useState('')

  useEffect(() => {
    loadTemplates()
  }, [categoryFilter, tagFilter])

  useEffect(() => {
    filterTemplates()
  }, [templates, searchQuery, categoryFilter, tagFilter])

  const loadTemplates = async () => {
    setLoading(true)
    try {
      const data = await api.listTemplates(categoryFilter || undefined, tagFilter || undefined)
      setTemplates(data.templates)
    } catch (error) {
      console.error('Failed to load templates:', error)
    } finally {
      setLoading(false)
    }
  }

  const filterTemplates = () => {
    let filtered = [...templates]

    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      filtered = filtered.filter(
        (t) =>
          t.name.toLowerCase().includes(query) ||
          t.description.toLowerCase().includes(query) ||
          t.category.toLowerCase().includes(query) ||
          (t.tags && t.tags.some((tag) => tag.toLowerCase().includes(query)))
      )
    }

    setFilteredTemplates(filtered)
  }

  const handleDeploy = async (templateId: string) => {
    const prefix = prompt('Enter prefix for container names (optional):')
    if (prefix === null) return

    try {
      const result = await api.deployTemplate(templateId, prefix || undefined, undefined, true)
      alert(`Template deployment initiated! ${result.message || 'Check containers page for status.'}`)
    } catch (error: any) {
      console.error('Failed to deploy:', error)
      alert(`Failed to deploy template: ${error.message}`)
    }
  }

  const handleDelete = async (templateId: string) => {
    if (!confirm(`Delete template "${templateId}"?`)) return

    try {
      await api.deleteTemplate(templateId)
      loadTemplates()
      if (selectedTemplate?.id === templateId) {
        setSelectedTemplate(null)
      }
    } catch (error: any) {
      console.error('Failed to delete:', error)
      alert(`Failed to delete template: ${error.message}`)
    }
  }

  const categories = Array.from(new Set(templates.map((t) => t.category)))
  const allTags = Array.from(new Set(templates.flatMap((t) => t.tags || [])))

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Stack Templates</h1>
        <p className="text-muted-foreground">
          Deploy pre-configured multi-container stacks with one click
        </p>
      </div>

      {/* Filters */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Filters
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="text-sm font-medium mb-2 block">Search</label>
              <div className="relative">
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search templates..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-8"
                  data-test-id={TEST_IDS.templates.deployTemplate}
                />
              </div>
            </div>
            <div>
              <label className="text-sm font-medium mb-2 block">Category</label>
              <select
                value={categoryFilter}
                onChange={(e) => setCategoryFilter(e.target.value)}
                className="w-full px-3 py-2 border rounded-md"
              >
                <option value="">All Categories</option>
                {categories.map((cat) => (
                  <option key={cat} value={cat}>
                    {cat}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-sm font-medium mb-2 block">Tag</label>
              <select
                value={tagFilter}
                onChange={(e) => setTagFilter(e.target.value)}
                className="w-full px-3 py-2 border rounded-md"
              >
                <option value="">All Tags</option>
                {allTags.map((tag) => (
                  <option key={tag} value={tag}>
                    {tag}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Templates Grid */}
      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">Loading templates...</p>
        </div>
      ) : filteredTemplates.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Layers className="h-12 w-12 mb-4 opacity-50" />
            <p className="text-muted-foreground">No templates found</p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredTemplates.map((template) => (
            <Card
              key={template.id}
              className={`cursor-pointer transition-all hover:shadow-lg ${
                selectedTemplate?.id === template.id ? 'ring-2 ring-primary' : ''
              }`}
              onClick={() => setSelectedTemplate(template)}
              data-test-id={TEST_IDS.templates.templateCard}
            >
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <CardTitle className="flex items-center gap-2 mb-2">
                      {template.name}
                      <Badge variant="outline">{template.category}</Badge>
                    </CardTitle>
                    <p className="text-sm text-muted-foreground">{template.description}</p>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div>
                    <div className="text-xs font-semibold text-muted-foreground mb-1">
                      Containers ({template.containers.length})
                    </div>
                    <div className="flex flex-wrap gap-1">
                      {template.containers.slice(0, 3).map((c) => (
                        <Badge key={c.name} variant="secondary" className="text-xs">
                          {c.name}
                        </Badge>
                      ))}
                      {template.containers.length > 3 && (
                        <Badge variant="secondary" className="text-xs">
                          +{template.containers.length - 3} more
                        </Badge>
                      )}
                    </div>
                  </div>

                  {template.tags && template.tags.length > 0 && (
                    <div>
                      <div className="text-xs font-semibold text-muted-foreground mb-1">Tags</div>
                      <div className="flex flex-wrap gap-1">
                        {template.tags.map((tag) => (
                          <Badge key={tag} variant="outline" className="text-xs">
                            {tag}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  )}

                  <div className="flex gap-2 pt-2">
                    <Button
                      size="sm"
                      className="flex-1"
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeploy(template.id)
                      }}
                      data-test-id={TEST_IDS.templates.deployTemplate}
                    >
                      <Play className="h-3 w-3 mr-1" />
                      Deploy
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={(e) => {
                        e.stopPropagation()
                        setSelectedTemplate(template)
                      }}
                    >
                      <Eye className="h-3 w-3" />
                    </Button>
                    {template.id !== 'lamp' &&
                      template.id !== 'mean' &&
                      template.id !== 'lemp' &&
                      template.id !== 'wordpress' &&
                      template.id !== 'microservices' &&
                      template.id !== 'django-postgres' && (
                        <Button
                          size="sm"
                          variant="destructive"
                          onClick={(e) => {
                            e.stopPropagation()
                            handleDelete(template.id)
                          }}
                          data-test-id={TEST_IDS.templates.deleteTemplate}
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      )}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Template Details Sidebar */}
      {selectedTemplate && (
        <Card className="mt-6">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>{selectedTemplate.name}</CardTitle>
              <Button variant="ghost" size="sm" onClick={() => setSelectedTemplate(null)}>
                ×
              </Button>
            </div>
          </CardHeader>
          <CardContent className="space-y-4 max-h-96 overflow-y-auto">
            <div>
              <div className="font-semibold mb-2">Description</div>
              <p className="text-sm text-muted-foreground">{selectedTemplate.description}</p>
            </div>

            <div>
              <div className="font-semibold mb-2">Category</div>
              <Badge>{selectedTemplate.category}</Badge>
            </div>

            <div>
              <div className="font-semibold mb-2">Containers</div>
              <div className="space-y-2">
                {selectedTemplate.containers.map((container) => (
                  <div key={container.name} className="p-3 border rounded-lg">
                    <div className="font-medium mb-1">{container.name}</div>
                    <div className="text-sm text-muted-foreground">
                      Image: {container.image}
                    </div>
                    {container.ports && Object.keys(container.ports).length > 0 && (
                      <div className="text-xs text-muted-foreground mt-1">
                        Ports: {Object.entries(container.ports).map(([c, h]) => `${h}:${c}`).join(', ')}
                      </div>
                    )}
                    {container.dependsOn && container.dependsOn.length > 0 && (
                      <div className="text-xs text-muted-foreground mt-1">
                        Depends on: {container.dependsOn.join(', ')}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {selectedTemplate.networks && selectedTemplate.networks.length > 0 && (
              <div>
                <div className="font-semibold mb-2">Networks</div>
                <div className="flex flex-wrap gap-1">
                  {selectedTemplate.networks.map((net) => (
                    <Badge key={net.name} variant="outline">
                      {net.name}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            {selectedTemplate.tags && selectedTemplate.tags.length > 0 && (
              <div>
                <div className="font-semibold mb-2">Tags</div>
                <div className="flex flex-wrap gap-1">
                  {selectedTemplate.tags.map((tag) => (
                    <Badge key={tag} variant="secondary">
                      {tag}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            <Button
              className="w-full"
              onClick={() => handleDeploy(selectedTemplate.id)}
              data-test-id={TEST_IDS.templates.deployTemplate}
            >
              <Play className="h-4 w-4 mr-2" />
              Deploy Stack
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

