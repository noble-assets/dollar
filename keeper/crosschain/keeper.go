package crosschain

import (
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	dollarv2 "dollar.noble.xyz/v2/types/v2"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// CrossChainProvider defines the interface for cross-chain providers
type CrossChainProvider interface {
	// GetProviderType returns the provider type (IBC or HYPERLANE)
	GetProviderType() dollarv2.Provider

	// SendMessage sends a cross-chain message
	SendMessage(ctx sdk.Context, route *vaultsv2.CrossChainRoute, msg CrossChainMessage) (*vaultsv2.ProviderTrackingInfo, error)

	// GetMessageStatus checks the status of a cross-chain message
	GetMessageStatus(ctx sdk.Context, tracking *vaultsv2.ProviderTrackingInfo) (MessageStatus, error)

	// GetConfirmations returns the number of confirmations for a message
	GetConfirmations(ctx sdk.Context, tracking *vaultsv2.ProviderTrackingInfo) (uint64, error)

	// EstimateGas estimates gas cost for a cross-chain operation
	EstimateGas(ctx sdk.Context, route *vaultsv2.CrossChainRoute, msg CrossChainMessage) (uint64, math.Int, error)

	// ValidateConfig validates provider-specific configuration
	ValidateConfig(config *vaultsv2.CrossChainProviderConfig) error
}

// CrossChainMessage represents a cross-chain message
type CrossChainMessage struct {
	Type      MessageType
	Sender    sdk.AccAddress
	Recipient string
	Amount    math.Int
	Data      []byte
}

// MessageType defines the type of cross-chain message
type MessageType int32

const (
	MessageTypeDeposit MessageType = iota
	MessageTypeWithdraw
	MessageTypeUpdate
	MessageTypeLiquidate
)

// MessageStatus represents the status of a cross-chain message
type MessageStatus int32

const (
	MessageStatusPending MessageStatus = iota
	MessageStatusSent
	MessageStatusDelivered
	MessageStatusConfirmed
	MessageStatusFailed
	MessageStatusTimeout
)

// CrossChainKeeper manages cross-chain operations for vaults
type CrossChainKeeper struct {
	// Provider implementations
	providers map[dollarv2.Provider]CrossChainProvider

	// Collections for state management
	routes      collections.Map[string, vaultsv2.CrossChainRoute]
	positions   collections.Map[collections.Pair[string, []byte], vaultsv2.RemotePosition]
	inFlight    collections.Map[uint64, vaultsv2.InFlightPosition]
	snapshots   collections.Map[collections.Pair[int32, int64], vaultsv2.CrossChainPositionSnapshot]
	driftAlerts collections.Map[collections.Pair[string, []byte], vaultsv2.DriftAlert]
	config      collections.Item[vaultsv2.CrossChainConfig]

	// Nonce counter for operations
	nonceCounter collections.Sequence
}

// NewCrossChainKeeper creates a new cross-chain keeper
func NewCrossChainKeeper(
	routes collections.Map[string, vaultsv2.CrossChainRoute],
	positions collections.Map[collections.Pair[string, []byte], vaultsv2.RemotePosition],
	inFlight collections.Map[uint64, vaultsv2.InFlightPosition],
	snapshots collections.Map[collections.Pair[int32, int64], vaultsv2.CrossChainPositionSnapshot],
	driftAlerts collections.Map[collections.Pair[string, []byte], vaultsv2.DriftAlert],
	config collections.Item[vaultsv2.CrossChainConfig],
	nonceCounter collections.Sequence,
) *CrossChainKeeper {
	return &CrossChainKeeper{
		providers:    make(map[dollarv2.Provider]CrossChainProvider),
		routes:       routes,
		positions:    positions,
		inFlight:     inFlight,
		snapshots:    snapshots,
		driftAlerts:  driftAlerts,
		config:       config,
		nonceCounter: nonceCounter,
	}
}

// RegisterProvider registers a cross-chain provider
func (k *CrossChainKeeper) RegisterProvider(provider CrossChainProvider) {
	k.providers[provider.GetProviderType()] = provider
}

// GetProvider returns a provider by type
func (k *CrossChainKeeper) GetProvider(providerType dollarv2.Provider) (CrossChainProvider, error) {
	provider, exists := k.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("provider %s not registered", providerType.String())
	}
	return provider, nil
}

