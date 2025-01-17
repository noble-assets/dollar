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
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	modulev1 "dollar.noble.xyz/api/module/v1"
	portalv1 "dollar.noble.xyz/api/portal/v1"
	dollarv1 "dollar.noble.xyz/api/v1"
	"dollar.noble.xyz/keeper"
	"dollar.noble.xyz/types"
	"dollar.noble.xyz/types/portal"
)

// ConsensusVersion defines the current Noble Dollar module consensus version.
const ConsensusVersion = 1

var (
	_ module.AppModuleBasic      = AppModule{}
	_ appmodule.AppModule        = AppModule{}
	_ module.HasConsensusVersion = AppModule{}
	_ module.HasGenesis          = AppModule{}
	_ module.HasGenesisBasics    = AppModuleBasic{}
	_ module.HasServices         = AppModule{}
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

	if err := portal.RegisterQueryHandlerClient(context.Background(), mux, portal.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genesis types.GenesisState
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
	var genesis types.GenesisState
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

	portal.RegisterMsgServer(cfg.MsgServer(), keeper.NewPortalMsgServer(m.keeper))
	portal.RegisterQueryServer(cfg.QueryServer(), keeper.NewPortalQueryServer(m.keeper))
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
							RpcMethod: "SetPeer",
							Use:       "set-peer [chain] [transceiver] [manager]",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "chain"},
								{ProtoField: "transceiver"},
								{ProtoField: "manager"},
							},
						},
						{
							RpcMethod: "Transfer",
							Use:       "transfer",
							Short:     "Transfer USDN from Noble and receive M cross-chain",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "chain"},
								{ProtoField: "recipient"},
								{ProtoField: "amount"},
							},
						},
					},
				},
			},
		},
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: dollarv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
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
							RpcMethod: "Peers",
							Use:       "peers",
						},
					},
				},
			},
		},
	}
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
	HeaderService header.Service
	EventService  event.Service

	Cdc            codec.Codec
	AddressCodec   address.Codec
	BankKeeper     types.BankKeeper
	WormholeKeeper portal.WormholeKeeper
}

type ModuleOutputs struct {
	depinject.Out

	Keeper       *keeper.Keeper
	Module       appmodule.AppModule
	Restrictions banktypes.SendRestrictionFn
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(
		in.Config.Denom,
		in.Cdc,
		in.StoreService,
		in.HeaderService,
		in.EventService,
		in.AddressCodec,
		in.BankKeeper,
		in.WormholeKeeper,
	)
	m := NewAppModule(in.AddressCodec, k)

	return ModuleOutputs{Keeper: k, Module: m, Restrictions: k.SendRestrictionFn}
}
