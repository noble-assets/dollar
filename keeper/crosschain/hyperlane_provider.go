package crosschain

import (
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	dollarv2 "dollar.noble.xyz/v2/types/v2"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// HyperlaneProvider implements CrossChainProvider for Hyperlane
type HyperlaneProvider struct {
	mailboxKeeper    HyperlaneMailboxKeeper
	gasPriceFeed     HyperlaneGasPriceFeed
	localDomain      uint32
	defaultGasLimit  uint64
	defaultTimeout   time.Duration
	confirmationsMap map[uint32]uint64 // domain -> required confirmations
}

// HyperlaneMailboxKeeper defines the expected Hyperlane mailbox keeper interface
type HyperlaneMailboxKeeper interface {
	// SendMessage sends a message through Hyperlane
	SendMessage(ctx sdk.Context, destinationDomain uint32, recipient common.Address, messageBody []byte, gasLimit uint64, gasPrice *big.Int) ([]byte, error)

	// GetMessageStatus returns the delivery status of a message
	GetMessageStatus(ctx sdk.Context, messageId []byte) (HyperlaneMessageStatus, error)

	// GetMessage retrieves a message by ID
	GetMessage(ctx sdk.Context, messageId []byte) (*HyperlaneMessage, error)

	// GetLatestCheckpoint returns the latest checkpoint for a domain
	GetLatestCheckpoint(ctx sdk.Context, domain uint32) (*HyperlaneCheckpoint, error)

	// VerifyMessageDelivery verifies that a message was delivered
	VerifyMessageDelivery(ctx sdk.Context, messageId []byte, proof []byte) (bool, error)
}

// HyperlaneGasPriceFeed provides gas price information for different domains
type HyperlaneGasPriceFeed interface {
	// GetGasPrice returns the current gas price for a domain
	GetGasPrice(ctx sdk.Context, domain uint32) (*big.Int, error)

	// EstimateGasCost estimates the gas cost for a cross-chain operation
	EstimateGasCost(ctx sdk.Context, destinationDomain uint32, gasLimit uint64) (*big.Int, error)
}

// HyperlaneMessageStatus represents the status of a Hyperlane message
type HyperlaneMessageStatus int32

const (
	HyperlaneMessageStatusPending HyperlaneMessageStatus = iota
	HyperlaneMessageStatusDispatched
	HyperlaneMessageStatusDelivered
	HyperlaneMessageStatusProcessed
	HyperlaneMessageStatusFailed
)

// HyperlaneMessage represents a Hyperlane message
type HyperlaneMessage struct {
	ID                []byte
	Nonce             uint64
	OriginDomain      uint32
	DestinationDomain uint32
	Sender            common.Address
	Recipient         common.Address
	MessageBody       []byte
	DispatchedBlock   uint64
	DeliveredBlock    uint64
	GasLimit          uint64
	GasPrice          *big.Int
	Status            HyperlaneMessageStatus
}

// HyperlaneCheckpoint represents a checkpoint for message verification
type HyperlaneCheckpoint struct {
	Domain      uint32
	Root        []byte
	Index       uint64
	BlockNumber uint64
	BlockHash   []byte
}

// NewHyperlaneProvider creates a new Hyperlane provider
func NewHyperlaneProvider(
	mailboxKeeper HyperlaneMailboxKeeper,
	gasPriceFeed HyperlaneGasPriceFeed,
	localDomain uint32,
	defaultGasLimit uint64,
	defaultTimeout time.Duration,
) *HyperlaneProvider {
	// Default confirmation requirements for common chains
	confirmationsMap := map[uint32]uint64{
		1:     12, // Ethereum mainnet
		137:   20, // Polygon
		42161: 1,  // Arbitrum
		10:    1,  // Optimism
		56:    3,  // BSC
		43114: 1,  // Avalanche
	}

	return &HyperlaneProvider{
		mailboxKeeper:    mailboxKeeper,
		gasPriceFeed:     gasPriceFeed,
		localDomain:      localDomain,
		defaultGasLimit:  defaultGasLimit,
		defaultTimeout:   defaultTimeout,
		confirmationsMap: confirmationsMap,
	}
}

// GetProviderType returns the Hyperlane provider type
func (p *HyperlaneProvider) GetProviderType() dollarv2.Provider {
	return dollarv2.Provider_HYPERLANE
}

// SendMessage sends a cross-chain message via Hyperlane
func (p *HyperlaneProvider) SendMessage(ctx sdk.Context, route *vaultsv2.CrossChainRoute, msg CrossChainMessage) (*vaultsv2.ProviderTrackingInfo, error) {
	// Validate route has Hyperlane config
	if route.ProviderConfig == nil || route.ProviderConfig.GetHyperlaneConfig() == nil {
		return nil, fmt.Errorf("route %s missing Hyperlane configuration", route.RouteId)
	}

	hyperlaneConfig := route.ProviderConfig.GetHyperlaneConfig()

	// Validate mailbox address
	if hyperlaneConfig.MailboxAddress == "" {
		return nil, fmt.Errorf("mailbox address is required for Hyperlane route %s", route.RouteId)
	}

	// Convert recipient address
	recipientAddr := common.HexToAddress(msg.Recipient)
	if recipientAddr == (common.Address{}) {
		return nil, fmt.Errorf("invalid recipient address: %s", msg.Recipient)
	}

	// Create message body
	messageBody, err := p.createMessageBody(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to create message body: %w", err)
	}

	// Determine gas parameters
	gasLimit := hyperlaneConfig.GasLimit
	if gasLimit == 0 {
		gasLimit = p.defaultGasLimit
	}

	var gasPrice *big.Int
	if !hyperlaneConfig.GasPrice.IsZero() {
		gasPrice = hyperlaneConfig.GasPrice.BigInt()
	} else {
		// Get current gas price from feed
		gasPriceFeed, err := p.gasPriceFeed.GetGasPrice(ctx, hyperlaneConfig.DomainId)
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price: %w", err)
		}
		gasPrice = gasPriceFeed
	}

	// Send message through mailbox
	messageId, err := p.mailboxKeeper.SendMessage(
		ctx,
		hyperlaneConfig.DomainId,
		recipientAddr,
		messageBody,
		gasLimit,
		gasPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send Hyperlane message: %w", err)
	}

	// Create tracking info
	tracking := &vaultsv2.ProviderTrackingInfo{
		TrackingInfo: &vaultsv2.ProviderTrackingInfo_HyperlaneTracking{
			HyperlaneTracking: &vaultsv2.HyperlaneTrackingInfo{
				MessageId:         messageId,
				OriginDomain:      p.localDomain,
				DestinationDomain: hyperlaneConfig.DomainId,
				Nonce:             0, // Will be set by mailbox
				OriginBlockNumber: uint64(ctx.BlockHeight()),
				Processed:         false,
				GasUsed:           gasLimit, // Initial estimate
			},
		},
	}

	return tracking, nil
}

