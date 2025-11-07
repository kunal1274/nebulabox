# Steps to Publish Updated npm Package

## Current Situation

✅ Install script updated with better Mac error handling  
⏳ Need to push to GitHub  
⏳ Need to republish npm package  

## Step-by-Step Process

### Step 1: Commit Changes to Git

```bash
# Check what's changed
git status

# Add npm-related files
git add package.json
git add scripts/install.js
git add scripts/uninstall.js
git add .npmignore
git add npm-publish-guide.md
git add NPM_INSTALLATION_METHODS.md
git add README_NPM.md
git add MAC_INSTALL_FIX.md
git add PUBLISH_NPM_STEPS.md

# Or add all at once
git add package.json scripts/install.js scripts/uninstall.js .npmignore *.md

# Commit
git commit -m "Add npm package support with improved Mac error handling"
```

### Step 2: Push to GitHub

```bash
# Push to main branch
git push origin main

# Verify on GitHub
# Check: https://github.com/kunal1274/nebulabox
```

### Step 3: Update npm Package Version

```bash
# Update version (patch increment)
npm version patch
# This will: 0.1.0 -> 0.1.1

# Or manually edit package.json
# "version": "0.1.1"
```

### Step 4: Test Locally (Optional but Recommended)

```bash
# Create package
npm pack

# Test install locally
npm install -g ./nebulabox-0.1.1.tgz

# Test on Mac (if you have access)
# Should show improved error messages
```

### Step 5: Publish to npm

```bash
# Make sure you're logged in
npm whoami

# If not logged in
npm login

# Publish
npm publish --access public

# Verify
npm view nebulabox
```

### Step 6: Verify Installation

On a Mac system:

```bash
# Try install (will show improved errors)
npm install -g nebulabox

# Should see:
# ⚠️  Note: Mac binaries may not be available yet...
# ❌ Installation failed: Binary not found...
# 💡 Mac Installation Alternatives...
```

## Quick Command Summary

```bash
# 1. Commit and push
git add package.json scripts/install.js scripts/uninstall.js .npmignore *.md
git commit -m "Add npm package support with improved Mac error handling"
git push origin main

# 2. Update version
npm version patch

# 3. Publish
npm publish --access public

# 4. Verify
npm view nebulabox
```

## What Gets Updated

After publishing:
- ✅ Updated install script with Mac error handling
- ✅ Better error messages for Mac users
- ✅ Helpful alternatives suggested
- ✅ Graceful 404 handling

## Testing After Publish

On Mac:
```bash
npm install -g nebulabox@latest
# Should see improved error messages
```

On Linux:
```bash
npm install -g nebulabox@latest
# Should work normally
```

## Important Notes

1. **Version Bump**: Always bump version before publishing
2. **Git First**: Push to GitHub before publishing npm (so install script can download from GitHub)
3. **Test Locally**: Test with `npm pack` first if possible
4. **Verify**: Check npm registry after publishing

---

**Ready to push and publish!** 🚀

