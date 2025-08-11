package crosschain

import (
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Import the actual bcp-innovations hyperlane-cosmos modules
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	corekeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	warpkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/warp/keeper"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"

	dollarv2 "dollar.noble.xyz/v2/types/v2"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// HyperlaneProvider implements CrossChainProvider for Hyperlane using bcp-innovations/hyperlane-cosmos
type HyperlaneProvider struct {
	coreKeeper       *corekeeper.Keeper
	warpKeeper       *warpkeeper.Keeper
	localDomain      uint32
	defaultGasLimit  uint64
	defaultTimeout   time.Duration
	confirmationsMap map[uint32]uint64 // domain -> required confirmations
	mailboxId        util.HexAddress   // The mailbox ID to use for this provider
}

// NewHyperlaneProvider creates a new Hyperlane provider using the actual cosmos SDK implementation
func NewHyperlaneProvider(
	coreKeeper *corekeeper.Keeper,
	warpKeeper *warpkeeper.Keeper,
	localDomain uint32,
	defaultGasLimit uint64,
	defaultTimeout time.Duration,
	mailboxId util.HexAddress,
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
		coreKeeper:       coreKeeper,
		warpKeeper:       warpKeeper,
		localDomain:      localDomain,
		defaultGasLimit:  defaultGasLimit,
		defaultTimeout:   defaultTimeout,
		confirmationsMap: confirmationsMap,
		mailboxId:        mailboxId,
	}
}

// GetProviderType returns the Hyperlane provider type
func (p *HyperlaneProvider) GetProviderType() dollarv2.Provider {
	return dollarv2.Provider_HYPERLANE
}

// SendMessage sends a cross-chain message via Hyperlane using the actual SDK
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

	// Convert addresses to HexAddress format
	senderHexAddr, err := util.DecodeHexAddress(msg.Sender.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode sender address: %w", err)
	}

	recipientHexAddr, err := util.DecodeHexAddress(msg.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to decode recipient address: %w", err)
	}

	// Create message body based on message type
	messageBody, err := p.createMessageBody(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to create message body: %w", err)
	}

	// Prepare max fee (for now, use a default or calculated value)
	maxFee := sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(100000))) // Example fee

	// Create empty metadata for now
	metadata := util.StandardHookMetadata{}

	// Dispatch message through the core keeper
	messageID, err := p.coreKeeper.DispatchMessage(
		ctx,
		p.mailboxId,              // originMailboxId
		senderHexAddr,            // sender
		maxFee,                   // maxFee
		hyperlaneConfig.DomainId, // destinationDomain
		recipientHexAddr,         // recipient
		messageBody,              // body
		metadata,                 // metadata
		nil,                      // postDispatchHookId (optional)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dispatch Hyperlane message: %w", err)
	}

	// Create tracking info
	tracking := &vaultsv2.ProviderTrackingInfo{
		TrackingInfo: &vaultsv2.ProviderTrackingInfo_HyperlaneTracking{
			HyperlaneTracking: &vaultsv2.HyperlaneTrackingInfo{
				MessageId:         messageID[:],
				OriginDomain:      p.localDomain,
				DestinationDomain: hyperlaneConfig.DomainId,
				Nonce:             0, // Will be set by the core keeper
				OriginBlockNumber: uint64(ctx.BlockHeight()),
				Processed:         false,
				GasUsed:           0, // Will be updated when processed
			},
		},
	}

	return tracking, nil
}

// GetMessageStatus checks the status of a Hyperlane message using the actual SDK
func (p *HyperlaneProvider) GetMessageStatus(ctx sdk.Context, tracking *vaultsv2.ProviderTrackingInfo) (MessageStatus, error) {
	hyperlaneTracking := tracking.GetHyperlaneTracking()
	if hyperlaneTracking == nil {
		return MessageStatusFailed, fmt.Errorf("invalid Hyperlane tracking info")
	}

	// Get the message ID (already in correct []byte format)
	messageID := hyperlaneTracking.MessageId

	// Check if the message has been delivered using the Messages collection
	delivered, err := p.coreKeeper.Messages.Has(ctx, collections.Join(p.mailboxId.GetInternalId(), messageID))
	if err != nil {
		return MessageStatusFailed, fmt.Errorf("failed to check message delivery status: %w", err)
	}

	if delivered {
		return MessageStatusConfirmed, nil
	}

	// If not delivered, assume it's still pending or sent
	// In a real implementation, you might want to check for additional states
	return MessageStatusSent, nil
}

