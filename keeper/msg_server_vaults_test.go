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

package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/core/header"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	vaultsv1 "dollar.noble.xyz/api/vaults/v1"
	"dollar.noble.xyz/keeper"
	"dollar.noble.xyz/types"
	"dollar.noble.xyz/types/vaults"
	"dollar.noble.xyz/utils"
	"dollar.noble.xyz/utils/mocks"
)

const ONE = 1_000_000

func TestPausing(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob := utils.TestAccount()
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	assert.Equal(t, vaults.NONE, k.GetVaultsPaused(ctx))

	// ARRANGE: Bob mints 100 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(100*ONE), nil)

	// ACT: Bob deposits 50 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Attempt to Pause with an invalid authority.
	_, err = vaultsServer.SetPausedState(ctx, &vaults.MsgSetPausedState{
		Signer: bob.Address,
		Paused: vaults.ALL,
	})
	assert.Error(t, err)

	// ACT: Pause ALL actions.
	_, err = vaultsServer.SetPausedState(ctx, &vaults.MsgSetPausedState{
		Signer: "authority",
		Paused: vaults.ALL,
	})
	assert.NoError(t, err)
	assert.Equal(t, vaults.ALL, k.GetVaultsPaused(ctx))

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 1, 0, 0, 0, time.UTC)})

	// ACT: Bob deposits 50 USDN into the Staked Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.Error(t, err)

	// ACT: Bob withdraws everything from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.Error(t, err)

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 2, 0, 0, 0, time.UTC)})

	// ACT: Pause only Unlock actions.
	_, err = vaultsServer.SetPausedState(ctx, &vaults.MsgSetPausedState{
		Signer: "authority",
		Paused: vaults.UNLOCK,
	})
	assert.NoError(t, err)
	assert.Equal(t, vaults.UNLOCK, k.GetVaultsPaused(ctx))

	// ACT: Bob deposits 50 USDN into the Staked Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob withdraws everything from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(100 * ONE),
	})
	assert.Error(t, err)

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 3, 0, 0, 0, time.UTC)})

	// ACT: Pause only Lock actions.
	_, err = vaultsServer.SetPausedState(ctx, &vaults.MsgSetPausedState{
		Signer: "authority",
		Paused: vaults.LOCK,
	})
	assert.NoError(t, err)
	assert.Equal(t, vaults.LOCK, k.GetVaultsPaused(ctx))

	// ACT: Bob withdraws everything from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob deposits 50 USDN into the Staked Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.Error(t, err)
}

func TestStakedVault(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	server := keeper.NewMsgServer(k)
	vaultsServer := keeper.NewVaultsMsgServer(k)
	vaultsQueryServer := keeper.NewVaultsQueryServer(k)
	bob := utils.TestAccount()
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	// ARRANGE: Bob mints 100 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(100*ONE), nil)

	// ACT: Bob deposits 50 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn"), math.NewInt(50*ONE)) // 50 USDN.

	// ASSERT: Matching Vaults Stats state.
	stats, _ := vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(50*ONE))

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ACT: Bob claims the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected to increase by a factor of 1.1.
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn"), math.NewInt(55*ONE))

	// ACT: Bob attempts to withdraw from the Staked Vault with an invalid amount.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
	})
	assert.Error(t, err)

	// ACT: Bob withdraws everything from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(0))
	assert.Equal(t, stats.StakedTotalPrincipal, math.ZeroInt())

	// ASSERT: Bob receives back the deposited amount.
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn"), math.NewInt(105*ONE))

	// ACT: Bob attempts to claim the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob does not have any yield to claim.
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn"), math.NewInt(105*ONE))

	// ARRANGE: Increase the index from 1.1 to 1.21 (~10%).
	_ = k.UpdateIndex(ctx, 1.21e12)

	// ACT: Bob claims the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected to increase by the yield.
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn").ToLegacyDec().TruncateInt(), math.NewInt(115499999))
}

