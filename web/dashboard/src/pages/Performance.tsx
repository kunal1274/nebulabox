import { useEffect, useRef, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

type PerfMetrics = {
  reqPerSec1m: number
  p95Ms1m: number
  errorRate1m: number
  reqPerSec5m: number
  p95Ms5m: number
  errorRate5m: number
  timestamp: number
}

export default function Performance() {
  const [m, setM] = useState<PerfMetrics | null>(null)
  const esRef = useRef<EventSource | null>(null)

  useEffect(() => {
    const base = (import.meta as any).env.VITE_API_URL || 'http://localhost:8081/api'
    const es = new EventSource(`${base}/perf/stream`)
    es.addEventListener('perf', (e: MessageEvent) => {
      setM(JSON.parse(e.data))
    })
    es.onerror = () => es.close()
    esRef.current = es
    return () => { es.close() }
  }, [])

  const Item = ({label, value, suffix}:{label:string;value:string; suffix?:string}) => (
    <Card>
      <CardHeader><CardTitle className="text-base">{label}</CardTitle></CardHeader>
      <CardContent>
        <div className="text-3xl font-bold">{value}{suffix && <span className="text-lg text-muted-foreground ml-1">{suffix}</span>}</div>
      </CardContent>
    </Card>
  )

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Performance</h1>
        <p className="text-muted-foreground">Request rates, latency, error rate</p>
      </div>
      {m ? (
        <div className="grid md:grid-cols-3 gap-4">
          <Item label="Req/s (1m)" value={m.reqPerSec1m.toFixed(2)} />
          <Item label="p95 (1m)" value={m.p95Ms1m.toFixed(0)} suffix="ms" />
          <Item label="Error Rate (1m)" value={(m.errorRate1m*100).toFixed(2)} suffix="%" />
          <Item label="Req/s (5m)" value={m.reqPerSec5m.toFixed(2)} />
          <Item label="p95 (5m)" value={m.p95Ms5m.toFixed(0)} suffix="ms" />
          <Item label="Error Rate (5m)" value={(m.errorRate5m*100).toFixed(2)} suffix="%" />
        </div>
      ) : (
        <div className="text-muted-foreground">Waiting for metrics...</div>
      )}
    </div>
  )
}


