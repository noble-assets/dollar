# Noble Dollar Vault Upgrade: Implementation & Migration Plan

## Overview

This document outlines the implementation plan for upgrading the Noble Dollar vaults from the current lock/unlock mechanism to a sophisticated share-based accounting system with cross-chain support, NAV bands, and yield redistribution capabilities. The migration requires explicit user action through signed transactions.

### Key Changes

1. **Share-based accounting** replacing direct balance tracking
2. **NAV band pricing** for fair deposit/withdrawal pricing
3. **Cross-chain position tracking** with conservative valuation
4. **Yield forgoing mechanism** with redistribution
5. **Fee management** through share dilution
6. **Emergency controls** and loss handling
7. **User-initiated migration** requiring explicit consent

## Migration Architecture

### Dual System Operation

During the migration period, both legacy and new vault systems will operate in parallel:

```
Legacy Vault (V1)          New Vault (V2)
    │                          │
    ├─ Lock/Unlock             ├─ Deposit/Withdraw
    ├─ Position tracking       ├─ Share-based accounting
    ├─ Reward distribution     ├─ NAV band pricing
    │                          ├─ Cross-chain support
    │                          │
    └──── Migration Tx ────────┘
```

### Migration States

```go
enum MigrationState {
    NOT_STARTED     // Pre-deployment
    ACTIVE          // Users can migrate
    CLOSING         // Grace period warning
    LOCKED          // No new migrations, withdrawals only
    DEPRECATED      // Legacy vault deprecated
}
```

## Implementation Phases

### Phase 1: Core Data Models (Weeks 1-2)

#### 1.1 Protocol Buffer Updates

**New Files:**
- `proto/noble/dollar/vaults/v1/nav.proto` - NAV band structures
- `proto/noble/dollar/vaults/v1/cross_chain.proto` - Remote position tracking
- `proto/noble/dollar/vaults/v1/fees.proto` - Fee parameters
- `proto/noble/dollar/vaults/v1/migration.proto` - Migration messages

**Migration Message Structure:**
```proto
message MsgMigratePosition {
  string signer = 1;
  VaultType vault_type = 2;  // FLEXIBLE or STAKED
  string amount = 3;         // Amount to migrate (0 = all)
}

message MsgMigratePositionResponse {
  string shares_received = 1;
  string principal_migrated = 2;
  string rewards_included = 3;
}
```

#### 1.2 State Structure Updates

**Keeper State Additions (V2 runs alongside V1):**
```go
// New V2 state
- VaultState (single item)
- UserPositions (simple map[address]Position)
- RemotePositions (map[routeId]RemotePosition)
- InFlightPositions (map[nonce]InFlight)
- ExitQueue (map[queueId]ExitRequest)
- LastFeeAccrual (timestamp)

// Migration tracking
- MigrationState (current state)
- MigrationStats (tracking progress)
- UserMigrationStatus (map[address]MigrationRecord)
```

### Phase 2: Migration Infrastructure (Weeks 3-4)

#### 2.1 Migration Handler

**New file: `keeper/migration_handler.go`**
```go
type MigrationHandler struct {
    legacyKeeper *LegacyKeeper
    newKeeper    *Keeper
}

type MigrationRecord struct {
    MigratedAt        time.Time
    LegacyPositions   []LegacyPosition
    SharesReceived    math.Int
    PrincipalMigrated math.Int
}
```

#### 2.2 Migration Process Flow

```
1. User initiates migration transaction
2. System validates legacy positions
3. Calculate migration values including rewards
4. Lock legacy position (prevent double-spend)
5. Mint shares in new system
6. Record migration details
7. Emit migration event
```

### Phase 3: User Migration Implementation (Weeks 5-6)

#### 3.1 Migration Message Handler

**New file: `keeper/msg_server_migrate.go`**
```go
func (k msgServer) MigratePosition(ctx context.Context, msg *MsgMigratePosition) (*MsgMigratePositionResponse, error) {
    signer, _ := k.addressCodec.StringToBytes(msg.Signer)
    
    // 1. Check migration state
    state := k.GetMigrationState(ctx)
    if state != ACTIVE && state != CLOSING {
        return nil, ErrMigrationNotActive
    }
    
    // 2. Get legacy positions
    legacyPositions, err := k.legacyKeeper.GetUserPositions(ctx, signer, msg.VaultType)
    if err != nil {
        return nil, err
    }
    
    if len(legacyPositions) == 0 {
        return nil, ErrNoPositionsToMigrate
    }
    
    // 3. Calculate migration amount
    totalAmount, totalPrincipal, rewards := k.calculateMigrationAmounts(ctx, legacyPositions, msg.Amount)
    
    // 4. Calculate shares to mint at current NAV
    shares := k.calculateMigrationShares(ctx, totalAmount)
    
    // 5. Lock legacy positions
    if err := k.lockLegacyPositions(ctx, signer, legacyPositions); err != nil {
        return nil, err
    }
    
    // 6. Mint shares in new system
    if err := k.mintMigrationShares(ctx, signer, shares, totalPrincipal); err != nil {
        // Rollback legacy lock on failure
        k.unlockLegacyPositions(ctx, signer, legacyPositions)
        return nil, err
    }
    
    // 7. Record migration
    k.recordMigration(ctx, signer, MigrationRecord{
        MigratedAt:        ctx.BlockTime(),
        LegacyPositions:   legacyPositions,
        SharesReceived:    shares,
        PrincipalMigrated: totalPrincipal,
    })
    
    // 8. Update stats
    k.updateMigrationStats(ctx, totalAmount, shares)
    
    return &MsgMigratePositionResponse{
        SharesReceived:    shares.String(),
        PrincipalMigrated: totalPrincipal.String(),
        RewardsIncluded:   rewards.String(),
    }, nil
}
```