// GetConfirmations returns the number of confirmations for a Hyperlane message
func (p *HyperlaneProvider) GetConfirmations(ctx sdk.Context, tracking *vaultsv2.ProviderTrackingInfo) (uint64, error) {
	hyperlaneTracking := tracking.GetHyperlaneTracking()
	if hyperlaneTracking == nil {
		return 0, fmt.Errorf("invalid Hyperlane tracking info")
	}

	// Check if message is delivered
	status, err := p.GetMessageStatus(ctx, tracking)
	if err != nil {
		return 0, err
	}

	// If message is confirmed, return required confirmations
	if status == MessageStatusConfirmed {
		requiredConfirmations := p.GetRequiredConfirmations(hyperlaneTracking.DestinationDomain)
		return requiredConfirmations, nil
	}

	// Otherwise, return 0 confirmations
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

	// Estimate total cost using configured gas price or default
	var totalCost *big.Int
	if !hyperlaneConfig.GasPrice.IsZero() {
		totalCost = hyperlaneConfig.GasPrice.BigInt()
		totalCost = totalCost.Mul(totalCost, big.NewInt(int64(baseGas)))
	} else {
		// Default fallback
		totalCost = big.NewInt(int64(baseGas) * 1000000) // 0.001 ETH per gas unit
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

	// Validate mailbox address format using HexAddress decoder
	_, err := util.DecodeHexAddress(hyperlaneConfig.MailboxAddress)
	if err != nil {
		return fmt.Errorf("invalid mailbox address format: %s, error: %w", hyperlaneConfig.MailboxAddress, err)
	}

	// Validate gas paymaster address if provided
	if hyperlaneConfig.GasPaymasterAddress != "" {
		_, err := util.DecodeHexAddress(hyperlaneConfig.GasPaymasterAddress)
		if err != nil {
			return fmt.Errorf("invalid gas paymaster address format: %s, error: %w", hyperlaneConfig.GasPaymasterAddress, err)
		}
	}

	// Validate hook address if provided
	if hyperlaneConfig.HookAddress != "" {
		_, err := util.DecodeHexAddress(hyperlaneConfig.HookAddress)
		if err != nil {
			return fmt.Errorf("invalid hook address format: %s, error: %w", hyperlaneConfig.HookAddress, err)
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

	// Update the tracking info
	hyperlaneTracking.DestinationTxHash = destinationTxHash
	hyperlaneTracking.DestinationBlockNumber = destinationBlockNumber
	hyperlaneTracking.GasUsed = gasUsed
	hyperlaneTracking.Processed = processed

	return nil
}

// GetDomainInfo returns information about a Hyperlane domain
func (p *HyperlaneProvider) GetDomainInfo(ctx sdk.Context, domain uint32) (*HyperlaneDomainInfo, error) {
	// Get mailbox info
	mailbox, err := p.coreKeeper.Mailboxes.Get(ctx, p.mailboxId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to get mailbox: %w", err)
	}

	requiredConfirmations := p.GetRequiredConfirmations(domain)

	return &HyperlaneDomainInfo{
		DomainID:              domain,
		LatestBlockNumber:     uint64(ctx.BlockHeight()), // Local block height
		LatestBlockHash:       fmt.Sprintf("%x", ctx.BlockHeader().LastBlockId.Hash),
		CurrentGasPrice:       big.NewInt(0), // Default gas price
		RequiredConfirmations: requiredConfirmations,
		MailboxId:             mailbox.Id.String(),
		MailboxDomain:         p.localDomain,
	}, nil
}

// HyperlaneDomainInfo contains information about a Hyperlane domain
type HyperlaneDomainInfo struct {
	DomainID              uint32
	LatestBlockNumber     uint64
	LatestBlockHash       string
	CurrentGasPrice       *big.Int
	RequiredConfirmations uint64
	MailboxId             string
	MailboxDomain         uint32
}

// Helper methods

func (p *HyperlaneProvider) createMessageBody(msg CrossChainMessage) ([]byte, error) {
	// Create a structured message body for vault operations
	// This should follow the specific encoding format expected by the receiving chain

	// For now, create a simple encoding - in production this would likely use
	// protobuf or another standardized encoding format
	messageBody := fmt.Sprintf("{\"operation\":\"%s\",\"sender\":\"%s\",\"recipient\":\"%s\",\"amount\":\"%s\",\"timestamp\":%d}",
		p.messageTypeToString(msg.Type),
		msg.Sender.String(),
		msg.Recipient,
		msg.Amount.String(),
		time.Now().Unix(),
	)

	if len(msg.Data) > 0 {
		messageBody = fmt.Sprintf("{\"operation\":\"%s\",\"sender\":\"%s\",\"recipient\":\"%s\",\"amount\":\"%s\",\"data\":\"%x\",\"timestamp\":%d}",
			p.messageTypeToString(msg.Type),
			msg.Sender.String(),
			msg.Recipient,
			msg.Amount.String(),
			msg.Data,
			time.Now().Unix(),
		)
	}

	return []byte(messageBody), nil
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

// SetMailboxId updates the mailbox ID used by this provider
func (p *HyperlaneProvider) SetMailboxId(mailboxId util.HexAddress) {
	p.mailboxId = mailboxId
}

// GetMailboxId returns the current mailbox ID
func (p *HyperlaneProvider) GetMailboxId() util.HexAddress {
	return p.mailboxId
}

// CreateCollateralToken creates a new collateral token using the warp keeper
func (p *HyperlaneProvider) CreateCollateralToken(ctx sdk.Context, owner string, originMailbox util.HexAddress, denom string) (util.HexAddress, error) {
	if p.warpKeeper == nil {
		return util.HexAddress{}, fmt.Errorf("warp keeper not available")
	}

	// Create message for collateral token creation
	msg := &warptypes.MsgCreateCollateralToken{
		Owner:         owner,
		OriginMailbox: originMailbox,
		OriginDenom:   denom,
	}

	// Create collateral token through warp keeper
	tokenID, err := p.warpKeeper.CreateCollateralToken(ctx, msg)
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("failed to create collateral token: %w", err)
	}

	return tokenID, nil
}

// TransferRemoteCollateral initiates a remote collateral token transfer using the warp keeper
func (p *HyperlaneProvider) TransferRemoteCollateral(ctx sdk.Context, tokenId util.HexAddress, cosmosSender string, destinationDomain uint32, recipient util.HexAddress, amount math.Int, customHookId *util.HexAddress, gasLimit math.Int, maxFee sdk.Coin, customHookMetadata []byte) (util.HexAddress, error) {
	if p.warpKeeper == nil {
		return util.HexAddress{}, fmt.Errorf("warp keeper not available")
	}

	// Get the token from the warp keeper
	token, err := p.warpKeeper.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("failed to get token: %w", err)
	}

	// Initiate remote collateral transfer through warp keeper
	messageID, err := p.warpKeeper.RemoteTransferCollateral(ctx, token, cosmosSender, destinationDomain, recipient, amount, customHookId, gasLimit, maxFee, customHookMetadata)
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("failed to initiate remote collateral transfer: %w", err)
	}

	return messageID, nil
}

// TransferRemoteSynthetic initiates a remote synthetic token transfer using the warp keeper
func (p *HyperlaneProvider) TransferRemoteSynthetic(ctx sdk.Context, tokenId util.HexAddress, cosmosSender string, destinationDomain uint32, recipient util.HexAddress, amount math.Int, customHookId *util.HexAddress, gasLimit math.Int, maxFee sdk.Coin, customHookMetadata []byte) (util.HexAddress, error) {
	if p.warpKeeper == nil {
		return util.HexAddress{}, fmt.Errorf("warp keeper not available")
	}

	// Get the token from the warp keeper
	token, err := p.warpKeeper.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("failed to get token: %w", err)
	}

	// Initiate remote synthetic transfer through warp keeper
	messageID, err := p.warpKeeper.RemoteTransferSynthetic(ctx, token, cosmosSender, destinationDomain, recipient, amount, customHookId, gasLimit, maxFee, customHookMetadata)
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("failed to initiate remote synthetic transfer: %w", err)
	}

	return messageID, nil
}
