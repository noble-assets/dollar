# Vaults Queries

## Paused

**Endpoint**: `/noble/dollar/vaults/v1/paused`

Retrieves the current Pause state of the Noble Dollar Vaults.

```json
{
  "paused": "lock"
}
```

### Response

- `paused` — The Pause state Noble Dollar Vaults module.

## PositionsByProvider

**Endpoint**: `/noble/dollar/portal/v1/positions/{address}`

Retrieves all the current Noble Dollar Vaults positions, filtered by a Provider address.

```json
{
  "positions" : [
    {
      "address": "F81sXJ7wwqXu0zUme2LTWSUXB1c=",
      "vault": "FLEXIBLE",
      "principal": "10000",
      "index": "1.000000000000000000",
      "amount": "10000",
      "time": "2020-01-01T08:00:00.000000Z"
    },
    {
      "address": "F81sXJ7wwqXu0zUme2LTWSUXB1c=",
      "vault": "STAKED",
      "principal": "10000",
      "index": "1.100000000000000000",
      "amount": "10000",
      "time": "2020-01-01T10:00:00.000000Z"
    }
  ]
}
```

### Response

- `positions` — An array of `vaults.PositionEntry` objects containing the user's position details.

## Stats

**Endpoint**: `/noble/dollar/vaults/v1/stats`

Retrieves the latest stats of the Vaults.

```json
{
    "flexible_total_principal": "100",
    "flexible_total_users": "10000",
    "flexible_total_distributed_rewards_principal": "42",
    "staked_total_principal": "1000",
    "staked_total_users": "50"
}
```

### Response

- `flexible_total_principal`: — The total principal amount currently held within the Flexible Vault.
- `flexible_total_users`: — The total number of users who have funds locked in the Flexible Vault.
- `flexible_total_distributed_rewards_principal`: — The total amount of boosted rewards principal, that has been distributed to users in the Flexible Vault.
- `staked_total_principal`: — The total principal amount that is locked inside the Staked Vault.
- `staked_total_users`: — The total number of users who have funds staked in the Staked Vault.
