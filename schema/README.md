# NebulaBox Schema-First Architecture

This directory contains the central schema definition that drives all UI, API, and Database generation.

## Overview

The schema-first approach ensures:
- **Single Source of Truth**: One schema file defines everything
- **Type Safety**: Auto-generated types for TypeScript, Go, and database
- **Consistency**: UI, API, and DB always stay in sync
- **Developer Experience**: Auto-complete, validation, and documentation from schema

## Schema Structure

```
nebulabox.schema.json
├── definitions/          # Type definitions (Container, Image, Workspace, etc.)
├── api/                 # API endpoint specifications
├── database/            # Database schema definitions
├── ui/                  # UI component specifications
├── testing/             # Testing hooks and coverage
└── collaboration/       # Real-time collaboration features
```

## Code Generation

### Generated Outputs

1. **TypeScript Types** (`generated/types.ts`)
   - React component props
   - API client types
   - Form validation schemas

2. **Go Types** (`internal/generated/types.go`)
   - API request/response types
   - Database models
   - Validation functions

3. **API Client SDK** (`generated/api-client.ts`)
   - Auto-generated API client with full type safety
   - Similar to tRPC or Blitz resolver pattern
   - Includes request/response validation

4. **Database Migrations** (`internal/database/migrations/`)
   - Auto-generated from schema definitions
   - Indexes, foreign keys, constraints

5. **UI Components** (`web/dashboard/src/components/generated/`)
   - Form components with validation
   - List components with filters/sort
   - Auto-generated from schema

## Usage

### Generate Code

```bash
# Generate all code from schema
npm run generate

# Generate specific outputs
npm run generate:types      # TypeScript + Go types
npm run generate:api        # API client SDK
npm run generate:db         # Database migrations
npm run generate:ui         # UI components
```

### Watch Mode (Auto-generate on schema changes)

```bash
npm run generate:watch
```

### Validate Schema

```bash
npm run validate:schema
```

## Live Testing Hooks

The schema includes testing hooks that trigger on save:

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

This enables:
- **Automatic test runs** when files are saved
- **Incremental testing** (only changed files)
- **Fast feedback** during development

## Collaboration Features

The schema defines collaboration features:

```json
{
  "collaboration": {
    "features": {
      "realtime": {
        "enabled": true,
        "protocol": "websocket",
        "sync": {
          "type": "operational-transform"
        }
      }
    }
  }
}
```

This enables:
- **Real-time updates** across all connected clients
- **Presence indicators** (cursors, selections)
- **File synchronization** with conflict resolution

## Schema Evolution

When updating the schema:

1. **Update schema file** (`schema/nebulabox.schema.json`)
2. **Run validation**: `npm run validate:schema`
3. **Generate code**: `npm run generate`
4. **Run tests**: `npm run test`
5. **Update documentation**: Schema changes auto-update docs

## Versioning

Schema versions follow semantic versioning:
- **Major**: Breaking changes (require migration)
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes (backward compatible)

Current version: `1.0.0`

## Next Steps

1. ✅ Schema definition created
2. ⏳ Code generation pipeline
3. ⏳ Live testing hooks integration
4. ⏳ Collaboration layer implementation
5. ⏳ Auto-generated SDKs

