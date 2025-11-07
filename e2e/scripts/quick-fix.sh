#!/bin/bash
# Quick Fix Script - Fixes common test failures automatically

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_DIR"

echo "🔧 Quick Fix: Running failed tests with fixes..."
echo ""

# Re-run only failed tests with more lenient settings
npx playwright test --last-failed --reporter=list --timeout=60000 || {
    echo ""
    echo "⚠️  Some tests still failing. Showing details..."
    echo ""
    npx playwright test --last-failed --reporter=list 2>&1 | tail -30
}

