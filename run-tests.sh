#!/bin/bash

set -e  # Exit on any error

# Colors for pretty output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create logs directory
mkdir -p logs

echo -e "${BLUE}🧹 Cleaning up previous state...${NC}"
docker compose down --volumes --remove-orphans >/dev/null 2>&1 || true
go clean -testcache

echo -e "${BLUE}🚀 Starting test database...${NC}"
docker compose up -d library-test-db

echo -e "${BLUE}⏳ Waiting for database to be ready...${NC}"
sleep 5

echo -e "${BLUE}🧪 Running tests...${NC}\n"

# Run tests and save output
echo -e "${YELLOW}Running Tests...${NC}"
go test -v ./internal/... 2>&1 | tee logs/test_output.log
test_status=${PIPESTATUS[0]}

# Parse test results from the log file
total_tests=$(grep -c "^=== RUN" logs/test_output.log || echo 0)
passed_tests=$(grep -c "^--- PASS" logs/test_output.log || echo 0)
failed_tests=$(grep -c "^--- FAIL" logs/test_output.log || echo 0)
skipped_tests=$(grep -c "^--- SKIP" logs/test_output.log || echo 0)
no_tests=$(grep -c "\[no test files\]" logs/test_output.log || echo 0)

# Get package results
packages_passed=$(grep "^ok" logs/test_output.log | wc -l || echo 0)
packages_failed=$(grep "^FAIL" logs/test_output.log | grep -v "FAIL    " | wc -l || echo 0)

# Calculate test coverage
echo -e "\n${YELLOW}Calculating test coverage...${NC}"
go test -v -coverprofile=logs/coverage.out ./internal/... 2>&1 | tee -a logs/test_output.log
coverage=$(go tool cover -func=logs/coverage.out | grep total | awk '{print $3}')

echo -e "\n${BLUE}📊 Test Statistics:${NC}"
echo "----------------------------------------"
echo -e "${BLUE}Test Cases:${NC}"
total_tests=${total_tests:-0}
passed_tests=${passed_tests:-0}
failed_tests=${failed_tests:-0}
skipped_tests=${skipped_tests:-0}

echo -e "  Total Run:            $total_tests"
echo -e "  ${GREEN}Passed:${NC}              $passed_tests"
if test "$failed_tests" -gt 0 2>/dev/null; then
    echo -e "  ${RED}Failed:${NC}              $failed_tests"
else
    echo -e "  Failed:                0"
fi
if test "$skipped_tests" -gt 0 2>/dev/null; then
    echo -e "  ${YELLOW}Skipped:${NC}             $skipped_tests"
else
    echo -e "  Skipped:               0"
fi

echo -e "\n${BLUE}Packages:${NC}"
echo -e "  Total:                $(($packages_passed + $packages_failed + $no_tests))"
echo -e "  ${GREEN}Passed:${NC}              $packages_passed"
if [ $packages_failed -gt 0 ]; then
    echo -e "  ${RED}Failed:${NC}              $packages_failed"
else
    echo -e "  Failed:                $packages_failed"
fi
echo -e "  Without Tests:         $no_tests"
echo -e "  Coverage:             $coverage"
echo "----------------------------------------"

# Display execution time
execution_time=$(grep "ok" logs/test_output.log | awk '{sum += $3} END {print sum}')
echo -e "Total Execution Time: ${execution_time}s"

echo -e "\n${BLUE}🧹 Cleaning up...${NC}"
docker compose down --volumes --remove-orphans >/dev/null 2>&1

# Save detailed logs
timestamp=$(date +%Y%m%d_%H%M%S)
mv logs/test_output.log logs/test_output_${timestamp}.log
echo -e "${BLUE}📝 Detailed logs saved to:${NC} logs/test_output_${timestamp}.log"

# Exit with proper status and message
if [ $test_status -eq 0 ]; then
    echo -e "\n${GREEN}✨ All tests passed successfully!${NC}"
    exit 0
else
    echo -e "\n${RED}💥 Some tests failed!${NC}"
    echo -e "Check logs for details: logs/test_output_${timestamp}.log"
    exit 1
fi