func TestStakedVaultMultiPositions(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	server := keeper.NewMsgServer(k)
	vaultsServer := keeper.NewVaultsMsgServer(k)
	vaultsQueryServer := keeper.NewVaultsQueryServer(k)
	bob := utils.TestAccount()
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	// ARRANGE: Bob mints 100 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(100*ONE), nil)

	// ACT: Bob deposits 50 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn"), math.NewInt(50*ONE)) // 50 USDN.
	// ASSERT: Matching Vaults Stats state.
	stats, _ := vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(50*ONE))

	// ACT: Bob attempts deposits 50 USDN into the Staked Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	// ASSERT: Should've failed to same block execution.
	assert.Error(t, err)

	// ARRANGE: Increase block time.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 1, 0, 0, 0, time.UTC)})

	// ACT: Bob deposits 50 USDN into the Staked Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Matching Positions state.
	bobPositions, err := k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.NoError(t, err)
	assert.Len(t, bobPositions, 2)
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn"), math.NewInt(0*ONE))
	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(100*ONE))

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ACT: Bob withdraws everything from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(100 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn").ToLegacyDec().TruncateInt(), math.NewInt(100*ONE))
	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(0))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(0))

	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob does not have any yield or rewards to claim.
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn").ToLegacyDec().TruncateInt(), math.NewInt(100*ONE))
}

