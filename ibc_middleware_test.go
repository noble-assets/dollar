package dollar

import (
	"testing"

	"dollar.noble.xyz/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	chantypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/stretchr/testify/require"
)

var _ porttypes.ICS4Wrapper = (*MockICS4Wrapper)(nil)
var _ porttypes.IBCModule = (*MockIBCModule)(nil)

type MockICS4Wrapper struct {
	t *testing.T
}

func (m MockICS4Wrapper) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (sequence uint64, err error) {
	return 0, nil
}

func (m MockICS4Wrapper) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	m.t.Fatal("WriteAcknowledgement should not have been called")
	return nil
}

func (m MockICS4Wrapper) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	m.t.Fatal("GetAppVersion should not have been called")
	return "", false
}

type MockIBCModule struct {
	t *testing.T
}

func (m MockIBCModule) OnChanOpenInit(sdk.Context, chantypes.Order, []string, string, string, *capabilitytypes.Capability, chantypes.Counterparty, string) (string, error) {
	m.t.Fatal("OnChanOpenInit should not have been called")
	return "", nil
}

func (m MockIBCModule) OnChanOpenTry(sdk.Context, chantypes.Order, []string, string, string, *capabilitytypes.Capability, chantypes.Counterparty, string) (version string, err error) {
	m.t.Fatal("OnChanOpenTry should not have been called")
	return "", nil
}

func (m MockIBCModule) OnChanOpenAck(sdk.Context, string, string, string, string) error {
	m.t.Fatal("OnChanOpenAck should not have been called")
	return nil
}

func (m MockIBCModule) OnChanOpenConfirm(sdk.Context, string, string) error {
	m.t.Fatal("OnChanOpenConfirm should not have been called")
	return nil
}

func (m MockIBCModule) OnChanCloseInit(sdk.Context, string, string) error {
	m.t.Fatal("OnChanCloseInit should not have been called")
	return nil
}

func (m MockIBCModule) OnChanCloseConfirm(sdk.Context, string, string) error {
	m.t.Fatal("OnChanCloseConfirm should not have been called")
	return nil
}

func (m MockIBCModule) OnRecvPacket(sdk.Context, chantypes.Packet, sdk.AccAddress) exported.Acknowledgement {
	m.t.Fatal("OnRecvPacket should not have been called")
	return nil
}

func (m MockIBCModule) OnAcknowledgementPacket(sdk.Context, chantypes.Packet, []byte, sdk.AccAddress) error {
	m.t.Fatal("OnAcknowledgementPacket should not have been called")
	return nil
}

func (m MockIBCModule) OnTimeoutPacket(sdk.Context, chantypes.Packet, sdk.AccAddress) error {
	m.t.Fatal("OnTimeoutPacket should not have been called")
	return nil
}

type MockDollarKeeper struct {
	denom string
}

func (m MockDollarKeeper) GetDenom() string {
	return m.denom
}

// TestSendPacket asserts that outgoing IBC transfers work as intended in cases where the underlying assets denom is
// the Noble Dollar, as well as cases where the denom is not the Noble Dollar.
func TestSendPacket(t *testing.T) {
	denom := "uusdn"

	tc := []struct {
		name string
		data transfertypes.FungibleTokenPacketData
		fail bool
	}{
		{
			"Outgoing IBC transfer should be blocked - should fail",
			transfertypes.NewFungibleTokenPacketData(denom, "1000", "test", "test", "test"),
			true,
		},
		{
			"Outgoing IBC transfer should not be blocked - should not fail",
			transfertypes.NewFungibleTokenPacketData("uusdc", "1000", "test", "test", "test"),
			false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			module := MockIBCModule{t}
			wrapper := MockICS4Wrapper{t}
			keeper := MockDollarKeeper{denom: denom}
			middleware := NewIBCMiddleware(module, wrapper, keeper)

			data, err := transfertypes.ModuleCdc.MarshalJSON(&tt.data)
			require.NoError(t, err)

			ctx := sdk.Context{}
			timeout := uint64(0)

			_, err = middleware.SendPacket(ctx, nil, "transfer", "channel-0", clienttypes.Height{}, timeout, data)

			if tt.fail {
				require.Error(t, err)
				require.ErrorIs(t, err, types.ErrCannotSendViaIBC)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
