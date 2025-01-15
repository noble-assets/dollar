package mocks

import (
	"testing"

	dollar "dollar.noble.xyz"
	wormholekeeper "github.com/noble-assets/wormhole/keeper"

	storetypes "cosmossdk.io/store/types"
	"dollar.noble.xyz/keeper"
	"dollar.noble.xyz/types"
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
)

func DollarKeeperWithKeepers(t testing.TB, bank BankKeeper, account AccountKeeper) (*keeper.Keeper, sdk.Context) {
	key := storetypes.NewKVStoreKey(types.ModuleName)
	tkey := storetypes.NewTransientStoreKey("transient_authority")
	wrapper := testutil.DefaultContextWithDB(t, key, tkey)

	cfg := MakeTestEncodingConfig("noble")
	types.RegisterInterfaces(cfg.InterfaceRegistry)

	storeService := runtime.NewKVStoreService(key)
	headerService := runtime.ProvideHeaderInfoService(&runtime.AppBuilder{})
	eventService := runtime.ProvideEventService()
	addressCdc := address.NewBech32Codec("noble")

	wormholeKeeper := wormholekeeper.NewKeeper(
		cfg.Codec,
		storeService,
		headerService,
		eventService,
		addressCdc,
		nil,
		nil,
		nil,
	)

	k := keeper.NewKeeper(
		"uusdn",
		cfg.Codec,
		storeService,
		headerService,
		eventService,
		addressCdc,
		bank,
		account,
		wormholeKeeper,
	)

	dollar.InitGenesis(wrapper.Ctx, k, address.NewBech32Codec("noble"), *types.DefaultGenesisState())
	return k, wrapper.Ctx
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
