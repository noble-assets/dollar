# Portal Queries

## Owner

**Endpoint**: `/noble/dollar/portal/v1/owner`

Retrieves the address of the current Noble Dollar Portal owner.

```json
{
  "address": "noble1owner"
}
```

### Response

- `owner` — The address of the current Noble Dollar Portal owner.

## Paused

**Endpoint**: `/noble/dollar/portal/v1/paused`

Retrieves the current paused state of the Noble Dollar Portal.

```json
{
  "paused": "true"
}
```

### Response

- `paused` — The paused state of the Noble Dollar Portal.

## Peers

**Endpoint**: `/noble/dollar/portal/v1/peers`

Retrieves all of the current Noble Dollar Portal external peers, filtered by Wormhole Chain ID.

```json
{
  "peers": {
    "1": {
      "transceiver": "",
      "manager": ""
    },
    "2": {...},
    ...
  }
}
```

### Response

- `peers` — A map containing external peer information, where the key is the Wormhole Chain ID (as `uint16`) and the value is a `portal.Peer` object containing the transceiver and manager addresses.

## Destination Tokens

**Endpoint**: `/noble/dollar/portal/v1/destination_tokens/{chain_id}`

Retrieves all supported destination tokens based off of the current Noble Dollar Portal bridging paths, filtered by Wormhole Chain ID.

```json
{
  "destination_tokens": ["...", "...", ...]
}
```

### Response

- `destination_tokens` — An array containing the list of supported destination tokens of the provided Wormhole Chain ID.

## Nonce

**Endpoint**: `/noble/dollar/portal/v1/nonce`

Retrieves the latest message sent nonce of the Noble Dollar Portal.

```json
{
  "nonce": 42
}
```

### Response

- `nonce` — A `uint64` of the latest message sent nonce.
