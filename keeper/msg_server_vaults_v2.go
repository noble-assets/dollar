package keeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/v2/types/vaults"
)

// vaultV2MsgServer is the server API for VaultV2Msg service
type vaultV2MsgServer struct {
	*Keeper
}

// NewVaultV2MsgServer returns an implementation of the V2 vault MsgServer interface
func NewVaultV2MsgServer(keeper *Keeper) vaults.VaultV2MsgServer {
	return &vaultV2MsgServer{Keeper: keeper}
}

var _ vaults.VaultV2MsgServer = vaultV2MsgServer{}

// Deposit implements vaults.VaultV2MsgServer.
func (k vaultV2MsgServer) Deposit(ctx context.Context, msg *vaults.MsgDeposit) (*vaults.MsgDepositResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address: %w", err)
	}

	// Validate vault type
	if msg.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Validate deposit amount
	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return nil, fmt.Errorf("deposit amount must be positive")
	}

	// Check if deposits are paused
	vaultState, err := k.getOrCreateVaultState(ctx, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	if vaultState.DepositsPaused {
		return nil, fmt.Errorf("deposits are currently paused for vault type %s", msg.VaultType.String())
	}

	if vaultState.CircuitBreakerActive {
		return nil, fmt.Errorf("circuit breaker is active, deposits not allowed")
	}

	// Calculate shares to mint and fees
	shareCalc, err := k.calculateDepositShares(ctx, msg.VaultType, msg.Amount, vaultState)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate deposit shares: %w", err)
	}

	// Check slippage protection
	if shareCalc.SharesAfterFees.LT(msg.MinShares) {
		return nil, fmt.Errorf("insufficient shares received: expected at least %s, got %s",
			msg.MinShares.String(), shareCalc.SharesAfterFees.String())
	}

	// Transfer tokens from user to vault
	vaultAddr := k.getVaultModuleAddress(msg.VaultType)
	if err := k.bank.SendCoins(ctx, signer, vaultAddr, sdk.NewCoins(sdk.NewCoin(k.denom, msg.Amount))); err != nil {
		return nil, fmt.Errorf("failed to transfer tokens to vault: %w", err)
	}

	// Update user position
	if err := k.updateUserPositionDeposit(ctx, signer, msg.VaultType, shareCalc, msg.ForgoYield); err != nil {
		// Rollback token transfer
		k.bank.SendCoins(ctx, vaultAddr, signer, sdk.NewCoins(sdk.NewCoin(k.denom, msg.Amount)))
		return nil, fmt.Errorf("failed to update user position: %w", err)
	}

	// Update vault state
	vaultState.TotalShares = vaultState.TotalShares.Add(shareCalc.SharesAfterFees)
	vaultState.TotalPrincipal = vaultState.TotalPrincipal.Add(msg.Amount)
	vaultState.NavPoint = vaultState.NavPoint.Add(msg.Amount) // NAV grows with deposits

	if err := k.V2Collections.VaultStates.Set(ctx, int32(msg.VaultType), vaultState); err != nil {
		return nil, fmt.Errorf("failed to update vault state: %w", err)
	}

	// Collect fees if any
	if shareCalc.FeeAmount.IsPositive() {
		if err := k.collectDepositFee(ctx, msg.VaultType, shareCalc.FeeAmount, shareCalc.FeeShares); err != nil {
			k.logger.Error("failed to collect deposit fee", "error", err)
		}
	}

	// Emit deposit event
	k.emitDepositEvent(ctx, signer, msg.VaultType, msg.Amount, shareCalc.SharesAfterFees, shareCalc.FeeAmount)

	return &vaults.MsgDepositResponse{
		SharesMinted: shareCalc.SharesAfterFees,
		SharePrice:   shareCalc.SharePrice,
		FeeCharged:   shareCalc.FeeAmount,
	}, nil
}