func TestStakedPartialRemoval(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	server := keeper.NewMsgServer(k)
	vaultsServer := keeper.NewVaultsMsgServer(k)
	vaultsQueryServer := keeper.NewVaultsQueryServer(k)
	bob := utils.TestAccount()

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ARRANGE: Bob mints 50 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(50*ONE), nil)

	// ACT: Bob deposits 50 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.STAKED,
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Matching Vaults Stats state.
	stats, _ := vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(45454545))

	// ARRANGE: Increase the index from 1.1 to 1.21 (~10%).
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.21e12)

	// ARRANGE: Bob mints other 50 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(50*ONE), nil)

	// ACT: Bob deposits other 50 USDN into the Staked Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.STAKED,
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Matching state.
	positions, _ := k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.Equal(t, 2, len(positions))
	assert.Equal(t, []vaults.PositionEntry{
		{
			Address:   bob.Bytes,
			Vault:     vaults.STAKED,
			Principal: math.NewInt(45454545),
			Index:     math.LegacyMustNewDecFromStr("1.1"),
			Amount:    math.NewInt(50 * ONE),
			Time:      time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC),
		},
		{
			Address:   bob.Bytes,
			Vault:     vaults.STAKED,
			Principal: math.NewInt(41322314),
			Index:     math.LegacyMustNewDecFromStr("1.21"),
			Amount:    math.NewInt(50 * ONE),
			Time:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}, positions)
	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(45454545+41322314))

	// ACT: Bob withdraws 10 USDN (partial first position) from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.STAKED,
		Amount: math.NewInt(10 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: matching state.
	assert.Equal(t, math.NewInt(10*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	positions, _ = k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.Equal(t, 2, len(positions))
	assert.Equal(t, []vaults.PositionEntry{
		{
			Address:   bob.Bytes,
			Vault:     vaults.STAKED,
			Principal: math.NewInt(36363636), // reduced (50-10)/1,1
			Index:     math.LegacyMustNewDecFromStr("1.1"),
			Amount:    math.NewInt(40 * ONE), // reduced
			Time:      time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC),
		},
		{
			Address:   bob.Bytes,
			Vault:     vaults.STAKED,
			Principal: math.NewInt(41322314),
			Index:     math.LegacyMustNewDecFromStr("1.21"),
			Amount:    math.NewInt(50 * ONE),
			Time:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}, positions)
	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(36363636+41322314))

	// ACT: Bob withdraws other 40 USDN (completes first position) from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.STAKED,
		Amount: math.NewInt(40 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Matching state.
	assert.Equal(t, math.NewInt(50*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	positions, _ = k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.Equal(t, 1, len(positions))
	assert.Equal(t, []vaults.PositionEntry{
		{
			Address:   bob.Bytes,
			Vault:     vaults.STAKED,
			Principal: math.NewInt(41322314),
			Index:     math.LegacyMustNewDecFromStr("1.21"),
			Amount:    math.NewInt(50 * ONE),
			Time:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}, positions)
	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(41322314))

	// ACT: Bob withdraws other 50 USDN from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.STAKED,
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Matching state.
	assert.Equal(t, math.NewInt(100*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	positions, _ = k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.Equal(t, 0, len(positions))
	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(0))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(0))

	// ACT: Bob claims the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob does not have any yield or rewards to claim.
	assert.Equal(t, math.NewInt(100*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
}

func TestStakedVaultRewardsMigration(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	server := keeper.NewMsgServer(k)
	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob := utils.TestAccount()

	// ARRANGE: Set the default index to 1.0 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	// ARRANGE: Bob mints 100 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(100*ONE), nil)

	// ACT: Bob deposits 50 USDN (half balance) into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(50*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	_ = k.UpdateIndex(ctx, 1.1e12)
	// ASSERT: Flexible vault balance is expected to increase by the yield.
	assert.Equal(t, math.NewInt(5*ONE), bank.Balances[vaults.FlexibleVaultAddress.String()].AmountOf("uusdn"))

	// ACT: Bob claims the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected to increase by the yield.
	assert.Equal(t, math.NewInt(55*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	// ASSERT: Flexible vault balance is expected to be the same.
	assert.Equal(t, math.NewInt(5*ONE), bank.Balances[vaults.FlexibleVaultAddress.String()].AmountOf("uusdn"))

	// ACT: Bob withdraws 50 USDN (total) from the Staked Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected to increase by only the deposited amount.
	assert.Equal(t, math.NewInt(105*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ARRANGE: Increase the index from 1.1 to 1.21 (~10%).
	_ = k.UpdateIndex(ctx, 1.21e12)
	assert.Equal(t, math.NewInt(5499999), bank.Balances[vaults.FlexibleVaultAddress.String()].AmountOf("uusdn"))

	// ACT: Bob claims the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(115499999), bank.Balances[bob.Address].AmountOf("uusdn"))
	assert.Equal(t, math.NewInt(0), bank.Balances[vaults.StakedVaultAddress.String()].AmountOf("uusdn"))
	assert.Equal(t, math.NewInt(5499999), bank.Balances[vaults.FlexibleVaultAddress.String()].AmountOf("uusdn")) // no change

	// ACT: Bob deposits 1 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Matching state.
	totalFlexiblePrincipal, _ := k.VaultsTotalFlexiblePrincipal.Get(ctx)
	assert.Equal(t, math.NewInt(826446), totalFlexiblePrincipal)
}

func TestFlexibleVaultMultiUser(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	server := keeper.NewMsgServer(k)
	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Set the default index to 1.0 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	// ARRANGE: Bob mints 1050 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1050*ONE), nil)
	// ARRANGE: Alice mints 50 USDN.
	_ = k.Mint(ctx, alice.Bytes, math.NewInt(50*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob deposits 50 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Alice deposits 50 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected to be empty.
	assert.Equal(t, math.NewInt(0*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	// ASSERT: Alice balance is expected to be empty.
	assert.Equal(t, math.NewInt(0*ONE), bank.Balances[alice.Address].AmountOf("uusdn"))

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ASSERT: Matching Staked Vault balance.
	assert.Equal(t, math.NewInt(1000*ONE), bank.Balances[vaults.StakedVaultAddress.String()].AmountOf("uusdn"))
	// ASSERT: Matching Flexible Vault balance.
	assert.Equal(t, math.NewInt((100)*ONE), bank.Balances[vaults.FlexibleVaultAddress.String()].AmountOf("uusdn"))

	// ASSERT: Matching Principal state.
	stakedPrincipal, _ := k.Principal.Get(ctx, vaults.StakedVaultAddress)
	flexiblePrincipal, _ := k.Principal.Get(ctx, vaults.FlexibleVaultAddress)
	assert.Equal(t, math.NewInt(909090910), stakedPrincipal)
	assert.Equal(t, math.NewInt(90909090), flexiblePrincipal)

	// ACT: Bob withdraws 50 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Alice withdraws 50 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected to increase by a factor of 1.1.
	assert.Equal(t, math.NewInt(55*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	// ASSERT: Alice balance is expected to increase by a factor of 1.1.
	assert.Equal(t, math.NewInt(55*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ACT: Bob attempts to claim the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ACT: Alice attempts to claim the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: alice.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt(55*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	// ASSERT: Alice balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt(55*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
}

func TestFlexibleVaultMultiUserEarlyExit(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	server := keeper.NewMsgServer(k)
	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Set the default index to 1.0 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	// ARRANGE: Bob mints 1050 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1050*ONE), nil)
	// ARRANGE: Alice mints 50 USDN.
	_ = k.Mint(ctx, alice.Bytes, math.NewInt(50*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob deposits 50 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Alice deposits 50 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected to be empty.
	assert.Equal(t, math.NewInt(0*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	// ASSERT: Alice balance is expected to be empty.
	assert.Equal(t, math.NewInt(0*ONE), bank.Balances[alice.Address].AmountOf("uusdn"))

	// ACT: Bob withdraws 50 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ACT: Alice withdraws 50 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(50 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected be as in the initial state.
	assert.Equal(t, math.NewInt(50*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	// ASSERT: Alice balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt(50*ONE), bank.Balances[alice.Address].AmountOf("uusdn"))

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	_ = k.UpdateIndex(ctx, 1.1e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})

	// ACT: Bob attempts to claim the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ACT: Alice attempts to claim the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: alice.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt((50*1.1)*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt((50*1.1)*ONE), bank.Balances[alice.Address].AmountOf("uusdn"))
}

func TestFlexibleVaultMultiUserEarlyExitCase2(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob := utils.TestAccount()

	// ARRANGE: Set the default index to 1.1 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ARRANGE: Bob mints 1000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Bob balance is expected be as in the initial state.
	assert.Equal(t, math.NewInt(1000*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
}

func TestFlexibleVaultMultiUserEarlyExitCase3(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	server := keeper.NewMsgServer(k)
	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob := utils.TestAccount()

	// ARRANGE: Set the default index to 1.1 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ARRANGE: Bob mints 1000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected be as in the initial state.
	assert.Equal(t, math.NewInt(1000*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ARRANGE: Increase the index from 1.1 to 1.21 (~10%).
	_ = k.UpdateIndex(ctx, 1.21e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})

	// ACT: Bob claims the yield.
	_, err = server.ClaimYield(ctx, &types.MsgClaimYield{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt((1100)*ONE-1), bank.Balances[bob.Address].AmountOf("uusdn"))
}

func TestFlexibleVaultBaseLockUnlock(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob := utils.TestAccount()

	// ARRANGE: Set the default index to 1.0 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	// ARRANGE: Bob mints 1000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt(1000*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ASSERT: Flexible Vault balance is expected to be empty.
	assert.Equal(t, math.NewInt(0), bank.Balances[vaults.FlexibleVaultAddress.String()].AmountOf("uusdn"))
}

func TestFlexibleVaultSimpleNoRewards(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob := utils.TestAccount()

	// ARRANGE: Set the default index to 1.0 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	// ARRANGE: Bob mints 1000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	_ = k.UpdateIndex(ctx, 1.1e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})

	// ACT: Bob withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt(1100*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))
	assert.Equal(t, math.NewInt(0), bank.Balances[vaults.FlexibleVaultAddress.String()].AmountOf("uusdn"))
}

func TestFlexibleVaultMultiUserFlexibleNoRewards(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Set the default index to 1.0 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1e12)

	// ARRANGE: Bob mints 1000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	_ = k.UpdateIndex(ctx, 1.1e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})

	// ARRANGE: Alice mints 1000 USDN.
	_ = k.Mint(ctx, alice.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Alice deposits 1000 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(0), bank.Balances[alice.Address].AmountOf("uusdn"))

	// ACT: Bob attempts to withdraw from the Flexible Vault with an invalid amount.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(1100*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ACT: Alice withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt(1100*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ASSERT: Alice balance is expected be as in the initial state.
	assert.Equal(t, math.NewInt(1000*ONE), bank.Balances[alice.Address].AmountOf("uusdn"))

	// ASSERT: Flexible Vault balance is expected be empty.
	assert.Equal(t, math.NewInt(0), bank.Balances[vaults.FlexibleVaultAddress.String()].AmountOf("uusdn"))
}

func TestFlexibleVaultMultiUserMultiEntry(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Set the default index to 1.01 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.01e12)

	// ACT: Bob deposits 2000 USDN into the Flexible Vault.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(2000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	_ = k.UpdateIndex(ctx, 1.1e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})

	// ARRANGE: Increase the index from 1.1 to 1.21 (~10%).
	_ = k.UpdateIndex(ctx, 1.21e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)})

	// ARRANGE: Increase the index from 1.21 to 1.33 (~10%).
	_ = k.UpdateIndex(ctx, 1.33e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)})

	// ARRANGE: Alice mints 1000 USDN.
	_ = k.Mint(ctx, alice.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Alice deposits 1000 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ARRANGE: Increase the index from 1.33 to 1.46 (~10%).
	_ = k.UpdateIndex(ctx, 1.46e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC)})

	// ASSERT: Matching Rewards state.
	rewards, err := k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(990099009),
			Rewards: math.NewInt(89108909),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.NewInt(990099009),
			Rewards: math.NewInt(108910891),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.33"),
			Total:   math.NewInt(990099009),
			Rewards: math.NewInt(118811881),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.46"),
			Total:   math.NewInt(1741978708),
			Rewards: math.NewInt(128712872),
		},
	}, rewards)

	// ASSERT: Matching Positions state.
	bobPositions, err := k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.NoError(t, err)
	alicePositions, err := k.GetVaultsPositionsByProvider(ctx, alice.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(bobPositions))
	assert.Equal(t, 1, len(alicePositions))
	assert.Equal(t, vaults.PositionEntry{
		Address:   bob.Bytes,
		Vault:     vaults.FLEXIBLE,
		Principal: math.LegacyNewDec(1000 * ONE).Quo(math.LegacyMustNewDecFromStr("1.01")).TruncateInt(),
		Index:     math.LegacyMustNewDecFromStr("1.01"),
		Amount:    math.NewInt(1000 * ONE),
		Time:      time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC),
	}, bobPositions[1])
	assert.Equal(t, vaults.PositionEntry{
		Address:   alice.Bytes,
		Vault:     vaults.FLEXIBLE,
		Principal: math.LegacyNewDec(1000 * ONE).Quo(math.LegacyMustNewDecFromStr("1.33")).TruncateInt(),
		Index:     math.LegacyMustNewDecFromStr("1.33"),
		Amount:    math.NewInt(1000 * ONE),
		Time:      time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
	}, alicePositions[0])

	// ACT: Bob withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Matching Rewards state.
	rewards, err = k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.33"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.46"),
			Total:   math.NewInt(751879699),
			Rewards: math.NewInt(128712872), // bob exited too early
		},
	}, rewards)

	// ASSERT: Matching Positions state.
	bobPositions, err = k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bobPositions))
	alicePositions, err = k.GetVaultsPositionsByProvider(ctx, alice.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(alicePositions))

	// ACT: Alice withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Matching Rewards state.
	rewards, err = k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.33"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.46"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(128712872), // bob and alice exited too early
		},
	}, rewards)

	// ASSERT: Matching Positions state.
	alicePositions, err = k.GetVaultsPositionsByProvider(ctx, alice.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(alicePositions))

	// ASSERT: Bob balance is expected be as in the initial state + standard yield + boosted yield. (1000/1,0*1,46)[yield] + 330[rewards] = ~ 1762
	assert.Equal(t, math.NewInt(1762376234), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ASSERT: Alice balance is expected be as in the initial state + standard yield. (1000/1,33*1,46)[yield] + 0[rewards] = ~ 1153
	assert.Equal(t, math.NewInt(1097744360), bank.Balances[alice.Address].AmountOf("uusdn"))
}

func TestFlexibleVaultRewardsSimple(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	vaultsQueryServer := keeper.NewVaultsQueryServer(k)
	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Set the default index to 1.01 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.01e12)

	// ARRANGE: Bob mints 1000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Matching Vaults Stats state.
	stats, _ := vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(990099009))

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ARRANGE: Bob mints 1000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(1000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Matching Positions state.
	bobPositions, err := k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(bobPositions))
	assert.Equal(t, vaults.PositionEntry{
		Address:   bob.Bytes,
		Vault:     vaults.FLEXIBLE,
		Principal: math.LegacyNewDec(1000 * ONE).Quo(math.LegacyMustNewDecFromStr("1.1")).TruncateInt(),
		Index:     math.LegacyMustNewDecFromStr("1.1"),
		Amount:    math.NewInt(1000 * ONE),
		Time:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}, bobPositions[1])
	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(990099009))
	assert.Equal(t, stats.FlexibleTotalUsers, uint64(1))
	assert.Equal(t, stats.FlexibleTotalPrincipal, math.NewInt(909090909))

	// ARRANGE: Increase the index from 1.1 to 1.21 (~10%).
	_ = k.UpdateIndex(ctx, 1.21e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)})

	// ARRANGE: Alice mints 9000 USDN.
	_ = k.Mint(ctx, alice.Bytes, math.NewInt(9000*ONE), nil)

	// ACT: Alice deposits 9000 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(9000 * ONE),
	})
	assert.NoError(t, err)
	// ASSERT: Matching Positions state.
	alicePositions, err := k.GetVaultsPositionsByProvider(ctx, alice.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(alicePositions))
	assert.Equal(t, vaults.PositionEntry{
		Address:   alice.Bytes,
		Vault:     vaults.FLEXIBLE,
		Principal: math.LegacyNewDec(9000 * ONE).Quo(math.LegacyMustNewDecFromStr("1.21")).TruncateInt(),
		Index:     math.LegacyMustNewDecFromStr("1.21"),
		Amount:    math.NewInt(9000 * ONE),
		Time:      time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
	}, alicePositions[0])
	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(990099009))
	assert.Equal(t, stats.FlexibleTotalUsers, uint64(2))
	assert.Equal(t, stats.FlexibleTotalPrincipal, math.NewInt(8347107437))

	// ASSERT: Matching Rewards state.
	rewards, err := k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(89108909),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.LegacyNewDec(1000 * ONE).Quo(math.LegacyMustNewDecFromStr("1.1")).TruncateInt(), // no alice yet
			Rewards: math.NewInt(108910891),
		},
	}, rewards)

	// ARRANGE: Increase the index from 1.21 to 1.33 (~10%).
	_ = k.UpdateIndex(ctx, 1.33e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)})

	// ACT: Bob withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Matching Positions state.
	bobPositions, err = k.GetVaultsPositionsByProvider(ctx, bob.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bobPositions))

	// ASSERT: Matching Rewards state.
	rewards, err = k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(89108909), // unclaimed, bob entered too late
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.33"),
			Total:   math.NewInt(7438016528),
			Rewards: math.NewInt(118811881), // unclaimed, bob exited too early
		},
	}, rewards)

	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(990099009))
	assert.Equal(t, stats.FlexibleTotalUsers, uint64(1))
	assert.Equal(t, stats.FlexibleTotalPrincipal, math.NewInt(7438016528))
	assert.Equal(t, stats.FlexibleTotalDistributedRewardsPrincipal, math.NewInt(108910891))

	// ACT: Alice withdraws 9000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(9000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Matching Positions state.
	alicePositions, err = k.GetVaultsPositionsByProvider(ctx, alice.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(alicePositions))

	// ASSERT: Matching Rewards state.
	rewards, err = k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(89108909), // unclaimed, bob entered too late
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.33"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(118811881), // unclaimed, bob & alice exited too early
		},
	}, rewards)

	// ASSERT: Matching Vaults Stats state.
	stats, _ = vaultsQueryServer.Stats(ctx, &vaults.QueryStats{})
	assert.Equal(t, stats.StakedTotalUsers, uint64(1))
	assert.Equal(t, stats.StakedTotalPrincipal, math.NewInt(990099009))
	assert.Equal(t, stats.FlexibleTotalUsers, uint64(0))
	assert.Equal(t, stats.FlexibleTotalPrincipal, math.NewInt(0))
	assert.Equal(t, stats.FlexibleTotalDistributedRewardsPrincipal, math.NewInt(108910891+0))

	// ASSERT: Bob balance is expected be as in the initial state + standard yield + boosted yield. (1000/1,1*1,33)[yield] + 110[rewards] = ~1318
	assert.Equal(t, math.NewInt(1318001799), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ASSERT: Alice balance is expected be as in the initial state + standard yield + boosted yield. (9000/1,21*1,33)[yield] + 0[rewards] = ~ 9892
	assert.Equal(t, math.NewInt(9892561982), bank.Balances[alice.Address].AmountOf("uusdn"))
}

