# Noble Dollar V2 Vault Testing Suite Documentation

## Overview

This document describes the comprehensive testing suite for the Noble Dollar V2 vault system. The test suite is designed to validate financial logic, detect vulnerabilities, and ensure system robustness under extreme conditions.

## Table of Contents

- [Test Architecture](#test-architecture)
- [Test Categories](#test-categories)
- [Financial Logic Tests](#financial-logic-tests)
- [Adversarial Attack Tests](#adversarial-attack-tests)
- [Stress Tests](#stress-tests)
- [Running the Tests](#running-the-tests)
- [Security Considerations](#security-considerations)
- [Known Attack Vectors](#known-attack-vectors)
- [Test Results Interpretation](#test-results-interpretation)
- [Contributing](#contributing)

## Test Architecture

The V2 vault testing suite consists of three main test files:

### 1. `msg_server_vaults_v2_test.go`
**Purpose**: Core functionality and financial logic validation
- Basic deposit/withdrawal operations
- Share price calculations
- Multi-user scenarios
- Precision handling
- Boundary conditions

### 2. `vaults_v2_adversarial_test.go`
**Purpose**: Attack vector detection and exploitation prevention
- First depositor attacks
- Sandwich attacks
- Flash loan simulations
- Precision exploitation
- Economic manipulation

### 3. `vaults_v2_stress_test.go`
**Purpose**: System robustness under extreme conditions
- High-frequency operations
- Massive scale testing
- Mathematical edge cases
- Extreme volatility scenarios

## Test Categories

### Financial Logic Tests

#### Core Invariants Tested
1. **Value Conservation**: `totalNAV = sum(userShares * sharePrice)`
2. **Share Price Consistency**: `sharePrice = totalNAV / totalShares`
3. **Non-negative Values**: All shares, NAV, and prices must be non-negative
4. **Precision Preservation**: Rounding errors must be bounded

#### Key Test Cases
- **Basic Operations**: Deposit, withdraw, yield distribution
- **Multi-user Scenarios**: Concurrent operations, fair share distribution
- **Precision Tests**: Dust amounts, fractional calculations
- **Large Numbers**: Operations at scale limits
- **Edge Cases**: Zero states, minimum values

### Adversarial Attack Tests

#### First Depositor Attack
```go
// Attack Pattern
1. Attacker deposits minimal amount (1 wei)
2. Artificial NAV inflation occurs
3. Victim deposits at inflated share price
4. Attacker withdraws, capturing victim's value
```

**Test Coverage**:
- Minimum deposit exploits
- NAV manipulation impact
- Share price inflation effects
- Victim fund capture scenarios

#### Sandwich Attacks
```go
// Attack Pattern
1. Attacker front-runs with large deposit
2. NAV manipulation occurs (MEV/oracle)
3. Victim's transaction executes at manipulated price
4. Attacker back-runs with withdrawal
```

**Test Coverage**:
- Front-running scenarios
- Price manipulation detection
- Slippage protection validation
- MEV resistance testing

#### Flash Loan Attacks
```go
// Attack Pattern
1. Borrow large amount
2. Deposit to gain shares
3. Manipulate NAV through external means
4. Withdraw at inflated value
5. Repay loan, keep profit
```

**Test Coverage**:
- Single-block attack profitability
- Large position impact
- Price manipulation resistance
- Capital efficiency exploits

### Stress Tests

#### Precision Stress
- Extreme decimal precision scenarios
- Rounding accumulation over many operations
- Mathematical overflow/underflow protection

#### Scale Stress
- Operations with maximum possible values
- Mixed scale operations (micro + macro amounts)
- High-frequency operation sequences

#### Volatility Stress
- Extreme NAV fluctuations
- Rapid price changes
- Mathematical stability under volatility

## Running the Tests

### Quick Start
```bash
# Run all tests
./scripts/run_comprehensive_tests.sh

# Run specific test suite
./scripts/run_comprehensive_tests.sh adversarial

# Run with coverage and verbose output
./scripts/run_comprehensive_tests.sh -v -c all
```

### Manual Test Execution
```bash
# Basic tests
go test ./keeper -run TestV2VaultSuite -v

# Adversarial tests
go test ./keeper -run TestAdversarialSuite -v

# Stress tests
go test ./keeper -run TestStressSuite -v
```

### Environment Variables
```bash
export VERBOSE=true     # Detailed output
export COVERAGE=true    # Generate coverage reports
export BENCHMARKS=true  # Run performance benchmarks
```

## Security Considerations

### Critical Security Properties

#### 1. Value Conservation
**Property**: Total value in the system must equal sum of all user positions
```
Invariant: ∑(userShares[i] * sharePrice) = totalNAV
```

#### 2. Share Price Monotonicity
**Property**: Share price should only increase (or stay same) with legitimate yield
```
Constraint: sharePrice(t+1) ≥ sharePrice(t) for yield events
```

#### 3. Atomicity
**Property**: All operations must be atomic and consistent
```
Guarantee: No partial state updates that violate invariants
```

#### 4. Access Control
**Property**: Only authorized entities can perform privileged operations
```
Control: NAV updates, vault configuration changes
```

### Vulnerability Indicators

The test suite looks for these warning signs:

1. **Unexpected Profits**: Attackers gaining value without legitimate source
2. **Value Leakage**: Total system value decreasing without withdrawals
3. **Precision Loss**: Accumulated rounding errors exceeding bounds
4. **Invariant Violations**: Core mathematical properties broken
5. **Overflow Conditions**: Mathematical operations causing unexpected results

## Known Attack Vectors

### 1. First Depositor Inflation Attack
**Severity**: High
**Description**: First depositor can inflate share price to steal from subsequent depositors
**Mitigation**: Minimum deposit requirements, dead shares, or share price bounds

### 2. Sandwich Attacks via MEV
**Severity**: Medium
**Description**: Attackers can front/back-run transactions to extract value
**Mitigation**: Slippage protection, time locks, or MEV-resistant design

### 3. Precision Exploitation
**Severity**: Low-Medium
**Description**: Exploiting rounding to accumulate dust amounts
**Mitigation**: Proper rounding favor, dust collection mechanisms

### 4. Flash Loan Manipulation
**Severity**: Medium
**Description**: Using borrowed capital to manipulate share prices
**Mitigation**: Multi-block time locks, oracle price validation

### 5. Yield Dilution
**Severity**: Medium
**Description**: Front-running yield distribution to dilute other users
**Mitigation**: Fair yield distribution mechanisms, time-weighted rewards

## Test Results Interpretation

### Success Criteria
✅ **All invariants maintained**
✅ **No profitable attacks detected**
✅ **Precision errors within acceptable bounds**
✅ **System stability under stress**

### Warning Signs
⚠️ **Precision accumulation approaching limits**
⚠️ **Edge cases causing unexpected behavior**
⚠️ **Performance degradation under load**

### Failure Indicators
❌ **Invariant violations**
❌ **Profitable attack scenarios**
❌ **Mathematical overflows/underflows**
❌ **System crashes or panics**

### Example Output Analysis
```
VULNERABILITY: First depositor attack successful. Attacker profit: 500000
```
This indicates a critical vulnerability where an attacker can profit from the first depositor attack pattern.

```
GOOD: Flash loan attack not profitable
```
This indicates the system successfully prevents flash loan exploitation.

## Test Data Analysis

### Coverage Metrics
- **Line Coverage**: Percentage of code lines executed
- **Branch Coverage**: Percentage of conditional branches tested
- **Function Coverage**: Percentage of functions called

### Performance Metrics
- **Gas Usage**: Cost of operations under various conditions
- **Execution Time**: Performance under load
- **Memory Usage**: Resource consumption patterns

### Vulnerability Metrics
- **Attack Success Rate**: Percentage of attacks that succeed
- **Profit Potential**: Maximum extractable value per attack
- **Detection Rate**: Percentage of attacks detected and prevented

## Common Test Failures and Solutions

### 1. Share Price Calculation Inconsistency
**Symptom**: `Share price calculation inconsistent` error
**Cause**: Floating point precision or integer overflow
**Solution**: Review calculation order, use proper decimal types

### 2. Value Conservation Violation
**Symptom**: Total value doesn't match sum of positions
**Cause**: Rounding errors, missing state updates, or logical bugs
**Solution**: Audit all value transfer operations, check state consistency

### 3. Invariant Violations Under Stress
**Symptom**: Tests pass individually but fail under stress
**Cause**: Race conditions, accumulated errors, or edge case interactions
**Solution**: Review concurrent operation handling, improve state management

## Best Practices for Test Development

### 1. Test Independence
Each test should be independent and not rely on side effects from other tests.

### 2. Comprehensive Coverage
Test both happy paths and error conditions, including:
- Valid inputs with expected outcomes
- Invalid inputs with proper error handling
- Edge cases and boundary conditions
- Adversarial inputs designed to exploit

### 3. Realistic Scenarios
Model tests after real-world usage patterns and known attack vectors.

### 4. Performance Awareness
Consider the gas cost and computational complexity of operations being tested.

### 5. Documentation
Document the purpose, expected behavior, and potential vulnerabilities each test addresses.

## Contributing

### Adding New Tests

1. **Identify the Test Category**: Financial logic, adversarial, or stress
2. **Choose the Appropriate File**: Based on the test category
3. **Follow Naming Conventions**: Use descriptive test names
4. **Document Expected Behavior**: Add comments explaining the test purpose
5. **Include Cleanup**: Ensure tests don't leave state artifacts

### Test Structure Template
```go
func (suite *TestSuite) TestSpecificScenario() {
    // Setup
    vaultType := vaults.FLEXIBLE
    user := suite.addresses[0].String()
    
    // Test execution
    // ... test logic ...
    
    // Assertions
    suite.Require().NoError(err)
    suite.checkFinancialInvariants(vaultType)
    
    // Cleanup (if needed)
}
```

### Vulnerability Testing Guidelines

1. **Model Real Attacks**: Base tests on documented attack patterns
2. **Check for Profit**: Verify attackers cannot extract unexpected value
3. **Test State Consistency**: Ensure all invariants hold after attacks
4. **Consider Gas Costs**: Include transaction costs in profitability analysis
5. **Document Findings**: Record both successful defenses and vulnerabilities

## Continuous Testing

### Integration with CI/CD
The test suite should be integrated into the continuous integration pipeline to:
- Run on every commit
- Block merges that break tests
- Generate coverage reports
- Alert on security test failures

### Regular Security Audits
Schedule regular runs of the full test suite, especially:
- Before major releases
- After significant code changes
- Following security incidents in similar systems
- When new attack vectors are discovered

### Test Maintenance
- Update tests when system requirements change
- Add tests for newly discovered attack vectors
- Remove or update obsolete tests
- Review test effectiveness periodically

## Conclusion

This comprehensive testing suite provides extensive coverage of the Noble Dollar V2 vault system's security, functionality, and robustness. Regular execution of these tests helps ensure the system remains secure against known attack vectors while maintaining correct financial behavior under all conditions.

The test suite is designed to evolve with the system and should be updated as new requirements, features, or security considerations emerge. By following the guidelines in this document, developers can contribute to maintaining a secure and robust vault system.