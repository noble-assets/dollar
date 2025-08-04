package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"dollar.noble.xyz/v2/types/vaults"
)

// msgServer is the server API for Msg service.
type migrationMsgServer struct {
	*Keeper
}

// NewMigrationMsgServer returns an implementation of the migration MsgServer interface
// for the provided Keeper.
func NewMigrationMsgServer(keeper *Keeper) vaults.MigrationMsgServer {
	return &migrationMsgServer{Keeper: keeper}
}

var _ vaults.MigrationMsgServer = migrationMsgServer{}

// MigratePosition implements vaults.MigrationMsgServer.
func (k migrationMsgServer) MigratePosition(ctx context.Context, msg *vaults.MsgMigratePosition) (*vaults.MsgMigratePositionResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address: %w", err)
	}

	// Check migration state
	migrationState, err := k.GetMigrationState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration state: %w", err)
	}

	if migrationState != vaults.MIGRATION_STATE_ACTIVE && migrationState != vaults.MIGRATION_STATE_CLOSING {
		return nil, fmt.Errorf("migration not active, current state: %s", migrationState.String())
	}

	// Validate vault type
	if msg.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Check if user has already migrated
	hasMigrated, err := k.HasUserMigrated(ctx, signer, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to check migration status: %w", err)
	}
	if hasMigrated {
		return nil, fmt.Errorf("user has already migrated for vault type %s", msg.VaultType.String())
	}

	// Get legacy positions
	legacyPositions, err := k.GetUserLegacyPositions(ctx, signer, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get legacy positions: %w", err)
	}

	if len(legacyPositions) == 0 {
		return nil, fmt.Errorf("no positions to migrate for vault type %s", msg.VaultType.String())
	}

	// Check migration eligibility
	if err := k.CanUserMigrate(ctx, signer, msg.VaultType); err != nil {
		return nil, fmt.Errorf("user cannot migrate: %w", err)
	}

	// Calculate migration amounts
	migrationAmounts, err := k.CalculateMigrationAmounts(ctx, legacyPositions, msg.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate migration amounts: %w", err)
	}

	// Validate minimum migration amount
	migrationConfig, err := k.GetMigrationConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration config: %w", err)
	}

	if migrationAmounts.TotalAmount.LT(migrationConfig.MinMigrationAmount) {
		return nil, fmt.Errorf("migration amount %s below minimum %s",
			migrationAmounts.TotalAmount.String(), migrationConfig.MinMigrationAmount.String())
	}

	// Check migration rate limit
	if err := k.CheckMigrationRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("migration rate limit exceeded: %w", err)
	}

	// Calculate shares to mint at current NAV
	shares, err := k.CalculateMigrationShares(ctx, msg.VaultType, migrationAmounts.TotalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate migration shares: %w", err)
	}

	// Execute migration
	migrationParams := MigrationExecutionParams{
		User:        signer,
		Positions:   legacyPositions,
		VaultType:   msg.VaultType,
		Principal:   migrationAmounts.Principal,
		Rewards:     migrationAmounts.Rewards,
		TotalShares: shares,
		ForgoYield:  msg.ForgoYield,
	}

	txHash, gasUsed, err := k.ExecuteMigration(ctx, &migrationParams)
	if err != nil {
		return nil, fmt.Errorf("failed to execute migration: %w", err)
	}

	// Increment block migration count
	if err := k.IncrementBlockMigrationCount(ctx); err != nil {
		k.logger.Error("failed to increment block migration count", "error", err)
	}

	return &vaults.MsgMigratePositionResponse{
		SharesReceived:    shares,
		PrincipalMigrated: migrationAmounts.Principal,
		RewardsMigrated:   migrationAmounts.Rewards,
		MigrationTxHash:   txHash,
		GasUsed:           gasUsed,
	}, nil
}

