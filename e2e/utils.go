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

	"cosmossdk.io/math"
	hyperlaneutil "github.com/bcp-innovations/hyperlane-cosmos/util"
	ismtypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	igptypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	hyperlanetypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/gogoproto/proto"
	wormholetypes "github.com/noble-assets/wormhole/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/relayer/rly"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"
	"go.uber.org/zap/zaptest"

	"dollar.noble.xyz/v2"
	portaltypes "dollar.noble.xyz/v2/types/portal"
	"dollar.noble.xyz/v2/utils"
)

// Suite is a utility for spinning up a new E2E testing suite.
func Suite(t *testing.T, ibcEnabled bool, hyperlaneEnabled bool) (ctx context.Context, noble *cosmos.CosmosChain, ibcSimapp *cosmos.CosmosChain, authority ibc.Wallet, guardians []utils.Guardian, tokenId string) {
	ctx = context.Background()
	logger := zaptest.NewLogger(t)
	reporter := testreporter.NewNopReporter()
	execReporter := reporter.RelayerExecReporter(t)
	client, network := interchaintest.DockerSetup(t)
	var relayer *rly.CosmosRelayer

	guardian := utils.NewGuardian(t)
	guardians = []utils.Guardian{guardian}

	numValidators, numFullNodes := 1, 0

	encodingConfig := testutil.MakeTestEncodingConfig(dollar.AppModule{}, warp.AppModule{})
	specs := []*interchaintest.ChainSpec{
		{
			Name:          "dollar",
			Version:       "local",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				Type:    "cosmos",
				Name:    "dollar",
				ChainID: "dollar-1",
				Images: []ibc.DockerImage{
					{
						Repository: "noble-dollar-simd",
						Version:    "local",
						UIDGID:     "1025:1025",
					},
				},
				Bin:            "simd",
				Bech32Prefix:   "noble",
				Denom:          "uusdc",
				GasPrices:      "0uusdc",
				GasAdjustment:  1.5,
				TrustingPeriod: "504h",
				ModifyGenesis: func(cc ibc.ChainConfig, genesis []byte) ([]byte, error) {
					peers := make(map[uint16]portaltypes.Peer)
					peers[uint16(vaautils.ChainIDEthereum)] = portaltypes.Peer{
						Transceiver: utils.SourceTransceiverAddress,
						Manager:     utils.SourceManagerAddress,
					}

					guardianSets := make(map[uint16]wormholetypes.GuardianSet)
					var addresses [][]byte
					for _, guardian := range guardians {
						addresses = append(addresses, guardian.Address.Bytes())
					}
					guardianSets[0] = wormholetypes.GuardianSet{
						Addresses:      addresses,
						ExpirationTime: 0,
					}

					updatedGenesis := []cosmos.GenesisKV{
						cosmos.NewGenesisKV("app_state.dollar.portal.peers", peers),
						cosmos.NewGenesisKV("app_state.wormhole.config", wormholetypes.Config{
							ChainId:          uint16(vaautils.ChainIDNoble),
							GuardianSetIndex: 0,
							GovChain:         uint16(vaautils.GovernanceChain),
							GovAddress:       vaautils.GovernanceEmitter.Bytes(),
						}),
						cosmos.NewGenesisKV("app_state.wormhole.guardian_sets", guardianSets),
					}

					return cosmos.ModifyGenesis(updatedGenesis)(cc, genesis)
				},
				EncodingConfig: &encodingConfig,
			},
		},
	}
	if ibcEnabled {
		specs = append(specs, &interchaintest.ChainSpec{
			Name:          "ibc-go-simd",
			Version:       "v8.7.0",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				ChainID: "ibc-go-simd-1",
			},
		})
	}
	factory := interchaintest.NewBuiltinChainFactory(logger, specs)

	chains, err := factory.Chains(t.Name())
	require.NoError(t, err)

	noble = chains[0].(*cosmos.CosmosChain)
	interchain := interchaintest.NewInterchain().AddChain(noble)
	if ibcEnabled {
		relayer = interchaintest.NewBuiltinRelayerFactory(
			ibc.CosmosRly,
			logger,
		).Build(t, client, network).(*rly.CosmosRelayer)

		ibcSimapp = chains[1].(*cosmos.CosmosChain)

		interchain = interchain.
			AddChain(ibcSimapp).
			AddRelayer(relayer, "relayer").
			AddLink(interchaintest.InterchainLink{
				Chain1:  noble,
				Chain2:  ibcSimapp,
				Relayer: relayer,
				Path:    "transfer",
			})
	}
	require.NoError(t, interchain.Build(ctx, execReporter, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
	}))

	t.Cleanup(func() {
		_ = interchain.Close()
	})

	if ibcEnabled {
		require.NoError(t, relayer.StartRelayer(ctx, execReporter))
	}

	// ARRANGE: Recover the authority wallet in order to perform gated actions.
	// This mnemonic must be for the authority account set in simapp/app.yaml!
	authority, err = interchaintest.GetAndFundTestUserWithMnemonic(ctx, "authority", "occur subway woman achieve deputy rapid museum point usual appear oil blue rate title claw debate flag gallery level object baby winner erase carbon", math.OneInt(), noble)
	require.NoError(t, err)

	if hyperlaneEnabled {
		validator := noble.Validators[0]

		_, err := validator.ExecTx(ctx, authority.KeyName(), "hyperlane", "ism", "create-noop")
		require.NoError(t, err)
		ismId := getHyperlaneIsmId(t, ctx, validator)

		_, err = validator.ExecTx(ctx, authority.KeyName(), "hyperlane", "hooks", "noop", "create")
		require.NoError(t, err)
		hookId := getHyperlaneHookId(t, ctx, validator)

		// TODO: Replace the Noble Hyperlane domain with the real value!
		_, err = validator.ExecTx(ctx, authority.KeyName(), "hyperlane", "mailbox", "create", ismId.String(), "1313817164")
		require.NoError(t, err)
		mailboxId := getHyperlaneMailboxId(t, ctx, validator)

		_, err = validator.ExecTx(ctx, authority.KeyName(), "hyperlane", "mailbox", "set", mailboxId.String(), "--required-hook", hookId.String(), "--default-hook", hookId.String())
		require.NoError(t, err)

		_, err = validator.ExecTx(ctx, authority.KeyName(), "hyperlane-transfer", "create-collateral-token", mailboxId.String(), "uusdn")
		require.NoError(t, err)
		tokenId = getHyperlaneTokenId(t, ctx, validator)

		_, err = validator.ExecTx(ctx, authority.KeyName(), "hyperlane-transfer", "enroll-remote-router", tokenId, "1", "0x0000000000000000000000000000000000000000000000000000000000000000", "0")
		require.NoError(t, err)
	}

	return
}

