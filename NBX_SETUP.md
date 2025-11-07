# Quick Setup - nbx Command

## ✅ Setup Complete!

The `nbx` command is now available. Here's how to use it:

## Option 1: Use from bin directory (No installation needed)

```bash
# Use the short command
./bin/nbx version
./bin/nbx ps
./bin/nbx --help

# Or use the full command
./bin/nebulabox version
```

## Option 2: Install to PATH (Recommended)

The installation script has already created symlinks in `~/.local/bin`:

```bash
# Add to PATH (if not already)
export PATH="$HOME/.local/bin:$PATH"

# Or add permanently to ~/.bashrc
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# Now you can use nbx from anywhere!
nbx version
nbx ps
nbx --help
```

## Option 3: Add project bin to PATH

```bash
# Add project bin to PATH
export PATH="$HOME/Documents/cursor-projects/nebulabox/bin:$PATH"

# Or add permanently
echo 'export PATH="$HOME/Documents/cursor-projects/nebulabox/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# Now use nbx from anywhere
nbx version
```

## Verify Installation

```bash
# Check if nbx is available
which nbx

# Test commands
nbx version
nbx ps
nbx images
nbx --help
```

## Both Commands Work

Both `nbx` and `nebulabox` point to the same binary:

```bash
nbx version          # Short form ✅
nebulabox version    # Full form ✅
```

## Quick Test

```bash
# Test from bin directory
./bin/nbx version

# Test if installed
nbx version
```

---

**That's it! You can now use `nbx` instead of `./bin/nebulabox`!** 🎉

