package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nebulabox/nebulabox/tests"
)

func main() {
	// Command line flags
	var (
		outputDir    = flag.String("output", "./test-results", "Output directory for test results")
		reportFormat = flag.String("format", "html", "Report format (html, json)")
		verbose      = flag.Bool("verbose", false, "Verbose output")
		filter       = flag.String("filter", "", "Filter tests by category or ID")
		timeout      = flag.Duration("timeout", 30*time.Minute, "Test timeout")
		parallel     = flag.Bool("parallel", false, "Run tests in parallel")
	)
	flag.Parse()

	// Create test config
	config := tests.TestConfig{
		OutputDir:    *outputDir,
		ReportFormat: *reportFormat,
		Parallel:     *parallel,
		Verbose:      *verbose,
		Filter:       *filter,
		Timeout:      *timeout,
	}

	// Create test runner
	runner := tests.NewTestRunner(config)

	// Check for specific test arguments
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "list":
			listTestCases(runner)
		case "run":
			if len(args) > 1 {
				// Run specific tests
				testIDs := args[1:]
				runner.RunSpecificTests(testIDs)
			} else {
				// Run all tests
				runner.RunTests()
			}
		case "category":
			if len(args) > 1 {
				// Run tests by category
				category := args[1]
				runner.RunByCategory(category)
			} else {
				fmt.Println("Usage: nebulabox-test category <category>")
				os.Exit(1)
			}
		case "help":
			showHelp()
		default:
			fmt.Printf("Unknown command: %s\n", args[0])
			showHelp()
			os.Exit(1)
		}
	} else {
		// Run all tests by default
		runner.RunTests()
	}
}

// listTestCases lists all available test cases
func listTestCases(tr *tests.TestRunner) {
	fmt.Println("📋 Available Test Cases")
	fmt.Println("======================")
	
	// Load test cases
	tr.LoadTestCases()
	
	// Group by category
	categories := make(map[string][]tests.TestCase)
	for _, tc := range tr.TestSuite.TestCases {
		categories[tc.Category] = append(categories[tc.Category], tc)
	}
	
	// Print by category
	for category, testCases := range categories {
		fmt.Printf("\n📂 %s (%d tests)\n", strings.ToUpper(category), len(testCases))
		fmt.Println(strings.Repeat("-", len(category)+20))
		
		for _, tc := range testCases {
			priority := "🔴"
			switch tc.Priority {
			case "high":
				priority = "🔴"
			case "medium":
				priority = "🟡"
			case "low":
				priority = "🟢"
			}
			
			fmt.Printf("  %s %s - %s\n", priority, tc.ID, tc.Name)
			fmt.Printf("      %s\n", tc.Description)
		}
	}
	
	fmt.Printf("\nTotal: %d test cases\n", len(tr.TestSuite.TestCases))
}

// showHelp shows help information
func showHelp() {
	fmt.Println("NebulaBox Test Runner")
	fmt.Println("====================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  nebulabox-test [command] [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run [test-ids...]  Run specific tests or all tests")
	fmt.Println("  list               List all available test cases")
	fmt.Println("  category <cat>     Run tests by category")
	fmt.Println("  help               Show this help message")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -output string     Output directory for test results (default: ./test-results)")
	fmt.Println("  -format string     Report format: html, json (default: html)")
	fmt.Println("  -verbose           Verbose output")
	fmt.Println("  -filter string     Filter tests by category or ID")
	fmt.Println("  -timeout duration  Test timeout (default: 30m)")
	fmt.Println("  -parallel          Run tests in parallel")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  nebulabox-test run                    # Run all tests")
	fmt.Println("  nebulabox-test run container_001      # Run specific test")
	fmt.Println("  nebulabox-test category containers    # Run container tests")
	fmt.Println("  nebulabox-test list                   # List all tests")
	fmt.Println("  nebulabox-test -verbose run           # Run with verbose output")
}
