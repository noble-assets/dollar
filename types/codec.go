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

package types

import (
	"dollar.noble.xyz/v3/types/portal"
	"dollar.noble.xyz/v3/types/v2"
	"dollar.noble.xyz/v3/types/vaults"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	portal.RegisterLegacyAminoCodec(cdc)
	vaults.RegisterLegacyAminoCodec(cdc)

	cdc.RegisterConcrete(&MsgClaimYield{}, "dollar/ClaimYield", nil)
	cdc.RegisterConcrete(&MsgSetPausedState{}, "dollar/SetPausedState", nil)

	v2.RegisterLegacyAminoCodec(cdc)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	portal.RegisterInterfaces(registry)
	vaults.RegisterInterfaces(registry)

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimYield{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSetPausedState{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	v2.RegisterInterfaces(registry)
}

var amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
