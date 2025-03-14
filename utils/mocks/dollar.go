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

package mocks

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	wormholekeeper "github.com/noble-assets/wormhole/keeper"
	wormholetypes "github.com/noble-assets/wormhole/types"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"dollar.noble.xyz/v2"
	"dollar.noble.xyz/v2/keeper"
	"dollar.noble.xyz/v2/types"
)

func DollarKeeperWithKeepers(t testing.TB, bank BankKeeper, account AccountKeeper) (*keeper.Keeper, *wormholekeeper.Keeper, sdk.Context) {
	keys := storetypes.NewKVStoreKeys(types.ModuleName, wormholetypes.ModuleName)
	ctx := testutil.DefaultContextWithKeys(keys, nil, nil)

	cfg := MakeTestEncodingConfig("noble")
	types.RegisterInterfaces(cfg.InterfaceRegistry)

	headerService := runtime.ProvideHeaderInfoService(&runtime.AppBuilder{})
	eventService := runtime.ProvideEventService()
	addressCdc := address.NewBech32Codec("noble")

	wormholeKeeper := wormholekeeper.NewKeeper(
		cfg.Codec,
		runtime.NewKVStoreService(keys[wormholetypes.ModuleName]),
		headerService,
		eventService,
		addressCdc,
		nil,
		nil,
		nil,
	)

	k := keeper.NewKeeper(
		"uusdn",
		"authority",
		1e6,
		1e6,
		cfg.Codec,
		runtime.NewKVStoreService(keys[types.ModuleName]),
		log.NewTestLogger(t),
		headerService,
		eventService,
		addressCdc,
		account,
		bank,
		nil,
		nil,
		wormholeKeeper,
	)

	bank = bank.WithSendCoinsRestriction(k.SendRestrictionFn)
	k.SetBankKeeper(bank)

	wormholeKeeper.Config.Set(ctx, wormholetypes.Config{
		ChainId:    uint16(vaautils.ChainIDNoble),
		GovChain:   uint16(vaautils.GovernanceChain),
		GovAddress: vaautils.GovernanceEmitter.Bytes(),
	})
	dollar.InitGenesis(ctx, k, addressCdc, *types.DefaultGenesisState())

	return k, wormholeKeeper, ctx
}

// MakeTestEncodingConfig is a modified testutil.MakeTestEncodingConfig that
// sets a custom Bech32 prefix in the interface registry.
func MakeTestEncodingConfig(prefix string, modules ...module.AppModuleBasic) moduletestutil.TestEncodingConfig {
	aminoCodec := codec.NewLegacyAmino()
	interfaceRegistry := codectestutil.CodecOptions{
		AccAddressPrefix: prefix,
	}.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)

	encCfg := moduletestutil.TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             aminoCodec,
	}

	mb := module.NewBasicManager(modules...)

	std.RegisterLegacyAminoCodec(encCfg.Amino)
	std.RegisterInterfaces(encCfg.InterfaceRegistry)
	mb.RegisterLegacyAminoCodec(encCfg.Amino)
	mb.RegisterInterfaces(encCfg.InterfaceRegistry)

	return encCfg
}
