package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// vaultV2QueryServer is the server API for VaultV2Query service
type vaultV2QueryServer struct {
	*Keeper
}

// NewVaultV2QueryServer returns an implementation of the V2 vault QueryServer interface
func NewVaultV2QueryServer(keeper *Keeper) vaultsv2.QueryServer {
	return &vaultV2QueryServer{Keeper: keeper}
}

var _ vaultsv2.QueryServer = vaultV2QueryServer{}

// VaultInfo implements vaultsv2.QueryServer
func (k vaultV2QueryServer) VaultInfo(ctx context.Context, req *vaultsv2.QueryVaultInfoRequest) (*vaultsv2.QueryVaultInfoResponse, error) {

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return default state for non-existent vault
			return &vaultsv2.QueryVaultInfoResponse{
				Config: vaultsv2.VaultConfig{
					Enabled: true,
				},
				TotalShares:     math.ZeroInt().String(),
				TotalNav:        math.ZeroInt().String(),
				SharePrice:      math.LegacyOneDec().String(),
				TotalDepositors: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// TODO: Get actual vault config
	config := vaultsv2.VaultConfig{
		Enabled: vaultState.DepositsEnabled,
	}

	return &vaultsv2.QueryVaultInfoResponse{
		Config:          config,
		TotalShares:     vaultState.TotalShares.String(),
		TotalNav:        vaultState.TotalNav.String(),
		SharePrice:      vaultState.SharePrice.String(),
		TotalDepositors: vaultState.TotalUsers,
	}, nil
}

// AllVaults implements vaultsv2.QueryServer
func (k vaultV2QueryServer) AllVaults(ctx context.Context, req *vaultsv2.QueryAllVaultsRequest) (*vaultsv2.QueryAllVaultsResponse, error) {
	var vaultList []vaultsv2.QueryVaultInfoResponse

	// Query for the single vault
	vaultResp, err := k.VaultInfo(ctx, &vaultsv2.QueryVaultInfoRequest{})
	if err == nil {
		vaultList = append(vaultList, *vaultResp)
	}

	return &vaultsv2.QueryAllVaultsResponse{
		Vaults: vaultList,
		// TODO: Implement pagination
	}, nil
}

// UserPosition implements vaultsv2.QueryServer
func (k vaultV2QueryServer) UserPosition(ctx context.Context, req *vaultsv2.QueryUserPositionRequest) (*vaultsv2.QueryUserPositionResponse, error) {
	// Validate address
	userAddr, err := k.address.StringToBytes(req.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid user address: %w", err)
	}

	// Get user position
	position, err := k.GetV2UserPosition(ctx, userAddr)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return empty position for user with no position
			return &vaultsv2.QueryUserPositionResponse{
				Position: &vaultsv2.UserPosition{
					Shares:             math.ZeroInt(),
					OriginalDeposit:    math.ZeroInt(),
					FirstDepositTime:   sdk.UnwrapSDKContext(ctx).BlockTime(),
					LastActivityTime:   sdk.UnwrapSDKContext(ctx).BlockTime(),
					ReceiveYield:       false,
					SharesPendingExit:  math.ZeroInt(),
					ActiveExitRequests: 0,
				},
				CurrentValue:    math.ZeroInt().String(),
				UnrealizedYield: math.ZeroInt().String(),
			}, nil
		}
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	// Calculate current value and unrealized gain
	vaultState, err := k.GetV2VaultState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	sharePrice := k.calculateV2SharePrice(vaultState)
	currentValue := sharePrice.MulInt(position.Shares).TruncateInt()
	unrealizedGain := currentValue.Sub(position.OriginalDeposit)

	return &vaultsv2.QueryUserPositionResponse{
		Position:        position,
		CurrentValue:    currentValue.String(),
		UnrealizedYield: unrealizedGain.String(),
	}, nil
}