// Withdraw implements vaults.VaultV2MsgServer.
func (k vaultV2MsgServer) Withdraw(ctx context.Context, msg *vaults.MsgWithdraw) (*vaults.MsgWithdrawResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address: %w", err)
	}

	// Validate vault type
	if msg.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get user position
	userPosition, err := k.GetV2UserPosition(ctx, msg.VaultType, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	// Determine shares to withdraw
	sharesToWithdraw := msg.Shares
	if sharesToWithdraw.IsZero() {
		sharesToWithdraw = userPosition.Shares // Withdraw all
	}

	// Validate shares amount
	if sharesToWithdraw.IsZero() || sharesToWithdraw.IsNegative() {
		return nil, fmt.Errorf("withdrawal shares must be positive")
	}

	if sharesToWithdraw.GT(userPosition.Shares) {
		return nil, fmt.Errorf("insufficient shares: user has %s, requested %s",
			userPosition.Shares.String(), sharesToWithdraw.String())
	}

	// Check if withdrawals are paused
	vaultState, err := k.GetV2VaultState(ctx, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	if vaultState.WithdrawalsPaused {
		return nil, fmt.Errorf("withdrawals are currently paused for vault type %s", msg.VaultType.String())
	}

	// Calculate tokens to withdraw and fees
	withdrawCalc, err := k.calculateWithdrawalTokens(ctx, msg.VaultType, sharesToWithdraw, vaultState)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate withdrawal tokens: %w", err)
	}

	// Check slippage protection
	if withdrawCalc.TokensAfterFees.LT(msg.MinTokens) {
		return nil, fmt.Errorf("insufficient tokens received: expected at least %s, got %s",
			msg.MinTokens.String(), withdrawCalc.TokensAfterFees.String())
	}

	// Update user position
	if err := k.updateUserPositionWithdraw(ctx, signer, msg.VaultType, sharesToWithdraw, withdrawCalc); err != nil {
		return nil, fmt.Errorf("failed to update user position: %w", err)
	}

	// Update vault state
	vaultState.TotalShares = vaultState.TotalShares.Sub(sharesToWithdraw)
	vaultState.NavPoint = vaultState.NavPoint.Sub(withdrawCalc.TotalTokenValue)

	if err := k.V2Collections.VaultStates.Set(ctx, int32(msg.VaultType), vaultState); err != nil {
		return nil, fmt.Errorf("failed to update vault state: %w", err)
	}

	// Transfer tokens from vault to user
	vaultAddr := k.getVaultModuleAddress(msg.VaultType)
	if err := k.bank.SendCoins(ctx, vaultAddr, signer, sdk.NewCoins(sdk.NewCoin(k.denom, withdrawCalc.TokensAfterFees))); err != nil {
		return nil, fmt.Errorf("failed to transfer tokens to user: %w", err)
	}

	// Collect fees if any
	if withdrawCalc.FeeAmount.IsPositive() {
		if err := k.collectWithdrawalFee(ctx, msg.VaultType, withdrawCalc.FeeAmount); err != nil {
			k.logger.Error("failed to collect withdrawal fee", "error", err)
		}
	}

	// Emit withdrawal event
	k.emitWithdrawalEvent(ctx, signer, msg.VaultType, sharesToWithdraw, withdrawCalc.TokensAfterFees, withdrawCalc.FeeAmount)

	return &vaults.MsgWithdrawResponse{
		TokensWithdrawn: withdrawCalc.TokensAfterFees,
		SharesBurned:    sharesToWithdraw,
		SharePrice:      withdrawCalc.SharePrice,
		FeeCharged:      withdrawCalc.FeeAmount,
	}, nil
}

