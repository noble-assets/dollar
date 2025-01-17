# Vaults Messages

## Lock

`noble.dollar.vaults.v1.MsgLock`

This message allows Noble Dollar users to lock a specified amount of $USDN into various types of vaults, each offering unique reward mechanisms.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v1.MsgLock",
        "signer": "noble1signer",
        "vault": "flexible",
        "amount": 1000000
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

- `amount` — The amount of $USDN to lock in the Vault.
- `vault` — The Vault Type where to lock $USDN (`Staked | Flexible`).

### State Changes

- The specified amount of $USDN is locked in the selected Vault account.
  - For a Staked Vault, the user does not earn any yield.
  - For a Flexible Vault, the user earns standard yield plus a boosted yield derived from the Staked Vault.

## Unlock

`noble.dollar.vaults.v1.MsgUnlock`

This message allows Noble Dollar users to unlock a specified amount of $USDN from a selected type of vault.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v1.MsgUnlock",
        "signer": "noble1signer",
        "vault": "flexible",
        "amount": 500000
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

- `amount` — The amount of $USDN to unlock from the Vault.
- `vault` — The type of Vault to unlock $USDN from (`Staked | Flexible`) .

### Requirements

- The user must have one or more active position in the specified `vault` with an `amount` equal to or greater than the specified value.

### State Changes

- The user closes all positions in the specified Vault up to the given amount, claiming any available rewards.

## SetPause

`noble.dollar.vaults.v1.MsgSetPause`

This message allows the authority to set the Vault Pause state to `LOCK` | `UNLOCK` | `ALL` | `NONE`, enabling or disabling the [Lock](#lock) and [Unlock](#unlock) actions. 

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.vaults.v1.MsgSetPause",
        "signer": "noble1signer",
        "paused": "LOCK"
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

- `paused` —  Specifies the pause state to set (`LOCK` | `UNLOCK` | `ALL` | `NONE`).

### Requirements

- Signer must be the current [`owner`](./01_state_portal.md#owner).

### State Changes

- [`paused`](./01_state_vaults.md#paused)