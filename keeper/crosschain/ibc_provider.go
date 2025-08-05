package crosschain

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"

	dollarv2 "dollar.noble.xyz/v2/types/v2"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// IBCProvider implements CrossChainProvider for IBC
type IBCProvider struct {
	channelKeeper  IBCChannelKeeper
	transferKeeper IBCTransferKeeper
	clientKeeper   IBCClientKeeper
	portID         string
	defaultTimeout time.Duration
}

// IBCChannelKeeper defines the expected IBC channel keeper interface
type IBCChannelKeeper interface {
	GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool)
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, sourcePort, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (uint64, error)
}

// IBCTransferKeeper defines the expected IBC transfer keeper interface
type IBCTransferKeeper interface {
	Transfer(ctx sdk.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error)
	DenomTrace(ctx sdk.Context, denomTraceHash string) (ibctransfertypes.DenomTrace, bool)
}

// IBCClientKeeper defines the expected IBC client keeper interface
type IBCClientKeeper interface {
	GetClientState(ctx sdk.Context, clientID string) (exported.ClientState, bool)
	GetClientStatus(ctx sdk.Context, clientState exported.ClientState, clientID string) string
}

// NewIBCProvider creates a new IBC provider
func NewIBCProvider(
	channelKeeper IBCChannelKeeper,
	transferKeeper IBCTransferKeeper,
	clientKeeper IBCClientKeeper,
	portID string,
	defaultTimeout time.Duration,
) *IBCProvider {
	return &IBCProvider{
		channelKeeper:  channelKeeper,
		transferKeeper: transferKeeper,
		clientKeeper:   clientKeeper,
		portID:         portID,
		defaultTimeout: defaultTimeout,
	}
}

// GetProviderType returns the IBC provider type
func (p *IBCProvider) GetProviderType() dollarv2.Provider {
	return dollarv2.Provider_IBC
}

// SendMessage sends a cross-chain message via IBC
func (p *IBCProvider) SendMessage(ctx sdk.Context, route *vaultsv2.CrossChainRoute, msg CrossChainMessage) (*vaultsv2.ProviderTrackingInfo, error) {
	// Validate route has IBC config
	if route.ProviderConfig == nil || route.ProviderConfig.GetIbcConfig() == nil {
		return nil, fmt.Errorf("route %s missing IBC configuration", route.RouteId)
	}

	ibcConfig := route.ProviderConfig.GetIbcConfig()

	// Validate channel exists and is open
	channel, found := p.channelKeeper.GetChannel(ctx, p.portID, ibcConfig.ChannelId)
	if !found {
		return nil, fmt.Errorf("IBC channel %s not found", ibcConfig.ChannelId)
	}

	if channel.State != channeltypes.OPEN {
		return nil, fmt.Errorf("IBC channel %s is not open, state: %s", ibcConfig.ChannelId, channel.State.String())
	}

	// Get sequence number
	sequence, found := p.channelKeeper.GetNextSequenceSend(ctx, p.portID, ibcConfig.ChannelId)
	if !found {
		return nil, fmt.Errorf("failed to get next sequence for channel %s", ibcConfig.ChannelId)
	}

	// Create IBC transfer message
	transferMsg := &ibctransfertypes.MsgTransfer{
		SourcePort:    p.portID,
		SourceChannel: ibcConfig.ChannelId,
		Token: sdk.Coin{
			Denom:  "uusdc", // This should be configurable
			Amount: msg.Amount,
		},
		Sender:   msg.Sender.String(),
		Receiver: msg.Recipient,
		TimeoutHeight: clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: ibcConfig.TimeoutHeight,
		},
		TimeoutTimestamp: ibcConfig.TimeoutTimestamp,
		Memo:             p.createMemo(msg),
	}

	// Execute transfer
	_, err := p.transferKeeper.Transfer(ctx, transferMsg)
	if err != nil {
		return nil, fmt.Errorf("IBC transfer failed: %w", err)
	}

	// Create tracking info
	tracking := &vaultsv2.ProviderTrackingInfo{
		TrackingInfo: &vaultsv2.ProviderTrackingInfo_IbcTracking{
			IbcTracking: &vaultsv2.IBCTrackingInfo{
				Sequence:           sequence,
				SourceChannel:      ibcConfig.ChannelId,
				SourcePort:         p.portID,
				DestinationChannel: channel.Counterparty.ChannelId,
				DestinationPort:    channel.Counterparty.PortId,
				TimeoutTimestamp:   ibcConfig.TimeoutTimestamp,
				TimeoutHeight:      ibcConfig.TimeoutHeight,
				AckReceived:        false,
			},
		},
	}

	return tracking, nil
}

// GetMessageStatus checks the status of an IBC message
func (p *IBCProvider) GetMessageStatus(ctx sdk.Context, tracking *vaultsv2.ProviderTrackingInfo) (MessageStatus, error) {
	ibcTracking := tracking.GetIbcTracking()
	if ibcTracking == nil {
		return MessageStatusFailed, fmt.Errorf("invalid IBC tracking info")
	}

	// Check if acknowledgment was received
	if ibcTracking.AckReceived {
		// Parse acknowledgment to determine success/failure
		if len(ibcTracking.AckData) > 0 {
			// Simple acknowledgment parsing - in practice this would be more sophisticated
			if string(ibcTracking.AckData) == "success" {
				return MessageStatusConfirmed, nil
			} else {
				return MessageStatusFailed, nil
			}
		}
		return MessageStatusDelivered, nil
	}

	// Check if message has timed out
	if ctx.BlockTime().Unix() > int64(ibcTracking.TimeoutTimestamp) {
		return MessageStatusTimeout, nil
	}

	// Check channel state
	channel, found := p.channelKeeper.GetChannel(ctx, ibcTracking.SourcePort, ibcTracking.SourceChannel)
	if !found || channel.State != channeltypes.OPEN {
		return MessageStatusFailed, nil
	}

	return MessageStatusSent, nil
}

