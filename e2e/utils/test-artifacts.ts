import { Page } from '@playwright/test'
import * as fs from 'fs'
import * as path from 'path'

/**
 * Test Artifact Management System
 * Manages hierarchical test artifacts: Suite > Test Case > Test Run > Artifacts
 */

export interface TestArtifact {
  type: 'screenshot' | 'video' | 'log' | 'network' | 'console' | 'report'
  name: string
  path: string
  metadata?: Record<string, any>
}

export interface TestRun {
  runId: string
  testCaseId: string
  suiteId: string
  status: 'passed' | 'failed' | 'skipped' | 'unexecuted'
  startTime: number
  endTime?: number
  artifacts: TestArtifact[]
  steps: TestStep[]
}

export interface TestStep {
  stepNumber: number
  name: string
  action: string
  status: 'passed' | 'failed' | 'skipped'
  beforeScreenshot?: string
  afterScreenshot?: string
  timestamp: number
  duration?: number
  error?: string
  metadata?: Record<string, any>
}

export class TestArtifactManager {
  private baseDir: string
  private currentRun: TestRun | null = null

  constructor(baseDir: string = 'test-results') {
    this.baseDir = baseDir
  }

  /**
   * Initialize test run structure
   */
  initializeRun(suiteId: string, testCaseId: string): TestRun {
    const runId = `run-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    
    this.currentRun = {
      runId,
      testCaseId,
      suiteId,
      status: 'unexecuted',
      startTime: Date.now(),
      artifacts: [],
      steps: [],
    }

    // Create directory structure
    const runDir = this.getRunDirectory(suiteId, testCaseId, runId)
    fs.mkdirSync(runDir, { recursive: true })
    fs.mkdirSync(path.join(runDir, 'screenshots'), { recursive: true })
    fs.mkdirSync(path.join(runDir, 'videos'), { recursive: true })
    fs.mkdirSync(path.join(runDir, 'logs'), { recursive: true })
    fs.mkdirSync(path.join(runDir, 'network'), { recursive: true })
    fs.mkdirSync(path.join(runDir, 'reports'), { recursive: true })

    return this.currentRun
  }

  /**
   * Get run directory path
   */
  getRunDirectory(suiteId: string, testCaseId: string, runId: string): string {
    return path.join(this.baseDir, suiteId, testCaseId, runId)
  }

  /**
   * Add test step with before/after screenshots
   */
  async addStep(
    page: Page,
    stepNumber: number,
    name: string,
    action: string
  ): Promise<TestStep> {
    if (!this.currentRun) {
      throw new Error('Test run not initialized')
    }

    const stepStartTime = Date.now()

    // Capture before screenshot
    const beforeScreenshot = await this.captureScreenshot(
      page,
      `step-${stepNumber.toString().padStart(2, '0')}_before.png`
    )

    // Perform action (caller should do this)
    // After action, capture after screenshot

    const step: TestStep = {
      stepNumber,
      name,
      action,
      status: 'passed',
      beforeScreenshot,
      timestamp: stepStartTime,
    }

    this.currentRun.steps.push(step)
    return step
  }

  /**
   * Complete step with after screenshot
   */
  async completeStep(
    page: Page,
    stepNumber: number,
    status: 'passed' | 'failed',
    error?: string,
    metadata?: Record<string, any>
  ): Promise<void> {
    if (!this.currentRun) {
      throw new Error('Test run not initialized')
    }

    const step = this.currentRun.steps.find((s) => s.stepNumber === stepNumber)
    if (!step) {
      throw new Error(`Step ${stepNumber} not found`)
    }

    // Capture after screenshot
    const afterScreenshot = await this.captureScreenshot(
      page,
      `step-${stepNumber.toString().padStart(2, '0')}_after.png`
    )

    step.afterScreenshot = afterScreenshot
    step.status = status
    step.error = error
    step.metadata = metadata
    step.duration = Date.now() - step.timestamp

    // Save step data
    await this.saveStepData(step)
  }

  /**
   * Capture screenshot with overlay annotation
   */
  async captureScreenshot(
    page: Page,
    filename: string,
    annotation?: { x: number; y: number; label: string }
  ): Promise<string> {
    if (!this.currentRun) {
      throw new Error('Test run not initialized')
    }

    const runDir = this.getRunDirectory(
      this.currentRun.suiteId,
      this.currentRun.testCaseId,
      this.currentRun.runId
    )
    const screenshotPath = path.join(runDir, 'screenshots', filename)

    // Add overlay if annotation provided
    if (annotation) {
      await this.addOverlay(page, annotation.x, annotation.y, annotation.label)
    }

    await page.screenshot({ path: screenshotPath, fullPage: true })

    // Remove overlay
    if (annotation) {
      await this.removeOverlay(page)
    }

    // Add to artifacts
    this.currentRun.artifacts.push({
      type: 'screenshot',
      name: filename,
      path: screenshotPath,
    })

    return screenshotPath
  }

  /**
   * Add visual overlay for annotation
   */
  private async addOverlay(
    page: Page,
    x: number,
    y: number,
    label: string
  ): Promise<void> {
    await page.evaluate(
      ({ x, y, label }) => {
        const id = 'test-overlay'
        let el = document.getElementById(id)
        if (!el) {
          el = document.createElement('div')
          el.id = id
          Object.assign(el.style, {
            position: 'fixed',
            left: '0',
            top: '0',
            width: '100%',
            height: '100%',
            pointerEvents: 'none',
            zIndex: '9999999',
          })
          document.body.appendChild(el)
        }

        const svgNS = 'http://www.w3.org/2000/svg'
        const svg = document.createElementNS(svgNS, 'svg')
        svg.setAttribute(
          'style',
          'position:absolute; left:0; top:0; width:100%; height:100%'
        )

        // Circle
        const circle = document.createElementNS(svgNS, 'circle')
        circle.setAttribute('cx', x.toString())
        circle.setAttribute('cy', y.toString())
        circle.setAttribute('r', '18')
        circle.setAttribute('fill', '#ffb300')
        circle.setAttribute('opacity', '0.95')
        svg.appendChild(circle)

        // Text
        const text = document.createElementNS(svgNS, 'text')
        text.setAttribute('x', x.toString())
        text.setAttribute('y', (y + 6).toString())
        text.setAttribute('text-anchor', 'middle')
        text.setAttribute('font-size', '14')
        text.setAttribute('font-weight', 'bold')
        text.setAttribute('fill', '#000')
        text.textContent = label
        svg.appendChild(text)

        el.appendChild(svg)
      },
      { x, y, label }
    )
  }

  /**
   * Remove overlay
   */
  private async removeOverlay(page: Page): Promise<void> {
    await page.evaluate(() => {
      const el = document.getElementById('test-overlay')
      if (el) {
        el.remove()
      }
    })
  }

  /**
   * Save step data
   */
  private async saveStepData(step: TestStep): Promise<void> {
    if (!this.currentRun) return

    const runDir = this.getRunDirectory(
      this.currentRun.suiteId,
      this.currentRun.testCaseId,
      this.currentRun.runId
    )
    const stepFile = path.join(
      runDir,
      'logs',
      `step-${step.stepNumber.toString().padStart(2, '0')}.json`
    )

    fs.writeFileSync(stepFile, JSON.stringify(step, null, 2))
  }

  /**
   * Finalize test run
   */
  async finalizeRun(status: 'passed' | 'failed' | 'skipped'): Promise<TestRun> {
    if (!this.currentRun) {
      throw new Error('Test run not initialized')
    }

    this.currentRun.status = status
    this.currentRun.endTime = Date.now()

    // Save run metadata
    const runDir = this.getRunDirectory(
      this.currentRun.suiteId,
      this.currentRun.testCaseId,
      this.currentRun.runId
    )
    const metadataFile = path.join(runDir, 'run-metadata.json')
    fs.writeFileSync(metadataFile, JSON.stringify(this.currentRun, null, 2))

    // Generate HTML report
    await this.generateReport(this.currentRun)

    const finalized = this.currentRun
    this.currentRun = null

    return finalized
  }

  /**
   * Generate HTML report
   */
  private async generateReport(run: TestRun): Promise<void> {
    const runDir = this.getRunDirectory(run.suiteId, run.testCaseId, run.runId)
    const reportPath = path.join(runDir, 'reports', 'report.html')

    const html = `
<!DOCTYPE html>
<html>
<head>
  <title>Test Report - ${run.testCaseId}</title>
  <style>
    body { font-family: sans-serif; margin: 20px; }
    .header { background: #f5f5f5; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
    .step { border: 1px solid #ddd; padding: 15px; margin: 10px 0; border-radius: 4px; }
    .step-header { font-weight: bold; margin-bottom: 10px; }
    .step-status { display: inline-block; padding: 4px 8px; border-radius: 4px; color: white; }
    .passed { background: #28a745; }
    .failed { background: #dc3545; }
    .skipped { background: #ffc107; color: #000; }
    .screenshots { display: flex; gap: 20px; margin-top: 10px; }
    .screenshot img { max-width: 400px; border: 1px solid #ddd; }
  </style>
</head>
<body>
  <div class="header">
    <h1>Test Report: ${run.testCaseId}</h1>
    <p><strong>Suite:</strong> ${run.suiteId}</p>
    <p><strong>Run ID:</strong> ${run.runId}</p>
    <p><strong>Status:</strong> <span class="step-status ${run.status}">${run.status.toUpperCase()}</span></p>
    <p><strong>Duration:</strong> ${((run.endTime || Date.now()) - run.startTime) / 1000}s</p>
    <p><strong>Steps:</strong> ${run.steps.length}</p>
  </div>

  ${run.steps
    .map(
      (step) => `
    <div class="step">
      <div class="step-header">
        Step ${step.stepNumber}: ${step.name}
        <span class="step-status ${step.status}">${step.status}</span>
      </div>
      <p><strong>Action:</strong> ${step.action}</p>
      <p><strong>Duration:</strong> ${step.duration ? step.duration + 'ms' : 'N/A'}</p>
      ${step.error ? `<p style="color: red;"><strong>Error:</strong> ${step.error}</p>` : ''}
      <div class="screenshots">
        ${step.beforeScreenshot ? `<div><p>Before</p><img src="${this.relativePath(step.beforeScreenshot, reportPath)}" /></div>` : ''}
        ${step.afterScreenshot ? `<div><p>After</p><img src="${this.relativePath(step.afterScreenshot, reportPath)}" /></div>` : ''}
      </div>
    </div>
  `
    )
    .join('')}
</body>
</html>
    `

    fs.writeFileSync(reportPath, html)

    run.artifacts.push({
      type: 'report',
      name: 'report.html',
      path: reportPath,
    })
  }

  /**
   * Get relative path for HTML images
   */
  private relativePath(targetPath: string, fromPath: string): string {
    const target = path.resolve(targetPath)
    const from = path.resolve(path.dirname(fromPath))
    return path.relative(from, target)
  }

  /**
   * Get current run
   */
  getCurrentRun(): TestRun | null {
    return this.currentRun
  }
}

