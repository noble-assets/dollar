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

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"dollar.noble.xyz/v3/types/v2"
)

var _ porttypes.IBCModule = &IBCModule{}

// IBCModule implements the porttypes.IBCModule interface. It implements custom
// logic in OnTimeoutPacket to increment the retry amount of Noble Dollar yield
// for external chains.
type IBCModule struct {
	underlying   porttypes.IBCModule
	dollarKeeper IBCModuleExpectedDollarKeeper
}

// IBCModuleExpectedDollarKeeper defines the interface expected by IBCModule for the Noble Dollar module.
type IBCModuleExpectedDollarKeeper interface {
	GetDenom() string
	IncrementRetryAmount(ctx context.Context, provider v2.Provider, identifier string, amount math.Int) error
}

// NewIBCModule returns a new instance of IBCModule.
func NewIBCModule(underlying porttypes.IBCModule, dollarKeeper IBCModuleExpectedDollarKeeper) porttypes.IBCModule {
	return IBCModule{
		underlying:   underlying,
		dollarKeeper: dollarKeeper,
	}
}

// AddAmountToRetry is a utility that adds an amount to retry on the next yield
// distribution. The provided packet must have either timed out or incurred an
// error acknowledgement.
func (m IBCModule) AddAmountToRetry(ctx sdk.Context, packet channeltypes.Packet) {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data); err != nil {
		return
	}

	denom := m.dollarKeeper.GetDenom()
	if data.Denom == denom {
		escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
		if data.Sender == escrowAddress.String() {
			amount, ok := math.NewIntFromString(data.Amount)
			if !ok {
				return
			}

			err := m.dollarKeeper.IncrementRetryAmount(ctx, v2.Provider_IBC, packet.SourceChannel, amount)
			if err != nil {
				return
			}
		}
	}

	return
}

func (m IBCModule) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return m.underlying.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (m IBCModule) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string) (version string, err error) {
	return m.underlying.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion)
}

func (m IBCModule) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return m.underlying.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (m IBCModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return m.underlying.OnChanOpenConfirm(ctx, portID, channelID)
}

func (m IBCModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return m.underlying.OnChanCloseInit(ctx, portID, channelID)
}

func (m IBCModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return m.underlying.OnChanCloseConfirm(ctx, portID, channelID)
}

func (m IBCModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	return m.underlying.OnRecvPacket(ctx, packet, relayer)
}

func (m IBCModule) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return m.underlying.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	if !ack.Success() {
		m.AddAmountToRetry(ctx, packet)
	}

	return m.underlying.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (m IBCModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	m.AddAmountToRetry(ctx, packet)

	return m.underlying.OnTimeoutPacket(ctx, packet, relayer)
}
