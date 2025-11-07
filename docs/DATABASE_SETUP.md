# Database Setup Guide

This guide explains how to set up PostgreSQL and MongoDB for NebulaBox.

## PostgreSQL Setup

### Installation

#### macOS (using Homebrew)
```bash
brew install postgresql@15
brew services start postgresql@15
```

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### Docker
```bash
docker run --name nebulabox-postgres \
  -e POSTGRES_USER=nebulabox \
  -e POSTGRES_PASSWORD=nebulabox \
  -e POSTGRES_DB=nebulabox \
  -p 5432:5432 \
  -d postgres:15
```

### Create Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database and user
CREATE DATABASE nebulabox;
CREATE USER nebulabox WITH PASSWORD 'nebulabox';
GRANT ALL PRIVILEGES ON DATABASE nebulabox TO nebulabox;
\q
```

### Environment Variables

```bash
export NEBULABOX_POSTGRES_HOST=localhost
export NEBULABOX_POSTGRES_PORT=5432
export NEBULABOX_POSTGRES_USER=nebulabox
export NEBULABOX_POSTGRES_PASSWORD=nebulabox
export NEBULABOX_POSTGRES_DB=nebulabox
export NEBULABOX_POSTGRES_SSLMODE=disable
```

## MongoDB Setup

### Installation

#### macOS (using Homebrew)
```bash
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

#### Ubuntu/Debian
```bash
wget -qO - https://www.mongodb.org/static/pgp/server-7.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
sudo apt update
sudo apt install -y mongodb-org
sudo systemctl start mongod
sudo systemctl enable mongod
```

#### Docker
```bash
docker run --name nebulabox-mongodb \
  -p 27017:27017 \
  -d mongo:7
```

### Environment Variables

```bash
export NEBULABOX_MONGODB_URI=mongodb://localhost:27017
# OR specify host/port separately
export NEBULABOX_MONGODB_HOST=localhost
export NEBULABOX_MONGODB_PORT=27017
export NEBULABOX_MONGODB_DB=nebulabox
```

## Database Schema

### PostgreSQL Tables

- `containers` - Container metadata and configuration
- `images` - Image metadata and tags
- `workspaces` - Shared runtime workspaces
- `workspace_members` - Workspace membership and roles
- `invites` - Workspace invitations
- `sessions` - Active workspace sessions
- `snapshots` - Container/workspace/volume snapshots
- `deployments` - Orchestrator deployments
- `nodes` - Cluster nodes
- `container_groups` - Container groups
- `templates` - Stack templates
- `users` - User accounts
- `teams` - Team management
- `tenants` - Tenant isolation
- `networks` - Custom networks
- `services` - Service discovery entries

### MongoDB Collections

- `audit_logs` - Audit events (90 days TTL)
- `container_logs` - Container stdout/stderr (30 days TTL)
- `api_metrics` - API request/response metrics (7 days TTL)
- `system_metrics` - System resource usage (30 days TTL)
- `build_logs` - Image build logs (14 days TTL)
- `test_runs` - E2E test execution history (90 days TTL)

## Auto-Migration

The database schema is automatically migrated when the server starts. The `AutoMigrate()` function in `internal/database/postgres.go` will create all tables and indexes if they don't exist.

## Manual Migration

If you need to run migrations manually:

```bash
# The AutoMigrate function is called automatically on server start
# Or you can call it programmatically:
go run cmd/api/main.go --migrate
```

## Connection Pooling

### PostgreSQL
- Max idle connections: 10
- Max open connections: 100
- Connection max lifetime: 1 hour

### MongoDB
- Max pool size: 100
- Min pool size: 10
- Max connection idle time: 30 seconds

## Health Checks

Both databases support health checks:

```go
// PostgreSQL
postgres := database.GetPostgreSQL()
err := postgres.HealthCheck()

// MongoDB
mongo := database.GetMongoDB()
err := mongo.HealthCheck()
```

## Troubleshooting

### PostgreSQL Connection Issues

1. Check if PostgreSQL is running:
   ```bash
   pg_isready
   ```

2. Check connection string:
   ```bash
   psql -h localhost -U nebulabox -d nebulabox
   ```

3. Check firewall:
   ```bash
   sudo ufw status
   ```

### MongoDB Connection Issues

1. Check if MongoDB is running:
   ```bash
   mongosh --eval "db.adminCommand('ping')"
   ```

2. Check connection:
   ```bash
   mongosh mongodb://localhost:27017/nebulabox
   ```

3. Check logs:
   ```bash
   # macOS
   tail -f /usr/local/var/log/mongodb/mongo.log
   
   # Linux
   sudo tail -f /var/log/mongodb/mongod.log
   ```

## Backup and Restore

### PostgreSQL Backup

```bash
# Backup
pg_dump -U nebulabox nebulabox > backup.sql

# Restore
psql -U nebulabox nebulabox < backup.sql
```

### MongoDB Backup

```bash
# Backup
mongodump --db nebulabox --out /backup/nebulabox

# Restore
mongorestore --db nebulabox /backup/nebulabox/nebulabox
```

