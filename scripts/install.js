#!/usr/bin/env node

/**
 * Post-install script for NebulaBox npm package
 * Downloads the appropriate binary for the current platform
 */

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

const PACKAGE_VERSION = require('../package.json').version;
const BIN_DIR = path.join(__dirname, '..', 'bin');
const PLATFORMS = {
  'linux': {
    'x64': 'linux-amd64',
    'arm64': 'linux-arm64'
  },
  'darwin': {
    'x64': 'darwin-amd64',
    'arm64': 'darwin-arm64'
  },
  'win32': {
    'x64': 'windows-amd64.exe',
    'arm64': 'windows-arm64.exe'
  }
};

function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;
  
  if (!PLATFORMS[platform]) {
    console.error(`❌ Unsupported platform: ${platform}`);
    process.exit(1);
  }
  
  const archMap = PLATFORMS[platform];
  if (!archMap[arch]) {
    console.error(`❌ Unsupported architecture: ${arch} on ${platform}`);
    process.exit(1);
  }
  
  return archMap[arch];
}

function downloadBinary(platform, callback) {
  const fileName = `nbx-${platform}`;
  const url = `https://github.com/kunal1274/nebulabox/releases/download/v${PACKAGE_VERSION}/${fileName}`;
  const isWindows = process.platform === 'win32';
  const filePath = path.join(BIN_DIR, isWindows ? 'nebulabox.exe' : 'nebulabox');
  
  console.log(`📦 Downloading NebulaBox for ${platform}...`);
  console.log(`   URL: ${url}`);
  
  // Create bin directory if it doesn't exist
  if (!fs.existsSync(BIN_DIR)) {
    fs.mkdirSync(BIN_DIR, { recursive: true });
  }
  
  const file = fs.createWriteStream(filePath);
  
  const handleResponse = (response) => {
    if (response.statusCode === 302 || response.statusCode === 301) {
      // Follow redirect
      https.get(response.headers.location, handleResponse);
      return;
    }
    
    if (response.statusCode === 404) {
      file.close();
      if (fs.existsSync(filePath)) {
        fs.unlinkSync(filePath);
      }
      callback(new Error(`Binary not found for platform: ${platform}\nAvailable platforms: linux-amd64, linux-arm64\nNote: Mac binaries not available yet. Use 'go build' or Linux VM.`));
      return;
    }
    
    if (response.statusCode !== 200) {
      file.close();
      if (fs.existsSync(filePath)) {
        fs.unlinkSync(filePath);
      }
      callback(new Error(`Failed to download: ${response.statusCode} ${response.statusMessage}\nURL: ${url}`));
      return;
    }
    
    response.pipe(file);
    file.on('finish', () => {
      file.close();
      if (!isWindows) {
        try {
          fs.chmodSync(filePath, '755');
        } catch (err) {
          console.warn(`⚠️  Could not set permissions: ${err.message}`);
        }
      }
      console.log(`✅ Downloaded: ${filePath}`);
      callback(null, filePath);
    });
  };
  
  https.get(url, handleResponse).on('error', (err) => {
    file.close();
    if (fs.existsSync(filePath)) {
      fs.unlinkSync(filePath);
    }
    callback(new Error(`Network error: ${err.message}\nURL: ${url}`));
  });
}

function createSymlinks(binaryPath) {
  const nbxPath = path.join(BIN_DIR, 'nbx');
  const nebulaboxPath = path.join(BIN_DIR, 'nebulabox');
  
  try {
    // Create symlinks (or copies on Windows)
    if (process.platform === 'win32') {
      // Windows doesn't support symlinks well, so copy
      if (fs.existsSync(binaryPath)) {
        fs.copyFileSync(binaryPath, nbxPath.replace('.exe', '') + '.exe');
        fs.copyFileSync(binaryPath, nebulaboxPath);
      }
    } else {
      // Unix-like systems: create symlinks
      if (fs.existsSync(binaryPath)) {
        if (fs.existsSync(nbxPath)) fs.unlinkSync(nbxPath);
        if (fs.existsSync(nebulaboxPath)) fs.unlinkSync(nebulaboxPath);
        
        fs.symlinkSync(path.basename(binaryPath), nbxPath);
        fs.symlinkSync(path.basename(binaryPath), nebulaboxPath);
        
        console.log(`✅ Created symlinks: nbx, nebulabox`);
      }
    }
  } catch (err) {
    console.warn(`⚠️  Could not create symlinks: ${err.message}`);
  }
}

// Main installation
const platform = getPlatform();
const currentPlatform = process.platform;
const currentArch = process.arch;

console.log(`\n🚀 Installing NebulaBox v${PACKAGE_VERSION}`);
console.log(`   Detected: ${currentPlatform} ${currentArch}`);
console.log(`   Platform: ${platform}\n`);

// Check if Mac binaries are available
if (currentPlatform === 'darwin') {
  console.log(`⚠️  Note: Mac binaries may not be available yet.`);
  console.log(`   The package will attempt to download, but if it fails:`);
  console.log(`   - Mac binaries require Linux-specific code to be ported`);
  console.log(`   - Use 'go build' from source for CLI testing`);
  console.log(`   - Or use Linux VM/container for full features\n`);
}

downloadBinary(platform, (err, binaryPath) => {
  if (err) {
    console.error(`\n❌ Installation failed: ${err.message}`);
    
    if (currentPlatform === 'darwin') {
      console.error(`\n💡 Mac Installation Alternatives:`);
      console.error(`\n   1. Build from source (CLI only):`);
      console.error(`      git clone https://github.com/kunal1274/nebulabox.git`);
      console.error(`      cd nebulabox`);
      console.error(`      go build -o nbx ./cmd/nebulabox`);
      console.error(`      sudo mv nbx /usr/local/bin/`);
      console.error(`\n   2. Use Linux VM/Container:`);
      console.error(`      docker run -it ubuntu:22.04`);
      console.error(`      # Then install inside Linux container`);
      console.error(`\n   3. Wait for Mac binaries (when available)`);
    } else {
      console.error(`\n💡 Alternative: Build from source`);
      console.error(`   git clone https://github.com/kunal1274/nebulabox.git`);
      console.error(`   cd nebulabox && go build -o nbx ./cmd/nebulabox\n`);
    }
    
    process.exit(1);
  }
  
  createSymlinks(binaryPath);
  
  console.log(`\n✅ NebulaBox installed successfully!`);
  console.log(`\n📋 Usage:`);
  console.log(`   npx nbx version`);
  console.log(`   npx nbx --help`);
  console.log(`\n   Or add to PATH:`);
  console.log(`   export PATH=$PATH:${BIN_DIR}`);
  console.log(`\n`);
});

