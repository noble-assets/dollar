# NAV Oracle System

## Overview

The NAV Oracle System provides real-time valuation data for remote positions held by vaults on Noble. Each remote position has a dedicated oracle that pushes share price and shares held data via Hyperlane messages.

## Architecture

### Oracle-to-Position Mapping

Each remote position is uniquely identified and mapped to a specific oracle:

```go
type RemotePosition struct {
    PositionID    string    // Unique identifier for the position
    ChainID       uint32    // Chain where position exists
    OracleAddress string    // Address of the oracle contract on remote chain
    VaultAddress  string    // Address of the vault holding this position
    SharesHeld    sdk.Int   // Number of shares held in this position
    SharePrice    sdk.Dec   // Current share price
    LastUpdate    time.Time // Timestamp of last oracle update
}

type OracleMapping struct {
    OracleAddress string // Oracle address on remote chain
    ChainID       uint32 // Source chain ID
    PositionID    string // Position this oracle reports for
}
```

### Data Flow

1. **Remote Oracle Contract** monitors the position and calculates share price
2. **Oracle** pushes update message via Hyperlane when price changes or on schedule
3. **Relayer** submits `MsgProcessMessage` to Noble with the oracle update
4. **Vault Module** decodes the `Message` field from `MsgProcessMessage`
5. **Vault Module** processes the NAV update and recalculates vault NAV

## Hyperlane Message Format

### MsgProcessMessage Structure

Oracle updates arrive through Hyperlane's `MsgProcessMessage`:

```go
type MsgProcessMessage struct {
    MailboxId HexAddress `json:"mailbox_id"` // Hyperlane mailbox identifier
    Relayer   string     `json:"relayer"`    // Address of relayer submitting
    Metadata  string     `json:"metadata"`   // Hyperlane metadata for verification
    Message   string     `json:"message"`    // Base64-encoded NAV oracle update
}
```

###  Position Update 

The `Message` field contains the NAV oracle update. All oracle updates are for individual positions. Each decoded message contains:

```
Message Structure (Fixed-Length: 97 bytes):
+----------------+------------------+-------------------+------------------+------------------+
| Field          | Size (bytes)     | Offset            | Type             | Description      |
+----------------+------------------+-------------------+------------------+------------------+
| Message Type   | 1                | 0                 | uint8            | 0x01 (NAV update)|
| Position ID    | 32               | 1                 | bytes32          | Unique position  |
| Share Price    | 32               | 33                | uint256          | Price per share  |
| Shares Held    | 32               | 65                | uint256          | Total shares     |
| Timestamp      | 8                | 97                | uint64           | Update timestamp |
+----------------+------------------+-------------------+------------------+------------------+
Total: 105 bytes
```

### Encoding Example

```go
func EncodeOracleUpdate(positionID [32]byte, sharePrice, sharesHeld *big.Int, timestamp uint64) []byte {
    buf := make([]byte, 105)
    
    // Message type
    buf[0] = 0x01
    
    // Position ID
    copy(buf[1:33], positionID[:])
    
    // Share price (32 bytes, big-endian)
    priceBytes := sharePrice.FillBytes(make([]byte, 32))
    copy(buf[33:65], priceBytes)
    
    // Shares held (32 bytes, big-endian)
    sharesBytes := sharesHeld.FillBytes(make([]byte, 32))
    copy(buf[65:97], sharesBytes)
    
    // Timestamp (8 bytes, big-endian)
    binary.BigEndian.PutUint64(buf[97:105], timestamp)
    
    return buf
}
```

## Update Mechanisms

### Push-Based Architecture

Oracles actively push updates to Noble - there is no pull mechanism:

1. **Scheduled Updates**: Regular intervals (e.g., every hour)
2. **Threshold Updates**: When price changes exceed configured threshold
3. **Event-Driven Updates**: On significant position events (deposits, withdrawals)

### Oracle Registration
### Registration

Each oracle must be registered and enrolled before sending updates:

