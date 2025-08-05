# Hyperlane Integration for Noble Dollar Vaults

## Overview

The Noble Dollar vault system now supports cross-chain operations through both **Hyperlane** and **IBC** protocols. This enables users to deploy their vault positions across multiple chains while maintaining unified risk management and accounting.

## Architecture

### Multi-Provider Design

The cross-chain vault system is built with a provider-agnostic architecture:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Noble Chain   │    │  Cross-Chain     │    │  Remote Chains  │
│                 │    │   Providers      │    │                 │
│  ┌───────────┐  │    │                  │    │  ┌───────────┐  │
│  │ Vault V2  │◄─┼────┤  IBC Provider    ├────┼─►│ Osmosis   │  │
│  │  System   │  │    │                  │    │  │ Stride    │  │
│  └───────────┘  │    │  Hyperlane       │    │  └───────────┘  │
│                 │    │  Provider        │    │                 │
│                 │    │                  │    │  ┌───────────┐  │
│                 │    │                  ├────┼─►│ Ethereum  │  │
│                 │    │                  │    │  │ Arbitrum  │  │
│                 │    │                  │    │  │ Polygon   │  │
│                 │    │                  │    │  └───────────┘  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Core Components

1. **CrossChainRoute**: Defines routing configuration for each chain
2. **RemotePosition**: Tracks user positions on remote chains
3. **InFlightPosition**: Manages operations in progress
4. **Provider Interface**: Abstracts IBC and Hyperlane specifics
5. **Risk Management**: Monitors drift and applies haircuts

## Provider Configuration

### Hyperlane Configuration

```protobuf
message HyperlaneConfig {
  uint32 domain_id = 1;                    // Chain domain ID
  string mailbox_address = 2;              // Mailbox contract address
  string gas_paymaster_address = 3;        // Gas paymaster for fees
  string hook_address = 4;                 // Optional hook contract
  uint64 gas_limit = 5;                    // Gas limit for execution
  string gas_price = 6;                    // Gas price for transactions
}
```

### IBC Configuration

```protobuf
message IBCConfig {
  string channel_id = 1;                   // IBC channel ID
  string port_id = 2;                      // IBC port (default: "transfer")
  uint64 timeout_timestamp = 3;            // Packet timeout
  uint64 timeout_height = 4;               // Timeout height
}
```

## Supported Chains

### Hyperlane Networks

| Chain | Domain ID | Gas Settings | Risk Level |
|-------|-----------|--------------|------------|
| Ethereum | 1 | 20-50 gwei | Medium |
| Arbitrum | 42161 | 0.1-1 gwei | Low |
| Polygon | 137 | 30-100 gwei | Medium |
| Optimism | 10 | 0.001-0.01 gwei | Low |
| BSC | 56 | 3-10 gwei | High |
| Avalanche | 43114 | 25-50 nAVAX | Medium |

### IBC Networks

| Chain | Channel | Timeout | Risk Level |
|-------|---------|---------|------------|
| Osmosis | channel-0 | 1-2 hours | Low |
| Stride | channel-8 | 30-60 min | Medium |
| Neutron | channel-12 | 1-2 hours | Low |
| Stargaze | channel-6 | 1 hour | Medium |

## Setup Instructions

### 1. Initialize Cross-Chain Keeper

```go
import (
    "dollar.noble.xyz/v2/keeper/crosschain"
    dollarv2 "dollar.noble.xyz/v2/types/v2"
    vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// Initialize keeper with collections
crossChainKeeper := crosschain.NewCrossChainKeeper(
    routesCollection,
    positionsCollection,
    inFlightCollection,
    snapshotsCollection,
    driftAlertsCollection,
    configCollection,
    nonceCounter,
)
```

### 2. Register Providers

```go
// Register IBC Provider
ibcProvider := crosschain.NewIBCProvider(
    channelKeeper,
    transferKeeper,
    clientKeeper,
    "transfer",
    time.Hour,
)
crossChainKeeper.RegisterProvider(ibcProvider)

// Register Hyperlane Provider
hyperlaneProvider := crosschain.NewHyperlaneProvider(
    mailboxKeeper,
    gasPriceFeed,
    4, // Noble domain ID
    200000,
    time.Hour,
)
crossChainKeeper.RegisterProvider(hyperlaneProvider)
```

