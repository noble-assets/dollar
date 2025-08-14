# Vaults V2 Queries

## NAV

**Endpoint**: `/noble/dollar/vaults/v2/nav`

Retrieves the current Net Asset Value and related metrics for the Noble vault.

```json
{
  "nav": "10500000",
  "nav_per_share": "1.050000000000000000",
  "last_update": "2024-01-15T14:30:00Z",
  "total_shares": "10000000",
  "local_assets": "2000000",
  "remote_positions_value": "8500000",
  "inflight_funds_value": "150000",
  "pending_withdrawals": "500000",
  "nav_breakdown": {
    "local": "2000000",
    "remote_positions": "8500000",
    "inflight": "150000",
    "liabilities": "-500000",
    "total": "10150000"
  }
}
```

### Response

- `nav` — Total Net Asset Value of the Noble vault in $USDN.
- `nav_per_share` — Current value per share.
- `last_update` — Timestamp of the last NAV update.
- `total_shares` — Total outstanding shares.
- `local_assets` — Value of assets held locally.
- `remote_positions_value` — Combined value of all remote positions.
- `inflight_funds_value` — Total value of funds currently in transit.
- `pending_withdrawals` — Total value locked for pending withdrawals.
- `nav_breakdown` — Detailed breakdown of NAV components.

## RemotePositions

**Endpoint**: `/noble/dollar/vaults/v2/remote_positions`

Retrieves all remote positions for the Noble vault.

```json
{
  "positions": [
    {
      "position_id": "1",
      "protocol": "aave-v3",
      "chain_id": 2,
      "asset_address": "0x...",
      "principal": "3000000",
      "current_value": "3150000",
      "apy": "5.2",
      "status": "ACTIVE",
      "last_update": "2024-01-15T14:00:00Z"
    },
    {
      "position_id": "2",
      "protocol": "compound-v3",
      "chain_id": 1,
      "asset_address": "0x...",
      "principal": "5000000",
      "current_value": "5350000",
      "apy": "7.1",
      "status": "ACTIVE",
      "last_update": "2024-01-15T14:00:00Z"
    }
  ],
  "total_positions": 2,
  "total_value": "8500000"
}
```

### Response

- `positions` — Array of remote position details.
- `total_positions` — Number of active remote positions.
- `total_value` — Combined value of all positions.

## WithdrawalQueue

**Endpoint**: `/noble/dollar/vaults/v2/withdrawal_queue`

Retrieves the current state of the withdrawal queue for the Noble vault.

```json
{
  "pending_requests": [
    {
      "request_id": "42",
      "user": "noble1user1",
      "shares": "1000",
      "requested_amount": "1050000",
      "nav_at_request": "1.050000000000000000",
      "timestamp": "2024-01-15T12:00:00Z",
      "status": "PENDING",
      "position": 1
    },
    {
      "request_id": "43",
      "user": "noble1user2",
      "shares": "500",
      "requested_amount": "525000",
      "nav_at_request": "1.050000000000000000",
      "timestamp": "2024-01-15T12:30:00Z",
      "status": "PENDING",
      "position": 2
    }
  ],
  "total_pending": "1575000",
  "available_liquidity": "500000",
  "queue_length": 2,
  "estimated_processing_time": "86400"
}
```

### Response

- `pending_requests` — Array of pending withdrawal requests in queue order.
- `total_pending` — Total $USDN value of pending withdrawals.
- `available_liquidity` — Current liquidity available for processing.
- `queue_length` — Number of requests in queue.
- `estimated_processing_time` — Estimated seconds until queue processing.

## UserWithdrawals

**Endpoint**: `/noble/dollar/vaults/v2/user_withdrawals/{user}`

Retrieves all withdrawal requests for a specific user.

```json
{
  "withdrawals": [
    {
      "request_id": "42",
      "shares": "1000",
      "requested_amount": "1050000",
      "fulfilled_amount": "1050000",
      "status": "CLAIMABLE",
      "timestamp": "2024-01-15T12:00:00Z",
      "claimable_at": "2024-01-16T12:00:00Z"
    },
    {
      "request_id": "38",
      "shares": "500",
      "requested_amount": "500000",
      "fulfilled_amount": "500000",
      "status": "CLAIMED",
      "timestamp": "2024-01-14T10:00:00Z",
      "claimed_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total_pending": "0",
  "total_claimable": "1050000",
  "total_claimed": "500000"
}
```

