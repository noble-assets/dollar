package simapp

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	transferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/noble-assets/wormhole"
	wormholetypes "github.com/noble-assets/wormhole/types"
)

func (app *SimApp) RegisterLegacyModules() error {
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(capabilitytypes.StoreKey),
		storetypes.NewMemoryStoreKey(capabilitytypes.MemStoreKey),
		storetypes.NewKVStoreKey(exported.StoreKey),
		storetypes.NewKVStoreKey(transfertypes.StoreKey),
	); err != nil {
		return err
	}

	app.ParamsKeeper.Subspace(exported.ModuleName).WithKeyTable(clienttypes.ParamKeyTable().RegisterParamSet(&connectiontypes.Params{}))

	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		app.appCodec,
		app.GetKey(capabilitytypes.StoreKey),
		app.GetMemKey(capabilitytypes.MemStoreKey),
	)

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(exported.ModuleName)
	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec,
		app.GetKey(exported.StoreKey),
		app.GetSubspace(exported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
		"noble1s7evsmath5f3ef7vk97ru2tez9k5rs00klunzu",
	)

	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(transfertypes.ModuleName)
	app.TransferKeeper = transferkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(transfertypes.StoreKey),
		app.GetSubspace(transfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, // TODO(@john): Think about migrating NobleICS4Wrapper
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
		"noble1s7evsmath5f3ef7vk97ru2tez9k5rs00klunzu",
	)
	app.DollarKeeper.SetIBCKeepers(app.IBCKeeper.ChannelKeeper, app.TransferKeeper)

	scopedWormholeKeeper := app.CapabilityKeeper.ScopeToModule(wormholetypes.ModuleName)
	app.WormholeKeeper.SetIBCKeepers(app.IBCKeeper.ChannelKeeper, app.IBCKeeper.PortKeeper, scopedWormholeKeeper)

	router := porttypes.NewRouter()
	router.AddRoute(transfertypes.PortID, transfer.NewIBCModule(app.TransferKeeper))
	router.AddRoute(wormholetypes.ModuleName, wormhole.NewIBCModule(app.WormholeKeeper))
	app.IBCKeeper.SetRouter(router)

	return app.RegisterModules(
		capability.NewAppModule(app.appCodec, *app.CapabilityKeeper, true),
		ibc.NewAppModule(app.IBCKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		tendermint.NewAppModule(),
		solomachine.NewAppModule(),
	)
}
