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
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec/address"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"

	portaltypes "dollar.noble.xyz/v2/types/portal"
	"dollar.noble.xyz/v2/types/portal/ntt"
	dollartypes "dollar.noble.xyz/v2/types/v2"
	"dollar.noble.xyz/v2/utils"
)

var (
	// sourceToken is the 32-byte representation of the $M token on Ethereum Mainnet.
	// https://github.com/m0-foundation/m-portal/blob/dbe93da561c94dfc04beec8a144b11b287957b7a/deployments/1.json#L2
	sourceToken = common.FromHex("0x000000000000000000000000866a2bf4e572cbcf37d5071a7a58503bfb36be1b")
	// destinationToken is the 32-byte representation of the "uusdn" denom.
	destinationToken = common.FromHex("0x000000000000000000000000000000000000000000000000000000757573646e")

	// recipientManagerAddress is the 32-byte representation of the "dollar/manager" module account.
	recipientManagerAddress = common.FromHex("0x0000000000000000000000002e859506ba229c183f8985d54fe7210923fb9bca")

	// channelID is the IBC channel identifier of ibc-go-simd.
	channelID = "channel-0"

	// ibcYieldRecipient is the "gov" module account on ibc-go-simd.
	ibcYieldRecipient = "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"
)

/**
Add two yield recipients to two different chains
Do a yield update to first chain
Make assertions
Do a yield update to second chain
*/