### Arguments

- `user` — The address of the user.

### Response

- `withdrawals` — Array of all user's withdrawal requests.
- `total_pending` — Total value of pending withdrawals.
- `total_claimable` — Total value ready to claim.
- `total_claimed` — Total value already claimed.

## UserShares

**Endpoint**: `/noble/dollar/vaults/v2/user_shares/{user}`

Retrieves share balance and value for a user in the Noble vault.

```json
{
  "shares": "10000",
  "share_value": "10500000",
  "nav_per_share": "1.050000000000000000",
  "unrealized_gain": "500000",
  "locked_shares": "1000"
}
```

### Arguments

- `user` — The address of the user.

### Response

- `shares` — User's share balance in the vault.
- `share_value` — Current value of user's shares.
- `nav_per_share` — Current NAV per share.
- `unrealized_gain` — User's unrealized gain.
- `locked_shares` — Shares locked for pending withdrawals.

## VaultStats

**Endpoint**: `/noble/dollar/vaults/v2/stats`

Retrieves comprehensive statistics for the Noble vault.

```json
{
  "total_deposits": "100000000",
  "total_withdrawals": "20000000",
  "total_shares": "76190476",
  "unique_depositors": "1250",
  "nav": "10500000",
  "apy_7d": "6.5",
  "apy_30d": "6.2",
  "management_fee": "0.020000000000000000",
  "performance_fee": "0.100000000000000000",
  "utilization_rate": "0.850000000000000000",
  "remote_positions_count": 3,
  "average_position_apy": "6.8",
  "deposit_limits": {
    "max_per_user": "1000000",
    "max_per_block": "5000000",
    "max_total": "100000000",
    "current_total": "80000000"
  }
}
```

### Response

- `total_deposits` — Historical total deposits.
- `total_withdrawals` — Historical total withdrawals.
- `total_shares` — Current outstanding shares.
- `unique_depositors` — Number of unique depositors.
- `nav` — Current Net Asset Value.
- `apy_7d` — 7-day average APY.
- `apy_30d` — 30-day average APY.
- `management_fee` — Annual management fee percentage.
- `performance_fee` — Performance fee percentage.
- `utilization_rate` — Percentage of capital deployed.
- `remote_positions_count` — Number of active remote positions.
- `average_position_apy` — Weighted average APY across positions.
- `deposit_limits` — Current deposit limit configuration and usage.

## DepositVelocity

**Endpoint**: `/noble/dollar/vaults/v2/deposit_velocity/{user}`

Retrieves deposit velocity metrics for malicious behavior detection.

```json
{
  "user": "noble1user",
  "last_deposit_block": "1234567",
  "recent_deposit_count": 5,
  "recent_deposit_volume": "5000000",
  "time_window_blocks": 1000,
  "suspicious_activity_flag": false,
  "cooldown_remaining_blocks": 0,
  "velocity_score": "0.250000000000000000"
}
```

### Arguments

- `user` — The address of the user.

### Response

- `user` — The user address.
- `last_deposit_block` — Block height of last deposit.
- `recent_deposit_count` — Number of deposits in time window.
- `recent_deposit_volume` — Total volume in time window.
- `time_window_blocks` — Size of the monitoring window.
- `suspicious_activity_flag` — Whether suspicious activity detected.
- `cooldown_remaining_blocks` — Blocks until next deposit allowed.
- `velocity_score` — Normalized velocity score (0-1).

## SimulateDeposit

**Endpoint**: `/noble/dollar/vaults/v2/simulate_deposit`

Simulates a deposit to show expected shares and checks.

```json
{
  "amount": "1000000",
  "user": "noble1user",
  "expected_shares": "952380",
  "nav_per_share": "1.050000000000000000",
  "checks": {
    "within_user_limit": true,
    "within_block_limit": true,
    "within_total_limit": true,
    "cooldown_passed": true,
    "velocity_check_passed": true
  },
  "warnings": []
}
```

### Query Parameters

- `amount` — The amount to simulate depositing.
- `user` — The user address for limit checks.

### Response

