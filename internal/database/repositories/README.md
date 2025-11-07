# Database Repositories

This package implements the repository pattern for database operations.

## Repositories

### ContainerRepository
- `CreateOrUpdate(container, opts)` - Create or update a container
- `Get(id)` - Get container by ID
- `List(all, workspaceID)` - List containers with optional filters
- `UpdateStatus(id, status)` - Update container status
- `Delete(id)` - Soft delete a container
- `AssociateWorkspace(containerID, workspaceID)` - Associate container with workspace
- `GetStorageOptions(id)` - Get additional container storage options

### ImageRepository
- `CreateOrUpdate(image)` - Create or update an image
- `Get(id)` - Get image by ID
- `List()` - List all images
- `FindByNameAndTag(name, tag)` - Find image by name and tag
- `Delete(id)` - Soft delete an image

### WorkspaceRepository
- `Create(workspace)` - Create a workspace
- `Get(id)` - Get workspace by ID
- `List(userID)` - List workspaces for a user
- `UpdateStatus(id, status)` - Update workspace status
- `Delete(id)` - Soft delete a workspace
- `AddMember(workspaceID, userID, role)` - Add member to workspace
- `RemoveMember(workspaceID, userID)` - Remove member from workspace

## Usage

```go
// Initialize repositories (after database is initialized)
repos, err := repositories.InitRepositories()
if err != nil {
    // Handle error or use nil (fallback to in-memory)
}

// Use repositories
if repos != nil && repos.Container != nil {
    container, err := repos.Container.Get("container-id")
    // ...
}
```

## Graceful Fallback

If repositories are nil (database not available), the system falls back to in-memory storage. Always check for nil before using repositories.