// getHyperlaneIsmId is a utility that returns the most recently creates ISM.
func getHyperlaneIsmId(t require.TestingT, ctx context.Context, validator *cosmos.ChainNode) hyperlaneutil.HexAddress {
	client := ismtypes.NewQueryClient(validator.GrpcConn)

	res, err := client.Isms(ctx, &ismtypes.QueryIsmsRequest{})
	require.NoError(t, err)
	require.Len(t, res.Isms, 1)
	require.Equal(t, "/hyperlane.core.interchain_security.v1.NoopISM", res.Isms[0].TypeUrl)

	var ism ismtypes.NoopISM
	err = proto.Unmarshal(res.Isms[0].Value, &ism)
	require.NoError(t, err)

	return ism.Id
}

// getHyperlaneHookId is a utility that returns the most recently created hook.
func getHyperlaneHookId(t require.TestingT, ctx context.Context, validator *cosmos.ChainNode) hyperlaneutil.HexAddress {
	client := igptypes.NewQueryClient(validator.GrpcConn)

	res, err := client.NoopHooks(ctx, &igptypes.QueryNoopHooksRequest{})
	require.NoError(t, err)
	require.Len(t, res.NoopHooks, 1)

	return res.NoopHooks[0].Id
}

// getHyperlaneMailboxId is a utility that returns the most recently created
func getHyperlaneMailboxId(t require.TestingT, ctx context.Context, validator *cosmos.ChainNode) hyperlaneutil.HexAddress {
	client := hyperlanetypes.NewQueryClient(validator.GrpcConn)

	res, err := client.Mailboxes(ctx, &hyperlanetypes.QueryMailboxesRequest{})
	require.NoError(t, err)
	require.Len(t, res.Mailboxes, 1)

	return res.Mailboxes[0].Id
}

// getHyperlaneTokenId is a utility that returns the most recently created warp
func getHyperlaneTokenId(t require.TestingT, ctx context.Context, validator *cosmos.ChainNode) string {
	client := warptypes.NewQueryClient(validator.GrpcConn)

	res, err := client.Tokens(ctx, &warptypes.QueryTokensRequest{})
	require.NoError(t, err)
	require.Len(t, res.Tokens, 1)

	return res.Tokens[0].Id
}
