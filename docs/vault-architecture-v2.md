# Noble Dollar Vault Architecture V2

## Overview

The Noble Dollar vault system has been redesigned with a simplified, clean architecture that separates legacy V1 vaults from the new V2 share-based system. This document outlines the architecture, migration approach, and technical implementation.

## System Architecture

### High-Level Design

```
┌─────────────────┐    Manual Migration    ┌─────────────────┐
│   V1 Vaults     │   (Withdraw + Deposit)  │   V2 Vaults     │
│  (Legacy)       │ ─────────────────────> │ (Share-based)   │
│                 │                        │                 │
│ • Withdrawal    │                        │ • Deposits      │
│ • Rewards       │                        │ • Withdrawals   │
│ • Queries       │                        │ • Share system  │
│                 │                        │ • Exit queues   │
└─────────────────┘                        └─────────────────┘
```

### V1 Vaults (Legacy System)

**Status: Withdrawal-Only**

V1 vaults represent the legacy direct-balance tracking system. They are now in maintenance mode with limited operations:

#### Supported Operations:
- ✅ **Withdrawals**: Full or partial withdrawals from existing positions
- ✅ **Rewards Claiming**: Claim accrued rewards
- ✅ **Balance Queries**: Check position balances and rewards
- ❌ **New Deposits**: No new deposits accepted

#### Key Characteristics:
- Direct balance tracking (no shares)
- Position-based accounting
- Immediate rewards calculation
- No automatic compounding

### V2 Vaults (Share-Based System)

**Status: Fully Operational**

V2 vaults implement a modern share-based accounting system with automatic yield compounding:

#### Supported Operations:
- ✅ **Deposits**: Receive shares based on current NAV
- ✅ **Withdrawals**: Burn shares to receive tokens (flexible vaults)
- ✅ **Exit Requests**: Queue-based withdrawal for staked vaults
- ✅ **Yield Management**: User-controlled yield preferences
- ✅ **NAV Updates**: Authority-controlled asset valuations

#### Key Characteristics:
- Share-based accounting: `sharePrice = totalNAV / totalShares`
- Automatic yield compounding
- Implicit yield tracking through share price appreciation
- Exit queues for staked vault unbonding
- Slippage protection
- Fee management through share dilution

## Migration Strategy

### Simplified Approach

Instead of complex automated migration, we implement a **manual migration** approach:

1. **User Action Required**: Users must actively migrate their positions
2. **Two-Step Process**: 
   - Step 1: Withdraw from V1 vault
   - Step 2: Deposit into V2 vault
3. **No Time Pressure**: V1 withdrawals remain available indefinitely
4. **Clean Separation**: No complex state synchronization needed

### Migration Benefits

#### For Users:
- **Full Control**: Users decide when and how much to migrate
- **Transparency**: Clear understanding of the process
- **No Rush**: No forced migration deadlines
- **Yield Optimization**: Choose timing based on market conditions

#### For Protocol:
- **Reduced Complexity**: No complex migration state machines
- **Lower Risk**: No atomic cross-system operations
- **Easier Testing**: Independent system validation
- **Cleaner Code**: Simplified keeper logic

## Technical Implementation

### Protocol Buffer Structure

```
noble/dollar/vaults/
├── v1/                    # Legacy system (withdrawal-only)
│   ├── vaults.proto       # V1 position definitions
│   ├── events.proto       # V1 events
│   ├── genesis.proto      # V1 genesis state
│   ├── query.proto        # V1 queries
│   └── tx.proto           # V1 transactions (withdrawals only)
└── v2/                    # New share-based system
    ├── vaults.proto       # V2 share and position definitions
    ├── nav.proto          # NAV calculation structures
    ├── fees.proto         # Fee management
    ├── cross_chain.proto  # Cross-chain operations
    ├── events.proto       # V2 events
    ├── genesis.proto      # V2 genesis state
    ├── query.proto        # V2 queries
    └── tx.proto           # V2 transactions
```

### Share-Based Accounting

#### Core Formula

```
sharePrice = totalNAV / totalShares
userValue = userShares × sharePrice
```

#### Deposit Process

```
1. Calculate current share price
2. Apply deposit fees: netDeposit = deposit - fees
3. Calculate shares: shares = netDeposit / sharePrice
4. Mint shares to user
5. Update total NAV and shares
```

