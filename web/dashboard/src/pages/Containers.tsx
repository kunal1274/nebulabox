import { useEffect, useState } from 'react'
import { Plus, Play, Square, FileText, RefreshCw } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { api, Container } from '@/lib/api'
import { Link } from 'react-router-dom'
import { TEST_IDS } from '@/lib/test-ids'

export function Containers() {
  const [containers, setContainers] = useState<Container[]>([])
  const [loading, setLoading] = useState(true)
  const [showAll, setShowAll] = useState(false)

  useEffect(() => {
    loadContainers()
  }, [showAll])

  const loadContainers = async () => {
    setLoading(true)
    try {
      const data = await api.listContainers(showAll)
      console.log('Loaded containers from API:', data)
      setContainers(data)
    } catch (error) {
      console.error('Failed to load containers:', error)
      // Only set empty array on error, don't show mock data
      // Mock data makes it confusing - we want to see the actual error
      setContainers([])
    } finally {
      setLoading(false)
    }
  }

  const handleStop = async (id: string) => {
    try {
      await api.stopContainer(id)
      loadContainers()
    } catch (error) {
      console.error('Failed to stop container:', error)
      alert('Failed to stop container')
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'running':
        return 'success'
      case 'stopped':
        return 'secondary'
      default:
        return 'default'
    }
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold mb-2">Containers</h1>
          <p className="text-muted-foreground">
            Manage your NebulaBox containers
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={() => setShowAll(!showAll)}
            data-test-id={TEST_IDS.containers.refresh}
          >
            <RefreshCw className="mr-2 h-4 w-4" />
            {showAll ? 'Show Running' : 'Show All'}
          </Button>
          <Link to="/containers/new">
            <Button data-test-id={TEST_IDS.containers.create}>
              <Plus className="mr-2 h-4 w-4" />
              New Container
            </Button>
          </Link>
        </div>
      </div>

      {loading ? (
        <div className="text-center py-12 text-muted-foreground">
          Loading containers...
        </div>
      ) : containers.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <p className="text-muted-foreground mb-4">No containers found</p>
            <Link to="/containers/new">
              <Button data-test-id={TEST_IDS.containers.create}>
                <Plus className="mr-2 h-4 w-4" />
                Create Your First Container
              </Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4" data-test-id={TEST_IDS.containers.list}>
          {containers.map((container) => (
            <Card key={container.id} data-test-id={TEST_IDS.containers.card} data-container-id={container.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="flex items-center gap-3">
                      {container.name || container.id}
                      <Badge variant={getStatusColor(container.status) as any}>
                        {container.status}
                      </Badge>
                    </CardTitle>
                    <div className="mt-2 text-sm text-muted-foreground">
                      <div>ID: {container.id}</div>
                      <div>Image: {container.image}</div>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    {container.status === 'running' ? (
                      <>
                        <Button
                          variant="outline"
                          size="icon"
                          onClick={() => handleStop(container.id)}
                          data-test-id={TEST_IDS.containers.stop}
                          data-container-id={container.id}
                        >
                          <Square className="h-4 w-4" />
                        </Button>
                        <Link to={`/containers/${container.id}/logs`}>
                          <Button variant="outline" size="icon" data-test-id={TEST_IDS.containers.logs} data-container-id={container.id}>
                            <FileText className="h-4 w-4" />
                          </Button>
                        </Link>
                      </>
                    ) : (
                      <Button variant="outline" size="icon" data-test-id={TEST_IDS.containers.start} data-container-id={container.id}>
                        <Play className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}