// EmergencyWithdrawLegacy implements vaults.MigrationMsgServer.
func (k migrationMsgServer) EmergencyWithdrawLegacy(ctx context.Context, msg *vaults.MsgEmergencyWithdrawLegacy) (*vaults.MsgEmergencyWithdrawLegacyResponse, error) {
	// Validate signer
	signer, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address: %w", err)
	}

	// Check if emergency withdrawals are allowed
	migrationState, err := k.GetMigrationState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration state: %w", err)
	}

	// Emergency withdrawals only allowed if migration is cancelled or failed
	if migrationState != vaults.MIGRATION_STATE_CANCELLED {
		return nil, fmt.Errorf("emergency withdrawals only allowed when migration is cancelled, current state: %s", migrationState.String())
	}

	// Validate vault type
	if msg.VaultType == vaults.UNSPECIFIED {
		return nil, fmt.Errorf("vault type must be specified")
	}

	// Get locked legacy positions for this user
	lockedPositions, err := k.GetUserLockedLegacyPositions(ctx, signer, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get locked legacy positions: %w", err)
	}

	if len(lockedPositions) == 0 {
		return nil, fmt.Errorf("no locked positions found for emergency withdrawal")
	}

	// Filter positions if specific indices requested
	var positionsToWithdraw []vaults.LockedLegacyPosition
	if len(msg.PositionIndices) > 0 {
		indexMap := make(map[int64]bool)
		for _, idx := range msg.PositionIndices {
			indexMap[idx] = true
		}

		for _, pos := range lockedPositions {
			if indexMap[pos.Position.Index] {
				positionsToWithdraw = append(positionsToWithdraw, pos)
			}
		}

		if len(positionsToWithdraw) == 0 {
			return nil, fmt.Errorf("none of the specified position indices found")
		}
	} else {
		positionsToWithdraw = lockedPositions
	}

	// Calculate total withdrawal amount
	totalWithdrawal := math.ZeroInt()
	for _, lockedPos := range positionsToWithdraw {
		// Calculate principal + accrued rewards
		rewards, err := k.CalculateLegacyRewards(ctx, lockedPos.Position)
		if err != nil {
			k.logger.Error("failed to calculate rewards for emergency withdrawal",
				"position", lockedPos.Position.Index, "error", err)
			rewards = math.ZeroInt() // Continue without rewards if calculation fails
		}
		totalWithdrawal = totalWithdrawal.Add(lockedPos.Position.Principal).Add(rewards)
	}

	// Execute emergency withdrawal
	if err := k.ExecuteEmergencyWithdrawal(ctx, signer, positionsToWithdraw, totalWithdrawal); err != nil {
		return nil, fmt.Errorf("failed to execute emergency withdrawal: %w", err)
	}

	// Emit emergency withdrawal event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"emergency_withdrawal_legacy",
			sdk.NewAttribute("user", sdk.AccAddress(signer).String()),
			sdk.NewAttribute("vault_type", msg.VaultType.String()),
			sdk.NewAttribute("amount_withdrawn", totalWithdrawal.String()),
			sdk.NewAttribute("positions_count", fmt.Sprintf("%d", len(positionsToWithdraw))),
		),
	)

	return &vaults.MsgEmergencyWithdrawLegacyResponse{
		AmountWithdrawn:    totalWithdrawal,
		PositionsWithdrawn: int64(len(positionsToWithdraw)),
	}, nil
}

