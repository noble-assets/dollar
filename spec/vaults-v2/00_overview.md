# Vaults V2 Overview

## Introduction

The Noble Dollar Vaults V2 system represents a significant evolution in decentralized yield optimization, introducing sophisticated mechanisms for managing multi-chain yield strategies while maintaining security and fair value distribution. This system enables $USDN holders to participate in diversified yield opportunities across multiple protocols and chains through a professionally managed vault structure.

## Core Architecture

### Multi-Position Vaults

Unlike traditional single-strategy vaults, V2 vaults can maintain multiple remote positions simultaneously across different protocols and chains. This architecture provides:

- **Diversification**: Risk is spread across multiple yield sources
- **Optimization**: Capital can be dynamically allocated to the highest-performing strategies
- **Resilience**: Issues with one protocol don't compromise the entire vault
- **Flexibility**: New opportunities can be added without vault migration

Each vault can configure:
- Maximum number of remote positions (e.g., 5-10)
- Allowed protocols (e.g., Aave, Compound, Morpho)
- Allowed chains (via IBC/Hyperlane Chain IDs)
- Rebalancing thresholds and strategies

### Net Asset Value (NAV) System

The NAV system provides accurate, real-time valuation of vault assets including all remote positions:

```
Total NAV = Local Assets + Σ(Remote Position Values) - Pending Liabilities
NAV per Share = Total NAV / Total Outstanding Shares
```

Key features:
- **Oracle Integration**: Remote position values are updated via trusted oracles
- **Staleness Protection**: Maximum age limits prevent using outdated values
- **Multi-Source Verification**: Critical updates require multiple oracle confirmations
- **Time-Weighted Averaging**: Reduces impact of temporary price spikes

### Withdrawal Queue Mechanism

The withdrawal queue is a critical innovation that prevents value extraction attacks while ensuring fair redemptions:

#### Queue Processing Flow

1. **Request Phase**: Users submit withdrawal requests, locking their shares
2. **Queue Entry**: Requests enter FIFO queue with NAV snapshot
3. **Liquidity Monitoring**: System tracks available liquidity from:
   - Local vault reserves
   - Maturing positions
   - Incoming deposits
4. **Fulfillment**: When liquidity available, requests are marked CLAIMABLE
5. **Claim Phase**: Users claim their $USDN after withdrawal delay period

#### Fair Value Protection

- Each request uses the NAV at time of request, not fulfillment
- Prevents front-running of NAV updates
- Eliminates sandwich attacks around large deposits/withdrawals
- Ensures all users receive fair value regardless of queue position

## Security Mechanisms

### Anti-Manipulation Framework

The system implements multiple layers of defense against malicious behavior:

#### 1. Deposit Velocity Controls

```
Velocity Score = (Recent Volume / Time Window) × (Deposit Count / Expected Count)
```

- Tracks deposit patterns over rolling time windows
- Flags suspicious rapid deposits that could indicate attacks
- Enforces cooldown periods for high-velocity users

#### 2. Deposit Limits

Multiple limit types work in concert:

- **Per-User Limits**: Maximum total deposit per address
- **Per-Block Limits**: Maximum deposits in single block
- **Per-Transaction Limits**: Maximum single deposit amount
- **Total Vault Limits**: Maximum total deposits across all users
- **Minimum Amounts**: Prevents dust attacks and griefing

#### 3. Share Price Manipulation Defense

- **Entry/Exit Delays**: Withdrawal queue prevents immediate round-trips
- **NAV Snapshots**: Lock in values at request time
- **Slippage Protection**: Maximum acceptable NAV deviation checks
- **Fee Structure**: Management and performance fees discourage short-term speculation

#### 4. Oracle Security

Remote position values are protected through:

- **Staleness Checks**: Reject updates older than configured threshold
- **Deviation Limits**: Flag suspicious large value changes
- **Proof Verification**: Cryptographic proofs for cross-chain values
- **Emergency Pauses**: Ability to halt on oracle failures

### Withdrawal Security

The withdrawal queue provides multiple security benefits:

- **No Instant Liquidity**: Prevents bank run scenarios
- **Orderly Unwinding**: Allows time to close remote positions
- **Fair Ordering**: FIFO processing prevents preferential treatment
- **Partial Fulfillment**: Can process withdrawals as liquidity arrives
- **Emergency Mode**: Can pause new deposits while honoring withdrawals

