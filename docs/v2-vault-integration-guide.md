# Noble Dollar V2 Vault Integration Guide

## Table of Contents

1. [Overview](#overview)
2. [Key Differences from V1](#key-differences-from-v1)
3. [Migration Process](#migration-process)
4. [V2 Vault Operations](#v2-vault-operations)
5. [Code Examples](#code-examples)
6. [API Reference](#api-reference)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)

## Overview

The Noble Dollar V2 vault system introduces a sophisticated share-based accounting mechanism that replaces the legacy lock/unlock system. This upgrade provides:

- **Share-based accounting** with implicit yield tracking
- **NAV (Net Asset Value) pricing** for fair deposit/withdrawal rates
- **Cross-chain position support** for multi-chain strategies
- **Flexible yield preferences** allowing users to forgo yield for redistribution
- **Advanced fee management** through share dilution
- **User-initiated migration** from V1 to V2 with no time pressure

### Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐
│   Legacy V1     │    │    New V2       │
│   Vault         │───▶│   Vault         │
├─────────────────┤    ├─────────────────┤
│ • Lock/Unlock   │    │ • Deposit/Withdraw │
│ • Direct Balance│    │ • Share-based   │
│ • Fixed Rewards │    │ • NAV Pricing   │
└─────────────────┘    └─────────────────┘
```

## Key Differences from V1

| Aspect | V1 (Legacy) | V2 (New) |
|--------|-------------|----------|
| **Accounting** | Direct balance tracking | Share-based accounting |
| **Yield** | Explicit reward distribution | Implicit through share appreciation |
| **Operations** | Lock/Unlock | Deposit/Withdraw |
| **Pricing** | Fixed 1:1 | Dynamic NAV-based |
| **Cross-chain** | Not supported | Full support |
| **Fees** | Limited | Comprehensive with share dilution |
| **Exit** | Immediate | Queue-based for large amounts |

### Share-Based Accounting Explained

In V2, users don't hold tokens directly. Instead:

1. **Deposits**: Users receive shares proportional to their deposit at current NAV
2. **Yield**: Tracked implicitly as share price appreciation
3. **Withdrawals**: Users burn shares to receive tokens at current NAV

**Formula:**
```
Share Price = Total NAV / Total Shares
User Value = User Shares × Share Price
```

## Migration Process

### Migration States

The migration follows a structured state machine:

```
NOT_STARTED → ACTIVE → CLOSING → LOCKED → DEPRECATED
                ↓
           CANCELLED (emergency)
```

### User Migration Flow

1. **Check Eligibility**: Verify user has legacy positions
2. **Preview Migration**: Calculate expected shares and value
3. **Execute Migration**: Submit signed migration transaction
4. **Confirmation**: Receive shares in V2 system

### Migration Commands

```bash
# Check migration status
nobeld query vaults migration-status

# Check user's migration eligibility
nobeld query vaults user-migration-status [user-address]

# Preview migration outcome
nobeld query vaults migration-preview [user-address] --vault-type=FLEXIBLE

# Execute migration
nobeld tx vaults migrate-position \
  --vault-type=FLEXIBLE \
  --amount=0 \
  --forgo-yield=false \
  --from=[user-key]
```

## V2 Vault Operations

### Deposits

Deposits in V2 mint shares based on current NAV:

```bash
# Deposit tokens to receive shares
nobeld tx vaults deposit \
  --vault-type=FLEXIBLE \
  --amount=1000000000 \
  --min-shares=990000000 \
  --forgo-yield=false \
  --from=[user-key]
```

**Parameters:**
- `amount`: Tokens to deposit (in base units)
- `min-shares`: Minimum shares expected (slippage protection)
- `forgo-yield`: Whether to forgo yield for redistribution

### Withdrawals

Two withdrawal methods are available:

#### 1. Immediate Withdrawal
For smaller amounts, immediate withdrawal:

```bash
# Withdraw by burning shares
nobeld tx vaults withdraw \
  --vault-type=FLEXIBLE \
  --shares=1000000000 \
  --min-tokens=990000000 \
  --from=[user-key]
```

#### 2. Exit Queue
For larger amounts, join the exit queue:

```bash
# Request to join exit queue
nobeld tx vaults request-exit \
  --vault-type=FLEXIBLE \
  --shares=10000000000 \
  --from=[user-key]

# Cancel exit request
nobeld tx vaults cancel-exit \
  --exit-id="exit_FLEXIBLE_abc12345_1234567_1640995200" \
  --from=[user-key]
```

### Yield Preferences

Users can choose to forgo yield, redistributing it to other vault participants:

```bash
# Set yield preference
nobeld tx vaults set-yield-preference \
  --vault-type=FLEXIBLE \
  --forgo-yield=true \
  --from=[user-key]
```

## Code Examples

### TypeScript/JavaScript Integration

```typescript
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate";
import { coins } from "@cosmjs/stargate";

class NobleV2VaultClient {
  constructor(
    private client: SigningCosmWasmClient,
    private signer: string
  ) {}

  // Deposit tokens to V2 vault
  async deposit(
    vaultType: "FLEXIBLE" | "STAKED",
    amount: string,
    minShares: string,
    forgoYield: boolean = false
  ) {
    const msg = {
      typeUrl: "/noble.dollar.vaults.v1.MsgDeposit",
      value: {
        signer: this.signer,
        vaultType,
        amount,
        minShares,
        forgoYield,
      },
    };

    return this.client.signAndBroadcast(this.signer, [msg], "auto");
  }

  // Withdraw tokens from V2 vault
  async withdraw(
    vaultType: "FLEXIBLE" | "STAKED",
    shares: string,
    minTokens: string
  ) {
    const msg = {
      typeUrl: "/noble.dollar.vaults.v1.MsgWithdraw",
      value: {
        signer: this.signer,
        vaultType,
        shares,
        minTokens,
      },
    };

    return this.client.signAndBroadcast(this.signer, [msg], "auto");
  }

  // Migrate from V1 to V2
  async migrate(
    vaultType: "FLEXIBLE" | "STAKED",
    amount: string = "0", // 0 = migrate all
    forgoYield: boolean = false
  ) {
    const msg = {
      typeUrl: "/noble.dollar.vaults.v1.MsgMigratePosition",
      value: {
        signer: this.signer,
        vaultType,
        amount,
        forgoYield,
      },
    };

    return this.client.signAndBroadcast(this.signer, [msg], "auto");
  }

  // Query user position
  async getUserPosition(vaultType: "FLEXIBLE" | "STAKED") {
    const response = await this.client.queryContractSmart(
      "vault-query-contract",
      {
        user_position: {
          address: this.signer,
          vault_type: vaultType,
        },
      }
    );
    return response;
  }

  // Preview deposit outcome
  async previewDeposit(vaultType: "FLEXIBLE" | "STAKED", amount: string) {
    const response = await this.client.queryContractSmart(
      "vault-query-contract",
      {
        deposit_preview: {
          vault_type: vaultType,
          amount,
        },
      }
    );
    return response;
  }
}
```

### Go Integration

```go
package main

import (
    "context"
    "fmt"
    
    "cosmossdk.io/math"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "dollar.noble.xyz/v2/types/vaults"
)

// V2VaultClient provides methods to interact with V2 vaults
type V2VaultClient struct {
    client   vaults.VaultV2MsgClient
    querier  vaults.VaultV2QueryClient
    signer   string
}

// Deposit tokens to receive shares
func (c *V2VaultClient) Deposit(ctx context.Context, vaultType vaults.VaultType, amount math.Int, minShares math.Int, forgoYield bool) (*vaults.MsgDepositResponse, error) {
    msg := &vaults.MsgDeposit{
        Signer:    c.signer,
        VaultType: vaultType,
        Amount:    amount,
        MinShares: minShares,
        ForgoYield: forgoYield,
    }
    
    return c.client.Deposit(ctx, msg)
}

// Withdraw tokens by burning shares
func (c *V2VaultClient) Withdraw(ctx context.Context, vaultType vaults.VaultType, shares math.Int, minTokens math.Int) (*vaults.MsgWithdrawResponse, error) {
    msg := &vaults.MsgWithdraw{
        Signer:    c.signer,
        VaultType: vaultType,
        Shares:    shares,
        MinTokens: minTokens,
    }
    
    return c.client.Withdraw(ctx, msg)
}

// Migrate from V1 to V2
func (c *V2VaultClient) Migrate(ctx context.Context, vaultType vaults.VaultType, amount math.Int, forgoYield bool) (*vaults.MsgMigratePositionResponse, error) {
    msg := &vaults.MsgMigratePosition{
        Signer:    c.signer,
        VaultType: vaultType,
        Amount:    amount,
        ForgoYield: forgoYield,
    }
    
    return c.client.MigratePosition(ctx, msg)
}

// Query user position and current value
func (c *V2VaultClient) GetUserPosition(ctx context.Context, address string, vaultType vaults.VaultType) (*vaults.QueryUserPositionResponse, error) {
    req := &vaults.QueryUserPositionRequest{
        Address:   address,
        VaultType: vaultType,
    }
    
    return c.querier.UserPosition(ctx, req)
}

// Preview deposit outcome
func (c *V2VaultClient) PreviewDeposit(ctx context.Context, vaultType vaults.VaultType, amount math.Int) (*vaults.QueryDepositPreviewResponse, error) {
    req := &vaults.QueryDepositPreviewRequest{
        VaultType: vaultType,
        Amount:    amount,
    }
    
    return c.querier.DepositPreview(ctx, req)
}

// Example usage
func ExampleUsage() {
    client := &V2VaultClient{
        // Initialize with gRPC clients
        signer: "noble1...",
    }
    
    ctx := context.Background()
    
    // 1. Preview deposit
    preview, err := client.PreviewDeposit(ctx, vaults.FLEXIBLE, math.NewInt(1000000))
    if err != nil {
        panic(err)
    }
    fmt.Printf("Deposit preview: %s shares for %s tokens\n", 
        preview.EstimatedShares.String(), "1000000")
    
    // 2. Execute deposit
    depositResp, err := client.Deposit(ctx, vaults.FLEXIBLE, 
        math.NewInt(1000000), preview.EstimatedShares, false)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Deposited successfully: %s shares minted\n", 
        depositResp.SharesMinted.String())
    
    // 3. Check position
    position, err := client.GetUserPosition(ctx, client.signer, vaults.FLEXIBLE)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Current position: %s shares worth %s tokens\n",
        position.Position.Shares.String(), position.CurrentValue.String())
}
```

## API Reference

### Message Types

#### MsgDeposit
Deposit tokens to receive shares in V2 vault.

**Fields:**
- `signer` (string): User address
- `vault_type` (VaultType): FLEXIBLE or STAKED
- `amount` (math.Int): Tokens to deposit
- `min_shares` (math.Int): Minimum shares expected
- `forgo_yield` (bool): Whether to forgo yield

#### MsgWithdraw
Withdraw tokens by burning shares.

**Fields:**
- `signer` (string): User address
- `vault_type` (VaultType): FLEXIBLE or STAKED
- `shares` (math.Int): Shares to burn (0 = withdraw all)
- `min_tokens` (math.Int): Minimum tokens expected

#### MsgMigratePosition
Migrate from V1 to V2 vault system.

**Fields:**
- `signer` (string): User address
- `vault_type` (VaultType): FLEXIBLE or STAKED
- `amount` (math.Int): Amount to migrate (0 = migrate all)
- `forgo_yield` (bool): Whether to forgo yield

### Query Types

#### QueryUserPosition
Get user's position and current value.

**Request:**
- `address` (string): User address
- `vault_type` (VaultType): FLEXIBLE or STAKED

**Response:**
- `position` (UserPosition): User's position details
- `current_value` (math.Int): Current token value of shares
- `unrealized_gain` (math.Int): Gain/loss since deposit

#### QueryDepositPreview
Preview deposit outcome before execution.

**Request:**
- `vault_type` (VaultType): FLEXIBLE or STAKED
- `amount` (math.Int): Tokens to deposit

**Response:**
- `estimated_shares` (math.Int): Expected shares to receive
- `share_price` (math.LegacyDec): Current share price
- `fee_amount` (math.Int): Fee to be charged

### Events

#### vault_deposit
Emitted when user deposits to V2 vault.

**Attributes:**
- `user`: User address
- `vault_type`: FLEXIBLE or STAKED
- `amount`: Tokens deposited
- `shares_minted`: Shares received
- `fee_charged`: Fee amount

#### vault_withdrawal
Emitted when user withdraws from V2 vault.

**Attributes:**
- `user`: User address
- `vault_type`: FLEXIBLE or STAKED
- `shares_burned`: Shares burned
- `tokens_withdrawn`: Tokens received
- `fee_charged`: Fee amount

#### vault_migration_completed
Emitted when user migrates from V1 to V2.

**Attributes:**
- `user`: User address
- `vault_type`: FLEXIBLE or STAKED
- `shares_received`: Shares minted in V2
- `principal_migrated`: Principal amount migrated
- `rewards_migrated`: Rewards amount migrated

## Best Practices

### 1. Always Use Slippage Protection

```bash
# Good: Set reasonable slippage tolerance
nobeld tx vaults deposit \
  --amount=1000000000 \
  --min-shares=990000000  # 1% slippage tolerance

# Bad: No slippage protection
nobeld tx vaults deposit \
  --amount=1000000000 \
  --min-shares=0
```

### 2. Preview Before Executing

Always preview operations before execution:

```typescript
// Preview first
const preview = await client.previewDeposit("FLEXIBLE", "1000000");
console.log(`Expected shares: ${preview.estimated_shares}`);

// Then execute with slippage protection
const minShares = BigInt(preview.estimated_shares) * 99n / 100n; // 1% slippage
await client.deposit("FLEXIBLE", "1000000", minShares.toString());
```

### 3. Monitor NAV Changes

NAV can change between preview and execution:

```go
// Check time since last NAV update
vaultState, _ := client.GetVaultState(ctx, vaults.FLEXIBLE)
timeSinceUpdate := time.Since(vaultState.LastNavUpdate)

if timeSinceUpdate > time.Hour {
    fmt.Println("Warning: NAV data may be stale")
}
```

### 4. Handle Exit Queue for Large Withdrawals

For large withdrawals, use the exit queue:

```bash
# For amounts > 10% of vault, use exit queue
nobeld tx vaults request-exit \
  --vault-type=FLEXIBLE \
  --shares=10000000000
```

### 5. Consider Yield Preferences

Think carefully about yield preferences:

```typescript
// For long-term holders who want to support the community
await client.setYieldPreference("FLEXIBLE", true); // Forgo yield

// For yield-seeking investors
await client.setYieldPreference("FLEXIBLE", false); // Keep yield
```

## Troubleshooting

### Common Errors

#### "insufficient shares received"
**Cause:** Slippage tolerance too tight or NAV changed between preview and execution.

**Solution:** 
- Increase slippage tolerance
- Preview again before execution
- Check for recent NAV updates

#### "deposits are currently paused"
**Cause:** Vault deposits are paused due to maintenance or circuit breaker.

**Solution:**
- Check vault status: `nobeld query vaults vault-state --vault-type=FLEXIBLE`
- Wait for deposits to be re-enabled
- Monitor official announcements

#### "migration not active"
**Cause:** Migration is not in ACTIVE or CLOSING state.

**Solution:**
- Check migration status: `nobeld query vaults migration-status`
- Wait for migration period to open
- Use emergency withdrawal if migration is cancelled

#### "circuit breaker is active"
**Cause:** Large NAV deviation triggered circuit breaker.

**Solution:**
- Wait for manual review and circuit breaker reset
- Monitor official channels for updates
- Consider emergency procedures if needed

### Debugging Tips

#### 1. Check Transaction Details
```bash
# Get detailed transaction info
nobeld query tx [tx-hash] --output=json
```

#### 2. Monitor Events
```bash
# Filter vault-related events
nobeld query tx [tx-hash] --output=json | jq '.logs[].events[] | select(.type | contains("vault"))'
```

#### 3. Verify Balances
```bash
# Check token balance
nobeld query bank balances [address]

# Check vault position
nobeld query vaults user-position [address] --vault-type=FLEXIBLE
```

#### 4. Check Share Price History
```bash
# Monitor share price changes
nobeld query vaults share-price --vault-type=FLEXIBLE
```

### Support Resources

- **Documentation**: [Noble Dollar V2 Docs](https://docs.noble.xyz/dollar/v2)
- **Discord**: [Noble Protocol Discord](https://discord.gg/noble)
- **GitHub**: [Noble Dollar Repository](https://github.com/noble-assets/dollar)
- **API Reference**: [gRPC API Docs](https://buf.build/noble/dollar)

### Emergency Procedures

If you encounter critical issues:

1. **Check Migration Status**: Ensure migration is still active
2. **Emergency Withdrawal**: Use if migration is cancelled
3. **Contact Support**: Reach out via Discord or GitHub
4. **Monitor Announcements**: Follow official channels for updates

```bash
# Emergency withdrawal from legacy vault (if migration cancelled)
nobeld tx vaults emergency-withdraw-legacy \
  --vault-type=FLEXIBLE \
  --position-indices=1,2,3 \
  --from=[user-key]
```

---

This guide covers the essential aspects of integrating with the Noble Dollar V2 vault system. For the latest updates and detailed API documentation, please refer to the official Noble Protocol documentation.