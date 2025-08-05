package keeper

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/v2/types/vaults"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// vaultV2MsgServer is the server API for VaultV2Msg service
type vaultV2MsgServer struct {
	*Keeper
}

// NewVaultV2MsgServer returns an implementation of the V2 vault MsgServer interface
func NewVaultV2MsgServer(keeper *Keeper) vaultsv2.MsgServer {
	return &vaultV2MsgServer{Keeper: keeper}
}

var _ vaultsv2.MsgServer = vaultV2MsgServer{}

// Deposit implements vaultsv2.MsgServer.
func (k vaultV2MsgServer) Deposit(ctx context.Context, msg *vaultsv2.MsgDeposit) (*vaultsv2.MsgDepositResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Depositor)
	if err != nil {
		return nil, fmt.Errorf("invalid depositor address: %w", err)
	}

	// Validate vault type
	if msg.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Validate deposit amount
	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return nil, fmt.Errorf("deposit amount must be positive")
	}

	// Check for potential overflow - prevent deposits that are too large
	maxSafeInt := math.NewIntFromBigInt(new(big.Int).Rsh(new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(255), nil), big.NewInt(1)), 1))
	if msg.Amount.GT(maxSafeInt) {
		return nil, fmt.Errorf("deposit amount too large, risk of overflow")
	}

	// Get or create vault state
	vaultState, err := k.getOrCreateV2VaultState(ctx, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// Check if deposits are enabled
	if !vaultState.DepositsEnabled {
		return nil, fmt.Errorf("deposits are currently disabled for vault type %s", msg.VaultType.String())
	}

	// Calculate shares to mint
	sharePrice := k.calculateV2SharePrice(vaultState)
	sharesToMint := math.LegacyNewDecFromInt(msg.Amount).Quo(sharePrice).TruncateInt()

	// Apply slippage protection
	if sharesToMint.LT(msg.MinShares) {
		return nil, fmt.Errorf("insufficient shares received: expected at least %s, got %s",
			msg.MinShares.String(), sharesToMint.String())
	}

	// Update or create user position
	userPosition, err := k.GetV2UserPosition(ctx, msg.VaultType, signer)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	blockTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	if errors.Is(err, collections.ErrNotFound) {
		// Create new position
		userPosition = &vaultsv2.UserPosition{
			Shares:             sharesToMint,
			OriginalDeposit:    msg.Amount,
			FirstDepositTime:   blockTime,
			LastActivityTime:   blockTime,
			ReceiveYield:       msg.ReceiveYield,
			SharesPendingExit:  math.ZeroInt(),
			ActiveExitRequests: 0,
		}
	} else {
		// Update existing position with overflow protection
		newShares := userPosition.Shares.Add(sharesToMint)
		newDeposit := userPosition.OriginalDeposit.Add(msg.Amount)

		// Verify no overflow occurred
		if newShares.LT(userPosition.Shares) || newDeposit.LT(userPosition.OriginalDeposit) {
			return nil, fmt.Errorf("deposit would cause integer overflow")
		}

		userPosition.Shares = newShares
		userPosition.OriginalDeposit = newDeposit
		userPosition.LastActivityTime = blockTime
		if msg.ReceiveYield {
			userPosition.ReceiveYield = true // User can opt into yield but not out once opted in
		}
	}

	// Save user position
	if err := k.SetV2UserPosition(ctx, msg.VaultType, signer, userPosition); err != nil {
		return nil, fmt.Errorf("failed to save user position: %w", err)
	}

	// Update vault state with overflow protection
	newTotalShares := vaultState.TotalShares.Add(sharesToMint)
	newTotalNav := vaultState.TotalNav.Add(msg.Amount)

	// Verify no overflow occurred
	if newTotalShares.LT(vaultState.TotalShares) || newTotalNav.LT(vaultState.TotalNav) {
		return nil, fmt.Errorf("deposit would cause vault state overflow")
	}

	vaultState.TotalShares = newTotalShares
	vaultState.TotalNav = newTotalNav

	// Recalculate share price (should be same or very close due to proportional increase)
	vaultState.SharePrice = k.calculateV2SharePrice(vaultState)
	vaultState.LastNavUpdate = blockTime

	// Count unique users (simplified - doesn't handle exact unique count)
	if errors.Is(err, collections.ErrNotFound) {
		vaultState.TotalUsers++
	}

	if err := k.SetV2VaultState(ctx, msg.VaultType, vaultState); err != nil {
		return nil, fmt.Errorf("failed to update vault state: %w", err)
	}

	// TODO: Transfer tokens from user to vault
	// TODO: Emit events

	return &vaultsv2.MsgDepositResponse{
		SharesReceived:  sharesToMint,
		AmountDeposited: msg.Amount,
		FeesPaid:        math.ZeroInt(), // TODO: Implement fees
		SharePrice:      sharePrice,
	}, nil
}

