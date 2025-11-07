import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { api } from '@/lib/api'

export function Settings() {
  const [events, setEvents] = useState<Array<{ id:string; event:string; repo:string; ref:string; action?:string; sender?:string; timestamp:number }>>([])
  const [glEvents, setGlEvents] = useState<Array<{ id:string; event:string; project:string; ref:string; user:string; timestamp:number }>>([])
  const apiBase = (import.meta as any).env.VITE_API_URL || 'http://localhost:8081/api'
  const [builds, setBuilds] = useState<Array<{ id:string; source:string; repo:string; ref:string; status:string; startedAt:number; endedAt:number }>>([])
  const [manualRepo, setManualRepo] = useState('')
  const [manualRef, setManualRef] = useState('main')
  const [tests, setTests] = useState<Array<{ id:string; repo:string; ref:string; suite:string; status:string; startedAt:number; endedAt:number }>>([])
  const [manualTestSuite, setManualTestSuite] = useState('default')
  const [deployments, setDeployments] = useState<Array<{ id:string; repo:string; ref:string; env:string; status:string; startedAt:number; endedAt:number }>>([])
  const [manualEnv, setManualEnv] = useState('staging')
  const [rollbacks, setRollbacks] = useState<Array<{ id:string; repo:string; fromRef:string; toRef:string; env:string; status:string; startedAt:number; endedAt:number }>>([])
  const [me, setMe] = useState<{ username: string; role: string } | null>(null)
  const [loginUser, setLoginUser] = useState('admin')
  const [loginPass, setLoginPass] = useState('admin')
  useEffect(()=>{
    const load = async ()=>{ try { const r = await api.getGitHubEvents(); setEvents(r.events||[]) } catch { setEvents([]) } }
    const loadGL = async ()=>{ try { const r = await api.getGitLabEvents(); setGlEvents(r.events||[]) } catch { setGlEvents([]) } }
    const loadBuilds = async ()=>{ try { const r = await api.getBuilds(); setBuilds(r.builds||[]) } catch { setBuilds([]) } }
    const loadTests = async ()=>{ try { const r = await api.getTests(); setTests(r.tests||[]) } catch { setTests([]) } }
    const loadDeploys = async ()=>{ try { const r = await api.getDeployments(); setDeployments(r.deployments||[]) } catch { setDeployments([]) } }
    const loadRollbacks = async ()=>{ try { const r = await api.getRollbacks(); setRollbacks(r.rollbacks||[]) } catch { setRollbacks([]) } }
    const loadMe = async ()=>{ try { const r = await api.me(); setMe(r.user || null) } catch { setMe(null) } }
    load(); loadGL(); loadBuilds(); loadTests(); loadDeploys(); loadRollbacks(); loadMe(); const t = setInterval(()=>{ load(); loadGL(); loadBuilds(); loadTests(); loadDeploys(); loadRollbacks(); loadMe() }, 5000); return ()=>clearInterval(t)
  }, [])
  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Settings</h1>
        <p className="text-muted-foreground">
          Configure NebulaBox dashboard
        </p>
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>General Settings</CardTitle>
          <CardDescription>
            Application configuration
          </CardDescription>
        </CardHeader>
        <CardContent>
      <div className="text-sm">Logged in: {me ? <><span className="font-medium">{me.username}</span> <span className="text-muted-foreground">({me.role})</span></> : <span className="text-muted-foreground">anonymous</span>}</div>
        </CardContent>
      </Card>

  <Card className="mb-6">
    <CardHeader>
      <CardTitle>Authentication</CardTitle>
      <CardDescription>Login to enable protected features</CardDescription>
    </CardHeader>
    <CardContent>
      <div className="grid grid-cols-12 gap-2 items-end">
        <div className="col-span-4">
          <div className="text-xs mb-1">Username</div>
          <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={loginUser} onChange={e=>setLoginUser(e.target.value)} data-test-id={TEST_IDS.settings.apiUrl} />
        </div>
        <div className="col-span-4">
          <div className="text-xs mb-1">Password</div>
          <input type="password" className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={loginPass} onChange={e=>setLoginPass(e.target.value)} />
        </div>
        <div className="col-span-2">
          <button className="h-10 px-3 rounded-md border text-sm" onClick={async()=>{ try { await api.login(loginUser, loginPass); const r = await api.me(); setMe(r.user||null) } catch { alert('Login failed') } }} data-test-id={TEST_IDS.settings.saveSettings}>Login</button>
        </div>
        <div className="col-span-2">
          <button className="h-10 px-3 rounded-md border text-sm" onClick={async()=>{ try { await api.logout(); setMe(null) } catch {} }} data-test-id={TEST_IDS.settings.saveSettings}>Logout</button>
        </div>
      </div>
      <div className="text-xs text-muted-foreground mt-2">Default credentials: admin/admin (override with NEBULABOX_ADMIN_USER/PASS)</div>
    </CardContent>
  </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Builds</CardTitle>
          <CardDescription>Automated builds on push and manual triggers</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end mb-3">
            <div className="col-span-5">
              <div className="text-xs mb-1">Repository</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualRepo} onChange={e=>setManualRepo(e.target.value)} placeholder="org/repo" />
            </div>
            <div className="col-span-3">
              <div className="text-xs mb-1">Ref</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualRef} onChange={e=>setManualRef(e.target.value)} placeholder="main" />
            </div>
            <div className="col-span-2">
              <button className="h-10 px-3 rounded-md border text-sm" onClick={async()=>{ if (!manualRepo.trim()) return; await api.triggerBuild({ source:'manual', repo: manualRepo.trim(), ref: manualRef.trim() }); setManualRepo(''); }} data-test-id={TEST_IDS.settings.saveSettings}>Trigger</button>
            </div>
          </div>
          <div className="border rounded divide-y">
            {builds.length===0 && (<div className="p-3 text-sm text-muted-foreground">No builds yet</div>)}
            {builds.map(b=> (
              <div key={b.id} className="p-3 text-sm flex justify-between">
                <div>{b.repo}@{b.ref} <span className="text-muted-foreground">({b.source})</span></div>
                <div className="text-muted-foreground">{b.status}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Tests</CardTitle>
          <CardDescription>Automated test runs (post-build) and manual triggers</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end mb-3">
            <div className="col-span-5">
              <div className="text-xs mb-1">Repository</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualRepo} onChange={e=>setManualRepo(e.target.value)} placeholder="org/repo" />
            </div>
            <div className="col-span-3">
              <div className="text-xs mb-1">Ref</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualRef} onChange={e=>setManualRef(e.target.value)} placeholder="main" />
            </div>
            <div className="col-span-2">
              <div className="text-xs mb-1">Suite</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualTestSuite} onChange={e=>setManualTestSuite(e.target.value)} placeholder="default" />
            </div>
            <div className="col-span-2">
              <button className="h-10 px-3 rounded-md border text-sm" onClick={async()=>{ if (!manualRepo.trim()) return; await api.runTests({ repo: manualRepo.trim(), ref: manualRef.trim(), suite: manualTestSuite.trim() }); }} data-test-id={TEST_IDS.settings.saveSettings}>Run</button>
            </div>
          </div>
          <div className="border rounded divide-y">
            {tests.length===0 && (<div className="p-3 text-sm text-muted-foreground">No tests yet</div>)}
            {tests.map(t=> (
              <div key={t.id} className="p-3 text-sm flex justify-between">
                <div>{t.repo}@{t.ref} · <span className="text-muted-foreground">{t.suite}</span></div>
                <div className="text-muted-foreground">{t.status}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Deployments</CardTitle>
          <CardDescription>Auto-deploy on green tests; trigger manual deployments</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end mb-3">
            <div className="col-span-5">
              <div className="text-xs mb-1">Repository</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualRepo} onChange={e=>setManualRepo(e.target.value)} placeholder="org/repo" />
            </div>
            <div className="col-span-3">
              <div className="text-xs mb-1">Ref</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualRef} onChange={e=>setManualRef(e.target.value)} placeholder="main" />
            </div>
            <div className="col-span-2">
              <div className="text-xs mb-1">Environment</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualEnv} onChange={e=>setManualEnv(e.target.value)} placeholder="staging" />
            </div>
            <div className="col-span-2">
              <button className="h-10 px-3 rounded-md border text-sm" onClick={async()=>{ if (!manualRepo.trim()) return; await api.triggerDeployment({ repo: manualRepo.trim(), ref: manualRef.trim(), env: manualEnv.trim() }); }} data-test-id={TEST_IDS.settings.saveSettings}>Deploy</button>
            </div>
          </div>
          <div className="border rounded divide-y">
            {deployments.length===0 && (<div className="p-3 text-sm text-muted-foreground">No deployments yet</div>)}
            {deployments.map(d=> (
              <div key={d.id} className="p-3 text-sm flex justify-between">
                <div>{d.repo}@{d.ref} → <span className="text-muted-foreground">{d.env}</span></div>
                <div className="text-muted-foreground">{d.status}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Rollbacks</CardTitle>
          <CardDescription>Rollback to previous successful deployment for an environment</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-12 gap-2 items-end mb-3">
            <div className="col-span-5">
              <div className="text-xs mb-1">Repository</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualRepo} onChange={e=>setManualRepo(e.target.value)} placeholder="org/repo" />
            </div>
            <div className="col-span-3">
              <div className="text-xs mb-1">Environment</div>
              <input className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" value={manualEnv} onChange={e=>setManualEnv(e.target.value)} placeholder="staging" />
            </div>
            <div className="col-span-2">
              <button className="h-10 px-3 rounded-md border text-sm" onClick={async()=>{ if (!manualRepo.trim()||!manualEnv.trim()) return; const res = await api.triggerRollback({ repo: manualRepo.trim(), env: manualEnv.trim() }); if (!(res as any).id) alert('No previous successful deployment to roll back to') }} data-test-id={TEST_IDS.settings.saveSettings}>Rollback</button>
            </div>
          </div>
          <div className="border rounded divide-y">
            {rollbacks.length===0 && (<div className="p-3 text-sm text-muted-foreground">No rollbacks</div>)}
            {rollbacks.map(r=> (
              <div key={r.id} className="p-3 text-sm flex justify-between">
                <div>{r.repo} · {r.env} · {r.fromRef} → {r.toRef}</div>
                <div className="text-muted-foreground">{r.status}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>GitHub Webhook</CardTitle>
          <CardDescription>Configure GitHub to POST to this endpoint. Optionally set NEBULABOX_GITHUB_SECRET on API.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-sm mb-3">Endpoint: <span className="font-mono bg-muted px-2 py-1 rounded">{apiBase}/webhooks/github</span></div>
          <div className="text-sm font-medium mb-2">Recent Events</div>
          <div className="border rounded divide-y">
            {events.length===0 && (<div className="p-3 text-sm text-muted-foreground">No events</div>)}
            {events.map(e=> (
              <div key={e.id+e.timestamp} className="p-3 text-sm flex justify-between">
                <div>{e.event}{e.action?`/${e.action}`:''}</div>
                <div className="text-muted-foreground">{e.repo} · {e.ref} · {e.sender || ''}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>GitLab Webhook</CardTitle>
          <CardDescription>Configure GitLab to POST to this endpoint. Optionally set NEBULABOX_GITLAB_SECRET on API.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-sm mb-3">Endpoint: <span className="font-mono bg-muted px-2 py-1 rounded">{apiBase}/webhooks/gitlab</span></div>
          <div className="text-sm font-medium mb-2">Recent Events</div>
          <div className="border rounded divide-y">
            {glEvents.length===0 && (<div className="p-3 text-sm text-muted-foreground">No events</div>)}
            {glEvents.map(e=> (
              <div key={e.id+e.timestamp} className="p-3 text-sm flex justify-between">
                <div>{e.event}</div>
                <div className="text-muted-foreground">{e.project} · {e.ref} · {e.user}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

