package examples

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/v2/keeper/crosschain"
	dollarv2 "dollar.noble.xyz/v2/types/v2"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// HyperlaneIntegrationExample demonstrates how to set up and use cross-chain vaults
// with both Hyperlane and IBC support
type HyperlaneIntegrationExample struct {
	keeper            *crosschain.CrossChainKeeper
	ibcProvider       *crosschain.IBCProvider
	hyperlaneProvider *crosschain.HyperlaneProvider
}

// SetupExample initializes the cross-chain vault system with both providers
func SetupExample(keeper *crosschain.CrossChainKeeper) *HyperlaneIntegrationExample {
	example := &HyperlaneIntegrationExample{
		keeper: keeper,
	}

	// Initialize providers
	example.setupProviders()

	return example
}

// setupProviders initializes both IBC and Hyperlane providers
func (e *HyperlaneIntegrationExample) setupProviders() {
	// Setup IBC Provider
	// Note: In real implementation, these would be injected dependencies
	e.ibcProvider = crosschain.NewIBCProvider(
		nil,        // IBCChannelKeeper - would be injected
		nil,        // IBCTransferKeeper - would be injected
		nil,        // IBCClientKeeper - would be injected
		"transfer", // Port ID
		time.Hour,  // Default timeout
	)

	// Setup Hyperlane Provider
	e.hyperlaneProvider = crosschain.NewHyperlaneProvider(
		nil,       // HyperlaneMailboxKeeper - would be injected
		nil,       // HyperlaneGasPriceFeed - would be injected
		4,         // Noble domain ID (example)
		200000,    // Default gas limit
		time.Hour, // Default timeout
	)

	// Register providers with keeper
	e.keeper.RegisterProvider(e.ibcProvider)
	e.keeper.RegisterProvider(e.hyperlaneProvider)
}

// CreateExampleRoutes creates example cross-chain routes for different scenarios
func (e *HyperlaneIntegrationExample) CreateExampleRoutes(ctx sdk.Context) error {
	routes := []vaultsv2.CrossChainRoute{
		// IBC Route to Osmosis
		{
			RouteId:          "noble-osmosis-ibc",
			SourceChain:      "noble-1",
			DestinationChain: "osmosis-1",
			Provider:         dollarv2.Provider_IBC,
			ProviderConfig: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_IbcConfig{
					IbcConfig: &vaultsv2.IBCConfig{
						ChannelId:        "channel-0",
						PortId:           "transfer",
						TimeoutTimestamp: uint64(time.Now().Add(time.Hour).Unix()),
						TimeoutHeight:    0,
					},
				},
			},
			Active:           true,
			MaxPositionValue: math.NewInt(1000000000000), // 1M USDC
			RiskParams: &vaultsv2.CrossChainRiskParams{
				PositionHaircut:      500,  // 5% haircut
				MaxDriftThreshold:    1000, // 10% max drift
				OperationTimeout:     3600, // 1 hour
				MaxRetries:           3,
				ConservativeDiscount: 200, // 2% additional discount
			},
		},

		// Hyperlane Route to Ethereum
		{
			RouteId:          "noble-ethereum-hyperlane",
			SourceChain:      "noble-1",
			DestinationChain: "ethereum-1",
			Provider:         dollarv2.Provider_HYPERLANE,
			ProviderConfig: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
					HyperlaneConfig: &vaultsv2.HyperlaneConfig{
						DomainId:            1,                                            // Ethereum mainnet
						MailboxAddress:      "0x2f2aFaE1139Ce54feFC03593FeE8AB2aDF4a85A7", // Example Hyperlane mailbox
						GasPaymasterAddress: "0x6cA0B6D22da47f091B7613223cD4BB03a2d77918", // Example gas paymaster
						HookAddress:         "",                                           // Optional
						GasLimit:            300000,
						GasPrice:            math.NewInt(20000000000), // 20 gwei
					},
				},
			},
			Active:           true,
			MaxPositionValue: math.NewInt(5000000000000), // 5M USDC
			RiskParams: &vaultsv2.CrossChainRiskParams{
				PositionHaircut:      300,  // 3% haircut
				MaxDriftThreshold:    800,  // 8% max drift
				OperationTimeout:     7200, // 2 hours
				MaxRetries:           5,
				ConservativeDiscount: 150, // 1.5% additional discount
			},
		},

		// Hyperlane Route to Arbitrum
		{
			RouteId:          "noble-arbitrum-hyperlane",
			SourceChain:      "noble-1",
			DestinationChain: "arbitrum-one",
			Provider:         dollarv2.Provider_HYPERLANE,
			ProviderConfig: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
					HyperlaneConfig: &vaultsv2.HyperlaneConfig{
						DomainId:            42161,                                        // Arbitrum One
						MailboxAddress:      "0x979Ca5202784112f4738403dBec5D0F3B9daabB9", // Example Arbitrum mailbox
						GasPaymasterAddress: "0xDd5E2bC2acF0d1A6B9dD7768D89eD8bE5a6Fc1D1", // Example gas paymaster
						GasLimit:            150000,
						GasPrice:            math.NewInt(100000000), // 0.1 gwei (cheaper on L2)
					},
				},
			},
			Active:           true,
			MaxPositionValue: math.NewInt(2000000000000), // 2M USDC
			RiskParams: &vaultsv2.CrossChainRiskParams{
				PositionHaircut:      200,  // 2% haircut (lower risk for L2)
				MaxDriftThreshold:    1200, // 12% max drift
				OperationTimeout:     1800, // 30 minutes (faster L2)
				MaxRetries:           3,
				ConservativeDiscount: 100, // 1% additional discount
			},
		},

		// IBC Route to Stride
		{
			RouteId:          "noble-stride-ibc",
			SourceChain:      "noble-1",
			DestinationChain: "stride-1",
			Provider:         dollarv2.Provider_IBC,
			ProviderConfig: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_IbcConfig{
					IbcConfig: &vaultsv2.IBCConfig{
						ChannelId:        "channel-8",
						PortId:           "transfer",
						TimeoutTimestamp: uint64(time.Now().Add(2 * time.Hour).Unix()),
						TimeoutHeight:    0,
					},
				},
			},
			Active:           true,
			MaxPositionValue: math.NewInt(500000000000), // 500K USDC
			RiskParams: &vaultsv2.CrossChainRiskParams{
				PositionHaircut:      400,  // 4% haircut
				MaxDriftThreshold:    600,  // 6% max drift
				OperationTimeout:     1800, // 30 minutes
				MaxRetries:           2,
				ConservativeDiscount: 300, // 3% additional discount
			},
		},
	}

	// Create all routes
	for _, route := range routes {
		if err := e.keeper.CreateRoute(ctx, &route); err != nil {
			return fmt.Errorf("failed to create route %s: %w", route.RouteId, err)
		}
		fmt.Printf("✅ Created %s route: %s\n", route.Provider.String(), route.RouteId)
	}

	return nil
}

