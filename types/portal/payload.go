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

package portal

import (
	"encoding/binary"
	"encoding/hex"

	"cosmossdk.io/math"

	"dollar.noble.xyz/types/portal/ntt"
)

type PayloadType int

const (
	Unknown PayloadType = iota
	Token
	Index
)

// GetPayloadType is a utility for determining the type of custom payload.
// Since we aren't implementing a registrar, the Key and List types are ignored.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L32-L43
func GetPayloadType(payload []byte) PayloadType {
	if len(payload) < 4 {
		return Unknown
	}

	switch hex.EncodeToString(payload[:4]) {
	// 0x99'N''T''T'
	case "994e5454":
		return Token
	// M0IT - M0 Index Transfer
	case "4d304954":
		return Index
	}

	return Unknown
}

// DecodeTokenPayload is a utility for decoding a custom payload of type Token.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L45-L62
func DecodeTokenPayload(payload []byte) (amount math.Int, index int64, recipient []byte, destination uint16) {
	ntt, _ := ntt.ParseNativeTokenTransfer(payload)

	amount = math.NewIntFromUint64(ntt.Amount)
	index = int64(binary.BigEndian.Uint64(ntt.AdditionalPayload))

	return amount, index, ntt.To[12:], ntt.ToChain
}

// DecodeIndexPayload is a utility for decoding a custom payload of type Index.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L68-L75
func DecodeIndexPayload(payload []byte) (index int64, destination uint16) {
	index = int64(binary.BigEndian.Uint64(payload[4:12]))
	destination = binary.BigEndian.Uint16(payload[12:14])

	return
}