func TestFlexibleVaultRewardsHacky(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Set the default index to 1.01 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.01e12)

	// ARRANGE: Bob mints 2000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(2000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Staking Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ARRANGE: Alice mints 9000 USDN.
	_ = k.Mint(ctx, alice.Bytes, math.NewInt(9000*ONE), nil)

	// ACT: Alice deposits 9000 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(9000 * ONE),
	})
	assert.NoError(t, err)

	// ARRANGE: Increase the index from 1.1 to 1.21 (~10%).
	_ = k.UpdateIndex(ctx, 1.21e12)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)})

	// ASSERT: Matching Rewards state.
	rewards, err := k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(990099009),
			Rewards: math.NewInt(89108909),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.NewInt(9171917190),
			Rewards: math.NewInt(108910891),
		},
	}, rewards)

	// ACT: Bob withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Alice withdraws 9000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: alice.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(9000 * ONE),
	})
	assert.NoError(t, err)

	// ASSERT: Rewards Positions state.
	rewards, err = k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(108910891),
		},
	}, rewards)

	// ASSERT: Bob balance is expected be as in the initial state + standard yield. (1000/1,0*1,21)[yield] + 100[rewards] = ~1287
	assert.Equal(t, math.NewInt(1287128709), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ASSERT: Alice balance is expected be as in the initial state + standard yield. (9000/1,1*1,21)[yield] + 0[rewards] = 9900
	assert.Equal(t, math.NewInt(9899999999), bank.Balances[alice.Address].AmountOf("uusdn"))
}