```go
// PositionOracleConfig stores oracle configuration for a position
type PositionOracleConfig struct {
    PositionID            string                // Unique position identifier
    OriginMailbox         util.HexAddress       // Expected origin mailbox
    MaxStaleness          time.Duration         // Maximum age before data is stale
}

// EnrolledOracleRouter maps position+origin to authorized oracle contract
type EnrolledOracleRouter struct {
    PositionID      string // Position this oracle reports for
    OriginDomain    uint32 // Source chain domain
    OracleContract  string // Authorized oracle contract address (hex)
}
```

## Oracle-Position Association

### Ensuring Correct Mapping

The system validates oracle-to-position mappings through enrolled routers:

```go
// ValidateOracleMessage ensures the message comes from an authorized oracle
func (k *Keeper) ValidateOracleMessage(
    ctx context.Context,
    mailboxId util.HexAddress,
    message util.HyperlaneMessage,
    positionID string,
) error {
    // Get oracle configuration
    oracle, err := k.RemotePositionOracles.Get(ctx, positionID)
    if err != nil {
        return fmt.Errorf("oracle not registered for position %s", positionID)
    }
    
    // Verify origin mailbox
    if oracle.OriginMailbox != mailboxId {
        return fmt.Errorf("invalid mailbox: expected %s, got %s", oracle.OriginMailbox, mailboxId)
    }
    
    // Get enrolled router for this origin
    enrolledOracle, err := k.EnrolledOracleRouters.Get(ctx, collections.Join(positionID, message.Origin))
    if err != nil {
        return fmt.Errorf("no enrolled oracle for origin %d", message.Origin)
    }
    
    // Verify sender is authorized
    if message.Sender.String() != strings.ToLower(enrolledOracle.OracleContract) {
        return fmt.Errorf("unauthorized sender: %s", message.Sender)
    }
    
    return nil
}
```

### Position State Updates

When an oracle update is received and validated:

```go
func (k Keeper) ApplyOracleUpdate(
    ctx sdk.Context,
    positionID string,
    sharePrice sdk.Dec,
    sharesHeld sdk.Int,
    timestamp time.Time,
) error {
    // Get the position
    position, found := k.GetRemotePosition(ctx, positionID)
    if !found {
        return ErrPositionNotFound
    }
    
    // Update position data
    position.SharePrice = sharePrice
    position.SharesHeld = sharesHeld
    position.LastUpdate = timestamp
    
    // Calculate position value
    positionValue := sharePrice.Mul(sdk.NewDecFromInt(sharesHeld))
    position.TotalValue = positionValue
    
    // Store updated position
    k.SetRemotePosition(ctx, position)
    
    // Trigger NAV recalculation for associated vault
    return k.RecalculateVaultNAV(ctx, position.VaultAddress)
}
```

## NAV Calculation Integration

### Vault NAV Aggregation

The vault's total NAV is calculated by summing all remote position values:

```go
func (k Keeper) CalculateVaultNAV(ctx sdk.Context, vaultAddr string) (sdk.Dec, error) {
    positions := k.GetVaultPositions(ctx, vaultAddr)
    totalNAV := sdk.ZeroDec()
    
    for _, position := range positions {
        // Check staleness
        if time.Since(position.LastUpdate) > position.MaxStaleness {
            // Use last known value with staleness flag
            ctx.Logger().Warn("Stale position data", 
                "position", position.PositionID,
                "last_update", position.LastUpdate)
        }
        
        // Calculate position value: shares * price
        positionValue := position.SharePrice.Mul(sdk.NewDecFromInt(position.SharesHeld))
        totalNAV = totalNAV.Add(positionValue)
    }
    
    return totalNAV, nil
}
```

## Validation and Security

### Hyperlane ISM Validation

Hyperlane's Interchain Security Module (ISM) handles all message validation including:
- Message authenticity verification
- Source chain validation
- Sender verification
- Replay protection

The Noble-side handler only needs to verify the oracle is enrolled for the specific position:

```go
// Simple enrollment check - ISM handles cryptographic validation
type EnrolledOracle struct {
    PositionID    string   // Position this oracle reports for
    OriginDomain  uint32   // Expected source domain
    OracleAddress string   // Expected sender address
}
```

