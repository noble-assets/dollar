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

import "encoding/hex"

var (
	// NativeTokenTransferPrefix is the prefix for all NativeTokenTransfer payloads.
	//
	// https://github.com/wormhole-foundation/native-token-transfers/blob/67df54701e0f4b3793b6c621719911804c9875a3/evm/src/libraries/TransceiverStructs.sol#L37-L39
	NativeTokenTransferPrefix, _ = hex.DecodeString("994e5454")
	// TransceiverMessagePrefix is the prefix for all TransceiverMessage payloads.
	//
	// https://github.com/wormhole-foundation/native-token-transfers/blob/67df54701e0f4b3793b6c621719911804c9875a3/evm/src/Transceiver/WormholeTransceiver/WormholeTransceiverState.sol#L38-L41
	TransceiverMessagePrefix, _ = hex.DecodeString("9945ff10")
)
