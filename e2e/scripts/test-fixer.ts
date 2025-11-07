#!/usr/bin/env node
/**
 * Playwright Test Fixer CLI
 * Utilities to check, fix, and monitor Playwright tests from console
 */

import { execSync } from 'child_process'
import * as fs from 'fs'
import * as path from 'path'

interface TestFailure {
  file: string
  test: string
  error: string
  line?: number
}

interface TestResult {
  passed: number
  failed: number
  skipped: number
  failures: TestFailure[]
}

/**
 * Parse Playwright test results
 */
function parseTestResults(): TestResult {
  const resultsFile = path.join(__dirname, '../test-results/results.json')
  
  if (!fs.existsSync(resultsFile)) {
    console.log('⚠️  No test results found. Run tests first.')
    return { passed: 0, failed: 0, skipped: 0, failures: [] }
  }

  const results = JSON.parse(fs.readFileSync(resultsFile, 'utf-8'))
  const failures: TestFailure[] = []

  for (const suite of results.suites || []) {
    for (const spec of suite.specs || []) {
      for (const test of spec.tests || []) {
        if (test.results && test.results.length > 0) {
          const lastResult = test.results[test.results.length - 1]
          if (lastResult.status === 'failed') {
            failures.push({
              file: spec.file || 'unknown',
              test: test.title || 'unknown',
              error: lastResult.error?.message || 'Unknown error',
              line: lastResult.error?.location?.line,
            })
          }
        }
      }
    }
  }

  const passed = results.stats?.ok || 0
  const failed = failures.length
  const skipped = results.stats?.skipped || 0

  return { passed, failed, skipped, failures }
}

/**
 * Check for test failures
 */
function checkFailures() {
  console.log('🔍 Checking for test failures...\n')
  
  const results = parseTestResults()
  
  console.log(`📊 Test Results:`)
  console.log(`   ✅ Passed: ${results.passed}`)
  console.log(`   ❌ Failed: ${results.failed}`)
  console.log(`   ⏭️  Skipped: ${results.skipped}`)
  
  if (results.failures.length > 0) {
    console.log(`\n❌ Failures Found:\n`)
    results.failures.forEach((failure, index) => {
      console.log(`${index + 1}. ${failure.test}`)
      console.log(`   File: ${failure.file}`)
      console.log(`   Error: ${failure.error}`)
      if (failure.line) {
        console.log(`   Line: ${failure.line}`)
      }
      console.log('')
    })
    process.exit(1)
  } else {
    console.log('\n✅ All tests passed!')
    process.exit(0)
  }
}

/**
 * Fix specific test failures
 */
function fixTest(testName?: string, file?: string) {
  console.log('🔧 Fixing test failures...\n')
  
  const results = parseTestResults()
  
  if (results.failures.length === 0) {
    console.log('✅ No failures to fix!')
    return
  }

  const failuresToFix = testName
    ? results.failures.filter(f => f.test.includes(testName))
    : file
    ? results.failures.filter(f => f.file.includes(file))
    : results.failures

  if (failuresToFix.length === 0) {
    console.log('⚠️  No matching failures found')
    return
  }

  console.log(`Found ${failuresToFix.length} failure(s) to analyze:\n`)
  
  failuresToFix.forEach((failure, index) => {
    console.log(`${index + 1}. ${failure.test}`)
    console.log(`   File: ${path.basename(failure.file)}`)
    console.log(`   Error: ${failure.error}`)
    console.log(`\n💡 Suggestions:`)
    
    // Provide fix suggestions based on error type
    if (failure.error.includes('timeout')) {
      console.log('   - Increase timeout or wait for element to be visible')
      console.log('   - Check if page is loading correctly')
    } else if (failure.error.includes('not visible')) {
      console.log('   - Element might not be rendered yet - add wait')
      console.log('   - Check if selector is correct')
    } else if (failure.error.includes('404') || failure.error.includes('Not Found')) {
      console.log('   - Check if route exists in App.tsx')
      console.log('   - Verify page component is imported')
    } else if (failure.error.includes('network')) {
      console.log('   - Check if API server is running')
      console.log('   - Verify network requests are correct')
    }
    console.log('')
  })
}

/**
 * Run tests and auto-fix common issues
 */
function runAndFix() {
  console.log('🚀 Running tests and checking for failures...\n')
  
  try {
    // Run tests
    execSync('npx playwright test --reporter=json --output=test-results', {
      stdio: 'inherit',
      cwd: path.join(__dirname, '..'),
    })
    
    // Check results
    checkFailures()
  } catch (error) {
    console.error('❌ Test run failed:', error)
    checkFailures()
  }
}

/**
 * Show test status summary
 */
function showStatus() {
  console.log('📊 Test Status Summary\n')
  
  const results = parseTestResults()
  const total = results.passed + results.failed + results.skipped
  const passRate = total > 0 ? ((results.passed / total) * 100).toFixed(1) : '0'
  
  console.log(`Total Tests: ${total}`)
  console.log(`✅ Passed: ${results.passed} (${passRate}%)`)
  console.log(`❌ Failed: ${results.failed}`)
  console.log(`⏭️  Skipped: ${results.skipped}`)
  
  if (results.failures.length > 0) {
    console.log(`\n🔴 Failing Tests:`)
    results.failures.forEach((f, i) => {
      console.log(`   ${i + 1}. ${path.basename(f.file)} - ${f.test}`)
    })
  }
}

/**
 * Open Playwright UI with specific test
 */
function openUI(testFile?: string) {
  console.log('🎨 Opening Playwright UI...\n')
  
  const cmd = testFile
    ? `npx playwright test --ui ${testFile}`
    : 'npx playwright test --ui'
  
  execSync(cmd, {
    stdio: 'inherit',
    cwd: path.join(__dirname, '..'),
  })
}

// CLI Interface
const command = process.argv[2]
const arg1 = process.argv[3]
const arg2 = process.argv[4]

switch (command) {
  case 'check':
  case 'status':
    showStatus()
    break
    
  case 'fix':
    fixTest(arg1, arg2)
    break
    
  case 'run':
  case 'test':
    runAndFix()
    break
    
  case 'ui':
    openUI(arg1)
    break
    
  default:
    console.log(`
🎭 Playwright Test Fixer CLI

Usage:
  npm run test:fix <command> [options]

Commands:
  check, status           Show test status summary
  fix [test] [file]      Analyze and suggest fixes for failures
  run, test              Run tests and check for failures
  ui [file]              Open Playwright UI (optionally for specific file)

Examples:
  npm run test:fix check              # Show status
  npm run test:fix fix                # Analyze all failures
  npm run test:fix fix "should load"  # Analyze specific test
  npm run test:fix ui                 # Open UI
  npm run test:fix ui tests/all-pages.spec.ts  # Open UI for specific file

`)
    process.exit(0)
}