#### 3.2 Share Calculation

```go
func (k Keeper) calculateMigrationShares(ctx context.Context, amount math.Int) math.Int {
    vaultState, _ := k.VaultState.Get(ctx)
    
    // Calculate shares at current NAV
    if vaultState.TotalShares.IsZero() {
        // First migrator gets 1:1
        return amount
    } else {
        // Use current NAV for fair pricing
        // shares = amount * total_shares / nav_point
        return amount.Mul(vaultState.TotalShares).Quo(vaultState.NavPoint)
    }
}
```

### Phase 4: Legacy Position Management (Weeks 7-8)

#### 4.1 Legacy Position Locking

```go
type LockedPosition struct {
    OriginalPosition Position
    LockedAt        time.Time
    MigratedTo      sdk.AccAddress
}

func (k LegacyKeeper) LockPosition(ctx context.Context, user sdk.AccAddress, positions []Position) error {
    for _, pos := range positions {
        // Mark position as locked
        locked := LockedPosition{
            OriginalPosition: pos,
            LockedAt:        ctx.BlockTime(),
            MigratedTo:      user,
        }
        
        key := collections.Join3(user.Bytes(), int32(pos.VaultType), pos.Index)
        if err := k.LockedPositions.Set(ctx, key, locked); err != nil {
            return err
        }
        
        // Remove from active positions
        if err := k.VaultsPositions.Remove(ctx, key); err != nil {
            return err
        }
    }
    
    return nil
}
```

#### 4.2 Emergency Unlock Mechanism

```go
func (k msgServer) EmergencyUnlock(ctx context.Context, msg *MsgEmergencyUnlock) (*MsgEmergencyUnlockResponse, error) {
    // Only allowed if migration fails or is cancelled
    state := k.GetMigrationState(ctx)
    if state != CANCELLED && state != FAILED {
        return nil, ErrCannotUnlock
    }
    
    // Restore positions and allow withdrawal
    return k.restoreLegacyPositions(ctx, msg.Signer)
}
```

### Phase 5: Transition Management (Weeks 9-10)

#### 5.1 Migration States and Transitions

```go
func (k Keeper) UpdateMigrationState(ctx context.Context, newState MigrationState) error {
    currentState := k.GetMigrationState(ctx)
    
    // Validate state transition
    validTransitions := map[MigrationState][]MigrationState{
        NOT_STARTED: {ACTIVE},
        ACTIVE:      {CLOSING, CANCELLED},
        CLOSING:     {LOCKED, CANCELLED},
        LOCKED:      {DEPRECATED},
        CANCELLED:   {NOT_STARTED},
    }
    
    allowed := false
    for _, valid := range validTransitions[currentState] {
        if valid == newState {
            allowed = true
            break
        }
    }
    
    if !allowed {
        return ErrInvalidStateTransition
    }
    
    // Execute transition actions
    switch newState {
    case CLOSING:
        k.announceClosingPeriod(ctx)
    case LOCKED:
        k.lockLegacyVault(ctx)
    case DEPRECATED:
        k.deprecateLegacyVault(ctx)
    }
    
    return k.MigrationState.Set(ctx, newState)
}
```

#### 5.2 Unmigrated Position Handling

```go
func (k Keeper) HandleUnmigratedPositions(ctx context.Context) error {
    state := k.GetMigrationState(ctx)
    if state != DEPRECATED {
        return ErrNotReady
    }
    
    // Option 1: Force migrate remaining positions
    unmigrated, _ := k.legacyKeeper.GetAllUnmigratedPositions(ctx)
    
    for _, pos := range unmigrated {
        // Migrate at current NAV
        k.forceMigratePosition(ctx, pos)
    }
    
    // Option 2: Move to withdrawal-only mode
    // Users can only withdraw principal, no yield
    
    return nil
}
```

### Phase 6: User Interface & Communication (Weeks 11-12)

#### 6.1 Migration Status Query

```proto
service Query {
    rpc MigrationStatus(QueryMigrationStatusRequest) returns (QueryMigrationStatusResponse);
    rpc UserMigrationStatus(QueryUserMigrationStatusRequest) returns (QueryUserMigrationStatusResponse);
}

message QueryMigrationStatusResponse {
    MigrationState state = 1;
    string total_migrated = 2;
    string total_remaining = 3;
    int64 users_migrated = 4;
    int64 users_remaining = 5;
}

message QueryUserMigrationStatusResponse {
    bool has_migrated = 1;
    repeated LegacyPosition legacy_positions = 2;
    MigrationRecord migration_record = 3;
    string estimated_shares = 4;
}
```

