# Quick Start Guide - NebulaBox CLI

## Installation

### Option 1: Quick Install (Recommended)

```bash
# Build the CLI
make build-cli-test

# Install to ~/.local/bin (adds to PATH)
./scripts/install-nbx.sh
```

Then add to your PATH (if not already):
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Option 2: Use Directly (No Installation)

```bash
# Build the CLI
make build-cli-test

# Use with full path
./bin/nbx version
./bin/nbx ps
./bin/nbx --help
```

### Option 3: Add to PATH Manually

```bash
# Build the CLI
make build-cli-test

# Add to PATH for current session
export PATH="$PWD/bin:$PATH"

# Or add to ~/.bashrc permanently
echo 'export PATH="$HOME/Documents/cursor-projects/nebulabox/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

## Verify Installation

```bash
# Check if nbx is available
which nbx

# Test version command
nbx version

# Test help
nbx --help
```

## Usage

Now you can use the shorter `nbx` command:

```bash
# Basic commands
nbx version
nbx --help
nbx ps
nbx images

# Container operations
nbx run nginx:latest --name web
nbx stop web
nbx logs web

# Image operations
nbx pull nginx:latest
nbx images
nbx rmi nginx:latest

# Build
nbx build . -f buildspec.json -t myapp:latest

# Hierarchy
nbx hierarchy create parent --file child-spec.json
nbx hierarchy tree
nbx hierarchy list

# Groups
nbx group create --file group-spec.json
nbx group list
```

## Troubleshooting

### Command Not Found

If you get "command not found":

1. **Check if binary exists:**
   ```bash
   ls -lh bin/nbx
   ```

2. **Check if in PATH:**
   ```bash
   echo $PATH | grep -i nebulabox
   ```

3. **Use full path:**
   ```bash
   ./bin/nbx version
   ```

4. **Reinstall:**
   ```bash
   ./scripts/install-nbx.sh
   ```

### Permission Denied

```bash
chmod +x bin/nbx
chmod +x bin/nebulabox
```

## Both Commands Available

Both `nbx` and `nebulabox` work the same:

```bash
nbx version          # Short form
nebulabox version    # Full form
```

Both point to the same binary!

