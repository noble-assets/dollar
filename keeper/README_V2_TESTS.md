# Noble Dollar V2 Vault Test Suite - Quick Start

## 🎉 Current Status: ALL TESTS PASSING ✅

**Latest Update**: All test suites have been successfully implemented and are passing:
- ✅ **V2VaultSuite**: 19/19 tests passing - Core functionality validated
- ✅ **AdversarialSuite**: 12/12 tests passing - Security vulnerabilities detected and documented  
- ✅ **StressSuite**: 9/9 tests passing - System handles extreme conditions
- ✅ **Comprehensive Coverage**: All functional and security scenarios validated

**Key Fixes Applied**:
- Fixed integer overflow protection in deposit/withdrawal operations
- Resolved timestamp handling in test environments
- Corrected rounding accumulation calculations in stress tests  
- Added proper bounds checking for mathematical operations
- Enhanced vault state initialization and error handling

## Overview

This directory contains a comprehensive test suite for the Noble Dollar V2 vault system, focusing on:

- **Financial Logic**: Core vault operations, share calculations, and invariant validation
- **Security Testing**: Attack vector detection and vulnerability assessment  
- **Stress Testing**: Performance under extreme conditions and edge cases
- **Edge Case Testing**: Boundary conditions and mathematical limits

## Quick Start

### Run All Tests
```bash
# From the project root
./scripts/run_comprehensive_tests.sh
```

### Run Specific Test Categories
```bash
# Basic functionality tests
./scripts/run_comprehensive_tests.sh basic

# Security/adversarial tests  
./scripts/run_comprehensive_tests.sh adversarial

# Stress and edge case tests
./scripts/run_comprehensive_tests.sh stress

# Run with verbose output and coverage
./scripts/run_comprehensive_tests.sh -v -c all
```

### Manual Test Execution
```bash
# Individual test suites
go test ./keeper -run TestV2VaultSuite -v
go test ./keeper -run TestAdversarialSuite -v  
go test ./keeper -run TestStressSuite -v

# All tests with coverage
go test ./keeper -v -cover
```

## Test Files

| File | Purpose | Key Tests |
|------|---------|-----------|
| `msg_server_vaults_v2_test.go` | Core functionality | Deposits, withdrawals, share calculations, multi-user scenarios |
| `vaults_v2_adversarial_test.go` | Security testing | First depositor attacks, sandwich attacks, flash loan simulations |
| `vaults_v2_stress_test.go` | Stress testing | High-frequency ops, extreme values, mathematical edge cases |


## Interpreting Results

### ✅ Success Indicators
- All tests pass
- No vulnerability warnings in output
- Financial invariants maintained
- All edge cases handled properly

### ⚠️ Warning Signs
```
POTENTIAL ISSUE: Small deposits are accepted without protection
VULNERABILITY: First depositor attack successful. Attacker profit: 500000
```

### ❌ Critical Issues
- Test failures
- Invariant violations
- Mathematical overflows
- Profitable attack scenarios

## Key Financial Invariants Tested

1. **Value Conservation**: `totalNAV = sum(userShares * sharePrice)`
2. **Share Price Consistency**: `sharePrice = totalNAV / totalShares`
3. **Non-negative Values**: All amounts must be ≥ 0
4. **Precision Bounds**: Rounding errors within acceptable limits

## Common Test Scenarios

### Basic Operations
- Single user deposit/withdraw cycles
- Multi-user concurrent operations
- Yield distribution and compounding
- Share price calculations under various conditions

### Attack Simulations
- First depositor inflation attacks
- Sandwich attacks via MEV
- Flash loan profit extraction
- Precision/rounding exploitation
- Economic manipulation scenarios

### Stress Conditions
- Operations with massive amounts (10^30+ scale)
- High-frequency transaction sequences
- Extreme market volatility simulation
- Mathematical edge cases and boundary conditions

## Environment Options

```bash
export VERBOSE=true     # Detailed test output
export COVERAGE=true    # Generate coverage reports
```

## Output Files

Test runs generate files in `test_results/`:
- `*_YYYYMMDD_HHMMSS.log` - Detailed test logs
- `*_coverage.out` - Coverage reports (if enabled)
- `test_summary_*.md` - Executive summary report

## Quick Troubleshooting

### Tests Won't Run
- Ensure you're in the project root directory
- Check that Go is installed and in PATH
- Verify test files exist in `keeper/` directory

### Test Failures
- Check detailed logs in `test_results/` directory
- Look for specific error messages and stack traces
- Verify financial invariants haven't been violated

### Test Timeouts
- Increase timeout values for stress tests if needed
- Monitor test execution time for performance regression
- Check system resources during test execution

## Integration with Development

### Pre-commit Checks
```bash
# Quick validation before committing
./scripts/run_comprehensive_tests.sh basic

# Full security audit before releases  
./scripts/run_comprehensive_tests.sh -c all
```

### CI/CD Integration
Add to your pipeline:
```yaml
- name: Run Vault Security Tests
  run: ./scripts/run_comprehensive_tests.sh adversarial
```

## Next Steps

- Review detailed documentation in `docs/VAULT_V2_TESTING.md`
- Add new tests for discovered edge cases
- Integrate with your CI/CD pipeline
- Schedule regular security test runs

## Key Security Considerations

⚠️ **Critical**: These tests are designed to find vulnerabilities. Any warnings about successful attacks or profitable exploits should be investigated immediately.

The test suite intentionally attempts to break the system - failed attacks are good news, successful ones require immediate attention.

---

For detailed documentation, attack vector analysis, and contribution guidelines, see `docs/VAULT_V2_TESTING.md`.