### 3. Create Cross-Chain Routes

```go
// Create Ethereum route via Hyperlane
ethRoute := &vaultsv2.CrossChainRoute{
    RouteId:          "noble-ethereum-hyperlane",
    SourceChain:      "noble-1",
    DestinationChain: "ethereum-1",
    Provider:         dollarv2.Provider_HYPERLANE,
    ProviderConfig: &vaultsv2.CrossChainProviderConfig{
        Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
            HyperlaneConfig: &vaultsv2.HyperlaneConfig{
                DomainId:            1,
                MailboxAddress:      "0x2f2aFaE1139Ce54feFC03593FeE8AB2aDF4a85A7",
                GasPaymasterAddress: "0x6cA0B6D22da47f091B7613223cD4BB03a2d77918",
                GasLimit:            300000,
                GasPrice:            math.NewInt(20000000000),
            },
        },
    },
    Active:           true,
    MaxPositionValue: math.NewInt(5000000000000),
    RiskParams: &vaultsv2.CrossChainRiskParams{
        PositionHaircut:      300,  // 3% haircut
        MaxDriftThreshold:    800,  // 8% max drift
        OperationTimeout:     7200, // 2 hours
        MaxRetries:           5,
        ConservativeDiscount: 150,
    },
}

err := crossChainKeeper.CreateRoute(ctx, ethRoute)
```

## Usage Examples

### Remote Deposit

```go
// Deposit to Ethereum via Hyperlane
nonce, err := crossChainKeeper.InitiateRemoteDeposit(
    ctx,
    userAddress,
    "noble-ethereum-hyperlane",
    math.NewInt(100000000), // 100 USDC
    "0x742d35Cc6474C451c4bE0C43D93C7424b1a4c3c4", // Ethereum address
    300000,                   // Gas limit
    math.NewInt(25000000000), // 25 gwei
)
```

### Remote Withdrawal

```go
// Withdraw from Arbitrum via Hyperlane
nonce, err := crossChainKeeper.InitiateRemoteWithdraw(
    ctx,
    userAddress,
    "noble-arbitrum-hyperlane",
    math.NewInt(50000000), // 50 shares
    150000,                // Gas limit
    math.NewInt(200000000), // 0.2 gwei
)
```

### Position Updates

```go
// Update position status (called by relayers)
hyperlaneTracking := &vaultsv2.ProviderTrackingInfo{
    TrackingInfo: &vaultsv2.ProviderTrackingInfo_HyperlaneTracking{
        HyperlaneTracking: &vaultsv2.HyperlaneTrackingInfo{
            MessageId:              []byte("0x1234567890abcdef"),
            OriginDomain:           4, // Noble
            DestinationDomain:      1, // Ethereum
            Nonce:                  12345,
            OriginTxHash:           "0xabcdef...",
            DestinationTxHash:      "0x123456...",
            OriginBlockNumber:      1000000,
            DestinationBlockNumber: 18500000,
            Processed:              true,
            GasUsed:                275000,
        },
    },
}

err := crossChainKeeper.UpdateRemotePosition(
    ctx,
    "noble-ethereum-hyperlane",
    userAddress.Bytes(),
    math.NewInt(102000000), // Updated value
    15,                     // Confirmations
    hyperlaneTracking,
    vaultsv2.REMOTE_POSITION_ACTIVE,
)
```

## API Reference

### Transaction Messages

#### Create Cross-Chain Route
```protobuf
message MsgCreateCrossChainRoute {
  string authority = 1;
  CrossChainRoute route = 2;
}
```

#### Remote Deposit
```protobuf
message MsgRemoteDeposit {
  string depositor = 1;
  string route_id = 2;
  VaultType vault_type = 3;
  string amount = 4;
  string remote_address = 5;
  string min_shares = 6;
  uint64 gas_limit = 7;        // Hyperlane only
  string gas_price = 8;        // Hyperlane only
}
```

#### Remote Withdrawal
```protobuf
message MsgRemoteWithdraw {
  string withdrawer = 1;
  string route_id = 2;
  VaultType vault_type = 3;
  string shares = 4;
  string min_amount = 5;
  uint64 gas_limit = 6;        // Hyperlane only
  string gas_price = 7;        // Hyperlane only
}
```

### Query Endpoints