// DemoRemoteDeposits demonstrates remote deposits to different chains
func (e *HyperlaneIntegrationExample) DemoRemoteDeposits(ctx sdk.Context, user sdk.AccAddress) error {
	fmt.Println("\n🚀 Demonstrating Remote Deposits...")

	// Example 1: Deposit to Ethereum via Hyperlane
	fmt.Println("\n1. Depositing to Ethereum via Hyperlane...")
	ethNonce, err := e.keeper.InitiateRemoteDeposit(
		ctx,
		user,
		"noble-ethereum-hyperlane",
		math.NewInt(100000000), // 100 USDC
		"0x742d35Cc6474C451c4bE0C43D93C7424b1a4c3c4", // Example Ethereum address
		300000,                   // Gas limit
		math.NewInt(25000000000), // 25 gwei
	)
	if err != nil {
		return fmt.Errorf("failed to initiate Ethereum deposit: %w", err)
	}
	fmt.Printf("   ✅ Initiated deposit with nonce: %d\n", ethNonce)

	// Example 2: Deposit to Arbitrum via Hyperlane
	fmt.Println("\n2. Depositing to Arbitrum via Hyperlane...")
	arbNonce, err := e.keeper.InitiateRemoteDeposit(
		ctx,
		user,
		"noble-arbitrum-hyperlane",
		math.NewInt(50000000), // 50 USDC
		"0x1234567890123456789012345678901234567890", // Example Arbitrum address
		150000,                 // Gas limit (lower for L2)
		math.NewInt(200000000), // 0.2 gwei
	)
	if err != nil {
		return fmt.Errorf("failed to initiate Arbitrum deposit: %w", err)
	}
	fmt.Printf("   ✅ Initiated deposit with nonce: %d\n", arbNonce)

	// Example 3: Deposit to Osmosis via IBC
	fmt.Println("\n3. Depositing to Osmosis via IBC...")
	osmoNonce, err := e.keeper.InitiateRemoteDeposit(
		ctx,
		user,
		"noble-osmosis-ibc",
		math.NewInt(75000000), // 75 USDC
		"osmo1abcdefghijklmnopqrstuvwxyzabcdefghijklm", // Example Osmosis address
		0,              // Gas limit (not used for IBC)
		math.ZeroInt(), // Gas price (not used for IBC)
	)
	if err != nil {
		return fmt.Errorf("failed to initiate Osmosis deposit: %w", err)
	}
	fmt.Printf("   ✅ Initiated deposit with nonce: %d\n", osmoNonce)

	return nil
}

