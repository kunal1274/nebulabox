# Schema-First Architecture

NebulaBox uses a **schema-first approach** where all UI, API, and Database derive from a single schema file. This ensures consistency, type safety, and enables powerful features like auto-generated SDKs, live testing hooks, and stateful collaboration.

## Overview

```
schema/nebulabox.schema.json
    ├── definitions/          → Type definitions (Container, Image, Workspace)
    ├── api/                  → API endpoint specifications
    ├── database/             → Database schema definitions
    ├── ui/                   → UI component specifications
    ├── testing/              → Testing hooks and coverage
    └── collaboration/        → Real-time collaboration features
```

## Key Benefits

### 1. Single Source of Truth
- One schema file defines everything
- No manual synchronization needed
- Automatic consistency across layers

### 2. Auto-Generated Code
- **TypeScript types** for frontend
- **Go types** for backend
- **API client SDK** (tRPC/Blitz-like)
- **Database migrations**
- **UI components**

### 3. Type Safety
- Compile-time validation
- Auto-complete in IDEs
- Runtime validation
- No type mismatches

### 4. Developer Experience
- Instant feedback
- Auto-generated documentation
- Live testing hooks
- Real-time collaboration

## Architecture Flow

```
┌─────────────────────────────────────────────────────────┐
│         schema/nebulabox.schema.json                    │
│         (Single Source of Truth)                        │
└─────────────────────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
        ▼           ▼           ▼
    ┌────────┐ ┌────────┐ ┌────────┐
    │   UI   │ │  API   │ │   DB   │
    └────────┘ └────────┘ └────────┘
        │           │           │
        └───────────┼───────────┘
                    │
        ┌───────────▼───────────┐
        │   Generated Code      │
        │   - Types             │
        │   - SDKs              │
        │   - Components        │
        └───────────────────────┘
```

## Usage

### Generate Code from Schema

```bash
# Generate all code
npm run generate

# Generate specific outputs
npm run generate:types      # TypeScript + Go types
npm run generate:api        # API client SDK
npm run generate:db         # Database migrations
npm run generate:ui         # UI components

# Watch mode (auto-generate on schema changes)
npm run generate:watch
```

### Validate Schema

```bash
npm run validate:schema
```

### Live Testing Hooks

```bash
# Start live testing (runs tests on file save)
npm run test:live
```

### Collaboration Mode

```bash
# Enable real-time collaboration
npm run collaborate
```

## Schema Structure

### Definitions

Type definitions for all entities:

```json
{
  "definitions": {
    "Container": {
      "type": "object",
      "properties": {
        "id": { "type": "string", "format": "uuid" },
        "name": { "type": "string" },
        "status": { "type": "string", "enum": ["running", "stopped"] }
      }
    }
  }
}
```

### API Endpoints

API endpoint specifications:

```json
{
  "api": {
    "endpoints": {
      "containers": {
        "list": {
          "method": "GET",
          "path": "/api/containers",
          "response": {
            "type": "array",
            "items": { "$ref": "#/definitions/Container" }
          }
        }
      }
    }
  }
}
```

### Database Schema

Database table definitions:

```json
{
  "database": {
    "tables": {
      "containers": {
        "primaryKey": "id",
        "columns": {
          "id": { "type": "string", "dbType": "VARCHAR(255)" },
          "name": { "type": "string", "dbType": "VARCHAR(255)" }
        }
      }
    }
  }
}
```

### UI Components

UI component specifications:

```json
{
  "ui": {
    "components": {
      "ContainerList": {
        "dataSource": "/api/containers",
        "fields": ["id", "name", "status"],
        "actions": ["start", "stop", "delete"]
      }
    }
  }
}
```

## Live Testing Hooks

The schema includes testing hooks that automatically run tests when files are saved:

```json
{
  "testing": {
    "hooks": {
      "onSave": {
        "triggers": ["unit", "integration", "e2e"],
        "scope": "changed-files",
        "timeout": 30000
      }
    }
  }
}
```

### Features

- **Automatic test runs** on file save
- **Incremental testing** (only changed files)
- **Fast feedback** during development
- **Configurable scope** (changed-files, related-files, all)

## Stateful Collaboration

Real-time collaboration like Cursor/Replit:

```json
{
  "collaboration": {
    "features": {
      "realtime": {
        "enabled": true,
        "protocol": "websocket",
        "sync": {
          "type": "operational-transform",
          "conflictResolution": "last-write-wins"
        }
      }
    }
  }
}
```

### Features

- **Real-time file synchronization**
- **Presence indicators** (cursors, selections)
- **Operational Transform** for conflict resolution
- **Multi-user collaboration** on same files

## Auto-Generated SDKs

Similar to tRPC or Blitz resolver:

```typescript
// Auto-generated from schema
import { NebulaBoxAPI } from './generated/api-client';

const api = new NebulaBoxAPI(client);

// Fully typed, auto-complete enabled
const containers = await api.containersList({ all: true });
const container = await api.containersGet({ id: '123' });
await api.containersStart({ id: '123' });
```

## Migration Path

### Current State
- ✅ Schema definition created
- ✅ Code generation scripts
- ✅ Live testing hooks
- ✅ Collaboration sync layer

### Next Steps
1. **Integration** - Wire up generated code to existing codebase
2. **Validation** - Add runtime validation from schema
3. **UI Components** - Generate React components from schema
4. **API Routes** - Auto-generate API routes from schema
5. **Database Migrations** - Auto-generate migrations from schema

## Versioning

Schema follows semantic versioning:
- **Major**: Breaking changes (require migration)
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes (backward compatible)

Current version: `1.0.0`

## Best Practices

1. **Always validate** schema before generating code
2. **Version control** schema changes carefully
3. **Run tests** after schema changes
4. **Update documentation** when schema changes
5. **Use schema-first** for all new features

## Future Enhancements

- [ ] GraphQL schema generation
- [ ] OpenAPI/Swagger generation
- [ ] TypeScript strict mode validation
- [ ] Visual schema editor
- [ ] Schema diff tooling
- [ ] Migration generator

