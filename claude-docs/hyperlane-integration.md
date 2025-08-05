# Hyperlane Integration Documentation

## Overview

This document describes the integration of Hyperlane cross-chain messaging protocol with the Dollar Noble vault system using the `github.com/bcp-innovations/hyperlane-cosmos` library.

## Architecture

The Hyperlane integration consists of several key components:

### Core Components

1. **HyperlaneProvider** - Implements the `CrossChainProvider` interface for Hyperlane operations
2. **Core Keeper** - Manages mailboxes and message dispatching
3. **Warp Keeper** - Handles token bridging and transfers
4. **Util Package** - Provides address encoding, message parsing, and other utilities

### Key Dependencies

- `github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper` - Core Hyperlane functionality
- `github.com/bcp-innovations/hyperlane-cosmos/x/warp/keeper` - Token bridging capabilities
- `github.com/bcp-innovations/hyperlane-cosmos/util` - Utility functions and types

## API Reference

### HyperlaneProvider

The `HyperlaneProvider` struct implements cross-chain operations for Hyperlane:

```go
type HyperlaneProvider struct {
    coreKeeper       *corekeeper.Keeper
    warpKeeper       *warpkeeper.Keeper
    localDomain      uint32
    defaultGasLimit  uint64
    defaultTimeout   time.Duration
    confirmationsMap map[uint32]uint64
    mailboxId        util.HexAddress
}
```

#### Constructor

```go
func NewHyperlaneProvider(
    coreKeeper *corekeeper.Keeper,
    warpKeeper *warpkeeper.Keeper,
    localDomain uint32,
    defaultGasLimit uint64,
    defaultTimeout time.Duration,
    mailboxId util.HexAddress,
) *HyperlaneProvider
```

#### Core Methods

##### SendMessage
Dispatches a cross-chain message via Hyperlane:

```go
func (p *HyperlaneProvider) SendMessage(
    ctx sdk.Context, 
    route *vaultsv2.CrossChainRoute, 
    msg CrossChainMessage,
) (*vaultsv2.ProviderTrackingInfo, error)
```

**Parameters:**
- `ctx`: SDK context
- `route`: Cross-chain route configuration containing Hyperlane config
- `msg`: Message containing operation type, sender, recipient, amount, and data

**Returns:**
- `ProviderTrackingInfo`: Tracking information with message ID and status
- `error`: Error if dispatch fails

##### GetMessageStatus
Checks the delivery status of a dispatched message:

```go
func (p *HyperlaneProvider) GetMessageStatus(
    ctx sdk.Context, 
    tracking *vaultsv2.ProviderTrackingInfo,
) (MessageStatus, error)
```

**Returns:**
- `MessageStatusSent`: Message dispatched but not yet delivered
- `MessageStatusConfirmed`: Message delivered to destination
- `MessageStatusFailed`: Message failed to deliver

##### EstimateGas
Calculates gas requirements and costs:

```go
func (p *HyperlaneProvider) EstimateGas(
    ctx sdk.Context, 
    route *vaultsv2.CrossChainRoute, 
    msg CrossChainMessage,
) (uint64, math.Int, error)
```

**Returns:**
- `uint64`: Estimated gas units
- `math.Int`: Total cost in tokens
- `error`: Error if estimation fails

#### Token Operations

##### CreateCollateralToken
Creates a new collateral token for cross-chain bridging:

```go
func (p *HyperlaneProvider) CreateCollateralToken(
    ctx sdk.Context, 
    owner string, 
    originMailbox util.HexAddress, 
    denom string,
) (util.HexAddress, error)
```

##### TransferRemoteCollateral
Initiates a collateral token transfer to a remote chain:

```go
func (p *HyperlaneProvider) TransferRemoteCollateral(
    ctx sdk.Context,
    tokenId util.HexAddress,
    cosmosSender string,
    destinationDomain uint32,
    recipient util.HexAddress,
    amount math.Int,
    customHookId *util.HexAddress,
    gasLimit math.Int,
    maxFee sdk.Coin,
    customHookMetadata []byte,
) (util.HexAddress, error)
```