// RequestExit implements vaults.VaultV2MsgServer.
func (k vaultV2MsgServer) RequestExit(ctx context.Context, msg *vaults.MsgRequestExit) (*vaults.MsgRequestExitResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address: %w", err)
	}

	// Validate shares amount
	if msg.Shares.IsZero() || msg.Shares.IsNegative() {
		return nil, fmt.Errorf("exit shares must be positive")
	}

	// Get user position
	userPosition, err := k.GetV2UserPosition(ctx, msg.VaultType, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	if msg.Shares.GT(userPosition.Shares) {
		return nil, fmt.Errorf("insufficient shares for exit request")
	}

	// Generate unique exit request ID
	exitID := k.generateExitRequestID(ctx, signer, msg.VaultType)

	// Calculate expected tokens (at current price)
	vaultState, err := k.GetV2VaultState(ctx, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	sharePrice := k.calculateSharePrice(vaultState)
	expectedTokens := sharePrice.MulInt(msg.Shares).TruncateInt()

	// Get current queue position
	queuePosition, err := k.getNextQueuePosition(ctx, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue position: %w", err)
	}

	// Create exit request
	exitRequest := vaults.ExitRequest{
		RequestId:      exitID,
		UserAddress:    signer,
		VaultType:      msg.VaultType,
		Shares:         msg.Shares,
		RequestedAt:    sdk.UnwrapSDKContext(ctx).BlockTime(),
		QueuePosition:  queuePosition,
		ExpectedTokens: expectedTokens,
		Status:         vaults.EXIT_STATUS_PENDING,
	}

	// Save exit request
	if err := k.V2Collections.ExitRequests.Set(ctx, exitID, exitRequest); err != nil {
		return nil, fmt.Errorf("failed to save exit request: %w", err)
	}

	// Add to exit queue
	queueKey := V2ExitQueueKey(msg.VaultType, queuePosition)
	if err := k.V2Collections.ExitQueue.Set(ctx, queueKey, exitID); err != nil {
		return nil, fmt.Errorf("failed to add to exit queue: %w", err)
	}

	// Reserve shares (reduce user's available shares)
	userPosition.Shares = userPosition.Shares.Sub(msg.Shares)
	userPosition.LastActivity = sdk.UnwrapSDKContext(ctx).BlockTime()
	userPosition.ExitRequests = append(userPosition.ExitRequests, &exitRequest)

	if err := k.SetV2UserPosition(ctx, msg.VaultType, signer, userPosition); err != nil {
		return nil, fmt.Errorf("failed to update user position: %w", err)
	}

	// Estimate completion time (simplified)
	estimatedCompletion := sdk.UnwrapSDKContext(ctx).BlockTime().Add(time.Hour * 24) // 1 day estimate

	return &vaults.MsgRequestExitResponse{
		QueuePosition:       queuePosition,
		EstimatedCompletion: estimatedCompletion,
		ExitId:              exitID,
	}, nil
}

// CancelExit implements vaults.VaultV2MsgServer.
func (k vaultV2MsgServer) CancelExit(ctx context.Context, msg *vaults.MsgCancelExit) (*vaults.MsgCancelExitResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address: %w", err)
	}

	// Get exit request
	exitRequest, err := k.V2Collections.ExitRequests.Get(ctx, msg.ExitId)
	if err != nil {
		return nil, fmt.Errorf("exit request not found: %w", err)
	}

	// Verify user owns this exit request
	if !sdk.AccAddress(exitRequest.UserAddress).Equals(sdk.AccAddress(signer)) {
		return nil, fmt.Errorf("user does not own this exit request")
	}

	// Check if request can be cancelled
	if exitRequest.Status != vaults.EXIT_STATUS_PENDING {
		return nil, fmt.Errorf("cannot cancel exit request with status: %s", exitRequest.Status.String())
	}

	// Update exit request status
	exitRequest.Status = vaults.EXIT_STATUS_CANCELLED

	// Return shares to user
	userPosition, err := k.GetV2UserPosition(ctx, exitRequest.VaultType, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	userPosition.Shares = userPosition.Shares.Add(exitRequest.Shares)
	userPosition.LastActivity = sdk.UnwrapSDKContext(ctx).BlockTime()

	// Remove from user's exit requests
	for i, req := range userPosition.ExitRequests {
		if req.RequestId == msg.ExitId {
			userPosition.ExitRequests = append(userPosition.ExitRequests[:i], userPosition.ExitRequests[i+1:]...)
			break
		}
	}

	if err := k.SetV2UserPosition(ctx, exitRequest.VaultType, signer, userPosition); err != nil {
		return nil, fmt.Errorf("failed to update user position: %w", err)
	}

	// Update exit request
	if err := k.V2Collections.ExitRequests.Set(ctx, msg.ExitId, exitRequest); err != nil {
		return nil, fmt.Errorf("failed to update exit request: %w", err)
	}

	// Remove from exit queue
	queueKey := V2ExitQueueKey(exitRequest.VaultType, exitRequest.QueuePosition)
	if err := k.V2Collections.ExitQueue.Remove(ctx, queueKey); err != nil {
		k.logger.Error("failed to remove from exit queue", "error", err)
	}

	return &vaults.MsgCancelExitResponse{
		SharesReturned: exitRequest.Shares,
	}, nil
}

// SetYieldPreference implements vaults.VaultV2MsgServer.
func (k vaultV2MsgServer) SetYieldPreference(ctx context.Context, msg *vaults.MsgSetYieldPreference) (*vaults.MsgSetYieldPreferenceResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address: %w", err)
	}

	// Get user position
	userPosition, err := k.GetV2UserPosition(ctx, msg.VaultType, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	previousPreference := userPosition.ForgoYield
	userPosition.ForgoYield = msg.ForgoYield
	userPosition.LastActivity = sdk.UnwrapSDKContext(ctx).BlockTime()

	if err := k.SetV2UserPosition(ctx, msg.VaultType, signer, userPosition); err != nil {
		return nil, fmt.Errorf("failed to update yield preference: %w", err)
	}

	// Emit yield preference change event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"yield_preference_updated",
			sdk.NewAttribute("user", sdk.AccAddress(signer).String()),
			sdk.NewAttribute("vault_type", msg.VaultType.String()),
			sdk.NewAttribute("previous_preference", fmt.Sprintf("%t", previousPreference)),
			sdk.NewAttribute("new_preference", fmt.Sprintf("%t", msg.ForgoYield)),
		),
	)

	return &vaults.MsgSetYieldPreferenceResponse{
		PreviousPreference: previousPreference,
		NewPreference:      msg.ForgoYield,
	}, nil
}

