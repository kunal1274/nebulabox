# Fix Large Files Issue - GitHub Push Error

## Problem

GitHub rejected the push because of large files:
- `node_modules/node/bin/node` is 118.02 MB
- GitHub's file size limit is 100 MB

## Solution

### Quick Fix (Recommended)

```bash
# Run the fix script
./scripts/fix-large-files.sh

# Review changes
git status

# Commit the fix
git add .gitignore
git commit -m "Remove node_modules from Git tracking"

# Push again
git push origin main
```

### Manual Fix

1. **Remove node_modules from Git:**
   ```bash
   git rm -r --cached node_modules
   git rm -r --cached */node_modules
   git rm -r --cached */*/node_modules
   ```

2. **Ensure .gitignore includes node_modules:**
   ```bash
   echo "node_modules/" >> .gitignore
   ```

3. **Commit and push:**
   ```bash
   git add .gitignore
   git commit -m "Remove node_modules from Git"
   git push origin main
   ```

### If Files Are Already in Git History

If large files are in previous commits, you need to remove them from history:

```bash
# Install git-filter-repo (if not installed)
# pip install git-filter-repo

# Remove node_modules from entire history
git filter-repo --path node_modules --invert-paths

# Force push (WARNING: This rewrites history)
git push origin main --force
```

**⚠️ Warning:** Force push rewrites history. Only do this if:
- You're the only one working on the repository, OR
- You've coordinated with your team

## Prevention

### Always Add to .gitignore

Make sure `.gitignore` includes:
```
node_modules/
dist/
bin/
*.log
```

### Check Before Committing

```bash
# Check for large files
find . -type f -size +10M -not -path "./.git/*" | head -10

# Check what will be committed
git status
git diff --cached --stat
```

## Verify Fix

```bash
# Check if node_modules is still tracked
git ls-files | grep node_modules

# Should return nothing if fixed correctly
```

## Next Steps After Fix

1. ✅ Run `./scripts/fix-large-files.sh`
2. ✅ Commit changes: `git add .gitignore && git commit -m "Remove node_modules"`
3. ✅ Push: `git push origin main`
4. ✅ Verify: Check GitHub repository

---

**After fixing, you can proceed with creating releases!** 🚀

