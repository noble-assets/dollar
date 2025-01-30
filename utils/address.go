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

package utils

import (
	"github.com/cometbft/cometbft/crypto/secp256k1"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmos "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Account struct {
	Key     secp256k1.PrivKey
	PubKey  *codectypes.Any
	Address string
	Invalid string
	Bytes   []byte
}

func TestAccount() Account {
	key := secp256k1.GenPrivKey()
	pubKey, _ := codectypes.NewAnyWithValue(&cosmos.PubKey{
		Key: key.PubKey().Bytes(),
	})
	bytes := key.PubKey().Address().Bytes()
	address, _ := sdk.Bech32ifyAddressBytes("noble", bytes)
	invalid, _ := sdk.Bech32ifyAddressBytes("cosmos", bytes)

	return Account{
		Key:     key,
		PubKey:  pubKey,
		Address: address,
		Invalid: invalid,
		Bytes:   bytes,
	}
}
