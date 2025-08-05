# User Migration Transaction Handlers

## Overview

This document provides detailed implementation of transaction handlers for user-initiated migration from the legacy Noble Dollar vault system to the new share-based accounting system.

## Core Migration Messages

### Message Definitions

```proto
// proto/noble/dollar/vaults/v1/migration.proto
syntax = "proto3";
package noble.dollar.vaults.v1;

import "cosmos/msg/v1/msg.proto";
import "cosmos_proto/cosmos.proto";
import "gogoproto/gogo.proto";
import "google/protobuf/timestamp.proto";
import "noble/dollar/vaults/v1/vaults.proto";

// MsgMigratePosition allows users to migrate their legacy positions to the new system
message MsgMigratePosition {
  option (cosmos.msg.v1.signer) = "signer";
  
  string signer = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  
  // Vault type to migrate from (FLEXIBLE or STAKED)
  VaultType vault_type = 2;
  
  // Amount to migrate (0 or empty = migrate all)
  string amount = 3 [
    (cosmos_proto.scalar) = "cosmos.Int",
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = true
  ];
  
  // Minimum shares to receive (slippage protection)
  string min_shares_out = 4 [
    (cosmos_proto.scalar) = "cosmos.Int",
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];
}

message MsgMigratePositionResponse {
  // Total shares received
  string shares_received = 1;
  
  // Principal amount migrated
  string principal_migrated = 2;
  
  // Rewards included in migration
  string rewards_included = 3;
  
  // Migration transaction ID for tracking
  string migration_id = 4;
}

// MsgEmergencyWithdrawLegacy for emergency withdrawal from locked legacy positions
message MsgEmergencyWithdrawLegacy {
  option (cosmos.msg.v1.signer) = "signer";
  
  string signer = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  VaultType vault_type = 2;
}

message MsgEmergencyWithdrawLegacyResponse {
  string amount_withdrawn = 1;
}
```

## Migration State Management

### State Structures

```go
package keeper

import (
    "time"
    "cosmossdk.io/math"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// MigrationState tracks the overall migration status
type MigrationState int32

const (
    MigrationState_NOT_STARTED MigrationState = 0
    MigrationState_ACTIVE      MigrationState = 1
    MigrationState_CLOSING     MigrationState = 2
    MigrationState_LOCKED      MigrationState = 3
    MigrationState_DEPRECATED  MigrationState = 4
    MigrationState_CANCELLED   MigrationState = 5
)

// MigrationConfig holds migration parameters
type MigrationConfig struct {
    StartTime            time.Time
    ClosingTime          time.Time
    FinalDeadline        time.Time
    
    // Safety parameters
    MaxMigrationPerBlock math.Int // Rate limiting
    MinMigrationAmount   math.Int // Dust prevention
    RequireFullMigration bool     // If true, no partial migrations
}

// UserMigrationRecord tracks individual user migration
type UserMigrationRecord struct {
    MigratedAt           time.Time
    FromVaultType        VaultType
    LegacyPositionCount  int32
    
    // Amounts
    PrincipalMigrated    math.Int
    RewardsMigrated      math.Int
    SharesReceived       math.Int
    
    // Tracking
    MigrationTxHash      string
    GasUsed              uint64
}

// MigrationStats tracks overall migration progress
type MigrationStats struct {
    TotalUsers           int64
    UsersMigrated        int64
    
    TotalValueLocked     math.Int
    ValueMigrated        math.Int
    
    TotalSharesIssued    math.Int
    
    LastMigrationTime    time.Time
    AverageGasPerMigration uint64
}
```

## Migration Handler Implementation

### Main Migration Handler