// GetMessageStatus checks the status of a Hyperlane message
func (p *HyperlaneProvider) GetMessageStatus(ctx sdk.Context, tracking *vaultsv2.ProviderTrackingInfo) (MessageStatus, error) {
	hyperlaneTracking := tracking.GetHyperlaneTracking()
	if hyperlaneTracking == nil {
		return MessageStatusFailed, fmt.Errorf("invalid Hyperlane tracking info")
	}

	// Get message status from mailbox
	status, err := p.mailboxKeeper.GetMessageStatus(ctx, hyperlaneTracking.MessageId)
	if err != nil {
		return MessageStatusFailed, fmt.Errorf("failed to get message status: %w", err)
	}

	// Convert Hyperlane status to MessageStatus
	switch status {
	case HyperlaneMessageStatusPending:
		return MessageStatusPending, nil
	case HyperlaneMessageStatusDispatched:
		return MessageStatusSent, nil
	case HyperlaneMessageStatusDelivered:
		return MessageStatusDelivered, nil
	case HyperlaneMessageStatusProcessed:
		return MessageStatusConfirmed, nil
	case HyperlaneMessageStatusFailed:
		return MessageStatusFailed, nil
	default:
		return MessageStatusFailed, fmt.Errorf("unknown Hyperlane message status: %d", status)
	}
}