// UserPositions implements vaultsv2.QueryServer
func (k vaultV2QueryServer) UserPositions(ctx context.Context, req *vaultsv2.QueryUserPositionsRequest) (*vaultsv2.QueryUserPositionsResponse, error) {
	var positions []vaultsv2.UserPositionWithVault

	// Check the single vault
	userPos, err := k.UserPosition(ctx, &vaultsv2.QueryUserPositionRequest{
		Address: req.Address,
	})
	if err == nil && !userPos.Position.Shares.IsZero() {
		positions = append(positions, vaultsv2.UserPositionWithVault{
			Position:     *userPos.Position,
			CurrentValue: userPos.CurrentValue,
		})
	}

	return &vaultsv2.QueryUserPositionsResponse{
		Positions: positions,
		// TODO: Implement pagination
	}, nil
}

// SharePrice implements vaultsv2.QueryServer
func (k vaultV2QueryServer) SharePrice(ctx context.Context, req *vaultsv2.QuerySharePriceRequest) (*vaultsv2.QuerySharePriceResponse, error) {
	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &vaultsv2.QuerySharePriceResponse{
				SharePrice:  math.LegacyOneDec().String(),
				TotalShares: math.ZeroInt().String(),
				TotalNav:    math.ZeroInt().String(),
			}, nil
		}
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	sharePrice := k.calculateV2SharePrice(vaultState)

	return &vaultsv2.QuerySharePriceResponse{
		SharePrice:  sharePrice.String(),
		TotalShares: vaultState.TotalShares.String(),
		TotalNav:    vaultState.TotalNav.String(),
	}, nil
}

// NAVInfo implements vaultsv2.QueryServer
func (k vaultV2QueryServer) NAVInfo(ctx context.Context, req *vaultsv2.QueryNAVInfoRequest) (*vaultsv2.QueryNAVInfoResponse, error) {

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	navInfo := vaultsv2.NAVInfo{
		CurrentNav:           vaultState.TotalNav,
		PreviousNav:          vaultState.TotalNav, // TODO: Track previous NAV
		LastUpdate:           vaultState.LastNavUpdate,
		ChangeBps:            0,     // TODO: Calculate change
		CircuitBreakerActive: false, // TODO: Implement circuit breaker
	}

	return &vaultsv2.QueryNAVInfoResponse{
		NavInfo: navInfo,
	}, nil
}

// DepositPreview implements vaultsv2.QueryServer
func (k vaultV2QueryServer) DepositPreview(ctx context.Context, req *vaultsv2.QueryDepositPreviewRequest) (*vaultsv2.QueryDepositPreviewResponse, error) {
	// Parse amount
	amount, ok := math.NewIntFromString(req.Amount)
	if !ok || amount.IsZero() || amount.IsNegative() {
		return nil, fmt.Errorf("invalid deposit amount")
	}

	// Get vault state
	vaultState, err := k.getOrCreateV2VaultState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// Calculate shares to receive
	sharePrice := k.calculateV2SharePrice(vaultState)
	sharesToReceive := math.LegacyNewDecFromInt(amount).Quo(sharePrice).TruncateInt()

	return &vaultsv2.QueryDepositPreviewResponse{
		SharesToReceive: sharesToReceive.String(),
		FeesToPay:       math.ZeroInt().String(), // TODO: Implement fees
		NetAmount:       amount.String(),
		SharePrice:      sharePrice.String(),
		FeeRateBps:      0, // TODO: Get from fee config
	}, nil
}

