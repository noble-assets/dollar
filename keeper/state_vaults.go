package keeper

import (
	"context"

	"cosmossdk.io/math"

	"dollar.noble.xyz/types/vaults"
)

func (k *Keeper) GetTotalFlexiblePrincipal(ctx context.Context) (math.Int, error) {
	value, err := k.TotalFlexiblePrincipal.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	return value, nil
}

func (k *Keeper) GetPaused(ctx context.Context) vaults.PausedType {
	value, err := k.Paused.Get(ctx)
	if err != nil {
		return vaults.NONE
	}
	return vaults.PausedType(value)
}

func (k *Keeper) GetPositions(ctx context.Context) ([]vaults.PositionEntry, error) {
	var positions []vaults.PositionEntry

	itr, err := k.Positions.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		position, _ := k.Positions.Get(ctx, key)
		positions = append(positions, vaults.PositionEntry{
			Address:   key.K1(),
			Vault:     vaults.VaultType(key.K2()),
			Index:     position.Index,
			Principal: position.Principal,
			Amount:    position.Amount,
			Time:      position.Time,
		})
	}

	return positions, err
}

func (k *Keeper) GetPositionsByProvider(ctx context.Context, provider []byte) ([]vaults.PositionEntry, error) {
	var positions []vaults.PositionEntry

	itr, err := k.Positions.Indexes.ByProvider.MatchExact(ctx, provider)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.PrimaryKey()
		position, _ := k.Positions.Get(ctx, key)
		positions = append(positions, vaults.PositionEntry{
			Address:   key.K1(),
			Vault:     vaults.VaultType(key.K2()),
			Index:     position.Index,
			Principal: position.Principal,
			Amount:    position.Amount,
			Time:      position.Time,
		})
	}

	return positions, err
}

func (k *Keeper) GetRewards(ctx context.Context) ([]vaults.Reward, error) {
	var rewards []vaults.Reward

	itr, err := k.Rewards.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		reward, _ := k.Rewards.Get(ctx, key)
		rewards = append(rewards, vaults.Reward{
			Index:   reward.Index,
			Total:   reward.Total,
			Rewards: reward.Rewards,
		})
	}

	return rewards, err
}