// ProcessExitQueue implements vaults.VaultV2MsgServer.
func (k vaultV2MsgServer) ProcessExitQueue(ctx context.Context, msg *vaults.MsgProcessExitQueue) (*vaults.MsgProcessExitQueueResponse, error) {
	// Validate authority
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	var requestsProcessed uint64
	var tokensDistributed math.Int = math.ZeroInt()
	var sharesBurned math.Int = math.ZeroInt()

	// Process up to maxRequests from the queue
	for requestsProcessed < msg.MaxRequests {
		// Get next request from queue
		exitRequest, queueKey, err := k.getNextExitRequest(ctx, msg.VaultType)
		if err != nil {
			if err == collections.ErrNotFound {
				break // No more requests in queue
			}
			return nil, fmt.Errorf("failed to get next exit request: %w", err)
		}

		// Process the exit request
		tokensForUser, err := k.processExitRequest(ctx, exitRequest)
		if err != nil {
			k.logger.Error("failed to process exit request", "request_id", exitRequest.RequestId, "error", err)
			continue
		}

		// Update totals
		requestsProcessed++
		tokensDistributed = tokensDistributed.Add(tokensForUser)
		sharesBurned = sharesBurned.Add(exitRequest.Shares)

		// Remove from queue
		if err := k.V2Collections.ExitQueue.Remove(ctx, queueKey); err != nil {
			k.logger.Error("failed to remove processed request from queue", "error", err)
		}
	}

	return &vaults.MsgProcessExitQueueResponse{
		RequestsProcessed: requestsProcessed,
		TokensDistributed: tokensDistributed,
		SharesBurned:      sharesBurned,
	}, nil
}

// UpdateNAV implements vaults.VaultV2MsgServer.
func (k vaultV2MsgServer) UpdateNAV(ctx context.Context, msg *vaults.MsgUpdateNAV) (*vaults.MsgUpdateNAVResponse, error) {
	// Validate authority
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	// Get current vault state
	vaultState, err := k.GetV2VaultState(ctx, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault state: %w", err)
	}

	previousNAV := vaultState.NavPoint

	// Validate NAV change
	if err := k.validateNAVUpdate(ctx, msg.VaultType, previousNAV, msg.NewNav); err != nil {
		return nil, fmt.Errorf("invalid NAV update: %w", err)
	}

	// Calculate change percentage
	var changePercentage int32
	if !previousNAV.IsZero() {
		changeBps := msg.NewNav.Sub(previousNAV).Mul(math.NewInt(10000)).Quo(previousNAV)
		changePercentage = int32(changeBps.Int64())
	}

	// Update vault state
	vaultState.NavPoint = msg.NewNav
	vaultState.LastNavUpdate = sdk.UnwrapSDKContext(ctx).BlockTime()

	// Check for circuit breaker conditions
	navConfig, err := k.getNAVConfig(ctx, msg.VaultType)
	if err == nil && navConfig.CircuitBreakerThreshold > 0 {
		if abs(changePercentage) > navConfig.CircuitBreakerThreshold {
			vaultState.CircuitBreakerActive = true
			k.logger.Warn("Circuit breaker activated due to large NAV change",
				"vault_type", msg.VaultType.String(),
				"change_bps", changePercentage,
				"threshold", navConfig.CircuitBreakerThreshold)
		}
	}

	if err := k.V2Collections.VaultStates.Set(ctx, int32(msg.VaultType), vaultState); err != nil {
		return nil, fmt.Errorf("failed to update vault state: %w", err)
	}

	// Record NAV update in history
	navUpdate := vaults.NAVUpdate{
		PreviousNav: previousNAV,
		NewNav:      msg.NewNav,
		Timestamp:   sdk.UnwrapSDKContext(ctx).BlockTime(),
		BlockHeight: sdk.UnwrapSDKContext(ctx).BlockHeight(),
		Reason:      msg.Reason,
	}

	historyKey := V2NAVHistoryKey(msg.VaultType, sdk.UnwrapSDKContext(ctx).BlockTime().Unix())
	if err := k.V2Collections.NAVUpdates.Set(ctx, historyKey, navUpdate); err != nil {
		k.logger.Error("failed to record NAV update in history", "error", err)
	}

	// Emit NAV update event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"nav_updated",
			sdk.NewAttribute("vault_type", msg.VaultType.String()),
			sdk.NewAttribute("previous_nav", previousNAV.String()),
			sdk.NewAttribute("new_nav", msg.NewNav.String()),
			sdk.NewAttribute("change_bps", fmt.Sprintf("%d", changePercentage)),
			sdk.NewAttribute("reason", msg.Reason),
			sdk.NewAttribute("authority", msg.Authority),
		),
	)

	return &vaults.MsgUpdateNAVResponse{
		PreviousNav:      previousNAV,
		NewNav:           msg.NewNav,
		ChangePercentage: changePercentage,
	}, nil
}

