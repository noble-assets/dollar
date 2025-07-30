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

package keeper

import (
	"context"
	"sort"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	hyperlaneutil "github.com/bcp-innovations/hyperlane-cosmos/util"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"google.golang.org/protobuf/runtime/protoiface"

	"dollar.noble.xyz/v2/types"
	"dollar.noble.xyz/v2/types/v2"
	"dollar.noble.xyz/v2/types/vaults"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) ClaimYield(ctx context.Context, msg *types.MsgClaimYield) (*types.MsgClaimYieldResponse, error) {
	if k.GetPaused(ctx) {
		return nil, types.ErrPaused
	}

	yield, account, err := k.GetYield(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.YieldName, account, sdk.NewCoins(sdk.NewCoin(k.denom, yield)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to distribute yield to user")
	}

	return &types.MsgClaimYieldResponse{}, k.event.EventManager(ctx).Emit(ctx, &types.YieldClaimed{
		Account: msg.Signer,
		Amount:  yield,
	})
}

func (k msgServer) SetPausedState(ctx context.Context, msg *types.MsgSetPausedState) (*types.MsgSetPausedStateResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != k.authority {
		return nil, errors.Wrapf(vaults.ErrInvalidAuthority, "expected %s, got %s", k.authority, msg.Signer)
	}

	if err := k.Paused.Set(ctx, msg.Paused); err != nil {
		return nil, err
	}

	event := protoiface.MessageV1(&types.Unpaused{})
	if msg.Paused {
		event = &types.Paused{}
	}

	return &types.MsgSetPausedStateResponse{}, k.event.EventManager(ctx).Emit(ctx, event)
}

func (k *Keeper) Burn(ctx context.Context, sender []byte, amount math.Int) error {
	coins := sdk.NewCoins(sdk.NewCoin(k.denom, amount))
	err := k.bank.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins)
	if err != nil {
		return err
	}
	err = k.bank.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) Mint(ctx context.Context, recipient []byte, amount math.Int, index *int64) error {
	if index != nil {
		_ = k.UpdateIndex(ctx, *index)
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.denom, amount))
	err := k.bank.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}
	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) UpdateIndex(ctx context.Context, index int64) error {
	oldIndex, err := k.Index.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get index from state")
	}
	if index <= oldIndex {
		return types.ErrDecreasingIndex
	}

	err = k.Index.Set(ctx, index)
	if err != nil {
		return errors.Wrap(err, "unable to set index in state")
	}

	totalPrincipal, err := k.GetTotalPrincipal(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get total principal from state")
	}

	currentSupply := k.bank.GetSupply(ctx, k.denom).Amount
	expectedSupply := k.GetPresentAmount(totalPrincipal, index)

	coins := sdk.NewCoins(sdk.NewCoin(k.denom, expectedSupply.Sub(currentSupply)))
	yield := math.ZeroInt()
	if coins.IsAllPositive() {
		err = k.bank.MintCoins(ctx, types.ModuleName, coins)
		if err != nil {
			return errors.Wrap(err, "unable to mint coins")
		}
		err = k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.YieldName, coins)
		if err != nil {
			return errors.Wrap(err, "unable to send coins")
		}
		yield = coins.AmountOf(k.denom)
		err = k.IncrementTotalYieldAccrued(ctx, yield)
		if err != nil {
			return errors.Wrap(err, "unable to increment total yield accrued")
		}
	}

	// Handle the vaults yield logic for Season One only if it has not ended.
	if !k.IsVaultsSeasonOneEnded(ctx) {
		if err := k.handleVaultsYieldSeasonOne(ctx, index); err != nil {
			return err
		}
	} else {
		// Handle the vaults yield logic for Season Two.
		if err := k.handleVaultsYieldSeasonTwo(ctx); err != nil {
			return err
		}
	}

	// Claim and transfer the yield of ibc external chains.
	err = k.claimExternalYieldIBC(ctx)
	if err != nil {
		return err
	}

	// Claim and transfer the yield of hyperlane external chains.
	err = k.claimExternalYieldHyperlane(ctx)
	if err != nil {
		return err
	}

	return k.event.EventManager(ctx).Emit(ctx, &types.IndexUpdated{
		OldIndex:       oldIndex,
		NewIndex:       index,
		TotalPrincipal: totalPrincipal,
		YieldAccrued:   yield,
	})
}

func (k *Keeper) claimModuleYield(ctx context.Context, addr sdk.Address) (math.Int, error) {
	// Get the Yield of the module address.
	yield, _, err := k.GetYield(ctx, addr.String())
	if err != nil {
		return math.ZeroInt(), nil
	}

	// If the Yield exists, claim it.
	if yield.IsPositive() {
		err = k.bank.SendCoinsFromModuleToAccount(ctx, types.YieldName, addr.Bytes(), sdk.NewCoins(sdk.NewCoin(k.denom, yield)))
		if err != nil {
			return math.ZeroInt(), err
		}
	}

	return yield, nil
}