// DemoRemoteWithdrawals demonstrates remote withdrawals from different chains
func (e *HyperlaneIntegrationExample) DemoRemoteWithdrawals(ctx sdk.Context, user sdk.AccAddress) error {
	fmt.Println("\n💰 Demonstrating Remote Withdrawals...")

	// Example 1: Withdraw from Ethereum via Hyperlane
	fmt.Println("\n1. Withdrawing from Ethereum via Hyperlane...")
	ethNonce, err := e.keeper.InitiateRemoteWithdraw(
		ctx,
		user,
		"noble-ethereum-hyperlane",
		math.NewInt(50000000),    // 50 shares
		350000,                   // Gas limit (higher for withdrawal)
		math.NewInt(30000000000), // 30 gwei
	)
	if err != nil {
		return fmt.Errorf("failed to initiate Ethereum withdrawal: %w", err)
	}
	fmt.Printf("   ✅ Initiated withdrawal with nonce: %d\n", ethNonce)

	// Example 2: Withdraw from Arbitrum via Hyperlane
	fmt.Println("\n2. Withdrawing from Arbitrum via Hyperlane...")
	arbNonce, err := e.keeper.InitiateRemoteWithdraw(
		ctx,
		user,
		"noble-arbitrum-hyperlane",
		math.NewInt(25000000),  // 25 shares
		200000,                 // Gas limit
		math.NewInt(500000000), // 0.5 gwei
	)
	if err != nil {
		return fmt.Errorf("failed to initiate Arbitrum withdrawal: %w", err)
	}
	fmt.Printf("   ✅ Initiated withdrawal with nonce: %d\n", arbNonce)

	return nil
}

// DemoStatusUpdates demonstrates how to update position statuses
func (e *HyperlaneIntegrationExample) DemoStatusUpdates(ctx sdk.Context, user sdk.AccAddress) error {
	fmt.Println("\n📊 Demonstrating Status Updates...")

	// Example: Update Ethereum position via Hyperlane
	fmt.Println("\n1. Updating Ethereum position status...")

	// Create example Hyperlane tracking info
	hyperlaneTracking := &vaultsv2.ProviderTrackingInfo{
		TrackingInfo: &vaultsv2.ProviderTrackingInfo_HyperlaneTracking{
			HyperlaneTracking: &vaultsv2.HyperlaneTrackingInfo{
				MessageId:              []byte("0x1234567890abcdef"),
				OriginDomain:           4, // Noble
				DestinationDomain:      1, // Ethereum
				Nonce:                  12345,
				OriginTxHash:           "0xabcdef1234567890abcdef1234567890abcdef12",
				DestinationTxHash:      "0x1234567890abcdef1234567890abcdef12345678",
				OriginBlockNumber:      1000000,
				DestinationBlockNumber: 18500000,
				Processed:              true,
				GasUsed:                275000,
			},
		},
	}

	err := e.keeper.UpdateRemotePosition(
		ctx,
		"noble-ethereum-hyperlane",
		user.Bytes(),
		math.NewInt(102000000), // Updated value: 102 USDC (gained 2 USDC)
		15,                     // 15 confirmations
		hyperlaneTracking,
		vaultsv2.REMOTE_POSITION_ACTIVE,
	)
	if err != nil {
		return fmt.Errorf("failed to update Ethereum position: %w", err)
	}
	fmt.Printf("   ✅ Updated Ethereum position with 15 confirmations\n")

	// Example: Update Osmosis position via IBC
	fmt.Println("\n2. Updating Osmosis position status...")

	// Create example IBC tracking info
	ibcTracking := &vaultsv2.ProviderTrackingInfo{
		TrackingInfo: &vaultsv2.ProviderTrackingInfo_IbcTracking{
			IbcTracking: &vaultsv2.IBCTrackingInfo{
				Sequence:           456,
				SourceChannel:      "channel-0",
				SourcePort:         "transfer",
				DestinationChannel: "channel-120",
				DestinationPort:    "transfer",
				TimeoutTimestamp:   uint64(time.Now().Add(time.Hour).Unix()),
				TimeoutHeight:      0,
				AckReceived:        true,
				AckData:            []byte("success"),
			},
		},
	}

	err = e.keeper.UpdateRemotePosition(
		ctx,
		"noble-osmosis-ibc",
		user.Bytes(),
		math.NewInt(76500000), // Updated value: 76.5 USDC (gained 1.5 USDC)
		1,                     // 1 confirmation (IBC finality)
		ibcTracking,
		vaultsv2.REMOTE_POSITION_ACTIVE,
	)
	if err != nil {
		return fmt.Errorf("failed to update Osmosis position: %w", err)
	}
	fmt.Printf("   ✅ Updated Osmosis position with IBC acknowledgment\n")

	return nil
}

