#!/bin/bash

# Fix large files issue in Git repository
# Removes node_modules and other large files from Git history

set -e

echo "🔧 Fixing large files issue in Git repository"
echo ""

# Check if node_modules is tracked
if git ls-files | grep -q "node_modules"; then
    echo "⚠️  Found node_modules in Git repository"
    echo ""
    
    # Ensure .gitignore has node_modules
    if ! grep -q "node_modules" .gitignore 2>/dev/null; then
        echo "📝 Adding node_modules to .gitignore..."
        echo "" >> .gitignore
        echo "# Node.js" >> .gitignore
        echo "node_modules/" >> .gitignore
        echo "npm-debug.log*" >> .gitignore
        echo "yarn-debug.log*" >> .gitignore
        echo "yarn-error.log*" >> .gitignore
    fi
    
    echo "🗑️  Removing node_modules from Git cache..."
    git rm -r --cached node_modules 2>/dev/null || true
    git rm -r --cached */node_modules 2>/dev/null || true
    git rm -r --cached */*/node_modules 2>/dev/null || true
    
    echo "✅ node_modules removed from Git tracking"
    echo ""
fi

# Check for other large files
echo "🔍 Checking for other large files (>50MB)..."
LARGE_FILES=$(git ls-files | xargs -I {} sh -c 'test -f {} && du -h {} 2>/dev/null' | awk '$1 ~ /[0-9]+[MG]/ && ($1+0 > 50 || $1 ~ /G/) {print}' || true)

if [ -n "$LARGE_FILES" ]; then
    echo "⚠️  Found large files:"
    echo "$LARGE_FILES"
    echo ""
    echo "Consider adding them to .gitignore or using Git LFS"
else
    echo "✅ No other large files found"
fi

# Check .gitignore
echo ""
echo "📋 Current .gitignore entries:"
grep -E "node_modules|dist|bin|\.exe|\.dmg|\.pkg" .gitignore 2>/dev/null || echo "  (checking...)"

echo ""
echo "✅ Large files fix complete!"
echo ""
echo "Next steps:"
echo "  1. Review changes: git status"
echo "  2. Commit the fix: git add .gitignore && git commit -m 'Remove node_modules from Git'"
echo "  3. Push again: git push origin main"
echo ""