// TestIBCYieldDistribution tests $USDN yield distribution across IBC channels.
func TestIBCYieldDistribution(t *testing.T) {
	ctx, _, chain, counterpartyChains, _, authority, guardians, _ := Suite(t, true, false, 2, 4)

	validator := chain.Validators[0]
	externalChain := counterpartyChains[0]
	externalChain2 := counterpartyChains[1]

	// ARRANGE: Create and fund test user accounts on both Noble and the external chain.
	wallets := interchaintest.GetAndFundTestUsers(t, ctx, "user", math.OneInt(), chain, externalChain, externalChain2)
	user, externalUser, externalUser2 := wallets[0], wallets[1], wallets[2]

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

	// ACT: Send 500,000 $USDN from the user on Noble to the external chain.
	_, err = chain.SendIBCTransfer(ctx, channelID, user.KeyName(), ibc.WalletAmount{
		Address: externalUser.FormattedAddress(),
		Denom:   "uusdn",
		Amount:  math.NewInt(500_000 * 1e6),
	}, ibc.TransferOptions{})

	// ASSERT: The transfer should've failed as a yield recipient hasn't been set.
	require.ErrorContains(t, err, fmt.Sprintf("ibc transfers of uusdn are currently disabled on %s", channelID))

	// ACT: Set the yield recipient for the external chains.
	_, err = validator.ExecTx(ctx, authority.KeyName(), "dollar", "set-yield-recipient", "IBC", channelID, ibcYieldRecipient)
	require.NoError(t, err)

	_, err = validator.ExecTx(ctx, authority.KeyName(), "dollar", "set-yield-recipient", "IBC", "channel-1", ibcYieldRecipient)
	require.NoError(t, err)

	// ASSERT: There are two yield recipient.
	yieldRecipients = getYieldRecipients(t, ctx, validator)
	key1 := fmt.Sprintf("%s/%s", dollartypes.Provider_IBC, channelID)
	key2 := fmt.Sprintf("%s/%s", dollartypes.Provider_IBC, "channel-1")
	require.Equal(t, ibcYieldRecipient, yieldRecipients[key1])
	require.Equal(t, ibcYieldRecipient, yieldRecipients[key2])

	// ACT: Send 500,000 $USDN from the user on Noble to the external chain.
	_, err = chain.SendIBCTransfer(ctx, channelID, user.KeyName(), ibc.WalletAmount{
		Address: externalUser.FormattedAddress(),
		Denom:   "uusdn",
		Amount:  math.NewInt(500_000 * 1e6),
	}, ibc.TransferOptions{})
	require.NoError(t, err)

	require.NoError(t, testutil.WaitForBlocks(ctx, 10, chain, externalChain))

	// ASSERT: The escrow account should now have 500,000 $USDN.
	expectedValue := math.NewInt(500_000 * 1e6)

	rawEscrowAddress := transfertypes.GetEscrowAddress(transfertypes.PortID, channelID)
	escrowAddress, err := address.NewBech32Codec(chain.Config().Bech32Prefix).BytesToString(rawEscrowAddress)
	require.NoError(t, err)

	balance, err = chain.BankQueryBalance(ctx, escrowAddress, "uusdn")
	require.NoError(t, err)
	require.True(t, expectedValue.Equal(balance), fmt.Sprintf("expected %s, got %s", expectedValue, balance))

	// ASSERT: The total supply should be 500,000 $USDN on the external chain.
	ibcDenom := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(transfertypes.PortID, channelID, "uusdn")).IBCDenom()
	totalSupply, err := externalChain.BankQueryTotalSupplyOf(ctx, ibcDenom)
	require.NoError(t, err)
	require.True(t, expectedValue.Equal(totalSupply.Amount), fmt.Sprintf("expected %s, got %s", expectedValue, totalSupply.Amount))

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
	_, err = validator.ExecTx(
		ctx, user.KeyName(),
		"dollar", "portal", "deliver", base64.StdEncoding.EncodeToString(bz),
		"--gas", "500000",
	)
	require.NoError(t, err)

	require.NoError(t, testutil.WaitForBlocks(ctx, 10, chain, externalChain))

	// ASSERT: The escrow account should now have 520,750 $USDN.
	balance, err = chain.BankQueryBalance(ctx, escrowAddress, "uusdn")
	require.NoError(t, err)
	require.True(t, math.NewInt(520_750*1e6).Equal(balance))

	// ASSERT: The total supply should be 520,750 $USDN on the external chain.
	totalSupply, err = externalChain.BankQueryTotalSupplyOf(ctx, ibcDenom)
	require.NoError(t, err)
	require.True(t, math.NewInt(520_750*1e6).Equal(totalSupply.Amount))

	// ----------- BEGIN YIELD UPDATE FOR SECOND YIELD RECIPIENT ---------------------------

	// ACT: Send 500,000 $USDN from the user on Noble to the external chain.
	_, err = chain.SendIBCTransfer(ctx, "channel-1", user.KeyName(), ibc.WalletAmount{
		Address: externalUser2.FormattedAddress(),
		Denom:   "uusdn",
		Amount:  math.NewInt(500_000 * 1e6),
	}, ibc.TransferOptions{})
	require.NoError(t, err)

	require.NoError(t, testutil.WaitForBlocks(ctx, 10, chain, externalChain2))

	// ASSERT: The escrow account should now have 500,000 $USDN.
	rawEscrowAddress2 := transfertypes.GetEscrowAddress(transfertypes.PortID, "channel-1")
	escrowAddress2, _ := address.NewBech32Codec(chain.Config().Bech32Prefix).BytesToString(rawEscrowAddress2)
	balance, err = chain.BankQueryBalance(ctx, escrowAddress2, "uusdn")
	require.NoError(t, err)
	require.True(t, expectedValue.Equal(balance), fmt.Sprintf("expected %s, got %s", expectedValue, balance))

	// ASSERT: The total supply should be 500,000 $USDN on the external chain.
	ibcDenom2 := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(transfertypes.PortID, "channel-0", "uusdn")).IBCDenom()
	totalSupply2, err := externalChain2.BankQueryTotalSupplyOf(ctx, ibcDenom2)
	require.NoError(t, err)
	require.True(t, expectedValue.Equal(totalSupply2.Amount), fmt.Sprintf("expected %s, got %s", expectedValue, totalSupply2.Amount))

	// ARRANGE: Prepare a VAA to be delivered that accrues a 4.15% yield.
	payload = portaltypes.EncodeIndexPayload(
		1041600000000,
		uint16(vaautils.ChainIDNoble),
	)

	transceiverMessage = buildTransceiverMessage(payload)
	vaa = utils.NewVAA(guardians, transceiverMessage)

	bz, err = vaa.Marshal()
	require.NoError(t, err)

	// ACT: Deliver the prepared VAA that accrues 4.15% yield.
	_, err = validator.ExecTx(
		ctx, user.KeyName(),
		"dollar", "portal", "deliver", base64.StdEncoding.EncodeToString(bz),
		"--gas", "500000",
	)
	require.NoError(t, err)

	require.NoError(t, testutil.WaitForBlocks(ctx, 10, chain, externalChain2))

	// ASSERT: The escrow account should now have 520,750 $USDN.
	expectedValue = math.NewInt(520_750 * 1e6)

	balance, err = chain.BankQueryBalance(ctx, escrowAddress2, "uusdn")
	require.NoError(t, err)
	require.True(t, expectedValue.Equal(balance), fmt.Sprintf("expected %s, got %s", expectedValue, balance))

	// ASSERT: The total supply should be 520,750 $USDN on the external chain.
	totalSupply, err = externalChain2.BankQueryTotalSupplyOf(ctx, ibcDenom2)
	require.NoError(t, err)
	require.True(t, expectedValue.Equal(totalSupply.Amount), fmt.Sprintf("expected %s, got %s", expectedValue, totalSupply.Amount))
}