// Helper functions

type ShareCalculation struct {
	SharesBeforeFees math.Int
	SharesAfterFees  math.Int
	FeeAmount        math.Int
	FeeShares        math.Int
	SharePrice       math.LegacyDec
}

type WithdrawCalculation struct {
	TotalTokenValue math.Int
	TokensAfterFees math.Int
	FeeAmount       math.Int
	SharePrice      math.LegacyDec
}

func (k *Keeper) calculateDepositShares(ctx context.Context, vaultType vaults.VaultType, amount math.Int, vaultState vaults.VaultState) (*ShareCalculation, error) {
	sharePrice := k.calculateSharePrice(vaultState)

	// Calculate shares before fees
	sharesBeforeFees := sharePrice.Quo(math.LegacyOneDec()).MulInt(amount).TruncateInt()

	// Calculate deposit fee
	feeConfig, err := k.getFeeConfig(ctx, vaultType)
	if err != nil {
		// Use zero fees if config not found
		return &ShareCalculation{
			SharesBeforeFees: sharesBeforeFees,
			SharesAfterFees:  sharesBeforeFees,
			FeeAmount:        math.ZeroInt(),
			FeeShares:        math.ZeroInt(),
			SharePrice:       sharePrice,
		}, nil
	}

	feeAmount := amount.Mul(math.NewInt(int64(feeConfig.DepositFeeRate))).Quo(math.NewInt(10000))
	feeShares := sharePrice.Quo(math.LegacyOneDec()).MulInt(feeAmount).TruncateInt()

	return &ShareCalculation{
		SharesBeforeFees: sharesBeforeFees,
		SharesAfterFees:  sharesBeforeFees.Sub(feeShares),
		FeeAmount:        feeAmount,
		FeeShares:        feeShares,
		SharePrice:       sharePrice,
	}, nil
}

func (k *Keeper) calculateWithdrawalTokens(ctx context.Context, vaultType vaults.VaultType, shares math.Int, vaultState vaults.VaultState) (*WithdrawCalculation, error) {
	sharePrice := k.calculateSharePrice(vaultState)

	// Calculate token value of shares
	totalTokenValue := sharePrice.MulInt(shares).TruncateInt()

	// Calculate withdrawal fee
	feeConfig, err := k.getFeeConfig(ctx, vaultType)
	if err != nil {
		// Use zero fees if config not found
		return &WithdrawCalculation{
			TotalTokenValue: totalTokenValue,
			TokensAfterFees: totalTokenValue,
			FeeAmount:       math.ZeroInt(),
			SharePrice:      sharePrice,
		}, nil
	}

	feeAmount := totalTokenValue.Mul(math.NewInt(int64(feeConfig.WithdrawalFeeRate))).Quo(math.NewInt(10000))

	return &WithdrawCalculation{
		TotalTokenValue: totalTokenValue,
		TokensAfterFees: totalTokenValue.Sub(feeAmount),
		FeeAmount:       feeAmount,
		SharePrice:      sharePrice,
	}, nil
}

func (k *Keeper) calculateSharePrice(vaultState vaults.VaultState) math.LegacyDec {
	if vaultState.TotalShares.IsZero() {
		return math.LegacyOneDec() // 1:1 ratio for first deposit
	}

	return math.LegacyNewDecFromInt(vaultState.NavPoint).Quo(math.LegacyNewDecFromInt(vaultState.TotalShares))
}