// DemoRiskManagement demonstrates risk management features
func (e *HyperlaneIntegrationExample) DemoRiskManagement(ctx sdk.Context, user sdk.AccAddress) error {
	fmt.Println("\n⚠️  Demonstrating Risk Management...")

	// Example: Simulate a position with high drift
	fmt.Println("\n1. Simulating high drift scenario...")

	highDriftTracking := &vaultsv2.ProviderTrackingInfo{
		TrackingInfo: &vaultsv2.ProviderTrackingInfo_HyperlaneTracking{
			HyperlaneTracking: &vaultsv2.HyperlaneTrackingInfo{
				MessageId:              []byte("0xdrift123456789abc"),
				OriginDomain:           4,
				DestinationDomain:      137, // Polygon
				Nonce:                  54321,
				OriginTxHash:           "0xdrift1234567890abcdef1234567890abcdef12",
				DestinationTxHash:      "0xdrift7890abcdef1234567890abcdef12345678",
				OriginBlockNumber:      1000100,
				DestinationBlockNumber: 45000000,
				Processed:              true,
				GasUsed:                180000,
			},
		},
	}

	// Update position with significant value change (triggering drift alert)
	err := e.keeper.UpdateRemotePosition(
		ctx,
		"noble-ethereum-hyperlane",
		user.Bytes(),
		math.NewInt(85000000), // Dropped to 85 USDC (15% loss, exceeds 10% drift threshold)
		20,
		highDriftTracking,
		vaultsv2.REMOTE_POSITION_DRIFT_EXCEEDED,
	)
	if err != nil {
		return fmt.Errorf("failed to update position with drift: %w", err)
	}
	fmt.Printf("   ⚠️  Position updated with drift alert generated\n")

	return nil
}

// DemoProviderSpecificFeatures demonstrates provider-specific features
func (e *HyperlaneIntegrationExample) DemoProviderSpecificFeatures(ctx sdk.Context) error {
	fmt.Println("\n🔧 Demonstrating Provider-Specific Features...")

	// Hyperlane gas estimation
	fmt.Println("\n1. Hyperlane Gas Estimation...")
	route, err := e.keeper.GetRoute(ctx, "noble-ethereum-hyperlane")
	if err != nil {
		return fmt.Errorf("failed to get Ethereum route: %w", err)
	}

	msg := crosschain.CrossChainMessage{
		Type:      crosschain.MessageTypeDeposit,
		Amount:    math.NewInt(100000000), // 100 USDC
		Recipient: "0x742d35Cc6474C451c4bE0C43D93C7424b1a4c3c4",
	}

	gasLimit, gasCost, err := e.hyperlaneProvider.EstimateGas(ctx, route, msg)
	if err != nil {
		return fmt.Errorf("failed to estimate Hyperlane gas: %w", err)
	}
	fmt.Printf("   📊 Estimated gas: %d units, cost: %s wei\n", gasLimit, gasCost.String())

	// IBC channel validation
	fmt.Println("\n2. IBC Channel Validation...")
	ibcRoute, err := e.keeper.GetRoute(ctx, "noble-osmosis-ibc")
	if err != nil {
		return fmt.Errorf("failed to get Osmosis route: %w", err)
	}

	err = e.ibcProvider.ValidateConfig(ibcRoute.ProviderConfig)
	if err != nil {
		return fmt.Errorf("IBC config validation failed: %w", err)
	}
	fmt.Printf("   ✅ IBC channel configuration validated\n")

	return nil
}