// GetConfirmations returns the number of confirmations for an IBC message
func (p *IBCProvider) GetConfirmations(ctx sdk.Context, tracking *vaultsv2.ProviderTrackingInfo) (uint64, error) {
	ibcTracking := tracking.GetIbcTracking()
	if ibcTracking == nil {
		return 0, fmt.Errorf("invalid IBC tracking info")
	}

	// For IBC, confirmations are based on acknowledgment
	if ibcTracking.AckReceived {
		return 1, nil // IBC has built-in finality
	}

	return 0, nil
}

// EstimateGas estimates gas cost for an IBC operation
func (p *IBCProvider) EstimateGas(ctx sdk.Context, route *vaultsv2.CrossChainRoute, msg CrossChainMessage) (uint64, math.Int, error) {
	// IBC gas estimation
	baseGas := uint64(200000) // Base gas for IBC transfer

	// Adjust based on message type
	switch msg.Type {
	case MessageTypeDeposit:
		baseGas += 50000
	case MessageTypeWithdraw:
		baseGas += 75000
	case MessageTypeUpdate:
		baseGas += 25000
	case MessageTypeLiquidate:
		baseGas += 100000
	}

	// Gas price estimation (this would typically come from a gas price oracle)
	gasPrice := math.NewInt(1000) // 0.001 USDC per gas unit

	totalCost := math.NewInt(int64(baseGas)).Mul(gasPrice)

	return baseGas, totalCost, nil
}

// ValidateConfig validates IBC-specific configuration
func (p *IBCProvider) ValidateConfig(config *vaultsv2.CrossChainProviderConfig) error {
	ibcConfig := config.GetIbcConfig()
	if ibcConfig == nil {
		return fmt.Errorf("IBC configuration is required")
	}

	if ibcConfig.ChannelId == "" {
		return fmt.Errorf("IBC channel ID is required")
	}

	if ibcConfig.PortId == "" {
		ibcConfig.PortId = "transfer" // Default to transfer port
	}

	if ibcConfig.TimeoutTimestamp == 0 && ibcConfig.TimeoutHeight == 0 {
		return fmt.Errorf("either timeout timestamp or timeout height must be set")
	}

	return nil
}

// UpdateMessageStatus updates the status of an IBC message (called by IBC callbacks)
func (p *IBCProvider) UpdateMessageStatus(
	ctx sdk.Context,
	tracking *vaultsv2.ProviderTrackingInfo,
	ackData []byte,
	success bool,
) error {
	ibcTracking := tracking.GetIbcTracking()
	if ibcTracking == nil {
		return fmt.Errorf("invalid IBC tracking info")
	}

	ibcTracking.AckReceived = true
	ibcTracking.AckData = ackData

	return nil
}

// GetChannelInfo returns information about an IBC channel
func (p *IBCProvider) GetChannelInfo(ctx sdk.Context, channelID string) (*IBCChannelInfo, error) {
	channel, found := p.channelKeeper.GetChannel(ctx, p.portID, channelID)
	if !found {
		return nil, fmt.Errorf("channel %s not found", channelID)
	}

	// Get client state
	_, clientState, err := p.channelKeeper.GetChannelClientState(ctx, p.portID, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client state: %w", err)
	}

	// Get client status
	clientStatus := p.clientKeeper.GetClientStatus(ctx, clientState, channel.ConnectionHops[0])

	return &IBCChannelInfo{
		ChannelID:    channelID,
		PortID:       p.portID,
		ConnectionID: channel.ConnectionHops[0],
		ClientID:     channel.ConnectionHops[0], // Simplified
		State:        channel.State.String(),
		ClientStatus: clientStatus,
		Counterparty: channel.Counterparty,
		Version:      channel.Version,
		Ordering:     channel.Ordering.String(),
	}, nil
}

// IBCChannelInfo contains information about an IBC channel
type IBCChannelInfo struct {
	ChannelID    string
	PortID       string
	ConnectionID string
	ClientID     string
	State        string
	ClientStatus string
	Counterparty channeltypes.Counterparty
	Version      string
	Ordering     string
}

// Helper methods

func (p *IBCProvider) createMemo(msg CrossChainMessage) string {
	// Create a memo that includes operation type and any relevant data
	memo := fmt.Sprintf("vault_op:%s", p.messageTypeToString(msg.Type))

	if len(msg.Data) > 0 {
		memo += fmt.Sprintf(",data:%x", msg.Data)
	}

	return memo
}

func (p *IBCProvider) messageTypeToString(msgType MessageType) string {
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

// Packet callback handlers (these would be called by IBC middleware)

// OnAcknowledgementPacket handles IBC packet acknowledgments
func (p *IBCProvider) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) error {
	// This would be called by IBC middleware when an acknowledgment is received
	// Implementation would update the tracking info and trigger any necessary callbacks
	return nil
}

// OnTimeoutPacket handles IBC packet timeouts
func (p *IBCProvider) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) error {
	// This would be called by IBC middleware when a packet times out
	// Implementation would update the tracking info and handle the timeout
	return nil
}
