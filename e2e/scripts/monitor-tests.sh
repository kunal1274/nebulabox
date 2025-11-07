#!/bin/bash
# Playwright Test Monitor
# Continuously monitors test results and reports failures

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
RESULTS_DIR="$PROJECT_DIR/test-results"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Playwright Test Monitor${NC}\n"

# Function to check for failures
check_failures() {
    if [ ! -f "$RESULTS_DIR/results.json" ]; then
        echo -e "${YELLOW}⚠️  No test results found. Run tests first.${NC}"
        return 1
    fi

    # Use Node.js to parse JSON (more reliable)
    node -e "
        const fs = require('fs');
        const results = JSON.parse(fs.readFileSync('$RESULTS_DIR/results.json', 'utf-8'));
        const failures = [];
        
        function extractFailures(suites) {
            suites.forEach(suite => {
                suite.specs?.forEach(spec => {
                    spec.tests?.forEach(test => {
                        const lastResult = test.results?.[test.results.length - 1];
                        if (lastResult?.status === 'failed') {
                            failures.push({
                                file: spec.file || 'unknown',
                                test: test.title || 'unknown',
                                error: lastResult.error?.message || 'Unknown error'
                            });
                        }
                    });
                });
                if (suite.suites) extractFailures(suite.suites);
            });
        }
        
        extractFailures(results.suites || []);
        
        const passed = results.stats?.ok || 0;
        const failed = failures.length;
        const skipped = results.stats?.skipped || 0;
        const total = passed + failed + skipped;
        
        console.log(JSON.stringify({ total, passed, failed, skipped, failures }));
    " > /tmp/test-results.json
    
    TOTAL=$(jq -r '.total' /tmp/test-results.json)
    PASSED=$(jq -r '.passed' /tmp/test-results.json)
    FAILED=$(jq -r '.failed' /tmp/test-results.json)
    SKIPPED=$(jq -r '.skipped' /tmp/test-results.json)
    
    echo -e "${BLUE}📊 Test Results:${NC}"
    echo -e "   Total: $TOTAL"
    echo -e "   ${GREEN}✅ Passed: $PASSED${NC}"
    echo -e "   ${RED}❌ Failed: $FAILED${NC}"
    echo -e "   ${YELLOW}⏭️  Skipped: $SKIPPED${NC}"
    
    if [ "$FAILED" -gt 0 ]; then
        echo -e "\n${RED}❌ Failures Found:${NC}\n"
        jq -r '.failures[] | "\(.test)\n   File: \(.file)\n   Error: \(.error)\n"' /tmp/test-results.json
        return 1
    else
        echo -e "\n${GREEN}✅ All tests passed!${NC}"
        return 0
    fi
}

# Function to run tests and monitor
run_and_monitor() {
    echo -e "${BLUE}🚀 Running tests...${NC}\n"
    
    cd "$PROJECT_DIR"
    npx playwright test --reporter=json,list
    
    echo -e "\n${BLUE}📊 Analyzing results...${NC}\n"
    check_failures
}

# Function to watch mode (re-run on file changes)
watch_mode() {
    echo -e "${BLUE}👀 Watch mode: Monitoring for changes...${NC}\n"
    echo -e "${YELLOW}Press Ctrl+C to stop${NC}\n"
    
    cd "$PROJECT_DIR"
    
    while true; do
        echo -e "${BLUE}[$(date +%H:%M:%S)] Running tests...${NC}"
        npx playwright test --reporter=json,list 2>&1 | tail -20
        
        check_failures
        
        echo -e "\n${YELLOW}Waiting for changes... (Press Ctrl+C to stop)${NC}\n"
        sleep 5
    done
}

# Main
case "${1:-check}" in
    check|status)
        check_failures
        ;;
    run|test)
        run_and_monitor
        ;;
    watch)
        watch_mode
        ;;
    *)
        echo "Usage: $0 [check|run|watch]"
        echo ""
        echo "Commands:"
        echo "  check, status  - Check existing test results"
        echo "  run, test      - Run tests and check results"
        echo "  watch          - Watch mode (auto-re-run on changes)"
        exit 1
        ;;
esac

