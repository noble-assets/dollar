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

Retrieves the current Pause state of the Noble Dollar Portal module.

```json
{
  "paused": "true"
}
```

### Response

- `paused` — The Pause state Noble Dollar Portal.

## Peers

**Endpoint**: `/noble/dollar/portal/v1/peers`

Retrieves all of the current Noble Dollar Portal external peers, filtered by Wormhole Chain ID.

```json
{
  "peers" : {
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