// UpdateMigrationState implements vaults.MigrationMsgServer.
func (k migrationMsgServer) UpdateMigrationState(ctx context.Context, msg *vaults.MsgUpdateMigrationState) (*vaults.MsgUpdateMigrationStateResponse, error) {
	// Validate authority
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	// Get current state
	currentState, err := k.GetMigrationState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current migration state: %w", err)
	}

	// Validate state transition
	if err := k.ValidateMigrationStateTransition(currentState, msg.NewState); err != nil {
		return nil, fmt.Errorf("invalid state transition: %w", err)
	}

	// Execute state transition
	if err := k.SetMigrationState(ctx, msg.NewState); err != nil {
		return nil, fmt.Errorf("failed to set migration state: %w", err)
	}

	// Execute state-specific actions
	switch msg.NewState {
	case vaults.MIGRATION_STATE_ACTIVE:
		if err := k.ActivateMigration(ctx); err != nil {
			k.logger.Error("failed to activate migration", "error", err)
		}

	case vaults.MIGRATION_STATE_CLOSING:
		if err := k.AnnounceClosingPeriod(ctx); err != nil {
			k.logger.Error("failed to announce closing period", "error", err)
		}

	case vaults.MIGRATION_STATE_LOCKED:
		if err := k.LockLegacyVault(ctx); err != nil {
			k.logger.Error("failed to lock legacy vault", "error", err)
		}

	case vaults.MIGRATION_STATE_DEPRECATED:
		if err := k.DeprecateLegacyVault(ctx); err != nil {
			k.logger.Error("failed to deprecate legacy vault", "error", err)
		}

	case vaults.MIGRATION_STATE_CANCELLED:
		if err := k.CancelMigration(ctx); err != nil {
			k.logger.Error("failed to cancel migration", "error", err)
		}
	}

	// Emit state change event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"migration_state_updated",
			sdk.NewAttribute("previous_state", currentState.String()),
			sdk.NewAttribute("new_state", msg.NewState.String()),
			sdk.NewAttribute("reason", msg.Reason),
			sdk.NewAttribute("authority", msg.Authority),
		),
	)

	return &vaults.MsgUpdateMigrationStateResponse{
		PreviousState: currentState,
		NewState:      msg.NewState,
		UpdatedAt:     sdk.UnwrapSDKContext(ctx).BlockTime(),
	}, nil
}

// Helper functions

// CanUserMigrate checks if a user is eligible to migrate
func (k *Keeper) CanUserMigrate(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType) error {
	// Check if user has already migrated
	hasMigrated, err := k.HasUserMigrated(ctx, user, vaultType)
	if err != nil {
		return err
	}
	if hasMigrated {
		return fmt.Errorf("user has already migrated")
	}

	// Check for any blocking conditions
	blockReason, err := k.GetMigrationBlockReason(ctx, user, vaultType)
	if err != nil {
		return err
	}
	if blockReason != nil {
		return fmt.Errorf("migration blocked: %s", blockReason.ReasonMessage)
	}

	return nil
}

// IncrementBlockMigrationCount tracks migrations per block for rate limiting
func (k *Keeper) IncrementBlockMigrationCount(ctx context.Context) error {
	currentBlock := sdk.UnwrapSDKContext(ctx).BlockHeight()

	count, err := k.GetBlockMigrationCount(ctx, currentBlock)
	if err != nil {
		return err
	}

	return k.SetBlockMigrationCount(ctx, currentBlock, count+1)
}

// ExecuteEmergencyWithdrawal executes an emergency withdrawal from locked positions
func (k *Keeper) ExecuteEmergencyWithdrawal(ctx context.Context, user sdk.AccAddress, positions []vaults.LockedLegacyPosition, totalAmount math.Int) error {
	// Restore positions to unlocked state
	for _, lockedPos := range positions {
		// Remove from locked positions
		key := V2LockedPositionKey(user, vaults.FLEXIBLE, lockedPos.Position.Index)
		if err := k.V2Collections.LockedLegacyPositions.Remove(ctx, key); err != nil {
			return fmt.Errorf("failed to remove locked position %d: %w", lockedPos.Position.Index, err)
		}

		// Don't restore to active positions - user is withdrawing
	}

	// Transfer tokens to user
	moduleAddr := sdk.AccAddress(authtypes.NewModuleAddress("vaults"))
	if err := k.bank.SendCoins(ctx, moduleAddr, user, sdk.NewCoins(sdk.NewCoin(k.denom, totalAmount))); err != nil {
		return fmt.Errorf("failed to transfer tokens: %w", err)
	}

	return nil
}

