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
	"encoding/json"
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	warpkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/warp/keeper"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	modulev1 "dollar.noble.xyz/v2/api/module/v1"
	portalv1 "dollar.noble.xyz/v2/api/portal/v1"
	dollarv1 "dollar.noble.xyz/v2/api/v1"
	vaultsv1 "dollar.noble.xyz/v2/api/vaults/v1"
	"dollar.noble.xyz/v2/client/cli"
	"dollar.noble.xyz/v2/keeper"
	"dollar.noble.xyz/v2/types"
	"dollar.noble.xyz/v2/types/portal"
	"dollar.noble.xyz/v2/types/v2"
	"dollar.noble.xyz/v2/types/vaults"
)

// ConsensusVersion defines the current Noble Dollar module consensus version.
const ConsensusVersion = 2

var (
	_ module.AppModuleBasic      = AppModule{}
	_ appmodule.AppModule        = AppModule{}
	_ module.HasConsensusVersion = AppModule{}
	_ module.HasGenesis          = AppModule{}
	_ module.HasGenesisBasics    = AppModuleBasic{}
	_ module.HasServices         = AppModule{}
	_ appmodule.HasBeginBlocker  = AppModule{}
)

//

type AppModuleBasic struct {
	addressCodec address.Codec
}

func NewAppModuleBasic(addressCodec address.Codec) AppModuleBasic {
	return AppModuleBasic{addressCodec: addressCodec}
}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
	if err := v2.RegisterQueryHandlerClient(context.Background(), mux, v2.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}

	if err := portal.RegisterQueryHandlerClient(context.Background(), mux, portal.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}

	if err := vaults.RegisterQueryHandlerClient(context.Background(), mux, vaults.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(v2.DefaultGenesisState())
}

func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genesis v2.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genesis); err != nil {
		return fmt.Errorf("failed to unmarshal Noble Dollar genesis state: %w", err)
	}

	return genesis.Validate()
}

//

type AppModule struct {
	AppModuleBasic

	keeper *keeper.Keeper
}

func NewAppModule(addressCodec address.Codec, keeper *keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(addressCodec),
		keeper:         keeper,
	}
}

func (AppModule) IsOnePerModuleType() {}

func (AppModule) IsAppModule() {}

func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

func (m AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var genesis v2.GenesisState
	cdc.MustUnmarshalJSON(bz, &genesis)

	InitGenesis(ctx, m.keeper, m.addressCodec, genesis)
}

func (m AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesis := ExportGenesis(ctx, m.keeper)
	return cdc.MustMarshalJSON(genesis)
}

func (m AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(m.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(m.keeper))
	v2.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerV2(m.keeper))
	v2.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerV2(m.keeper))

	portal.RegisterMsgServer(cfg.MsgServer(), keeper.NewPortalMsgServer(m.keeper))
	portal.RegisterQueryServer(cfg.QueryServer(), keeper.NewPortalQueryServer(m.keeper))

	vaults.RegisterMsgServer(cfg.MsgServer(), keeper.NewVaultsMsgServer(m.keeper))
	vaults.RegisterQueryServer(cfg.QueryServer(), keeper.NewVaultsQueryServer(m.keeper))

	migrator := keeper.NewMigrator(m.keeper)
	if err := cfg.RegisterMigration(types.ModuleName, 1, migrator.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate Noble Dollar from version 1 to 2: %v", err))
	}
}

func (m AppModule) BeginBlock(ctx context.Context) error {
	return m.keeper.BeginBlocker(ctx)
}

//

