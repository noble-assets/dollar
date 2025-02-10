# Portal Messages

## Deliver

`noble.dollar.portal.v1.MsgDeliver`

This message acts as the core mechanism for delivering Noble Dollar Portal messages from other chains. While its primary use case is within vote extension processes by validators, it remains publicly accessible to support permissionless, manual message relaying.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.portal.v1.MsgDeliver",
        "signer": "noble1validator",
        "vaa": "base64_encoded_vaa"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `vaa` — The VAA (Verifiable Action Approval) message to be delivered.

### Requirements

- The VAA is a valid Noble Dollar Portal message, signed by the current Wormhole Guardian Set.

### State Changes

- An issuance of $USDN via the `x/bank` module.
  - In the case of an $M transfer, this is minted directly to a user.
  - In the case of an index update, this is minted to the module yield accrual account.

## Transfer

`noble.dollar.portal.v1.MsgTransfer`

This message allows transferring \$USDN from Noble and receiving \$M cross-chain.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.portal.v1.MsgTransfer",
        "signer": "noble1user",
        "chain": "2",
        "recipient": "base64_encoded_recipient_address",
        "amount": "1000000"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `chain` — The Wormhole Chain ID of the destination chain.
- `recipient` — The 32-byte encoded recipient address on the destination chain.
- `amount` — The amount of $USDN to be transferred.

### State Changes

- The specified amount of $USDN is burned on Noble and sent via Wormhole to the destination chain.

## SetPause

`noble.dollar.portal.v1.MsgSetPause`

This message allows the owner to set the PortalPause state to `true` or `false`, enabling or disabling the [Deliver](#deliver) and [Transfer](#transfer) actions.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.portal.v1.MsgSetPause",
        "signer": "noble1signer",
        "paused": "true"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `paused` —  Specifies the pause state to set (`true` | `false`).

### Requirements

- Signer must be the current [Owner](01_state_portal.md#owner).

### State Changes

- [`paused`](./01_state_portal.md#paused)

## Set Peer

`noble.dollar.portal.v1.MsgSetPeer`

This message allows the owner of the Noble Dollar Portal to set external peers.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.portal.v1.MsgSetPeer",
        "signer": "noble1owner",
        "chain": 1,
        "transceiver": "base64_encoded_transceiver",
        "manager": "base64_encoded_manager"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `chain` — The Wormhole Chain ID of the peer being set.
- `transceiver`— The transceiver address for cross-chain communication.
- `manager` — The manager address responsible for the peer's operation.

### Requirements

- Signer must be the current [`owner`](./01_state_portal.md#owner).

### State Changes

- [`peers`](./01_state_portal.md#peers)


## Transfer Ownership

`noble.dollar.portal.v1.MsgTransferOwnership`

This message allows the transfer of ownership of the Noble Dollar Portal to a new address.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.portal.v1.MsgTransferOwnership",
        "signer": "noble1owner",
        "new_owner": "noble1newowner"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `new_owner` — The Noble address to transfer ownership to.

### Requirements

- Signer must be the current [`owner`](./01_state_portal.md#owner).
- `new_owner` must not be the current [`owner`](./01_state_portal.md#owner).

### State Changes

- [`owner`](./01_state_portal.md#owner)