// ValidateMigrationStateTransition validates if a state transition is allowed
func (k *Keeper) ValidateMigrationStateTransition(current, new vaults.MigrationState) error {
	validTransitions := map[vaults.MigrationState][]vaults.MigrationState{
		vaults.MIGRATION_STATE_NOT_STARTED: {
			vaults.MIGRATION_STATE_ACTIVE,
		},
		vaults.MIGRATION_STATE_ACTIVE: {
			vaults.MIGRATION_STATE_CLOSING,
			vaults.MIGRATION_STATE_CANCELLED,
		},
		vaults.MIGRATION_STATE_CLOSING: {
			vaults.MIGRATION_STATE_LOCKED,
			vaults.MIGRATION_STATE_CANCELLED,
		},
		vaults.MIGRATION_STATE_LOCKED: {
			vaults.MIGRATION_STATE_DEPRECATED,
		},
		vaults.MIGRATION_STATE_CANCELLED: {
			vaults.MIGRATION_STATE_NOT_STARTED, // Allow restart
		},
		vaults.MIGRATION_STATE_DEPRECATED: {
			// Final state - no transitions allowed
		},
	}

	allowed, exists := validTransitions[current]
	if !exists {
		return fmt.Errorf("unknown current state: %s", current.String())
	}

	for _, validNew := range allowed {
		if validNew == new {
			return nil
		}
	}

	return fmt.Errorf("invalid transition from %s to %s", current.String(), new.String())
}

// State transition action functions

func (k *Keeper) ActivateMigration(ctx context.Context) error {
	// Initialize migration statistics
	stats := vaults.MigrationStats{
		TotalUsers:             0,
		UsersMigrated:          0,
		TotalValueLocked:       math.ZeroInt(),
		ValueMigrated:          math.ZeroInt(),
		TotalSharesIssued:      math.ZeroInt(),
		LastMigrationTime:      sdk.UnwrapSDKContext(ctx).BlockTime(),
		AverageGasPerMigration: 0,
		CompletionPercentage:   0,
	}

	return k.SetMigrationStats(ctx, stats)
}

func (k *Keeper) AnnounceClosingPeriod(ctx context.Context) error {
	// Logic to announce closing period (e.g., emit events, update configs)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"migration_closing_period_announced",
			sdk.NewAttribute("timestamp", sdkCtx.BlockTime().String()),
		),
	)
	return nil
}

func (k *Keeper) LockLegacyVault(ctx context.Context) error {
	// Lock legacy vault operations (prevent new deposits)
	return k.VaultsPaused.Set(ctx, int32(vaults.LOCK))
}

func (k *Keeper) DeprecateLegacyVault(ctx context.Context) error {
	// Fully deprecate legacy vault
	return k.VaultsPaused.Set(ctx, int32(vaults.ALL))
}

func (k *Keeper) CancelMigration(ctx context.Context) error {
	// Enable emergency unlock for all locked positions
	// This would be implemented to enable unlock_enabled flag on all locked positions
	return nil
}

// Placeholder implementations for state management functions
// These would be fully implemented in state_vaults_v2.go

func (k *Keeper) GetUserLockedLegacyPositions(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType) ([]vaults.LockedLegacyPosition, error) {
	// Implementation would iterate through V2LockedLegacyPositionPrefix for user
	return nil, nil // Placeholder
}

func (k *Keeper) SetMigrationState(ctx context.Context, state vaults.MigrationState) error {
	return k.V2Collections.MigrationState.Set(ctx, int32(state))
}

func (k *Keeper) SetBlockMigrationCount(ctx context.Context, blockHeight, count int64) error {
	// Implementation would track per-block migration counts
	return nil // Placeholder
}

func (k *Keeper) SetMigrationStats(ctx context.Context, stats vaults.MigrationStats) error {
	return k.V2Collections.MigrationStats.Set(ctx, stats)
}

func (k *Keeper) GetMigrationBlockReason(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType) (*vaults.MigrationBlockReason, error) {
	// Implementation would check for blocking conditions
	return nil, nil // Placeholder
}
