# Queries

## Principal

**Endpoint**: `/dollar/v1/principal/{account}`

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

## Yield

**Endpoint**: `/dollar/v1/yield/{account}`

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
