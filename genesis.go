package dollar

import (
	"context"
	"fmt"

	"cosmossdk.io/core/address"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	"dollar.noble.xyz/keeper"
	"dollar.noble.xyz/types"
	"dollar.noble.xyz/types/portal"
)

func InitGenesis(ctx context.Context, k *keeper.Keeper, address address.Codec, genesis types.GenesisState) {
	var err error

	err = k.Index.Set(ctx, genesis.Index)
	if err != nil {
		panic(errors.Wrap(err, "unable to set genesis index"))
	}

	for rawAccount, rawPrincipal := range genesis.Principal {
		account, err := address.StringToBytes(rawAccount)
		if err != nil {
			panic(errors.Wrapf(err, "unable to decode account %s", rawAccount))
		}

		principal, ok := math.NewIntFromString(rawPrincipal)
		if !ok {
			panic(fmt.Errorf("unable to parse principal %s", rawPrincipal))
		}

		err = k.Principal.Set(ctx, account, principal)
		if err != nil {
			panic(errors.Wrapf(err, "unable to set genesis principal (%s:%s)", rawAccount, rawPrincipal))
		}
	}

	if err = k.Owner.Set(ctx, genesis.Portal.Owner); err != nil {
		panic(errors.Wrap(err, "unable to set genesis owner"))
	}

	for chain, peer := range genesis.Portal.Peers {
		err = k.Peers.Set(ctx, chain, peer)
		if err != nil {
			panic(errors.Wrapf(err, "unable to set genesis peer (%d:%s)", chain, peer))
		}
	}
}

func ExportGenesis(ctx context.Context, k *keeper.Keeper) *types.GenesisState {
	index, _ := k.Index.Get(ctx)
	principal, _ := k.GetPrincipal(ctx)

	owner, _ := k.Owner.Get(ctx)
	peers, _ := k.GetPeers(ctx)

	return &types.GenesisState{
		Portal: portal.GenesisState{
			Owner: owner,
			Peers: peers,
		},
		Index:     index,
		Principal: principal,
	}
}
