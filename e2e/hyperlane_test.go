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
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	hyperlaneutil "github.com/bcp-innovations/hyperlane-cosmos/util"
	hyperlanetypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	portaltypes "dollar.noble.xyz/v2/types/portal"
	"dollar.noble.xyz/v2/types/portal/ntt"
	dollartypes "dollar.noble.xyz/v2/types/v2"
	"dollar.noble.xyz/v2/utils"
)

// hyperlaneYieldRecipient is the empty address for testing purposes.
var hyperlaneYieldRecipient = "0x0000000000000000000000000000000000000000000000000000000000000000"

// TestHyperlaneYieldDistribution tests $USDN yield distribution across Hyperlane routes.
func TestHyperlaneYieldDistribution(t *testing.T) {
	ctx, _, chain, _, _, authority, guardians, tokenId := Suite(t, false, true)
	validator := chain.Validators[0]

	// ARRANGE: Create and fund a test user account on Noble.
	wallets := interchaintest.GetAndFundTestUsers(t, ctx, "user", math.OneInt(), chain)
	user := wallets[0]

	// ARRANGE: Pad the user address to be compatible with Wormhole's NTT standard.
	paddedAddress := make([]byte, 32)
	copy(paddedAddress[12:], user.Address())

	// ARRANGE: Prepare a VAA to be delivered that mints the user 1,000,000 $USDN.
	additionalPayload := portaltypes.EncodeAdditionalPayload(1e12, destinationToken)
	payload := ntt.EncodeNativeTokenTransfer(ntt.NativeTokenTransfer{
		Amount:            1_000_000 * 1e6,
		SourceToken:       sourceToken,
		To:                paddedAddress,
		ToChain:           uint16(vaautils.ChainIDNoble),
		AdditionalPayload: additionalPayload,
	})

	transceiverMessage := buildTransceiverMessage(payload)
	vaa := utils.NewVAA(guardians, transceiverMessage)

	bz, err := vaa.Marshal()
	require.NoError(t, err)

	// ACT: Deliver the prepared VAA that mints 1,000,000 $USDN.
	_, err = validator.ExecTx(
		ctx, authority.KeyName(),
		"dollar", "portal", "deliver", base64.StdEncoding.EncodeToString(bz),
	)
	require.NoError(t, err)

	// ASSERT: The user should now have 1,000,000 $USDN.
	balance, err := chain.BankQueryBalance(ctx, user.FormattedAddress(), "uusdn")
	require.NoError(t, err)
	require.True(t, math.NewInt(1_000_000*1e6).Equal(balance))

	// ACT: Query all yield recipients.
	yieldRecipients := getYieldRecipients(t, ctx, validator)

	// ASSERT: There are no yield recipients.
	require.Empty(t, yieldRecipients)

	// ACT: Set the yield recipient for the external chain.
	_, err = validator.ExecTx(ctx, authority.KeyName(), "dollar", "set-yield-recipient", "HYPERLANE", tokenId, hyperlaneYieldRecipient)
	require.NoError(t, err)

	// ASSERT: There is one yield recipient.
	yieldRecipients = getYieldRecipients(t, ctx, validator)
	key := fmt.Sprintf("%s/%s", dollartypes.Provider_HYPERLANE, tokenId)
	require.Equal(t, hyperlaneYieldRecipient, yieldRecipients[key])

	// ACT: Send 500,000 $USDN from the user on Noble to the external chain.
	hash, err := validator.ExecTx(ctx, user.KeyName(), "warp", "transfer", tokenId, "1", "0x"+common.Bytes2Hex(paddedAddress), math.NewInt(500_000*1e6).String(), "--max-hyperlane-fee", "0uusdn")
	require.NoError(t, err)

	// ASSERT: The escrow account should now have 500,000 $USDN.
	rawEscrowAddress := authtypes.NewModuleAddress(warptypes.ModuleName).Bytes()
	escrowAddress, _ := address.NewBech32Codec(chain.Config().Bech32Prefix).BytesToString(rawEscrowAddress)
	balance, err = chain.BankQueryBalance(ctx, escrowAddress, "uusdn")
	require.NoError(t, err)
	require.True(t, math.NewInt(500_000*1e6).Equal(balance))
	// ASSERT: A valid Hyperlane Warp message should have been dispatched.
	warpPayload := getWarpPayload(t, validator, hash)
	require.Equal(t, big.NewInt(500_000*1e6), warpPayload.Amount())

	// ARRANGE: Prepare a VAA to be delivered that accrues a 4.15% yield.
	payload = portaltypes.EncodeIndexPayload(
		1041500000000,
		uint16(vaautils.ChainIDNoble),
	)

	transceiverMessage = buildTransceiverMessage(payload)
	vaa = utils.NewVAA(guardians, transceiverMessage)

	bz, err = vaa.Marshal()
	require.NoError(t, err)

	// ACT: Deliver the prepared VAA that accrues 4.15% yield.
	hash, err = validator.ExecTx(
		ctx, user.KeyName(),
		"dollar", "portal", "deliver", base64.StdEncoding.EncodeToString(bz),
		"--gas", "500000",
	)
	require.NoError(t, err)

	// ASSERT: The escrow account should now have 520,750 $USDN.
	balance, err = chain.BankQueryBalance(ctx, escrowAddress, "uusdn")
	require.NoError(t, err)
	require.True(t, math.NewInt(520_750*1e6).Equal(balance))
	// ASSERT: A valid Hyperlane Warp message should have been dispatched.
	warpPayload = getWarpPayload(t, validator, hash)
	require.Equal(t, big.NewInt(20_750*1e6), warpPayload.Amount())
}

func getWarpPayload(t require.TestingT, validator *cosmos.ChainNode, hash string) warptypes.WarpPayload {
	tx, err := validator.GetTransaction(validator.CliContext(), hash)
	require.NoError(t, err)

	for _, rawEvent := range tx.Events {
		event, err := sdk.ParseTypedEvent(rawEvent)
		if err == nil {
			dispatch, ok := event.(*hyperlanetypes.EventDispatch)
			if ok {
				message, err := hyperlaneutil.ParseHyperlaneMessage(common.FromHex(dispatch.Message))
				require.NoError(t, err)
				payload, err := warptypes.ParseWarpPayload(message.Body)
				require.NoError(t, err)

				return payload
			}
		}
	}

	return warptypes.WarpPayload{}
}
