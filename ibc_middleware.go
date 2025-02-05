package dollar

import (
	sdkerrors "cosmossdk.io/errors"
	"dollar.noble.xyz/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	chantypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ porttypes.Middleware = (*IBCMiddleware)(nil)

// IBCMiddleware implements the Middleware interface specified via ICS-30.
// The middleware checks outgoing IBC transfers to ensure that Noble Dollar cannot be sent via IBC to other chains.
type IBCMiddleware struct {
	app          porttypes.IBCModule
	ics4Wrapper  porttypes.ICS4Wrapper
	dollarKeeper ExpectedKeeper
}

// NewIBCMiddleware returns a new instance of IBCMiddleware.
func NewIBCMiddleware(
	app porttypes.IBCModule,
	ics4Wrapper porttypes.ICS4Wrapper,
	keeper ExpectedKeeper,
) IBCMiddleware {
	return IBCMiddleware{
		app:          app,
		ics4Wrapper:  ics4Wrapper,
		dollarKeeper: keeper,
	}
}

// ExpectedKeeper defines the methods the IBCMiddleware expects in order to block outgoing IBC transfers of the Noble Dollar.
type ExpectedKeeper interface {
	GetDenom() string
}

// OnChanOpenInit implements the IBCModule interface.
func (i IBCMiddleware) OnChanOpenInit(ctx sdk.Context, order chantypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty chantypes.Counterparty, version string) (string, error) {
	return i.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

// OnChanOpenTry implements the IBCModule interface.
func (i IBCMiddleware) OnChanOpenTry(ctx sdk.Context, order chantypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty chantypes.Counterparty, counterpartyVersion string) (version string, err error) {
	return i.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion)
}

// OnChanOpenAck implements the IBCModule interface.
func (i IBCMiddleware) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return i.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

// OnChanOpenConfirm implements the IBCModule interface.
func (i IBCMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return i.app.OnChanOpenConfirm(ctx, portID, channelID)
}

// OnChanCloseInit implements the IBCModule interface.
func (i IBCMiddleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return i.app.OnChanCloseInit(ctx, portID, channelID)
}

// OnChanCloseConfirm implements the IBCModule interface.
func (i IBCMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return i.app.OnChanCloseConfirm(ctx, portID, channelID)
}

// OnRecvPacket implements the IBCModule interface.
func (i IBCMiddleware) OnRecvPacket(ctx sdk.Context, packet chantypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	return i.app.OnRecvPacket(ctx, packet, relayer)
}

// OnAcknowledgementPacket implements the IBCModule interface.
func (i IBCMiddleware) OnAcknowledgementPacket(ctx sdk.Context, packet chantypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return i.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface.
func (i IBCMiddleware) OnTimeoutPacket(ctx sdk.Context, packet chantypes.Packet, relayer sdk.AccAddress) error {
	return i.app.OnTimeoutPacket(ctx, packet, relayer)
}

// SendPacket attempts to unmarshal the packet data into the ICS-20 FungibleTokenPacketData type. If the data is a
// ICS-20 transfer packet then a check is done on the denom to ensure that Noble Dollar cannot be transferred out
// of Noble via IBC.
func (i IBCMiddleware) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (sequence uint64, err error) {
	var packetData transfertypes.FungibleTokenPacketData

	if err := transfertypes.ModuleCdc.UnmarshalJSON(data, &packetData); err != nil {
		return i.ics4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// If the packet data denom matches the denom currently initialized in the Dollar module we want to return an error,
	// so that outgoing IBC transfers of Noble Dollar fail.
	denom := i.dollarKeeper.GetDenom()

	if packetData.Denom == denom {
		return 0, sdkerrors.Wrapf(types.ErrCannotSendViaIBC, "transfers of %s are currently disabled", denom)
	}

	return i.ics4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// WriteAcknowledgement implements the ICS4 Wrapper interface.
func (i IBCMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	return i.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

// GetAppVersion implements the ICS4 Wrapper interface.
func (i IBCMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return i.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}