// DemoFullUserJourney demonstrates a complete user journey across both providers
func (e *HyperlaneIntegrationExample) DemoFullUserJourney(ctx sdk.Context, user sdk.AccAddress) error {
	fmt.Println("\n🌟 Complete User Journey Demo...")
	fmt.Printf("User: %s\n", user.String())

	// Step 1: Create routes
	fmt.Println("\n📋 Step 1: Setting up cross-chain routes...")
	if err := e.CreateExampleRoutes(ctx); err != nil {
		return fmt.Errorf("failed to create routes: %w", err)
	}

	// Step 2: Make deposits to different chains
	fmt.Println("\n💸 Step 2: Making remote deposits...")
	if err := e.DemoRemoteDeposits(ctx, user); err != nil {
		return fmt.Errorf("failed to demo deposits: %w", err)
	}

	// Step 3: Simulate position updates
	fmt.Println("\n🔄 Step 3: Updating position statuses...")
	if err := e.DemoStatusUpdates(ctx, user); err != nil {
		return fmt.Errorf("failed to demo status updates: %w", err)
	}

	// Step 4: Make withdrawals
	fmt.Println("\n💰 Step 4: Making remote withdrawals...")
	if err := e.DemoRemoteWithdrawals(ctx, user); err != nil {
		return fmt.Errorf("failed to demo withdrawals: %w", err)
	}

	// Step 5: Demonstrate risk management
	fmt.Println("\n⚠️  Step 5: Risk management scenarios...")
	if err := e.DemoRiskManagement(ctx, user); err != nil {
		return fmt.Errorf("failed to demo risk management: %w", err)
	}

	// Step 6: Show provider-specific features
	fmt.Println("\n🔧 Step 6: Provider-specific features...")
	if err := e.DemoProviderSpecificFeatures(ctx); err != nil {
		return fmt.Errorf("failed to demo provider features: %w", err)
	}

	fmt.Println("\n✨ User journey demo completed successfully!")
	return nil
}

// GetUserPositionSummary provides a summary of a user's cross-chain positions
func (e *HyperlaneIntegrationExample) GetUserPositionSummary(ctx sdk.Context, user sdk.AccAddress) (*UserPositionSummary, error) {
	summary := &UserPositionSummary{
		UserAddress:         user.String(),
		TotalPositions:      0,
		TotalValue:          math.ZeroInt(),
		TotalShares:         math.ZeroInt(),
		PositionsByChain:    make(map[string]PositionInfo),
		PositionsByProvider: make(map[string]ProviderSummary),
	}

	// Walk through all remote positions for this user
	// Note: This would need to be implemented as a method in the crosschain keeper
	// For now, we'll return a placeholder summary
	err := error(nil)
	// TODO: Implement actual position walking through keeper methods
	// This would typically involve adding a GetUserRemotePositions method to the crosschain keeper

	return summary, err
}

// UserPositionSummary provides a comprehensive view of a user's cross-chain positions
type UserPositionSummary struct {
	UserAddress         string                     `json:"user_address"`
	TotalPositions      int                        `json:"total_positions"`
	TotalValue          math.Int                   `json:"total_value"`
	TotalShares         math.Int                   `json:"total_shares"`
	PositionsByChain    map[string]PositionInfo    `json:"positions_by_chain"`
	PositionsByProvider map[string]ProviderSummary `json:"positions_by_provider"`
}

// PositionInfo contains aggregated position information for a chain
type PositionInfo struct {
	Value     math.Int `json:"value"`
	Shares    math.Int `json:"shares"`
	Positions int      `json:"positions"`
}

// ProviderSummary contains aggregated information for a provider
type ProviderSummary struct {
	TotalValue      math.Int `json:"total_value"`
	TotalShares     math.Int `json:"total_shares"`
	ActivePositions int      `json:"active_positions"`
}

