#!/bin/bash

# Comprehensive Noble Dollar V2 Vault Test Runner
# This script runs all vault-related tests including financial logic, adversarial scenarios, and stress tests

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
VERBOSE=${VERBOSE:-false}
COVERAGE=${COVERAGE:-false}

OUTPUT_DIR="test_results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_section() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE} $1 ${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

# Test suite definitions
declare -A TEST_SUITES=(
    ["basic"]="./keeper -run TestV2VaultSuite -timeout 10m"
    ["adversarial"]="./keeper -run TestAdversarialSuite -timeout 15m"
    ["stress"]="./keeper -run TestStressSuite -timeout 20m"
    ["existing"]="./keeper -run TestPausing,TestStakedVault,TestFlexibleVault -timeout 10m"
)

# Performance tracking
declare -A TEST_TIMES
declare -A TEST_RESULTS

# Function to run a single test suite
run_test_suite() {
    local suite_name=$1
    local test_command=$2
    local output_file="$OUTPUT_DIR/${suite_name}_${TIMESTAMP}.log"

    log_info "Running $suite_name tests..."

    local start_time=$(date +%s)

    # Build test command
    local full_command="go test $test_command"

    if [ "$VERBOSE" = true ]; then
        full_command+=" -v"
    fi

    if [ "$COVERAGE" = true ]; then
        full_command+=" -coverprofile=$OUTPUT_DIR/${suite_name}_coverage.out"
    fi

    # Run the test
    if $full_command 2>&1 | tee "$output_file"; then
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        TEST_TIMES[$suite_name]=$duration
        TEST_RESULTS[$suite_name]="PASS"
        log_success "$suite_name tests completed in ${duration}s"
    else
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        TEST_TIMES[$suite_name]=$duration
        TEST_RESULTS[$suite_name]="FAIL"
        log_error "$suite_name tests failed after ${duration}s"
        return 1
    fi
}



# Function to analyze test results
analyze_results() {
    log_section "TEST ANALYSIS"

    local total_time=0
    local passed_tests=0
    local failed_tests=0

    for suite in "${!TEST_RESULTS[@]}"; do
        local result=${TEST_RESULTS[$suite]}
        local time=${TEST_TIMES[$suite]}
        total_time=$((total_time + time))

        if [ "$result" = "PASS" ]; then
            passed_tests=$((passed_tests + 1))
            log_success "$suite: $result (${time}s)"
        else
            failed_tests=$((failed_tests + 1))
            log_error "$suite: $result (${time}s)"
        fi
    done

    echo ""
    log_info "Total test time: ${total_time}s"
    log_info "Passed suites: $passed_tests"
    log_info "Failed suites: $failed_tests"

    # Check for specific vulnerability indicators
    check_vulnerabilities

    # Generate summary report
    generate_summary_report $passed_tests $failed_tests $total_time
}

