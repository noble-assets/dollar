# Vaults V2 Overview

## Introduction

The Noble Dollar Vault V2 system represents a significant evolution in decentralized yield optimization, introducing sophisticated mechanisms for managing multi-chain yield strategies while maintaining security and fair value distribution. This system features a single vault on the Noble chain that enables $USDN holders to participate in diversified yield opportunities by managing multiple remote positions across different protocols and chains through a professionally managed structure.

## Core Architecture

### Single Vault, Multiple Remote Positions

The V2 system features a single vault on the Noble chain that can maintain multiple remote positions simultaneously across different protocols and chains. This architecture provides:

- **Diversification**: Risk is spread across multiple yield sources managed by the single vault
- **Optimization**: The vault can dynamically allocate capital to the highest-performing strategies
- **Resilience**: Issues with one protocol don't compromise the entire vault
- **Flexibility**: New opportunities can be added without system migration

The Noble vault can configure:
- Maximum number of remote positions (e.g., 5-10)
- Allowed protocols (e.g., Aave, Compound, Morpho)
- Allowed chains (via IBC/Hyperlane Chain IDs)
- Rebalancing thresholds and strategies

### Net Asset Value (NAV) System

The NAV system provides accurate, real-time valuation of the Noble vault's assets including all its remote positions and inflight funds:

```
Total NAV = Local Assets 
          + Σ(Remote Position Values) 
          + Σ(Inflight Funds Values)
          - Pending Liabilities
          
NAV per Share = Total NAV / Total Outstanding Shares
```

Key features:
- **Oracle Integration**: Remote position values are updated via trusted oracles
- **Inflight Tracking**: Funds in transit remain counted in NAV at last known value
- **Staleness Protection**: Maximum age limits prevent using outdated values
- **Multi-Source Verification**: Critical updates require multiple oracle confirmations
- **Time-Weighted Averaging**: Reduces impact of temporary price spikes
- **Bridge Completion**: Tracking of completed bridge transactions

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

Remote positions leverage Hyperlane for secure cross-chain operations:

- **Deployment**: Capital bridged via Hyperlane to target chain (always in USDN)
- **Inflight Tracking**: USDN marked as inflight during bridge transit
- **Updates**: Position values relayed through Hyperlane oracles
- **Withdrawals**: Yields and principal bridged back to Noble as USDN
- **Completion Tracking**: Monitoring of bridge transaction completions
- **Emergency Recovery**: Fallback mechanisms for bridge failures

### Risk Management

Each remote position is subject to:

- **Concentration Limits**: No single position > X% of vault
- **Protocol Limits**: Maximum exposure per protocol
- **Chain Limits**: Maximum exposure per blockchain
- **Correlation Analysis**: Avoid over-concentration in correlated strategies
- **Health Monitoring**: Automatic alerts for position degradation

## Inflight Funds Management

### Overview

Inflight funds represent capital that is temporarily in transit between the Noble vault and its remote positions, or between positions during rebalancing. This capital remains fully accounted for in the NAV to ensure accurate vault valuation at all times.

### Inflight Fund Types

1. **Deposit to Position**: USDN being deployed from the Noble vault to a remote protocol
2. **Withdrawal from Position**: USDN returning from a remote protocol to the Noble vault
3. **Rebalance Between Positions**: USDN moving between the vault's remote positions (via Noble)
4. **Pending Deployment**: USDN from deposits awaiting allocation to positions
5. **Pending Withdrawal Distribution**: USDN returned from positions awaiting distribution to withdrawal queue
6. **Yield Collection**: Periodic harvest of accumulated yields in USDN

### Tracking Mechanism

Each inflight transaction maintains:
- **Transaction ID**: Unique identifier for tracking
- **Expected Value**: USDN amount sent including estimated bridge fees
- **Current Value**: Last known USDN value for NAV calculation
- **Time Bounds**: Expected arrival time and maximum duration
- **Bridge Details**: Hyperlane protocol and confirmation data
- **Status Updates**: PENDING → CONFIRMED → COMPLETED lifecycle

### NAV Impact

Inflight funds are included in NAV calculations to prevent artificial value fluctuations:

```
During Transit:
- USDN leaves source → Marked as inflight
- NAV unchanged (USDN still counted)
- Hyperlane confirms → Status: CONFIRMED
- USDN arrives → Status: COMPLETED

Special States:
- New deposits → PENDING_DEPLOYMENT until allocated
- Returned funds → PENDING_WITHDRAWAL_DISTRIBUTION until claimed
- Rebalancing → WITHDRAWAL then PENDING_DEPLOYMENT then DEPOSIT
```

### Transaction Completion

1. **Status Tracking**: Bridge transactions monitored for completion
2. **Timeout Management**: Stale transactions flagged for investigation
3. **Manual Intervention**: Failed transactions require governance action

### Risk Mitigation

- **Maximum Duration Limits**: Funds cannot remain inflight indefinitely
- **Value Caps**: Limits on total inflight exposure for the vault
- **Bridge Diversification**: Use multiple bridges to reduce single point of failure
- **Insurance Reserve**: Coverage for potential bridge failures
- **Proof Requirements**: Cryptographic verification of bridge completions

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
Mark Funds as Inflight
    ↓
Deploy Capital via Bridge
    ↓
Monitor Bridge Confirmation
    ↓
Reconcile on Arrival
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

### Remote Position Capital Return Flow

```
Initiate Position Withdrawal
    ↓
Mark Expected USDN as Inflight
    ↓
NAV Includes Inflight USDN Value
    ↓
Hyperlane Confirmation Received
    ↓
Mark Transaction as Completed
    ↓
Mark as PENDING_WITHDRAWAL_DISTRIBUTION
    ↓
Update Vault Liquidity
    ↓
Process Withdrawal Queue
```

### Rebalancing Between Positions Flow

```
Initiate Rebalance Strategy
    ↓
Withdraw USDN from Source Position
    ↓
Mark as WITHDRAWAL_FROM_POSITION Inflight
    ↓
USDN Arrives at Noble
    ↓
Mark as PENDING_DEPLOYMENT
    ↓
Deploy to Target Position
    ↓
Mark as DEPOSIT_TO_POSITION Inflight
    ↓
Confirm Arrival at Target
    ↓
Update Position Values
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

- **Single Token**: Users hold shares in the single Noble vault, complexity abstracted
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
- **Inflight Recovery**: Manual intervention for stuck bridge transactions
- **Oracle Fallback**: Manual NAV updates if oracles fail
- **Transaction Override**: Force-complete stale inflight funds
- **Emergency Withdrawal**: Direct redemption at last known NAV
- **Governance Override**: Multi-sig can intervene if needed

## Conclusion

The Noble Dollar Vault V2 system represents a sophisticated approach to decentralized yield optimization, balancing the need for high returns with robust security mechanisms. Through innovations like the withdrawal queue, the single vault's ability to manage multiple remote positions, and comprehensive anti-manipulation controls, the system provides users with professional-grade yield strategies while maintaining the security and fairness expected in DeFi.

The architecture's flexibility allows for future enhancements and new strategy additions without requiring user migration, ensuring the system can evolve with the rapidly changing DeFi landscape while protecting user value.