// WithdrawalPreview implements vaultsv2.QueryServer
func (k vaultV2QueryServer) WithdrawalPreview(ctx context.Context, req *vaultsv2.QueryWithdrawalPreviewRequest) (*vaultsv2.QueryWithdrawalPreviewResponse, error) {

	// Parse shares
	shares, ok := math.NewIntFromString(req.Shares)
	if !ok || shares.IsZero() || shares.IsNegative() {
		return nil, fmt.Errorf("invalid withdrawal shares")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// Calculate amount to receive
	sharePrice := k.calculateV2SharePrice(vaultState)
	amountToReceive := sharePrice.MulInt(shares).TruncateInt()

	return &vaultsv2.QueryWithdrawalPreviewResponse{
		AmountToReceive: amountToReceive.String(),
		FeesToPay:       math.ZeroInt().String(), // TODO: Implement fees
		GrossAmount:     amountToReceive.String(),
		SharePrice:      sharePrice.String(),
		FeeRateBps:      0, // TODO: Get from fee config
	}, nil
}

// ExitQueue implements vaultsv2.QueryServer
func (k vaultV2QueryServer) ExitQueue(ctx context.Context, req *vaultsv2.QueryExitQueueRequest) (*vaultsv2.QueryExitQueueResponse, error) {
	// In the simplified design, we don't implement exit queues
	return &vaultsv2.QueryExitQueueResponse{
		ExitRequests: []vaultsv2.ExitRequestWithUser{},
		// TODO: Implement pagination
	}, nil
}

// UserExitRequests implements vaultsv2.QueryServer
func (k vaultV2QueryServer) UserExitRequests(ctx context.Context, req *vaultsv2.QueryUserExitRequestsRequest) (*vaultsv2.QueryUserExitRequestsResponse, error) {
	// In the simplified design, we don't implement exit queues
	return &vaultsv2.QueryUserExitRequestsResponse{
		ExitRequests: []vaultsv2.ExitRequestWithVault{},
		// TODO: Implement pagination
	}, nil
}

// FeeInfo implements vaultsv2.QueryServer
func (k vaultV2QueryServer) FeeInfo(ctx context.Context, req *vaultsv2.QueryFeeInfoRequest) (*vaultsv2.QueryFeeInfoResponse, error) {
	// TODO: Get actual fee config
	feeConfig := vaultsv2.FeeConfig{
		DepositFeeRate:    0,
		WithdrawalFeeRate: 0,
	}

	return &vaultsv2.QueryFeeInfoResponse{
		FeeConfig: feeConfig,
	}, nil
}

// Stats implements vaultsv2.QueryServer
func (k vaultV2QueryServer) Stats(ctx context.Context, req *vaultsv2.QueryStatsRequest) (*vaultsv2.QueryStatsResponse, error) {
	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return zero stats for non-existent vault
			stats := vaultsv2.VaultStatsEntry{
				TotalDepositors:       0,
				TotalDeposited:        math.ZeroInt(),
				TotalWithdrawn:        math.ZeroInt(),
				TotalFeesCollected:    math.ZeroInt(),
				TotalYieldDistributed: math.ZeroInt(),
				ActivePositions:       0,
				AveragePositionSize:   math.LegacyZeroDec(),
			}
			return &vaultsv2.QueryStatsResponse{Stats: stats}, nil
		}
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// TODO: Calculate actual stats from historical data
	stats := vaultsv2.VaultStatsEntry{
		TotalDepositors:       vaultState.TotalUsers,
		TotalDeposited:        vaultState.TotalNav,   // Simplified
		TotalWithdrawn:        math.ZeroInt(),        // TODO: Track withdrawals
		TotalFeesCollected:    math.ZeroInt(),        // TODO: Track fees
		TotalYieldDistributed: math.ZeroInt(),        // TODO: Track yield
		ActivePositions:       vaultState.TotalUsers, // Simplified
		AveragePositionSize:   math.LegacyZeroDec(),  // TODO: Calculate
	}

	return &vaultsv2.QueryStatsResponse{
		Stats: stats,
	}, nil
}

// AllStats implements vaultsv2.QueryServer
func (k vaultV2QueryServer) AllStats(ctx context.Context, req *vaultsv2.QueryAllStatsRequest) (*vaultsv2.QueryAllStatsResponse, error) {
	var allStats []vaultsv2.VaultStatsEntry

	// Get stats for the single vault
	vaultResp, err := k.Stats(ctx, &vaultsv2.QueryStatsRequest{})
	if err == nil {
		allStats = append(allStats, vaultResp.Stats)
	}

	return &vaultsv2.QueryAllStatsResponse{
		Stats: allStats,
		// TODO: Implement pagination
	}, nil
}

