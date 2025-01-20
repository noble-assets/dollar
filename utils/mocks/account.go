package mocks

import (
	"context"

	"cosmossdk.io/core/address"
	"dollar.noble.xyz/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/codec"
)

var _ types.AccountKeeper = AccountKeeper{}

type AccountKeeper struct {
	Accounts      map[string]sdk.AccountI
	AccountNumber uint64
}

func (k AccountKeeper) SetModuleAccount(ctx context.Context, acc sdk.ModuleAccountI) {
	k.Accounts[acc.GetAddress().String()] = acc
}

func (k AccountKeeper) NewAccount(ctx context.Context, acc sdk.AccountI) sdk.AccountI {
	k.AccountNumber += 1
	if err := acc.SetAccountNumber(k.AccountNumber); err != nil {
		panic(err)
	}

	return acc
}

func (AccountKeeper) AddressCodec() address.Codec {
	return codec.NewBech32Codec("noble")
}

func (k AccountKeeper) GetAccount(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	// NOTE: The mock bankKeeper already sets the Bech32 prefix.
	return k.Accounts[addr.String()]
}
