package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/v2/types/vaults"
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
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return default state for non-existent vault
			return &vaultsv2.QueryVaultInfoResponse{
				Config: vaultsv2.VaultConfig{
					VaultType: req.VaultType,
					Enabled:   true,
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
		VaultType: req.VaultType,
		Enabled:   vaultState.DepositsEnabled,
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

	// Query for STAKED vault
	stakedResp, err := k.VaultInfo(ctx, &vaultsv2.QueryVaultInfoRequest{VaultType: vaults.STAKED})
	if err == nil {
		vaultList = append(vaultList, *stakedResp)
	}

	// Query for FLEXIBLE vault
	flexibleResp, err := k.VaultInfo(ctx, &vaultsv2.QueryVaultInfoRequest{VaultType: vaults.FLEXIBLE})
	if err == nil {
		vaultList = append(vaultList, *flexibleResp)
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

	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get user position
	position, err := k.GetV2UserPosition(ctx, req.VaultType, userAddr)
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
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
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

	// Check STAKED vault
	stakedPos, err := k.UserPosition(ctx, &vaultsv2.QueryUserPositionRequest{
		Address:   req.Address,
		VaultType: vaults.STAKED,
	})
	if err == nil && !stakedPos.Position.Shares.IsZero() {
		positions = append(positions, vaultsv2.UserPositionWithVault{
			VaultType:    vaults.STAKED,
			Position:     *stakedPos.Position,
			CurrentValue: stakedPos.CurrentValue,
		})
	}

	// Check FLEXIBLE vault
	flexiblePos, err := k.UserPosition(ctx, &vaultsv2.QueryUserPositionRequest{
		Address:   req.Address,
		VaultType: vaults.FLEXIBLE,
	})
	if err == nil && !flexiblePos.Position.Shares.IsZero() {
		positions = append(positions, vaultsv2.UserPositionWithVault{
			VaultType:    vaults.FLEXIBLE,
			Position:     *flexiblePos.Position,
			CurrentValue: flexiblePos.CurrentValue,
		})
	}

	return &vaultsv2.QueryUserPositionsResponse{
		Positions: positions,
		// TODO: Implement pagination
	}, nil
}

// SharePrice implements vaultsv2.QueryServer
func (k vaultV2QueryServer) SharePrice(ctx context.Context, req *vaultsv2.QuerySharePriceRequest) (*vaultsv2.QuerySharePriceResponse, error) {
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
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
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
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
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Parse amount
	amount, ok := math.NewIntFromString(req.Amount)
	if !ok || amount.IsZero() || amount.IsNegative() {
		return nil, fmt.Errorf("invalid deposit amount")
	}

	// Get vault state
	vaultState, err := k.getOrCreateV2VaultState(ctx, req.VaultType)
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
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Parse shares
	shares, ok := math.NewIntFromString(req.Shares)
	if !ok || shares.IsZero() || shares.IsNegative() {
		return nil, fmt.Errorf("invalid withdrawal shares")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
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
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

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
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return zero stats for non-existent vault
			stats := vaultsv2.VaultStatsEntry{
				VaultType:             req.VaultType,
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
		VaultType:             req.VaultType,
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

	// Get stats for STAKED vault
	stakedResp, err := k.Stats(ctx, &vaultsv2.QueryStatsRequest{VaultType: vaults.STAKED})
	if err == nil {
		allStats = append(allStats, stakedResp.Stats)
	}

	// Get stats for FLEXIBLE vault
	flexibleResp, err := k.Stats(ctx, &vaultsv2.QueryStatsRequest{VaultType: vaults.FLEXIBLE})
	if err == nil {
		allStats = append(allStats, flexibleResp.Stats)
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