### Staleness Protection

```go
type StatenessConfig struct {
    WarningThreshold    time.Duration // 1 hour - log warning
    CriticalThreshold   time.Duration // 4 hours - alert operators
}

func (k Keeper) CheckDataFreshness(ctx sdk.Context, positionID string) DataStatus {
    position := k.GetRemotePosition(ctx, positionID)
    age := time.Since(position.LastUpdate)
    
    config := k.GetStatenessConfig(ctx)
    
    switch {
    case age > config.CriticalThreshold:
        return DataStatusCritical
    case age > config.WarningThreshold:
        return DataStatusWarning
    default:
        return DataStatusFresh
    }
}
```

## Hyperlane Message Processing

### Handler Implementation

The vault module implements a handler that processes NAV oracle updates from Hyperlane messages:

```go
// Handle processes NAV oracle updates from Hyperlane messages
func (k *Keeper) Handle(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error {
    // Parse NAV payload from message body
    payload, err := types.ParseNAVPayload(message.Body)
    if err != nil {
        return fmt.Errorf("failed to parse NAV payload: %w", err)
    }
    
    // Get the oracle configuration for this position
    oracle, err := k.RemotePositionOracles.Get(ctx, payload.PositionID)
    if err != nil {
        return fmt.Errorf("no oracle registered for position %s: %w", payload.PositionID, err)
    }
    
    // Verify the message comes from the correct mailbox
    if oracle.OriginMailbox != mailboxId {
        return fmt.Errorf("invalid origin mailbox: expected %s, got %s", oracle.OriginMailbox, mailboxId)
    }
    
    // Get enrolled oracle router for this origin
    enrolledOracle, err := k.EnrolledOracleRouters.Get(ctx, collections.Join(payload.PositionID, message.Origin))
    if err != nil {
        return fmt.Errorf("no enrolled oracle found for origin %d and position %s", message.Origin, payload.PositionID)
    }
    
    // Verify sender is the authorized oracle contract
    if message.Sender.String() != strings.ToLower(enrolledOracle.OracleContract) {
        return fmt.Errorf("unauthorized oracle: expected %s, got %s", enrolledOracle.OracleContract, message.Sender.String())
    }
    
    // Validate data freshness
    // Apply the NAV update
    err = k.ApplyNAVUpdate(ctx, oracle, payload)
    if err != nil {
        return fmt.Errorf("failed to apply NAV update: %w", err)
    }
    
    // Emit event
    _ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventNAVOracleUpdate{
        Sender:       message.Sender,
        PositionId:   payload.PositionID,
        OriginDomain: message.Origin,
        SharePrice:   payload.SharePrice.String(),
        SharesHeld:   payload.SharesHeld.String(),
        Timestamp:    payload.Timestamp,
    })
    
    return nil
}

// ParseNAVPayload extracts NAV oracle data from Hyperlane message body
func ParseNAVPayload(body []byte) (*NAVPayload, error) {
    if len(body) != NAV_PAYLOAD_SIZE { // 105 bytes
        return nil, fmt.Errorf("invalid NAV payload size: expected %d, got %d", NAV_PAYLOAD_SIZE, len(body))
    }
    
    // Check message type
    if body[0] != NAV_UPDATE_MESSAGE_TYPE { // 0x01
        return nil, fmt.Errorf("invalid message type: 0x%02x", body[0])
    }
    
    // Extract position ID (32 bytes)
    positionID := hex.EncodeToString(body[1:33])
    
    // Extract share price (32 bytes, big-endian)
    sharePriceBig := new(big.Int).SetBytes(body[33:65])
    sharePrice := sdk.NewDecFromBigInt(sharePriceBig).Quo(sdk.NewDec(1e18))
    
    // Extract shares held (32 bytes, big-endian)
    sharesHeldBig := new(big.Int).SetBytes(body[65:97])
    sharesHeld := sdk.NewIntFromBigInt(sharesHeldBig)
    
    // Extract timestamp (8 bytes, big-endian)
    timestamp := time.Unix(int64(binary.BigEndian.Uint64(body[97:105])), 0)
    
    return &NAVPayload{
        MessageType: body[0],
        PositionID:  positionID,
        SharePrice:  sharePrice,
        SharesHeld:  sharesHeld,
        Timestamp:   timestamp,
    }, nil
}

// NAVPayload represents the decoded NAV oracle update
type NAVPayload struct {
    MessageType uint8
    PositionID  string
    SharePrice  sdk.Dec
    SharesHeld  sdk.Int
    Timestamp   time.Time
}

// ApplyNAVUpdate updates the position state with new NAV data
func (k *Keeper) ApplyNAVUpdate(ctx context.Context, oracle types.PositionOracleConfig, payload *NAVPayload) error {
    // Get the remote position
    position, err := k.RemotePositions.Get(ctx, payload.PositionID)
    if err != nil {
        return fmt.Errorf("position not found: %w", err)
    }
    
    // Update position state
    position.SharePrice = payload.SharePrice
    position.SharesHeld = payload.SharesHeld
    position.LastUpdate = payload.Timestamp
    position.TotalValue = payload.SharePrice.Mul(sdk.NewDecFromInt(payload.SharesHeld))
    
    // Store updated position
    if err := k.RemotePositions.Set(ctx, payload.PositionID, position); err != nil {
        return err
    }
    
    // Update last oracle update time
    if err := k.LastOracleUpdate.Set(ctx, payload.PositionID, payload.Timestamp); err != nil {
        return err
    }
    
    // Trigger NAV recalculation for the vault
    return k.RecalculateVaultNAV(ctx, position.VaultAddress)
}
```

