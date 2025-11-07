import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { api } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Services() {
  const [services, setServices] = useState<Record<string, Array<{ id:string; address:string; port:number; version?:string; network?:string; createdAt:number }>>>({})
  const [name, setName] = useState('')
  const [resolveName, setResolveName] = useState('')
  const [resolved, setResolved] = useState<Array<{ id:string; address:string; port:number; version?:string; network?:string; createdAt:number }>>([])
  const [nextInstance, setNextInstance] = useState<{ id:string; address:string; port:number; version?:string; network?:string; createdAt:number } | null>(null)
  const [dnsName, setDnsName] = useState('')
  const [dnsResult, setDnsResult] = useState<string[]>([])
  const [dnsRecords, setDnsRecords] = useState<Record<string, string[]>>({})

  const load = async () => {
    try {
      const res = await api.listServices()
      setServices(res.services || {})
    } catch { setServices({}) }
  }
  useEffect(()=>{ load(); loadDns(); const t=setInterval(()=>{ load(); loadDns() }, 10000); return ()=>clearInterval(t) }, [])

  const register = async () => {
    if (!name.trim()) return
    await api.registerService({ name: name.trim(), address: '127.0.0.1', port: 0 })
    setName(''); load()
  }
  const dereg = async (svc: string, id: string) => { await api.deregisterService({ name: svc, id }); load() }
  const doResolve = async () => {
    if (!resolveName.trim()) { setResolved([]); return }
    const r = await api.resolveService(resolveName.trim())
    setResolved(r.instances || [])
  }
  const doNext = async () => {
    if (!resolveName.trim()) { setNextInstance(null); return }
    const r = await api.resolveServiceNext(resolveName.trim())
    setNextInstance((r as any).instance || null)
  }
  const loadDns = async () => {
    try {
      const r = await api.listDNSRecords()
      setDnsRecords(r.records || {})
    } catch { setDnsRecords({}) }
  }
  const doDnsResolve = async () => {
    if (!dnsName.trim()) { setDnsResult([]); return }
    const r = await api.dnsResolve(dnsName.trim())
    setDnsResult(r.a || [])
  }
  const addDns = async () => {
    if (!dnsName.trim()) return
    await api.addDNSRecord(dnsName.trim(), dnsResult.length ? dnsResult : ['127.0.0.1'])
    setDnsName(''); setDnsResult([]); loadDns()
  }
  const delDns = async (name: string) => { await api.deleteDNSRecord(name); loadDns() }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Service Discovery</h1>
        <p className="text-muted-foreground">Register and resolve service instances</p>
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Register Service (manual)</CardTitle>
          <CardDescription>Useful for testing discovery flow</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end">
            <div className="col-span-6">
              <div className="text-xs mb-1">Service Name</div>
              <Input value={name} onChange={e=>setName(e.target.value)} placeholder="api" data-test-id={TEST_IDS.services.serviceName} />
            </div>
            <div className="col-span-2">
              <Button onClick={register} data-test-id={TEST_IDS.services.registerService}>Register</Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Resolve Service</CardTitle>
          <CardDescription>Lookup instances by service name</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end">
            <div className="col-span-6">
              <Input value={resolveName} onChange={e=>setResolveName(e.target.value)} placeholder="api" />
            </div>
            <div className="col-span-2">
              <Button onClick={doResolve} data-test-id={TEST_IDS.services.createService}>Resolve</Button>
            </div>
            <div className="col-span-2">
              <Button variant="secondary" onClick={doNext}>Next</Button>
            </div>
          </div>
          {nextInstance && (
            <div className="mt-3 text-sm">Next instance: <span className="text-muted-foreground">{nextInstance.id} · {nextInstance.address}:{nextInstance.port}</span></div>
          )}
          <div className="mt-3 border rounded divide-y">
            {resolved.length===0 && (<div className="p-3 text-sm text-muted-foreground">No instances</div>)}
            {resolved.map(x=> (
              <div key={x.id} className="p-3 text-sm flex justify-between">
                <div>{x.id}</div>
                <div className="text-muted-foreground">{x.address}:{x.port} {x.network?`· ${x.network}`:''}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>DNS Resolve</CardTitle>
          <CardDescription>Resolve names via DNS API (.svc/.local use services)</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end">
            <div className="col-span-6">
              <Input value={dnsName} onChange={e=>setDnsName(e.target.value)} placeholder="api.svc" />
            </div>
            <div className="col-span-2">
              <Button onClick={doDnsResolve} data-test-id={TEST_IDS.services.createService}>Resolve</Button>
            </div>
            <div className="col-span-2">
              <Button variant="secondary" onClick={addDns} data-test-id={TEST_IDS.services.registerService}>Add A Record</Button>
            </div>
          </div>
          <div className="mt-3 text-sm">A records: {dnsResult.length? dnsResult.join(', ') : <span className="text-muted-foreground">none</span>}</div>
          <div className="mt-4">
            <div className="text-sm font-medium mb-2">Custom Records</div>
            <div className="border rounded divide-y">
              {Object.keys(dnsRecords).length===0 && (<div className="p-3 text-sm text-muted-foreground">No custom records</div>)}
              {Object.entries(dnsRecords).map(([n, ips])=> (
                <div key={n} className="p-3 text-sm flex justify-between">
                  <div>{n}</div>
                  <div className="text-muted-foreground">{(ips||[]).join(', ')}</div>
                  <Button variant="outline" onClick={()=>delDns(n)}>Delete</Button>
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Registered Services</CardTitle>
          <CardDescription>All known services and instances</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="border rounded divide-y">
            {Object.keys(services).length===0 && (<div className="p-3 text-sm text-muted-foreground">No services registered</div>)}
            {Object.entries(services).map(([svc, list]) => (
              <div key={svc} className="p-3">
                <div className="font-medium mb-2">{svc}</div>
                <div className="divide-y border rounded">
                  {list.map(x=> (
                    <div key={x.id} className="p-2 text-sm flex justify-between">
                      <div>{x.id}</div>
                      <div className="text-muted-foreground">{x.address}:{x.port} {x.network?`· ${x.network}`:''}</div>
                      <Button variant="outline" onClick={()=>dereg(svc, x.id)}>Remove</Button>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}