// TestIBCYieldDistributionTimeout tests $USDN yield distribution across an IBC channel, while triggering a timeout.
func TestIBCYieldDistributionTimeout(t *testing.T) {
	ctx, execReporter, chain, counterpartyChains, relayer, authority, guardians, _ := Suite(t, true, false, 2, 4)
	validator := chain.Validators[0]
	externalChain := counterpartyChains[0]

	// ARRANGE: Create and fund test user accounts on both Noble and the external chain.
	wallets := interchaintest.GetAndFundTestUsers(t, ctx, "user", math.OneInt(), chain, externalChain)
	user, externalUser := wallets[0], wallets[1]

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

	// ACT: Send 500,000 $USDN from the user on Noble to the external chain.
	_, err = chain.SendIBCTransfer(ctx, channelID, user.KeyName(), ibc.WalletAmount{
		Address: externalUser.FormattedAddress(),
		Denom:   "uusdn",
		Amount:  math.NewInt(500_000 * 1e6),
	}, ibc.TransferOptions{})

	// ASSERT: The transfer should've failed as a yield recipient hasn't been set.
	require.ErrorContains(t, err, fmt.Sprintf("ibc transfers of uusdn are currently disabled on %s", channelID))

	// ACT: Set the yield recipient for the external chain.
	_, err = validator.ExecTx(ctx, authority.KeyName(), "dollar", "set-yield-recipient", "IBC", channelID, ibcYieldRecipient)
	require.NoError(t, err)

	// ASSERT: There is one yield recipient.
	yieldRecipients = getYieldRecipients(t, ctx, validator)
	key := fmt.Sprintf("%s/%s", dollartypes.Provider_IBC, channelID)
	require.Equal(t, ibcYieldRecipient, yieldRecipients[key])

	// ACT: Send 500,000 $USDN from the user on Noble to the external chain.
	_, err = chain.SendIBCTransfer(ctx, channelID, user.KeyName(), ibc.WalletAmount{
		Address: externalUser.FormattedAddress(),
		Denom:   "uusdn",
		Amount:  math.NewInt(500_000 * 1e6),
	}, ibc.TransferOptions{})
	require.NoError(t, err)

	require.NoError(t, testutil.WaitForBlocks(ctx, 10, chain, externalChain))

	// ASSERT: The escrow account should now have 500,000 $USDN.
	rawEscrowAddress := transfertypes.GetEscrowAddress(transfertypes.PortID, channelID)
	escrowAddress, _ := address.NewBech32Codec(chain.Config().Bech32Prefix).BytesToString(rawEscrowAddress)
	balance, err = chain.BankQueryBalance(ctx, escrowAddress, "uusdn")
	require.NoError(t, err)
	require.True(t, math.NewInt(500_000*1e6).Equal(balance))
	// ASSERT: The total supply should be 500,000 $USDN on the external chain.
	ibcDenom := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(transfertypes.PortID, channelID, "uusdn")).IBCDenom()
	totalSupply, err := externalChain.BankQueryTotalSupplyOf(ctx, ibcDenom)
	require.NoError(t, err)
	require.True(t, math.NewInt(500_000*1e6).Equal(totalSupply.Amount))

	// ARRANGE: Stop the relayer before accruing yield.
	require.NoError(t, relayer.StopRelayer(ctx, execReporter))

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
	_, err = validator.ExecTx(
		ctx, user.KeyName(),
		"dollar", "portal", "deliver", base64.StdEncoding.EncodeToString(bz),
		"--gas", "500000",
	)
	require.NoError(t, err)

	// ARRANGE: Start the relayer after 10 minutes to trigger a timeout.
	time.Sleep(10 * time.Minute)
	require.NoError(t, relayer.StartRelayer(ctx, execReporter))
	require.NoError(t, testutil.WaitForBlocks(ctx, 10, chain, externalChain))

	// ASSERT: The retry amount should be 20,750 $USDN.
	retryAmounts := getRetryAmounts(t, ctx, validator)
	require.Equal(t, "20750000000", retryAmounts[key])
	// ASSERT: The escrow account should now have 520,750 $USDN.
	balance, err = chain.BankQueryBalance(ctx, escrowAddress, "uusdn")
	require.NoError(t, err)
	require.True(t, math.NewInt(520_750*1e6).Equal(balance))
	// ASSERT: The total supply should still be 500,000 $USDN on the external chain.
	totalSupply, err = externalChain.BankQueryTotalSupplyOf(ctx, ibcDenom)
	require.NoError(t, err)
	require.True(t, math.NewInt(500_000*1e6).Equal(totalSupply.Amount))

	// ARRANGE: Prepare a VAA to be delivered that accrues another 4.15% yield.
	payload = portaltypes.EncodeIndexPayload(
		1083000000000,
		uint16(vaautils.ChainIDNoble),
	)

	transceiverMessage = buildTransceiverMessage(payload)
	vaa = utils.NewVAA(guardians, transceiverMessage)

	bz, err = vaa.Marshal()
	require.NoError(t, err)

	// ACT: Deliver the prepared VAA that accrues another 4.15% yield.
	_, err = validator.ExecTx(
		ctx, user.KeyName(),
		"dollar", "portal", "deliver", base64.StdEncoding.EncodeToString(bz),
		"--gas", "500000",
	)
	require.NoError(t, err)

	require.NoError(t, testutil.WaitForBlocks(ctx, 10, chain, externalChain))

	// ASSERT: The retry amounts should be empty.
	retryAmounts = getRetryAmounts(t, ctx, validator)
	require.Empty(t, retryAmounts)
	// ASSERT: The escrow account should now have 541,500 $USDN.
	balance, err = chain.BankQueryBalance(ctx, escrowAddress, "uusdn")
	require.NoError(t, err)
	require.True(t, math.NewInt(541_500*1e6).Equal(balance))
	// ASSERT: The total supply should be 541,500 $USDN on the external chain.
	totalSupply, err = externalChain.BankQueryTotalSupplyOf(ctx, ibcDenom)
	require.NoError(t, err)
	require.True(t, math.NewInt(541_500*1e6).Equal(totalSupply.Amount))
}

