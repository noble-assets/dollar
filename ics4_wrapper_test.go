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
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/stretchr/testify/require"
)

var _ porttypes.ICS4Wrapper = (*MockICS4Wrapper)(nil)

type MockICS4Wrapper struct {
	t *testing.T
}

func (m MockICS4Wrapper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	return 0, nil
}

func (m MockICS4Wrapper) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet exported.PacketI,
	ack exported.Acknowledgement,
) error {
	m.t.Fatal("WriteAcknowledgement should not have been called")
	return nil
}

func (m MockICS4Wrapper) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	m.t.Fatal("GetAppVersion should not have been called")
	return "", false
}

type MockDollarKeeper struct {
	denom string
}

func (m MockDollarKeeper) GetDenom() string {
	return m.denom
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
			wrapper := MockICS4Wrapper{t}
			keeper := MockDollarKeeper{denom: denom}
			nobleWrapper := NewICS4Wrapper(wrapper, keeper)

			data, err := transfertypes.ModuleCdc.MarshalJSON(&tt.data)
			require.NoError(t, err)

			ctx := sdk.Context{}
			timeout := uint64(0)

			_, err = nobleWrapper.SendPacket(ctx, nil, "transfer", "channel-0", clienttypes.Height{}, timeout, data)

			if tt.fail {
				require.Error(t, err)
				require.ErrorContains(t, err, fmt.Sprintf("ibc transfers of %s are currently disabled", denom))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
