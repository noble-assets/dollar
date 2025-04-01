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
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"dollar.noble.xyz/v2/types/v2"
)

var _ porttypes.ICS4Wrapper = &ICS4Wrapper{}

// ICS4Wrapper implements the porttypes.ICS4Wrapper interface. It implements
// custom logic in SendPacket in order to check all outgoing IBC transfers so
// that $USDN cannot be sent to another chain.
type ICS4Wrapper struct {
	ics4Wrapper  porttypes.ICS4Wrapper
	dollarKeeper ExpectedDollarKeeper
}

// ExpectedDollarKeeper defines the interface expected by ICS4Wrapper for the Noble Dollar module.
type ExpectedDollarKeeper interface {
	GetDenom() string
	HasYieldRecipient(ctx context.Context, provider v2.Provider, identifier string) bool
}

// NewICS4Wrapper returns a new instance of ICS4Wrapper.
func NewICS4Wrapper(app porttypes.ICS4Wrapper, dollarKeeper ExpectedDollarKeeper) porttypes.ICS4Wrapper {
	return ICS4Wrapper{
		ics4Wrapper:  app,
		dollarKeeper: dollarKeeper,
	}
}

// SendPacket attempts to unmarshal the provided packet data as the ICS-20
// FungibleTokenPacketData type. If the packet is a valid ICS-20 transfer, then
// a check is performed on the denom to ensure that $USDN cannot be transferred
// out of Noble via IBC.
func (w ICS4Wrapper) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (sequence uint64, err error) {
	var packetData transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(data, &packetData); err != nil {
		return w.ics4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	denom := w.dollarKeeper.GetDenom()
	if packetData.Denom == denom {
		enabled := w.dollarKeeper.HasYieldRecipient(ctx, v2.Provider_IBC, sourceChannel)
		if !enabled {
			return 0, fmt.Errorf("ibc transfers of %s are currently disabled on %s", denom, sourceChannel)
		}
	}

	return w.ics4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// WriteAcknowledgement implements the porttypes.ICS4Wrapper interface.
func (w ICS4Wrapper) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI, ack ibcexported.Acknowledgement) error {
	return w.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

// GetAppVersion implements the porttypes.ICS4Wrapper interface.
func (w ICS4Wrapper) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return w.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}