// CreateRoute creates a new cross-chain route
func (k *CrossChainKeeper) CreateRoute(ctx sdk.Context, route *vaultsv2.CrossChainRoute) error {
	// Validate route
	if err := k.validateRoute(ctx, route); err != nil {
		return fmt.Errorf("invalid route: %w", err)
	}

	// Check if route already exists
	has, err := k.routes.Has(ctx, route.RouteId)
	if err != nil {
		return fmt.Errorf("failed to check if route exists: %w", err)
	}
	if has {
		return fmt.Errorf("route %s already exists", route.RouteId)
	}

	// Store route
	return k.routes.Set(ctx, route.RouteId, *route)
}

// UpdateRoute updates an existing cross-chain route
func (k *CrossChainKeeper) UpdateRoute(ctx sdk.Context, routeId string, route *vaultsv2.CrossChainRoute) error {
	// Check if route exists
	has, err := k.routes.Has(ctx, routeId)
	if err != nil {
		return fmt.Errorf("failed to check if route exists: %w", err)
	}
	if !has {
		return fmt.Errorf("route %s not found", routeId)
	}

	// Validate updated route
	if err := k.validateRoute(ctx, route); err != nil {
		return fmt.Errorf("invalid route update: %w", err)
	}

	// Update route
	route.RouteId = routeId
	return k.routes.Set(ctx, routeId, *route)
}

// DisableRoute disables a cross-chain route
func (k *CrossChainKeeper) DisableRoute(ctx sdk.Context, routeId string) error {
	route, err := k.routes.Get(ctx, routeId)
	if err != nil {
		return fmt.Errorf("route %s not found: %w", routeId, err)
	}

	route.Active = false
	return k.routes.Set(ctx, routeId, route)
}

// GetRoute retrieves a cross-chain route
func (k *CrossChainKeeper) GetRoute(ctx sdk.Context, routeId string) (*vaultsv2.CrossChainRoute, error) {
	route, err := k.routes.Get(ctx, routeId)
	if err != nil {
		return nil, fmt.Errorf("route %s not found: %w", routeId, err)
	}
	return &route, nil
}

// GetAllRoutes retrieves all cross-chain routes
func (k *CrossChainKeeper) GetAllRoutes(ctx sdk.Context) ([]*vaultsv2.CrossChainRoute, error) {
	var routes []*vaultsv2.CrossChainRoute

	err := k.routes.Walk(ctx, nil, func(key string, value vaultsv2.CrossChainRoute) (bool, error) {
		routes = append(routes, &value)
		return false, nil
	})

	return routes, err
}

// InitiateRemoteDeposit initiates a deposit to a remote chain
func (k *CrossChainKeeper) InitiateRemoteDeposit(
	ctx sdk.Context,
	depositor sdk.AccAddress,
	routeId string,
	amount math.Int,
	remoteAddress string,
	gasLimit uint64,
	gasPrice math.Int,
) (uint64, error) {
	// Get route
	route, err := k.GetRoute(ctx, routeId)
	if err != nil {
		return 0, err
	}

	if !route.Active {
		return 0, fmt.Errorf("route %s is not active", routeId)
	}

	// Validate amount against route limits
	if amount.GT(route.MaxPositionValue) {
		return 0, fmt.Errorf("amount %s exceeds route limit %s", amount.String(), route.MaxPositionValue.String())
	}

	// Get provider
	provider, err := k.GetProvider(route.Provider)
	if err != nil {
		return 0, err
	}

	// Generate nonce
	nonce, err := k.nonceCounter.Next(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Create cross-chain message
	msg := CrossChainMessage{
		Type:      MessageTypeDeposit,
		Sender:    depositor,
		Recipient: remoteAddress,
		Amount:    amount,
	}

	// Send message
	tracking, err := provider.SendMessage(ctx, route, msg)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %w", err)
	}

	// Create in-flight position
	inFlightPos := vaultsv2.InFlightPosition{
		Nonce:                 nonce,
		RouteId:               routeId,
		UserAddress:           depositor.Bytes(),
		OperationType:         vaultsv2.OPERATION_REMOTE_DEPOSIT,
		Amount:                amount,
		InitiatedAt:           ctx.BlockTime(),
		ExpectedCompletion:    ctx.BlockTime().Add(time.Duration(route.RiskParams.OperationTimeout) * time.Second),
		RetryCount:            0,
		Status:                vaultsv2.INFLIGHT_PENDING,
		Provider:              route.Provider,
		ProviderTracking:      tracking,
		Confirmations:         0,
		RequiredConfirmations: k.getRequiredConfirmations(route.Provider),
	}

	// Store in-flight position
	if err := k.inFlight.Set(ctx, nonce, inFlightPos); err != nil {
		return 0, fmt.Errorf("failed to store in-flight position: %w", err)
	}

	return nonce, nil
}

