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

// MigrationHandler handles the migration from legacy vaults to V2 share-based vaults
type MigrationHandler struct {
	keeper *Keeper
}

// NewMigrationHandler creates a new migration handler
func NewMigrationHandler(keeper *Keeper) *MigrationHandler {
	return &MigrationHandler{
		keeper: keeper,
	}
}

// MigrationExecutionParams contains parameters for executing a migration
type MigrationExecutionParams struct {
	User        sdk.AccAddress
	Positions   []vaults.Position
	VaultType   vaults.VaultType
	Principal   math.Int
	Rewards     math.Int
	TotalShares math.Int
	ForgoYield  bool
}

// MigrationAmounts contains calculated migration amounts
type MigrationAmounts struct {
	Principal     math.Int
	Rewards       math.Int
	TotalAmount   math.Int
	PositionCount int64
}

// MigratePosition executes a user's position migration from legacy to V2 vault
func (k *Keeper) MigratePosition(ctx context.Context, msg *vaults.MsgMigratePosition) (*vaults.MsgMigratePositionResponse, error) {
	signer, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address: %w", err)
	}

	// 1. Check migration state
	migrationState, err := k.GetMigrationState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration state: %w", err)
	}

	if migrationState != vaults.MIGRATION_STATE_ACTIVE && migrationState != vaults.MIGRATION_STATE_CLOSING {
		return nil, fmt.Errorf("migration not active, current state: %d", migrationState)
	}

	// 2. Check if user has already migrated
	hasMigrated, err := k.HasUserMigrated(ctx, signer, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to check migration status: %w", err)
	}
	if hasMigrated {
		return nil, fmt.Errorf("user has already migrated for vault type %s", msg.VaultType.String())
	}

	// 3. Get legacy positions
	legacyPositions, err := k.GetUserLegacyPositions(ctx, signer, msg.VaultType)
	if err != nil {
		return nil, fmt.Errorf("failed to get legacy positions: %w", err)
	}

	if len(legacyPositions) == 0 {
		return nil, fmt.Errorf("no positions to migrate for vault type %s", msg.VaultType.String())
	}

	// 4. Calculate migration amounts
	migrationAmounts, err := k.CalculateMigrationAmounts(ctx, legacyPositions, msg.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate migration amounts: %w", err)
	}

	// 5. Check migration rate limit
	if err := k.CheckMigrationRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("migration rate limit exceeded: %w", err)
	}

	// 6. Calculate shares to mint at current NAV
	shares, err := k.CalculateMigrationShares(ctx, msg.VaultType, migrationAmounts.TotalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate migration shares: %w", err)
	}

	// 7. Execute migration with atomic operations
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

	return &vaults.MsgMigratePositionResponse{
		SharesReceived:    shares,
		PrincipalMigrated: migrationAmounts.Principal,
		RewardsMigrated:   migrationAmounts.Rewards,
		MigrationTxHash:   txHash,
		GasUsed:           gasUsed,
	}, nil
}

// ExecuteMigration performs the atomic migration operation
func (k *Keeper) ExecuteMigration(ctx context.Context, params *MigrationExecutionParams) (string, uint64, error) {
	// Start gas tracking
	gasStart := sdk.UnwrapSDKContext(ctx).GasMeter().GasConsumed()

	// 1. Lock legacy positions first (prevent double-spend)
	migrationID := k.GenerateMigrationID(ctx, params.User)
	if err := k.LockLegacyPositions(ctx, params.User, params.Positions, migrationID); err != nil {
		return "", 0, fmt.Errorf("failed to lock legacy positions: %w", err)
	}

	// 2. Mint shares in new system
	if err := k.MintMigrationShares(ctx, params.User, params.VaultType, params.TotalShares, params.Principal, params.ForgoYield); err != nil {
		// Rollback: unlock legacy positions
		k.UnlockLegacyPositions(ctx, params.User, params.Positions)
		return "", 0, fmt.Errorf("failed to mint migration shares: %w", err)
	}

	// 3. Record migration details
	migrationRecord := vaults.UserMigrationRecord{
		MigratedAt:          sdk.UnwrapSDKContext(ctx).BlockTime(),
		FromVaultType:       params.VaultType,
		LegacyPositionCount: int64(len(params.Positions)),
		PrincipalMigrated:   params.Principal,
		RewardsMigrated:     params.Rewards,
		SharesReceived:      params.TotalShares,
		MigrationTxHash:     migrationID, // Will be updated with actual tx hash
		GasUsed:             0,           // Will be calculated at end
		YieldForgone:        params.ForgoYield,
	}

	if err := k.RecordUserMigration(ctx, params.User, migrationRecord); err != nil {
		// This is non-critical, log but don't fail
		k.logger.Error("failed to record migration", "user", params.User.String(), "error", err)
	}

	// 4. Update migration statistics
	if err := k.UpdateMigrationStats(ctx, params.VaultType, params.Principal.Add(params.Rewards), params.TotalShares); err != nil {
		// This is non-critical, log but don't fail
		k.logger.Error("failed to update migration stats", "error", err)
	}

	// 5. Emit migration event
	k.EmitMigrationEvent(ctx, params.User, params.VaultType, params.TotalShares, params.Principal, params.Rewards)

	// Calculate gas used
	gasEnd := sdk.UnwrapSDKContext(ctx).GasMeter().GasConsumed()
	gasUsed := gasEnd - gasStart

	return migrationID, gasUsed, nil
}