```go
// msg_server_migrate.go
package keeper

import (
    "context"
    "fmt"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    
    "dollar.noble.xyz/v2/types/vaults"
)

func (k msgServer) MigratePosition(ctx context.Context, msg *vaults.MsgMigratePosition) (*vaults.MsgMigratePositionResponse, error) {
    signer, err := k.addressCodec.StringToBytes(msg.Signer)
    if err != nil {
        return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid signer address: %s", err)
    }
    
    // 1. Validate migration state
    if err := k.validateMigrationState(ctx); err != nil {
        return nil, err
    }
    
    // 2. Check rate limiting
    if err := k.checkMigrationRateLimit(ctx); err != nil {
        return nil, sdkerrors.ErrInvalidRequest.Wrap("migration rate limit exceeded")
    }
    
    // 3. Get and validate legacy positions
    positions, err := k.getLegacyPositionsForMigration(ctx, signer, msg.VaultType)
    if err != nil {
        return nil, err
    }
    
    if len(positions) == 0 {
        return nil, sdkerrors.ErrInvalidRequest.Wrap("no positions to migrate")
    }
    
    // 4. Check if user already migrated
    if k.hasUserMigrated(ctx, signer, msg.VaultType) {
        return nil, sdkerrors.ErrInvalidRequest.Wrap("position already migrated")
    }
    
    // 5. Calculate migration amounts
    migrationCalc, err := k.calculateMigrationAmounts(ctx, positions, msg.Amount)
    if err != nil {
        return nil, err
    }
    
    // 6. Calculate shares at current NAV
    shares, err := k.calculateMigrationShares(ctx, migrationCalc.TotalAmount)
    if err != nil {
        return nil, err
    }
    
    // 7. Validate minimum shares
    if shares.LT(msg.MinSharesOut) {
        return nil, sdkerrors.ErrInvalidRequest.Wrapf(
            "shares received (%s) less than minimum requested (%s)",
            shares, msg.MinSharesOut,
        )
    }
    
    // 8. Execute migration in atomic operation
    migrationID, err := k.executeMigration(ctx, executeMigrationParams{
        User:         signer,
        Positions:    positions,
        VaultType:    msg.VaultType,
        Principal:    migrationCalc.Principal,
        Rewards:      migrationCalc.Rewards,
        TotalShares:  shares,
    })
    
    if err != nil {
        return nil, sdkerrors.ErrInvalidRequest.Wrapf("migration execution failed: %s", err)
    }
    
    // 9. Emit migration event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "position_migrated",
            sdk.NewAttribute("user", msg.Signer),
            sdk.NewAttribute("vault_type", msg.VaultType.String()),
            sdk.NewAttribute("principal", migrationCalc.Principal.String()),
            sdk.NewAttribute("rewards", migrationCalc.Rewards.String()),
            sdk.NewAttribute("shares", shares.String()),
            sdk.NewAttribute("migration_id", migrationID),
        ),
    )
    
    return &vaults.MsgMigratePositionResponse{
        SharesReceived:    shares.String(),
        PrincipalMigrated: migrationCalc.Principal.String(),
        RewardsIncluded:   migrationCalc.Rewards.String(),
        MigrationId:       migrationID,
    }, nil
}

type executeMigrationParams struct {
    User        sdk.AccAddress
    Positions   []vaults.LegacyPosition
    VaultType   vaults.VaultType
    Principal   math.Int
    Rewards     math.Int
    TotalShares math.Int
}

func (k Keeper) executeMigration(ctx context.Context, params executeMigrationParams) (string, error) {
    // Generate migration ID
    migrationID := k.generateMigrationID(ctx, params.User)
    
    // 1. Lock legacy positions
    for _, pos := range params.Positions {
        if err := k.lockLegacyPosition(ctx, params.User, pos); err != nil {
            return "", fmt.Errorf("failed to lock position: %w", err)
        }
    }
    
    // 2. Create new position in V2 system
    v2Position := vaults.UserPosition{
        Shares:              params.TotalShares,
        Principal:           params.Principal,
        ForgoingYield:       false, // Default to receiving yield
        CheckpointPrincipal: params.Principal.Add(params.Rewards),
        CheckpointShares:    params.TotalShares,
    }
    
    if err := k.UserPositions.Set(ctx, params.User, v2Position); err != nil {
        // Rollback: unlock legacy positions
        k.unlockLegacyPositions(ctx, params.User, params.Positions)
        return "", fmt.Errorf("failed to create V2 position: %w", err)
    }
    
    // 3. Update vault state
    vaultState, err := k.VaultState.Get(ctx)
    if err != nil {
        return "", fmt.Errorf("failed to get vault state: %w", err)
    }
    
    vaultState.TotalShares = vaultState.TotalShares.Add(params.TotalShares)
    vaultState.NavPoint = vaultState.NavPoint.Add(params.Principal).Add(params.Rewards)
    
    if err := k.VaultState.Set(ctx, vaultState); err != nil {
        return "", fmt.Errorf("failed to update vault state: %w", err)
    }
    
    // 4. Record migration
    record := UserMigrationRecord{
        MigratedAt:          ctx.BlockTime(),
        FromVaultType:       params.VaultType,
        LegacyPositionCount: int32(len(params.Positions)),
        PrincipalMigrated:   params.Principal,
        RewardsMigrated:     params.Rewards,
        SharesReceived:      params.TotalShares,
        MigrationTxHash:     migrationID,
        GasUsed:             ctx.GasMeter().GasConsumed(),
    }
    
    if err := k.UserMigrationRecords.Set(ctx, params.User, record); err != nil {
        return "", fmt.Errorf("failed to record migration: %w", err)
    }
    
    // 5. Update migration stats
    if err := k.updateMigrationStats(ctx, params); err != nil {
        return "", fmt.Errorf("failed to update stats: %w", err)
    }
    
    return migrationID, nil
}
```