// InitiateRemoteWithdraw initiates a withdrawal from a remote chain
func (k *CrossChainKeeper) InitiateRemoteWithdraw(
	ctx sdk.Context,
	withdrawer sdk.AccAddress,
	routeId string,
	shares math.Int,
	gasLimit uint64,
	gasPrice math.Int,
) (uint64, error) {
	// Get route
	route, err := k.GetRoute(ctx, routeId)
	if err != nil {
		return 0, err
	}

	if !route.Active {
		return 0, fmt.Errorf("route %s is not active", routeId)
	}

	// Check if user has remote position
	positionKey := collections.Join(routeId, withdrawer.Bytes())
	position, err := k.positions.Get(ctx, positionKey)
	if err != nil {
		return 0, fmt.Errorf("no remote position found for user %s on route %s", withdrawer.String(), routeId)
	}

	if shares.GT(position.AllocatedShares) {
		return 0, fmt.Errorf("insufficient shares: requested %s, available %s", shares.String(), position.AllocatedShares.String())
	}

	// Get provider
	provider, err := k.GetProvider(route.Provider)
	if err != nil {
		return 0, err
	}

	// Generate nonce
	nonce, err := k.nonceCounter.Next(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Calculate amount based on shares
	amount := k.calculateAmountFromShares(shares, position.RemoteValue, position.AllocatedShares)

	// Create cross-chain message
	msg := CrossChainMessage{
		Type:      MessageTypeWithdraw,
		Sender:    withdrawer,
		Recipient: position.RemoteAddress,
		Amount:    amount,
	}

	// Send message
	tracking, err := provider.SendMessage(ctx, route, msg)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %w", err)
	}

	// Create in-flight position
	inFlightPos := vaultsv2.InFlightPosition{
		Nonce:                 nonce,
		RouteId:               routeId,
		UserAddress:           withdrawer.Bytes(),
		OperationType:         vaultsv2.OPERATION_REMOTE_WITHDRAW,
		Amount:                amount,
		Shares:                shares,
		InitiatedAt:           ctx.BlockTime(),
		ExpectedCompletion:    ctx.BlockTime().Add(time.Duration(route.RiskParams.OperationTimeout) * time.Second),
		RetryCount:            0,
		Status:                vaultsv2.INFLIGHT_PENDING,
		Provider:              route.Provider,
		ProviderTracking:      tracking,
		Confirmations:         0,
		RequiredConfirmations: k.getRequiredConfirmations(route.Provider),
	}

	// Store in-flight position
	if err := k.inFlight.Set(ctx, nonce, inFlightPos); err != nil {
		return 0, fmt.Errorf("failed to store in-flight position: %w", err)
	}

	return nonce, nil
}

// UpdateRemotePosition updates a remote position
func (k *CrossChainKeeper) UpdateRemotePosition(
	ctx sdk.Context,
	routeId string,
	userAddress []byte,
	remoteValue math.Int,
	confirmations uint64,
	tracking *vaultsv2.ProviderTrackingInfo,
	status vaultsv2.RemotePositionStatus,
) error {
	positionKey := collections.Join(routeId, userAddress)

	// Get existing position or create new one
	position, err := k.positions.Get(ctx, positionKey)
	if err != nil {
		// Create new position
		route, err := k.GetRoute(ctx, routeId)
		if err != nil {
			return err
		}

		position = vaultsv2.RemotePosition{
			RouteId:               routeId,
			LocalAddress:          userAddress,
			RemoteValue:           remoteValue,
			LastUpdate:            ctx.BlockTime(),
			Status:                status,
			Provider:              route.Provider,
			ProviderTracking:      tracking,
			Confirmations:         confirmations,
			RequiredConfirmations: k.getRequiredConfirmations(route.Provider),
		}
	} else {
		// Update existing position
		position.RemoteValue = remoteValue
		position.LastUpdate = ctx.BlockTime()
		position.Status = status
		position.ProviderTracking = tracking
		position.Confirmations = confirmations
	}

	// Calculate conservative value with haircut
	route, err := k.GetRoute(ctx, routeId)
	if err != nil {
		return err
	}

	haircut := math.NewInt(int64(route.RiskParams.PositionHaircut))
	conservative := remoteValue.Mul(math.NewInt(10000).Sub(haircut)).Quo(math.NewInt(10000))
	position.ConservativeValue = conservative

	// Calculate drift
	if !position.AllocatedShares.IsZero() {
		expectedValue := k.calculateExpectedValue(ctx, position.AllocatedShares, routeId)
		drift := k.calculateDrift(remoteValue, expectedValue)
		position.CurrentDrift = drift

		// Check for drift alerts
		if drift > route.RiskParams.MaxDriftThreshold {
			k.createDriftAlert(ctx, routeId, userAddress, drift, route.RiskParams.MaxDriftThreshold)
		}
	}

	return k.positions.Set(ctx, positionKey, position)
}