// CalculateMigrationAmounts calculates the total amounts to be migrated
func (k *Keeper) CalculateMigrationAmounts(ctx context.Context, positions []vaults.Position, requestedAmount math.Int) (*MigrationAmounts, error) {
	if len(positions) == 0 {
		return nil, fmt.Errorf("no positions provided")
	}

	totalPrincipal := math.ZeroInt()
	totalRewards := math.ZeroInt()
	positionCount := int64(len(positions))

	// Calculate total principal and rewards from all positions
	for _, position := range positions {
		totalPrincipal = totalPrincipal.Add(position.Principal)

		// Calculate accrued rewards for this position
		rewards, err := k.CalculateLegacyRewards(ctx, &position)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate rewards for position %d: %w", position.Index, err)
		}
		totalRewards = totalRewards.Add(rewards)
	}

	totalAmount := totalPrincipal.Add(totalRewards)

	// If specific amount requested, validate it doesn't exceed total
	if !requestedAmount.IsZero() {
		if requestedAmount.GT(totalAmount) {
			return nil, fmt.Errorf("requested amount %s exceeds available amount %s", requestedAmount.String(), totalAmount.String())
		}
		// For partial migrations, pro-rate the amounts
		ratio := math.LegacyNewDecFromInt(requestedAmount).Quo(math.LegacyNewDecFromInt(totalAmount))
		totalPrincipal = ratio.MulInt(totalPrincipal).TruncateInt()
		totalRewards = ratio.MulInt(totalRewards).TruncateInt()
		totalAmount = requestedAmount
	}

	return &MigrationAmounts{
		Principal:     totalPrincipal,
		Rewards:       totalRewards,
		TotalAmount:   totalAmount,
		PositionCount: positionCount,
	}, nil
}

// CalculateMigrationShares calculates shares to mint based on current NAV
func (k *Keeper) CalculateMigrationShares(ctx context.Context, vaultType vaults.VaultType, amount math.Int) (math.Int, error) {
	vaultState, err := k.GetV2VaultState(ctx, vaultType)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("failed to get vault state: %w", err)
	}

	// If vault is empty, first migrator gets 1:1 share ratio
	if vaultState.TotalShares.IsZero() {
		return amount, nil
	}

	// Calculate shares using current NAV: shares = amount * total_shares / nav_point
	if vaultState.NavPoint.IsZero() {
		return math.ZeroInt(), fmt.Errorf("vault NAV is zero")
	}

	shares := amount.Mul(vaultState.TotalShares).Quo(vaultState.NavPoint)

	if shares.IsZero() {
		return math.ZeroInt(), fmt.Errorf("calculated shares is zero, amount too small")
	}

	return shares, nil
}

// CalculateLegacyRewards calculates accrued rewards for a legacy position
func (k *Keeper) CalculateLegacyRewards(ctx context.Context, position *vaults.Position) (math.Int, error) {
	// Get current rewards for the position index
	reward, err := k.VaultsRewards.Get(ctx, position.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return math.ZeroInt(), nil // No rewards accrued
		}
		return math.ZeroInt(), err
	}

	// Calculate rewards based on position's share of the reward pool
	if reward.Total.IsZero() {
		return math.ZeroInt(), nil
	}

	// Rewards = (position_amount / total_in_reward_pool) * available_rewards
	positionRewards := position.Amount.Mul(reward.Rewards).Quo(reward.Total)
	return positionRewards, nil
}

