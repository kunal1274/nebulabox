import { useEffect, useState } from 'react'
import { Layers, Plus, Trash2, Play, Square, GitBranch, FolderTree, X } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { api, ContainerGroup, GroupHierarchy, ContainerTree, CreateGroupRequest, SharedResources } from '@/lib/api'
import { TEST_IDS } from '@/lib/test-ids'

export function ContainerGroups() {
  const [activeTab, setActiveTab] = useState('groups')
  const [groups, setGroups] = useState<ContainerGroup[]>([])
  const [selectedGroup, setSelectedGroup] = useState<ContainerGroup | null>(null)
  const [hierarchy, setHierarchy] = useState<GroupHierarchy | null>(null)
  const [containerTree, setContainerTree] = useState<ContainerTree | null>(null)

  // Group creation form
  const [groupName, setGroupName] = useState('')
  const [groupDescription, setGroupDescription] = useState('')
  const [parentGroupId, setParentGroupId] = useState('')
  const [sharedNetwork, setSharedNetwork] = useState('')
  const [creating, setCreating] = useState(false)

  // Relationship form
  const [parentId, setParentId] = useState('')
  const [childId, setChildId] = useState('')
  const [relType, setRelType] = useState('dependency')
  const [creatingRel, setCreatingRel] = useState(false)

  // Container selection for tree view
  const [treeContainerId, setTreeContainerId] = useState('')

  useEffect(() => {
    loadGroups()
  }, [])

  const loadGroups = async () => {
    try {
      const data = await api.listGroups()
      setGroups(data.groups)
    } catch (error) {
      console.error('Failed to load groups:', error)
    }
  }

  const loadHierarchy = async (groupId: string) => {
    try {
      const h = await api.getGroupHierarchy(groupId)
      setHierarchy(h)
    } catch (error) {
      console.error('Failed to load hierarchy:', error)
    }
  }

  const loadContainerTree = async (containerId: string) => {
    if (!containerId.trim()) {
      setContainerTree(null)
      return
    }
    try {
      const tree = await api.getContainerTree(containerId)
      setContainerTree(tree)
    } catch (error) {
      console.error('Failed to load container tree:', error)
      setContainerTree(null)
    }
  }

  const handleCreateGroup = async () => {
    if (!groupName.trim()) {
      alert('Please enter a group name')
      return
    }

    setCreating(true)
    try {
      const sharedResources: SharedResources | undefined = sharedNetwork ? {
        network: sharedNetwork,
      } : undefined

      const req: CreateGroupRequest = {
        name: groupName,
        description: groupDescription || undefined,
        parentGroupId: parentGroupId || undefined,
        sharedResources,
      }

      await api.createGroup(req)
      setGroupName('')
      setGroupDescription('')
      setParentGroupId('')
      setSharedNetwork('')
      loadGroups()
    } catch (error: any) {
      console.error('Failed to create group:', error)
      alert(`Failed to create group: ${error.message}`)
    } finally {
      setCreating(false)
    }
  }

  const handleDeleteGroup = async (groupId: string) => {
    if (!confirm('Are you sure you want to delete this group?')) {
      return
    }

    try {
      await api.deleteGroup(groupId, true)
      loadGroups()
      if (selectedGroup?.id === groupId) {
        setSelectedGroup(null)
        setHierarchy(null)
      }
    } catch (error: any) {
      console.error('Failed to delete group:', error)
      alert(`Failed to delete group: ${error.message}`)
    }
  }

  const handleStartGroup = async (groupId: string) => {
    try {
      await api.startGroup(groupId)
      alert('Group start initiated')
    } catch (error: any) {
      console.error('Failed to start group:', error)
      alert(`Failed to start group: ${error.message}`)
    }
  }

  const handleStopGroup = async (groupId: string) => {
    try {
      await api.stopGroup(groupId)
      alert('Group stop initiated')
    } catch (error: any) {
      console.error('Failed to stop group:', error)
      alert(`Failed to stop group: ${error.message}`)
    }
  }

  const handleCreateRelationship = async () => {
    if (!parentId.trim() || !childId.trim()) {
      alert('Please enter both parent and child container IDs')
      return
    }

    setCreatingRel(true)
    try {
      await api.createRelationship(parentId, childId, relType)
      setParentId('')
      setChildId('')
      if (treeContainerId) {
        loadContainerTree(treeContainerId)
      }
    } catch (error: any) {
      console.error('Failed to create relationship:', error)
      alert(`Failed to create relationship: ${error.message}`)
    } finally {
      setCreatingRel(false)
    }
  }

  const handleAddContainer = async (groupId: string) => {
    const containerId = prompt('Enter container ID to add:')
    if (!containerId) return

    try {
      await api.addContainerToGroup(groupId, containerId)
      loadGroups()
      if (selectedGroup?.id === groupId) {
        const updated = await api.getGroup(groupId)
        setSelectedGroup(updated)
        loadHierarchy(groupId)
      }
    } catch (error: any) {
      console.error('Failed to add container:', error)
      alert(`Failed to add container: ${error.message}`)
    }
  }

  const handleRemoveContainer = async (groupId: string, containerId: string) => {
    try {
      await api.removeContainerFromGroup(groupId, containerId)
      loadGroups()
      if (selectedGroup?.id === groupId) {
        const updated = await api.getGroup(groupId)
        setSelectedGroup(updated)
        loadHierarchy(groupId)
      }
    } catch (error: any) {
      console.error('Failed to remove container:', error)
      alert(`Failed to remove container: ${error.message}`)
    }
  }

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleString()
    } catch {
      return dateStr
    }
  }

  const renderHierarchy = (h: GroupHierarchy, level = 0): JSX.Element => {
    const indent = level * 20
    return (
      <div key={h.group.id} className="ml-4" style={{ marginLeft: `${indent}px` }}>
        <div className="p-3 border rounded-lg mb-2 bg-blue-50 dark:bg-blue-900/20">
          <div className="flex items-center justify-between">
            <div>
              <div className="font-semibold flex items-center gap-2">
                <FolderTree className="h-4 w-4" />
                {h.group.name}
              </div>
              <div className="text-sm text-muted-foreground">
                {h.group.description || 'No description'}
              </div>
              <div className="text-xs text-muted-foreground mt-1">
                {h.containers.length} containers
              </div>
            </div>
          </div>
        </div>
        {h.children.map(child => renderHierarchy(child, level + 1))}
      </div>
    )
  }

  const renderContainerTree = (tree: ContainerTree, level = 0): JSX.Element => {
    const indent = level * 20
    return (
      <div key={tree.containerId} className="mb-2" style={{ marginLeft: `${indent}px` }}>
        <div className="p-2 border rounded bg-gray-50 dark:bg-gray-900">
          <div className="flex items-center gap-2">
            <GitBranch className="h-3 w-3" />
            <span className="font-mono text-sm">{tree.containerId.substring(0, 12)}</span>
            {tree.type && (
              <Badge variant="outline" className="text-xs">{tree.type}</Badge>
            )}
          </div>
        </div>
        {tree.children.map(child => renderContainerTree(child, level + 1))}
      </div>
    )
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Container Groups & Hierarchy</h1>
        <p className="text-muted-foreground">
          Organize containers into groups with shared resources and parent-child relationships
        </p>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList>
          <TabsTrigger value="groups">
            <Layers className="h-4 w-4 mr-2" />
            Groups
          </TabsTrigger>
          <TabsTrigger value="hierarchy">
            <FolderTree className="h-4 w-4 mr-2" />
            Hierarchy
          </TabsTrigger>
          <TabsTrigger value="relationships">
            <GitBranch className="h-4 w-4 mr-2" />
            Relationships
          </TabsTrigger>
        </TabsList>

        {/* Groups Tab */}
        <TabsContent value="groups" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Plus className="h-5 w-5" />
                Create Group
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Group Name *</Label>
                  <Input
                    value={groupName}
                    onChange={(e) => setGroupName(e.target.value)}
                    placeholder="my-app-group"
                    data-test-id={TEST_IDS.groups.groupName}
                  />
                </div>
                <div>
                  <Label>Parent Group ID (optional)</Label>
                  <Input
                    value={parentGroupId}
                    onChange={(e) => setParentGroupId(e.target.value)}
                    placeholder="parent-group-id"
                  />
                </div>
                <div className="col-span-2">
                  <Label>Description</Label>
                  <Input
                    value={groupDescription}
                    onChange={(e) => setGroupDescription(e.target.value)}
                    placeholder="Group description"
                  />
                </div>
                <div>
                  <Label>Shared Network (optional)</Label>
                  <Input
                    value={sharedNetwork}
                    onChange={(e) => setSharedNetwork(e.target.value)}
                    placeholder="shared-network"
                  />
                </div>
              </div>
              <Button onClick={handleCreateGroup} disabled={creating} className="mt-4" data-test-id={TEST_IDS.groups.createGroup}>
                <Plus className="h-4 w-4 mr-2" />
                {creating ? 'Creating...' : 'Create Group'}
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Container Groups ({groups.length})</CardTitle>
            </CardHeader>
            <CardContent>
              {groups.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Layers className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>No groups created. Create one to get started.</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {groups.map((group) => (
                    <div key={group.id} className="p-4 border rounded-lg">
                      <div className="flex items-start justify-between mb-2">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <span className="font-semibold text-lg">{group.name}</span>
                            {group.parentGroupId && (
                              <Badge variant="outline">Child of {group.parentGroupId.substring(0, 8)}</Badge>
                            )}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            {group.description || 'No description'}
                          </div>
                          <div className="text-xs text-muted-foreground mt-1">
                            {group.containerIds.length} containers | Created: {formatDate(group.createdAt)}
                          </div>
                          {group.sharedResources && (
                            <div className="mt-2 text-xs">
                              {group.sharedResources.network && (
                                <Badge variant="outline" className="mr-1">Network: {group.sharedResources.network}</Badge>
                              )}
                              {group.sharedResources.volumes && group.sharedResources.volumes.length > 0 && (
                                <Badge variant="outline" className="mr-1">{group.sharedResources.volumes.length} volumes</Badge>
                              )}
                            </div>
                          )}
                        </div>
                        <div className="flex gap-2">
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => {
                              setSelectedGroup(group)
                              loadHierarchy(group.id)
                            }}
                          >
                            <FolderTree className="h-4 w-4" />
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleStartGroup(group.id)}
                            data-test-id={TEST_IDS.groups.startGroup}
                          >
                            <Play className="h-4 w-4" />
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleStopGroup(group.id)}
                            data-test-id={TEST_IDS.groups.stopGroup}
                          >
                            <Square className="h-4 w-4" />
                          </Button>
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => handleDeleteGroup(group.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                      <div className="mt-2">
                        {group.containerIds.length > 0 ? (
                          <div className="flex flex-wrap gap-1">
                            {group.containerIds.map((cid) => (
                              <Badge key={cid} variant="secondary" className="font-mono text-xs">
                                {cid.substring(0, 12)}
                                <button
                                  onClick={() => handleRemoveContainer(group.id, cid)}
                                  className="ml-1 hover:text-red-500"
                                  data-test-id={TEST_IDS.groups.removeContainer}
                                >
                                  <X className="h-3 w-3" />
                                </button>
                              </Badge>
                            ))}
                          </div>
                        ) : (
                          <div className="text-sm text-muted-foreground">No containers</div>
                        )}
                        <Button
                          size="sm"
                          variant="outline"
                          className="mt-2"
                          onClick={() => handleAddContainer(group.id)}
                          data-test-id={TEST_IDS.groups.addContainer}
                        >
                          <Plus className="h-3 w-3 mr-1" />
                          Add Container
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Hierarchy Tab */}
        <TabsContent value="hierarchy" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Group Hierarchy</CardTitle>
            </CardHeader>
            <CardContent>
              {selectedGroup ? (
                <div>
                  <div className="mb-4">
                    <Button
                      variant="outline"
                      onClick={() => {
                        setSelectedGroup(null)
                        setHierarchy(null)
                      }}
                    >
                      Clear Selection
                    </Button>
                  </div>
                  {hierarchy && renderHierarchy(hierarchy)}
                </div>
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <FolderTree className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>Select a group from the Groups tab to view its hierarchy</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Relationships Tab */}
        <TabsContent value="relationships" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <GitBranch className="h-5 w-5" />
                Create Container Relationship
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <Label>Parent Container ID *</Label>
                  <Input
                    value={parentId}
                    onChange={(e) => setParentId(e.target.value)}
                    placeholder="parent-container-id"
                  />
                </div>
                <div>
                  <Label>Child Container ID *</Label>
                  <Input
                    value={childId}
                    onChange={(e) => setChildId(e.target.value)}
                    placeholder="child-container-id"
                  />
                </div>
                <div className="col-span-2">
                  <Label>Relationship Type</Label>
                  <select
                    value={relType}
                    onChange={(e) => setRelType(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md"
                  >
                    <option value="dependency">Dependency</option>
                    <option value="composition">Composition</option>
                    <option value="nested">Nested</option>
                  </select>
                </div>
              </div>
              <Button onClick={handleCreateRelationship} disabled={creatingRel} data-test-id={TEST_IDS.groups.createGroup}>
                <Plus className="h-4 w-4 mr-2" />
                {creatingRel ? 'Creating...' : 'Create Relationship'}
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Container Tree View</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="mb-4">
                <Label>Container ID</Label>
                <div className="flex gap-2">
                  <Input
                    value={treeContainerId}
                    onChange={(e) => setTreeContainerId(e.target.value)}
                    placeholder="container-id"
                    onKeyPress={(e) => e.key === 'Enter' && loadContainerTree(treeContainerId)}
                  />
                  <Button onClick={() => loadContainerTree(treeContainerId)}>
                    Load Tree
                  </Button>
                </div>
              </div>
              {containerTree ? (
                <div className="border rounded-lg p-4 max-h-96 overflow-y-auto">
                  {renderContainerTree(containerTree)}
                </div>
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <GitBranch className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p>Enter a container ID to view its relationship tree</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