- `expected_shares` — Number of shares that would be minted.
- `nav_per_share` — Current NAV per share used in calculation.
- `checks` — Results of various limit and safety checks.
- `warnings` — Any warnings about the deposit.

## SimulateWithdrawal

**Endpoint**: `/noble/dollar/vaults/v2/simulate_withdrawal`

Simulates a withdrawal to show expected proceeds and queue time.

```json
{
  "shares": "1000",
  "user": "noble1user",
  "expected_amount": "1050000",
  "nav_per_share": "1.050000000000000000",
  "queue_position": 5,
  "estimated_fulfillment_time": "172800",
  "available_liquidity": "500000",
  "ahead_in_queue": "2500000"
}
```

### Query Parameters

- `shares` — The number of shares to simulate redeeming.
- `user` — The user address.

### Response

- `expected_amount` — Expected $USDN to receive.
- `nav_per_share` — Current NAV per share.
- `queue_position` — Position in withdrawal queue.
- `estimated_fulfillment_time` — Estimated seconds until fulfillment.
- `available_liquidity` — Current available liquidity.
- `ahead_in_queue` — Total value of requests ahead in queue.

## InflightFunds

**Endpoint**: `/noble/dollar/vaults/v2/inflight_funds`

Retrieves all inflight funds for the Noble vault, including funds in transit between the vault and remote positions.

```json
{
  "inflight_funds": [
    {
      "transaction_id": "hyperlane-tx-123",
      "type": "DEPOSIT_TO_POSITION",
      "amount": "100000",
      "current_value": "100000",
      "source": "noble",
      "destination": "ethereum-aave-v3",
      "initiated_at": "2024-01-15T14:00:00Z",
      "expected_at": "2024-01-15T14:30:00Z",
      "status": "PENDING",
      "bridge": "hyperlane",
      "time_remaining": "600"
    },
    {
      "transaction_id": "hyperlane-tx-456",
      "type": "WITHDRAWAL_FROM_POSITION",
      "amount": "50000",
      "current_value": "50000",
      "source": "arbitrum-compound-v3",
      "destination": "noble",
      "initiated_at": "2024-01-15T13:45:00Z",
      "expected_at": "2024-01-15T14:15:00Z",
      "status": "CONFIRMED",
      "bridge": "hyperlane",
      "time_remaining": "0"
    }
  ],
  "total_inflight": "150000",
  "pending_deployment": "200000",
  "pending_withdrawal_distribution": "75000",
  "by_type": {
    "deposits_to_position": "100000",
    "withdrawals_from_position": "50000",
    "rebalance_between_positions": "0",
    "pending_deployment": "200000",
    "pending_withdrawal_distribution": "75000"
  },
  "by_status": {
    "pending": "100000",
    "confirmed": "50000",
    "completed": "0"
  }
}
```

### Response

- `inflight_funds` — Array of all inflight fund transactions (all amounts in USDN).
- `total_inflight` — Total value of all funds in transit (USDN).
- `pending_deployment` — USDN awaiting deployment to remote positions.
- `pending_withdrawal_distribution` — USDN awaiting distribution to withdrawal queue.
- `by_type` — Breakdown by transaction type including pending states.
- `by_status` — Breakdown by current status.

## StaleInflightFunds

**Endpoint**: `/noble/dollar/vaults/v2/stale_inflight_funds`

Retrieves inflight funds that have exceeded their expected arrival time and may require intervention.

```json
{
  "stale_funds": [
    {
      "transaction_id": "hyperlane-tx-789",
      "amount": "75000",
      "type": "DEPOSIT_TO_POSITION",
      "initiated_at": "2024-01-14T10:00:00Z",
      "expected_at": "2024-01-14T10:30:00Z",
      "hours_overdue": "28.5",
      "bridge": "hyperlane",
      "last_known_status": "PENDING",
      "recommended_action": "investigate"
    }
  ],
  "total_stale_value": "75000",
  "total_stale_count": 1,
  "oldest_stale_hours": "28.5"
}
```

### Response

- `stale_funds` — Array of stale inflight fund entries (amounts in USDN).
- `total_stale_value` — Total value of stale funds in USDN.
- `total_stale_count` — Number of stale transactions.
- `oldest_stale_hours` — Hours since oldest stale transaction.

