# Messages

## Claim Yield

`noble.dollar.v1.MsgClaimYield`

A message allowing holders of the Noble Dollar to claim their accumulated yield from the protocol. This yield is transferred from the module yield accrual account to the account of the transaction signer.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.dollar.v1.MsgClaimYield",
        "signer": "noble1user"
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

- `signer` — The address of the user claiming yield.

### Requirements

- Signer must be a holder of $USDN with unclaimed yield.

### State Changes

- A transfer of $USDN from the module yield accrual to the transaction signer accounts.
