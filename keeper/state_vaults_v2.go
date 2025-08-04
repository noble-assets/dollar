package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	"dollar.noble.xyz/v2/types/vaults"
)

// V2VaultState represents the core state of a V2 vault with share-based accounting
type V2VaultState struct {
	// Total shares outstanding
	TotalShares math.Int
	// Current NAV point (total value of all assets)
	NavPoint math.Int
	// Last NAV update timestamp
	LastNavUpdate *timestamppb.Timestamp
	// Total principal deposited (for tracking)
	TotalPrincipal math.Int
	// Accumulated yield
	AccumulatedYield math.Int
	// Emergency flags
	CircuitBreakerActive bool
	DepositsPaused       bool
	WithdrawalsPaused    bool
}

// V2UserPosition represents a user's position in the V2 vault system
type V2UserPosition struct {
	// User's share balance
	Shares math.Int
	// Principal deposited (for tracking)
	PrincipalDeposited math.Int
	// Average entry price
	AvgEntryPrice math.LegacyDec
	// First deposit timestamp
	FirstDeposit *timestamppb.Timestamp
	// Last activity timestamp
	LastActivity *timestamppb.Timestamp
	// Yield preference setting
	ForgoYield bool
}

// V2ExitRequest represents a request to exit the vault via the queue
type V2ExitRequest struct {
	// Request ID
	RequestID string
	// User address
	UserAddress sdk.AccAddress
	// Vault type
	VaultType vaults.VaultType
	// Shares to exit with
	Shares math.Int
	// Request timestamp
	RequestedAt *timestamppb.Timestamp
	// Queue position
	QueuePosition uint64
	// Expected tokens at request time
	ExpectedTokens math.Int
	// Status
	Status ExitRequestStatus
}

// ExitRequestStatus represents the status of an exit request
type ExitRequestStatus int32

const (
	ExitStatusPending ExitRequestStatus = iota
	ExitStatusProcessing
	ExitStatusCompleted
	ExitStatusCancelled
	ExitStatusExpired
)

// V2RemotePosition represents a position on another chain
type V2RemotePosition struct {
	// Route used for this position
	RouteID string
	// Remote address
	RemoteAddress string
	// Local user address
	LocalAddress sdk.AccAddress
	// Position value on remote chain
	RemoteValue math.Int
	// Conservative valuation (with haircut)
	ConservativeValue math.Int
	// Last update timestamp
	LastUpdate *timestamppb.Timestamp
	// Current drift (basis points)
	CurrentDrift int32
	// Status
	Status RemotePositionStatus
	// Allocated shares
	AllocatedShares math.Int
}

// RemotePositionStatus represents the status of a remote position
type RemotePositionStatus int32

const (
	RemotePositionActive RemotePositionStatus = iota
	RemotePositionPendingUpdate
	RemotePositionDriftExceeded
	RemotePositionLiquidating
	RemotePositionClosed
	RemotePositionError
)

// V2InFlightPosition represents a position operation in progress
type V2InFlightPosition struct {
	// Unique nonce
	Nonce uint64
	// Route ID
	RouteID string
	// User address
	UserAddress sdk.AccAddress
	// Operation type
	OperationType InFlightOperationType
	// Amount involved
	Amount math.Int
	// Shares involved
	Shares math.Int
	// Initiated timestamp
	InitiatedAt *timestamppb.Timestamp
	// Expected completion
	ExpectedCompletion *timestamppb.Timestamp
	// Retry count
	RetryCount int32
	// Status
	Status InFlightStatus
	// Error message
	ErrorMessage string
}

// InFlightOperationType defines the type of cross-chain operation
type InFlightOperationType int32

const (
	OperationRemoteDeposit InFlightOperationType = iota
	OperationRemoteWithdraw
	OperationRebalance
	OperationLiquidate
)

// InFlightStatus represents the status of an in-flight operation
type InFlightStatus int32

const (
	InFlightPending InFlightStatus = iota
	InFlightProcessing
	InFlightCompleted
	InFlightFailed
	InFlightTimeout
	InFlightCancelled
)

// V2MigrationState represents the current migration state
type V2MigrationState int32

const (
	MigrationStateNotStarted V2MigrationState = iota
	MigrationStateActive
	MigrationStateClosing
	MigrationStateLocked
	MigrationStateDeprecated
	MigrationStateCancelled
)

// V2MigrationStats tracks migration progress
type V2MigrationStats struct {
	// Total users in legacy system
	TotalUsers uint64
	// Users migrated
	UsersMigrated uint64
	// Total value locked in legacy
	TotalValueLocked math.Int
	// Value migrated
	ValueMigrated math.Int
	// Total shares issued through migration
	TotalSharesIssued math.Int
	// Last migration timestamp
	LastMigrationTime *timestamppb.Timestamp
	// Average gas per migration
	AverageGasPerMigration uint64
	// Completion percentage (basis points)
	CompletionPercentage int32
}

