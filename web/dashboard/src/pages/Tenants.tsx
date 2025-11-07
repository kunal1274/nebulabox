import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { api } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function Tenants() {
  const [tenants, setTenants] = useState<Array<{
    id: string
    name: string
    domain: string
    created: string
    createdBy: string
    quota: { maxContainers: number; maxNetworks: number; maxTeams: number; maxStorageGB: number }
  }>>([])
  const [selectedTenant, setSelectedTenant] = useState<{
    id: string
    name: string
    domain: string
    quota: any
  } | null>(null)
  const [usage, setUsage] = useState<{ containers: number; networks: number; teams: number } | null>(null)
  const [newTenantName, setNewTenantName] = useState('')
  const [newTenantDomain, setNewTenantDomain] = useState('')
  const [newQuota, setNewQuota] = useState({ maxContainers: 100, maxNetworks: 50, maxTeams: 20, maxStorageGB: 1000 })
  const [assignUsername, setAssignUsername] = useState('')
  const [me, setMe] = useState<{ username: string; role: string } | null>(null)

  useEffect(() => {
    loadTenants()
    api.me().then(r => setMe(r.user || null)).catch(() => {})
  }, [])

  const loadTenants = async () => {
    try {
      const data = await api.listTenants()
      setTenants(data)
    } catch {
      setTenants([])
    }
  }

  const loadTenantDetails = async (id: string) => {
    try {
      const tenant = await api.getTenant(id)
      setSelectedTenant(tenant)
      try {
        const usageData = await api.getTenantUsage(id)
        setUsage(usageData.usage)
        setSelectedTenant({ ...tenant, quota: usageData.quota })
      } catch {}
    } catch {
      setSelectedTenant(null)
      setUsage(null)
    }
  }

  const createTenant = async () => {
    if (!newTenantName.trim()) return
    try {
      await api.createTenant(newTenantName.trim(), newTenantDomain.trim() || undefined, newQuota)
      setNewTenantName('')
      setNewTenantDomain('')
      setNewQuota({ maxContainers: 100, maxNetworks: 50, maxTeams: 20, maxStorageGB: 1000 })
      loadTenants()
    } catch (err: any) {
      alert(err.message || 'Failed to create tenant')
    }
  }

  const deleteTenant = async (id: string) => {
    if (!confirm('Delete this tenant? All associated resources will be orphaned.')) return
    try {
      await api.deleteTenant(id)
      if (selectedTenant?.id === id) { setSelectedTenant(null); setUsage(null) }
      loadTenants()
    } catch (err: any) {
      alert(err.message || 'Failed to delete tenant')
    }
  }

  const assignUser = async () => {
    if (!selectedTenant || !assignUsername.trim()) return
    try {
      await api.assignUserToTenant(selectedTenant.id, assignUsername.trim())
      setAssignUsername('')
      alert('User assigned to tenant')
    } catch (err: any) {
      alert(err.message || 'Failed to assign user')
    }
  }

  const isSystemAdmin = me?.role === 'admin'

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Tenants</h1>
        <p className="text-muted-foreground">Manage multi-tenant organizations</p>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div>
          {isSystemAdmin && (
            <Card className="mb-6">
              <CardHeader>
                <CardTitle>Create Tenant</CardTitle>
                <CardDescription>Create a new organization tenant</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div>
                    <div className="text-xs mb-1">Name</div>
                    <Input value={newTenantName} onChange={e => setNewTenantName(e.target.value)} placeholder="acme-corp" data-test-id={TEST_IDS.tenants.tenantName} />
                  </div>
                  <div>
                    <div className="text-xs mb-1">Domain (optional)</div>
                    <Input value={newTenantDomain} onChange={e => setNewTenantDomain(e.target.value)} placeholder="acme.com" />
                  </div>
                  <div className="grid grid-cols-2 gap-2 text-xs">
                    <div>
                      <div className="mb-1">Max Containers</div>
                      <Input type="number" value={newQuota.maxContainers} onChange={e => setNewQuota({ ...newQuota, maxContainers: parseInt(e.target.value) || 0 })} />
                    </div>
                    <div>
                      <div className="mb-1">Max Networks</div>
                      <Input type="number" value={newQuota.maxNetworks} onChange={e => setNewQuota({ ...newQuota, maxNetworks: parseInt(e.target.value) || 0 })} />
                    </div>
                    <div>
                      <div className="mb-1">Max Teams</div>
                      <Input type="number" value={newQuota.maxTeams} onChange={e => setNewQuota({ ...newQuota, maxTeams: parseInt(e.target.value) || 0 })} />
                    </div>
                    <div>
                      <div className="mb-1">Max Storage (GB)</div>
                      <Input type="number" value={newQuota.maxStorageGB} onChange={e => setNewQuota({ ...newQuota, maxStorageGB: parseInt(e.target.value) || 0 })} />
                    </div>
                  </div>
                  <Button onClick={createTenant} data-test-id={TEST_IDS.tenants.createTenant}>Create Tenant</Button>
                </div>
              </CardContent>
            </Card>
          )}

          <Card>
            <CardHeader>
              <CardTitle>Tenants</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {tenants.length === 0 && <div className="text-sm text-muted-foreground">No tenants found</div>}
                {tenants.map(t => (
                  <div
                    key={t.id}
                    className={`p-3 border rounded cursor-pointer hover:bg-accent ${selectedTenant?.id === t.id ? 'bg-accent' : ''}`}
                    onClick={() => loadTenantDetails(t.id)}
                  >
                    <div className="font-medium">{t.name}</div>
                    {t.domain && <div className="text-xs text-muted-foreground">{t.domain}</div>}
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {selectedTenant && (
          <Card>
            <CardHeader>
              <CardTitle>{selectedTenant.name}</CardTitle>
              <CardDescription>
                {selectedTenant.domain || 'No domain'}
                {isSystemAdmin && (
                  <Button variant="destructive" size="sm" className="ml-2" onClick={() => deleteTenant(selectedTenant.id)}>
                    Delete Tenant
                  </Button>
                )}
              </CardDescription>
            </CardHeader>
            <CardContent>
              {usage && (
                <div className="mb-4">
                  <div className="text-sm font-medium mb-2">Resource Usage</div>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span>Containers:</span>
                      <span>{usage.containers} / {selectedTenant.quota.maxContainers}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Networks:</span>
                      <span>{usage.networks} / {selectedTenant.quota.maxNetworks}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Teams:</span>
                      <span>{usage.teams} / {selectedTenant.quota.maxTeams}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Storage:</span>
                      <span>{selectedTenant.quota.maxStorageGB} GB limit</span>
                    </div>
                  </div>
                </div>
              )}

              {isSystemAdmin && (
                <div className="border-t pt-4">
                  <div className="text-sm font-medium mb-2">Assign User</div>
                  <div className="grid grid-cols-12 gap-2">
                    <div className="col-span-9">
                      <Input
                        value={assignUsername}
                        onChange={e => setAssignUsername(e.target.value)}
                        placeholder="username"
                      />
                    </div>
                    <div className="col-span-3">
                      <Button onClick={assignUser} data-test-id={TEST_IDS.tenants.setQuota}>Assign</Button>
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  )
}