## Configuration

### Route Configuration

Cross-chain routes must include Hyperlane-specific configuration:

```go
route := &vaultsv2.CrossChainRoute{
    RouteId: "hyperlane-ethereum",
    ProviderConfig: &vaultsv2.CrossChainProviderConfig{
        Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
            HyperlaneConfig: &vaultsv2.HyperlaneConfig{
                DomainId:              1, // Ethereum mainnet
                MailboxAddress:        "0x...", // Hyperlane mailbox contract
                GasPaymasterAddress:   "0x...", // Optional gas paymaster
                HookAddress:           "0x...", // Optional post-dispatch hook
                GasLimit:              200000,
                GasPrice:              math.NewInt(20000000000), // 20 gwei
            },
        },
    },
}
```

### Required Fields

- `DomainId`: Hyperlane domain identifier for the destination chain
- `MailboxAddress`: Address of the Hyperlane mailbox contract

### Optional Fields

- `GasPaymasterAddress`: Contract address for gas payment automation
- `HookAddress`: Custom post-dispatch hook contract
- `GasLimit`: Gas limit for destination chain execution
- `GasPrice`: Gas price for cost estimation

## Address Encoding

Hyperlane uses a specific address encoding format via `util.HexAddress`:

```go
// Decode from string
hexAddr, err := util.DecodeHexAddress("0x1234567890123456789012345678901234567890")

// Encode to string
addrString := hexAddr.String()

// Get internal ID for storage
internalId := hexAddr.GetInternalId()
```

## Message Tracking

### Tracking Structure

```go
type HyperlaneTrackingInfo struct {
    MessageId              []byte  // Hyperlane message ID
    OriginDomain           uint32  // Source domain
    DestinationDomain      uint32  // Target domain
    Nonce                  uint64  // Message nonce
    OriginTxHash           string  // Source transaction hash
    DestinationTxHash      string  // Destination transaction hash
    OriginBlockNumber      uint64  // Source block height
    DestinationBlockNumber uint64  // Destination block height
    Processed              bool    // Processing status
    GasUsed                uint64  // Gas consumed
}
```

### Status Updates

Message status can be updated by relayers:

```go
err := provider.UpdateMessageStatus(
    ctx,
    tracking,
    "0xdestinationTxHash",
    destinationBlockNumber,
    gasUsed,
    true, // processed
)
```

## Example Usage

### Basic Setup

```go
import (
    "github.com/bcp-innovations/hyperlane-cosmos/util"
    corekeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
    warpkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/warp/keeper"
    "dollar.noble.xyz/v2/keeper/crosschain"
)

// Create provider
mailboxId, _ := util.DecodeHexAddress("0x1234567890123456789012345678901234567890")
provider := crosschain.NewHyperlaneProvider(
    coreKeeper,  // injected dependency
    warpKeeper,  // injected dependency
    4,           // Noble domain ID
    200000,      // default gas limit
    time.Hour,   // default timeout
    mailboxId,   // mailbox ID
)

// Register with cross-chain keeper
keeper.RegisterProvider(provider)
```

### Sending Messages

```go
// Create cross-chain message
msg := crosschain.CrossChainMessage{
    Type:      crosschain.MessageTypeDeposit,
    Sender:    senderAddr,
    Recipient: "0xrecipientAddress",
    Amount:    math.NewInt(1000000), // 1 USDC
    Data:      []byte("additional data"),
}

// Send via route
tracking, err := provider.SendMessage(ctx, route, msg)
if err != nil {
    return err
}

// Check status
status, err := provider.GetMessageStatus(ctx, tracking)
```

### Gas Estimation

