// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package dollar

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	chantypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/stretchr/testify/require"
)

var (
	_ porttypes.IBCModule   = (*MockIBCModule)(nil)
	_ porttypes.ICS4Wrapper = (*MockICS4Wrapper)(nil)
	_ DollarKeeper          = (*MockDollarKeeper)(nil)
)

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

type MockICS4Wrapper struct {
	t *testing.T
}

func (w MockICS4Wrapper) SendPacket(sdk.Context, *capabilitytypes.Capability, string, string, clienttypes.Height, uint64, []byte) (uint64, error) {
	return 0, nil
}

func (w MockICS4Wrapper) WriteAcknowledgement(sdk.Context, *capabilitytypes.Capability, exported.PacketI, exported.Acknowledgement) error {
	w.t.Fatal("WriteAcknowledgement should not have been called")
	return nil
}

func (w MockICS4Wrapper) GetAppVersion(sdk.Context, string, string) (string, bool) {
	w.t.Fatal("GetAppVersion should not have been called")
	return "", false
}

type MockDollarKeeper struct {
	denom string
}

func (k MockDollarKeeper) GetDenom() string {
	return k.denom
}

// TestSendPacket asserts that outgoing IBC transfers work as expected in cases
// where the denom is $USDN, as well as cases where the denom is not.
func TestSendPacket(t *testing.T) {
	denom := "uusdn"

	tc := []struct {
		name string
		data transfertypes.FungibleTokenPacketData
		fail bool
	}{
		{
			"Outgoing IBC transfer of USDN - should be blocked",
			transfertypes.NewFungibleTokenPacketData(denom, "1000000", "test", "test", "test"),
			true,
		},
		{
			"Outgoing IBC transfer of USDC - should not be blocked",
			transfertypes.NewFungibleTokenPacketData("uusdc", "1000000", "test", "test", "test"),
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
				require.ErrorContains(t, err, fmt.Sprintf("ibc transfers of %s are currently disabled", denom))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
