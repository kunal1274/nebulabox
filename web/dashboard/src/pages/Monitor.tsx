import { useEffect, useState } from 'react'
import { Activity, Cpu, HardDrive, MemoryStick } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { api, SystemStats, SystemHistoryPoint } from '@/lib/api'
import type { AlertThresholds } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'
 

export function Monitor() {
  const [stats, setStats] = useState<SystemStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [alerts, setAlerts] = useState<AlertThresholds | null>(null)
  const [alertEvents, setAlertEvents] = useState<Array<{ type:string; value:number; threshold:number; timestamp:number }>>([])
  const [range, setRange] = useState<'1h'|'6h'|'24h'>('1h')
  const [history, setHistory] = useState<SystemHistoryPoint[]>([])

  useEffect(() => {
    // load alert thresholds
    api.getAlerts().then(setAlerts).catch(()=>{})
    // prefer SSE stream; fallback to polling
    let es: EventSource | null = null
    try {
      es = new EventSource(`${(import.meta as any).env.VITE_API_URL || 'http://localhost:8081/api'}/system/stream`)
      es.addEventListener('stats', (e: MessageEvent) => {
        const data = JSON.parse(e.data) as SystemStats & { timestamp?: number }
        setStats({
          cpuUsage: data.cpuUsage,
          memoryUsage: data.memoryUsage,
          diskUsage: data.diskUsage,
          containersRunning: data.containersRunning,
          containersTotal: data.containersTotal,
        })
        setLoading(false)
      })
      es.onerror = () => {
        es && es.close()
      }
    } catch (_) {}

    // alerts stream
    let esAlerts: EventSource | null = null
    try {
      esAlerts = new EventSource(`${(import.meta as any).env.VITE_API_URL || 'http://localhost:8081/api'}/alerts/stream`)
      esAlerts.addEventListener('alert', (e: MessageEvent) => {
        const a = JSON.parse(e.data) as { type:string; value:number; threshold:number; timestamp:number }
        setAlertEvents(prev => [a, ...prev].slice(0, 100))
      })
    } catch(_) {}
    const interval = setInterval(loadStats, 5000)
    return () => { clearInterval(interval); es && es.close(); esAlerts && esAlerts.close() }
  }, [])

  useEffect(() => {
    loadHistory()
    // refresh history periodically
    const t = setInterval(loadHistory, 15000)
    return () => clearInterval(t)
  }, [range])

  const loadStats = async () => {
    try {
      const data = await api.getSystemStats()
      setStats(data)
    } catch (error) {
      console.error('Failed to load stats:', error)
      // Mock data for development
      setStats({
        cpuUsage: 45.2,
        memoryUsage: 62.8,
        diskUsage: 38.5,
        containersRunning: 3,
        containersTotal: 5,
      })
    } finally {
      setLoading(false)
    }
  }

  const loadHistory = async () => {
    const r = range === '1h' ? 3600 : range === '6h' ? 21600 : 86400
    try {
      const res = await api.getSystemHistory({ range: r, step: Math.max(5, Math.floor(r/120)) })
      setHistory(res.points)
    } catch (_) {
      setHistory([])
    }
  }

  const Sparkline = ({ values, color }: { values: number[]; color: string }) => {
    const width = 300
    const height = 60
    if (values.length === 0) return <div className="text-muted-foreground">No history</div>
    const max = Math.max(100, ...values)
    const min = 0
    const pts = values.map((v,i)=>{
      const x = (i/(values.length-1)) * width
      const y = height - ((v - min)/(max - min)) * height
      return `${x},${y}`
    }).join(' ')
    return (
      <svg width={width} height={height} className="w-full h-16">
        <polyline fill="none" stroke={color} strokeWidth="2" points={pts} />
      </svg>
    )
  }

  const metrics = [
    {
      title: 'CPU Usage',
      value: stats?.cpuUsage.toFixed(1) || '0',
      unit: '%',
      icon: Cpu,
      color: 'bg-green-500',
      usage: stats?.cpuUsage || 0,
    },
    {
      title: 'Memory Usage',
      value: stats?.memoryUsage.toFixed(1) || '0',
      unit: '%',
      icon: MemoryStick,
      color: 'bg-orange-500',
      usage: stats?.memoryUsage || 0,
    },
    {
      title: 'Disk Usage',
      value: stats?.diskUsage.toFixed(1) || '0',
      unit: '%',
      icon: HardDrive,
      color: 'bg-purple-500',
      usage: stats?.diskUsage || 0,
    },
    {
      title: 'Active Containers',
      value: stats?.containersRunning.toString() || '0',
      unit: `/${stats?.containersTotal || 0} total`,
      icon: Activity,
      color: 'bg-blue-500',
      usage: stats
        ? (stats.containersRunning / stats.containersTotal) * 100
        : 0,
    },
  ]

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">System Monitor</h1>
        <p className="text-muted-foreground">
          Real-time system resource monitoring
        </p>
      </div>

      {loading && !stats ? (
        <div className="text-center py-12 text-muted-foreground">
          Loading metrics...
        </div>
      ) : (
        <div className="grid gap-6 md:grid-cols-2">
          {metrics.map((metric) => {
            const Icon = metric.icon
            return (
              <Card key={metric.title} data-test-id={TEST_IDS.monitor.metricCard}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle className="text-lg">{metric.title}</CardTitle>
                      <CardDescription>Current usage</CardDescription>
                    </div>
                    <Icon className="h-8 w-8 text-muted-foreground" />
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="text-3xl font-bold">
                      {metric.value}
                      <span className="text-lg text-muted-foreground ml-1">
                        {metric.unit}
                      </span>
                    </div>
                    <div>
                      <div className="w-full bg-secondary rounded-full h-3">
                        <div
                          className={`${metric.color} h-3 rounded-full transition-all duration-500`}
                          style={{ width: `${Math.min(metric.usage, 100)}%` }}
                        />
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}

      <Card className="mt-6">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Resource Usage History</CardTitle>
              <CardDescription>Historic CPU, Memory, and Disk usage</CardDescription>
            </div>
            <div className="flex gap-2">
              <Button variant={range==='1h'?'default':'secondary'} onClick={()=>setRange('1h')}>1h</Button>
              <Button variant={range==='6h'?'default':'secondary'} onClick={()=>setRange('6h')}>6h</Button>
              <Button variant={range==='24h'?'default':'secondary'} onClick={()=>setRange('24h')}>24h</Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <div>
              <div className="text-sm mb-2">CPU %</div>
              <Sparkline values={history.map(h=>h.cpu)} color="#22c55e" />
            </div>
            <div>
              <div className="text-sm mb-2">Memory %</div>
              <Sparkline values={history.map(h=>h.mem)} color="#f97316" />
            </div>
            <div>
              <div className="text-sm mb-2">Disk %</div>
              <Sparkline values={history.map(h=>h.disk)} color="#a855f7" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Container Status</CardTitle>
          <CardDescription>
            Overview of container activity
          </CardDescription>
        </CardHeader>
        <CardContent>
          {stats ? (
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <div className="text-sm text-muted-foreground mb-2">
                  Running Containers
                </div>
                <div className="text-2xl font-bold">
                  {stats.containersRunning}
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground mb-2">
                  Total Containers
                </div>
                <div className="text-2xl font-bold">
                  {stats.containersTotal}
                </div>
              </div>
            </div>
          ) : (
            <div className="text-sm text-muted-foreground">
              No data available
            </div>
          )}
        </CardContent>
      </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Alerts</CardTitle>
          <CardDescription>Configure thresholds and view recent alerts</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end">
            <div className="col-span-3">
              <div className="text-xs text-muted-foreground mb-1">CPU % High</div>
              <Input value={alerts?.cpuHigh ?? 85} onChange={(e)=>setAlerts({ ...(alerts||{cpuHigh:85,memoryHigh:80,diskHigh:90}), cpuHigh: Number(e.target.value) })} data-test-id={TEST_IDS.monitor.filterLogs} />
            </div>
            <div className="col-span-3">
              <div className="text-xs text-muted-foreground mb-1">Memory % High</div>
              <Input value={alerts?.memoryHigh ?? 80} onChange={(e)=>setAlerts({ ...(alerts||{cpuHigh:85,memoryHigh:80,diskHigh:90}), memoryHigh: Number(e.target.value) })} data-test-id={TEST_IDS.monitor.filterLogs} />
            </div>
            <div className="col-span-3">
              <div className="text-xs text-muted-foreground mb-1">Disk % High</div>
              <Input value={alerts?.diskHigh ?? 90} onChange={(e)=>setAlerts({ ...(alerts||{cpuHigh:85,memoryHigh:80,diskHigh:90}), diskHigh: Number(e.target.value) })} data-test-id={TEST_IDS.monitor.filterLogs} />
            </div>
            <div className="col-span-3">
              <Button onClick={() => alerts && api.setAlerts(alerts)} data-test-id={TEST_IDS.monitor.refresh}>Save</Button>
            </div>
          </div>
          <div className="mt-4 border rounded p-2 h-48 overflow-auto font-mono text-sm bg-muted">
            {alertEvents.map((a,i)=>(
              <div key={i}>[{new Date(a.timestamp*1000).toLocaleTimeString()}] {a.type.toUpperCase()} {a.value.toFixed(1)}% ≥ {a.threshold}%</div>
            ))}
            {alertEvents.length===0 && <div className="text-muted-foreground">No alerts</div>}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

