import { useEffect, useState } from 'react'
import { Package, Tag, Trash2, Calendar, User, Hash, FileText, RefreshCw, AlertCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { api, VersionSummary, ImageVersion } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { TEST_IDS } from '@/lib/test-ids'

export function Registry() {
  const [repositories, setRepositories] = useState<string[]>([])
  const [selectedRepo, setSelectedRepo] = useState<string | null>(null)
  const [versions, setVersions] = useState<ImageVersion[]>([])
  const [summary, setSummary] = useState<VersionSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [versionsLoading, setVersionsLoading] = useState(false)
  const [deletingTag, setDeletingTag] = useState<string | null>(null)

  useEffect(() => {
    loadRepositories()
  }, [])

  useEffect(() => {
    if (selectedRepo) {
      loadVersions(selectedRepo)
      loadSummary(selectedRepo)
    }
  }, [selectedRepo])

  const loadRepositories = async () => {
    setLoading(true)
    try {
      const data = await api.listRegistryRepositories()
      setRepositories(data.repositories)
      if (data.repositories.length > 0 && !selectedRepo) {
        setSelectedRepo(data.repositories[0])
      }
    } catch (error) {
      console.error('Failed to load repositories:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadVersions = async (repo: string) => {
    setVersionsLoading(true)
    try {
      const data = await api.listRepositoryVersions(repo)
      setVersions(data.versions || [])
    } catch (error) {
      console.error('Failed to load versions:', error)
      setVersions([])
    } finally {
      setVersionsLoading(false)
    }
  }

  const loadSummary = async (repo: string) => {
    try {
      const data = await api.getRepositorySummary(repo)
      setSummary(data)
    } catch (error) {
      console.error('Failed to load summary:', error)
    }
  }

  const handleDeleteVersion = async (repo: string, tag: string) => {
    if (!confirm(`Delete version ${tag} from ${repo}?`)) return
    
    setDeletingTag(tag)
    try {
      await api.deleteVersion(repo, tag)
      await loadVersions(repo)
      await loadSummary(repo)
    } catch (error) {
      console.error('Failed to delete version:', error)
      alert('Failed to delete version')
    } finally {
      setDeletingTag(null)
    }
  }

  const formatDate = (dateStr: string) => {
    if (!dateStr) return 'Unknown'
    try {
      return new Date(dateStr).toLocaleString()
    } catch {
      return dateStr
    }
  }

  const formatSize = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  const isLatest = (tag: string) => {
    return summary?.latest === tag || tag === 'latest'
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Registry</h1>
        <p className="text-muted-foreground">Manage container image repositories and versions</p>
      </div>

      <div className="grid lg:grid-cols-3 gap-6">
        {/* Repositories List */}
        <div className="lg:col-span-1">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <Package className="h-5 w-5" />
                  Repositories
                </CardTitle>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={loadRepositories}
                  disabled={loading}
                  data-test-id={TEST_IDS.registry.loadCatalog}
                >
                  <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-muted-foreground">Loading...</div>
              ) : repositories.length === 0 ? (
                <div className="text-muted-foreground text-sm">
                  No repositories found
                </div>
              ) : (
                <div className="space-y-2" data-test-id={TEST_IDS.registry.list}>
                  {repositories.map((repo) => (
                    <div
                      key={repo}
                      onClick={() => setSelectedRepo(repo)}
                      className={`p-3 rounded-lg border cursor-pointer transition-colors ${
                        selectedRepo === repo
                          ? 'bg-primary text-primary-foreground'
                          : 'hover:bg-accent'
                      }`}
                      data-test-id={TEST_IDS.registry.card}
                      data-repo-name={repo}
                    >
                      <div className="font-medium truncate">{repo}</div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Versions and Details */}
        <div className="lg:col-span-2 space-y-6">
          {selectedRepo ? (
            <>
              {/* Repository Summary */}
              {summary && (
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <FileText className="h-5 w-5" />
                      {selectedRepo}
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-3 gap-4">
                      <div>
                        <div className="text-sm text-muted-foreground">Total Versions</div>
                        <div className="text-2xl font-bold">{summary.totalVersions}</div>
                      </div>
                      <div>
                        <div className="text-sm text-muted-foreground">Latest</div>
                        <div className="text-lg font-semibold">
                          {summary.latest ? (
                            <Badge variant="default">{summary.latest}</Badge>
                          ) : (
                            <span className="text-muted-foreground">-</span>
                          )}
                        </div>
                      </div>
                      <div>
                        <div className="text-sm text-muted-foreground">Last Updated</div>
                        <div className="text-sm">{formatDate(summary.updatedAt)}</div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              )}

              {/* Versions List */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Tag className="h-5 w-5" />
                    Versions
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {versionsLoading ? (
                    <div className="text-muted-foreground">Loading versions...</div>
                  ) : versions.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground">
                      <AlertCircle className="h-12 w-12 mx-auto mb-2 opacity-50" />
                      <p>No versions found for this repository</p>
                    </div>
                  ) : (
                    <div className="space-y-4" data-test-id={TEST_IDS.registry.list}>
                      {versions.map((version) => (
                        <div
                          key={version.tag}
                          className="p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                          data-test-id={TEST_IDS.registry.card}
                          data-version-tag={version.tag}
                        >
                          <div className="flex items-start justify-between">
                            <div className="flex-1">
                              <div className="flex items-center gap-2 mb-2">
                                <span className="font-semibold text-lg">{version.tag}</span>
                                {isLatest(version.tag) && (
                                  <Badge variant="default">Latest</Badge>
                                )}
                              </div>
                              
                              <div className="grid grid-cols-2 gap-4 text-sm text-muted-foreground">
                                <div className="flex items-center gap-2">
                                  <Hash className="h-4 w-4" />
                                  <span className="font-mono text-xs truncate">
                                    {version.digest.substring(0, 20)}...
                                  </span>
                                </div>
                                
                                <div className="flex items-center gap-2">
                                  <Calendar className="h-4 w-4" />
                                  <span>{formatDate(version.createdAt)}</span>
                                </div>
                                
                                <div className="flex items-center gap-2">
                                  <User className="h-4 w-4" />
                                  <span>{version.createdBy || 'Unknown'}</span>
                                </div>
                                
                                <div className="flex items-center gap-2">
                                  <Package className="h-4 w-4" />
                                  <span>{formatSize(version.size)}</span>
                                </div>
                              </div>
                              
                              {version.description && (
                                <div className="mt-2 text-sm text-muted-foreground">
                                  {version.description}
                                </div>
                              )}
                              
                              {version.metadata && Object.keys(version.metadata).length > 0 && (
                                <div className="mt-2 flex flex-wrap gap-1">
                                  {Object.entries(version.metadata).map(([key, value]) => (
                                    <Badge key={key} variant="outline" className="text-xs">
                                      {key}: {value}
                                    </Badge>
                                  ))}
                                </div>
                              )}
                            </div>
                            
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleDeleteVersion(selectedRepo, version.tag)}
                              disabled={deletingTag === version.tag}
                              className="text-destructive hover:text-destructive"
                              data-test-id={TEST_IDS.registry.deleteTag}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </>
          ) : (
            <Card>
              <CardContent className="py-12 text-center text-muted-foreground">
                <Package className="h-12 w-12 mx-auto mb-4 opacity-50" />
                <p>Select a repository to view versions</p>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}

