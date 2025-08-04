package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/v2/types/vaults"
)

// vaultV2QueryServer is the server API for VaultV2Query service
type vaultV2QueryServer struct {
	*Keeper
}

// NewVaultV2QueryServer returns an implementation of the V2 vault QueryServer interface
func NewVaultV2QueryServer(keeper *Keeper) vaults.VaultV2QueryServer {
	return &vaultV2QueryServer{Keeper: keeper}
}

var _ vaults.VaultV2QueryServer = vaultV2QueryServer{}

// VaultState implements vaults.VaultV2QueryServer.
func (k vaultV2QueryServer) VaultState(ctx context.Context, req *vaults.QueryVaultStateRequest) (*vaults.QueryVaultStateResponse, error) {
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// Get NAV config
	navConfig, err := k.getNAVConfig(ctx, req.VaultType)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, fmt.Errorf("failed to get NAV config: %w", err)
	}

	// Get fee config
	feeConfig, err := k.getFeeConfig(ctx, req.VaultType)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, fmt.Errorf("failed to get fee config: %w", err)
	}

	return &vaults.QueryVaultStateResponse{
		VaultState: &vaultState,
		NavConfig:  &navConfig,
		FeeConfig:  &feeConfig,
	}, nil
}

// UserPosition implements vaults.VaultV2QueryServer.
func (k vaultV2QueryServer) UserPosition(ctx context.Context, req *vaults.QueryUserPositionRequest) (*vaults.QueryUserPositionResponse, error) {
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
			return &vaults.QueryUserPositionResponse{
				Position: &vaults.UserPosition{
					Shares:             math.ZeroInt(),
					PrincipalDeposited: math.ZeroInt(),
					AvgEntryPrice:      math.LegacyZeroDec(),
					FirstDeposit:       sdk.UnwrapSDKContext(ctx).BlockTime(),
					LastActivity:       sdk.UnwrapSDKContext(ctx).BlockTime(),
					ForgoYield:         false,
					ExitRequests:       []*vaults.ExitRequest{},
				},
				CurrentValue:   math.ZeroInt(),
				UnrealizedGain: math.ZeroInt(),
			}, nil
		}
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	// Calculate current value and unrealized gain
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	sharePrice := k.calculateSharePrice(vaultState)
	currentValue := sharePrice.MulInt(position.Shares).TruncateInt()
	unrealizedGain := currentValue.Sub(position.PrincipalDeposited)

	return &vaults.QueryUserPositionResponse{
		Position:       &position,
		CurrentValue:   currentValue,
		UnrealizedGain: unrealizedGain,
	}, nil
}

// SharePrice implements vaults.VaultV2QueryServer.
func (k vaultV2QueryServer) SharePrice(ctx context.Context, req *vaults.QuerySharePriceRequest) (*vaults.QuerySharePriceResponse, error) {
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	sharePrice := k.calculateSharePrice(vaultState)

	return &vaults.QuerySharePriceResponse{
		SharePrice: sharePrice,
		LastUpdate: vaultState.LastNavUpdate,
	}, nil
}

// PricingInfo implements vaults.VaultV2QueryServer.
func (k vaultV2QueryServer) PricingInfo(ctx context.Context, req *vaults.QueryPricingInfoRequest) (*vaults.QueryPricingInfoResponse, error) {
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	if req.Amount.IsZero() || req.Amount.IsNegative() {
		return nil, fmt.Errorf("amount must be positive")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	var pricingInfo vaults.PricingInfo

	if req.IsDeposit {
		// Calculate deposit pricing
		shareCalc, err := k.calculateDepositShares(ctx, req.VaultType, req.Amount, vaultState)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate deposit shares: %w", err)
		}

		pricingInfo = vaults.PricingInfo{
			SharePrice:       shareCalc.SharePrice,
			EffectiveFeeRate: k.getEffectiveDepositFeeRate(ctx, req.VaultType),
			ExpectedAmount:   shareCalc.SharesAfterFees,
		}
	} else {
		// Calculate withdrawal pricing
		// Convert amount to shares first
		sharePrice := k.calculateSharePrice(vaultState)
		shares := sharePrice.Quo(math.LegacyOneDec()).MulInt(req.Amount).TruncateInt()

		withdrawCalc, err := k.calculateWithdrawalTokens(ctx, req.VaultType, shares, vaultState)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate withdrawal tokens: %w", err)
		}

		pricingInfo = vaults.PricingInfo{
			SharePrice:       withdrawCalc.SharePrice,
			EffectiveFeeRate: k.getEffectiveWithdrawalFeeRate(ctx, req.VaultType),
			ExpectedAmount:   withdrawCalc.TokensAfterFees,
		}
	}

	return &vaults.QueryPricingInfoResponse{
		PricingInfo: &pricingInfo,
	}, nil
}