```go
gasUnits, totalCost, err := provider.EstimateGas(ctx, route, msg)
if err != nil {
    return err
}

fmt.Printf("Estimated gas: %d units, cost: %s\n", gasUnits, totalCost.String())
```

## Error Handling

Common error scenarios and handling:

### Address Validation Errors
```go
if err := provider.ValidateConfig(config); err != nil {
    // Handle invalid addresses or configuration
    log.Error("Invalid Hyperlane configuration", "error", err)
}
```

### Message Dispatch Errors
```go
tracking, err := provider.SendMessage(ctx, route, msg)
if err != nil {
    // Handle insufficient fees, invalid recipients, etc.
    log.Error("Failed to dispatch message", "error", err)
}
```

### Status Query Errors
```go
status, err := provider.GetMessageStatus(ctx, tracking)
if err != nil {
    // Handle message not found, invalid tracking info, etc.
    log.Error("Failed to query message status", "error", err)
}
```

## Testing

### Unit Tests

The provider includes comprehensive unit tests covering:
- Message dispatch and tracking
- Gas estimation
- Configuration validation
- Error handling

### Integration Tests

End-to-end tests demonstrate:
- Cross-chain vault operations
- Yield distribution via Hyperlane
- Token bridging workflows

### Example Test

```go
func TestHyperlaneIntegration(t *testing.T) {
    // Setup test environment
    ctx, chain, validator := setupTestChain(t)
    
    // Create and fund test user
    user := createTestUser(t, ctx, chain)
    
    // Test message dispatch
    tracking, err := provider.SendMessage(ctx, route, msg)
    require.NoError(t, err)
    require.NotNil(t, tracking)
    
    // Verify message status
    status, err := provider.GetMessageStatus(ctx, tracking)
    require.NoError(t, err)
    require.Equal(t, crosschain.MessageStatusSent, status)
}
```

## Security Considerations

### Address Validation
- All addresses must be validated using `util.DecodeHexAddress`
- Invalid addresses will cause transaction failures

### Gas Limits
- Set appropriate gas limits to prevent failed executions
- Monitor gas price fluctuations on destination chains

### Fee Management
- Implement maximum fee limits to prevent excessive costs
- Use gas paymasters for automated fee handling

### Message Verification
- Hyperlane includes built-in security modules (ISMs)
- Configure appropriate security thresholds for your use case

## Troubleshooting

### Common Issues

1. **Invalid Address Format**
   - Ensure addresses are properly hex-encoded
   - Use `util.DecodeHexAddress` for validation

2. **Insufficient Gas**
   - Increase gas limits in route configuration
   - Monitor destination chain gas prices

3. **Message Not Delivered**
   - Check relayer status
   - Verify ISM configuration
   - Confirm sufficient gas payment

4. **Configuration Errors**
   - Validate all required fields are set
   - Test configuration with `ValidateConfig`

### Debug Tools

```go
// Get domain information
domainInfo, err := provider.GetDomainInfo(ctx, domainId)
if err == nil {
    fmt.Printf("Domain %d: Block %d, Gas Price: %s\n", 
        domainInfo.DomainID, 
        domainInfo.LatestBlockNumber, 
        domainInfo.CurrentGasPrice.String())
}

// Check required confirmations
confirmations := provider.GetRequiredConfirmations(domainId)
fmt.Printf("Required confirmations for domain %d: %d\n", domainId, confirmations)
```

## Future Enhancements

### Planned Features
- Automatic gas price optimization
- Enhanced error recovery mechanisms
- Multi-route message aggregation
- Advanced security module configurations

### Migration Path
- The current implementation uses release candidate `v1.0.0-rc0`
- Monitor for stable releases and upgrade paths
- Test thoroughly before production deployments

## References

- [Hyperlane Documentation](https://docs.hyperlane.xyz/)
- [BCP Innovations Hyperlane Cosmos](https://github.com/bcp-innovations/hyperlane-cosmos)
- [Dollar Noble Repository](https://github.com/noble-assets/dollar)