func (k *Keeper) claimStakedVaultYield(ctx context.Context) (math.Int, error) {
	// Get the Yield of the Staked Vault.
	yield, _, err := k.GetYield(ctx, vaults.StakedVaultAddress.String())
	if err != nil {
		return math.ZeroInt(), nil
	}
	if !yield.IsPositive() {
		return math.ZeroInt(), nil
	}

	// Redirect the Yield from the Yield Module to the Flexible Vault instead of the Staked Vault.
	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.YieldName, vaults.FlexibleVaultAddress, sdk.NewCoins(sdk.NewCoin(k.denom, yield)))
	if err != nil {
		return math.ZeroInt(), err
	}

	// Get the current Index.
	rawIndex, err := k.Index.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	index := math.LegacyNewDec(rawIndex).QuoInt64(1e12)
	// Calculate the Yield Principal.
	yieldPrincipal := yield.ToLegacyDec().Quo(index).TruncateInt()

	// Reduce the Staked Vault Principal by the Yield Principal.
	stakedPrincipal, err := k.Principal.Get(ctx, vaults.StakedVaultAddress)
	if err != nil {
		stakedPrincipal = math.ZeroInt()
	}
	if err = k.Principal.Set(ctx, vaults.StakedVaultAddress, stakedPrincipal.Sub(yieldPrincipal)); err != nil {
		return math.ZeroInt(), err
	}

	// Add the Yield Principal to the Flexible Vault Principal.
	flexiblePrincipal, err := k.Principal.Get(ctx, vaults.FlexibleVaultAddress)
	if err != nil {
		flexiblePrincipal = math.ZeroInt()
	}
	if err = k.Principal.Set(ctx, vaults.FlexibleVaultAddress, flexiblePrincipal.Add(yieldPrincipal)); err != nil {
		return math.ZeroInt(), err
	}
	return yield, nil
}

func (k *Keeper) claimExternalYieldIBC(ctx context.Context) error {
	provider := v2.Provider_IBC
	yieldRecipients, err := k.GetYieldRecipientsByProvider(ctx, provider)
	if err != nil {
		return errors.Wrap(err, "unable to get ibc yield recipients from state")
	}

	channelIds := make([]string, 0, len(yieldRecipients))
	for channelId := range yieldRecipients {
		channelIds = append(channelIds, channelId)
	}
	sort.Strings(channelIds)

	for _, channelId := range channelIds {
		yieldRecipient := yieldRecipients[channelId]

		escrowAddress := transfertypes.GetEscrowAddress(transfertypes.PortID, channelId)
		yield, err := k.claimModuleYield(ctx, escrowAddress)
		if err != nil {
			return errors.Wrapf(err, "unable to claim yield for %s/%s", provider, channelId)
		}
		retryAmount, err := k.GetRetryAmountAndRemove(ctx, provider, channelId)
		if err != nil {
			return errors.Wrapf(err, "unable to get and remove retry amount for %s/%s", provider, channelId)
		}
		accruedYield := yield.Add(retryAmount)
		if !accruedYield.IsPositive() {
			continue
		}

		timeout := uint64(k.header.GetHeaderInfo(ctx).Time.UnixNano()) + transfertypes.DefaultRelativePacketTimeoutTimestamp
		_, transferErr := k.transfer.Transfer(ctx, &transfertypes.MsgTransfer{
			SourcePort:       transfertypes.PortID,
			SourceChannel:    channelId,
			Token:            sdk.NewCoin(k.denom, accruedYield),
			Sender:           escrowAddress.String(),
			Receiver:         yieldRecipient,
			TimeoutHeight:    clienttypes.ZeroHeight(),
			TimeoutTimestamp: timeout,
			Memo:             "",
		})
		if transferErr != nil {
			k.logger.Error("unable to transfer ibc yield", "identifier", channelId, "err", transferErr)

			err = k.IncrementRetryAmount(ctx, provider, channelId, accruedYield)
			if err != nil {
				return errors.Wrapf(err, "unable to increment retry amount for %s/%s", provider, channelId)
			}
		}

		err = k.IncrementTotalExternalYield(ctx, provider, channelId, yield)
		if err != nil {
			return errors.Wrapf(err, "unable to increment total yield for %s/%s", provider, channelId)
		}

		if transferErr == nil {
			k.logger.Info("claimed and transferred ibc yield", "amount", accruedYield, "identifier", channelId)
		}
	}

	return nil
}

