#!/bin/bash

# Simple script to remove node_modules from Git history
# Uses BFG Repo-Cleaner approach or git filter-branch

set -e

echo "🧹 Cleaning Git history of large files..."
echo ""

# Check if we're on main/master
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"
echo ""

# Create a backup branch first
BACKUP_BRANCH="backup-before-clean-$(date +%Y%m%d-%H%M%S)"
echo "📦 Creating backup branch: $BACKUP_BRANCH"
git branch "$BACKUP_BRANCH" || true

echo ""
echo "🗑️  Removing node_modules from all commits..."

# Remove node_modules from entire history
git filter-branch --force --index-filter \
    'git rm -rf --cached --ignore-unmatch node_modules' \
    --prune-empty --tag-name-filter cat -- --all

echo ""
echo "🗑️  Removing go.tar.gz from all commits..."
git filter-branch --force --index-filter \
    'git rm -rf --cached --ignore-unmatch go.tar.gz' \
    --prune-empty --tag-name-filter cat -- --all

echo ""
echo "🧹 Cleaning up..."
rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now --aggressive

echo ""
echo "✅ History cleaned!"
echo ""
echo "📋 Summary:"
echo "  - Backup branch created: $BACKUP_BRANCH"
echo "  - node_modules removed from history"
echo "  - go.tar.gz removed from history"
echo ""
echo "⚠️  Next: Force push to update remote"
echo "  git push origin $CURRENT_BRANCH --force"
echo ""
echo "If something goes wrong, restore from backup:"
echo "  git reset --hard $BACKUP_BRANCH"

