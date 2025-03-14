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

package e2e

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	portaltypes "dollar.noble.xyz/types/portal"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/stretchr/testify/require"
)

// TestMsgDeliverInjection ensures that the injected message can't be executed publicly.
func TestMsgDeliverInjection(t *testing.T) {
	ctx, chain, _, _ := Suite(t, false)

	broadcaster := cosmos.NewBroadcaster(t, chain)
	user := interchaintest.GetAndFundTestUsers(t, ctx, "wallet", math.OneInt(), chain)[0]

	bCtx, bCancel := context.WithTimeout(ctx, 20*time.Second)
	defer bCancel()

	mockVAA := "AQAAAAABAPel1AcBA57rIzaTw70Qqlta9SxhuBYByiTv3viGqwgfFq4Wfx/EN0Mb8D71aTIwBz36NUmI98Q2fCEQyFlFSqQAZ1vRXAAAAAAnEgAAAAAAAAAAAAAAAHsb16a05hwqEjrGvCy/xhRDfQRwAAAAAAAAsrwPAScUAAAAAAAAAAAAAAAAKcvx4HFm0xRGMHrgeZn6bRYiOZAAAADjmUX/EAAAAAAAAAAAAAAAABt64ZSyDFVbnZmcg190zc42pnp0AAAAAAAAAAAAAAAAG3rhlLIMVVudmZyDX3TNzjamenQAmwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABrAAAAAAAAAAAAAAAAlO0ORGvBexrFFbj2bBk6ZU0driQAWZlOVFQGAAAAAAAAJw8AAAAAAAAAAAAAAAAMlBrZTKSlLtrqvyA7Yb3RgHzuwAAAAAAAAAAAAAAAAJTtDkRrwXsaxRW49mwZOmVNHa4kJxQACAAAAO++YLGeAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGpv1m2ycUAAAAAAAAAAAAAAAAlO0ORGvBexrFFbj2bBk6ZU0driQAAAAAAAAAAAAAAAB6ClOEd3b36UzDV0KXGssiF7DbgQAAAAAAAAAAAAAAAHoKU4R3dvfpTMNXQpcayyIXsNuBAAAAAAAAAAAAAAAAKcvx4HFm0xRGMHrgeZn6bRYiOZAA"

	_, err := cosmos.BroadcastTx(
		bCtx,
		broadcaster,
		user,
		&portaltypes.MsgDeliverInjection{
			Vaa: []byte(mockVAA),
		},
	)

	require.Error(t, err)
	require.ErrorContains(t, err, "no message handler found")

	// TODO: Query for the VAA and ensure it was not processed
}
