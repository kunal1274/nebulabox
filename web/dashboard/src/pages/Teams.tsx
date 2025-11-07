import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { api } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { TEST_IDS } from '@/lib/test-ids'

export function Teams() {
  const [teams, setTeams] = useState<Array<{ id: string; name: string; description: string; created: string; createdBy: string }>>([])
  const [selectedTeam, setSelectedTeam] = useState<{ id: string; name: string; description: string; created: string; createdBy: string } | null>(null)
  const [members, setMembers] = useState<Array<{ username: string; role: string; joinedAt: string }>>([])
  const [newTeamName, setNewTeamName] = useState('')
  const [newTeamDesc, setNewTeamDesc] = useState('')
  const [inviteUsername, setInviteUsername] = useState('')
  const [inviteRole, setInviteRole] = useState<'admin' | 'editor' | 'viewer'>('viewer')
  const [me, setMe] = useState<{ username: string; role: string } | null>(null)

  useEffect(() => {
    loadTeams()
    api.me().then(r => setMe(r.user || null)).catch(() => {})
  }, [])

  const loadTeams = async () => {
    try {
      const data = await api.listTeams()
      setTeams(data)
    } catch {
      setTeams([])
    }
  }

  const loadTeamDetails = async (id: string) => {
    try {
      const data = await api.getTeam(id)
      setSelectedTeam(data.team)
      setMembers(data.members)
    } catch {
      setSelectedTeam(null)
      setMembers([])
    }
  }

  const createTeam = async () => {
    if (!newTeamName.trim()) return
    try {
      await api.createTeam(newTeamName.trim(), newTeamDesc.trim() || undefined)
      setNewTeamName('')
      setNewTeamDesc('')
      loadTeams()
    } catch (err: any) {
      alert(err.message || 'Failed to create team')
    }
  }

  const deleteTeam = async (id: string) => {
    if (!confirm('Delete this team? This will not delete containers/networks, but will remove access.')) return
    try {
      await api.deleteTeam(id)
      if (selectedTeam?.id === id) { setSelectedTeam(null); setMembers([]) }
      loadTeams()
    } catch (err: any) {
      alert(err.message || 'Failed to delete team')
    }
  }

  const invite = async () => {
    if (!selectedTeam || !inviteUsername.trim()) return
    try {
      await api.inviteMember(selectedTeam.id, inviteUsername.trim(), inviteRole)
      setInviteUsername('')
      loadTeamDetails(selectedTeam.id)
    } catch (err: any) {
      alert(err.message || 'Failed to invite member')
    }
  }

  const removeMember = async (username: string) => {
    if (!selectedTeam) return
    if (!confirm(`Remove ${username} from team?`)) return
    try {
      await api.removeMember(selectedTeam.id, username)
      loadTeamDetails(selectedTeam.id)
    } catch (err: any) {
      alert(err.message || 'Failed to remove member')
    }
  }

  const updateRole = async (username: string, role: 'admin' | 'editor' | 'viewer') => {
    if (!selectedTeam) return
    try {
      await api.updateMemberRole(selectedTeam.id, username, role)
      loadTeamDetails(selectedTeam.id)
    } catch (err: any) {
      alert(err.message || 'Failed to update role')
    }
  }

  const isTeamAdmin = selectedTeam && me && members.find(m => m.username === me.username && m.role === 'admin')

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Teams & Workspaces</h1>
        <p className="text-muted-foreground">Organize resources into shared workspaces</p>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div>
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>Create Team</CardTitle>
              <CardDescription>Create a new workspace</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div>
                  <div className="text-xs mb-1">Name</div>
                  <Input value={newTeamName} onChange={e => setNewTeamName(e.target.value)} placeholder="my-team" data-test-id={TEST_IDS.teams.teamName} />
                </div>
                <div>
                  <div className="text-xs mb-1">Description (optional)</div>
                  <Input value={newTeamDesc} onChange={e => setNewTeamDesc(e.target.value)} placeholder="Team description" />
                </div>
                <Button onClick={createTeam} data-test-id={TEST_IDS.teams.createTeam}>Create</Button>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>My Teams</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {teams.length === 0 && <div className="text-sm text-muted-foreground">No teams yet</div>}
                {teams.map(t => (
                  <div
                    key={t.id}
                    className={`p-3 border rounded cursor-pointer hover:bg-accent ${selectedTeam?.id === t.id ? 'bg-accent' : ''}`}
                    onClick={() => loadTeamDetails(t.id)}
                  >
                    <div className="font-medium">{t.name}</div>
                    {t.description && <div className="text-xs text-muted-foreground">{t.description}</div>}
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {selectedTeam && (
          <Card>
            <CardHeader>
              <CardTitle>{selectedTeam.name}</CardTitle>
              <CardDescription>
                {selectedTeam.description || 'No description'}
                {isTeamAdmin && (
                  <Button variant="destructive" size="sm" className="ml-2" onClick={() => deleteTeam(selectedTeam.id)}>
                    Delete Team
                  </Button>
                )}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="mb-4">
                <div className="text-sm font-medium mb-2">Members</div>
                <div className="space-y-2 mb-4">
                  {members.map(m => (
                    <div key={m.username} className="flex items-center justify-between p-2 border rounded text-sm">
                      <div>
                        <span className="font-medium">{m.username}</span>
                        <Badge variant="outline" className="ml-2">{m.role}</Badge>
                      </div>
                      {isTeamAdmin && m.username !== me?.username && (
                        <div className="space-x-1">
                          <select
                            className="text-xs border rounded px-2 py-1"
                            value={m.role}
                            onChange={e => updateRole(m.username, e.target.value as any)}
                          >
                            <option value="admin">Admin</option>
                            <option value="editor">Editor</option>
                            <option value="viewer">Viewer</option>
                          </select>
                          <Button variant="outline" size="sm" onClick={() => removeMember(m.username)} data-test-id={TEST_IDS.teams.addMember}>Remove</Button>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </div>

              {isTeamAdmin && (
                <div className="border-t pt-4">
                  <div className="text-sm font-medium mb-2">Invite Member</div>
                  <div className="grid grid-cols-12 gap-2">
                    <div className="col-span-5">
                      <Input
                        value={inviteUsername}
                        onChange={e => setInviteUsername(e.target.value)}
                        placeholder="username"
                        data-test-id={TEST_IDS.teams.teamName}
                      />
                    </div>
                    <div className="col-span-4">
                      <select
                        className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                        value={inviteRole}
                        onChange={e => setInviteRole(e.target.value as any)}
                      >
                        <option value="viewer">Viewer</option>
                        <option value="editor">Editor</option>
                        <option value="admin">Admin</option>
                      </select>
                    </div>
                    <div className="col-span-3">
                      <Button onClick={invite} data-test-id={TEST_IDS.teams.addMember}>Invite</Button>
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

