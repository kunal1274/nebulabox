# NebulaBox Build Specification

NebulaBox provides a simpler, more structured alternative to Dockerfiles through the Build Specification format. This JSON/YAML-based format makes it easier to define container images with better validation and tooling support.

## Overview

The Build Specification is a declarative format that describes how to build a container image. It's converted to a Dockerfile internally, making it compatible with standard container tooling while providing a better developer experience.

## Specification Format

### Basic Structure

```json
{
  "version": "1.0",
  "name": "my-app",
  "tag": "my-app:latest",
  "base": {
    "image": "alpine",
    "tag": "3.19"
  },
  "steps": [
    { "type": "run", "command": "apk add --no-cache nodejs npm" }
  ]
}
```

### Required Fields

- **`name`** (string): Application name
- **`base`** (object): Base image configuration
  - **`image`** (string): Base image name (e.g., "alpine", "node", "ubuntu")
  - **`tag`** (string, optional): Base image tag (default: "latest")
- **`steps`** (array): Build steps to execute

### Optional Fields

- **`version`** (string): Specification version (default: "1.0")
- **`tag`** (string): Output image tag (default: "latest")
- **`env`** (object): Environment variables as key-value pairs
- **`workdir`** (string): Working directory path
- **`expose`** (array): Array of port numbers to expose
- **`labels`** (object): Image labels as key-value pairs
- **`health`** (object): Health check configuration
- **`user`** (string): User to run commands as

## Step Types

### `run`

Execute a shell command:

```json
{
  "type": "run",
  "command": "npm install",
  "comment": "Install dependencies"
}
```

- **`command`** (required): Command to execute
- **`comment`** (optional): Step description
- **`workdir`** (optional): Working directory for this step
- **`env`** (optional): Environment variables for this step
- **`user`** (optional): User to run this step as

### `copy`

Copy files from build context:

```json
{
  "type": "copy",
  "source": "package.json",
  "dest": "/app/package.json",
  "comment": "Copy package file"
}
```

- **`source`** (required): Source path (relative to build context)
- **`dest`** (required): Destination path in image
- **`comment`** (optional): Step description

### `add`

Similar to `copy` but supports URLs and archives:

```json
{
  "type": "add",
  "source": "https://example.com/file.tar.gz",
  "dest": "/app",
  "comment": "Download and extract archive"
}
```

### `cmd`

Set the default command to run:

```json
{
  "type": "cmd",
  "command": "[\"node\", \"index.js\"]",
  "comment": "Start application"
}
```

- **`command`** (required): Command in JSON array format or shell string

### `arg`

Define build arguments:

```json
{
  "type": "arg",
  "command": "NODE_VERSION=18",
  "comment": "Node.js version"
}
```

### `volume`

Create a volume mount point:

```json
{
  "type": "volume",
  "dest": "/data",
  "comment": "Data directory"
}
```

## Health Checks

Configure container health checks:

```json
{
  "health": {
    "type": "http",
    "path": "/health",
    "port": 3000,
    "interval": 30,
    "timeout": 10,
    "retries": 3
  }
}
```

### Health Check Types

- **`http`**: HTTP endpoint check
  - **`path`** (required): Health endpoint path
  - **`port`** (required): Port number
- **`tcp`**: TCP port check
  - **`port`** (required): Port number
- **`cmd`**: Custom command check
  - **`command`** (required): Command to execute

### Health Check Options

- **`interval`** (int, optional): Check interval in seconds (default: 30)
- **`timeout`** (int, optional): Timeout in seconds (default: 10)
- **`retries`** (int, optional): Retry count before marking unhealthy (default: 3)

## Complete Example

