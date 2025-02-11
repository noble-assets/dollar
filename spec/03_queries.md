# Queries

## Paused

**Endpoint**: `/noble/dollar/v1/paused`

Retrieves the current paused state of the Noble Dollar.

```json
{
  "paused": "true"
}
```

### Response

- `paused` — The paused state of the Noble Dollar.

## Principal

**Endpoint**: `/noble/dollar/v1/principal/{account}`

Retrieves the principal amount associated with a $USDN holders account.

```json
{
  "principal": "1000000"
}
```

### Arguments

- `account` — The address of the holder you wish to request the principal of.

### Response

- `principal` — The current principal amount held by the requested account.

## Stats

**Endpoint**: `/noble/dollar/v1/stats`

Retrieves the latest stats of Noble Dollar.

```json
{
  "total_holders": "1000",
  "total_principal": "100000",
  "total_yield_accrued": "100"
}
```

### Response

- `total_holders`:  — The total number of $USDN holders.
- `total_principal`:  — The total principal amount in the system.
- `total_yield_accrued`:  — The total amount of yield that has been accrued over time.


## Yield

**Endpoint**: `/noble/dollar/v1/yield/{account}`

Retrieves the amount of yield that is claimable for a $USDN holder.

```json
{
  "claimable_amount": "50000"
}
```

### Arguments

- `account` — The address of the holder you wish to request the yield of.

### Response

- `claimable_amount` — The current amount of yield claimable by the requested account.