// Withdraw implements vaultsv2.MsgServer.
func (k vaultV2MsgServer) Withdraw(ctx context.Context, msg *vaultsv2.MsgWithdraw) (*vaultsv2.MsgWithdrawResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Withdrawer)
	if err != nil {
		return nil, fmt.Errorf("invalid withdrawer address: %w", err)
	}

	// Validate vault type
	if msg.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Validate withdrawal shares
	if msg.Shares.IsZero() || msg.Shares.IsNegative() {
		return nil, fmt.Errorf("withdrawal shares must be positive")
	}

	// Get vault state
	vaultState, err := k.GetV2VaultState(ctx, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	// Check if withdrawals are enabled
	if !vaultState.WithdrawalsEnabled {
		return nil, fmt.Errorf("withdrawals are currently disabled for vault type %s", msg.VaultType.String())
	}

	// Get user position
	userPosition, err := k.GetV2UserPosition(ctx, msg.VaultType, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	// Check if user has enough shares
	if msg.Shares.GT(userPosition.Shares) {
		return nil, fmt.Errorf("insufficient shares: have %s, requested %s",
			userPosition.Shares.String(), msg.Shares.String())
	}

	// Calculate withdrawal amount
	sharePrice := k.calculateV2SharePrice(vaultState)
	withdrawalAmount := sharePrice.MulInt(msg.Shares).TruncateInt()

	// Apply slippage protection
	if withdrawalAmount.LT(msg.MinAmount) {
		return nil, fmt.Errorf("insufficient amount received: expected at least %s, got %s",
			msg.MinAmount.String(), withdrawalAmount.String())
	}

	// Update user position
	userPosition.Shares = userPosition.Shares.Sub(msg.Shares)
	userPosition.LastActivityTime = sdk.UnwrapSDKContext(ctx).BlockTime()

	// If user has no shares left, we could remove the position entirely
	// For now, keep it to maintain history
	if err := k.SetV2UserPosition(ctx, msg.VaultType, signer, userPosition); err != nil {
		return nil, fmt.Errorf("failed to update user position: %w", err)
	}

	// Update vault state
	vaultState.TotalShares = vaultState.TotalShares.Sub(msg.Shares)
	vaultState.TotalNav = vaultState.TotalNav.Sub(withdrawalAmount)
	vaultState.SharePrice = k.calculateV2SharePrice(vaultState)
	vaultState.LastNavUpdate = sdk.UnwrapSDKContext(ctx).BlockTime()

	if err := k.SetV2VaultState(ctx, msg.VaultType, vaultState); err != nil {
		return nil, fmt.Errorf("failed to update vault state: %w", err)
	}

	// TODO: Transfer tokens from vault to user
	// TODO: Emit events

	return &vaultsv2.MsgWithdrawResponse{
		AmountWithdrawn: withdrawalAmount,
		SharesRedeemed:  msg.Shares,
		FeesPaid:        math.ZeroInt(), // TODO: Implement fees
		SharePrice:      sharePrice,
	}, nil
}

// RequestExit implements vaultsv2.MsgServer (placeholder for staked vaults)
func (k vaultV2MsgServer) RequestExit(ctx context.Context, msg *vaultsv2.MsgRequestExit) (*vaultsv2.MsgRequestExitResponse, error) {
	// For the simplified design, staked vaults could work differently
	// For now, return not implemented
	return nil, fmt.Errorf("exit requests not implemented in simplified design - use direct withdrawal for flexible vaults")
}

// CancelExit implements vaultsv2.MsgServer (placeholder)
func (k vaultV2MsgServer) CancelExit(ctx context.Context, msg *vaultsv2.MsgCancelExit) (*vaultsv2.MsgCancelExitResponse, error) {
	return nil, fmt.Errorf("exit requests not implemented in simplified design")
}

// SetYieldPreference implements vaultsv2.MsgServer
func (k vaultV2MsgServer) SetYieldPreference(ctx context.Context, msg *vaultsv2.MsgSetYieldPreference) (*vaultsv2.MsgSetYieldPreferenceResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.User)
	if err != nil {
		return nil, fmt.Errorf("invalid user address: %w", err)
	}

	// Get user position
	userPosition, err := k.GetV2UserPosition(ctx, msg.VaultType, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	previousPreference := userPosition.ReceiveYield
	userPosition.ReceiveYield = msg.ReceiveYield
	userPosition.LastActivityTime = sdk.UnwrapSDKContext(ctx).BlockTime()

	if err := k.SetV2UserPosition(ctx, msg.VaultType, signer, userPosition); err != nil {
		return nil, fmt.Errorf("failed to update yield preference: %w", err)
	}

	return &vaultsv2.MsgSetYieldPreferenceResponse{
		PreviousPreference: previousPreference,
		NewPreference:      msg.ReceiveYield,
	}, nil
}

// ProcessExitQueue implements vaultsv2.MsgServer (placeholder)
func (k vaultV2MsgServer) ProcessExitQueue(ctx context.Context, msg *vaultsv2.MsgProcessExitQueue) (*vaultsv2.MsgProcessExitQueueResponse, error) {
	return nil, fmt.Errorf("exit queue processing not implemented in simplified design")
}

// UpdateNAV implements vaultsv2.MsgServer
func (k vaultV2MsgServer) UpdateNAV(ctx context.Context, msg *vaultsv2.MsgUpdateNAV) (*vaultsv2.MsgUpdateNAVResponse, error) {
	// Validate authority
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	// Get current vault state
	vaultState, err := k.GetV2VaultState(ctx, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	previousNav := vaultState.TotalNav
	vaultState.TotalNav = msg.NewNav
	vaultState.SharePrice = k.calculateV2SharePrice(vaultState)
	vaultState.LastNavUpdate = sdk.UnwrapSDKContext(ctx).BlockTime()

	// Calculate change in basis points
	var changeBps int32
	if !previousNav.IsZero() {
		change := math.LegacyNewDecFromInt(msg.NewNav.Sub(previousNav)).Quo(math.LegacyNewDecFromInt(previousNav))
		changeBpsDec := change.MulInt64(10000)

		// Check bounds to prevent int64 overflow
		maxInt64 := math.LegacyNewDec(9223372036854775807)  // math.MaxInt64
		minInt64 := math.LegacyNewDec(-9223372036854775808) // math.MinInt64

		if changeBpsDec.GT(maxInt64) {
			changeBps = int32(9999) // Cap at 99.99% change
		} else if changeBpsDec.LT(minInt64) {
			changeBps = int32(-9999) // Cap at -99.99% change
		} else {
			changeBps = int32(changeBpsDec.TruncateInt64()) // Convert to basis points
		}
	}

	if err := k.SetV2VaultState(ctx, msg.VaultType, vaultState); err != nil {
		return nil, fmt.Errorf("failed to update vault state: %w", err)
	}

	return &vaultsv2.MsgUpdateNAVResponse{
		PreviousNav:   previousNav,
		NewNav:        msg.NewNav,
		ChangeBps:     changeBps,
		NewSharePrice: vaultState.SharePrice,
	}, nil
}

// UpdateVaultConfig implements vaultsv2.MsgServer (placeholder)
func (k vaultV2MsgServer) UpdateVaultConfig(ctx context.Context, msg *vaultsv2.MsgUpdateVaultConfig) (*vaultsv2.MsgUpdateVaultConfigResponse, error) {
	return nil, fmt.Errorf("vault config updates not implemented yet")
}

// UpdateParams implements vaultsv2.MsgServer (placeholder)
func (k vaultV2MsgServer) UpdateParams(ctx context.Context, msg *vaultsv2.MsgUpdateParams) (*vaultsv2.MsgUpdateParamsResponse, error) {
	return nil, fmt.Errorf("parameter updates not implemented yet")
}

// Cross-chain message handlers

// CreateCrossChainRoute implements vaultsv2.MsgServer
func (k vaultV2MsgServer) CreateCrossChainRoute(ctx context.Context, msg *vaultsv2.MsgCreateCrossChainRoute) (*vaultsv2.MsgCreateCrossChainRouteResponse, error) {
	// Validate authority
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Create route
	if err := crossChainKeeper.CreateRoute(sdk.UnwrapSDKContext(ctx), &msg.Route); err != nil {
		return nil, fmt.Errorf("failed to create cross-chain route: %w", err)
	}

	return &vaultsv2.MsgCreateCrossChainRouteResponse{
		RouteId:     msg.Route.RouteId,
		RouteConfig: fmt.Sprintf("Created route from %s to %s via %s", msg.Route.SourceChain, msg.Route.DestinationChain, msg.Route.Provider.String()),
	}, nil
}

// UpdateCrossChainRoute implements vaultsv2.MsgServer
func (k vaultV2MsgServer) UpdateCrossChainRoute(ctx context.Context, msg *vaultsv2.MsgUpdateCrossChainRoute) (*vaultsv2.MsgUpdateCrossChainRouteResponse, error) {
	// Validate authority
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Update route
	if err := crossChainKeeper.UpdateRoute(sdk.UnwrapSDKContext(ctx), msg.RouteId, &msg.Route); err != nil {
		return nil, fmt.Errorf("failed to update cross-chain route: %w", err)
	}

	return &vaultsv2.MsgUpdateCrossChainRouteResponse{
		RouteId:        msg.RouteId,
		PreviousConfig: "Previous configuration", // TODO: Get actual previous config
		NewConfig:      "New configuration",      // TODO: Get actual new config
	}, nil
}

// DisableCrossChainRoute implements vaultsv2.MsgServer
func (k vaultV2MsgServer) DisableCrossChainRoute(ctx context.Context, msg *vaultsv2.MsgDisableCrossChainRoute) (*vaultsv2.MsgDisableCrossChainRouteResponse, error) {
	// Validate authority
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Disable route
	if err := crossChainKeeper.DisableRoute(sdk.UnwrapSDKContext(ctx), msg.RouteId); err != nil {
		return nil, fmt.Errorf("failed to disable cross-chain route: %w", err)
	}

	// TODO: Count affected positions
	affectedPositions := int64(0)

	return &vaultsv2.MsgDisableCrossChainRouteResponse{
		RouteId:           msg.RouteId,
		AffectedPositions: affectedPositions,
	}, nil
}

// RemoteDeposit implements vaultsv2.MsgServer
func (k vaultV2MsgServer) RemoteDeposit(ctx context.Context, msg *vaultsv2.MsgRemoteDeposit) (*vaultsv2.MsgRemoteDepositResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Depositor)
	if err != nil {
		return nil, fmt.Errorf("invalid depositor address: %w", err)
	}

	// Validate deposit amount
	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return nil, fmt.Errorf("deposit amount must be positive")
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Initiate remote deposit
	nonce, err := crossChainKeeper.InitiateRemoteDeposit(
		sdk.UnwrapSDKContext(ctx),
		signer,
		msg.RouteId,
		msg.Amount,
		msg.RemoteAddress,
		msg.GasLimit,
		msg.GasPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate remote deposit: %w", err)
	}

	// Get the created in-flight position for response
	inFlightPos, err := crossChainKeeper.GetInFlightPosition(sdk.UnwrapSDKContext(ctx), nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to get in-flight position: %w", err)
	}

	return &vaultsv2.MsgRemoteDepositResponse{
		Nonce:              nonce,
		RouteId:            msg.RouteId,
		SharesAllocated:    msg.Amount, // TODO: Calculate actual shares
		AmountSent:         msg.Amount,
		ExpectedCompletion: inFlightPos.ExpectedCompletion,
		ProviderTracking:   inFlightPos.ProviderTracking,
	}, nil
}

// RemoteWithdraw implements vaultsv2.MsgServer
func (k vaultV2MsgServer) RemoteWithdraw(ctx context.Context, msg *vaultsv2.MsgRemoteWithdraw) (*vaultsv2.MsgRemoteWithdrawResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Withdrawer)
	if err != nil {
		return nil, fmt.Errorf("invalid withdrawer address: %w", err)
	}

	// Validate withdrawal shares
	if msg.Shares.IsZero() || msg.Shares.IsNegative() {
		return nil, fmt.Errorf("withdrawal shares must be positive")
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Initiate remote withdrawal
	nonce, err := crossChainKeeper.InitiateRemoteWithdraw(
		sdk.UnwrapSDKContext(ctx),
		signer,
		msg.RouteId,
		msg.Shares,
		msg.GasLimit,
		msg.GasPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate remote withdrawal: %w", err)
	}

	// Get the created in-flight position for response
	inFlightPos, err := crossChainKeeper.GetInFlightPosition(sdk.UnwrapSDKContext(ctx), nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to get in-flight position: %w", err)
	}

	return &vaultsv2.MsgRemoteWithdrawResponse{
		Nonce:              nonce,
		RouteId:            msg.RouteId,
		SharesWithdrawn:    msg.Shares,
		ExpectedAmount:     inFlightPos.Amount,
		ExpectedCompletion: inFlightPos.ExpectedCompletion,
		ProviderTracking:   inFlightPos.ProviderTracking,
	}, nil
}

// UpdateRemotePosition implements vaultsv2.MsgServer
func (k vaultV2MsgServer) UpdateRemotePosition(ctx context.Context, msg *vaultsv2.MsgUpdateRemotePosition) (*vaultsv2.MsgUpdateRemotePositionResponse, error) {
	// Validate relayer (for now, allow any address - in production this should be restricted)
	_, err := k.address.StringToBytes(msg.Relayer)
	if err != nil {
		return nil, fmt.Errorf("invalid relayer address: %w", err)
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Get previous position for comparison
	previousPosition, err := crossChainKeeper.GetRemotePosition(sdk.UnwrapSDKContext(ctx), msg.RouteId, msg.UserAddress)
	var previousValue math.Int
	if err != nil {
		previousValue = math.ZeroInt()
	} else {
		previousValue = previousPosition.RemoteValue
	}

	// Update remote position
	if err := crossChainKeeper.UpdateRemotePosition(
		sdk.UnwrapSDKContext(ctx),
		msg.RouteId,
		msg.UserAddress,
		msg.RemoteValue,
		msg.Confirmations,
		msg.ProviderTracking,
		msg.Status,
	); err != nil {
		return nil, fmt.Errorf("failed to update remote position: %w", err)
	}

	// Get updated position for conservative value
	updatedPosition, err := crossChainKeeper.GetRemotePosition(sdk.UnwrapSDKContext(ctx), msg.RouteId, msg.UserAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated position: %w", err)
	}

	return &vaultsv2.MsgUpdateRemotePositionResponse{
		RouteId:           msg.RouteId,
		UserAddress:       msg.UserAddress,
		PreviousValue:     previousValue,
		NewValue:          msg.RemoteValue,
		ConservativeValue: updatedPosition.ConservativeValue,
	}, nil
}

// ProcessInFlightPosition implements vaultsv2.MsgServer
func (k vaultV2MsgServer) ProcessInFlightPosition(ctx context.Context, msg *vaultsv2.MsgProcessInFlightPosition) (*vaultsv2.MsgProcessInFlightPositionResponse, error) {
	// Validate authority
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	// Get cross-chain keeper
	crossChainKeeper := k.V2Collections.GetCrossChainKeeper()
	if crossChainKeeper == nil {
		return nil, fmt.Errorf("cross-chain keeper not initialized")
	}

	// Get in-flight position before processing
	inFlightPos, err := crossChainKeeper.GetInFlightPosition(sdk.UnwrapSDKContext(ctx), msg.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to get in-flight position: %w", err)
	}

	// Process the in-flight position
	if err := crossChainKeeper.ProcessInFlightPosition(
		sdk.UnwrapSDKContext(ctx),
		msg.Nonce,
		msg.ResultStatus,
		msg.ResultAmount,
		msg.ErrorMessage,
		msg.ProviderTracking,
	); err != nil {
		return nil, fmt.Errorf("failed to process in-flight position: %w", err)
	}

	return &vaultsv2.MsgProcessInFlightPositionResponse{
		Nonce:           msg.Nonce,
		FinalStatus:     msg.ResultStatus,
		AmountProcessed: msg.ResultAmount,
		SharesAffected:  inFlightPos.Shares,
	}, nil
}

// Helper functions

// getOrCreateV2VaultState gets an existing vault state or creates a new one with defaults
func (k *Keeper) getOrCreateV2VaultState(ctx context.Context, vaultType vaults.VaultType) (*vaultsv2.VaultState, error) {
	state, err := k.GetV2VaultState(ctx, vaultType)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Create default vault state
			blockTime := sdk.UnwrapSDKContext(ctx).BlockTime()
			defaultState := &vaultsv2.VaultState{
				VaultType:              vaultType,
				TotalShares:            math.ZeroInt(),
				TotalNav:               math.ZeroInt(),
				SharePrice:             math.LegacyOneDec(),
				TotalUsers:             0,
				DepositsEnabled:        true,
				WithdrawalsEnabled:     true,
				LastNavUpdate:          blockTime,
				TotalSharesPendingExit: math.ZeroInt(),
				PendingExitRequests:    0,
			}
			if err := k.SetV2VaultState(ctx, vaultType, defaultState); err != nil {
				return nil, err
			}
			return defaultState, nil
		}
		return nil, err
	}
	return state, nil
}

// calculateV2SharePrice calculates the current share price based on NAV and total shares
func (k *Keeper) calculateV2SharePrice(vaultState *vaultsv2.VaultState) math.LegacyDec {
	if vaultState.TotalShares.IsZero() {
		return math.LegacyOneDec() // 1:1 ratio for first deposit
	}

	return math.LegacyNewDecFromInt(vaultState.TotalNav).Quo(math.LegacyNewDecFromInt(vaultState.TotalShares))
}