func (k *Keeper) getOrCreateVaultState(ctx context.Context, vaultType vaults.VaultType) (vaults.VaultState, error) {
	state, err := k.V2Collections.VaultStates.Get(ctx, int32(vaultType))
	if err == nil {
		return state, nil
	}

	if errors.Is(err, collections.ErrNotFound) {
		return vaults.VaultState{}, err
	}

	// Create new vault state
	newState := vaults.VaultState{
		TotalShares:          math.ZeroInt(),
		NavPoint:             math.ZeroInt(),
		TotalPrincipal:       math.ZeroInt(),
		AccumulatedYield:     math.ZeroInt(),
		LastNavUpdate:        sdk.UnwrapSDKContext(ctx).BlockTime(),
		CircuitBreakerActive: false,
		DepositsPaused:       false,
		WithdrawalsPaused:    false,
	}

	return newState, nil
}

func (k *Keeper) getVaultModuleAddress(vaultType vaults.VaultType) sdk.AccAddress {
	switch vaultType {
	case vaults.STAKED:
		return vaults.StakedVaultAddress
	case vaults.FLEXIBLE:
		return vaults.FlexibleVaultAddress
	default:
		return vaults.FlexibleVaultAddress // fallback
	}
}

func (k *Keeper) getFeeConfig(ctx context.Context, vaultType vaults.VaultType) (vaults.FeeConfig, error) {
	return k.V2Collections.FeeConfigs.Get(ctx, int32(vaultType))
}

func (k *Keeper) getNAVConfig(ctx context.Context, vaultType vaults.VaultType) (vaults.NAVConfig, error) {
	return k.V2Collections.NAVConfigs.Get(ctx, int32(vaultType))
}

func (k *Keeper) updateUserPositionDeposit(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType, shareCalc *ShareCalculation, forgoYield bool) error {
	key := V2VaultUserKey(vaultType, user)
	userPosition, err := k.V2Collections.UserPositions.Get(ctx, key)

	if errors.Is(err, collections.ErrNotFound) {
		// Create new position
		userPosition = vaults.UserPosition{
			Shares:             shareCalc.SharesAfterFees,
			PrincipalDeposited: shareCalc.SharesAfterFees, // For deposits, shares equal principal initially
			AvgEntryPrice:      shareCalc.SharePrice,
			FirstDeposit:       sdk.UnwrapSDKContext(ctx).BlockTime(),
			LastActivity:       sdk.UnwrapSDKContext(ctx).BlockTime(),
			ForgoYield:         forgoYield,
		}
	} else if err != nil {
		return fmt.Errorf("failed to get user position: %w", err)
	} else {
		// Update existing position
		oldShares := userPosition.Shares
		newShares := oldShares.Add(shareCalc.SharesAfterFees)

		// Calculate new average entry price
		if !newShares.IsZero() {
			oldValue := userPosition.AvgEntryPrice.MulInt(oldShares)
			newValue := shareCalc.SharePrice.MulInt(shareCalc.SharesAfterFees)
			totalValue := oldValue.Add(newValue)
			userPosition.AvgEntryPrice = totalValue.QuoInt(newShares)
		}

		userPosition.Shares = newShares
		userPosition.PrincipalDeposited = userPosition.PrincipalDeposited.Add(shareCalc.SharesAfterFees)
		userPosition.LastActivity = sdk.UnwrapSDKContext(ctx).BlockTime()
		userPosition.ForgoYield = forgoYield
	}

	return k.V2Collections.UserPositions.Set(ctx, key, userPosition)
}

func (k *Keeper) updateUserPositionWithdraw(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType, sharesToWithdraw math.Int, withdrawCalc *WithdrawCalculation) error {
	key := V2VaultUserKey(vaultType, user)
	userPosition, err := k.V2Collections.UserPositions.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to get user position: %w", err)
	}

	// Update shares
	userPosition.Shares = userPosition.Shares.Sub(sharesToWithdraw)
	userPosition.LastActivity = sdk.UnwrapSDKContext(ctx).BlockTime()

	// Proportionally reduce principal deposited
	if !userPosition.Shares.IsZero() {
		principalReduction := userPosition.PrincipalDeposited.Mul(sharesToWithdraw).Quo(userPosition.Shares.Add(sharesToWithdraw))
		userPosition.PrincipalDeposited = userPosition.PrincipalDeposited.Sub(principalReduction)
	} else {
		// Full withdrawal
		userPosition.PrincipalDeposited = math.ZeroInt()
	}

	return k.V2Collections.UserPositions.Set(ctx, key, userPosition)
}

