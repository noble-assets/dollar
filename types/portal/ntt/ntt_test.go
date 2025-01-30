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
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	nttpkg "dollar.noble.xyz/types/portal/ntt"
)

func TestNativeTokenTransferUtilities(t *testing.T) {
	raw, err := hex.DecodeString("994e54540600000000000f42400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000153b01572d58e7d24a4611b2bb31d9e8de9535a70fa900080000000000000000")
	require.NoError(t, err, "unable to decode input")

	ntt, err := nttpkg.ParseNativeTokenTransfer(raw)
	require.NoError(t, err, "unable to decode native token transfer")

	require.Equal(t, 1_000_000, int(ntt.Amount))
	require.Equal(t, make([]byte, 32), ntt.SourceToken)
	require.Equal(t, "noble1z5asz4edtrnayjjxzxetkvwear0f2dd8amwkhd", MustEncodeRecipient(ntt.To))
	require.Equal(t, 4009, int(ntt.ToChain))
	require.Equal(t, 0, int(binary.BigEndian.Uint64(ntt.AdditionalPayload)))

	bz := nttpkg.EncodeNativeTokenTransfer(ntt)
	require.Equal(t, raw, bz)
}