#### 6.2 Migration Preview

```go
func (k Keeper) GetMigrationPreview(ctx context.Context, user sdk.AccAddress) (*MigrationPreview, error) {
    positions, _ := k.legacyKeeper.GetUserPositions(ctx, user, ALL)
    if len(positions) == 0 {
        return nil, ErrNoPositions
    }
    
    total, principal, rewards := k.calculateMigrationAmounts(ctx, positions, math.ZeroInt())
    shares := k.calculateMigrationShares(ctx, total)
    
    return &MigrationPreview{
        TotalValue:       total,
        Principal:        principal,
        AccruedRewards:   rewards,
        EstimatedShares:  shares,
    }, nil
}
```

## Migration Timeline

### Pre-Launch Phase (2 weeks before)
- Deploy V2 contracts (inactive)
- Announce migration plan
- Release migration UI/tools
- Partner integration testing

### Phase 1: Migration Opens (Week 1)
- **State**: ACTIVE
- **Target**: Begin user migrations
- **Actions**: User education, documentation

### Phase 2: Main Migration (Weeks 2-12)
- **State**: ACTIVE
- **Target**: Majority of TVL migrated
- **Actions**: Regular reminders, support

### Phase 3: Grace Period (Weeks 13-16)
- **State**: CLOSING
- **Target**: Final migrations
- **Actions**: Urgent notifications

### Phase 4: Legacy Lockdown (Week 17+)
- **State**: LOCKED
- **Legacy**: Withdrawals only, no new deposits
- **Actions**: Handle remaining positions

### Phase 5: Deprecation (Month 6+)
- **State**: DEPRECATED
- **Legacy**: Emergency withdrawals only
- **Actions**: Final cleanup

## Risk Mitigation

### Technical Risks

1. **Double-spend Protection**
   - Lock positions before migration
   - Atomic state transitions
   - Comprehensive testing

2. **Migration Failure Handling**
   - Automatic rollback on error
   - Emergency unlock mechanism
   - Position recovery tools

3. **State Consistency**
   - Continuous reconciliation
   - Migration audit trail
   - Checksum verification

### User Experience Risks

1. **Confusion/Errors**
   - Clear UI with previews
   - Step-by-step guides
   - Support documentation

2. **Procrastination**
   - Regular reminders
   - Clear deadlines
   - Communication plan

3. **Loss of Funds Fear**
   - Transparent process
   - Test migrations
   - Security guarantees

## Monitoring & Metrics

### Key Metrics

```go
type MigrationMetrics struct {
    // Progress
    TotalValueMigrated   math.Int
    TotalValueRemaining  math.Int
    UsersMigrated        int64
    UsersRemaining       int64
    
    // Performance
    AverageMigrationTime time.Duration
    FailureRate          sdk.Dec
    GasUsedPerMigration  uint64
    
    // Health
    ErrorCount           int64
    RollbackCount        int64
}
```

### Monitoring Dashboard

- Real-time migration progress
- User migration funnel
- Error rates and types
- Gas usage patterns

## Emergency Procedures

### Migration Pause

```go
func (k msgServer) PauseMigration(ctx context.Context, msg *MsgPauseMigration) error {
    // Authority only
    if msg.Authority != k.authority {
        return sdkerrors.ErrUnauthorized
    }
    
    // Save current state
    k.SaveMigrationCheckpoint(ctx)
    
    // Pause migrations
    k.MigrationPaused.Set(ctx, true)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "migration_paused",
            sdk.NewAttribute("reason", msg.Reason),
            sdk.NewAttribute("block", fmt.Sprintf("%d", ctx.BlockHeight())),
        ),
    )
    
    return nil
}
```

### Recovery Procedures

1. **Position Recovery**: Restore from locked positions
2. **State Rollback**: Revert to checkpoint
3. **Emergency Withdrawal**: Allow principal recovery
4. **Manual Migration**: Admin-assisted migration

## Success Criteria

### Technical Success
- [ ] 95%+ positions migrated successfully
- [ ] Zero loss of user funds
- [ ] Migration completed within timeline
- [ ] All integrations updated

### Business Success
- [ ] Minimal TVL loss (<5%)
- [ ] User satisfaction maintained
- [ ] Partner relationships intact
- [ ] New features accessible

## Post-Migration

### Cleanup Tasks
1. Archive legacy code
2. Optimize state storage
3. Update documentation

### Long-term Support
1. Historical data access
2. Tax reporting support
3. Legacy API wrapper (if needed)
4. User support tools

## Migration Transaction Example

```bash
# User checks their migration status
nobeld query vaults user-migration-status --address noble1abc...

# User previews migration (no state change)
nobeld query vaults migration-preview --address noble1abc...

# User executes migration
nobeld tx vaults migrate-position \
  --vault-type FLEXIBLE \
  --from noble1abc...
```

This user-initiated migration approach ensures explicit consent while maintaining system integrity throughout the transition period.