## Error Handling

### Oracle Failure Modes

1. **Stale Data**: Use last known value with warning
2. **Missing Oracle**: Position marked as unavailable for NAV
3. **Invalid Data**: Reject update, maintain last good value
4. **Unauthorized Sender**: Reject and log security event

### Fallback Strategy

```go
type FallbackStrategy struct {
    UseLastKnownGood    bool          // Use cached value if fresh enough
    MaxCacheAge         time.Duration // Maximum age for cached data
    AlertThreshold      time.Duration // When to alert operators
}
```

## Governance Parameters

### Configurable Parameters

```go
type OracleGovernanceParams struct {
    // Update intervals
    MaxUpdateInterval    time.Duration // Maximum time between updates (freshness)
}
```

### Oracle Registration Governance

Adding or updating oracle registrations requires governance:

```go
type MsgRegisterOracle struct {
    Authority     string // Governance module account
    PositionID    string 
    OracleAddress string
    ChainID       uint32
    MaxStaleness  time.Duration
}
```

## Monitoring and Alerts

### Key Metrics

- **Update Frequency**: Track updates per position
- **Data Staleness**: Age of last update per position
- **Price Volatility**: Detect unusual price movements
- **Oracle Health**: Track failed updates and errors

### Alert Triggers

```go
type AlertConfig struct {
    StaleDataThreshold      time.Duration // Alert if data older than threshold
    PriceDeviationThreshold sdk.Dec       // Alert on large price changes
    FailedUpdateThreshold   uint32        // Alert after N failed updates
}
```

## Security Considerations

### Attack Vectors and Mitigations

1. **Oracle Manipulation**
   - Mitigation: Enrolled oracle verification, ISM validates message authenticity

2. **Replay Attacks**
   - Mitigation: Handled by Hyperlane ISM's nonce tracking

3. **Data Staleness Attacks**
   - Mitigation: Staleness checks, fallback to last known value

4. **Position ID Confusion**
   - Mitigation: Strict oracle-to-position enrollment validation

### Fallback Procedures

1. **Stop Oracle Updates**: Governance can disable oracle updates
2. **Oracle Replacement**: Swap to backup oracle via governance
3. **Position Removal**: Remove compromised position from NAV

## Conclusion

This oracle system provides secure, efficient NAV updates for vault positions by:
- Maintaining one oracle per remote position
- Using fixed-length message encoding for gas efficiency
- Enforcing strict oracle-to-position enrollment
- Leveraging Hyperlane ISM for message validation
- Implementing staleness checks and fallback mechanisms