#### Withdrawal Process

```
1. Calculate current share price
2. Calculate gross amount: grossAmount = shares × sharePrice
3. Apply withdrawal fees: netAmount = grossAmount - fees
4. Burn user shares
5. Transfer tokens to user
6. Update total NAV and shares
```

### Vault Types

#### Flexible Vaults (V2)
- **Immediate liquidity**: Instant withdrawals
- **Share-based**: Automatic yield compounding
- **Lower yields**: Trade-off for liquidity

#### Staked Vaults (V2)
- **Higher yields**: Rewards for illiquidity
- **Exit queues**: Unbonding period required
- **Batch processing**: Efficient withdrawal handling

## Usage Examples

### User Migration Workflow

```bash
# Step 1: Check V1 position
dollard query vaults position $USER_ADDRESS FLEXIBLE

# Step 2: Withdraw from V1
dollard tx vaults withdraw $USER_ADDRESS FLEXIBLE $AMOUNT

# Step 3: Deposit into V2
dollard tx vaults-v2 deposit $USER_ADDRESS FLEXIBLE $AMOUNT --receive-yield=true
```

### V2 Operations

```bash
# Deposit into V2 flexible vault
dollard tx vaults-v2 deposit $USER FLEXIBLE 1000000 --receive-yield=true

# Check position
dollard query vaults-v2 position $USER FLEXIBLE

# Withdraw (flexible vault)
dollard tx vaults-v2 withdraw $USER FLEXIBLE 500000

# Request exit (staked vault)
dollard tx vaults-v2 request-exit $USER STAKED 1000000

# Cancel exit request
dollard tx vaults-v2 cancel-exit $USER STAKED $EXIT_REQUEST_ID
```

## Administrative Operations

### NAV Management

```bash
# Update vault NAV (authority only)
dollard tx vaults-v2 update-nav $AUTHORITY FLEXIBLE 1050000000 "Daily yield distribution"

# Process exit queue (authority only)
dollard tx vaults-v2 process-exit-queue $AUTHORITY STAKED 10
```

### Configuration

```bash
# Update vault parameters
dollard tx vaults-v2 update-params $AUTHORITY --min-deposit=1000000

# Update vault configuration
dollard tx vaults-v2 update-vault-config $AUTHORITY FLEXIBLE --deposit-fee=50
```

## Safety Features

### V2 System Protections

1. **Slippage Protection**: Minimum amount/share guarantees
2. **NAV Change Limits**: Maximum NAV change per update (basis points)
3. **Rate Limiting**: Maximum exit requests per block
4. **Circuit Breakers**: Emergency pause mechanisms
5. **Fee Caps**: Maximum fee rates

### Operational Security

1. **Authority Controls**: Multi-sig authority for critical operations
2. **Parameter Governance**: On-chain governance for parameter updates
3. **Emergency Procedures**: Quick response for critical issues
4. **Audit Trail**: Comprehensive event logging

## Benefits of New Architecture

### User Experience
- **Simpler Understanding**: Clear V1 vs V2 distinction
- **Automatic Compounding**: No manual reward claiming needed
- **Flexible Timing**: Self-paced migration
- **Yield Optimization**: Share price appreciation

### Protocol Benefits
- **Reduced Complexity**: Simpler keeper logic
- **Better Capital Efficiency**: Share-based pooling
- **Scalability**: More efficient operations
- **Maintainability**: Clean code separation

### Developer Experience
- **Clear APIs**: Distinct V1/V2 interfaces
- **Independent Testing**: Separate test suites
- **Easier Integration**: Well-defined boundaries
- **Future Extensibility**: Clean foundation for new features

## Future Considerations

### V1 Sunset Planning
- Monitor V1 usage metrics
- Communicate deprecation timeline
- Provide migration incentives
- Plan final sunset procedures

### V2 Enhancements
- Cross-chain vault support
- Advanced yield strategies
- Governance integration
- DeFi protocol integrations

## Conclusion

The simplified V1/V2 architecture provides a clean migration path while maintaining user control and reducing system complexity. The manual migration approach eliminates complex state synchronization while providing users with full transparency and control over their vault positions.

This design establishes a solid foundation for future vault system enhancements while ensuring a smooth transition from the legacy system.