// GetConfirmations returns the number of confirmations for a Hyperlane message
func (p *HyperlaneProvider) GetConfirmations(ctx sdk.Context, tracking *vaultsv2.ProviderTrackingInfo) (uint64, error) {
	hyperlaneTracking := tracking.GetHyperlaneTracking()
	if hyperlaneTracking == nil {
		return 0, fmt.Errorf("invalid Hyperlane tracking info")
	}

	// Get message details
	message, err := p.mailboxKeeper.GetMessage(ctx, hyperlaneTracking.MessageId)
	if err != nil {
		return 0, fmt.Errorf("failed to get message: %w", err)
	}

	// If message hasn't been delivered yet, return 0
	if message.Status < HyperlaneMessageStatusDelivered {
		return 0, nil
	}

	// Calculate confirmations based on destination domain
	requiredConfirmations, exists := p.confirmationsMap[hyperlaneTracking.DestinationDomain]
	if !exists {
		requiredConfirmations = 12 // Default for unknown chains
	}

	// Get latest checkpoint for destination domain
	checkpoint, err := p.mailboxKeeper.GetLatestCheckpoint(ctx, hyperlaneTracking.DestinationDomain)
	if err != nil {
		return 0, fmt.Errorf("failed to get checkpoint: %w", err)
	}

	// Calculate confirmations
	if checkpoint.BlockNumber > message.DeliveredBlock {
		confirmations := checkpoint.BlockNumber - message.DeliveredBlock
		if confirmations >= requiredConfirmations {
			return requiredConfirmations, nil
		}
		return confirmations, nil
	}

	return 0, nil
}

// EstimateGas estimates gas cost for a Hyperlane operation
func (p *HyperlaneProvider) EstimateGas(ctx sdk.Context, route *vaultsv2.CrossChainRoute, msg CrossChainMessage) (uint64, math.Int, error) {
	hyperlaneConfig := route.ProviderConfig.GetHyperlaneConfig()
	if hyperlaneConfig == nil {
		return 0, math.ZeroInt(), fmt.Errorf("missing Hyperlane configuration")
	}

	// Base gas calculation
	baseGas := p.defaultGasLimit

	// Adjust based on message type
	switch msg.Type {
	case MessageTypeDeposit:
		baseGas += 50000 // Additional gas for deposit processing
	case MessageTypeWithdraw:
		baseGas += 75000 // Additional gas for withdrawal processing
	case MessageTypeUpdate:
		baseGas += 25000 // Additional gas for update processing
	case MessageTypeLiquidate:
		baseGas += 100000 // Additional gas for liquidation processing
	}

	// Use configured gas limit if provided
	if hyperlaneConfig.GasLimit > 0 {
		baseGas = hyperlaneConfig.GasLimit
	}

	// Estimate total cost
	totalCost, err := p.gasPriceFeed.EstimateGasCost(ctx, hyperlaneConfig.DomainId, baseGas)
	if err != nil {
		// Fallback to configured gas price
		if !hyperlaneConfig.GasPrice.IsZero() {
			totalCost = hyperlaneConfig.GasPrice.BigInt()
			totalCost = totalCost.Mul(totalCost, big.NewInt(int64(baseGas)))
		} else {
			// Default fallback
			totalCost = big.NewInt(int64(baseGas) * 1000000) // 0.001 ETH per gas unit
		}
	}

	return baseGas, math.NewIntFromBigInt(totalCost), nil
}

// ValidateConfig validates Hyperlane-specific configuration
func (p *HyperlaneProvider) ValidateConfig(config *vaultsv2.CrossChainProviderConfig) error {
	hyperlaneConfig := config.GetHyperlaneConfig()
	if hyperlaneConfig == nil {
		return fmt.Errorf("Hyperlane configuration is required")
	}

	if hyperlaneConfig.DomainId == 0 {
		return fmt.Errorf("Hyperlane domain ID is required")
	}

	if hyperlaneConfig.MailboxAddress == "" {
		return fmt.Errorf("Hyperlane mailbox address is required")
	}

	// Validate mailbox address format
	mailboxAddr := common.HexToAddress(hyperlaneConfig.MailboxAddress)
	if mailboxAddr == (common.Address{}) {
		return fmt.Errorf("invalid mailbox address format: %s", hyperlaneConfig.MailboxAddress)
	}

	// Validate gas paymaster address if provided
	if hyperlaneConfig.GasPaymasterAddress != "" {
		paymasterAddr := common.HexToAddress(hyperlaneConfig.GasPaymasterAddress)
		if paymasterAddr == (common.Address{}) {
			return fmt.Errorf("invalid gas paymaster address format: %s", hyperlaneConfig.GasPaymasterAddress)
		}
	}

	// Validate hook address if provided
	if hyperlaneConfig.HookAddress != "" {
		hookAddr := common.HexToAddress(hyperlaneConfig.HookAddress)
		if hookAddr == (common.Address{}) {
			return fmt.Errorf("invalid hook address format: %s", hyperlaneConfig.HookAddress)
		}
	}

	// Validate gas parameters
	if hyperlaneConfig.GasLimit > 0 && hyperlaneConfig.GasLimit < 21000 {
		return fmt.Errorf("gas limit too low: %d (minimum 21000)", hyperlaneConfig.GasLimit)
	}

	if !hyperlaneConfig.GasPrice.IsZero() && hyperlaneConfig.GasPrice.IsNegative() {
		return fmt.Errorf("gas price cannot be negative")
	}

	return nil
}