func (k *Keeper) collectDepositFee(ctx context.Context, vaultType vaults.VaultType, feeAmount, feeShares math.Int) error {
	// Get fee config to determine recipient
	feeConfig, err := k.getFeeConfig(ctx, vaultType)
	if err != nil {
		return fmt.Errorf("failed to get fee config: %w", err)
	}

	if !feeConfig.FeesEnabled {
		return nil
	}

	// Record fee collection
	feeCollection := vaults.FeeCollection{
		Timestamp:     sdk.UnwrapSDKContext(ctx).BlockTime(),
		TotalAmount:   feeAmount,
		SharesDiluted: feeShares,
		Method:        vaults.COLLECTION_SHARE_DILUTION,
		BlockHeight:   sdk.UnwrapSDKContext(ctx).BlockHeight(),
		Breakdown: []*vaults.FeeTypeBreakdown{
			{
				FeeType:       vaults.FEE_TYPE_DEPOSIT,
				Amount:        feeAmount,
				SharesDiluted: feeShares,
				RateApplied:   feeConfig.DepositFeeRate,
			},
		},
	}

	// Save fee collection record
	collectionKey := V2FeeCollectionKey(vaultType, sdk.UnwrapSDKContext(ctx).BlockTime().Unix())
	return k.V2Collections.FeeCollections.Set(ctx, collectionKey, feeCollection)
}

func (k *Keeper) collectWithdrawalFee(ctx context.Context, vaultType vaults.VaultType, feeAmount math.Int) error {
	// Get fee config
	feeConfig, err := k.getFeeConfig(ctx, vaultType)
	if err != nil {
		return fmt.Errorf("failed to get fee config: %w", err)
	}

	if !feeConfig.FeesEnabled {
		return nil
	}

	// Record fee collection
	feeCollection := vaults.FeeCollection{
		Timestamp:     sdk.UnwrapSDKContext(ctx).BlockTime(),
		TotalAmount:   feeAmount,
		SharesDiluted: math.ZeroInt(), // Withdrawal fees don't dilute shares
		Method:        vaults.COLLECTION_TOKEN_DIRECT,
		BlockHeight:   sdk.UnwrapSDKContext(ctx).BlockHeight(),
		Breakdown: []*vaults.FeeTypeBreakdown{
			{
				FeeType:       vaults.FEE_TYPE_WITHDRAWAL,
				Amount:        feeAmount,
				SharesDiluted: math.ZeroInt(),
				RateApplied:   feeConfig.WithdrawalFeeRate,
			},
		},
	}

	// Save fee collection record
	collectionKey := V2FeeCollectionKey(vaultType, sdk.UnwrapSDKContext(ctx).BlockTime().Unix())
	return k.V2Collections.FeeCollections.Set(ctx, collectionKey, feeCollection)
}

func (k *Keeper) emitDepositEvent(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType, amount, shares, fee math.Int) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"vault_deposit",
			sdk.NewAttribute("user", user.String()),
			sdk.NewAttribute("vault_type", vaultType.String()),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("shares_minted", shares.String()),
			sdk.NewAttribute("fee_charged", fee.String()),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", sdkCtx.BlockHeight())),
		),
	)
}

func (k *Keeper) emitWithdrawalEvent(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType, shares, tokens, fee math.Int) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"vault_withdrawal",
			sdk.NewAttribute("user", user.String()),
			sdk.NewAttribute("vault_type", vaultType.String()),
			sdk.NewAttribute("shares_burned", shares.String()),
			sdk.NewAttribute("tokens_withdrawn", tokens.String()),
			sdk.NewAttribute("fee_charged", fee.String()),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", sdkCtx.BlockHeight())),
		),
	)
}

func (k *Keeper) generateExitRequestID(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType) string {
	blockHeight := sdk.UnwrapSDKContext(ctx).BlockHeight()
	timestamp := sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
	return fmt.Sprintf("exit_%s_%s_%d_%d", vaultType.String(), user.String()[:8], blockHeight, timestamp)
}

func (k *Keeper) getNextQueuePosition(ctx context.Context, vaultType vaults.VaultType) (uint64, error) {
	// Find the highest queue position for this vault type
	var maxPosition uint64 = 0

	err := k.V2Collections.ExitQueue.Walk(ctx, collections.NewPrefixedPairRange[int32, uint64](int32(vaultType)), func(key collections.Pair[int32, uint64], value string) (bool, error) {
		if key.K2() > maxPosition {
			maxPosition = key.K2()
		}
		return false, nil // Continue iteration
	})

	if err != nil {
		return 0, fmt.Errorf("failed to find max queue position: %w", err)
	}

	return maxPosition + 1, nil
}

