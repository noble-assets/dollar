# Vaults Queries

## Paused

**Endpoint**: `/dollar/vaults/v1/paused`

Retrieves the current Pause state of the Noble Dollar Vaults.

```json
{
  "paused": "lock"
}
```

### Response

- `paused` — The Pause state Noble Dollar Vaults module.

## PositionsByProvider

**Endpoint**: `/dollar/portal/v1/positions/{address}`

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