#### Get Cross-Chain Routes
```bash
# Get all routes
GET /noble/dollar/vaults/v2/crosschain/routes

# Get specific route
GET /noble/dollar/vaults/v2/crosschain/route/{route_id}
```

#### Get Remote Positions
```bash
# Get user's remote position on a route
GET /noble/dollar/vaults/v2/crosschain/position/{route_id}/{address}

# Get all remote positions for a user
GET /noble/dollar/vaults/v2/crosschain/positions/{address}
```

#### Get In-Flight Operations
```bash
# Get specific in-flight operation
GET /noble/dollar/vaults/v2/crosschain/inflight/{nonce}

# Get all in-flight operations for a user
GET /noble/dollar/vaults/v2/crosschain/inflight/user/{address}
```

## Risk Management

### Position Haircuts

Each route applies conservative valuations through haircuts:

```go
type CrossChainRiskParams struct {
    PositionHaircut       int32 // Basis points (100 = 1%)
    MaxDriftThreshold     int32 // Maximum allowed drift
    OperationTimeout      int64 // Timeout in seconds
    MaxRetries            int32 // Retry attempts
    ConservativeDiscount  int32 // Additional safety margin
}
```

### Drift Monitoring

Positions are monitored for value drift:

- **Drift Calculation**: `(actual_value - expected_value) / expected_value * 10000`
- **Alert Threshold**: Configurable per route (typically 5-12%)
- **Actions**: Automatic alerts, potential liquidation

### Confirmation Requirements

| Provider | Confirmations | Finality Time |
|----------|---------------|---------------|
| IBC | 1 | ~30 seconds |
| Hyperlane (Ethereum) | 12 | ~3 minutes |
| Hyperlane (L2s) | 1-3 | ~30 seconds |

## Integration Examples

### Basic Integration

```go
// 1. Setup routes for your chains
routes := []string{"ethereum", "arbitrum", "osmosis"}
for _, chain := range routes {
    route := ExampleConfigurations[chain]
    err := keeper.CreateRoute(ctx, &route)
    if err != nil {
        panic(err)
    }
}

// 2. Make a remote deposit
nonce, err := keeper.InitiateRemoteDeposit(
    ctx,
    user,
    "noble-ethereum-hyperlane",
    amount,
    remoteAddress,
    gasLimit,
    gasPrice,
)

// 3. Track the operation
inFlight, err := keeper.GetInFlightPosition(ctx, nonce)
fmt.Printf("Status: %s", inFlight.Status.String())
```

### Advanced Risk Management

```go
// Monitor drift alerts
alerts, err := queryServer.DriftAlerts(ctx, &vaultsv2.QueryDriftAlertsRequest{
    Address: userAddress,
})

for _, alert := range alerts.Alerts {
    if alert.Alert.CurrentDrift > 1000 { // > 10%
        // Trigger rebalancing or liquidation
        handleHighDrift(alert)
    }
}

// Check position health
position, err := keeper.GetRemotePosition(ctx, routeId, userAddress)
healthScore := calculateHealthScore(position)
```

## Gas Management

### Hyperlane Gas Estimation

```go
// Estimate gas for operation
provider := keeper.GetProvider(dollarv2.Provider_HYPERLANE)
gasLimit, gasCost, err := provider.EstimateGas(ctx, route, message)

// Dynamic gas pricing
gasPrice := gasPriceFeed.GetCurrentPrice(route.DestinationDomain)
totalCost := gasPrice.Mul(math.NewInt(int64(gasLimit)))
```

### Gas Optimization Strategies

1. **Batch Operations**: Combine multiple operations when possible
2. **Dynamic Pricing**: Use gas price feeds for optimal pricing
3. **Route Selection**: Choose cheaper routes when available
4. **Timing**: Execute during low-congestion periods

## Security Considerations

### Message Verification

#### Hyperlane
- Messages are cryptographically verified through validator signatures
- Requires threshold of validator attestations
- Built-in fraud proofs for disputed messages

#### IBC
- Channel handshake ensures secure connection
- Client verification for each packet
- Built-in timeout and acknowledgment mechanisms

### Risk Mitigation

1. **Position Limits**: Maximum position sizes per route
2. **Haircuts**: Conservative valuations (2-5% discount)
3. **Drift Monitoring**: Real-time position tracking
4. **Emergency Controls**: Circuit breakers for critical situations

## Monitoring and Alerts