// buildTransceiverMessage is a utility that builds a transceiver message.
func buildTransceiverMessage(payload []byte) []byte {
	managerMessage := ntt.ManagerMessage{
		Id:      make([]byte, 32),
		Sender:  make([]byte, 32),
		Payload: payload,
	}

	transceiverMessage := ntt.TransceiverMessage{
		SourceManagerAddress:    utils.SourceManagerAddress,
		RecipientManagerAddress: recipientManagerAddress,
		ManagerPayload:          ntt.EncodeManagerMessage(managerMessage),
		TransceiverPayload:      nil,
	}

	return ntt.EncodeTransceiverMessage(transceiverMessage)
}

// getYieldRecipients is a utility that queries the yield recipients.
func getYieldRecipients(t require.TestingT, ctx context.Context, validator *cosmos.ChainNode) map[string]string {
	client := dollartypes.NewQueryClient(validator.GrpcConn)

	res, err := client.YieldRecipients(ctx, &dollartypes.QueryYieldRecipients{})
	require.NoError(t, err)

	return res.YieldRecipients
}

// getRetryAmounts is a utility that queries the retry amounts.
func getRetryAmounts(t require.TestingT, ctx context.Context, validator *cosmos.ChainNode) map[string]string {
	client := dollartypes.NewQueryClient(validator.GrpcConn)

	res, err := client.RetryAmounts(ctx, &dollartypes.QueryRetryAmounts{})
	require.NoError(t, err)

	return res.RetryAmounts
}