// V2UserMigrationRecord tracks migration details for a user
type V2UserMigrationRecord struct {
	// Migration timestamp
	MigratedAt *timestamppb.Timestamp
	// Source vault type
	FromVaultType vaults.VaultType
	// Legacy position count
	LegacyPositionCount int64
	// Principal migrated
	PrincipalMigrated math.Int
	// Rewards migrated
	RewardsMigrated math.Int
	// Shares received
	SharesReceived math.Int
	// Migration transaction hash
	MigrationTxHash string
	// Gas used
	GasUsed uint64
	// Whether yield was forgone
	YieldForgone bool
}

// V2LockedLegacyPosition represents a locked legacy position
type V2LockedLegacyPosition struct {
	// Original position
	Position vaults.Position
	// Lock timestamp
	LockedAt *timestamppb.Timestamp
	// Migrated to address
	MigratedTo sdk.AccAddress
	// Migration ID
	MigrationID string
	// Whether unlock is enabled
	UnlockEnabled bool
	// Lock reason
	LockReason string
}

// V2NAVConfig defines NAV configuration
type V2NAVConfig struct {
	// Minimum update interval (seconds)
	MinNavUpdateInterval int64
	// Maximum deviation threshold (basis points)
	MaxNavDeviation int32
	// Circuit breaker threshold (basis points)
	CircuitBreakerThreshold int32
}

// V2FeeConfig defines fee configuration
type V2FeeConfig struct {
	// Management fee rate (annual basis points)
	ManagementFeeRate int32
	// Performance fee rate (basis points)
	PerformanceFeeRate int32
	// Deposit fee rate (basis points)
	DepositFeeRate int32
	// Withdrawal fee rate (basis points)
	WithdrawalFeeRate int32
	// Fee recipient
	FeeRecipient sdk.AccAddress
	// Whether fees are enabled
	FeesEnabled bool
	// High water mark for performance fees
	HighWaterMark math.LegacyDec
}

// V2FeeAccrual represents accrued fees
type V2FeeAccrual struct {
	// Fee type
	FeeType FeeType
	// Accrued amount
	AccruedAmount math.Int
	// Shares to dilute
	SharesToDilute math.Int
	// Period start
	PeriodStart *timestamppb.Timestamp
	// Period end
	PeriodEnd *timestamppb.Timestamp
	// Whether collected
	Collected bool
}

// FeeType defines fee types
type FeeType int32

const (
	FeeTypeManagement FeeType = iota
	FeeTypePerformance
	FeeTypeDeposit
	FeeTypeWithdrawal
	FeeTypeCrossChain
	FeeTypeEmergency
)

// V2CrossChainRoute defines a cross-chain route
type V2CrossChainRoute struct {
	// Route ID
	RouteID string
	// Source chain
	SourceChain string
	// Destination chain
	DestinationChain string
	// IBC channel
	IBCChannel string
	// Whether active
	Active bool
	// Maximum position value
	MaxPositionValue math.Int
	// Position haircut (basis points)
	PositionHaircut int32
	// Max drift threshold (basis points)
	MaxDriftThreshold int32
	// Operation timeout (seconds)
	OperationTimeout int64
	// Max retries
	MaxRetries int32
	// Conservative discount (basis points)
	ConservativeDiscount int32
}

// V2 State Keys for collections
var (
	// V2 Vault State Keys
	V2VaultStatePrefix       = collections.NewPrefix(200) // V2 vault states by type
	V2UserPositionPrefix     = collections.NewPrefix(201) // User positions: (vault_type, address) -> Position
	V2ExitRequestPrefix      = collections.NewPrefix(202) // Exit requests: request_id -> ExitRequest
	V2ExitQueuePrefix        = collections.NewPrefix(203) // Exit queue: (vault_type, queue_pos) -> request_id
	V2RemotePositionPrefix   = collections.NewPrefix(204) // Remote positions: (route_id, address) -> RemotePosition
	V2InFlightPositionPrefix = collections.NewPrefix(205) // In-flight operations: nonce -> InFlightPosition

	// Migration State Keys
	V2MigrationStateKey          = collections.NewPrefix(210) // Current migration state
	V2MigrationStatsKey          = collections.NewPrefix(211) // Migration statistics
	V2UserMigrationRecordPrefix  = collections.NewPrefix(212) // User migration records: address -> MigrationRecord
	V2LockedLegacyPositionPrefix = collections.NewPrefix(213) // Locked legacy positions: (address, vault_type, index) -> LockedPosition

	// Configuration Keys
	V2NAVConfigPrefix       = collections.NewPrefix(220) // NAV configs by vault type
	V2FeeConfigPrefix       = collections.NewPrefix(221) // Fee configs by vault type
	V2CrossChainRoutePrefix = collections.NewPrefix(222) // Cross-chain routes: route_id -> Route

	// Fee Management Keys
	V2FeeAccrualPrefix  = collections.NewPrefix(230) // Fee accruals: (vault_type, fee_type, period) -> Accrual
	V2LastFeeAccrualKey = collections.NewPrefix(231) // Last fee accrual timestamp by vault type

	// NAV and Pricing Keys
	V2LastNAVUpdatePrefix = collections.NewPrefix(240) // Last NAV update by vault type
	V2NAVHistoryPrefix    = collections.NewPrefix(241) // NAV history: (vault_type, timestamp) -> NAV
)
