import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { api } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Networks() {
  const [userRole, setUserRole] = useState<string | null>(null)
  const [items, setItems] = useState<Array<{ id:string; name:string; driver:string; subnet:string; created:string }>>([])
  const [name, setName] = useState('')
  const [driver, setDriver] = useState('bridge')
  const [subnet, setSubnet] = useState('')

  const load = async () => {
    try { setItems(await api.listNetworks() as any) } catch { setItems([]) }
  }
  useEffect(()=>{
    load()
    api.me().then(r => setUserRole(r.user?.role || null)).catch(()=>{})
  }, [])

  const create = async () => {
    if (!name.trim()) return
    await api.createNetwork({ name: name.trim(), driver: driver || undefined, subnet: subnet || undefined })
    setName(''); setSubnet('')
    load()
  }
  const del = async (id: string) => { await api.deleteNetwork(id); load() }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Networks</h1>
        <p className="text-muted-foreground">Create and manage custom networks</p>
      </div>
      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Create Network</CardTitle>
          <CardDescription>Bridge networks for grouping containers</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end">
            <div className="col-span-4">
              <div className="text-xs mb-1">Name</div>
              <Input value={name} onChange={e=>setName(e.target.value)} placeholder="mynet" data-test-id={TEST_IDS.networks.networkName} />
            </div>
            <div className="col-span-3">
              <div className="text-xs mb-1">Driver</div>
              <select className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={driver} onChange={e=>setDriver(e.target.value)} data-test-id={TEST_IDS.networks.networkDriver}>
                <option value="bridge">bridge</option>
              </select>
            </div>
            <div className="col-span-3">
              <div className="text-xs mb-1">Subnet (optional)</div>
              <Input value={subnet} onChange={e=>setSubnet(e.target.value)} placeholder="10.10.0.0/24" />
            </div>
            <div className="col-span-2">
              <Button onClick={create} disabled={userRole ? userRole !== 'admin' && userRole !== 'editor' : false} data-test-id={TEST_IDS.networks.createNetwork}>Create</Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Existing Networks</CardTitle>
          <CardDescription>Click delete to remove</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="divide-y border rounded">
            {items.length === 0 && (
              <div className="p-3 text-sm text-muted-foreground">No networks yet</div>
            )}
            {items.map(n=> (
              <div key={n.id} className="p-3 flex items-center justify-between text-sm">
                <div>
                  <div className="font-medium">{n.name}</div>
                  <div className="text-muted-foreground">{n.driver}{n.subnet?` · ${n.subnet}`:''}</div>
                </div>
                <Button variant="outline" onClick={()=>del(n.id)} disabled={userRole ? userRole !== 'admin' && userRole !== 'editor' : false} data-test-id={TEST_IDS.networks.createNetwork}>Delete</Button>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}


