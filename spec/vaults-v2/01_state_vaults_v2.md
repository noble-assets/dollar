# Vaults V2 State

## NAV (Net Asset Value)

The `NAV` field is a [`collections.Item`][item] that stores the current Net Asset Value of the single Noble vault, representing the total value of all assets including its remote positions and inflight funds (`math.LegacyDec`).

The NAV calculation includes:
- Local assets (held in the Noble vault)
- Remote position values (capital deployed from the vault)
- Inflight funds (in transit between the vault and remote positions)
- Minus pending liabilities (withdrawals to be paid)

```go
const NAVKey = []byte("vaults/v2/nav")
```

### NAV Calculation Formula

```
Total NAV = Local Assets 
          + Σ(Remote Position Values) 
          + Σ(Inflight Funds Values)
          - Pending Withdrawal Liabilities

NAV per Share = Total NAV / Total Outstanding Shares
```

## LastNAVUpdate

The `LastNAVUpdate` field is a [`collections.Item`][item] that stores the timestamp of the last NAV calculation (`time.Time`).

```go
const LastNAVUpdateKey = []byte("vaults/v2/last_nav_update")
```

## RemotePositions

The `RemotePositions` field is a mapping ([`collections.Map`][map]) between position ID (`uint64`) and a `vaults.v2.RemotePosition` value. The single Noble vault maintains multiple remote positions across different protocols and chains.

```go
const RemotePositionPrefix = []byte("vaults/v2/remote_position/")
```

### RemotePosition Structure

```go
type RemotePosition struct {
    PositionID       uint64
    VaultAddress     []byte  // Address of the ERC-4626 compatible vault (Boring Vault or other)
    ChainID          uint32  // Hyperlane Domain ID (e.g., 998 for Hyperliquid, 8453 for Base, 4000261 for Noble App Layer)
    SharesHeld       math.Int // Number of vault shares held
    Principal        math.Int // USDN amount initially deposited
    LastUpdatedNAV   math.LegacyDec
    LastUpdateTime   time.Time
    Status           PositionStatus // ACTIVE, WITHDRAWING, CLOSED
}
```

## InflightFunds

The `InflightFunds` field is a mapping ([`collections.Map`][map]) between a Hyperlane route identifier (`uint32`) and an `vaults.v2.InflightFund` value. These represent funds that are currently in transit between the Noble vault and its remote positions or between positions, tracked per Hyperlane route.

```go
const InflightFundsPrefix = []byte("vaults/v2/inflight_funds/")
```

### InflightFund Structure

```go
type InflightFund struct {
    HyperlaneRouteID uint32       // Hyperlane route identifier
    TransactionID    string       // Hyperlane message ID
    Type             InflightType // DEPOSIT_TO_POSITION, WITHDRAWAL_FROM_POSITION, REBALANCE_BETWEEN_POSITIONS, PENDING_DEPLOYMENT, PENDING_WITHDRAWAL_DISTRIBUTION
    Amount           math.Int     // Always in USDN
    SourceDomain     uint32       // Hyperlane source domain
    DestDomain       uint32       // Hyperlane destination domain
    SourcePosition   *uint64      // Optional: source position ID if from position
    DestPosition     *uint64      // Optional: destination position ID if to position
    InitiatedAt      time.Time
    ExpectedAt       time.Time
    Status           InflightStatus // PENDING, CONFIRMED, COMPLETED, FAILED
    ValueAtInitiation math.Int    // USDN value when initiated
}
```

## InflightRoutes

The `InflightRoutes` field is a [`collections.Item`][item] that stores a list of all active Hyperlane route IDs (`[]uint32`) with inflight funds. This allows quick lookup of all inflight funds for NAV calculation.

```go
const InflightRoutesKey = []byte("vaults/v2/inflight_routes")
```

## TotalInflightValue

The `TotalInflightValue` field is a [`collections.Item`][item] that stores the total value of all inflight funds across all Hyperlane routes (`math.Int`). This is cached for efficient NAV calculations and always denominated in USDN.

```go
const TotalInflightValueKey = []byte("vaults/v2/total_inflight_value")
```

## InflightValueByRoute

The `InflightValueByRoute` field is a mapping ([`collections.Map`][map]) between Hyperlane route ID (`uint32`) and the total inflight value on that route (`math.Int`). This enables per-route exposure tracking.

```go
const InflightValueByRoutePrefix = []byte("vaults/v2/inflight_value_by_route/")
```

## PendingDeploymentFunds

The `PendingDeploymentFunds` field is a [`collections.Item`][item] that stores the amount of USDN received from deposits but not yet deployed to remote positions (`math.Int`).

```go
const PendingDeploymentFundsKey = []byte("vaults/v2/pending_deployment")
```

## PendingWithdrawalDistribution

The `PendingWithdrawalDistribution` field is a [`collections.Item`][item] that stores the amount of USDN returned from remote positions but not yet distributed to withdrawal claimants (`math.Int`).

```go
const PendingWithdrawalDistributionKey = []byte("vaults/v2/pending_withdrawal_dist")
```

## WithdrawalQueue

The `WithdrawalQueue` field is a mapping ([`collections.Map`][map]) between a withdrawal request ID (`uint64`) and a `vaults.v2.WithdrawalRequest` value.

