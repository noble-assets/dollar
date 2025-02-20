# Portal State

## Owner

The `Owner` field is a [`collections.Item`][item] that stores the address of the current owner of this submodule.

```go
const OwnerKey = []byte("portal/owner")
```

## Paused

The `Paused` field is a [`collections.Item`][item] that stores the current paused state (`boolean`).

```go
const PausedKey = []byte("portal/paused")
```

## Peers

The `Peers` field is a mapping ([`collections.Map`][map]) between Wormhole Chain IDs (`uint16`) and a `portal.Peer` value.

```go
const PeerPrefix = []byte("portal/peer/")
```

## Bridging Paths

The `BridgingPaths` field is a mapping ([`collections.Map`][map]) between a pair ([`collections.Pair`][pair]), Wormhole Chain ID + destination token, and a `bool` value.

```go
const BridgingPathPrefix = []byte("portal/bridging_path/")
```

## Nonce

The `Nonce` field is a [`collections.Item`][item] that stores the latest sent message nonce (`uint32`).

```go
const NonceKey = []byte("portal/nonce")
```

[item]: https://docs.cosmos.network/v0.50/build/packages/collections#item
[pair]: https://docs.cosmos.network/v0.50/build/packages/collections#composite-keys
[map]: https://docs.cosmos.network/v0.50/build/packages/collections#map
