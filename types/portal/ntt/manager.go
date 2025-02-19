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

package ntt

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

// EncodeManagerMessage is a utility that encodes a manager message.
func EncodeManagerMessage(msg ManagerMessage) (bz []byte) {
	bz = append(bz, msg.Id...)
	bz = append(bz, msg.Sender...)

	bz = binary.BigEndian.AppendUint16(bz, uint16(len(msg.Payload)))
	bz = append(bz, msg.Payload...)

	return
}

// ParseManagerMessage is a utility that parses a manager message.
func ParseManagerMessage(bz []byte) (msg ManagerMessage, err error) {
	if len(bz) < 66 {
		return ManagerMessage{}, errors.New("manager message is malformed")
	}

	msg.Id = bz[:32]
	msg.Sender = bz[32:64]

	payloadLength := binary.BigEndian.Uint16(bz[64:66])
	msg.Payload = bz[66:]
	if len(msg.Payload) != int(payloadLength) {
		return ManagerMessage{}, errors.New("manager message payload is invalid length")
	}

	return
}

func ManagerMessageDigest(sourceChainId uint16, msg ManagerMessage) []byte {
	rawSourceChainId := make([]byte, 2)
	binary.BigEndian.PutUint16(rawSourceChainId, sourceChainId)

	return crypto.Keccak256(append(rawSourceChainId, EncodeManagerMessage(msg)...))
}
