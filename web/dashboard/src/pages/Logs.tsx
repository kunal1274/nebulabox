import { useEffect, useRef, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { api, LogEntry } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export default function Logs() {
  const [query, setQuery] = useState('')
  const [containerId, setContainerId] = useState('')
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [streaming, setStreaming] = useState(false)
  const esRef = useRef<EventSource | null>(null)

  const search = async () => {
    const data = await api.searchLogs({ query, containerId, limit: 200 })
    setLogs(data)
  }

  const startStream = () => {
    stopStream()
    const base = (import.meta as any).env.VITE_API_URL || 'http://localhost:8081/api'
    const url = `${base}/logs/stream${containerId ? `?containerId=${encodeURIComponent(containerId)}` : ''}`
    const es = new EventSource(url)
    es.addEventListener('log', (e: MessageEvent) => {
      const entry = JSON.parse(e.data) as LogEntry
      setLogs(prev => [entry, ...prev].slice(0, 500))
    })
    es.onerror = () => { es.close(); setStreaming(false) }
    esRef.current = es
    setStreaming(true)
  }

  const stopStream = () => {
    esRef.current?.close(); esRef.current = null; setStreaming(false)
  }

  useEffect(() => { return () => stopStream() }, [])

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Logs</h1>
        <p className="text-muted-foreground">Aggregate search and live streaming</p>
      </div>

      <Card className="mb-4">
        <CardHeader>
          <CardTitle>Search</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end">
            <div className="col-span-6">
              <Input placeholder="query text" value={query} onChange={(e)=>setQuery(e.target.value)} />
            </div>
            <div className="col-span-4">
              <Input placeholder="container id (optional)" value={containerId} onChange={(e)=>setContainerId(e.target.value)} />
            </div>
            <div className="col-span-2">
              <Button onClick={search} data-test-id={TEST_IDS.monitor.filterLogs}>Search</Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Live Stream</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2 mb-3">
            <Button onClick={startStream} disabled={streaming} data-test-id={TEST_IDS.monitor.refresh}>Start</Button>
            <Button variant="outline" onClick={stopStream} disabled={!streaming}>Stop</Button>
          </div>
          <div className="border rounded p-2 h-72 overflow-auto font-mono text-sm bg-muted">
            {logs.map((l, i) => (
              <div key={i}>
                [{new Date(l.timestamp*1000).toLocaleTimeString()}] {l.container} {l.level} - {l.message}
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}