```go
const WithdrawalQueuePrefix = []byte("vaults/v2/withdrawal_queue/")
```

### WithdrawalRequest Structure

```go
type WithdrawalRequest struct {
    RequestID        uint64
    User             []byte
    SharesAmount     math.Int
    RequestedAmount  math.Int  // Amount in $USDN requested
    NAVAtRequest     math.LegacyDec
    Timestamp        time.Time
    Status           WithdrawalStatus // PENDING, PROCESSING, CLAIMABLE, CLAIMED
    FulfilledAmount  math.Int  // Actual amount available for claim
}
```

## WithdrawalQueueSequence

The `WithdrawalQueueSequence` field is a [`collections.Sequence`][sequence] that generates unique withdrawal request IDs.

```go
const WithdrawalQueueSequenceKey = []byte("vaults/v2/withdrawal_queue_seq")
```

## PendingWithdrawals

The `PendingWithdrawals` field is a [`collections.Item`][item] that stores the total amount of pending withdrawals for the vault (`math.Int`).

```go
const PendingWithdrawalsKey = []byte("vaults/v2/pending_withdrawals")
```

## UserShares

The `UserShares` field is a mapping ([`collections.Map`][map]) between user address (`[]byte`) and their share balance in the vault (`math.Int`).

```go
const UserSharesPrefix = []byte("vaults/v2/user_shares/")
```

## TotalShares

The `TotalShares` field is a [`collections.Item`][item] that stores the total outstanding shares for the vault (`math.Int`).

```go
const TotalSharesKey = []byte("vaults/v2/total_shares")
```

## DepositLimits

The `DepositLimits` field is a [`collections.Item`][item] that stores the deposit limits for the vault (`vaults.v2.DepositLimit`).

```go
const DepositLimitsKey = []byte("vaults/v2/deposit_limits")
```

### DepositLimit Structure

```go
type DepositLimit struct {
    MaxDepositPerUser     math.Int
    MaxDepositPerBlock    math.Int
    MaxTotalDeposits      math.Int
    MinDepositAmount      math.Int
    DepositCooldownBlocks uint64
}
```

## UserDepositHistory

The `UserDepositHistory` field is a mapping ([`collections.Map`][map]) between a composite key of user address and block height (`[]byte`, `int64`) and deposit amount (`math.Int`). Used for malicious deposit detection.

```go
const UserDepositHistoryPrefix = []byte("vaults/v2/user_deposit_history/")
```

## DepositVelocity

The `DepositVelocity` field is a mapping ([`collections.Map`][map]) between user address (`[]byte`) and a `vaults.v2.DepositVelocity` value for tracking deposit patterns and defending against manipulation.

```go
const DepositVelocityPrefix = []byte("vaults/v2/deposit_velocity/")
```

### DepositVelocity Structure

```go
type DepositVelocity struct {
    LastDepositBlock      int64
    RecentDepositCount    uint32  // Rolling count over time window
    RecentDepositVolume   math.Int // Rolling volume over time window
    SuspiciousActivityFlag bool
    TimeWindowBlocks      int64   // e.g., 1000 blocks
}
```

## RemotePositionOracles

The `RemotePositionOracles` field is a mapping ([`collections.Map`][map]) between a composite key of vault address and chain ID (`[]byte`, `uint32`) and oracle configuration (`vaults.v2.OracleConfig`). This is used to track the value of shares in remote ERC-4626 compatible vaults.

```go
const RemotePositionOraclesPrefix = []byte("vaults/v2/remote_position_oracles/")
```

### OracleConfig Structure

```go
type OracleConfig struct {
    OracleType       string  // "hyperlane", "chainlink", "internal"
    UpdateFrequency  int64   // blocks between required updates
    MaxStaleness     int64   // maximum blocks before data considered stale
    MinConfirmations uint32  // minimum confirmations required
    TrustedSources   [][]byte // addresses of trusted oracle sources
}
```

## VaultConfiguration

The `VaultConfiguration` field is a [`collections.Item`][item] that stores the configuration for the Noble vault (`vaults.v2.VaultConfig`).

```go
const VaultConfigurationKey = []byte("vaults/v2/vault_config")
```

### VaultConfig Structure

```go
type VaultConfig struct {
    Name                  string
    MaxRemotePositions    uint32
    AllowedChains         []uint32  // Hyperlane domain IDs
    AllowedVaultAddresses map[uint32][][]byte // Chain ID -> List of approved vault addresses
    RebalanceThreshold    math.LegacyDec // Percentage deviation before rebalance
    WithdrawalDelayBlocks int64
    ManagementFee         math.LegacyDec // Annual fee as percentage
    PerformanceFee        math.LegacyDec // Fee on profits
    EmergencyMode         bool
    MaxInflightDuration   int64          // Max blocks funds can be inflight
    InflightValueCap      math.Int       // Maximum value allowed inflight
}
```

[item]: https://docs.cosmos.network/v0.50/build/packages/collections#item
[map]: https://docs.cosmos.network/v0.50/build/packages/collections#map
[sequence]: https://docs.cosmos.network/v0.50/build/packages/collections#sequence
