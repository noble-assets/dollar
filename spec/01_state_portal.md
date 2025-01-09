# Portal State

## Owner

The `Owner` field is a [`collections.Item`][item] that stores the address of the current owner of this submodule.

```go
const OwnerKey = []byte("owner")
```

## Peers

The `Peers` field is a mapping ([`collections.Map`][map]) between Wormhole Chain IDs (`uint16`) and a `portal.Peer` values.

```go
const PeerPrefix = []byte("peer/")
```

[item]: https://docs.cosmos.network/v0.50/build/packages/collections#item
[map]: https://docs.cosmos.network/v0.50/build/packages/collections#map
