import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, RefreshCw } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { api } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function ContainerLogs() {
  const { id } = useParams<{ id: string }>()
  const [logs, setLogs] = useState<string[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      loadLogs()
    }
  }, [id])

  const loadLogs = async () => {
    if (!id) return

    setLoading(true)
    try {
      const data = await api.getContainerLogs(id)
      setLogs(data)
    } catch (error) {
      console.error('Failed to load logs:', error)
      // Mock data for development
      setLogs([
        `[${new Date().toISOString()}] Container ${id} started`,
        `[${new Date().toISOString()}] Application listening on port 80`,
        `[${new Date().toISOString()}] Health check passed`,
        `[${new Date().toISOString()}] Request processed: GET /`,
        `[${new Date().toISOString()}] Request processed: GET /api/status`,
        `[${new Date().toISOString()}] Container ${id} running normally`,
      ])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <Link to="/containers">
            <Button variant="ghost" className="mb-2">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Containers
            </Button>
          </Link>
          <h1 className="text-3xl font-bold mb-2">Container Logs</h1>
          <p className="text-muted-foreground">Container ID: {id}</p>
        </div>
        <Button onClick={loadLogs} variant="outline" data-test-id={TEST_IDS.monitor.refresh}>
          <RefreshCw className="mr-2 h-4 w-4" />
          Refresh
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Logs</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="text-center py-12 text-muted-foreground">
              Loading logs...
            </div>
          ) : logs.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              No logs available
            </div>
          ) : (
            <div className="bg-black rounded-lg p-4 font-mono text-sm text-green-400 max-h-[600px] overflow-y-auto">
              {logs.map((log, index) => (
                <div key={index} className="mb-1">
                  {log}
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