// Params implements vaultsv2.QueryServer
func (k vaultV2QueryServer) Params(ctx context.Context, req *vaultsv2.QueryParamsRequest) (*vaultsv2.QueryParamsResponse, error) {
	// TODO: Get actual params from state
	params := vaultsv2.Params{
		Authority:                k.authority,
		DefaultDepositFeeRate:    0,
		DefaultWithdrawalFeeRate: 0,
		MinDepositAmount:         math.NewInt(1000000), // 1 USDC
		MinWithdrawalAmount:      math.NewInt(1000000), // 1 USDC
		MaxNavChangeBps:          1000,                 // 10%
		ExitRequestTimeout:       86400,                // 24 hours
		MaxExitRequestsPerBlock:  100,
		VaultsEnabled:            true,
	}

	return &vaultsv2.QueryParamsResponse{
		Params: params,
	}, nil
}

// Cross-chain query handlers

// CrossChainRoutes implements vaultsv2.QueryServer
func (k vaultV2QueryServer) CrossChainRoutes(ctx context.Context, req *vaultsv2.QueryCrossChainRoutesRequest) (*vaultsv2.QueryCrossChainRoutesResponse, error) {
	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Get all routes
	routes, err := crossChainKeeper.GetAllRoutes(sdk.UnwrapSDKContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get cross-chain routes: %w", err)
	}

	// Convert to response format
	responseRoutes := make([]vaultsv2.CrossChainRoute, len(routes))
	for i, route := range routes {
		responseRoutes[i] = *route
	}

	return &vaultsv2.QueryCrossChainRoutesResponse{
		Routes: responseRoutes,
		// TODO: Implement pagination
	}, nil
}

