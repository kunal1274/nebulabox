#!/bin/bash

# Setup Git Hooks for NebulaBox
# Installs pre-commit hook for automated testing

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

HOOKS_DIR=".git/hooks"
PRE_COMMIT_HOOK="$HOOKS_DIR/pre-commit"

echo "🔧 Setting up Git hooks..."

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Create pre-commit hook
cat > "$PRE_COMMIT_HOOK" << 'EOF'
#!/bin/bash
# Pre-commit hook for NebulaBox
# Automatically runs tests before commits

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../scripts" && pwd)"
exec "$SCRIPT_DIR/pre-commit.sh"
EOF

chmod +x "$PRE_COMMIT_HOOK"

echo "✅ Git hooks installed!"
echo ""
echo "Pre-commit hook will now run tests before each commit."
echo "To skip hooks: git commit --no-verify"

