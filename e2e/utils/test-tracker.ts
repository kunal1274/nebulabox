import { Page } from '@playwright/test'
import * as fs from 'fs'
import * as path from 'path'
import { TestArtifactManager, TestRun, TestStep } from './test-artifacts'

/**
 * Test Tracker
 * Tracks test execution status (passed/failed/unexecuted)
 * Integrates with artifact manager
 */

export interface TestCase {
  suiteId: string
  testCaseId: string
  name: string
  description?: string
  preconditions?: string[]
  steps: TestCaseStep[]
}

export interface TestCaseStep {
  stepNumber: number
  name: string
  action: string
  expected: string
}

export interface TestExecutionResult {
  runId: string
  testCaseId: string
  suiteId: string
  status: 'passed' | 'failed' | 'skipped' | 'unexecuted'
  steps: TestStepResult[]
  startTime: number
  endTime?: number
  duration?: number
}

export interface TestStepResult {
  stepNumber: number
  name: string
  status: 'passed' | 'failed' | 'skipped'
  error?: string
  duration?: number
}

export class TestTracker {
  private artifactManager: TestArtifactManager
  private testRuns: Map<string, TestExecutionResult> = new Map()
  private baseDir: string

  constructor(baseDir: string = 'test-results') {
    this.artifactManager = new TestArtifactManager(baseDir)
    this.baseDir = baseDir
  }

  /**
   * Start test execution
   */
  async startTest(suiteId: string, testCaseId: string): Promise<string> {
    const run = this.artifactManager.initializeRun(suiteId, testCaseId)
    
    const execution: TestExecutionResult = {
      runId: run.runId,
      testCaseId,
      suiteId,
      status: 'unexecuted',
      steps: [],
      startTime: Date.now(),
    }

    this.testRuns.set(run.runId, execution)
    return run.runId
  }

  /**
   * Record step execution
   */
  async recordStep(
    runId: string,
    stepNumber: number,
    status: 'passed' | 'failed' | 'skipped',
    error?: string,
    duration?: number
  ): Promise<void> {
    const execution = this.testRuns.get(runId)
    if (!execution) {
      throw new Error(`Test run ${runId} not found`)
    }

    const stepResult: TestStepResult = {
      stepNumber,
      name: `Step ${stepNumber}`,
      status,
      error,
      duration,
    }

    execution.steps.push(stepResult)

    // Update overall status
    if (status === 'failed') {
      execution.status = 'failed'
    } else if (execution.status === 'unexecuted' && status === 'passed') {
      execution.status = 'passed'
    }

    // Save execution state
    await this.saveExecution(execution)
  }

  /**
   * Complete test execution
   */
  async completeTest(
    runId: string,
    finalStatus?: 'passed' | 'failed' | 'skipped'
  ): Promise<TestExecutionResult> {
    const execution = this.testRuns.get(runId)
    if (!execution) {
      throw new Error(`Test run ${runId} not found`)
    }

    execution.endTime = Date.now()
    execution.duration = execution.endTime - execution.startTime

    if (finalStatus) {
      execution.status = finalStatus
    } else {
      // Determine status from steps
      const hasFailed = execution.steps.some((s) => s.status === 'failed')
      const allPassed = execution.steps.length > 0 && execution.steps.every((s) => s.status === 'passed')
      
      if (hasFailed) {
        execution.status = 'failed'
      } else if (allPassed) {
        execution.status = 'passed'
      } else {
        execution.status = 'skipped'
      }
    }

    // Finalize artifact manager
    await this.artifactManager.finalizeRun(execution.status)

    // Save final execution
    await this.saveExecution(execution)

    // Generate summary
    await this.generateSummary(execution)

    return execution
  }

  /**
   * Save execution result
   */
  private async saveExecution(execution: TestExecutionResult): Promise<void> {
    const runDir = this.artifactManager.getRunDirectory(
      execution.suiteId,
      execution.testCaseId,
      execution.runId
    )
    const executionFile = path.join(runDir, 'execution-result.json')
    fs.writeFileSync(executionFile, JSON.stringify(execution, null, 2))
  }

  /**
   * Generate test summary
   */
  private async generateSummary(execution: TestExecutionResult): Promise<void> {
    const summaryDir = path.join(this.baseDir, execution.suiteId, 'summaries')
    fs.mkdirSync(summaryDir, { recursive: true })

    const summary = {
      suiteId: execution.suiteId,
      testCaseId: execution.testCaseId,
      runId: execution.runId,
      status: execution.status,
      startTime: execution.startTime,
      endTime: execution.endTime,
      duration: execution.duration,
      totalSteps: execution.steps.length,
      passedSteps: execution.steps.filter((s) => s.status === 'passed').length,
      failedSteps: execution.steps.filter((s) => s.status === 'failed').length,
      skippedSteps: execution.steps.filter((s) => s.status === 'skipped').length,
    }

    const summaryFile = path.join(summaryDir, `${execution.testCaseId}-${execution.runId}.json`)
    fs.writeFileSync(summaryFile, JSON.stringify(summary, null, 2))
  }

  /**
   * Get execution result
   */
  getExecution(runId: string): TestExecutionResult | undefined {
    return this.testRuns.get(runId)
  }

  /**
   * List all executions for a test case
   */
  listExecutions(suiteId: string, testCaseId: string): TestExecutionResult[] {
    return Array.from(this.testRuns.values()).filter(
      (e) => e.suiteId === suiteId && e.testCaseId === testCaseId
    )
  }

  /**
   * Generate test case report
   */
  async generateTestCaseReport(suiteId: string, testCaseId: string): Promise<string> {
    const executions = this.listExecutions(suiteId, testCaseId)
    
    const report = {
      suiteId,
      testCaseId,
      totalRuns: executions.length,
      passedRuns: executions.filter((e) => e.status === 'passed').length,
      failedRuns: executions.filter((e) => e.status === 'failed').length,
      skippedRuns: executions.filter((e) => e.status === 'skipped').length,
      unexecutedRuns: executions.filter((e) => e.status === 'unexecuted').length,
      executions: executions.map((e) => ({
        runId: e.runId,
        status: e.status,
        startTime: e.startTime,
        endTime: e.endTime,
        duration: e.duration,
        steps: e.steps.length,
      })),
    }

    const reportDir = path.join(this.baseDir, suiteId, testCaseId)
    fs.mkdirSync(reportDir, { recursive: true })
    const reportFile = path.join(reportDir, 'test-case-report.json')
    fs.writeFileSync(reportFile, JSON.stringify(report, null, 2))

    return reportFile
  }
}