// LockLegacyPositions locks legacy positions to prevent double-spend during migration
func (k *Keeper) LockLegacyPositions(ctx context.Context, user sdk.AccAddress, positions []vaults.Position, migrationID string) error {
	for _, position := range positions {
		lockedPosition := vaults.LockedLegacyPosition{
			Position:      &position,
			LockedAt:      sdk.UnwrapSDKContext(ctx).BlockTime(),
			MigratedTo:    user.Bytes(),
			MigrationId:   migrationID,
			UnlockEnabled: false,
		}

		key := V2LockedPositionKey(user, vaults.FLEXIBLE, position.Index)
		if err := k.V2Collections.LockedLegacyPositions.Set(ctx, key, lockedPosition); err != nil {
			return fmt.Errorf("failed to lock position %d: %w", position.Index, err)
		}

		// Remove from active positions
		positionKey := collections.Join3(user.Bytes(), int32(1), position.Index) // Adjust vault type
		if err := k.VaultsPositions.Remove(ctx, positionKey); err != nil {
			return fmt.Errorf("failed to remove position %d from active positions: %w", position.Index, err)
		}
	}

	return nil
}

// UnlockLegacyPositions unlocks positions (used for rollback)
func (k *Keeper) UnlockLegacyPositions(ctx context.Context, user sdk.AccAddress, positions []vaults.Position) error {
	for _, position := range positions {
		key := V2LockedPositionKey(user, vaults.FLEXIBLE, position.Index)
		if err := k.V2Collections.LockedLegacyPositions.Remove(ctx, key); err != nil {
			k.logger.Error("failed to unlock position during rollback", "position", position.Index, "error", err)
		}

		// Restore to active positions
		positionKey := collections.Join3(user.Bytes(), int32(1), position.Index)
		if err := k.VaultsPositions.Set(ctx, positionKey, position); err != nil {
			k.logger.Error("failed to restore position during rollback", "position", position.Index, "error", err)
		}
	}

	return nil
}

// MintMigrationShares mints shares in the V2 system for migrated users
func (k *Keeper) MintMigrationShares(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType, shares, principal math.Int, forgoYield bool) error {
	// Get or create user position
	key := V2VaultUserKey(vaultType, user)
	userPosition, err := k.V2Collections.UserPositions.Get(ctx, key)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return fmt.Errorf("failed to get user position: %w", err)
	}

	// Create new position or update existing
	if errors.Is(err, collections.ErrNotFound) {
		userPosition = vaults.UserPosition{
			Shares:             shares,
			PrincipalDeposited: principal,
			AvgEntryPrice:      math.LegacyOneDec(),
			FirstDeposit:       sdk.UnwrapSDKContext(ctx).BlockTime(),
			LastActivity:       sdk.UnwrapSDKContext(ctx).BlockTime(),
			ForgoYield:         forgoYield,
		}
	} else {
		// Update existing position
		oldShares := userPosition.Shares
		newShares := oldShares.Add(shares)

		// Calculate new average entry price
		if !newShares.IsZero() {
			oldValue := userPosition.AvgEntryPrice.MulInt(oldShares)
			newValue := math.LegacyOneDec().MulInt(shares) // Migration at 1:1
			totalValue := oldValue.Add(newValue)
			userPosition.AvgEntryPrice = totalValue.QuoInt(newShares)
		}

		userPosition.Shares = newShares
		userPosition.PrincipalDeposited = userPosition.PrincipalDeposited.Add(principal)
		userPosition.LastActivity = sdk.UnwrapSDKContext(ctx).BlockTime()
		userPosition.ForgoYield = forgoYield // Update yield preference
	}

	// Save user position
	if err := k.V2Collections.UserPositions.Set(ctx, key, userPosition); err != nil {
		return fmt.Errorf("failed to save user position: %w", err)
	}

	// Update vault state
	vaultState, err := k.V2Collections.VaultStates.Get(ctx, int32(vaultType))
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return fmt.Errorf("failed to get vault state: %w", err)
	}

	if errors.Is(err, collections.ErrNotFound) {
		// Initialize new vault state
		vaultState = vaults.VaultState{
			TotalShares:    shares,
			NavPoint:       shares,
			TotalPrincipal: principal,
			LastNavUpdate:  sdk.UnwrapSDKContext(ctx).BlockTime(),
		}
	} else {
		// Update existing vault state
		vaultState.TotalShares = vaultState.TotalShares.Add(shares)
		vaultState.TotalPrincipal = vaultState.TotalPrincipal.Add(principal)
		vaultState.NavPoint = vaultState.NavPoint.Add(shares) // For migration, NAV grows with shares
	}

	if err := k.V2Collections.VaultStates.Set(ctx, int32(vaultType), vaultState); err != nil {
		return fmt.Errorf("failed to update vault state: %w", err)
	}

	return nil
}

