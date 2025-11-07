#!/usr/bin/env node

/**
 * Live Testing Hooks
 * 
 * Triggers tests automatically on file save/change
 * Similar to Cursor's live testing feature
 */

const chokidar = require('chokidar');
const { exec } = require('child_process');
const path = require('path');

const TEST_CONFIG = {
  unit: {
    command: 'npm run test:unit',
    timeout: 30000,
    scope: 'changed-files'
  },
  integration: {
    command: 'npm run test:integration',
    timeout: 60000,
    scope: 'related-files'
  },
  e2e: {
    command: 'npm run test:e2e',
    timeout: 120000,
    scope: 'all'
  }
};

// Track file changes
const changedFiles = new Set();
let testTimeout = null;

function runTests(scope = 'changed-files') {
  // Clear existing timeout
  if (testTimeout) {
    clearTimeout(testTimeout);
  }
  
  // Debounce test runs
  testTimeout = setTimeout(() => {
    console.log(`\n🧪 Running tests for: ${Array.from(changedFiles).join(', ')}\n`);
    
    // Determine which tests to run based on changed files
    const testTypes = determineTestTypes(Array.from(changedFiles));
    
    // Run tests in parallel
    const promises = testTypes.map(testType => {
      return new Promise((resolve, reject) => {
        console.log(`Running ${testType} tests...`);
        exec(TEST_CONFIG[testType].command, (error, stdout, stderr) => {
          if (error) {
            console.error(`❌ ${testType} tests failed:`);
            console.error(stderr);
            reject(error);
          } else {
            console.log(`✅ ${testType} tests passed`);
            resolve(stdout);
          }
        });
      });
    });
    
    Promise.all(promises)
      .then(() => {
        console.log('\n✅ All tests passed!\n');
        changedFiles.clear();
      })
      .catch((error) => {
        console.error('\n❌ Some tests failed\n');
        // Don't clear changed files so tests can be re-run
      });
  }, 1000); // 1 second debounce
}

function determineTestTypes(files) {
  const testTypes = new Set();
  
  files.forEach(file => {
    // Backend files
    if (file.includes('internal/api/')) {
      testTypes.add('unit');
      testTypes.add('integration');
    }
    
    // Frontend files
    if (file.includes('web/dashboard/src/')) {
      testTypes.add('unit');
      testTypes.add('e2e');
    }
    
    // Database files
    if (file.includes('internal/database/')) {
      testTypes.add('integration');
    }
    
    // Schema changes trigger all tests
    if (file.includes('schema/')) {
      testTypes.add('unit');
      testTypes.add('integration');
      testTypes.add('e2e');
    }
  });
  
  return Array.from(testTypes);
}

// Watch for file changes
const watcher = chokidar.watch([
  'internal/**/*.go',
  'web/dashboard/src/**/*.{ts,tsx}',
  'schema/**/*.json',
  'e2e/**/*.{ts,spec.ts}'
], {
  ignored: /(^|[\/\\])\../, // ignore dotfiles
  persistent: true,
  ignoreInitial: true
});

watcher
  .on('change', (filePath) => {
    console.log(`📝 File changed: ${filePath}`);
    changedFiles.add(filePath);
    runTests();
  })
  .on('add', (filePath) => {
    console.log(`➕ File added: ${filePath}`);
    changedFiles.add(filePath);
    runTests();
  });

console.log('👀 Watching for file changes...');
console.log('Press Ctrl+C to stop\n');

// Handle graceful shutdown
process.on('SIGINT', () => {
  console.log('\n\n👋 Stopping file watcher...');
  watcher.close();
  process.exit(0);
});

