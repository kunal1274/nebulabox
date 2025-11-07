# Clean Git History - Remove Large Files

## Problem

Even after removing files from Git tracking, they remain in Git history. GitHub checks the entire history and rejects pushes if any commit contains files >100MB.

## Solution Options

### Option 1: Simple Clean (Recommended)

```bash
# This creates a backup and cleans history
./scripts/clean-history-simple.sh

# Then force push
git push origin main --force
```

**What it does:**
- Creates a backup branch
- Removes `node_modules/` from all commits
- Removes `go.tar.gz` from all commits
- Cleans up Git references
- You can restore from backup if needed

### Option 2: Manual Clean

```bash
# Interactive script with more control
./scripts/remove-from-history.sh

# Then force push
git push origin main --force
```

### Option 3: Manual Commands

If you prefer to run commands manually:

```bash
# 1. Create backup
git branch backup-before-clean

# 2. Remove node_modules from all commits
git filter-branch --force --index-filter \
    'git rm -rf --cached --ignore-unmatch node_modules' \
    --prune-empty --tag-name-filter cat -- --all

# 3. Remove go.tar.gz from all commits
git filter-branch --force --index-filter \
    'git rm -rf --cached --ignore-unmatch go.tar.gz' \
    --prune-empty --tag-name-filter cat -- --all

# 4. Clean up
rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now --aggressive

# 5. Force push
git push origin main --force
```

## ⚠️ Important Warnings

1. **Force Push Rewrites History**: This changes all commit hashes
2. **Team Coordination**: If others are working on this repo, coordinate with them
3. **Backup First**: The script creates a backup branch, but you can also:
   ```bash
   git clone --mirror https://github.com/kunal1274/nebulabox.git backup-repo.git
   ```

## After Cleaning

1. ✅ Verify: `git log --oneline | head -5`
2. ✅ Check size: `du -sh .git`
3. ✅ Test push: `git push origin main --force`
4. ✅ Verify on GitHub: Check that the push succeeds

## If Something Goes Wrong

```bash
# Restore from backup branch
git reset --hard backup-before-clean

# Or restore from remote
git fetch origin
git reset --hard origin/main
```

## Alternative: Start Fresh (Nuclear Option)

If history cleaning is too complex:

```bash
# 1. Create orphan branch (no history)
git checkout --orphan new-main

# 2. Add all current files
git add .
git commit -m "Initial commit - cleaned history"

# 3. Delete old main
git branch -D main

# 4. Rename new branch
git branch -m main

# 5. Force push
git push origin main --force
```

**Warning**: This loses all commit history!

---

**Recommended**: Use `./scripts/clean-history-simple.sh` - it's the safest option with automatic backup.

