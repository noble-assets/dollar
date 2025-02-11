# State

## Index

The `Index` field is a [`collections.Item`][item] that stores the current $USDN rebasing multiplier (`math.LegacyDec`).

```go
const IndexKey = []byte("index")
```

## Paused

The `Paused` field is a [`collections.Item`][item] that stores the current paused state (`boolean`).

```go
const PausedKey = []byte("paused")
```

## Principal

The `Principal` field is a mapping ([`collections.Map`][map]) between user addresses (`[]byte`) and their principal amount (`math.Int`) of their $USDN balance.

```go
const PrincipalPrefix = []byte("principal/")
```

[item]: https://docs.cosmos.network/v0.50/build/packages/collections#item
[map]: https://docs.cosmos.network/v0.50/build/packages/collections#map
