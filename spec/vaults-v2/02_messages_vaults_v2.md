# Vaults V2 Messages

## Deposit

`noble.dollar.vaults.v2.MsgDeposit`

This message allows Noble Dollar users to deposit $USDN into the Noble vault and receive shares representing their proportional ownership. The system implements multiple safeguards against malicious deposit behavior including velocity checks, cooldown periods, and per-block limits.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgDeposit",
        "signer": "noble1user",
        "amount": "1000000"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `amount` — The amount of $USDN to deposit into the Noble vault.

### Requirements

- The `amount` must meet the minimum deposit requirement.
- User must not exceed per-user deposit limits.
- Total deposits in the current block must not exceed per-block limits.
- User must have passed any required cooldown period since last deposit.
- Vault must not be in emergency mode.
- User's deposit velocity must not trigger suspicious activity flags.

- User's $USDN balance is decreased by the deposit amount.
- Funds are marked as PENDING_DEPLOYMENT until allocated to remote positions.
- Shares are minted to the user based on current NAV.
- User's deposit history and velocity metrics are updated.
- Total vault shares are increased.
- NAV includes pending deployment funds as local assets.

### Anti-Manipulation Mechanisms

- **Velocity Tracking**: Monitors frequency and volume of deposits over rolling time windows.
- **Cooldown Enforcement**: Prevents rapid successive deposits that could manipulate share prices.
- **Block Limits**: Caps total deposits per block to prevent flash loan attacks.
- **Share Price Protection**: Uses time-weighted average NAV for share calculations during high volatility.

## RequestWithdrawal

`noble.dollar.vaults.v2.MsgRequestWithdrawal`

This message initiates a withdrawal request that enters a queue for processing. The queued approach prevents manipulation through sandwich attacks and ensures fair NAV-based redemptions.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgRequestWithdrawal",
        "signer": "noble1user",
        "shares": "1000"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `shares` — The number of shares to redeem from the Noble vault.

### Requirements

- User must have sufficient share balance in the vault.
- Vault must not be in emergency withdrawal-only mode.
- Shares are immediately locked and cannot be transferred.

### State Changes

- User's shares are locked (not burned yet).
- A withdrawal request is created and added to the queue with PENDING status.
- Request ID is generated and returned to the user.
- NAV at time of request is recorded for fair value calculation.

## ClaimWithdrawal

`noble.dollar.vaults.v2.MsgClaimWithdrawal`

This message allows users to claim fulfilled withdrawal requests from the queue.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgClaimWithdrawal",
        "signer": "noble1user",
        "request_id": "42"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `request_id` — The unique identifier of the withdrawal request to claim.

### Requirements

- Request must belong to the signer.
- Request status must be CLAIMABLE.
- Sufficient time must have passed since request (withdrawal delay).

### State Changes

- $USDN is transferred to the user based on the fulfilled amount.
- User's shares are burned.
- Request status is updated to CLAIMED.
- Pending withdrawals counter is decreased.

## UpdateNAV

`noble.dollar.vaults.v2.MsgUpdateNAV`