func (k *Keeper) claimExternalYieldHyperlane(ctx context.Context) error {
	provider := v2.Provider_HYPERLANE
	yieldRecipients, err := k.GetYieldRecipientsByProvider(ctx, provider)
	if err != nil {
		return errors.Wrap(err, "unable to get hyperlane yield recipients from state")
	}

	identifiers := make([]string, 0, len(yieldRecipients))
	for identifier := range yieldRecipients {
		identifiers = append(identifiers, identifier)
	}
	sort.Strings(identifiers)

	address := authtypes.NewModuleAddress(warptypes.ModuleName)
	yield, err := k.claimModuleYield(ctx, address)
	if err != nil {
		return errors.Wrapf(err, "unable to claim yield for hyperlane")
	}
	if !yield.IsPositive() {
		return nil
	}

	// NOTE: We iterate over the yield recipients twice to first calculate the
	// total collateral across all supported routes. This is done so that we
	// can safely calculate the yield portion of each route.

	totalCollateral := math.ZeroInt()
	tokens := make(map[string]warptypes.HypToken)
	for _, identifier := range identifiers {
		rawIdentifier, err := hyperlaneutil.DecodeHexAddress(identifier)
		if err != nil {
			return errors.Wrap(err, "unable to decode hyperlane identifier")
		}
		tokenId := rawIdentifier.GetInternalId()

		token, err := k.warp.HypTokens.Get(ctx, tokenId)
		if err != nil {
			return errors.Wrap(err, "unable to get hyperlane token from state")
		}

		totalCollateral = totalCollateral.Add(token.CollateralBalance)
		tokens[identifier] = token
	}

	if !totalCollateral.IsPositive() {
		return nil
	}

	for _, identifier := range identifiers {
		yieldRecipient := yieldRecipients[identifier]

		token := tokens[identifier]
		collateral := token.CollateralBalance
		collateralPortion := math.LegacyNewDecFromInt(collateral).QuoInt(totalCollateral)
		yieldPortion := collateralPortion.MulInt(yield).TruncateInt()
		retryAmount, err := k.GetRetryAmountAndRemove(ctx, provider, identifier)
		if err != nil {
			return errors.Wrapf(err, "unable to get and remove retry amount for %s/%s", provider, identifier)
		}
		accruedYield := yieldPortion.Add(retryAmount)
		if !accruedYield.IsPositive() {
			continue
		}

		router, err := k.getHyperlaneRouter(ctx, token.Id.GetInternalId())
		if err != nil {
			return err
		}

		yieldRecipientBz, err := hyperlaneutil.DecodeHexAddress(yieldRecipient)
		if err != nil {
			return errors.Wrap(err, "unable to decode hyperlane yield recipient")
		}

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		_, transferErr := k.warp.RemoteTransferCollateral(
			sdkCtx,
			token,
			address.String(),
			router.ReceiverDomain,
			yieldRecipientBz,
			accruedYield,
			nil,
			math.ZeroInt(),
			sdk.NewCoin(k.denom, math.ZeroInt()),
			nil,
		)
		if transferErr != nil {
			k.logger.Error("unable to transfer hyperlane yield", "identifier", identifier, "err", transferErr)

			err = k.IncrementRetryAmount(ctx, provider, identifier, accruedYield)
			if err != nil {
				return errors.Wrapf(err, "unable to increment retry amount for %s/%s", provider, identifier)
			}
		}

		err = k.IncrementTotalExternalYield(ctx, provider, identifier, yieldPortion)
		if err != nil {
			return errors.Wrapf(err, "unable to increment total yield for %s/%s", provider, identifier)
		}

		if transferErr == nil {
			k.logger.Info("claimed and transferred hyperlane yield", "amount", accruedYield, "identifier", identifier)
		}
	}

	return nil
}

// handleVaultsYieldSeasonOne handles the logic of the vaults for Season One.
// Yield from the Staked vault gets redirected to the Flexible vault.
func (k *Keeper) handleVaultsYieldSeasonOne(ctx context.Context, index int64) error {
	// Claim the yield of the Flexible vault.
	flexibleYield, err := k.claimModuleYield(ctx, vaults.FlexibleVaultAddress)
	if err != nil {
		return err
	}

	// Claim the yield of the Staked vault and redirect it to the Flexible vault.
	stakedYield, err := k.claimStakedVaultYield(ctx)
	if err != nil {
		return err
	}

	// get the current Flexible total principal.
	totalFlexiblePrincipal := math.ZeroInt()
	if has, _ := k.VaultsTotalFlexiblePrincipal.Has(ctx); has {
		current, err := k.VaultsTotalFlexiblePrincipal.Get(ctx)
		if err != nil {
			return err
		}
		totalFlexiblePrincipal = totalFlexiblePrincipal.Add(current)
	}

	// Register the new Rewards record.
	rewards := stakedYield.Add(flexibleYield)
	if err = k.VaultsRewards.Set(ctx, index, vaults.Reward{
		Index:   index,
		Total:   totalFlexiblePrincipal,
		Rewards: rewards,
	}); err != nil {
		return err
	}

	return nil
}

// handleVaultsYieldSeasonTwo handles the logic of the vaults for Season Two.
// Yield from the Staked vault gets redirected to a configured collector address.
func (k *Keeper) handleVaultsYieldSeasonTwo(ctx context.Context) error {
	// Claim the yield of the Staked vault.
	yield, err := k.claimModuleYield(ctx, vaults.StakedVaultAddress)
	if err != nil {
		return err
	}

	// Ensure that there is a valid amount of yield to send.
	if !yield.IsPositive() {
		return nil
	}

	// Send the Staked vault yield to the Collector address.
	err = k.bank.SendCoins(ctx, vaults.StakedVaultAddress, k.vaultsSeasonTwoYieldCollector, sdk.NewCoins(sdk.NewCoin(k.denom, yield)))
	if err != nil {
		return err
	}

	return nil
}
