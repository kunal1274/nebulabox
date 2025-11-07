# Nebula Registry

The Nebula Registry is a Docker Registry HTTP API V2 compatible private registry server for managing container images.

## Features

- **Docker Registry V2 API Compatible**: Full support for standard Docker/OCI registry operations
- **Persistent Storage**: File-based storage with automatic metadata persistence
- **Authentication**: Token-based and Basic authentication support
- **Image Versioning**: Semantic versioning with metadata tracking
- **RESTful API**: Additional management endpoints for repositories and versions
- **Multi-tenant Ready**: Supports user authentication and authorization

## Quick Start

### Build Registry

```bash
make build-registry
```

### Run Registry

```bash
make registry
```

Or manually:

```bash
./nebulabox-registry
```

### Environment Variables

- `NEBULABOX_REGISTRY_PORT`: Registry port (default: 5001)
- `NEBULABOX_REGISTRY_STORAGE`: Storage directory (default: ./registry-storage)

### Testing

```bash
# Run basic registry tests
make registry-test

# Run registry unit tests
make registry-test-unit

# Run integration tests (requires API server)
make registry-test-integration

# Check registry health
make registry-health
```

## API Endpoints

### Docker Registry V2 API

- `GET /v2/` - Registry version check
- `GET /v2/_catalog` - List all repositories
- `GET /v2/<repo>/tags/list` - List tags for a repository
- `POST /v2/<repo>/blobs/uploads/` - Start blob upload
- `PATCH /v2/<repo>/blobs/uploads/<uuid>` - Upload blob chunk
- `PUT /v2/<repo>/blobs/uploads/<uuid>?digest=<digest>` - Finalize blob
- `GET /v2/<repo>/blobs/<digest>` - Get blob
- `PUT /v2/<repo>/manifests/<tag>` - Put manifest
- `GET /v2/<repo>/manifests/<tag>` - Get manifest
- `DELETE /v2/<repo>/manifests/<tag>` - Delete tag

### Authentication

- `POST /auth/login` - Login and get token
  ```json
  {
    "username": "admin",
    "password": "admin"
  }
  ```
  Response:
  ```json
  {
    "token": "nb-...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "scope": "pull push"
  }
  ```

- `POST /auth/token` - Generate token with custom scopes

### Management API

- `GET /api/registry/repositories` - List all repositories
- `GET /api/registry/repositories/:repo/versions` - List versions for a repo
- `GET /api/registry/repositories/:repo/versions/:tag` - Get version info
- `DELETE /api/registry/repositories/:repo/versions/:tag` - Delete version
- `GET /api/registry/repositories/:repo/summary` - Get repository summary

## Storage

The registry supports two storage backends:

### File-based Storage (Default)

- Stores blobs, manifests, and metadata on disk
- Configurable via `NEBULABOX_REGISTRY_STORAGE`
- Automatically persists version metadata

Storage structure:
```
registry-storage/
├── metadata.json          # Version metadata
└── repos/
    └── <repo>/
        ├── blobs/
        │   └── <digest>
        └── manifests/
            └── <tag>
```

### In-Memory Storage

- Fallback when file storage fails
- Data is lost on restart
- Suitable for testing

## Authentication

### Default Users

- Username: `admin`
- Password: `admin`
- Roles: `admin`, `push`, `pull`

### Adding Users

Users can be added programmatically via the `AuthConfig` API:

```go
auth := NewAuthConfig()
auth.AddUser("username", "password", []string{"pull", "push"})
```

### Token Authentication

Tokens are Bearer tokens with expiration (default: 24 hours).

Example:
```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:5001/v2/_catalog
```

### Basic Authentication

Supports HTTP Basic Auth:

```bash
curl -u admin:admin \
  http://localhost:5001/v2/_catalog
```

## Image Versioning

The registry tracks version metadata for each image:

- **Tag**: Image tag (e.g., `latest`, `v1.0.0`)
- **Digest**: Content-addressable digest
- **CreatedAt**: Timestamp when version was created
- **CreatedBy**: Username who created the version
- **Size**: Image size in bytes
- **Metadata**: Additional custom metadata
- **Description**: Optional description

### Semantic Versioning

The registry automatically detects semantic versions (e.g., `v1.0.0`, `v2.3.1`) and tracks the latest version.

### Version Comparison

```go
// Compare two versions
result := CompareVersions("v1.2.3", "v1.2.4")
// Returns: -1 (v1.2.3 < v1.2.4)
```

## Integration with API Server

The registry is automatically integrated with the main API server:

1. API server connects to registry on startup
2. Auto-authenticates using `NEBULABOX_ADMIN_USER` and `NEBULABOX_ADMIN_PASS`
3. Image listing pulls from registry
4. Push operations route to registry

### Configuration

Set in API server:
```bash
NEBULABOX_REGISTRY_URL=http://localhost:5001
NEBULABOX_ADMIN_USER=admin
NEBULABOX_ADMIN_PASS=admin
```

## Examples

### Push an Image (via Docker)

```bash
# Login
docker login localhost:5001

# Tag image
docker tag myimage:latest localhost:5001/myimage:v1.0.0

# Push
docker push localhost:5001/myimage:v1.0.0
```

### List Repositories

```bash
curl http://localhost:5001/v2/_catalog
```

### Get Version Information

```bash
curl http://localhost:5001/api/registry/repositories/myimage/versions
```

### Delete a Version

```bash
curl -X DELETE \
  -H "Authorization: Bearer <token>" \
  http://localhost:5001/api/registry/repositories/myimage/versions/v1.0.0
```

## Troubleshooting

### Registry Won't Start

1. Check if port is in use:
   ```bash
   lsof -i :5001
   ```

2. Check storage directory permissions:
   ```bash
   ls -la registry-storage/
   ```

3. Check logs:
   ```bash
   ./nebulabox-registry
   ```

### Authentication Fails

1. Verify credentials:
   ```bash
   curl -X POST http://localhost:5001/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"admin"}'
   ```

2. Check token expiration (default: 24 hours)

### Storage Issues

1. Verify storage directory exists and is writable
2. Check disk space
3. Review `metadata.json` for corruption

## Development

### Running Tests

```bash
# Unit tests
make registry-test-unit

# Integration tests
make registry-test-integration

# Manual testing
make registry-test
```

### Adding Features

The registry uses a pluggable storage interface:

```go
type Storage interface {
    ListRepositories() []string
    ListTags(repo string) []string
    GetBlob(repo, digest string) (*blobData, bool)
    PutBlob(repo, digest string, data []byte) error
    // ... more methods
}
```

Implement this interface to add custom storage backends (e.g., S3, database).

## Security Considerations

1. **Production**: Hash passwords using bcrypt or similar
2. **HTTPS**: Use TLS in production
3. **Token Expiration**: Configure appropriate token lifetimes
4. **Access Control**: Implement proper RBAC for multi-user scenarios
5. **Storage Encryption**: Encrypt stored blobs in production

## See Also

- [Docker Registry HTTP API V2 Specification](https://docs.docker.com/registry/spec/api/)
- [OCI Distribution Specification](https://github.com/opencontainers/distribution-spec)
- [CI/CD Documentation](./CI_CD.md)