This message updates the Net Asset Value of the Noble vault based on the performance of its remote positions. This is typically called by authorized keepers or through automated processes.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgUpdateNAV",
        "signer": "noble1keeper",
        "position_updates": [
          {
            "position_id": "1",
            "new_value": "1050000",
            "proof": "base64_encoded_proof"
          },
          {
            "position_id": "2",
            "new_value": "2100000",
            "proof": "base64_encoded_proof"
          }
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `position_updates` — Array of remote position value updates with proofs for the vault's positions.

### Requirements

- Signer must be an authorized keeper or the system itself.
- Proofs must be valid according to configured oracle requirements.
- Updates must not be stale based on max staleness configuration.
- All active remote positions should be included or have recent updates.

### State Changes

- Individual remote position values are updated.
- Inflight funds are included in NAV calculation with their last known values.
- Vault's total NAV is recalculated as: Local Assets + Σ(Remote Positions) + Σ(Inflight Funds) - Pending Liabilities.
- Last NAV update timestamp is recorded.
- Pending withdrawal requests may be processed if liquidity available.
- Stale inflight funds (beyond max duration) may be marked for investigation.

## CreateRemotePosition

`noble.dollar.vaults.v2.MsgCreateRemotePosition`

This message deploys capital from the Noble vault to a remote yield-generating position on another chain or protocol. The vault can maintain multiple remote positions for diversification.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgCreateRemotePosition",
        "signer": "noble1manager",
        "protocol": "aave-v3",
        "chain_id": 2,
        "amount": "1000000",
        "asset_address": "base64_encoded_asset_address",
        "parameters": {
          "pool_address": "base64_encoded_pool_address",
          "strategy": "stable_lending"
        }
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `protocol` — The target protocol identifier.
- `chain_id` — The Wormhole Chain ID of the destination.
- `amount` — The amount of capital to deploy.
- `asset_address` — The address of the asset on the remote chain.
- `parameters` — Protocol-specific deployment parameters.

### Requirements

- Signer must be the vault manager or authorized operator.
- The vault must have sufficient unallocated capital.
- Protocol and chain must be in the vault's allowed list.
- Number of positions must not exceed the vault's max remote positions limit.
- Combined remote position exposure must maintain diversification requirements.

### State Changes

- Capital is marked as inflight with DEPOSIT_TO_POSITION type.
- Inflight fund entry created with expected arrival time (amount in USDN).
- Capital is bridged to the destination chain via Hyperlane/IBC.
- Vault's available liquidity is decreased.
- NAV continues to include the inflight USDN value during transit.
- Upon confirmation, inflight status transitions to CONFIRMED.
- New remote position entry is created with ACTIVE status once funds arrive.
- Inflight fund entry is marked COMPLETED and archived.

## CloseRemotePosition

`noble.dollar.vaults.v2.MsgCloseRemotePosition`

This message initiates the withdrawal of capital from a remote position back to the vault.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgCloseRemotePosition",
        "signer": "noble1manager",
        "position_id": "1",
        "partial_amount": "500000"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `position_id` — The unique identifier of the remote position to close.
- `partial_amount` — Optional amount for partial withdrawal (omit for full closure).

### Requirements

- Signer must be the vault manager or authorized operator.
- Position must be in ACTIVE status.
- For partial withdrawals, amount must leave minimum position size.

### State Changes

- Withdrawal is initiated on the remote chain.
- Position status updated to WITHDRAWING.
- Inflight fund entry created with WITHDRAWAL_FROM_POSITION type (amount in USDN).
- Expected return amount and arrival time are tracked.
- NAV continues to include the inflight USDN value during transit.
- Funds marked as PENDING_WITHDRAWAL_DISTRIBUTION upon arrival.
- Withdrawal requests in queue may be marked for processing upon receipt.
- Upon bridge confirmation via Hyperlane, inflight status transitions to CONFIRMED.
- Once received, inflight entry is marked COMPLETED.

## ProcessWithdrawalQueue

`noble.dollar.vaults.v2.MsgProcessWithdrawalQueue`

This message processes pending withdrawal requests in the queue when liquidity becomes available.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgProcessWithdrawalQueue",
        "signer": "noble1keeper",
        "max_requests": 100
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `max_requests` — Maximum number of withdrawal requests to process in this transaction.

### Requirements

- The vault must have available liquidity.
- Requests are processed in FIFO order.
- Each request uses the NAV at time of request for fair value.

### State Changes

- Pending requests are marked as CLAIMABLE up to available liquidity.
- Fulfilled amounts are calculated based on NAV.
- Available vault liquidity is reserved for claims.
- Users are notified their withdrawals are ready to claim.



## HandleStaleInflight

`noble.dollar.vaults.v2.MsgHandleStaleInflight`

This message handles inflight funds that have exceeded their maximum duration without reconciliation, potentially indicating bridge failures or other issues.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgHandleStaleInflight",
        "signer": "noble1governance",
        "transaction_id": "hyperlane-tx-456",
        "action": "write_off",
        "justification": "Bridge timeout confirmed, funds unrecoverable"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `transaction_id` — The identifier of the stale inflight transaction.
- `action` — Action to take: "write_off", "extend", "manual_recovery".
- `justification` — Explanation for the action taken.

### Requirements

- Transaction must have exceeded maximum inflight duration.
- Signer must be governance or emergency multisig.
- Action must be appropriate for the situation.
- Write-offs require governance approval for amounts above threshold.

### State Changes

- Inflight fund status updated based on action.
- For write-offs: NAV is reduced by the lost USDN amount.
- For extensions: New expected arrival time is set.
- For manual recovery: Fund is marked for manual intervention.
- Incident is logged for tracking purposes.
- Affected users may be compensated from insurance fund.

## Rebalance

`noble.dollar.vaults.v2.MsgRebalance`

This message rebalances capital across the Noble vault's multiple remote positions based on target allocations.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgRebalance",
        "signer": "noble1manager",
        "target_allocations": [
          {
            "position_id": "1",
            "target_percentage": "40"
          },
          {
            "position_id": "2",
            "target_percentage": "35"
          },
          {
            "position_id": "3",
            "target_percentage": "25"
          }
        ],
        "rebalance_strategy": "GRADUAL"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `target_allocations` — Desired allocation percentages across the vault's remote positions.

### Requirements

- Signer must be the vault manager.
- Target allocations must sum to 100% or less (remainder stays local).
- Rebalancing must respect minimum position sizes.
- Current allocations must deviate beyond rebalance threshold.

### State Changes

- Capital movements are initiated between positions (all in USDN).
- For position-to-position rebalancing:
  - Source position withdrawal creates WITHDRAWAL_FROM_POSITION inflight.
  - Upon arrival at Noble, funds marked as PENDING_DEPLOYMENT.
  - Deployment to target position creates DEPOSIT_TO_POSITION inflight.
- For position-to-Noble rebalancing (for withdrawals):
  - Creates WITHDRAWAL_FROM_POSITION inflight.
  - Upon arrival, funds marked as PENDING_WITHDRAWAL_DISTRIBUTION.
- Temporary liquidity constraints may delay withdrawal processing.
- Rebalancing transactions are tracked.
- NAV remains constant during rebalancing (inflight USDN counted at initiation value).
