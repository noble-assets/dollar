# Vaults V2 State

## NAV (Net Asset Value)

The `NAV` field is a [`collections.Item`][item] that stores the current Net Asset Value of the vault, representing the total value of all assets including remote positions (`math.LegacyDec`).

```go
const NAVKey = []byte("vaults/v2/nav")
```

## LastNAVUpdate

The `LastNAVUpdate` field is a [`collections.Item`][item] that stores the timestamp of the last NAV calculation (`time.Time`).

```go
const LastNAVUpdateKey = []byte("vaults/v2/last_nav_update")
```

## RemotePositions

The `RemotePositions` field is a mapping ([`collections.Map`][map]) between a composite key of vault ID and position ID (`string`, `uint64`) and a `vaults.v2.RemotePosition` value. A single vault can maintain multiple remote positions across different protocols and chains.

```go
const RemotePositionPrefix = []byte("vaults/v2/remote_position/")
```

### RemotePosition Structure

```go
type RemotePosition struct {
    PositionID       uint64
    VaultID          string
    Protocol         string  // e.g., "aave", "compound", "morpho"
    ChainID          uint16  // Hyperlane/IBC Chain ID
    AssetAddress     []byte  // Remote asset address
    Principal        math.Int
    LastUpdatedNAV   math.LegacyDec
    LastUpdateTime   time.Time
    Status           PositionStatus // ACTIVE, WITHDRAWING, CLOSED
}
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
    VaultID          string
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

The `PendingWithdrawals` field is a [`collections.Item`][item] that stores the total amount of pending withdrawals across all vaults (`math.Int`).

```go
const PendingWithdrawalsKey = []byte("vaults/v2/pending_withdrawals")
```

## VaultShares

The `VaultShares` field is a mapping ([`collections.Map`][map]) between a composite key of user address and vault ID (`[]byte`, `string`) and their share balance (`math.Int`).

```go
const VaultSharesPrefix = []byte("vaults/v2/vault_shares/")
```

## TotalShares

The `TotalShares` field is a mapping ([`collections.Map`][map]) between vault ID (`string`) and the total outstanding shares for that vault (`math.Int`).

```go
const TotalSharesPrefix = []byte("vaults/v2/total_shares/")
```

## DepositLimits

The `DepositLimits` field is a mapping ([`collections.Map`][map]) between vault ID (`string`) and a `vaults.v2.DepositLimit` value.

```go
const DepositLimitsPrefix = []byte("vaults/v2/deposit_limits/")
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

The `RemotePositionOracles` field is a mapping ([`collections.Map`][map]) between a composite key of protocol and chain ID (`string`, `uint16`) and oracle configuration (`vaults.v2.OracleConfig`).

```go
const RemotePositionOraclesPrefix = []byte("vaults/v2/remote_position_oracles/")
```

### OracleConfig Structure

```go
type OracleConfig struct {
    OracleType       string  // "wormhole", "chainlink", "internal"
    UpdateFrequency  int64   // blocks between required updates
    MaxStaleness     int64   // maximum blocks before data considered stale
    MinConfirmations uint32  // minimum confirmations required
    TrustedSources   [][]byte // addresses of trusted oracle sources
}
```

## VaultConfiguration

The `VaultConfiguration` field is a mapping ([`collections.Map`][map]) between vault ID (`string`) and a `vaults.v2.VaultConfig` value.

```go
const VaultConfigurationPrefix = []byte("vaults/v2/vault_config/")
```

### VaultConfig Structure

```go
type VaultConfig struct {
    VaultID               string
    Name                  string
    MaxRemotePositions    uint32
    AllowedProtocols      []string
    AllowedChains         []uint16
    RebalanceThreshold    math.LegacyDec // Percentage deviation before rebalance
    WithdrawalDelayBlocks int64
    ManagementFee         math.LegacyDec // Annual fee as percentage
    PerformanceFee        math.LegacyDec // Fee on profits
    EmergencyMode         bool
}
```

[item]: https://docs.cosmos.network/v0.50/build/packages/collections#item
[map]: https://docs.cosmos.network/v0.50/build/packages/collections#map
[sequence]: https://docs.cosmos.network/v0.50/build/packages/collections#sequence
