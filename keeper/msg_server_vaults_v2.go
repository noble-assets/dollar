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
		// Update existing position
		userPosition.Shares = userPosition.Shares.Add(sharesToMint)
		userPosition.OriginalDeposit = userPosition.OriginalDeposit.Add(msg.Amount)
		userPosition.LastActivityTime = blockTime
		if msg.ReceiveYield {
			userPosition.ReceiveYield = true // User can opt into yield but not out once opted in
		}
	}

	// Save user position
	if err := k.SetV2UserPosition(ctx, msg.VaultType, signer, userPosition); err != nil {
		return nil, fmt.Errorf("failed to save user position: %w", err)
	}

	// Update vault state
	vaultState.TotalShares = vaultState.TotalShares.Add(sharesToMint)
	vaultState.TotalNav = vaultState.TotalNav.Add(msg.Amount)

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
		changeBps = int32(change.MulInt64(10000).TruncateInt64()) // Convert to basis points
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