### Migration Calculation Logic

```go
// migration_calculations.go
package keeper

import (
    "context"
    "fmt"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "cosmossdk.io/math"
)

type MigrationAmounts struct {
    Principal    math.Int
    Rewards      math.Int
    TotalAmount  math.Int
    PositionCount int
}

func (k Keeper) calculateMigrationAmounts(ctx context.Context, positions []vaults.LegacyPosition, requestedAmount math.Int) (*MigrationAmounts, error) {
    totalPrincipal := math.ZeroInt()
    totalRewards := math.ZeroInt()
    
    for _, pos := range positions {
        totalPrincipal = totalPrincipal.Add(pos.Principal)
        
        // Calculate accrued rewards
        rewards, err := k.calculateLegacyRewards(ctx, pos)
        if err != nil {
            return nil, err
        }
        totalRewards = totalRewards.Add(rewards)
    }
    
    totalAmount := totalPrincipal.Add(totalRewards)
    
    // Handle partial migration if requested
    if !requestedAmount.IsNil() && !requestedAmount.IsZero() && requestedAmount.LT(totalAmount) {
        config := k.GetMigrationConfig(ctx)
        if config.RequireFullMigration {
            return nil, fmt.Errorf("partial migration not allowed")
        }
        
        // Calculate proportional amounts
        ratio := sdk.NewDecFromInt(requestedAmount).Quo(sdk.NewDecFromInt(totalAmount))
        totalPrincipal = ratio.Mul(sdk.NewDecFromInt(totalPrincipal)).TruncateInt()
        totalRewards = ratio.Mul(sdk.NewDecFromInt(totalRewards)).TruncateInt()
        totalAmount = requestedAmount
    }
    
    return &MigrationAmounts{
        Principal:     totalPrincipal,
        Rewards:       totalRewards,
        TotalAmount:   totalAmount,
        PositionCount: len(positions),
    }, nil
}

func (k Keeper) calculateMigrationShares(ctx context.Context, amount math.Int) (math.Int, error) {
    vaultState, err := k.VaultState.Get(ctx)
    if err != nil {
        return math.Int{}, err
    }
    
    // Calculate shares at current NAV
    if vaultState.TotalShares.IsZero() {
        // First migrator gets 1:1
        return amount, nil
    }
    
    // shares = amount * total_shares / nav_point
    shares := amount.Mul(vaultState.TotalShares).Quo(vaultState.NavPoint)
    
    if shares.IsZero() {
        return math.Int{}, fmt.Errorf("calculated shares would be zero")
    }
    
    return shares, nil
}

func (k Keeper) calculateLegacyRewards(ctx context.Context, pos vaults.LegacyPosition) (math.Int, error) {
    // For flexible vault, calculate based on reward index
    if pos.VaultType == vaults.FLEXIBLE {
        currentReward, exists := k.legacyKeeper.VaultsRewards.Get(ctx, pos.Index)
        if !exists {
            return math.ZeroInt(), nil
        }
        
        // Calculate proportional rewards
        if currentReward.Index > pos.Index {
            rewardRate := currentReward.Rewards.Quo(currentReward.Total)
            return pos.Amount.Mul(rewardRate).Quo(math.NewInt(1e18)), nil
        }
    }
    
    // Staked vault or no rewards
    return math.ZeroInt(), nil
}
```

