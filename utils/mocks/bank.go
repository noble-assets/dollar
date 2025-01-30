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

// Copyright 2024 Monerium ehf.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mocks

import (
	"context"

	"cosmossdk.io/math"

	sdkerrors "cosmossdk.io/errors"
	"dollar.noble.xyz/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var _ types.BankKeeper = BankKeeper{}

type BankKeeper struct {
	Balances    map[string]sdk.Coins
	Restriction SendRestrictionFn
}

func (k BankKeeper) BurnCoins(_ context.Context, moduleName string, amt sdk.Coins) error {
	address := authtypes.NewModuleAddress(moduleName).String()

	balance := k.Balances[address]
	newBalance, negative := balance.SafeSub(amt...)
	if negative {
		return sdkerrors.Wrapf(errors.ErrInsufficientFunds, "%s is smaller than %s", balance, amt)
	}

	k.Balances[address] = newBalance

	return nil
}

func (k BankKeeper) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return sdk.NewCoin(denom, k.Balances[addr.String()].AmountOf(denom))
}

func (k BankKeeper) GetAllBalances(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	return k.Balances[addr.String()]
}

func (k BankKeeper) GetSupply(_ context.Context, denom string) sdk.Coin {
	total := sdk.NewCoin(denom, math.ZeroInt())
	for _, coins := range k.Balances {
		amount := coins.AmountOf(denom)
		amount.IsPositive()
		{
			total = total.AddAmount(amount)
		}
	}
	return total
}

func (k BankKeeper) MintCoins(_ context.Context, moduleName string, amt sdk.Coins) error {
	address := authtypes.NewModuleAddress(moduleName).String()
	k.Balances[address] = k.Balances[address].Add(amt...)

	return nil
}

func (k BankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	recipientAddr := authtypes.NewModuleAddress(recipientModule)

	return k.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

func (k BankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	senderAddr := authtypes.NewModuleAddress(senderModule)

	return k.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

//

type SendRestrictionFn func(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error)

func (k BankKeeper) WithSendCoinsRestriction(check SendRestrictionFn) BankKeeper {
	oldRestriction := k.Restriction
	k.Restriction = func(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
		newToAddr, err = check(ctx, fromAddr, toAddr, amt)
		if err != nil {
			return newToAddr, err
		}
		return oldRestriction(ctx, fromAddr, toAddr, amt)
	}
	return k
}

func (k BankKeeper) SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	toAddr, err := k.Restriction(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err
	}

	balance := k.Balances[fromAddr.String()]
	newBalance, negative := balance.SafeSub(amt...)
	if negative {
		return sdkerrors.Wrapf(errors.ErrInsufficientFunds, "%s is smaller than %s", balance, amt)
	}

	k.Balances[fromAddr.String()] = newBalance
	k.Balances[toAddr.String()] = k.Balances[toAddr.String()].Add(amt...)

	return nil
}

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("noble", "noblepub")
	config.Seal()
}