## Remote Position Management

### Position Lifecycle

1. **Creation**: Deploy capital to approved protocol/chain
2. **Monitoring**: Continuous NAV updates via oracles
3. **Rebalancing**: Adjust allocations based on performance
4. **Harvesting**: Collect yields back to vault
5. **Closure**: Withdraw capital when needed

### Cross-Chain Coordination

Remote positions leverage Hyperlane/IBC for secure cross-chain operations:

- **Deployment**: Capital bridged via Hyperlane/IBC to target chain
- **Updates**: Position values relayed through Hyperlane/IBC oracles
- **Withdrawals**: Yields and principal bridged back to Noble
- **Emergency Recovery**: Fallback mechanisms for bridge failures

### Risk Management

Each remote position is subject to:

- **Concentration Limits**: No single position > X% of vault
- **Protocol Limits**: Maximum exposure per protocol
- **Chain Limits**: Maximum exposure per blockchain
- **Correlation Analysis**: Avoid over-concentration in correlated strategies
- **Health Monitoring**: Automatic alerts for position degradation

## Operational Flows

### Deposit Flow with Security Checks

```
User Deposit Request
    ↓
Velocity Check → [Fail: Reject]
    ↓ Pass
Limit Checks → [Fail: Reject]
    ↓ Pass
Cooldown Check → [Fail: Reject]
    ↓ Pass
Calculate Shares (Current NAV)
    ↓
Mint Shares to User
    ↓
Update Velocity Metrics
    ↓
Deploy Capital to Positions
```

### Withdrawal Flow with Queue

```
Withdrawal Request
    ↓
Lock User Shares
    ↓
Enter Queue (FIFO)
    ↓
Record NAV Snapshot
    ↓
Wait for Liquidity
    ↓
Process Queue (Keeper/Auto)
    ↓
Mark as CLAIMABLE
    ↓
User Claims (After Delay)
    ↓
Burn Shares & Transfer $USDN
```

## Economic Model

### Fee Structure

- **Management Fee**: Annual percentage on total AUM
- **Performance Fee**: Percentage of profits above hurdle rate
- **No Deposit/Withdrawal Fees**: Encourages participation
- **Fee Distribution**: To protocol treasury and/or stakers

### Yield Distribution

- **Automatic Compounding**: Yields increase NAV per share
- **No Manual Claims**: Value accrues to share price
- **Fair Distribution**: Proportional to share ownership
- **Tax Efficiency**: Unrealized gains until withdrawal

## Advantages Over V1

### Enhanced Capital Efficiency

- **Multi-Strategy**: Higher yields through diversification
- **Dynamic Allocation**: Respond to market opportunities
- **Cross-Chain Reach**: Access best yields anywhere
- **Professional Management**: Expert strategy selection

### Improved Security

- **Queue System**: Eliminates many attack vectors
- **Velocity Controls**: Prevents rapid manipulation
- **NAV Snapshots**: Fair value for all users
- **Oracle Redundancy**: Multiple verification sources

### Better User Experience

- **Single Token**: Users hold vault shares, complexity abstracted
- **Automatic Optimization**: No manual strategy switching
- **Transparent NAV**: Clear value per share
- **Predictable Withdrawals**: Queue provides certainty

## Emergency Procedures

### Pause Mechanisms

- **Deposit Pause**: Stop new deposits, allow withdrawals
- **Withdrawal Pause**: Emergency only, requires governance
- **Position Pause**: Halt specific remote positions
- **Global Pause**: Complete system halt (extreme cases)

### Recovery Procedures

- **Position Recovery**: Retrieve funds from failed positions
- **Oracle Fallback**: Manual NAV updates if oracles fail
- **Emergency Withdrawal**: Direct redemption at last known NAV
- **Governance Override**: Multi-sig can intervene if needed

## Conclusion

The Noble Dollar Vaults V2 system represents a sophisticated approach to decentralized yield optimization, balancing the need for high returns with robust security mechanisms. Through innovations like the withdrawal queue, multi-position architecture, and comprehensive anti-manipulation controls, the system provides users with professional-grade yield strategies while maintaining the security and fairness expected in DeFi.

The architecture's flexibility allows for future enhancements and new strategy additions without requiring user migration, ensuring the system can evolve with the rapidly changing DeFi landscape while protecting user value.