// ExitQueue implements vaults.VaultV2QueryServer.
func (k vaultV2QueryServer) ExitQueue(ctx context.Context, req *vaults.QueryExitQueueRequest) (*vaults.QueryExitQueueResponse, error) {
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	var userRequests []*vaults.ExitRequest
	var totalQueueLength uint64
	var totalQueuedShares math.Int = math.ZeroInt()

	// If user address is provided, get their specific requests
	if req.UserAddress != "" {
		userAddr, err := k.address.StringToBytes(req.UserAddress)
		if err != nil {
			return nil, fmt.Errorf("invalid user address: %w", err)
		}

		position, err := k.GetV2UserPosition(ctx, req.VaultType, userAddr)
		if err != nil && !errors.Is(err, collections.ErrNotFound) {
			return nil, fmt.Errorf("failed to get user position: %w", err)
		}

		if !errors.Is(err, collections.ErrNotFound) {
			userRequests = position.ExitRequests
		}
	}

	// Count total queue length and shares
	err := k.V2Collections.ExitQueue.Walk(ctx, collections.NewPrefixedPairRange[int32, uint64](int32(req.VaultType)), func(key collections.Pair[int32, uint64], exitRequestID string) (bool, error) {
		totalQueueLength++

		// Get the exit request to sum shares
		exitRequest, err := k.V2Collections.ExitRequests.Get(ctx, exitRequestID)
		if err == nil && exitRequest.Status == vaults.EXIT_STATUS_PENDING {
			totalQueuedShares = totalQueuedShares.Add(exitRequest.Shares)
		}

		return false, nil // Continue iteration
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate exit queue: %w", err)
	}

	return &vaults.QueryExitQueueResponse{
		UserRequests:      userRequests,
		TotalQueueLength:  totalQueueLength,
		TotalQueuedShares: totalQueuedShares,
	}, nil
}

// DepositPreview implements vaults.VaultV2QueryServer.
func (k vaultV2QueryServer) DepositPreview(ctx context.Context, req *vaults.QueryDepositPreviewRequest) (*vaults.QueryDepositPreviewResponse, error) {
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	if req.Amount.IsZero() || req.Amount.IsNegative() {
		return nil, fmt.Errorf("deposit amount must be positive")
	}

	// Get vault state
	vaultState, err := k.getOrCreateVaultState(ctx, req.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// Calculate shares and fees
	shareCalc, err := k.calculateDepositShares(ctx, req.VaultType, req.Amount, vaultState)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate deposit shares: %w", err)
	}

	return &vaults.QueryDepositPreviewResponse{
		EstimatedShares: shareCalc.SharesAfterFees,
		SharePrice:      shareCalc.SharePrice,
		FeeAmount:       shareCalc.FeeAmount,
	}, nil
}

// WithdrawalPreview implements vaults.VaultV2QueryServer.
func (k vaultV2QueryServer) WithdrawalPreview(ctx context.Context, req *vaults.QueryWithdrawalPreviewRequest) (*vaults.QueryWithdrawalPreviewResponse, error) {
	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	if req.Shares.IsZero() || req.Shares.IsNegative() {
		return nil, fmt.Errorf("withdrawal shares must be positive")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// Calculate tokens and fees
	withdrawCalc, err := k.calculateWithdrawalTokens(ctx, req.VaultType, req.Shares, vaultState)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate withdrawal tokens: %w", err)
	}

	return &vaults.QueryWithdrawalPreviewResponse{
		EstimatedTokens: withdrawCalc.TokensAfterFees,
		SharePrice:      withdrawCalc.SharePrice,
		FeeAmount:       withdrawCalc.FeeAmount,
	}, nil
}

// Helper functions for query server

func (k *Keeper) getEffectiveDepositFeeRate(ctx context.Context, vaultType vaults.VaultType) int32 {
	feeConfig, err := k.getFeeConfig(ctx, vaultType)
	if err != nil || !feeConfig.FeesEnabled {
		return 0
	}
	return feeConfig.DepositFeeRate
}

func (k *Keeper) getEffectiveWithdrawalFeeRate(ctx context.Context, vaultType vaults.VaultType) int32 {
	feeConfig, err := k.getFeeConfig(ctx, vaultType)
	if err != nil || !feeConfig.FeesEnabled {
		return 0
	}
	return feeConfig.WithdrawalFeeRate
}

// Migration-related query methods

// MigrationStatus returns the current migration status
func (k vaultV2QueryServer) MigrationStatus(ctx context.Context, req *vaults.QueryMigrationStatusRequest) (*vaults.QueryMigrationStatusResponse, error) {
	// Get migration state
	migrationState, err := k.GetMigrationState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration state: %w", err)
	}

	// Get migration stats
	stats, err := k.V2Collections.MigrationStats.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, fmt.Errorf("failed to get migration stats: %w", err)
	}

	if errors.Is(err, collections.ErrNotFound) {
		// Return default stats if not found
		stats = vaults.MigrationStats{
			TotalUsers:        0,
			UsersMigrated:     0,
			TotalValueLocked:  math.ZeroInt(),
			ValueMigrated:     math.ZeroInt(),
			TotalSharesIssued: math.ZeroInt(),
		}
	}

	return &vaults.QueryMigrationStatusResponse{
		State:          migrationState,
		TotalMigrated:  stats.ValueMigrated,
		TotalRemaining: stats.TotalValueLocked.Sub(stats.ValueMigrated),
		UsersMigrated:  int64(stats.UsersMigrated),
		UsersRemaining: int64(stats.TotalUsers - stats.UsersMigrated),
	}, nil
}

