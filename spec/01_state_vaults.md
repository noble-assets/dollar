# Vaults State

## Paused

The `Paused` field is a [`collections.Item`][item] that stores the current paused state (`vaults.PauseType`).

```go
const PausedKey = []byte("paused")
```

## Rewards

The `Rewards` field is a mapping ([`collections.Map`][map]) between indexes (`string`) and a `vaults.Reward` value.

```go
const RewardPrefix = []byte("reward/")
```

## Positions

The `Positions` field is a mapping ([`collections.Map`][map]) between the keys <address, vault, timestamp> (`[]byte`, `vaults.VaultType`, `int64`) and a `vaults.Position` value.

```go
const PositionPrefix = []byte("position/")
```

## TotalFlexiblePrincipal

The `TotalFlexiblePrincipal` field is a [`collections.Item`][item] that stores the current total principal stored in the flexible vault (`math.Int`).

```go
const TotalFlexiblePrincipalKey = []byte("paused")
```

[item]: https://docs.cosmos.network/v0.50/build/packages/collections#item
[map]: https://docs.cosmos.network/v0.50/build/packages/collections#map