// PrintSummary prints a formatted summary of user positions
func (summary *UserPositionSummary) PrintSummary() {
	fmt.Printf("\n📊 Position Summary for %s\n", summary.UserAddress)
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Total Positions: %d\n", summary.TotalPositions)
	fmt.Printf("Total Value: %s USDC\n", summary.TotalValue.String())
	fmt.Printf("Total Shares: %s\n", summary.TotalShares.String())

	fmt.Printf("\n🌐 Positions by Chain:\n")
	for chain, info := range summary.PositionsByChain {
		fmt.Printf("  %s: %s USDC (%s shares) across %d positions\n",
			chain, info.Value.String(), info.Shares.String(), info.Positions)
	}

	fmt.Printf("\n🔗 Positions by Provider:\n")
	for provider, info := range summary.PositionsByProvider {
		fmt.Printf("  %s: %s USDC (%s shares) across %d positions\n",
			provider, info.TotalValue.String(), info.TotalShares.String(), info.ActivePositions)
	}
	fmt.Printf("═══════════════════════════════════════════════════\n")
}

// Configuration examples for different chains
var ExampleConfigurations = map[string]vaultsv2.CrossChainRoute{
	"ethereum": {
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
					GasPrice:            math.NewInt(20000000000), // 20 gwei
				},
			},
		},
		Active:           true,
		MaxPositionValue: math.NewInt(5000000000000), // 5M USDC
		RiskParams: &vaultsv2.CrossChainRiskParams{
			PositionHaircut:      300,
			MaxDriftThreshold:    800,
			OperationTimeout:     7200,
			MaxRetries:           5,
			ConservativeDiscount: 150,
		},
	},

	"arbitrum": {
		RouteId:          "noble-arbitrum-hyperlane",
		SourceChain:      "noble-1",
		DestinationChain: "arbitrum-one",
		Provider:         dollarv2.Provider_HYPERLANE,
		ProviderConfig: &vaultsv2.CrossChainProviderConfig{
			Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
				HyperlaneConfig: &vaultsv2.HyperlaneConfig{
					DomainId:            42161,
					MailboxAddress:      "0x979Ca5202784112f4738403dBec5D0F3B9daabB9",
					GasPaymasterAddress: "0xDd5E2bC2acF0d1A6B9dD7768D89eD8bE5a6Fc1D1",
					GasLimit:            150000,
					GasPrice:            math.NewInt(100000000), // 0.1 gwei
				},
			},
		},
		Active:           true,
		MaxPositionValue: math.NewInt(2000000000000), // 2M USDC
		RiskParams: &vaultsv2.CrossChainRiskParams{
			PositionHaircut:      200,
			MaxDriftThreshold:    1200,
			OperationTimeout:     1800,
			MaxRetries:           3,
			ConservativeDiscount: 100,
		},
	},

	"polygon": {
		RouteId:          "noble-polygon-hyperlane",
		SourceChain:      "noble-1",
		DestinationChain: "polygon-mainnet",
		Provider:         dollarv2.Provider_HYPERLANE,
		ProviderConfig: &vaultsv2.CrossChainProviderConfig{
			Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
				HyperlaneConfig: &vaultsv2.HyperlaneConfig{
					DomainId:            137,
					MailboxAddress:      "0x5d934f4e2f797775e53561bB72aca21ba36B96BB",
					GasPaymasterAddress: "0x8105a095368f3c0c80d3D0D52C84b027E5E3F078",
					GasLimit:            200000,
					GasPrice:            math.NewInt(30000000000), // 30 gwei
				},
			},
		},
		Active:           true,
		MaxPositionValue: math.NewInt(1000000000000), // 1M USDC
		RiskParams: &vaultsv2.CrossChainRiskParams{
			PositionHaircut:      350,
			MaxDriftThreshold:    1000,
			OperationTimeout:     3600,
			MaxRetries:           4,
			ConservativeDiscount: 200,
		},
	},

	"osmosis": {
		RouteId:          "noble-osmosis-ibc",
		SourceChain:      "noble-1",
		DestinationChain: "osmosis-1",
		Provider:         dollarv2.Provider_IBC,
		ProviderConfig: &vaultsv2.CrossChainProviderConfig{
			Config: &vaultsv2.CrossChainProviderConfig_IbcConfig{
				IbcConfig: &vaultsv2.IBCConfig{
					ChannelId:        "channel-0",
					PortId:           "transfer",
					TimeoutTimestamp: uint64(time.Now().Add(2 * time.Hour).Unix()),
					TimeoutHeight:    0,
				},
			},
		},
		Active:           true,
		MaxPositionValue: math.NewInt(500000000000), // 500K USDC
		RiskParams: &vaultsv2.CrossChainRiskParams{
			PositionHaircut:      400,
			MaxDriftThreshold:    600,
			OperationTimeout:     1800,
			MaxRetries:           2,
			ConservativeDiscount: 300,
		},
	},
}
