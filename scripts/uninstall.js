#!/usr/bin/env node

/**
 * Pre-uninstall script for NebulaBox npm package
 * Cleans up binaries
 */

const fs = require('fs');
const path = require('path');

const BIN_DIR = path.join(__dirname, '..', 'bin');

console.log('🧹 Cleaning up NebulaBox binaries...');

if (fs.existsSync(BIN_DIR)) {
  try {
    const files = fs.readdirSync(BIN_DIR);
    files.forEach(file => {
      const filePath = path.join(BIN_DIR, file);
      try {
        fs.unlinkSync(filePath);
        console.log(`   Removed: ${file}`);
      } catch (err) {
        console.warn(`   Could not remove: ${file}`);
      }
    });
    
    fs.rmdirSync(BIN_DIR);
    console.log('✅ Cleanup complete');
  } catch (err) {
    console.warn(`⚠️  Cleanup warning: ${err.message}`);
  }
}