// ProcessInFlightPosition processes an in-flight operation
func (k *CrossChainKeeper) ProcessInFlightPosition(
	ctx sdk.Context,
	nonce uint64,
	status vaultsv2.InFlightStatus,
	resultAmount math.Int,
	errorMessage string,
	tracking *vaultsv2.ProviderTrackingInfo,
) error {
	inFlightPos, err := k.inFlight.Get(ctx, nonce)
	if err != nil {
		return fmt.Errorf("in-flight position %d not found: %w", nonce, err)
	}

	// Update status and tracking
	inFlightPos.Status = status
	inFlightPos.ProviderTracking = tracking
	if errorMessage != "" {
		inFlightPos.ErrorMessage = errorMessage
	}

	switch status {
	case vaultsv2.INFLIGHT_COMPLETED:
		// Process successful operation
		if err := k.processSuccessfulOperation(ctx, &inFlightPos, resultAmount); err != nil {
			return fmt.Errorf("failed to process successful operation: %w", err)
		}

	case vaultsv2.INFLIGHT_FAILED, vaultsv2.INFLIGHT_TIMEOUT:
		// Process failed operation
		if err := k.processFailedOperation(ctx, &inFlightPos); err != nil {
			return fmt.Errorf("failed to process failed operation: %w", err)
		}
	}

	return k.inFlight.Set(ctx, nonce, inFlightPos)
}

// GetRemotePosition retrieves a remote position
func (k *CrossChainKeeper) GetRemotePosition(ctx sdk.Context, routeId string, userAddress []byte) (*vaultsv2.RemotePosition, error) {
	positionKey := collections.Join(routeId, userAddress)
	position, err := k.positions.Get(ctx, positionKey)
	if err != nil {
		return nil, fmt.Errorf("remote position not found: %w", err)
	}
	return &position, nil
}

// GetInFlightPosition retrieves an in-flight position
func (k *CrossChainKeeper) GetInFlightPosition(ctx sdk.Context, nonce uint64) (*vaultsv2.InFlightPosition, error) {
	position, err := k.inFlight.Get(ctx, nonce)
	if err != nil {
		return nil, fmt.Errorf("in-flight position %d not found: %w", nonce, err)
	}
	return &position, nil
}

// Helper methods

func (k *CrossChainKeeper) validateRoute(ctx sdk.Context, route *vaultsv2.CrossChainRoute) error {
	if route.RouteId == "" {
		return fmt.Errorf("route ID cannot be empty")
	}

	if route.SourceChain == "" || route.DestinationChain == "" {
		return fmt.Errorf("source and destination chains cannot be empty")
	}

	if route.MaxPositionValue.IsZero() || route.MaxPositionValue.IsNegative() {
		return fmt.Errorf("max position value must be positive")
	}

	// Validate provider-specific configuration
	provider, err := k.GetProvider(route.Provider)
	if err != nil {
		return fmt.Errorf("unsupported provider: %w", err)
	}

	return provider.ValidateConfig(route.ProviderConfig)
}

func (k *CrossChainKeeper) getRequiredConfirmations(provider dollarv2.Provider) uint64 {
	switch provider {
	case dollarv2.Provider_IBC:
		return 1 // IBC has built-in finality
	case dollarv2.Provider_HYPERLANE:
		return 12 // Ethereum-like confirmations
	default:
		return 6 // Default
	}
}

func (k *CrossChainKeeper) calculateAmountFromShares(shares, totalValue, totalShares math.Int) math.Int {
	if totalShares.IsZero() {
		return math.ZeroInt()
	}
	return shares.Mul(totalValue).Quo(totalShares)
}

func (k *CrossChainKeeper) calculateExpectedValue(ctx sdk.Context, shares math.Int, routeId string) math.Int {
	// This would calculate expected value based on vault NAV and share price
	// Placeholder implementation
	return shares
}