// CrossChainRoute implements vaultsv2.QueryServer
func (k vaultV2QueryServer) CrossChainRoute(ctx context.Context, req *vaultsv2.QueryCrossChainRouteRequest) (*vaultsv2.QueryCrossChainRouteResponse, error) {
	if req.RouteId == "" {
		return nil, fmt.Errorf("route ID must be specified")
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Get route
	route, err := crossChainKeeper.GetRoute(sdk.UnwrapSDKContext(ctx), req.RouteId)
	if err != nil {
		return nil, fmt.Errorf("failed to get cross-chain route: %w", err)
	}

	return &vaultsv2.QueryCrossChainRouteResponse{
		Route: *route,
	}, nil
}

// RemotePosition implements vaultsv2.QueryServer
func (k vaultV2QueryServer) RemotePosition(ctx context.Context, req *vaultsv2.QueryRemotePositionRequest) (*vaultsv2.QueryRemotePositionResponse, error) {
	if req.RouteId == "" {
		return nil, fmt.Errorf("route ID must be specified")
	}

	// Validate address
	userAddr, err := k.address.StringToBytes(req.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid user address: %w", err)
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Get remote position
	position, err := crossChainKeeper.GetRemotePosition(sdk.UnwrapSDKContext(ctx), req.RouteId, userAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote position: %w", err)
	}

	return &vaultsv2.QueryRemotePositionResponse{
		Position: *position,
	}, nil
}

// RemotePositions implements vaultsv2.QueryServer
func (k vaultV2QueryServer) RemotePositions(ctx context.Context, req *vaultsv2.QueryRemotePositionsRequest) (*vaultsv2.QueryRemotePositionsResponse, error) {
	// Validate address
	userAddr, err := k.address.StringToBytes(req.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid user address: %w", err)
	}

	var positions []vaultsv2.RemotePositionWithRoute

	// Walk through all remote positions and filter by user
	err = k.V2Collections.RemotePositions.Walk(sdk.UnwrapSDKContext(ctx), nil, func(key collections.Pair[string, []byte], value vaultsv2.RemotePosition) (bool, error) {
		// Check if this position belongs to the user
		if string(key.K2()) == string(userAddr) {
			positions = append(positions, vaultsv2.RemotePositionWithRoute{
				RouteId:  key.K1(),
				Position: value,
			})
		}
		return false, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get remote positions: %w", err)
	}

	return &vaultsv2.QueryRemotePositionsResponse{
		Positions: positions,
		// TODO: Implement pagination
	}, nil
}

// InFlightPosition implements vaultsv2.QueryServer
func (k vaultV2QueryServer) InFlightPosition(ctx context.Context, req *vaultsv2.QueryInFlightPositionRequest) (*vaultsv2.QueryInFlightPositionResponse, error) {
	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Get in-flight position
	position, err := crossChainKeeper.GetInFlightPosition(sdk.UnwrapSDKContext(ctx), req.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to get in-flight position: %w", err)
	}

	return &vaultsv2.QueryInFlightPositionResponse{
		Position: *position,
	}, nil
}

// InFlightPositions implements vaultsv2.QueryServer
func (k vaultV2QueryServer) InFlightPositions(ctx context.Context, req *vaultsv2.QueryInFlightPositionsRequest) (*vaultsv2.QueryInFlightPositionsResponse, error) {
	// Validate address
	userAddr, err := k.address.StringToBytes(req.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid user address: %w", err)
	}

	var positions []vaultsv2.InFlightPosition

	// Walk through all in-flight positions and filter by user
	err = k.V2Collections.InFlightPositions.Walk(sdk.UnwrapSDKContext(ctx), nil, func(key uint64, value vaultsv2.InFlightPosition) (bool, error) {
		// Check if this position belongs to the user
		if string(value.UserAddress) == string(userAddr) {
			positions = append(positions, value)
		}
		return false, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get in-flight positions: %w", err)
	}

	return &vaultsv2.QueryInFlightPositionsResponse{
		Positions: positions,
		// TODO: Implement pagination
	}, nil
}

// CrossChainSnapshot implements vaultsv2.QueryServer
func (k vaultV2QueryServer) CrossChainSnapshot(ctx context.Context, req *vaultsv2.QueryCrossChainSnapshotRequest) (*vaultsv2.QueryCrossChainSnapshotResponse, error) {

	// Determine timestamp to query
	timestamp := req.Timestamp
	if timestamp == 0 {
		timestamp = sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
	}

	// Get snapshot
	snapshot, err := k.V2Collections.CrossChainSnapshots.Get(sdk.UnwrapSDKContext(ctx), timestamp)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return empty snapshot
			snapshot = vaultsv2.CrossChainPositionSnapshot{
				TotalRemoteValue:       math.ZeroInt(),
				TotalConservativeValue: math.ZeroInt(),
				ActivePositions:        0,
				DriftExceededPositions: 0,
				Timestamp:              sdk.UnwrapSDKContext(ctx).BlockTime(),
				TotalRemoteShares:      math.ZeroInt(),
			}
		} else {
			return nil, fmt.Errorf("failed to get cross-chain snapshot: %w", err)
		}
	}

	return &vaultsv2.QueryCrossChainSnapshotResponse{
		Snapshot: snapshot,
	}, nil
}

// DriftAlerts implements vaultsv2.QueryServer
func (k vaultV2QueryServer) DriftAlerts(ctx context.Context, req *vaultsv2.QueryDriftAlertsRequest) (*vaultsv2.QueryDriftAlertsResponse, error) {
	var alerts []vaultsv2.DriftAlertWithDetails

	// Walk through all drift alerts and filter as requested
	err := k.V2Collections.DriftAlerts.Walk(sdk.UnwrapSDKContext(ctx), nil, func(key collections.Pair[string, []byte], value vaultsv2.DriftAlert) (bool, error) {
		// Apply filters
		if req.RouteId != "" && key.K1() != req.RouteId {
			return false, nil
		}

		if req.Address != "" {
			userAddr, err := k.address.StringToBytes(req.Address)
			if err == nil && string(key.K2()) != string(userAddr) {
				return false, nil
			}
		}

		// Get route and position details
		crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
		if crossChainKeeper == nil {
			return true, fmt.Errorf("cross-chain keeper not initialized")
		}

		route, err := crossChainKeeper.GetRoute(sdk.UnwrapSDKContext(ctx), key.K1())
		if err != nil {
			return false, nil // Skip if route not found
		}

		position, err := crossChainKeeper.GetRemotePosition(sdk.UnwrapSDKContext(ctx), key.K1(), key.K2())
		if err != nil {
			return false, nil // Skip if position not found
		}

		alerts = append(alerts, vaultsv2.DriftAlertWithDetails{
			Alert:    value,
			Route:    *route,
			Position: *position,
		})

		return false, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get drift alerts: %w", err)
	}

	return &vaultsv2.QueryDriftAlertsResponse{
		Alerts: alerts,
		// TODO: Implement pagination
	}, nil
}
