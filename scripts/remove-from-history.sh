#!/bin/bash

# Remove large files from entire Git history
# WARNING: This rewrites Git history. Use with caution!

set -e

echo "⚠️  WARNING: This will rewrite Git history!"
echo "   Make sure you have a backup or are working on a branch."
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

echo ""
echo "🔍 Finding large files in Git history..."

# Files to remove from history
FILES_TO_REMOVE=(
    "node_modules/node/bin/node"
    "node_modules/"
    "go.tar.gz"
)

echo ""
echo "🗑️  Removing files from Git history..."

# Use git filter-branch to remove files from all commits
for FILE in "${FILES_TO_REMOVE[@]}"; do
    echo "   Removing: $FILE"
    
    # Remove file from all commits
    git filter-branch --force --index-filter \
        "git rm -rf --cached --ignore-unmatch '$FILE'" \
        --prune-empty --tag-name-filter cat -- --all 2>/dev/null || true
done

# Clean up backup refs
echo ""
echo "🧹 Cleaning up backup refs..."
rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now --aggressive

echo ""
echo "✅ History cleaned!"
echo ""
echo "📋 Next steps:"
echo "  1. Verify: git log --all --oneline | head -5"
echo "  2. Force push: git push origin main --force"
echo ""
echo "⚠️  WARNING: Force push rewrites remote history!"
echo "   Only do this if you're the only one working on the repo."

