package simapp

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	"github.com/noble-assets/wormhole"
	"github.com/noble-assets/wormhole/types"
)

func (app *SimApp) RegisterLegacyModules() error {
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(capabilitytypes.StoreKey),
		storetypes.NewMemoryStoreKey(capabilitytypes.MemStoreKey),
		storetypes.NewKVStoreKey(exported.StoreKey),
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

	scopedWormholeKeeper := app.CapabilityKeeper.ScopeToModule(types.ModuleName)
	app.WormholeKeeper.SetIBCKeepers(app.IBCKeeper.ChannelKeeper, app.IBCKeeper.PortKeeper, scopedWormholeKeeper)

	router := porttypes.NewRouter()
	router.AddRoute(types.ModuleName, wormhole.NewIBCModule(app.WormholeKeeper))
	app.IBCKeeper.SetRouter(router)

	return app.RegisterModules(
		capability.NewAppModule(app.appCodec, *app.CapabilityKeeper, true),
		ibc.NewAppModule(app.IBCKeeper),
		tendermint.NewAppModule(),
		solomachine.NewAppModule(),
	)
}