### Legacy Position Management

```go
// legacy_position_manager.go
package keeper

import (
    "context"
    "fmt"
    "time"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "dollar.noble.xyz/v2/types/vaults"
)

// LockedLegacyPosition represents a migrated legacy position
type LockedLegacyPosition struct {
    Position      vaults.Position
    LockedAt      time.Time
    MigratedTo    sdk.AccAddress
    MigrationID   string
    UnlockEnabled bool // For emergency unlock
}

func (k Keeper) lockLegacyPosition(ctx context.Context, user sdk.AccAddress, pos vaults.LegacyPosition) error {
    // Create lock record
    locked := LockedLegacyPosition{
        Position:      pos.Position,
        LockedAt:      ctx.BlockTime(),
        MigratedTo:    user,
        MigrationID:   k.generateMigrationID(ctx, user),
        UnlockEnabled: false,
    }
    
    // Store locked position
    key := k.getLegacyPositionKey(user, pos.VaultType, pos.Index)
    if err := k.LockedLegacyPositions.Set(ctx, key, locked); err != nil {
        return err
    }
    
    // Remove from active legacy positions
    if err := k.legacyKeeper.RemovePosition(ctx, user, pos.VaultType, pos.Index); err != nil {
        return err
    }
    
    // Update legacy vault stats
    if err := k.legacyKeeper.DecrementVaultTotalPrincipal(ctx, pos.VaultType, pos.Principal); err != nil {
        return err
    }
    
    return nil
}

func (k Keeper) unlockLegacyPositions(ctx context.Context, user sdk.AccAddress, positions []vaults.LegacyPosition) error {
    for _, pos := range positions {
        key := k.getLegacyPositionKey(user, pos.VaultType, pos.Index)
        
        // Get locked position
        locked, err := k.LockedLegacyPositions.Get(ctx, key)
        if err != nil {
            continue // Skip if not found
        }
        
        // Restore to legacy system
        if err := k.legacyKeeper.RestorePosition(ctx, user, pos); err != nil {
            return err
        }
        
        // Remove lock
        if err := k.LockedLegacyPositions.Remove(ctx, key); err != nil {
            return err
        }
    }
    
    return nil
}

// Emergency unlock for failed migrations
func (k msgServer) EmergencyWithdrawLegacy(ctx context.Context, msg *vaults.MsgEmergencyWithdrawLegacy) (*vaults.MsgEmergencyWithdrawLegacyResponse, error) {
    signer, _ := k.addressCodec.StringToBytes(msg.Signer)
    
    // Check if emergency withdrawals are enabled
    state := k.GetMigrationState(ctx)
    if state != MigrationState_CANCELLED && state != MigrationState_DEPRECATED {
        return nil, fmt.Errorf("emergency withdrawals not enabled")
    }
    
    // Find all locked positions for user
    lockedPositions := k.getLockedPositionsForUser(ctx, signer, msg.VaultType)
    if len(lockedPositions) == 0 {
        return nil, fmt.Errorf("no locked positions found")
    }
    
    totalWithdrawn := math.ZeroInt()
    
    for _, locked := range lockedPositions {
        // Enable unlock
        locked.UnlockEnabled = true
        k.LockedLegacyPositions.Set(ctx, locked.Key, locked)
        
        // Calculate withdrawal amount (principal only, no rewards)
        amount := locked.Position.Principal
        totalWithdrawn = totalWithdrawn.Add(amount)
        
        // Transfer funds
        if err := k.bankKeeper.SendCoinsFromModuleToAccount(
            ctx, 
            types.ModuleName, 
            signer,
            sdk.NewCoins(sdk.NewCoin(k.denom, amount)),
        ); err != nil {
            return nil, err
        }
    }
    
    return &vaults.MsgEmergencyWithdrawLegacyResponse{
        AmountWithdrawn: totalWithdrawn.String(),
    }, nil
}
```