// UpdateMessageStatus updates the status of a Hyperlane message (called by relayers)
func (p *HyperlaneProvider) UpdateMessageStatus(
	ctx sdk.Context,
	tracking *vaultsv2.ProviderTrackingInfo,
	destinationTxHash string,
	destinationBlockNumber uint64,
	gasUsed uint64,
	processed bool,
) error {
	hyperlaneTracking := tracking.GetHyperlaneTracking()
	if hyperlaneTracking == nil {
		return fmt.Errorf("invalid Hyperlane tracking info")
	}

	hyperlaneTracking.DestinationTxHash = destinationTxHash
	hyperlaneTracking.DestinationBlockNumber = destinationBlockNumber
	hyperlaneTracking.GasUsed = gasUsed
	hyperlaneTracking.Processed = processed

	return nil
}

// GetDomainInfo returns information about a Hyperlane domain
func (p *HyperlaneProvider) GetDomainInfo(ctx sdk.Context, domain uint32) (*HyperlaneDomainInfo, error) {
	checkpoint, err := p.mailboxKeeper.GetLatestCheckpoint(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpoint for domain %d: %w", domain, err)
	}

	gasPrice, err := p.gasPriceFeed.GetGasPrice(ctx, domain)
	if err != nil {
		gasPrice = big.NewInt(0) // Set to zero if not available
	}

	requiredConfirmations, exists := p.confirmationsMap[domain]
	if !exists {
		requiredConfirmations = 12 // Default
	}

	return &HyperlaneDomainInfo{
		DomainID:              domain,
		LatestBlockNumber:     checkpoint.BlockNumber,
		LatestBlockHash:       hexutil.Encode(checkpoint.BlockHash),
		CurrentGasPrice:       gasPrice,
		RequiredConfirmations: requiredConfirmations,
	}, nil
}

// HyperlaneDomainInfo contains information about a Hyperlane domain
type HyperlaneDomainInfo struct {
	DomainID              uint32
	LatestBlockNumber     uint64
	LatestBlockHash       string
	CurrentGasPrice       *big.Int
	RequiredConfirmations uint64
}

// Helper methods

func (p *HyperlaneProvider) createMessageBody(msg CrossChainMessage) ([]byte, error) {
	// Create a structured message body for vault operations
	messageBody := map[string]interface{}{
		"operation": p.messageTypeToString(msg.Type),
		"sender":    msg.Sender.String(),
		"recipient": msg.Recipient,
		"amount":    msg.Amount.String(),
		"timestamp": time.Now().Unix(),
	}

	if len(msg.Data) > 0 {
		messageBody["data"] = hexutil.Encode(msg.Data)
	}

	// In a real implementation, this would be properly serialized (JSON, protobuf, etc.)
	// For now, we'll create a simple byte representation
	return []byte(fmt.Sprintf("%+v", messageBody)), nil
}

func (p *HyperlaneProvider) messageTypeToString(msgType MessageType) string {
	switch msgType {
	case MessageTypeDeposit:
		return "deposit"
	case MessageTypeWithdraw:
		return "withdraw"
	case MessageTypeUpdate:
		return "update"
	case MessageTypeLiquidate:
		return "liquidate"
	default:
		return "unknown"
	}
}

// SetRequiredConfirmations updates the required confirmations for a domain
func (p *HyperlaneProvider) SetRequiredConfirmations(domain uint32, confirmations uint64) {
	p.confirmationsMap[domain] = confirmations
}

// GetRequiredConfirmations returns the required confirmations for a domain
func (p *HyperlaneProvider) GetRequiredConfirmations(domain uint32) uint64 {
	if confirmations, exists := p.confirmationsMap[domain]; exists {
		return confirmations
	}
	return 12 // Default
}