// UserMigrationStatus returns a user's migration status
func (k vaultV2QueryServer) UserMigrationStatus(ctx context.Context, req *vaults.QueryUserMigrationStatusRequest) (*vaults.QueryUserMigrationStatusResponse, error) {
	// Validate address
	userAddr, err := k.address.StringToBytes(req.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid user address: %w", err)
	}

	// Check if user has migrated
	migrationRecord, err := k.GetUserMigrationRecord(ctx, userAddr)
	hasMigrated := !errors.Is(err, collections.ErrNotFound)

	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, fmt.Errorf("failed to get migration record: %w", err)
	}

	response := &vaults.QueryUserMigrationStatusResponse{
		HasMigrated: hasMigrated,
	}

	if hasMigrated {
		response.MigrationRecord = &migrationRecord
	} else {
		// Get legacy positions for estimation
		legacyPositions, err := k.GetUserLegacyPositions(ctx, userAddr, vaults.FLEXIBLE)
		if err != nil {
			return nil, fmt.Errorf("failed to get legacy positions: %w", err)
		}

		// Convert slice of values to slice of pointers
		legacyPositionPtrs := make([]*vaults.Position, len(legacyPositions))
		for i := range legacyPositions {
			legacyPositionPtrs[i] = &legacyPositions[i]
		}
		response.LegacyPositions = legacyPositionPtrs

		// Calculate estimated shares if they have positions
		if len(legacyPositions) > 0 {
			migrationAmounts, err := k.CalculateMigrationAmounts(ctx, legacyPositions, math.ZeroInt())
			if err == nil {
				estimatedShares, err := k.CalculateMigrationShares(ctx, vaults.FLEXIBLE, migrationAmounts.TotalAmount)
				if err == nil {
					response.EstimatedShares = estimatedShares
				}
			}
		}
	}

	return response, nil
}

// MigrationPreview calculates what a user would receive from migration
func (k vaultV2QueryServer) MigrationPreview(ctx context.Context, req *vaults.QueryMigrationPreviewRequest) (*vaults.QueryMigrationPreviewResponse, error) {
	// Validate address
	userAddr, err := k.address.StringToBytes(req.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid user address: %w", err)
	}

	if req.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get legacy positions
	legacyPositions, err := k.GetUserLegacyPositions(ctx, userAddr, req.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get legacy positions: %w", err)
	}

	if len(legacyPositions) == 0 {
		return &vaults.QueryMigrationPreviewResponse{
			TotalValue:      math.ZeroInt(),
			PrincipalAmount: math.ZeroInt(),
			AccruedRewards:  math.ZeroInt(),
			EstimatedShares: math.ZeroInt(),
			CurrentNav:      math.LegacyZeroDec(),
			PositionCount:   0,
		}, nil
	}

	// Calculate migration amounts
	migrationAmounts, err := k.CalculateMigrationAmounts(ctx, legacyPositions, math.ZeroInt())
	if err != nil {
		return nil, fmt.Errorf("failed to calculate migration amounts: %w", err)
	}

	// Calculate estimated shares
	estimatedShares, err := k.CalculateMigrationShares(ctx, req.VaultType, migrationAmounts.TotalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate migration shares: %w", err)
	}

	// Get current NAV
	vaultState, err := k.GetV2VaultState(ctx, req.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	currentNAV := k.calculateSharePrice(vaultState)

	return &vaults.QueryMigrationPreviewResponse{
		TotalValue:      migrationAmounts.TotalAmount,
		PrincipalAmount: migrationAmounts.Principal,
		AccruedRewards:  migrationAmounts.Rewards,
		EstimatedShares: estimatedShares,
		CurrentNav:      currentNAV,
		PositionCount:   migrationAmounts.PositionCount,
	}, nil
}
