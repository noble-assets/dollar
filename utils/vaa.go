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
	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"
)

var sequence = uint64(0)

// Guardian defines the structure of a guardian for testing.
type Guardian struct {
	Key     *ecdsa.PrivateKey
	Address common.Address
}

// NewGuardian creates a new Guardian for testing.
func NewGuardian(t require.TestingT) Guardian {
	key, err := ecdsa.GenerateKey(ethcrypto.S256(), rand.Reader)
	require.NoError(t, err)

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
		EmitterChain:     vaautils.ChainIDEthereum,
		EmitterAddress:   vaautils.Address(SourceTransceiverAddress),
		Payload:          payload,
	}
	sequence += 1

	for idx, guardian := range guardians {
		vaa.AddSignature(guardian.Key, uint8(idx))
	}

	return vaa
}