### Drift Alerts

```go
// Automatic alert generation
type DriftAlert struct {
    RouteId           string
    UserAddress       []byte
    CurrentDrift      int32     // Basis points
    ThresholdExceeded int32     // Threshold that was exceeded
    Timestamp         time.Time
    RecommendedAction string
}
```

### Operational Metrics

- **Total Remote Value**: Sum of all cross-chain positions
- **Active Positions**: Number of active remote positions
- **Failed Operations**: Count of failed cross-chain operations
- **Average Confirmation Time**: Per provider/chain
- **Gas Efficiency**: Cost per operation by chain

## Error Handling

### Common Error Scenarios

| Error Type | Hyperlane | IBC | Recovery |
|------------|-----------|-----|----------|
| Gas Exhaustion | Retry with higher gas | N/A | Manual retry |
| Network Congestion | Wait or increase gas price | Packet delay | Automatic timeout |
| Invalid Address | Immediate failure | Immediate failure | User correction |
| Insufficient Funds | Immediate failure | Immediate failure | User action required |
| Timeout | Message expires | Packet timeout | Automatic refund |

### Recovery Mechanisms

```go
// Process failed operation
err := keeper.ProcessInFlightPosition(
    ctx,
    nonce,
    vaultsv2.INFLIGHT_FAILED,
    math.ZeroInt(),
    "timeout: operation expired",
    nil,
)
```

## Performance Optimization

### Batch Processing

```go
// Process multiple operations together
func ProcessBatchOperations(ctx sdk.Context, operations []Operation) error {
    for _, op := range operations {
        // Validate operation
        if err := validateOperation(op); err != nil {
            continue // Skip invalid operations
        }
        
        // Process operation
        processOperation(ctx, op)
    }
    return nil
}
```

### Caching Strategies

- **Route Information**: Cache active routes for fast lookup
- **Gas Prices**: Cache gas prices with TTL
- **Confirmation Status**: Cache message statuses to reduce queries

## Development Guide

### Adding New Chains

1. **Determine Provider**: Choose IBC or Hyperlane based on chain support
2. **Configure Route**: Set up provider-specific configuration
3. **Test Integration**: Verify message delivery and confirmations
4. **Set Risk Parameters**: Configure haircuts and drift thresholds
5. **Deploy**: Add route through governance

### Custom Provider Implementation

```go
type CustomProvider struct {
    // Provider-specific fields
}

func (p *CustomProvider) GetProviderType() dollarv2.Provider {
    return dollarv2.Provider_CUSTOM // Would need to add to enum
}

func (p *CustomProvider) SendMessage(ctx sdk.Context, route *vaultsv2.CrossChainRoute, msg CrossChainMessage) (*vaultsv2.ProviderTrackingInfo, error) {
    // Custom implementation
    return nil, nil
}

// Implement remaining interface methods...
```

## Testing

### Integration Tests

```go
func TestHyperlaneIntegration(t *testing.T) {
    // Setup test environment
    keeper := setupTestKeeper()
    user := generateTestUser()
    
    // Create test route
    route := createTestRoute("ethereum")
    err := keeper.CreateRoute(ctx, route)
    require.NoError(t, err)
    
    // Test deposit
    nonce, err := keeper.InitiateRemoteDeposit(ctx, user, route.RouteId, amount, remoteAddr, gasLimit, gasPrice)
    require.NoError(t, err)
    
    // Verify in-flight position
    inFlight, err := keeper.GetInFlightPosition(ctx, nonce)
    require.NoError(t, err)
    assert.Equal(t, vaultsv2.INFLIGHT_PENDING, inFlight.Status)
}
```

### Security Tests

```go
func TestCrossChainSecurity(t *testing.T) {
    // Test position spoofing protection
    testPositionSpoofing(t)
    
    // Test replay attack protection
    testReplayAttack(t)
    
    // Test drift limits
    testDriftLimits(t)
    
    // Test emergency controls
    testEmergencyControls(t)
}
```

## CLI Commands

### Create Route

```bash
# Create Hyperlane route
noblex tx vaults v2 create-cross-chain-route \
  --route-id="noble-ethereum-hyperlane" \
  --source-chain="noble-1" \
  --destination-chain="ethereum-1" \
  --provider="HYPERLANE" \
  --domain-id=1 \
  --mailbox="0x2f2aFaE1139Ce54feFC03593FeE8AB2aDF4a85A7" \
  --gas-limit=300000 \
  --max-position="5000000000000" \
  --from=authority
```

