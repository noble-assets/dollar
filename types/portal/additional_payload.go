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
	"errors"
)

// EncodeAdditionalPayload is a utility for encoding an M Portal additional payload.
//
// https://github.com/m0-foundation/m-portal/blob/9c72c13d8416ef77f3ea89316167697133d1eeee/src/libs/PayloadEncoder.sol#L68-L73
func EncodeAdditionalPayload(index int64, destinationToken []byte) (bz []byte) {
	bz = binary.BigEndian.AppendUint64(bz, uint64(index))
	bz = append(bz, destinationToken...)

	return
}

// DecodeAdditionalPayload is a utility for decoding an M Portal additional payload.
//
// https://github.com/m0-foundation/m-portal/blob/9c72c13d8416ef77f3ea89316167697133d1eeee/src/libs/PayloadEncoder.sol#L75-L86
func DecodeAdditionalPayload(bz []byte) (index int64, destinationToken []byte, err error) {
	indexLen := 8
	destinationTokenLen := 32

	if len(bz) < indexLen+destinationTokenLen {
		return 0, []byte{}, errors.New("native token transfer additional payload is invalid length")
	}

	offset := 0
	index = int64(binary.BigEndian.Uint64(bz[offset : offset+indexLen]))

	offset += indexLen
	destinationToken = bz[offset : offset+destinationTokenLen]

	return index, destinationToken, nil
}