### Migration Status Queries

```go
// query_migration.go
package keeper

import (
    "context"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "dollar.noble.xyz/v2/types/vaults"
)

func (k Keeper) MigrationStatus(ctx context.Context, req *vaults.QueryMigrationStatusRequest) (*vaults.QueryMigrationStatusResponse, error) {
    state := k.GetMigrationState(ctx)
    config := k.GetMigrationConfig(ctx)
    stats := k.GetMigrationStats(ctx)
    
    return &vaults.QueryMigrationStatusResponse{
        State:          state,
        TotalMigrated:  stats.ValueMigrated.String(),
        TotalRemaining: stats.TotalValueLocked.Sub(stats.ValueMigrated).String(),
        UsersMigrated:  stats.UsersMigrated,
        UsersRemaining: stats.TotalUsers - stats.UsersMigrated,
        Config:         config,
        Stats:          stats,
    }, nil
}

func (k Keeper) UserMigrationStatus(ctx context.Context, req *vaults.QueryUserMigrationStatusRequest) (*vaults.QueryUserMigrationStatusResponse, error) {
    user, err := sdk.AccAddressFromBech32(req.Address)
    if err != nil {
        return nil, err
    }
    
    // Check if already migrated
    record, hasMigrated := k.UserMigrationRecords.Get(ctx, user)
    
    // Get legacy positions
    legacyPositions := k.getAllLegacyPositions(ctx, user)
    
    // Calculate preview if not migrated
    var preview *MigrationPreview
    if !hasMigrated && len(legacyPositions) > 0 {
        preview, err = k.calculateMigrationPreview(ctx, user, legacyPositions)
        if err != nil {
            return nil, err
        }
    }
    
    return &vaults.QueryUserMigrationStatusResponse{
        HasMigrated:     hasMigrated,
        MigrationRecord: record,
        LegacyPositions: legacyPositions,
        Preview:         preview,
        CanMigrate:      k.canUserMigrate(ctx, user),
        BlockedReason:   k.getMigrationBlockReason(ctx, user),
    }, nil
}

func (k Keeper) calculateMigrationPreview(ctx context.Context, user sdk.AccAddress, positions []vaults.LegacyPosition) (*MigrationPreview, error) {
    // Calculate total amounts
    amounts, err := k.calculateMigrationAmounts(ctx, positions, math.ZeroInt())
    if err != nil {
        return nil, err
    }
    
    // Calculate shares
    shares, err := k.calculateMigrationShares(ctx, amounts.TotalAmount)
    if err != nil {
        return nil, err
    }
    
    // Get current NAV for value calculation
    vaultState, _ := k.VaultState.Get(ctx)
    estimatedValue := shares.Mul(vaultState.NavPoint).Quo(vaultState.TotalShares)
    
    return &MigrationPreview{
        Principal:      amounts.Principal,
        Rewards:        amounts.Rewards,
        TotalAmount:    amounts.TotalAmount,
        EstimatedShares: shares,
        EstimatedValue: estimatedValue,
    }, nil
}
```

### Migration Validation and Safety