### Remote Deposit

```bash
# Deposit to Ethereum
noblex tx vaults v2 remote-deposit \
  --route-id="noble-ethereum-hyperlane" \
  --vault-type="FLEXIBLE" \
  --amount="100000000" \
  --remote-address="0x742d35Cc6474C451c4bE0C43D93C7424b1a4c3c4" \
  --gas-limit=300000 \
  --gas-price="25000000000" \
  --from=user
```

### Query Positions

```bash
# Query remote positions
noblex query vaults v2 remote-positions \
  --address="noble1abc..."

# Query specific position
noblex query vaults v2 remote-position \
  --route-id="noble-ethereum-hyperlane" \
  --address="noble1abc..."
```

## Troubleshooting

### Common Issues

#### High Gas Costs
- **Problem**: Ethereum gas fees too high
- **Solution**: Use L2 routes (Arbitrum, Polygon) or wait for lower gas periods

#### Position Drift
- **Problem**: Remote position value significantly different from expected
- **Solution**: Check for MEV attacks, network issues, or legitimate market movements

#### Failed Transactions
- **Problem**: Cross-chain transaction failed
- **Solution**: Check gas limits, network status, and retry with higher gas

#### IBC Timeouts
- **Problem**: IBC packets timing out
- **Solution**: Increase timeout parameters or check channel health

### Debug Commands

```bash
# Check route status
noblex query vaults v2 cross-chain-route --route-id="noble-ethereum-hyperlane"

# Check in-flight operations
noblex query vaults v2 inflight-position --nonce=12345

# Check drift alerts
noblex query vaults v2 drift-alerts --address="noble1abc..."

# Check cross-chain snapshot
noblex query vaults v2 cross-chain-snapshot --vault-type="FLEXIBLE"
```

### Log Analysis

Monitor these log patterns:

```bash
# Successful operations
grep "cross_chain_success" logs/

# Failed operations
grep "cross_chain_error" logs/

# Gas estimation issues
grep "gas_estimation_failed" logs/

# Drift alerts
grep "drift_alert_generated" logs/
```

## Governance and Upgrades

### Route Management

Routes can be managed through governance:

```bash
# Propose new route
noblex tx gov submit-proposal create-cross-chain-route proposal.json

# Update existing route
noblex tx gov submit-proposal update-cross-chain-route proposal.json

# Disable route
noblex tx gov submit-proposal disable-cross-chain-route proposal.json
```

### Parameter Updates

```bash
# Update global cross-chain config
noblex tx gov submit-proposal update-cross-chain-config \
  --global-haircut=400 \
  --max-remote-exposure=2000 \
  --default-timeout=7200
```

## Future Enhancements

### Planned Features

1. **Cross-Chain Yield Farming**: Automatic yield optimization across chains
2. **Dynamic Route Selection**: Automatic route selection based on costs and risks
3. **Advanced Risk Modeling**: Machine learning-based drift prediction
4. **Liquidity Aggregation**: Pool liquidity across multiple chains
5. **Cross-Chain Governance**: Vote on proposals from any connected chain

### Research Areas

- **Zero-Knowledge Proofs**: Privacy-preserving cross-chain operations
- **Optimistic Rollups**: Integration with optimistic rollup chains
- **Cross-Chain AMMs**: Automated market making across chains
- **Intent-Based Operations**: User intent matching across chains

## Support

### Resources

- **Documentation**: https://docs.noble.xyz/vaults/cross-chain
- **Discord**: https://discord.gg/noble
- **GitHub**: https://github.com/noble-assets/dollar
- **Testnet**: Connect to Noble testnet for testing

### Reporting Issues

When reporting cross-chain issues, include:

1. Route ID and provider type
2. Transaction hashes (origin and destination)
3. User address and operation details
4. Error messages and timestamps
5. Gas parameters used

### Emergency Contacts

- **Security Issues**: security@noble.xyz
- **Critical Bugs**: bugs@noble.xyz
- **Integration Support**: dev@noble.xyz

---

*This documentation covers the cross-chain vault integration supporting both Hyperlane and IBC protocols. For the latest updates, check the official Noble documentation.*