func (k *CrossChainKeeper) calculateDrift(actual, expected math.Int) int32 {
	if expected.IsZero() {
		return 0
	}
	diff := actual.Sub(expected).Abs()
	drift := diff.Mul(math.NewInt(10000)).Quo(expected)
	return int32(drift.Int64())
}

func (k *CrossChainKeeper) createDriftAlert(ctx sdk.Context, routeId string, userAddress []byte, drift, threshold int32) {
	alertKey := collections.Join(routeId, userAddress)
	alert := vaultsv2.DriftAlert{
		RouteId:           routeId,
		UserAddress:       userAddress,
		CurrentDrift:      drift,
		ThresholdExceeded: threshold,
		Timestamp:         ctx.BlockTime(),
		RecommendedAction: "Rebalance position to reduce drift",
	}
	k.driftAlerts.Set(ctx, alertKey, alert)
}

func (k *CrossChainKeeper) processSuccessfulOperation(ctx sdk.Context, inFlightPos *vaultsv2.InFlightPosition, resultAmount math.Int) error {
	switch inFlightPos.OperationType {
	case vaultsv2.OPERATION_REMOTE_DEPOSIT:
		return k.processSuccessfulDeposit(ctx, inFlightPos, resultAmount)
	case vaultsv2.OPERATION_REMOTE_WITHDRAW:
		return k.processSuccessfulWithdraw(ctx, inFlightPos, resultAmount)
	default:
		return fmt.Errorf("unsupported operation type: %s", inFlightPos.OperationType.String())
	}
}

func (k *CrossChainKeeper) processSuccessfulDeposit(ctx sdk.Context, inFlightPos *vaultsv2.InFlightPosition, resultAmount math.Int) error {
	// Create or update remote position
	positionKey := collections.Join(inFlightPos.RouteId, inFlightPos.UserAddress)

	position, err := k.positions.Get(ctx, positionKey)
	if err != nil {
		// Create new position
		route, err := k.GetRoute(ctx, inFlightPos.RouteId)
		if err != nil {
			return err
		}

		position = vaultsv2.RemotePosition{
			RouteId:               inFlightPos.RouteId,
			LocalAddress:          inFlightPos.UserAddress,
			RemoteValue:           resultAmount,
			AllocatedShares:       inFlightPos.Amount, // Assuming 1:1 for now
			LastUpdate:            ctx.BlockTime(),
			Status:                vaultsv2.REMOTE_POSITION_ACTIVE,
			Provider:              route.Provider,
			ProviderTracking:      inFlightPos.ProviderTracking,
			RequiredConfirmations: k.getRequiredConfirmations(route.Provider),
		}
	} else {
		// Update existing position
		position.RemoteValue = position.RemoteValue.Add(resultAmount)
		position.AllocatedShares = position.AllocatedShares.Add(inFlightPos.Amount)
		position.LastUpdate = ctx.BlockTime()
		position.ProviderTracking = inFlightPos.ProviderTracking
	}

	return k.positions.Set(ctx, positionKey, position)
}

func (k *CrossChainKeeper) processSuccessfulWithdraw(ctx sdk.Context, inFlightPos *vaultsv2.InFlightPosition, resultAmount math.Int) error {
	// Update remote position
	positionKey := collections.Join(inFlightPos.RouteId, inFlightPos.UserAddress)

	position, err := k.positions.Get(ctx, positionKey)
	if err != nil {
		return fmt.Errorf("remote position not found for withdrawal: %w", err)
	}

	// Reduce position
	position.RemoteValue = position.RemoteValue.Sub(resultAmount)
	position.AllocatedShares = position.AllocatedShares.Sub(inFlightPos.Shares)
	position.LastUpdate = ctx.BlockTime()
	position.ProviderTracking = inFlightPos.ProviderTracking

	// Close position if no shares left
	if position.AllocatedShares.IsZero() {
		position.Status = vaultsv2.REMOTE_POSITION_CLOSED
	}

	return k.positions.Set(ctx, positionKey, position)
}

func (k *CrossChainKeeper) processFailedOperation(ctx sdk.Context, inFlightPos *vaultsv2.InFlightPosition) error {
	// For failed operations, we might need to:
	// 1. Return funds to user (for deposits)
	// 2. Restore shares (for withdrawals)
	// 3. Log the failure for analysis

	// This is a placeholder - actual implementation would depend on specific requirements
	return nil
}