```go
// migration_validation.go
package keeper

import (
    "context"
    "fmt"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) validateMigrationState(ctx context.Context) error {
    state := k.GetMigrationState(ctx)
    
    switch state {
    case MigrationState_NOT_STARTED:
        return fmt.Errorf("migration has not started")
    case MigrationState_ACTIVE, MigrationState_CLOSING:
        return nil // Valid states for migration
    case MigrationState_LOCKED:
        return fmt.Errorf("migration period has ended")
    case MigrationState_DEPRECATED:
        return fmt.Errorf("legacy system deprecated")
    case MigrationState_CANCELLED:
        return fmt.Errorf("migration cancelled")
    default:
        return fmt.Errorf("unknown migration state")
    }
}

func (k Keeper) checkMigrationRateLimit(ctx context.Context) error {
    config := k.GetMigrationConfig(ctx)
    
    // Get current block migrations
    blockHeight := ctx.BlockHeight()
    currentBlockMigrations := k.getBlockMigrationCount(ctx, blockHeight)
    currentBlockValue := k.getBlockMigrationValue(ctx, blockHeight)
    
    // Check count limit
    if currentBlockMigrations >= 10 { // Max 10 migrations per block
        return fmt.Errorf("too many migrations in current block")
    }
    
    // Check value limit
    if currentBlockValue.GTE(config.MaxMigrationPerBlock) {
        return fmt.Errorf("migration value limit exceeded for current block")
    }
    
    return nil
}

func (k Keeper) hasUserMigrated(ctx context.Context, user sdk.AccAddress, vaultType vaults.VaultType) bool {
    record, exists := k.UserMigrationRecords.Get(ctx, user)
    if !exists {
        return false
    }
    
    // Check if specific vault type was migrated
    return record.FromVaultType == vaultType || record.FromVaultType == vaults.ALL
}

func (k Keeper) canUserMigrate(ctx context.Context, user sdk.AccAddress) bool {
    // Check migration state
    if err := k.validateMigrationState(ctx); err != nil {
        return false
    }
    
    // Check if already migrated
    if _, migrated := k.UserMigrationRecords.Get(ctx, user); migrated {
        return false
    }
    
    // Check if has positions
    positions := k.getAllLegacyPositions(ctx, user)
    return len(positions) > 0
}

func (k Keeper) getMigrationBlockReason(ctx context.Context, user sdk.AccAddress) string {
    if err := k.validateMigrationState(ctx); err != nil {
        return err.Error()
    }
    
    if _, migrated := k.UserMigrationRecords.Get(ctx, user); migrated {
        return "already migrated"
    }
    
    positions := k.getAllLegacyPositions(ctx, user)
    if len(positions) == 0 {
        return "no positions to migrate"
    }
    
    return ""
}
```

### Migration State Transitions

```go
// migration_state_manager.go
package keeper

import (
    "context"
    "fmt"
)

func (k Keeper) UpdateMigrationState(ctx context.Context, newState MigrationState) error {
    currentState := k.GetMigrationState(ctx)
    
    // Validate state transition
    validTransitions := map[MigrationState][]MigrationState{
        MigrationState_NOT_STARTED: {MigrationState_ACTIVE},
        MigrationState_ACTIVE:      {MigrationState_CLOSING, MigrationState_CANCELLED},
        MigrationState_CLOSING:     {MigrationState_LOCKED, MigrationState_CANCELLED},
        MigrationState_LOCKED:      {MigrationState_DEPRECATED},
        MigrationState_CANCELLED:   {MigrationState_NOT_STARTED},
    }
    
    allowed := false
    for _, valid := range validTransitions[currentState] {
        if valid == newState {
            allowed = true
            break
        }
    }
    
    if !allowed {
        return fmt.Errorf("invalid state transition from %v to %v", currentState, newState)
    }
    
    // Execute transition actions
    switch newState {
    case MigrationState_CLOSING:
        k.announceClosingPeriod(ctx)
    case MigrationState_LOCKED:
        k.lockLegacyVault(ctx)
    case MigrationState_DEPRECATED:
        k.deprecateLegacyVault(ctx)
    }
    
    return k.MigrationState.Set(ctx, newState)
}

func (k Keeper) announceClosingPeriod(ctx context.Context) {
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "migration_closing_soon",
            sdk.NewAttribute("deadline", k.GetMigrationConfig(ctx).ClosingTime.String()),
        ),
    )
}

func (k Keeper) lockLegacyVault(ctx context.Context) {
    // Set legacy vault to withdrawal-only mode
    k.legacyKeeper.SetPausedState(ctx, vaults.LOCK)
}

func (k Keeper) deprecateLegacyVault(ctx context.Context) {
    // Final deprecation
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "legacy_vault_deprecated",
            sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
        ),
    )
}
```

This implementation provides a complete user-initiated migration system without any incentive structures. Users migrate at their own pace based on the current NAV, with safety mechanisms to prevent double-spending and ensure data integrity throughout the process.