func TestFlexibleVaultRewardsEarlyExit(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.DollarKeeperWithKeepers(t, bank, account)
	bank.Restriction = k.SendRestrictionFn
	k.SetBankKeeper(bank)

	vaultsServer := keeper.NewVaultsMsgServer(k)
	bob := utils.TestAccount()

	// ARRANGE: Set the default index to 1.01 .
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.01e12)

	// ARRANGE: Bob mints 2000 USDN.
	_ = k.Mint(ctx, bob.Bytes, math.NewInt(2000*ONE), nil)

	// ACT: Bob deposits 1000 USDN into the Staked Vault.
	_, err := vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_STAKED),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ARRANGE: Increase the index from 1.0 to 1.1 (~10%).
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.1e12)

	// ARRANGE: Increase the index from 1.1 to 1.21 (~10%).
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.21e12)

	// ACT: Bob deposits 1000 USDN into the Flexible Vault.
	_, err = vaultsServer.Lock(ctx, &vaults.MsgLock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ACT: Bob withdraws 1000 USDN from the Flexible Vault.
	_, err = vaultsServer.Unlock(ctx, &vaults.MsgUnlock{
		Signer: bob.Address,
		Vault:  vaults.VaultType(vaultsv1.VaultType_FLEXIBLE),
		Amount: math.NewInt(1000 * ONE),
	})
	assert.NoError(t, err)

	// ARRANGE: Increase the index from 1.21 to 1.33 (~10%).
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)})
	_ = k.UpdateIndex(ctx, 1.33e12)

	// ASSERT: Bob balance is expected be as in the initial state + standard yield.
	assert.Equal(t, math.NewInt(1000*ONE), bank.Balances[bob.Address].AmountOf("uusdn"))

	// ASSERT: Matching Rewards state.
	rewards, err := k.GetVaultsRewards(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []vaults.Reward{
		{
			Index:   math.LegacyMustNewDecFromStr("1.01"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(0),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.1"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(89108909),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.21"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(108910891),
		},
		{
			Index:   math.LegacyMustNewDecFromStr("1.33"),
			Total:   math.NewInt(0),
			Rewards: math.NewInt(118811881),
		},
	}, rewards)
}
