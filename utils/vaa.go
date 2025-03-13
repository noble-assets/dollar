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
	"crypto/ecdsa"
	"crypto/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"
)

var sequence = uint64(0)

// Guardian defines the structure of a guardian for unit testing.
type Guardian struct {
	Key     *ecdsa.PrivateKey
	Address common.Address
}

// NewGuardian creates a new Guardian for unit testing.
func NewGuardian() Guardian {
	ethcrypto.S256()
	key, _ := ecdsa.GenerateKey(ethcrypto.S256(), rand.Reader)

	return Guardian{
		Key:     key,
		Address: ethcrypto.PubkeyToAddress(key.PublicKey),
	}
}

// NewVAA creates and signs a new VAA based on a specified payload and guardian set.
func NewVAA(guardians []Guardian, payload []byte) vaautils.VAA {
	vaa := vaautils.VAA{
		Version:          vaautils.SupportedVAAVersion,
		GuardianSetIndex: 0,
		Timestamp:        time.Now().Local().Truncate(time.Second),
		Sequence:         sequence,
		Payload:          payload,

		// For testing purposes, we assume that VAAs are emitted by the
		// Wormhole Transceiver contract for the Noble Portal deployed on
		// Ethereum Mainnet.
		//
		// https://github.com/m0-foundation/m-portal/blob/dbe93da561c94dfc04beec8a144b11b287957b7a/deployments/noble/1.json#L3
		EmitterChain:   vaautils.ChainIDEthereum,
		EmitterAddress: vaautils.Address(common.FromHex("0x000000000000000000000000c7dd372c39e38bf11451ab4a8427b4ae38cef644")),
	}
	sequence += 1

	var addresses []common.Address
	for idx, guardian := range guardians {
		vaa.AddSignature(guardian.Key, uint8(idx))
		addresses = append(addresses, guardian.Address)
	}

	return vaa
}
