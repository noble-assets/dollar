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

package ntt_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	"dollar.noble.xyz/v3/types/portal/ntt"
)

func TestManagerMessageUtilities(t *testing.T) {
	raw, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000090f8bf6a479f320ead074411a4b0e7944ea8c9c10059994e54540600000000000f42400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000153b01572d58e7d24a4611b2bb31d9e8de9535a70fa900080000000000000000")
	require.NoError(t, err, "unable to decode input")

	msg, err := ntt.ParseManagerMessage(raw)
	require.NoError(t, err, "unable to parse manager message")

	require.Equal(t, make([]byte, 32), msg.Id)
	require.Equal(t, "90f8bf6a479f320ead074411a4b0e7944ea8c9c1", hex.EncodeToString(msg.Sender[12:]))
	require.Equal(t, "994e54540600000000000f42400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000153b01572d58e7d24a4611b2bb31d9e8de9535a70fa900080000000000000000", hex.EncodeToString(msg.Payload))

	bz := ntt.EncodeManagerMessage(msg)
	require.Equal(t, raw, bz)
}

func TestManagerMessageDigest(t *testing.T) {
	id, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000002")
	sender, _ := hex.DecodeString("00000000000000000000000012b1a4226ba7d9ad492779c924b0fc00bdcb6217")
	payload, _ := hex.DecodeString("994e54540600000000000f423f000000000000000000000000866a2bf4e572cbcf37d5071a7a58503bfb36be1b00000000000000000000000012b1a4226ba7d9ad492779c924b0fc00bdcb621727130028000000e8d58b8383000000000000000000000000866a2bf4e572cbcf37d5071a7a58503bfb36be1b")

	msg := ntt.ManagerMessage{
		Id:      id,
		Sender:  sender,
		Payload: payload,
	}

	digest := ntt.ManagerMessageDigest(uint16(vaautils.ChainIDSepolia), msg)
	// https://sepolia.etherscan.io/tx/0xee004a1d65c6d1bfb12e83347432a17a0a2da5031cc81e5f53526eacc43fa9a1#eventlog#395
	require.Equal(t, "E98541A7019A1D0092A127981AD6C47C30C72C80965B11EC7072886649A9E791", strings.ToUpper(hex.EncodeToString(digest)), "expected a different digest")
}