// Validation and safety functions

// CheckMigrationRateLimit checks if migration rate limit is exceeded
func (k *Keeper) CheckMigrationRateLimit(ctx context.Context) error {
	config, err := k.GetMigrationConfig(ctx)
	if err != nil {
		return err
	}

	if config.MaxMigrationPerBlock == 0 {
		return nil // No rate limit
	}

	currentBlock := sdk.UnwrapSDKContext(ctx).BlockHeight()
	migrationCount, err := k.GetBlockMigrationCount(ctx, currentBlock)
	if err != nil {
		return err
	}

	if migrationCount >= config.MaxMigrationPerBlock {
		return fmt.Errorf("migration rate limit exceeded: %d/%d", migrationCount, config.MaxMigrationPerBlock)
	}

	return nil
}

// HasUserMigrated checks if a user has already migrated for a vault type
func (k *Keeper) HasUserMigrated(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType) (bool, error) {
	_, err := k.GetUserMigrationRecord(ctx, user)
	if errors.Is(err, collections.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Helper functions for state management

// GenerateMigrationID generates a unique migration ID
func (k *Keeper) GenerateMigrationID(ctx context.Context, user sdk.AccAddress) string {
	blockHeight := sdk.UnwrapSDKContext(ctx).BlockHeight()
	timestamp := sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
	return fmt.Sprintf("mig_%s_%d_%d", user.String()[:8], blockHeight, timestamp)
}

// EmitMigrationEvent emits a migration completion event
func (k *Keeper) EmitMigrationEvent(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType, shares, principal, rewards math.Int) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"vault_migration_completed",
			sdk.NewAttribute("user", user.String()),
			sdk.NewAttribute("vault_type", vaultType.String()),
			sdk.NewAttribute("shares_received", shares.String()),
			sdk.NewAttribute("principal_migrated", principal.String()),
			sdk.NewAttribute("rewards_migrated", rewards.String()),
		),
	)
}

// State accessor functions (these would be implemented in state_vaults_v2.go)

func (k *Keeper) GetMigrationState(ctx context.Context) (vaults.MigrationState, error) {
	state, err := k.V2Collections.MigrationState.Get(ctx)
	if errors.Is(err, collections.ErrNotFound) {
		return vaults.MIGRATION_STATE_NOT_STARTED, nil
	}
	return vaults.MigrationState(state), err
}

func (k *Keeper) GetMigrationConfig(ctx context.Context) (vaults.MigrationConfig, error) {
	config, err := k.V2Collections.MigrationConfig.Get(ctx)
	if errors.Is(err, collections.ErrNotFound) {
		// Return default config
		return vaults.MigrationConfig{
			MaxMigrationPerBlock: 10,
			MinMigrationAmount:   math.NewInt(1000),
		}, nil
	}
	return config, err
}

func (k *Keeper) GetBlockMigrationCount(ctx context.Context, blockHeight int64) (int64, error) {
	count, err := k.V2Collections.BlockMigrationCounts.Get(ctx, blockHeight)
	if errors.Is(err, collections.ErrNotFound) {
		return 0, nil
	}
	return count, err
}