```json
{
  "version": "1.0",
  "name": "node-api",
  "tag": "myregistry/node-api:v1.0",
  "base": {
    "image": "node",
    "tag": "18-alpine"
  },
  "steps": [
    {
      "type": "run",
      "command": "npm install -g npm@latest",
      "comment": "Update npm"
    },
    {
      "type": "copy",
      "source": "package.json",
      "dest": "/app/package.json",
      "comment": "Copy package file"
    },
    {
      "type": "run",
      "command": "npm install --production",
      "comment": "Install dependencies",
      "workdir": "/app"
    },
    {
      "type": "copy",
      "source": ".",
      "dest": "/app",
      "comment": "Copy application code"
    },
    {
      "type": "cmd",
      "command": "[\"node\", \"server.js\"]",
      "comment": "Start application"
    }
  ],
  "env": {
    "NODE_ENV": "production",
    "PORT": "3000"
  },
  "workdir": "/app",
  "expose": [3000],
  "labels": {
    "maintainer": "NebulaBox",
    "version": "1.0",
    "description": "Node.js API server"
  },
  "health": {
    "type": "http",
    "path": "/health",
    "port": 3000,
    "interval": 30,
    "timeout": 10,
    "retries": 3
  },
  "user": "node"
}
```

## API Endpoints

### Validate Specification

```bash
POST /api/buildspec/validate
Content-Type: application/json

{
  "spec": { ... }
}
```

Response:
```json
{
  "valid": true,
  "dockerfile": "FROM node:18-alpine\n...",
  "message": "Specification is valid"
}
```

### Convert to Dockerfile

```bash
POST /api/buildspec/convert
Content-Type: application/json

{
  "spec": { ... }
}
```

Response:
```json
{
  "valid": true,
  "dockerfile": "FROM node:18-alpine\n...",
  "message": "Converted to Dockerfile"
}
```

### Build from Specification

```bash
POST /api/buildspec/build
Content-Type: application/json

{
  "spec": { ... },
  "tag": "my-app:latest"
}
```

Response:
```json
{
  "valid": true,
  "dockerfile": "FROM node:18-alpine\n...",
  "tag": "my-app:latest",
  "logs": [
    "[+] Building from NebulaBox build spec",
    "..."
  ],
  "message": "Build initiated from specification"
}
```

## Advantages Over Dockerfiles

1. **Structured Format**: JSON/YAML is easier to parse and validate
2. **Better Tooling**: IDE support, autocomplete, validation
3. **Reusability**: Can be templated and versioned easily
4. **Validation**: Early error detection before build time
5. **Documentation**: Built-in comments and metadata
6. **Consistency**: Enforces best practices and structure

## Conversion to Dockerfile

The specification is automatically converted to a standard Dockerfile:

```dockerfile
FROM node:18-alpine
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
LABEL maintainer="NebulaBox"
LABEL version="1.0"
EXPOSE 3000
RUN npm install -g npm@latest
COPY package.json /app/package.json
RUN npm install --production
COPY . /app
CMD ["node", "server.js"]
HEALTHCHECK --interval=30s --timeout=10s --retries=3 CMD curl -f http://localhost:3000/health || exit 1
USER node
```

## Best Practices

1. **Use semantic versioning** for the specification version
2. **Add comments** to steps for documentation
3. **Group related steps** together logically
4. **Use build args** for configurable values
5. **Set appropriate users** for security (avoid root when possible)
6. **Include health checks** for production images
7. **Use .dockerignore** equivalent by being selective with copy steps
8. **Keep base images** updated and minimal
9. **Use multi-stage builds** by specifying multiple base images (future feature)
10. **Version your specs** in version control

## Limitations

- Currently converts to single-stage Dockerfile (multi-stage coming soon)
- Some advanced Dockerfile features may not be fully supported yet
- Requires conversion to Dockerfile before actual build (transparent to user)

## Future Enhancements

- Multi-stage build support
- Template variables and inheritance
- YAML format support (currently JSON-only)
- Visual builder UI
- AI-assisted specification generation
- Integration with CI/CD pipelines
- Specification versioning and migration tools

## See Also

- [Registry Documentation](./REGISTRY.md)
- [Image Building Guide](./IMAGE_BUILDING.md)