func (k *Keeper) getNextExitRequest(ctx context.Context, vaultType vaults.VaultType) (vaults.ExitRequest, collections.Pair[int32, uint64], error) {
	// Find the lowest queue position (first in queue) for this vault type
	var minPosition uint64 = ^uint64(0) // Max uint64
	var minKey collections.Pair[int32, uint64]
	var exitRequestID string

	err := k.V2Collections.ExitQueue.Walk(ctx, collections.NewPrefixedPairRange[int32, uint64](int32(vaultType)), func(key collections.Pair[int32, uint64], value string) (bool, error) {
		if key.K2() < minPosition {
			minPosition = key.K2()
			minKey = key
			exitRequestID = value
		}
		return false, nil // Continue iteration
	})

	if err != nil {
		return vaults.ExitRequest{}, collections.Pair[int32, uint64]{}, fmt.Errorf("failed to find next exit request: %w", err)
	}

	if exitRequestID == "" {
		return vaults.ExitRequest{}, collections.Pair[int32, uint64]{}, collections.ErrNotFound
	}

	// Get the actual exit request
	exitRequest, err := k.V2Collections.ExitRequests.Get(ctx, exitRequestID)
	if err != nil {
		return vaults.ExitRequest{}, collections.Pair[int32, uint64]{}, fmt.Errorf("failed to get exit request: %w", err)
	}

	return exitRequest, minKey, nil
}

func (k *Keeper) processExitRequest(ctx context.Context, exitRequest vaults.ExitRequest) (math.Int, error) {
	// Get current vault state to calculate token value
	vaultState, err := k.GetV2VaultState(ctx, exitRequest.VaultType)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("failed to get vault state: %w", err)
	}

	// Calculate current token value of shares
	sharePrice := k.calculateSharePrice(vaultState)
	tokensToDistribute := sharePrice.MulInt(exitRequest.Shares).TruncateInt()

	// Update exit request status
	exitRequest.Status = vaults.EXIT_STATUS_COMPLETED
	if err := k.V2Collections.ExitRequests.Set(ctx, exitRequest.RequestId, exitRequest); err != nil {
		return math.ZeroInt(), fmt.Errorf("failed to update exit request: %w", err)
	}

	// Transfer tokens to user
	user := sdk.AccAddress(exitRequest.UserAddress)
	vaultAddr := k.getVaultModuleAddress(exitRequest.VaultType)
	if err := k.bank.SendCoins(ctx, vaultAddr, user, sdk.NewCoins(sdk.NewCoin(k.denom, tokensToDistribute))); err != nil {
		return math.ZeroInt(), fmt.Errorf("failed to transfer tokens to user: %w", err)
	}

	// Update vault state
	vaultState.TotalShares = vaultState.TotalShares.Sub(exitRequest.Shares)
	vaultState.NavPoint = vaultState.NavPoint.Sub(tokensToDistribute)

	if err := k.V2Collections.VaultStates.Set(ctx, int32(exitRequest.VaultType), vaultState); err != nil {
		return math.ZeroInt(), fmt.Errorf("failed to update vault state: %w", err)
	}

	return tokensToDistribute, nil
}

func (k *Keeper) validateNAVUpdate(ctx context.Context, vaultType vaults.VaultType, previousNAV, newNAV math.Int) error {
	// Check minimum time between updates
	navConfig, err := k.getNAVConfig(ctx, vaultType)
	if err == nil && navConfig.MinNavUpdateInterval > 0 {
		lastUpdate, err := k.V2Collections.LastNAVUpdate.Get(ctx, int32(vaultType))
		if err == nil {
			timeSinceUpdate := sdk.UnwrapSDKContext(ctx).BlockTime().Unix() - lastUpdate
			if timeSinceUpdate < navConfig.MinNavUpdateInterval {
				return fmt.Errorf("NAV update too soon: %d seconds since last update, minimum %d required",
					timeSinceUpdate, navConfig.MinNavUpdateInterval)
			}
		}
	}

	// Check maximum deviation
	if err == nil && navConfig.MaxNavDeviation > 0 && !previousNAV.IsZero() {
		changeBps := newNAV.Sub(previousNAV).Mul(math.NewInt(10000)).Quo(previousNAV)
		if abs(int32(changeBps.Int64())) > navConfig.MaxNavDeviation {
			return fmt.Errorf("NAV change exceeds maximum deviation: %d basis points, maximum %d allowed",
				changeBps.Int64(), navConfig.MaxNavDeviation)
		}
	}

	// NAV cannot be negative
	if newNAV.IsNegative() {
		return fmt.Errorf("NAV cannot be negative")
	}

	return nil
}

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
