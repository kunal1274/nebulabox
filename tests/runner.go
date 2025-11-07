package tests

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// TestRunner represents the main test runner
type TestRunner struct {
	TestSuite *TestSuite
	Config    TestConfig
}

// TestConfig represents test configuration
type TestConfig struct {
	OutputDir     string
	ReportFormat  string
	Parallel      bool
	Verbose       bool
	Filter        string
	Timeout       time.Duration
}

// NewTestRunner creates a new test runner
func NewTestRunner(config TestConfig) *TestRunner {
	// Set default config
	if config.OutputDir == "" {
		config.OutputDir = "./test-results"
	}
	if config.ReportFormat == "" {
		config.ReportFormat = "json"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Minute
	}
	
	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	
	return &TestRunner{
		TestSuite: NewTestSuite(),
		Config:    config,
	}
}

// LoadTestCases loads all test cases
func (tr *TestRunner) LoadTestCases() {
	log.Println("📋 Loading test cases...")
	
	// Load all test cases
	testCases := GetAllTestCases()
	
	for _, tc := range testCases {
		tr.TestSuite.AddTestCase(tc)
	}
	
	log.Printf("✅ Loaded %d test cases", len(testCases))
}

// RunTests runs all tests
func (tr *TestRunner) RunTests() {
	log.Println("🚀 Starting test execution...")
	
	// Load test cases
	tr.LoadTestCases()
	
	// Run all tests
	tr.TestSuite.RunAllTests()
	
	// Generate report
	report := tr.TestSuite.GetTestReport()
	
	// Print report
	report.PrintReport()
	
	// Save report
	tr.saveReport(report)
	
	// Exit with appropriate code
	if report.FailedTests > 0 {
		os.Exit(1)
	}
}

// RunSpecificTests runs specific test cases
func (tr *TestRunner) RunSpecificTests(testIDs []string) {
	log.Printf("🎯 Running specific tests: %v", testIDs)
	
	// Load test cases
	tr.LoadTestCases()
	
	var results []TestRun
	
	// Run specific tests
	for _, tc := range tr.TestSuite.TestCases {
		for _, testID := range testIDs {
			if tc.ID == testID {
				tr := tr.TestSuite.RunTestCase(tc)
				results = append(results, tr)
				break
			}
		}
	}
	
	// Generate report
	report := TestReport{
		TotalTests:    len(results),
		PassedTests:   0,
		FailedTests:   0,
		SkippedTests:  0,
		TotalDuration: 0,
		TestRuns:      results,
	}
	
	for _, tr := range results {
		report.TotalDuration += tr.Duration
		switch tr.Status {
		case "passed":
			report.PassedTests++
		case "failed":
			report.FailedTests++
		case "skipped":
			report.SkippedTests++
		}
	}
	
	// Print report
	report.PrintReport()
	
	// Save report
	tr.saveReport(report)
}

// RunByCategory runs tests by category
func (tr *TestRunner) RunByCategory(category string) {
	log.Printf("📂 Running tests by category: %s", category)
	
	// Load test cases
	tr.LoadTestCases()
	
	var results []TestRun
	
	// Run tests by category
	for _, tc := range tr.TestSuite.TestCases {
		if tc.Category == category {
			tr := tr.TestSuite.RunTestCase(tc)
			results = append(results, tr)
		}
	}
	
	// Generate report
	report := TestReport{
		TotalTests:    len(results),
		PassedTests:   0,
		FailedTests:   0,
		SkippedTests:  0,
		TotalDuration: 0,
		TestRuns:      results,
	}
	
	for _, tr := range results {
		report.TotalDuration += tr.Duration
		switch tr.Status {
		case "passed":
			report.PassedTests++
		case "failed":
			report.FailedTests++
		case "skipped":
			report.SkippedTests++
		}
	}
	
	// Print report
	report.PrintReport()
	
	// Save report
	tr.saveReport(report)
}

// saveReport saves the test report
func (tr *TestRunner) saveReport(report TestReport) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	
	// Save JSON report
	jsonFile := filepath.Join(tr.Config.OutputDir, fmt.Sprintf("test-report-%s.json", timestamp))
	if err := report.SaveReport(jsonFile); err != nil {
		log.Printf("❌ Failed to save JSON report: %v", err)
	} else {
		log.Printf("📄 JSON report saved: %s", jsonFile)
	}
	
	// Save HTML report
	htmlFile := filepath.Join(tr.Config.OutputDir, fmt.Sprintf("test-report-%s.html", timestamp))
	if err := tr.generateHTMLReport(report, htmlFile); err != nil {
		log.Printf("❌ Failed to save HTML report: %v", err)
	} else {
		log.Printf("📄 HTML report saved: %s", htmlFile)
	}
}

// generateHTMLReport generates an HTML test report
func (tr *TestRunner) generateHTMLReport(report TestReport, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>NebulaBox Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f0f0f0; padding: 20px; border-radius: 5px; }
        .summary { margin: 20px 0; }
        .test-case { margin: 10px 0; padding: 10px; border: 1px solid #ddd; border-radius: 3px; }
        .passed { background: #d4edda; border-color: #c3e6cb; }
        .failed { background: #f8d7da; border-color: #f5c6cb; }
        .skipped { background: #fff3cd; border-color: #ffeaa7; }
        .details { margin-top: 10px; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="header">
        <h1>NebulaBox Test Report</h1>
        <p>Generated: %s</p>
    </div>
    
    <div class="summary">
        <h2>Summary</h2>
        <p>Total Tests: %d</p>
        <p>Passed: %d</p>
        <p>Failed: %d</p>
        <p>Skipped: %d</p>
        <p>Success Rate: %.2f%%</p>
        <p>Total Duration: %v</p>
    </div>
    
    <div class="test-results">
        <h2>Test Results</h2>
`, time.Now().Format("2006-01-02 15:04:05"), 
		report.TotalTests, report.PassedTests, report.FailedTests, report.SkippedTests,
		float64(report.PassedTests)/float64(report.TotalTests)*100, report.TotalDuration)
	
	for _, tr := range report.TestRuns {
		statusClass := tr.Status
		html += fmt.Sprintf(`
        <div class="test-case %s">
            <h3>%s</h3>
            <p><strong>Status:</strong> %s</p>
            <p><strong>Duration:</strong> %v</p>
            <p><strong>Description:</strong> %s</p>
`, statusClass, tr.TestCase.Name, tr.Status, tr.Duration, tr.TestCase.Description)
		
		if tr.Error != "" {
			html += fmt.Sprintf(`            <p><strong>Error:</strong> %s</p>`, tr.Error)
		}
		
		html += `        </div>`
	}
	
	html += `
    </div>
</body>
</html>`
	
	_, err = file.WriteString(html)
	return err
}
