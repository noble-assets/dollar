# Vaults V2 Messages

## Deposit

`noble.dollar.vaults.v2.MsgDeposit`

This message allows Noble Dollar users to deposit $USDN into a vault and receive shares representing their proportional ownership. The system implements multiple safeguards against malicious deposit behavior including velocity checks, cooldown periods, and per-block limits.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgDeposit",
        "signer": "noble1user",
        "vault_id": "vault-001",
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

- `vault_id` — The unique identifier of the vault to deposit into.
- `amount` — The amount of $USDN to deposit.

### Requirements

- The `amount` must meet the minimum deposit requirement for the vault.
- User must not exceed per-user deposit limits.
- Total deposits in the current block must not exceed per-block limits.
- User must have passed any required cooldown period since last deposit.
- Vault must not be in emergency mode.
- User's deposit velocity must not trigger suspicious activity flags.

### State Changes

- User's $USDN balance is decreased by the deposit amount.
- Vault shares are minted to the user based on current NAV.
- User's deposit history and velocity metrics are updated.
- Total vault shares are increased.

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
        "vault_id": "vault-001",
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

- `vault_id` — The unique identifier of the vault to withdraw from.
- `shares` — The number of shares to redeem.

### Requirements

- User must have sufficient share balance in the specified vault.
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

This message updates the Net Asset Value of a vault based on the performance of its remote positions. This is typically called by authorized keepers or through automated processes.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgUpdateNAV",
        "signer": "noble1keeper",
        "vault_id": "vault-001",
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

- `vault_id` — The unique identifier of the vault to update.
- `position_updates` — Array of remote position value updates with proofs.

### Requirements

- Signer must be an authorized keeper or the system itself.
- Proofs must be valid according to configured oracle requirements.
- Updates must not be stale based on max staleness configuration.
- All active positions for the vault should be included or have recent updates.

### State Changes

- Individual remote position values are updated.
- Vault's total NAV is recalculated as sum of all positions plus local assets.
- Last NAV update timestamp is recorded.
- Pending withdrawal requests may be processed if liquidity available.

## CreateRemotePosition

`noble.dollar.vaults.v2.MsgCreateRemotePosition`

This message deploys capital from a vault to a remote yield-generating position on another chain or protocol. Vaults can maintain multiple remote positions for diversification.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgCreateRemotePosition",
        "signer": "noble1manager",
        "vault_id": "vault-001",
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

- `vault_id` — The vault providing the capital.
- `protocol` — The target protocol identifier.
- `chain_id` — The Wormhole Chain ID of the destination.
- `amount` — The amount of capital to deploy.
- `asset_address` — The address of the asset on the remote chain.
- `parameters` — Protocol-specific deployment parameters.

### Requirements

- Signer must be the vault manager or authorized operator.
- Vault must have sufficient unallocated capital.
- Protocol and chain must be in the vault's allowed list.
- Number of positions must not exceed vault's max remote positions limit.
- Combined remote position exposure must maintain diversification requirements.

### State Changes

- Capital is bridged to the destination chain via Hyperlane or IBC.
- New remote position entry is created with ACTIVE status.
- Vault's available liquidity is decreased.
- Position is added to NAV tracking with initial value.

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
        "vault_id": "vault-001",
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

- `vault_id` — The vault that owns the position.
- `position_id` — The unique identifier of the position to close.
- `partial_amount` — Optional amount for partial withdrawal (omit for full closure).

### Requirements

- Signer must be the vault manager or authorized operator.
- Position must be in ACTIVE status.
- For partial withdrawals, amount must leave minimum position size.

### State Changes

- Withdrawal is initiated on the remote chain.
- Position status updated to WITHDRAWING.
- Expected return amount is tracked for reconciliation.
- Withdrawal requests in queue may be marked for processing upon receipt.

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
        "vault_id": "vault-001",
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

- `vault_id` — The vault whose queue to process.
- `max_requests` — Maximum number of requests to process in this transaction.

### Requirements

- Vault must have available liquidity.
- Requests are processed in FIFO order.
- Each request uses the NAV at time of request for fair value.

### State Changes

- Pending requests are marked as CLAIMABLE up to available liquidity.
- Fulfilled amounts are calculated based on NAV.
- Available vault liquidity is reserved for claims.
- Users are notified their withdrawals are ready to claim.

## Rebalance

`noble.dollar.vaults.v2.MsgRebalance`

This message rebalances capital across a vault's multiple remote positions based on target allocations.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v2.MsgRebalance",
        "signer": "noble1manager",
        "vault_id": "vault-001",
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

- `vault_id` — The vault to rebalance.
- `target_allocations` — Desired allocation percentages across positions.

### Requirements

- Signer must be the vault manager.
- Target allocations must sum to 100% or less (remainder stays local).
- Rebalancing must respect minimum position sizes.
- Current allocations must deviate beyond rebalance threshold.

### State Changes

- Capital movements are initiated between positions.
- Temporary liquidity constraints may delay withdrawal processing.
- Rebalancing transactions are tracked for audit.
- NAV remains constant (minus any transaction costs).
