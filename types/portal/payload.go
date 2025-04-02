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

	"github.com/ethereum/go-ethereum/common"

	"cosmossdk.io/math"

	"dollar.noble.xyz/v2/types/portal/ntt"
)

type PayloadType int

const (
	Unknown PayloadType = iota
	Token
	Index
)

// M0IT - M0 Index Transfer
const M0IT = "4d304954"

// GetPayloadType is a utility for determining the type of custom payload.
// Since we aren't implementing a registrar, the Key and List types are ignored.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L32-L43
func GetPayloadType(payload []byte) PayloadType {
	if len(payload) < 4 {
		return Unknown
	}

	switch hex.EncodeToString(payload[:4]) {
	case "994e5454":
		return Token
	case M0IT:
		return Index
	}

	return Unknown
}

// TokenPayload is a data structure that holds the fields decoded from
// a token transfer payload.
type TokenPayload struct {
	Amount             math.Int
	Index              int64
	Recipient          []byte
	DestinationChainId uint16
	DestinationToken   []byte
}

// DecodeTokenPayload is a utility for decoding a custom payload of type Token.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L45-L62
func DecodeTokenPayload(payload []byte) TokenPayload {
	ntt, _ := ntt.ParseNativeTokenTransfer(payload)

	amount := math.NewIntFromUint64(ntt.Amount)

	// NOTE: the error here is ignored like in the parsing of ntt.
	index, destinationToken, _ := DecodeAdditionalPayload(ntt.AdditionalPayload)

	return TokenPayload{amount, index, ntt.To[12:], ntt.ToChain, destinationToken}
}

// DecodeIndexPayload is a utility for decoding a custom payload of type Index.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L68-L75
func DecodeIndexPayload(payload []byte) (index int64, destination uint16) {
	index = int64(binary.BigEndian.Uint64(payload[4:12]))
	destination = binary.BigEndian.Uint16(payload[12:14])

	return
}

// EncodeIndexPayload is a utility for encoding a custom payload of type Index.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L64-L66
func EncodeIndexPayload(index int64, destination uint16) (bz []byte) {
	bz = append(bz, common.FromHex(M0IT)...)

	indexBz := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBz, uint64(index))
	bz = append(bz, indexBz...)

	destinationBz := make([]byte, 2)
	binary.BigEndian.PutUint16(destinationBz, destination)
	bz = append(bz, destinationBz...)

	return
}

// EventsPayload is a data structure used to hold information required to emit complete
// events during the handling of a vaa.
type EventsPayload struct {
	SourceChainId uint32
	Sender        []byte
	MessageId     []byte
}
