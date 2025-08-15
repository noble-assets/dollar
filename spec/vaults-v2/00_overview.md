# Vaults V2 Overview

## Introduction

The Noble Dollar Vault V2 system represents a significant evolution in decentralized yield optimization, introducing sophisticated mechanisms for managing multi-chain yield strategies while maintaining security and fair value distribution. This system features a single vault on the Noble chain that enables $USDN holders to participate in diversified yield opportunities by managing multiple remote positions across different protocols and chains through a professionally managed structure.

## Core Architecture

### Single Vault, Multiple Remote Positions

The V2 system features a single vault on the Noble chain that can maintain multiple remote positions simultaneously across different chains. These remote positions are ERC-4626 compatible vaults. This architecture provides:

- **Diversification**: Risk is spread across multiple yield sources managed by the single vault
- **Optimization**: The vault can dynamically allocate capital to the highest-performing strategies
- **Resilience**: Issues with one protocol don't compromise the entire vault
- **Flexibility**: New opportunities can be added without system migration

The Noble vault can configure:
- Maximum number of remote positions (e.g., 5-10)
- Allowed protocols (e.g., Hyperliquid, Base lending protocols, Noble App Layer vaults)
- Allowed chains (via Hyperlane Domain IDs: 998 for Hyperliquid, 8453 for Base, 4000261 for Noble App Layer)
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
- **Push-Based Oracle Integration**: Remote chains actively push position values to Noble via Hyperlane using fixed-length byte encoding
- **Automatic NAV Updates**: Noble receives and processes byte-encoded price updates as they arrive from remote chains
- **Inflight Tracking**: USDN marked as inflight per Hyperlane route ID during bridge transit
- **Staleness Protection**: Maximum age limits prevent using outdated values when pushes stop arriving
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
- **Message Verification**: Hyperlane message authentication for cross-chain values
- **Fallback Values**: Use last known good values on oracle issues

### Withdrawal Security

The withdrawal queue provides multiple security benefits:

- **No Instant Liquidity**: Prevents bank run scenarios
- **Orderly Unwinding**: Allows time to close remote positions
- **Fair Ordering**: FIFO processing prevents preferential treatment
- **Partial Fulfillment**: Can process withdrawals as liquidity arrives
- **Emergency Mode**: Can pause new deposits while honoring withdrawals

## Remote Position Management

### Position Lifecycle

1. **Creation**: Deploy USDN to approved ERC-4626 compatible vault on target chain
2. **Share Tracking**: Receive and track vault shares representing the position
3. **Price Reception**: Remote chains push share price updates to Noble via Hyperlane
4. **Automatic NAV Updates**: Noble processes pushed price data to maintain current valuations
5. **Rebalancing**: Redeem shares from one vault and deposit to another
6. **Harvesting**: Yields compound within remote vaults
7. **Closure**: Redeem remote vault shares for USDN and withdraw to Noble


### Cross-Chain Coordination

Remote positions leverage Hyperlane for secure cross-chain operations:

- **Deployment**: USDN bridged via specific Hyperlane routes to deposit into ERC-4626 compatible vaults
- **Share Management**: Remote shares received and tracked for each remote position
- **Inflight Tracking**: USDN marked as inflight per Hyperlane route ID during bridge transit
- **Route Management**: Each Hyperlane route (e.g., Noble→Hyperliquid, Base→Noble) tracked separately
- **Push-Based Price Updates**: Remote chains proactively push share prices to Noble via Hyperlane using fixed-length byte encoding for efficiency
- **Automatic Value Updates**: Noble continuously receives and applies byte-encoded price data directly from the Hyperlane Mailbox
- **Redemptions**: Remote vault shares redeemed for USDN and bridged back to Noble via specific return routes
- **Completion Tracking**: Monitoring of bridge transaction completions per route
- **Emergency Recovery**: Fallback mechanisms for bridge failures on specific routes

### Risk Management

Each remote position is subject to:

- **Concentration Limits**: No single position > X% of total vault
- **Chain Limits**: Maximum exposure per blockchain
- **Approved Vaults Only**: Only deploy to pre-approved vault addresses
- **Push-Based Price Monitoring**: Receive continuous share price updates pushed from remote chains
- **Health Monitoring**: Automatic alerts when price pushes stop arriving or show degradation

## Inflight Funds Management

### Overview

Inflight funds represent capital that is temporarily in transit between the Noble vault and its remote positions, or between positions during rebalancing. Each inflight transaction is tracked by its specific Hyperlane route identifier, allowing precise monitoring of capital flows across different bridge paths. This capital remains fully accounted for in the NAV to ensure accurate vault valuation at all times.

### Inflight Fund Types

1. **Deposit to Position**: USDN being deployed from the Noble vault to a remote ERC-4626 compatible vault
2. **Withdrawal from Position**: USDN returning from redeemed remote vault shares to the Noble vault
3. **Rebalance Between Positions**: USDN moving between remote vaults (via Noble after share redemption)
4. **Pending Deployment**: USDN from deposits awaiting allocation to remote vaults
5. **Pending Withdrawal Distribution**: USDN from redeemed shares awaiting distribution to withdrawal queue
6. **Yield Collection**: Automatic compounding within remote vaults

### Tracking Mechanism

Each inflight transaction maintains:
- **Hyperlane Route ID**: Unique route identifier (e.g., 4000260998 for Noble→Hyperliquid)
- **Transaction ID**: Hyperlane message ID for the specific transfer
- **Source/Destination Domains**: Hyperlane domain IDs for the route endpoints (4000260 for Noble, 998 for Hyperliquid, 8453 for Base, 4000261 for Noble App Layer)
- **Expected Value**: USDN amount sent including estimated bridge fees
- **Current Value**: Last known USDN value for NAV calculation
- **Time Bounds**: Expected arrival time and maximum duration per route
- **Status Updates**: PENDING → CONFIRMED → COMPLETED lifecycle per route
- **Route-Specific Limits**: Maximum exposure allowed per Hyperlane route

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

- **Maximum Duration Limits**: Funds cannot remain inflight indefinitely on any route
- **Per-Route Value Caps**: Limits on inflight exposure for each Hyperlane route
- **Total Value Caps**: Aggregate limits across all routes for the vault
- **Route Diversification**: Use multiple Hyperlane routes to reduce concentration risk
- **Insurance Reserve**: Coverage for potential failures on specific routes
- **Message Authentication**: Verification of Hyperlane message origin and integrity

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
Deposit into Remote Vault
    ↓
Track Remote Vault Shares
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
Initiate Vault Share Redemption
    ↓
Redeem Shares for USDN
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
Redeem Shares from Source Vault
    ↓
Mark as WITHDRAWAL_FROM_POSITION Inflight (Route A)
    ↓
Track via Hyperlane Route ID (e.g., 998_4000260)
    ↓
USDN Arrives at Noble
    ↓
Mark as PENDING_DEPLOYMENT
    ↓
Deploy to Target Vault
    ↓
Mark as DEPOSIT_TO_POSITION Inflight (Route B)
    ↓
Track via Hyperlane Route ID (e.g., 4000260_8453)
    ↓
Deposit into Target Vault & Receive Remote Position Shares
    ↓
Update Share Balances & Clear Route Tracking
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