func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: dollarv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "ClaimYield",
					Use:       "claim-yield",
				},
				{
					RpcMethod: "SetPausedState",
					Use:       "set-paused-state [paused]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "paused"},
					},
				},
			},
			SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
				"portal": {
					Service: portalv1.Msg_ServiceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod:      "Deliver",
							Use:            "deliver [vaa]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "vaa"}},
						},
						{
							RpcMethod: "Transfer",
							Use:       "transfer [amount] [destination-chain-id] [destination-token] [recipient]",
							Short:     "Transfer USDN from Noble and receive M cross-chain",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "amount"},
								{ProtoField: "destination_chain_id"},
								{ProtoField: "destination_token"},
								{ProtoField: "recipient"},
							},
						},
						{
							RpcMethod: "SetPausedState",
							Use:       "set-paused-state [paused]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "paused"},
							},
						},
						{
							RpcMethod: "SetPeer",
							Use:       "set-peer [chain] [transceiver] [manager]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "chain"},
								{ProtoField: "transceiver"},
								{ProtoField: "manager"},
							},
						},
						{
							RpcMethod: "SetBridgingPath",
							Use:       "set-bridging-path [destination-chain-id] [destination-token] [supported]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "destination_chain_id"},
								{ProtoField: "destination_token"},
								{ProtoField: "supported"},
							},
						},
						{
							RpcMethod:      "TransferOwnership",
							Use:            "transfer-ownership [new-owner]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "new_owner"}},
						},
					},
				},
				"vaults": {
					Service: vaultsv1.Msg_ServiceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod: "Lock",
							Use:       "lock [vault] [amount]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "vault"},
								{ProtoField: "amount"},
							},
						},
						{
							RpcMethod: "Unlock",
							Use:       "unlock [vault] [amount]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "vault"},
								{ProtoField: "amount"},
							},
						},
						{
							RpcMethod: "SetPausedState",
							Use:       "set-paused-state [paused]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "paused"},
							},
						},
					},
				},
			},
			EnhanceCustomCommand: true,
		},
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: dollarv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Index",
					Use:       "index",
				},
				{
					RpcMethod:      "Principal",
					Use:            "principal [account]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "account"}},
				},
				{
					RpcMethod:      "Yield",
					Use:            "yield [account]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "account"}},
				},
				{
					// NOTE: This is skipped as it is overridden by the v2 command.
					RpcMethod: "Stats",
					Skip:      true,
				},
			},
			SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
				"portal": {
					Service: portalv1.Query_ServiceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod: "Owner",
							Use:       "owner",
						},
						{
							RpcMethod: "Paused",
							Use:       "paused",
						},
						{
							RpcMethod: "Peers",
							Use:       "peers",
						},
						{
							RpcMethod:      "DestinationTokens",
							Use:            "destination-tokens [chain-id]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "chain_id"}},
						},
						{
							RpcMethod: "Nonce",
							Use:       "nonce",
						},
					},
				},
				"vaults": {
					Service: vaultsv1.Query_ServiceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod: "Paused",
							Use:       "paused",
							Short:     "Retrieves the current pausing state of the Vault module",
						},
						{
							RpcMethod: "PendingRewards",
							Use:       "pending-rewards",
							Short:     "Retrieves the total amount of rewards pending distribution",
						},
						{
							RpcMethod:      "PendingRewardsByProvider",
							Use:            "pending-rewards-by-provider [provider]",
							Short:          "Retrieves the total amount of pending rewards for a specified provider",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "provider"}},
						},
						{
							RpcMethod:      "PositionsByProvider",
							Use:            "positions-by-provider [provider]",
							Short:          "List Vaults positions by a specific provider",
							Long:           "Retrieves all the active Vaults positions attributed to provider",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "provider"}},
						},
						{
							RpcMethod: "Stats",
							Use:       "stats",
						},
					},
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}

func (AppModule) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

func (AppModule) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config        *modulev1.Module
	StoreService  store.KVStoreService
	Logger        log.Logger
	HeaderService header.Service
	EventService  event.Service

	Cdc            codec.Codec
	AddressCodec   address.Codec
	BankKeeper     types.BankKeeper
	AccountKeeper  types.AccountKeeper
	WarpKeeper     warpkeeper.Keeper
	WormholeKeeper portal.WormholeKeeper
}

type ModuleOutputs struct {
	depinject.Out

	Keeper       *keeper.Keeper
	Module       appmodule.AppModule
	Restrictions banktypes.SendRestrictionFn
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	if in.Config.Authority == "" {
		panic("authority for Noble Dollar module must be set")
	}

	if in.Config.VaultsSeasonOneEndTimestamp == 0 {
		panic("vaults season one end timestamp must be set")
	}

	if in.Config.VaultsSeasonTwoYieldCollector == "" {
		panic("vaults season two yield collector must be set")
	}
	// Get the address bytes of the Collector address.
	vaultsSeasonTwoYieldCollectorAddress, err := in.AddressCodec.StringToBytes(in.Config.VaultsSeasonTwoYieldCollector)
	if err != nil {
		panic("vaults season two yield collector must be a valid address")
	}

	authority := authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	k := keeper.NewKeeper(
		in.Config.Denom,
		authority.String(),
		in.Config.VaultsMinimumLock,
		in.Config.VaultsMinimumUnlock,
		in.Config.VaultsSeasonOneEndTimestamp,
		vaultsSeasonTwoYieldCollectorAddress,
		in.Cdc,
		in.StoreService,
		in.Logger,
		in.HeaderService,
		in.EventService,
		in.AddressCodec,
		in.AccountKeeper,
		in.BankKeeper,
		nil,
		nil,
		&in.WarpKeeper,
		in.WormholeKeeper,
	)
	m := NewAppModule(in.AddressCodec, k)

	return ModuleOutputs{Keeper: k, Module: m, Restrictions: k.SendRestrictionFn}
}