func (k *Keeper) GetUserLegacyPositions(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType) ([]vaults.Position, error) {
	var positions []vaults.Position

	// Iterate through all positions for this user and vault type
	err := k.VaultsPositions.Walk(ctx, collections.NewPrefixedTripleRange[[]byte, int32, int64](user.Bytes()), func(key collections.Triple[[]byte, int32, int64], position vaults.Position) (bool, error) {
		// Check if this position matches the requested vault type
		if vaults.VaultType(key.K2()) == vaultType {
			positions = append(positions, position)
		}
		return false, nil // Continue iteration
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate user positions: %w", err)
	}

	return positions, nil
}

func (k *Keeper) GetV2VaultState(ctx context.Context, vaultType vaults.VaultType) (vaults.VaultState, error) {
	state, err := k.V2Collections.VaultStates.Get(ctx, int32(vaultType))
	if errors.Is(err, collections.ErrNotFound) {
		// Return default/empty state
		return vaults.VaultState{
			TotalShares:    math.ZeroInt(),
			NavPoint:       math.ZeroInt(),
			TotalPrincipal: math.ZeroInt(),
			LastNavUpdate:  sdk.UnwrapSDKContext(ctx).BlockTime(),
		}, nil
	}
	return state, err
}

func (k *Keeper) SetV2VaultState(ctx context.Context, vaultType vaults.VaultType, state vaults.VaultState) error {
	return k.V2Collections.VaultStates.Set(ctx, int32(vaultType), state)
}

func (k *Keeper) GetV2UserPosition(ctx context.Context, vaultType vaults.VaultType, user sdk.AccAddress) (vaults.UserPosition, error) {
	key := V2VaultUserKey(vaultType, user)
	return k.V2Collections.UserPositions.Get(ctx, key)
}

func (k *Keeper) SetV2UserPosition(ctx context.Context, vaultType vaults.VaultType, user sdk.AccAddress, position vaults.UserPosition) error {
	key := V2VaultUserKey(vaultType, user)
	return k.V2Collections.UserPositions.Set(ctx, key, position)
}

func (k *Keeper) SetLockedLegacyPosition(ctx context.Context, key collections.Triple[[]byte, int32, int64], position vaults.LockedLegacyPosition) error {
	return k.V2Collections.LockedLegacyPositions.Set(ctx, key, position)
}

func (k *Keeper) RemoveLockedLegacyPosition(ctx context.Context, key collections.Triple[[]byte, int32, int64]) error {
	return k.V2Collections.LockedLegacyPositions.Remove(ctx, key)
}

func (k *Keeper) RecordUserMigration(ctx context.Context, user sdk.AccAddress, record vaults.UserMigrationRecord) error {
	return k.V2Collections.UserMigrationRecords.Set(ctx, user.Bytes(), record)
}

func (k *Keeper) GetUserMigrationRecord(ctx context.Context, user sdk.AccAddress) (vaults.UserMigrationRecord, error) {
	return k.V2Collections.UserMigrationRecords.Get(ctx, user.Bytes())
}

func (k *Keeper) UpdateMigrationStats(ctx context.Context, vaultType vaults.VaultType, migratedValue, shares math.Int) error {
	// Get current migration stats
	stats, err := k.V2Collections.MigrationStats.Get(ctx)
	if errors.Is(err, collections.ErrNotFound) {
		// Initialize new stats if not found
		stats = vaults.MigrationStats{
			TotalUsers:             0,
			UsersMigrated:          0,
			TotalValueLocked:       math.ZeroInt(),
			ValueMigrated:          math.ZeroInt(),
			TotalSharesIssued:      math.ZeroInt(),
			LastMigrationTime:      sdk.UnwrapSDKContext(ctx).BlockTime(),
			AverageGasPerMigration: 0,
			CompletionPercentage:   0,
		}
	} else if err != nil {
		return fmt.Errorf("failed to get migration stats: %w", err)
	}

	// Update stats
	stats.UsersMigrated++
	stats.ValueMigrated = stats.ValueMigrated.Add(migratedValue)
	stats.TotalSharesIssued = stats.TotalSharesIssued.Add(shares)
	stats.LastMigrationTime = sdk.UnwrapSDKContext(ctx).BlockTime()

	// Calculate completion percentage if we have total users
	if stats.TotalUsers > 0 {
		stats.CompletionPercentage = int32((stats.UsersMigrated * 10000) / stats.TotalUsers) // basis points
	}

	// Save updated stats
	return k.V2Collections.MigrationStats.Set(ctx, stats)
}