# Function to check for potential vulnerabilities in test output
check_vulnerabilities() {
    log_section "VULNERABILITY ANALYSIS"

    local vuln_count=0

    # Check for vulnerability indicators in test logs
    for log_file in "$OUTPUT_DIR"/*.log; do
        if [ -f "$log_file" ]; then
            local suite_name=$(basename "$log_file" .log | sed "s/_${TIMESTAMP}//")

            # Look for vulnerability keywords
            local vulnerabilities=$(grep -i "vulnerability\|exploit\|attack.*successful\|profit.*positive" "$log_file" || true)

            if [ -n "$vulnerabilities" ]; then
                vuln_count=$((vuln_count + 1))
                log_warning "Potential vulnerabilities detected in $suite_name:"
                echo "$vulnerabilities" | sed 's/^/  /'
                echo ""
            fi
        fi
    done

    if [ $vuln_count -eq 0 ]; then
        log_success "No obvious vulnerabilities detected in test outputs"
    else
        log_error "$vuln_count test suite(s) detected potential vulnerabilities"
    fi
}

# Function to generate summary report
generate_summary_report() {
    local passed=$1
    local failed=$2
    local total_time=$3
    local report_file="$OUTPUT_DIR/test_summary_${TIMESTAMP}.md"

    log_info "Generating summary report: $report_file"

    cat > "$report_file" << EOF
# Noble Dollar V2 Vault Test Summary

**Test Run:** $(date)
**Duration:** ${total_time} seconds
**Passed Suites:** $passed
**Failed Suites:** $failed

## Test Suite Results

| Suite | Result | Duration (s) | Notes |
|-------|--------|--------------|-------|
EOF

    for suite in "${!TEST_RESULTS[@]}"; do
        local result=${TEST_RESULTS[$suite]}
        local time=${TEST_TIMES[$suite]}
        local status_emoji="✅"
        local notes=""

        if [ "$result" = "FAIL" ]; then
            status_emoji="❌"
            notes="Check logs for details"
        fi

        echo "| $suite | $status_emoji $result | $time | $notes |" >> "$report_file"
    done

    cat >> "$report_file" << EOF

## Coverage Information

EOF

    if [ "$COVERAGE" = true ]; then
        echo "Coverage reports generated for each test suite." >> "$report_file"
        echo "Use \`go tool cover -html=<coverage_file>\` to view detailed coverage." >> "$report_file"
    else
        echo "Coverage analysis was not enabled. Use COVERAGE=true to enable." >> "$report_file"
    fi

    cat >> "$report_file" << EOF

## Key Test Areas Covered

### Financial Logic Tests
- Basic deposit/withdrawal operations
- Share price calculations
- Precision and rounding behavior
- Large number handling
- Multi-user scenarios

### Adversarial Tests
- First depositor inflation attacks
- Sandwich attacks
- Flash loan simulations
- Precision exploitation
- Slippage protection
- Value conservation

### Stress Tests
- Extreme precision scenarios
- High-frequency operations
- Massive scale operations
- Extreme volatility
- Mathematical edge cases
- Concurrent user simulation

## Recommendations

1. Review any failed tests immediately
2. Investigate vulnerability warnings
3. Consider additional tests for any edge cases discovered
4. Regularly run these tests during development
5. Update tests as the system evolves

## Files Generated

- Test logs: \`$OUTPUT_DIR/*_${TIMESTAMP}.log\`
- Coverage reports: \`$OUTPUT_DIR/*_coverage.out\` (if enabled)


EOF
}

# Function to check prerequisites
check_prerequisites() {
    log_section "CHECKING PREREQUISITES"

    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        log_error "go.mod not found. Please run this script from the project root."
        exit 1
    fi

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi

    # Check Go version
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go version: $go_version"

    # Check if test files exist
    local test_files=(
        "keeper/msg_server_vaults_v2_test.go"
        "keeper/vaults_v2_adversarial_test.go"
        "keeper/vaults_v2_stress_test.go"
    )

    for test_file in "${test_files[@]}"; do
        if [ ! -f "$test_file" ]; then
            log_warning "Test file not found: $test_file"
        else
            log_info "Found test file: $test_file"
        fi
    done

    log_success "Prerequisites check completed"
}

# Function to print usage
print_usage() {
    echo "Usage: $0 [OPTIONS] [SUITE_NAME]"
    echo ""
    echo "Options:"
    echo "  -v, --verbose     Enable verbose output"
    echo "  -c, --coverage    Enable coverage analysis"

    echo "  -h, --help        Show this help message"
    echo ""
    echo "Suite Names:"
    echo "  basic        - Basic V2 vault functionality tests"
    echo "  adversarial  - Adversarial attack scenario tests"
    echo "  stress       - Stress and edge case tests"
    echo "  existing     - Existing vault tests"
    echo "  all          - Run all test suites (default)"
    echo ""
    echo "Environment Variables:"
    echo "  VERBOSE=true     Same as -v"
    echo "  COVERAGE=true    Same as -c"

    echo ""
    echo "Examples:"
    echo "  $0                    # Run all tests"
    echo "  $0 basic              # Run only basic tests"
    echo "  $0 -v -c adversarial  # Run adversarial tests with verbose output and coverage"
    echo "  COVERAGE=true $0      # Run all tests with coverage"
}

# Main execution function
main() {
    local run_suite="all"

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -c|--coverage)
                COVERAGE=true
                shift
                ;;

            -h|--help)
                print_usage
                exit 0
                ;;
            basic|adversarial|stress|existing|all)
                run_suite=$1
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                print_usage
                exit 1
                ;;
        esac
    done

    # Print header
    log_section "NOBLE DOLLAR V2 VAULT COMPREHENSIVE TEST RUNNER"
    log_info "Timestamp: $(date)"
    log_info "Output directory: $OUTPUT_DIR"
    log_info "Verbose: $VERBOSE"
    log_info "Coverage: $COVERAGE"

    log_info "Target suite: $run_suite"

    # Check prerequisites
    check_prerequisites

    # Run tests
    local overall_success=true

    if [ "$run_suite" = "all" ]; then
        # Run all test suites
        for suite_name in "${!TEST_SUITES[@]}"; do
            if ! run_test_suite "$suite_name" "${TEST_SUITES[$suite_name]}"; then
                overall_success=false
            fi
        done
    else
        # Run specific suite
        if [ -n "${TEST_SUITES[$run_suite]}" ]; then
            if ! run_test_suite "$run_suite" "${TEST_SUITES[$run_suite]}"; then
                overall_success=false
            fi
        else
            log_error "Unknown test suite: $run_suite"
            print_usage
            exit 1
        fi
    fi



    # Analyze results
    analyze_results

    # Final status
    log_section "FINAL RESULTS"
    if [ "$overall_success" = true ]; then
        log_success "All requested tests completed successfully!"
        exit 0
    else
        log_error "Some tests failed. Check the logs for details."
        exit 1
    fi
}

# Run main function with all arguments
main "$@"
